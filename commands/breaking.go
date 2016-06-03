package commands

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/michigan-com/gannett-newsfetch/config"
	api "github.com/michigan-com/gannett-newsfetch/gannettApi"
	m "github.com/michigan-com/gannett-newsfetch/model"
	"github.com/michigan-com/newsfetch/lib"
)

var breakingCommand = &cobra.Command{
	Use:   "breaking-news",
	Short: "Check the Gannett API for breaking news",
	Run:   breakingCmdRun,
}

func breakingCmdRun(command *cobra.Command, argv []string) {
	apiConfig, _ := config.GetApiConfig()
	envConfig, _ := config.GetEnv()

	if envConfig.MongoUri == "" {
		log.Print("No mongo uri specified, no articles will be saved")
		return
	} else if len(apiConfig.SiteCodes) == 0 {
		log.Fatalf("No site codes input, please set the SITE_CODES env variable")
		return
	}

	FetchBreakingNews(envConfig.MongoUri, apiConfig.SiteCodes, envConfig.GnapiDomain)
}

func FetchBreakingNews(mongoUri string, siteCodes []string, gnapiDomain string) {
	session := lib.DBConnect(mongoUri)
	defer session.Close()

	for {
		var breakingWait sync.WaitGroup
		breakingChannel := make(chan *m.SearchArticle, len(siteCodes)*100)
		log.Info("Fetching breaking news...")
		for _, siteCode := range siteCodes {
			breakingWait.Add(1)
			go func(siteCode string) {
				defer breakingWait.Done()
				articles := api.GetBreakingNews(siteCode)

				for _, article := range articles {
					breakingChannel <- article
				}
			}(siteCode)
		}
		breakingWait.Wait()
		close(breakingChannel)
		log.Info("...Done fetching articles")

		SaveArticles(breakingChannel, session)

		if gnapiDomain != "" {
			gnapiUrl := fmt.Sprintf("%s/%s/", gnapiDomain, "breaking-news")
			resp, err := http.Get(gnapiUrl)
			if err == nil {
				resp.Body.Close()
			}
		}

		if loop > 0 {
			log.Infof("Sleeping for %d seconds...", loop)
			time.Sleep(time.Duration(loop) * time.Second)
			log.Info("...and now I'm awake!")
		} else {
			break
		}
	}
}

func SaveArticles(breakingChannel chan *m.SearchArticle, session *mgo.Session) {
	log.Info("Collecting breaking articles to save...")
	breakingNewsSnapshot := m.BreakingNewsSnapshot{}
	breakingArticles := make([]*m.BreakingNewsArticle, 0, 100)
	toScrape := make([]interface{}, 0, 100)
	for article := range breakingChannel {
		articleId := lib.GetArticleId(article.Urls.LongUrl)
		if articleId == -1 {
			log.Warningf(`Failed to get id for url %s`, article.Urls.LongUrl)
			continue
		}

		breakingArticle := &m.BreakingNewsArticle{
			ArticleId:   articleId,
			Headline:    article.Headline,
			Subheadline: article.PromoBrief,
		}

		breakingArticles = append(breakingArticles, breakingArticle)

		if shouldSummarizeArticle(article, session) {
			articleIdQuery := bson.M{"article_id": article.AssetId}
			toScrape = append(toScrape, articleIdQuery)
			toScrape = append(toScrape, articleIdQuery)
		}
	}

	if len(toScrape) > 0 {
		bulk := session.DB("").C("ToScrape").Bulk()
		bulk.Upsert(toScrape...)
		_, err := bulk.Run()
		if err != nil {
			log.Errorf("Failed to store articles to be scraped: %v", err)
		}
	}

	breakingNewsSnapshot.Articles = breakingArticles
	breakingCol := session.DB("").C("BreakingNews")
	err := breakingCol.Insert(breakingNewsSnapshot)

	if err != nil {
		log.Warningf(`

			Failed to save breaking news snapshot:

				Err: %v

		`, err)
		return
	}

	RemoveOldBreakingSnapshot(breakingCol)

	log.Infof("...Done saving breaking articles, count: %d", len(breakingArticles))
}

func RemoveOldBreakingSnapshot(col *mgo.Collection) {

	var snapshot = bson.M{
		"_id": -1,
	}
	// Remove old snapshots
	col.Find(bson.M{}).
		Select(bson.M{"_id": 1}).
		Sort("-_id").
		One(&snapshot)

	_, err := col.RemoveAll(bson.M{
		"_id": bson.M{
			"$ne": snapshot["_id"],
		},
	})

	if err != nil {
		log.Errorf("Error while removing old breaking news snapshots %v", err)
	}
}

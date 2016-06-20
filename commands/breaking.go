package commands

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	api "github.com/michigan-com/gannett-newsfetch/gannettApi"
	m "github.com/michigan-com/gannett-newsfetch/model"
	"github.com/michigan-com/newsfetch/lib"
)

func FetchBreakingNews(session *mgo.Session, siteCodes []string, gnapiDomain string, loopInterval time.Duration, gannettAPIKey string) {
	for {
		var breakingWait sync.WaitGroup
		breakingChannel := make(chan *m.SearchArticle, len(siteCodes)*100)
		log.Info("Fetching breaking news...")
		for _, siteCode := range siteCodes {
			breakingWait.Add(1)
			go func(siteCode string) {
				defer breakingWait.Done()
				articles := api.GetBreakingNews(siteCode, gannettAPIKey)

				for _, article := range articles {
					breakingChannel <- article
				}
			}(siteCode)
		}
		breakingWait.Wait()
		close(breakingChannel)
		log.Info("...Done fetching articles")

		SaveBreakingArticles(breakingChannel, session)

		if gnapiDomain != "" {
			gnapiUrl := fmt.Sprintf("%s/%s/", gnapiDomain, "breaking-news")
			resp, err := http.Get(gnapiUrl)
			if err == nil {
				resp.Body.Close()
			}
		}

		if loopInterval > 0 {
			log.Infof("Sleeping for %d ms...", loopInterval/time.Millisecond)
			time.Sleep(loopInterval)
			log.Info("...and now I'm awake!")
		} else {
			break
		}
	}
}

func SaveBreakingArticles(breakingChannel chan *m.SearchArticle, session *mgo.Session) []*m.BreakingNewsArticle {
	log.Info("Collecting breaking articles to save...")
	breakingNewsSnapshot := m.BreakingNewsSnapshot{}
	breakingArticles := make([]*m.BreakingNewsArticle, 0, 100)
	toScrape := make([]interface{}, 0, 100)
	articleCol := session.DB("").C("Article")
	for breaking := range breakingChannel {
		articleId := lib.GetArticleId(breaking.Urls.LongUrl)
		if articleId == -1 {
			log.Warningf(`Failed to get id for url %s`, breaking.Urls.LongUrl)
			continue
		}

		// Determine if the breaking article still needs to be scraped. If so,
		// wait until we find the article before sending out the breaking news alert
		storedArticle := &m.Article{}
		err := articleCol.Find(bson.M{"article_id": articleId}).One(storedArticle)
		if err == mgo.ErrNotFound || len(storedArticle.Summary) != 3 {
			articleIdQuery := bson.M{"article_id": articleId}
			toScrape = append(toScrape, articleIdQuery)
			toScrape = append(toScrape, articleIdQuery)
		} else {
			// only add a breaking news article if we've summarized and scraped it
			breakingArticle := &m.BreakingNewsArticle{
				ArticleId:   storedArticle.ArticleId,
				Headline:    storedArticle.Headline,
				Subheadline: storedArticle.Subheadline,
				Photo:       storedArticle.Photo,
				Video:       storedArticle.Video,
			}
			breakingArticles = append(breakingArticles, breakingArticle)
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
		return breakingArticles
	}

	RemoveOldBreakingSnapshot(breakingCol)

	log.Infof("...Done saving breaking articles, count: %d", len(breakingArticles))
	return breakingArticles
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

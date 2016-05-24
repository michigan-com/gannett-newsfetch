package commands

import (
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/cobra"
	"gopkg.in/mgo.v2/bson"

	"github.com/michigan-com/gannett-newsfetch/config"
	api "github.com/michigan-com/gannett-newsfetch/gannettApi"
	"github.com/michigan-com/gannett-newsfetch/lib"
)

var cleanupCommand = &cobra.Command{
	Use:   "scrape-and-summarize",
	Short: "Grab stories that we see in chartbeat but not the Gannett API",
	Run:   scrapeAndSummarize,
}

func scrapeAndSummarize(command *cobra.Command, args []string) {
	var envConfig, _ = config.GetEnv()
	var articleWait sync.WaitGroup
	if envConfig.MongoUri == "" {
		log.Warning("No mongo URI specified, this command is basically useless")
		return
	}

	session := lib.DBConnect(envConfig.MongoUri)
	toScrape := session.DB("").C("ToScrape")
	defer session.Close()

	for {
		toSummarize := make([]interface{}, 0, 100)
		articleIdsToScrape := make([]map[string]int, 0, 100)

		log.Info("Finding articles in need of scraping...")
		err := toScrape.Find(bson.M{}).Select(bson.M{"article_id": true, "_id": false}).All(&articleIdsToScrape)
		if err != nil {
			log.Errorf("Error getting articles IDs from ToScrape collection: %v", err)
		}

		if len(articleIdsToScrape) > 0 {
			log.Infof("...scraping %d articles...", len(articleIdsToScrape))
			for _, articleIdObj := range articleIdsToScrape {
				articleWait.Add(1)
				articleId := articleIdObj["article_id"]
				go func(articleId int) {
					defer articleWait.Done()
					assetArticle, assetPhoto := api.GetAssetArticleAndPhoto(articleId)

					mongoArticle := api.FormatAssetArticleForSaving(assetArticle, assetPhoto)
					mongoArticle.Save(session)

					articleIdQuery := bson.M{"article_id": articleId}
					toSummarize = append(toSummarize, articleIdQuery)
					toSummarize = append(toSummarize, articleIdQuery)

					toScrape.Remove(articleIdQuery)
				}(articleId)
			}
			log.Infof("...Done scraping articles")
			articleWait.Wait()
		} else {
			log.Infof("...no articles in need of scraping")
		}

		if len(toSummarize) > 0 {
			log.Info("Summarizing articles...")
			_, err := ProcessSummaries(toSummarize, session)
			if err != nil {
				log.Errorf("Failed to process summaries: %v", err)
			}
			log.Info("...Done processing summaries")
		} else {
			log.Info("No articles to summarize.")
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

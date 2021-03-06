package commands

import (
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/andreyvit/sem"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/michigan-com/brvty-api/mongoqueue"
	api "github.com/michigan-com/gannett-newsfetch/gannettApi"
	m "github.com/michigan-com/gannett-newsfetch/model"
)

func ScrapeAndSummarize(session *mgo.Session, queue *mongoqueue.Queue, loopInterval time.Duration, mongoUri string, assetApiKey string) {
	var articleWait sync.WaitGroup
	s := sem.New(100)

	toScrape := session.DB("").C("ToScrape")

	for {
		toSummarize := make([]interface{}, 0, 100)
		var requests []m.ScrapeRequest

		log.Info("Finding articles in need of scraping...")
		err := toScrape.Find(bson.M{}).Select(bson.M{"article_id": true, "article_url": true, "_id": false}).All(&requests)
		if err != nil {
			log.Errorf("Error loading ToScrape collection: %v", err)
		}

		if len(requests) > 0 {
			log.Infof("...scraping %d articles...", len(requests))
			for _, request := range requests {
				articleWait.Add(1)

				s.Exec((func(request m.ScrapeRequest) func() {
					return func() {
						defer articleWait.Done()
						assetArticleContent := api.GetAssetArticleContent(request.ArticleID, assetApiKey)

						// if queue != nil && request.ArticleURL != "" {
						// 	err := queue.Add(mongoqueue.Request{
						// 		Name: fmt.Sprintf("brvty-%v", request.ArticleURL),
						// 		Op:   OpBrvty,
						// 		Args: map[string]interface{}{
						// 			ParamArticleID: request.ArticleID,
						// 			ParamURL:       request.ArticleURL,
						// 		},
						// 	})
						// 	if err != nil {
						// 		log.Errorf("Failed to enqueue brvty job for article at %v: %v", request.ArticleURL, err)
						// 		os.Exit(22)
						// 	}
						// }

						mongoArticle := api.FormatAssetArticleForSaving(assetArticleContent)
						mongoArticle.Save(session)

						articleIdQuery := bson.M{"article_id": request.ArticleID}
						toSummarize = append(toSummarize, articleIdQuery)
						toSummarize = append(toSummarize, articleIdQuery)

						toScrape.Remove(articleIdQuery)
					}
				})(request))
			}

			log.Infof("...Done scraping articles")
			articleWait.Wait()

			// enqueue postback jobs to deliver newsfetch bodies and summaries to Brvty (for analysis and introspection)
			// if queue != nil {
			// 	for _, request := range requests {
			// 		if request.ArticleURL != "" {
			// 			err := queue.Add(mongoqueue.Request{
			// 				Name: fmt.Sprintf("brvty-postback-%v", request.ArticleURL),
			// 				Op:   OpBrvtyPostback,
			// 				Args: map[string]interface{}{
			// 					ParamArticleID: request.ArticleID,
			// 					ParamURL:       request.ArticleURL,
			// 				},
			// 			})
			// 			if err != nil {
			// 				log.Errorf("Failed to enqueue brvty postback job for article at %v: %v", request.ArticleURL, err)
			// 				os.Exit(22)
			// 			}
			// 		}
			// 	}
			// }
		} else {
			log.Infof("...no articles in need of scraping")
		}

		// if len(toSummarize) > 0 {
		// 	log.Info("Summarizing articles...")
		// 	_, err := ProcessSummaries(session, toSummarize, mongoUri, summaryVEnv)
		// 	if err != nil {
		// 		log.Errorf("Failed to process summaries: %v", err)
		// 	}
		// 	log.Info("...Done processing summaries")
		// } else {
		// 	log.Info("No articles to summarize.")
		// }

		if loopInterval > 0 {
			log.Infof("Sleeping for %d ms...", loopInterval/time.Millisecond)
			time.Sleep(loopInterval)
			log.Info("...and now I'm awake!")
		} else {
			break
		}
	}
}

func pluckRequestURLs(requests []m.ScrapeRequest) []string {
	result := make([]string, 0, len(requests))
	for _, request := range requests {
		result = append(result, request.ArticleURL)
	}
	return result
}

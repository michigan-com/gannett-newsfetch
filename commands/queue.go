package commands

import (
	"time"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/michigan-com/brvty-api/brvtyclient"
	"github.com/michigan-com/brvty-api/mongoqueue"
	m "github.com/michigan-com/gannett-newsfetch/model"
)

type extractor struct {
	session      *mgo.Session
	client       *brvtyclient.Client
	brvtyTimeout time.Duration
}

func RunQueuedJobs(session *mgo.Session, client *brvtyclient.Client, queue *mongoqueue.Queue, brvtyTimeout time.Duration) {
	ex := extractor{session, client, brvtyTimeout}

	queue.Run(map[string]mongoqueue.Worker{
		OpBrvty: mongoqueue.Worker{
			Func: mongoqueue.WorkerFunc(ex.extract),
		},
		OpBrvtyPostback: mongoqueue.Worker{
			Func: mongoqueue.WorkerFunc(ex.postback),
		},
	}, mongoqueue.RunParams{
		PollInterval: 1 * time.Second,
		Shutdown:     nil,
	})
}

func (ex *extractor) extract(op string, args map[string]interface{}) error {
	aid := args[ParamArticleID].(int)
	if aid == 0 {
		return errors.New("Invalid params: missing article ID")
	}
	url := args[ParamURL].(string)
	if url == "" {
		return errors.New("Invalid params: missing URL")
	}

	resources, err := ex.client.Add([]string{url}, ex.brvtyTimeout)
	if err != nil {
		return errors.Wrap(err, "brvty.Add failed")
	}
	resource := resources[0]

	articlesC := ex.session.DB("").C("Article")

	body := resource.OptimalBody()
	if body == nil {
		return errors.New("No body yet")
	}

	summary := resource.OptimalSummary()
	if summary == nil {
		return errors.New("No summary yet")
	}

	err = articlesC.Update(bson.M{"article_id": aid}, bson.M{
		"$set": bson.M{
			"brvty": bson.M{
				"headline": body.Headline,
				"text":     body.Text,
				"summary":  summary.Sentences,
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "saving failed")
	}

	return nil
}

func (ex *extractor) postback(op string, args map[string]interface{}) error {
	aid := args[ParamArticleID].(int)
	if aid == 0 {
		return errors.New("Invalid params: missing article ID")
	}
	url := args[ParamURL].(string)
	if url == "" {
		return errors.New("Invalid params: missing URL")
	}

	var article *m.Article
	articlesC := ex.session.DB("").C("Article")
	err := articlesC.Find(bson.M{"article_id": aid}).One(&article)
	if err != nil {
		return errors.Wrap(err, "failed to load newsfetch article")
	}

	if article.Body != "" {
		err = ex.client.UpdateBody(url, "newsfetch", article.Headline, article.Body)
		if err != nil {
			return errors.Wrap(err, "UpdateBody failed")
		}
	}

	if len(article.Summary) > 0 {
		err = ex.client.UpdateSummary(url, "newsfetch", article.Summary)
		if err != nil {
			return errors.Wrap(err, "UpdateSummary failed")
		}
	}

	return nil
}

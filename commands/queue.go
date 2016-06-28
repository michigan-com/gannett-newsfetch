package commands

import (
	"time"

	"github.com/pkg/errors"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/michigan-com/brvty-api/brvtyclient"
	"github.com/michigan-com/brvty-api/mongoqueue"
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

	if resource.Summary == nil {
		return errors.New("No summary yet")
	}

	err = articlesC.Update(bson.M{"article_id": aid}, bson.M{
		"$set": bson.M{
			"brvty": bson.M{
				"headline": resource.Body.Headline,
				"text":     resource.Body.Text,
				"summary":  resource.Summary.Sentences,
			},
		},
	})
	if err != nil {
		return errors.Wrap(err, "saving failed")
	}

	return nil
}

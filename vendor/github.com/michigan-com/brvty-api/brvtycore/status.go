package brvtycore

import (
	// 	"time"

	// 	"gopkg.in/mgo.v2"
	"github.com/michigan-com/brvty-api/mongoqueue"
	"gopkg.in/mgo.v2/bson"
)

type Status struct {
	mongoqueue.Stats
	ResourceCount int
}

func (c *Core) CheckStatus() (*Status, error) {
	resCount, err := c.collection.Find(bson.M{}).Count()

	stats, err := c.Queue.Stats()
	if err != nil {
		return nil, err
	}

	return &Status{*stats, resCount}, nil
}

func (c *Core) GetRecent(n int) ([]*Resource, error) {
	var rr []*Resource
	// didn't work: .Select(bson.M{"url": 1, "body": 1, "ctime": 1, "summary": 1})
	err := c.collection.Find(bson.M{}).Sort("-ctime").Limit(n).All(&rr)
	return rr, err
}

package brvtycore

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/michigan-com/brvty-api/canonicalurl"
	"github.com/michigan-com/brvty-api/gannett"
)

type AddParams struct {
	// Tags []string
}

func (c *Core) Add(urls []string, params AddParams) error {
	for _, url := range urls {
		err := c.AddOne(url, params)
		if err != nil {
			return err
		}
	}
	return nil
}

func (c *Core) AddOne(url string, params AddParams) error {
	url, err := canonicalurl.CleanURLString(url, canonicalurl.Safe)
	if err != nil {
		// don't allow to enqueue invalid URLs
		return err
	}

	id := computeID(url)
	now := bson.Now()

	var tags []string
	gannettId := gannett.FindArticleID(url)
	if gannettId != gannett.IDNotFound {
		tags = append(tags, TagGannett)
	}

	var res *Resource
	_, err = c.collection.FindId(id).Apply(mgo.Change{
		Upsert:    true,
		ReturnNew: true,
		Update: bson.M{
			"$setOnInsert": bson.M{
				"v":    1,
				"url":  url,
				"tags": tags,

				"ctime": now,
			},
			// "$set": bson.M{
			// 	"extr.pending": true,
			// },
			// "$inc": bson.M{
			// 	"rev": 1,
			// },
			// "$max": bson.M{
			// 	"reqtime": now,
			// },
		},
	}, &res)
	if err != nil {
		return err
	}

	err = c.enqueueExtraction(res)
	if err != nil {
		return err
	}

	// err = collection.Update(bson.M{
	// 	"_id": id,
	// 	"rev": bson.M{"$ne": rev},
	// }, bson.M{
	// 	"$set": bson.M{
	// 		"rev":      rev,
	// 		"text":     payload.Text,
	// 		"headline": payload.Headline,
	// 		"url":      payload.URL,
	// 		"tags":     payload.Tags,
	// 		"revtime":  now,

	// 		"pending": true,
	// 		"ptime":   now,
	// 	},
	// 	"$inc": bson.M{"nrevs": 1},
	// })
	// if err == nil {
	// 	return true, nil
	// }
	// if err != mgo.ErrNotFound {
	// 	return false, err
	// }

	// err = collection.Update(bson.M{
	// 	"_id":     id,
	// 	"rev":     rev,
	// 	"pending": false,
	// 	"failed":  true,
	// }, bson.M{
	// 	"$set": bson.M{
	// 		"pending": true,
	// 		"ptime":   now,
	// 	},
	// 	"$inc": bson.M{"nretries": 1},
	// })
	// if err == nil {
	// 	return true, nil
	// }
	// if err != mgo.ErrNotFound {
	// 	return false, err
	// }

	return nil
}

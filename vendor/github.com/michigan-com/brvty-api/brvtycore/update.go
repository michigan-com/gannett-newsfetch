package brvtycore

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

func (c *Core) UpdateBody(url string, strategy string, body ResourceBody) error {
	bodyHash := body.Hash()
	id := computeID(url)

	var res *Resource
	err := c.collection.FindId(id).One(&res)
	if err != nil {
		return err
	}

	isOptimalStrategy := (strategy == res.OptimalExtractionStrategy())

	set := bson.M{}
	set[fmt.Sprintf("bodies.%s", strategy)] = body

	if isOptimalStrategy {
		set["bodyhash"] = bodyHash
	}
	if isOptimalStrategy || !res.HasBody() {
		set["bodykey"] = strategy
	}

	_, err = c.collection.FindId(id).Apply(mgo.Change{
		Update: bson.M{
			"$set": set,
		},
		ReturnNew: true,
	}, &res)
	if err != nil {
		return err
	}

	if isOptimalStrategy {
		err = c.enqueueSummarization(id, url, bodyHash)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c *Core) UpdateSummary(url string, strategy string, summary ResourceSummary) error {
	id := computeID(url)

	var res *Resource
	err := c.collection.FindId(id).One(&res)
	if err != nil {
		return err
	}

	set := bson.M{}
	set[fmt.Sprintf("summaries.%s", strategy)] = summary

	if strategy == res.OptimalSummarizationStrategy() || !res.HasSummary() {
		set["summarykey"] = strategy
	}

	err = c.collection.UpdateId(id, bson.M{
		"$set": set,
	})
	return err
}

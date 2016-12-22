package brvtycore

import (
	"gopkg.in/mgo.v2/bson"
)

func (c *Core) GetOne(url string) (*Resource, error) {
	rr, err := c.Get([]string{url})
	if err != nil {
		return nil, err
	} else {
		return rr[0], nil
	}
}

// getUnordered returns the known resources corresponding to the given urls.
// It skips any previously unknown urls and returns the results in arbitrary order.
// All returned items are non-nil.
func (c *Core) getUnordered(urls []string) ([]*Resource, error) {
	ids := computeIDsForURLs(urls)

	var rr []*Resource
	err := c.collection.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&rr)
	return rr, err
}

// Get returns the resources corresponding to the given urls.
//
// Upon successful completion, returns exactly the same number of items as
// the input slice of urls. Each returned item is a previously added resource
// for the corresponding url, or nil if the url doesn't match any known
// resources.
func (c *Core) Get(urls []string) ([]*Resource, error) {
	ids := computeIDsForURLs(urls)

	var unordered []*Resource
	err := c.collection.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&unordered)
	if err != nil {
		return nil, err
	}

	indexed := make(map[string]*Resource, len(unordered))
	for _, res := range unordered {
		indexed[res.Id] = res
	}

	var ordered []*Resource
	for _, id := range ids {
		ordered = append(ordered, indexed[id])
	}

	return ordered, nil
}

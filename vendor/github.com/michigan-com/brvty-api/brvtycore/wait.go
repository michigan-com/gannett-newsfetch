package brvtycore

import (
	// "fmt"
	"gopkg.in/mgo.v2/bson"

	"github.com/michigan-com/brvty-api/mongoqueue"
)

type WaitParams mongoqueue.WaitParams

func (c *Core) Wait(urls []string, params WaitParams) (bool, error) {
	qparams := mongoqueue.WaitParams(params)
	qparams.StartDeadline()

	ids := computeIDsForURLs(urls)

	// Repeat this in a loop because computeJobIDs may only be able to produce
	// further job IDs after prior jobs finish.
	prevJobIds := []string{}
	for {
		var rr []*Resource
		err := c.collection.Find(bson.M{"_id": bson.M{"$in": ids}}).All(&rr)
		if err != nil {
			return false, err
		}

		jobIds := computeJobIDs(rr)
		// fmt.Printf("brvtycore Wait: jobIds = %v, prevJobIds = %v\n", jobIds, prevJobIds)
		if stringsEqual(prevJobIds, jobIds) {
			return true, nil
		}

		done, err := c.Queue.Wait(jobIds, qparams)
		if err != nil {
			return false, err
		}
		if !done || qparams.IsAfterDeadline() {
			// fmt.Printf("timeout!!\n")
			return false, nil
		}

		prevJobIds = jobIds
	}
}

func stringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, el := range a {
		if el != b[i] {
			return false
		}
	}
	return true
}

package brvtycore

import (
	"fmt"

	"github.com/michigan-com/brvty-api/gannett"
	"github.com/michigan-com/brvty-api/mongoqueue"
)

func computeExtractionJobName(id, strategy string) string {
	return fmt.Sprintf("%s-extract-%s", id, strategy)
}

func ExtractionJobOp(strategy string) string {
	return fmt.Sprintf("%s-%s", OpExtractBasename, strategy)
}

func computeSummarizationJobName(id, bodyHash string) string {
	return fmt.Sprintf("%s-summarize-%s", id, bodyHash)
}

func computeJobIDs(resources []*Resource) []string {
	jobIds := make([]string, 0, 3*len(resources))
	for _, res := range resources {
		if res.HasTag(TagGannett) {
			jobIds = append(jobIds, computeExtractionJobName(res.Id, ExtractionStrategyGannett))
		} else {
			jobIds = append(jobIds, computeExtractionJobName(res.Id, ExtractionStrategyGeneric))
		}
		if res.BodyHash != "" {
			jobIds = append(jobIds, computeSummarizationJobName(res.Id, res.BodyHash))
		}
	}
	return jobIds
}

func (c *Core) enqueueExtraction(res *Resource) error {
	id := res.Id

	gannettId := gannett.FindArticleID(res.URL)
	if gannettId != gannett.IDNotFound {
		err := c.Queue.Add(mongoqueue.Request{
			Name: computeExtractionJobName(id, ExtractionStrategyGannett),
			Op:   ExtractionJobOp(ExtractionStrategyGannett),
			Args: map[string]interface{}{
				ParamResourceURL:      res.URL,
				ParamGannettArticleID: gannettId,
			},
		})
		if err != nil {
			return err
		}
	}

	err := c.Queue.Add(mongoqueue.Request{
		Name: computeExtractionJobName(id, ExtractionStrategyGeneric),
		Op:   ExtractionJobOp(ExtractionStrategyGeneric),
		Args: map[string]interface{}{
			ParamResourceURL: res.URL,
		},
	})
	if err != nil {
		return err
	}

	return nil
}

func (c *Core) enqueueSummarization(id, url, bodyHash string) error {
	return c.Queue.Add(mongoqueue.Request{
		Name: computeSummarizationJobName(id, bodyHash),
		Op:   OpSummarize,
		Args: map[string]interface{}{
			ParamResourceURL: url,
		},
	})
}

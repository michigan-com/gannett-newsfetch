package brvtycore

import (
	"time"

	"gopkg.in/mgo.v2"

	"github.com/michigan-com/brvty-api/canonicalurl"
	"github.com/michigan-com/brvty-api/mongoqueue"
)

type Core struct {
	DB    *mgo.Database
	Queue *mongoqueue.Queue

	collection *mgo.Collection
}

type Priority int

const (
	PriorityBackground    Priority = 0
	PriorityUserInitiated Priority = 500
	PriorityRealtime      Priority = 1000
)

type ResourceId struct {
	Namespace string
	Name      string
}

type ResourcePayload struct {
	Text     string
	Headline string
	URL      string
	Tags     []string
}

const (
	ResourceVersionCurrent = 1
)

type ExtractionRequest struct {
	ResourceId string `json:"id"`
	Revision   int    `json:"rev"`
	Attempt    int    `json:"attempt"`

	URL string `json:"url"`
}

type SummarizationRequest struct {
	ResourceId string `json:"id"`
	Revision   int    `json:"rev"`
	Attempt    int    `json:"attempt"`

	Headline string `json:"headline"`
	Text     string `json:"text"`
}

type ResourceBody struct {
	Headline string `json:"headline" bson:"headline"`
	Text     string `json:"text" bson:"text"`
}

type ResourceSummary struct {
	Sentences []string `json:"sentences,omitempty" bson:"sentences,omitempty"`
}

type Resource struct {
	Id      string `json:"id" bson:"_id"`
	Version int    `json:"v" bson:"v"`

	URL string `json:"url" bson:"url"`
	// Tags []string `json:"tags" bson:"tags"`

	// ExtractionPending indicates whether there's an extraction job scheduled or running.
	// ExtractionPending bool `json:"pending"`

	// ExtractionFailed bool `json:"failed"`

	// CreationTime is the time when the URL has been encountered for the first time.
	CreationTime time.Time `json:"ctime" bson:"ctime"`

	// // Revision is incremented every time there's a request to scrape/re-scrape the article.
	// Revision int `json:"rev" bson:"rev"`

	// // RevisionTime is the last time Revision has been incremented.
	// RevisionTime time.Time `json:"revtime" bson:"revtime"`

	// // ExtractionTime is the time of the last successful extraction.
	// ExtractionTime time.Time `json:"extrtime" bson:"extrtime"`

	// extraction result
	// Body     ResourceBody `json:"body,omitempty" bson:",inline"`
	// HasBody  bool         `json:"hasbody" bson:"hasbody"`

	Tags []string `json:"tags" bson:"tags"`

	Bodies    map[string]ResourceBody    `json:"bodies" bson:"bodies"`
	Summaries map[string]ResourceSummary `json:"summaries" bson:"summaries"`

	// for the optimal body only (the one that will be summarized)
	BodyHash string `json:"bodyhash,omitempty" bson:"bodyhash"`

	BestBodyKey    string `json:"bodykey" bson:"bodykey"`
	BestSummaryKey string `json:"summarykey" bson:"summarykey"`

	// Summary    *ResourceSummary `json:"summary,omitempty" bson:"summary,omitempty"`
	// HasSummary bool             `json:"hassummary" bson:"hassummary"`
}

func (r *Resource) HasTag(tag string) bool {
	for _, t := range r.Tags {
		if t == tag {
			return true
		}
	}
	return false
}

func (r *Resource) PrettyURL() string {
	url, err := canonicalurl.CleanURLString(r.URL, canonicalurl.Short)
	if err != nil {
		return r.URL
	}
	return url
}

func (r *Resource) OptimalExtractionStrategy() string {
	if r.HasTag(TagGannett) {
		return ExtractionStrategyGannett
	} else {
		return ExtractionStrategyGeneric
	}
}

func (r *Resource) OptimalSummarizationStrategy() string {
	return SummarizationStrategyNLTKBasedV1
}

func (r *Resource) OptimalBody() *ResourceBody {
	key := r.OptimalExtractionStrategy()
	if result, ok := r.Bodies[key]; ok {
		return &result
	} else {
		return nil
	}
}

func (r *Resource) OptimalSummary() *ResourceSummary {
	key := r.OptimalSummarizationStrategy()
	if result, ok := r.Summaries[key]; ok {
		return &result
	} else {
		return nil
	}
}

func (r *Resource) HasBody() bool {
	return r.BestBodyKey != ""
}

func (r *Resource) HasSummary() bool {
	return r.BestSummaryKey != ""
}

func (r *Resource) Body() *ResourceBody {
	if r.BestBodyKey == "" {
		return nil
	}
	result := r.Bodies[r.BestBodyKey]
	return &result
}

func (r *Resource) Summary() *ResourceSummary {
	if r.BestSummaryKey == "" {
		return nil
	}
	result := r.Summaries[r.BestSummaryKey]
	return &result
}

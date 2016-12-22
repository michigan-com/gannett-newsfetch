package brvtycore

import (
	"gopkg.in/mgo.v2"
)

var ErrNotFound error = mgo.ErrNotFound

const (
	OpExtractBasename = "extract"
	OpSummarize       = "summarize"
)

const (
	ParamResourceURL      = "url"
	ParamGannettArticleID = "articleID"
)

const (
	ExtractionStrategyGeneric      = "readability"
	ExtractionStrategyGannett      = "gannett"
	ExtractionStrategyOldNewsfetch = "newsfetch"
)

const (
	SummarizationStrategyNLTKBasedV1  = "nltk-1"
	SummarizationStrategyOldNewsfetch = "newsfetch"
)

const (
	TagGannett = "gannett"
)

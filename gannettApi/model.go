package gannettApi

import (
	m "github.com/michigan-com/gannett-newsfetch/model"
)

/**
 * For use with the /asset/ API
 */
type AssetArticle struct {
	AssetId            int              `json:"id"`
	Headline           string           `json:"title"`
	Ssts               Ssts             `json:"ssts"`
	Links              links            `json:"_links"`
	PublishDate        string           `json:"publishDate"`
	InitialPublishDate string           `json:"initialPublishDate"`
	PromoBrief         string           `json:"promoBrief"`
	Attribution        AssetAttribution `json:"attribution"`
	FullText           string           `json:"fullText"`
	StoryHighlights    []string         `json:"storyHighlights"`
}

type ArticlesResponse struct {
	TotalNumResults int
	NumResults      int
	Results         []*SearchArticle
}

/**
 * Use for the /search/v4 api
 */
type SearchArticle struct {
	AssetId       int    `json:"assetId"`
	Headline      string `json:"headline"`
	Ssts          Ssts   `json:"ssts"`
	Urls          urls   `json:"urls"`
	DatePublished string `json:"datePublished"`
	SolrFields    Solr   `json:"requestedSolrFields"`
	PromoBrief    string `json:"promoBrief"`
}

type Ssts struct {
	Section    string `json:"section"`
	SubSection string `json:"subSection"`
	Topic      string `json:"topic"`
	SubTopic   string `json:"subTopic"`
}

type links struct {
	LongUrl  AssetUrl    `json:"longUrl"`
	ShortUrl AssetUrl    `json:"shortUrl"`
	Photo    *AssetPhoto `json:"photo"` // have to hit a second API to get photo info
}

type AssetUrl struct {
	Href string `json:"href"`
}

type urls struct {
	LongUrl  string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}

type attribution struct {
	Author string `json:"author"`
}

type AssetAttribution struct {
	Author string `json:"byline"`
}

type PhotoInfo struct {
	AbsoluteUrl    string                 `json:"absoluteUrl"`
	Crops          map[string]m.PhotoInfo `json:"crops"`
	Caption        string                 `json:"caption"`
	Credit         string                 `json:"credit"`
	OriginalWidth  int                    `json:"originalWidth"`
	OriginalHeight int                    `json:"originalHeight"`
}

type AssetPhotoInfo struct {
	AbsoluteUrl    string          `json:"absoluteUrl"`
	PublishUrl     string          `json:"publishUrl"`
	Crops          []Crop          `json:"crops"`
	Caption        string          `json:"caption"`
	Credit         string          `json:"credit"`
	OriginalWidth  string          `json:"originalWidth"`
	OriginalHeight string          `json:"originalHeight"`
	Attributes     PhotoAttributes `json:"_attributes"`
}

type PhotoAttributes struct {
	SmallBaseName string `json:"smallbasename"`
	SImageHeight  string `json:"simageheight"`
	SImageWidth   string `json:"simagewidth"`
}

type Crop struct {
	Name   string `json:"name"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Path   string `json:"path"`
}

type AssetPhoto struct {
	Id int `json:"id"`
}

type Solr struct {
	InitalPublished string `json:"initialpublished"`
}

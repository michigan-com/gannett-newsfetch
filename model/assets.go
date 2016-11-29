package model

import (
	"encoding/json"
)

/* Wrapper around all the content we get returned in API calls */
type AssetArticleContent struct {
	Article *AssetArticle
	Assets  *ArticleAssets
}

type ArticleAssets struct {
	Photo         *AssetPhoto
	Video         *AssetVideo
	VideoPlaylist *AssetVideoPlaylist
}

/**
 * For use with the /asset/ API
 */
type AssetArticle struct {
	AssetId            int              `json:"id"`
	Headline           string           `json:"title"`
	Ssts               Ssts             `json:"ssts"`
	Links              Links            `json:"_links"`
	PublishDate        string           `json:"publishDate"`
	InitialPublishDate string           `json:"initialPublishDate"`
	PromoBrief         string           `json:"promoBrief"`
	Attribution        AssetAttribution `json:"attribution"`
	FullText           string           `json:"fullText"`
	StoryHighlights    []string         `json:"storyHighlights"`
}

type AssetPhoto struct {
	AbsoluteUrl    string          `json:"absoluteUrl"`
	PublishUrl     string          `json:"publishUrl"`
	Crops          []Crop          `json:"crops"`
	Caption        string          `json:"caption"`
	Credit         string          `json:"credit"`
	OriginalWidth  string          `json:"originalWidth"`
	OriginalHeight string          `json:"originalHeight"`
	Attributes     PhotoAttributes `json:"_attributes"`
}

type AssetVideo struct {
	Thumbnail  string           `json:"thumbnail"`
	VideoStill string           `json:"videoStill"`
	Length     string           `json:"length"`
	Renditions []VideoRendition `json:"renditions"`
}

type AssetVideoPlaylist struct {
	AssetId int `json:"id"`
}

type Ssts struct {
	Section    string `json:"section"`
	SubSection string `json:"subSection"`
	Topic      string `json:"topic"`
	SubTopic   string `json:"subTopic"`
}

type Links struct {
	LongUrl  AssetUrl           `json:"longUrl"`
	ShortUrl AssetUrl           `json:"shortUrl"`
	Assets   []*GannettApiAsset `json:"assets"`
	Photo    *PhotoAsset        `json:"photo"`
}

type AssetUrl struct {
	Href string `json:"href"`
}

type Urls struct {
	LongUrl  string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}

type attribution struct {
	Author string `json:"author"`
}

type AssetAttribution struct {
	Author string `json:"byline"`
}

// type AssetPhotoInfo struct {
// 	AbsoluteUrl    string               `json:"absoluteUrl"`
// 	Crops          map[string]PhotoInfo `json:"crops"`
// 	Caption        string               `json:"caption"`
// 	Credit         string               `json:"credit"`
// 	OriginalWidth  int                  `json:"originalWidth"`
// 	OriginalHeight int                  `json:"originalHeight"`
// }

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

type Solr struct {
	InitalPublished string `json:"initialpublished"`
}

type VideoRendition struct {
	EncodingRate int    `json:"encodingRate"` // idk, might be useful
	Height       int    `json:"height"`
	Width        int    `json:"width"`
	Size         int    `json:"size"`
	Url          string `json:"url"`
	Duration     int    `json:"duration"`
}

/** Generic type for pointers to other assets (e.g. photos, videos, galleries, etc) */
type GannettApiAsset struct {
	Id                    int    `json:"id"`
	Type                  string `json:"type"`
	RelationshipTypeFlags string `json:"relationshipTypeFlags"` // used to find if this is the primary image
}

type PhotoAsset struct {
	Id int `json:"id"`
}

/**
 * Use for the /search/v4 api
 */
type ArticlesResponse struct {
	TotalNumResults int
	NumResults      int
	Results         []*SearchArticle
}

type SearchArticle struct {
	AssetId       int    `json:"assetId"`
	Headline      string `json:"headline"`
	Ssts          Ssts   `json:"ssts"`
	Urls          Urls   `json:"urls"`
	DatePublished string `json:"datePublished"`
	SolrFields    Solr   `json:"requestedSolrFields"`
	PromoBrief    string `json:"promoBrief"`
}

/** Interfaces  for reading from JSON */
type AssetResp interface {
	Decode(*json.Decoder) error
}

func (a *AssetArticle) Decode(decoder *json.Decoder) error {
	return decoder.Decode(a)
}

func (p *AssetPhoto) Decode(decoder *json.Decoder) error {
	return decoder.Decode(p)
}

func (v *AssetVideo) Decode(decoder *json.Decoder) error {
	return decoder.Decode(v)
}

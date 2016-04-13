package gannettApi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	api "github.com/michigan-com/gannett-newsfetch/gannettApi"
)

type ArticlesResponse struct {
	TotalNumResults int
	NumResults      int
	Results         []*ArticleIn
}

type ArticleIn struct {
	AssetId       int         `json:"assetId"`
	Headline      string      `json:"headline"`
	Ssts          ssts        `json:"ssts"`
	Urls          urls        `json:"urls"`
	DatePublished string      `json:"datePublished"`
	SolrFields    Solr        `json:"requestedSolrFields"`
	PromoBrief    string      `json:"promoBrief"`
	Attribution   attribution `json:"attribution"`
	Photo         PhotoInfo   `json:"photo"`
}

type ssts struct {
	Section    string `json:"section"`
	SubSection string `json:"subSection"`
	Topic      string `json:"topic"`
	SubTopic   string `json:"subTopic"`
}

type urls struct {
	LongUrl  string `json:"longUrl"`
	ShortUrl string `json:"shortUrl"`
}

type attribution struct {
	Author string `json:"author"`
}

type PhotoInfo struct {
	AbsoluteUrl    string            `json:"absoluteUrl"`
	Crops          map[string]string `json:"crops"`
	Caption        string            `json:"caption"`
	Credit         string            `json:"credit"`
	OriginalWidth  int               `json:"originalWidth"`
	OriginalHeight int               `json:"originalHeight"`
}

type Solr struct {
	InitalPublished string `json:"initialpublished"`
}

/*
	Get articles within an entire day
*/
func GetArticlesByDay(siteCode string, date time.Time) []*ArticleIn {
	year, monthObj, day := date.Date()
	month := int(monthObj)

	startDate := api.FormatAsDateString(year, month, day, 0, 0, 0)
	endDate := api.FormatAsDateString(year, month, day, 23, 59, 59)

	return getArticles(siteCode, startDate, endDate)
}

/*
	Get articles within a date range
	Date strings should be formatted via gannettApi.FormatAsDateString
*/
func getArticles(siteCode string, startDate string, endDate string) []*ArticleIn {
	var articleResponse *ArticlesResponse = &ArticlesResponse{}
	var queryParams = api.GetDefaultSearchValues(siteCode)
	queryParams.Add("fq", fmt.Sprintf("initialpublished:[%s TO %s]", startDate, endDate))
	queryParams.Add("fq", "assettypename:text")
	queryParams.Add("fl", "initialpublished")
	queryParams.Add("sort", "lastpublished desc")

	url := fmt.Sprintf("%s?%s", api.GannettApiSearchRoot, queryParams.Encode())
	log.Println(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Warningf(`

		Failed to get articles for site %s. http.Get failed:

			Err: %v
		`, siteCode, err)
		return articleResponse.Results
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(articleResponse)

	if err != nil {
		log.Warningf(`

		Failed to get articles for site %s. Json decoding failed

			Err: %v
		`, siteCode, err)
		return []*ArticleIn{}
	}

	return articleResponse.Results
}

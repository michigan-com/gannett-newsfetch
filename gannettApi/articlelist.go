package gannettApi

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"

	"github.com/andreyvit/debugflag"

	m "github.com/michigan-com/gannett-newsfetch/model"
)

/*
	Get articles within an entire day
*/
func GetArticlesByDay(siteCode string, date time.Time) []*m.SearchArticle {
	year, monthObj, day := date.Date()
	month := int(monthObj)

	startDate := FormatAsDateString(year, month, day, 0, 0, 0)
	endDate := FormatAsDateString(year, month, day, 23, 59, 59)

	return getArticles(siteCode, startDate, endDate)
}

/*
	Get articles within a date range
	Date strings should be formatted via gannettApi.FormatAsDateString
*/
func getArticles(siteCode string, startDate string, endDate string) []*m.SearchArticle {
	var articleResponse *m.ArticlesResponse = &m.ArticlesResponse{}
	var queryParams = GetDefaultSearchValues(siteCode)
	queryParams.Add("fq", fmt.Sprintf("lastpublished:[%s TO %s]", startDate, endDate))
	queryParams.Add("fq", "assettypename:text")
	queryParams.Add("fl", "initialpublished")
	queryParams.Add("sort", "lastpublished desc")

	url := fmt.Sprintf("%s?%s", GannettApiSearchRoot, queryParams.Encode())
	resp, err := http.Get(url)
	if err != nil {
		log.Warningf(`

		Failed to get articles for site %s. http.Get failed:

			Err: %v
		`, siteCode, err)
		return articleResponse.Results
	}
	defer resp.Body.Close()

	var body io.Reader = resp.Body
	if debugflag.IsEnabled("json:articles") {
		body = dumpJSONFromReader(fmt.Sprintf("Articles JSON from %s", url), "    ", resp.Body)
	}

	decoder := json.NewDecoder(body)
	err = decoder.Decode(articleResponse)

	if err != nil {
		log.Warningf(`

		Failed to get articles for site %s. Json decoding failed

			Err: %v
		`, siteCode, err)
		return []*m.SearchArticle{}
	}

	return articleResponse.Results
}

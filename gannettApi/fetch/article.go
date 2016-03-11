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
	Results         []ArticleIn
}

type ArticleIn struct {
	AssetId int
}

/*
	Get articles within an entire day
*/
func GetArticlesByDay(siteCode string, date time.Time) {
	year, monthObj, day := date.Date()
	month := int(monthObj)

	startDate := api.FormatAsDateString(year, month, day, 0, 0, 0)
	endDate := api.FormatAsDateString(year, month, day, 23, 59, 59)

	getArticles(siteCode, startDate, endDate)
}

/*
	Get articles within a date range
	Date strings should be formatted via gannettApi.FormatAsDateString
*/
func getArticles(siteCode string, startDate string, endDate string) {
	log.Info(startDate, endDate)
	var queryParams = api.GetDefaultValues(siteCode)
	queryParams.Add("fq", fmt.Sprintf("initialpublished:[%s TO %s]", startDate, endDate))
	queryParams.Add("fq", "assettypename:text")
	queryParams.Add("fl", "initialpublished")
	queryParams.Add("sort", "initialpublished desc")

	url := fmt.Sprintf("%s?%s", api.GannettApiUrlRoot, queryParams.Encode())
	log.Info(url)
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}

	articleResponse := &ArticlesResponse{}
	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(articleResponse)

	if err != nil {
		panic(err)
	}
	log.Info(articleResponse.Results)
}

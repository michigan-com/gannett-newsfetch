package gannettApi

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"

	m "github.com/michigan-com/gannett-newsfetch/model"
)

func GetBreakingNews(siteCode string, gannettSearchAPIKey string) []*m.SearchArticle {
	breakingResponse := &m.ArticlesResponse{}
	queryParams := GetDefaultSearchValues(siteCode, gannettSearchAPIKey)
	queryParams.Add("fl", "sitecode")
	queryParams.Add("fl", "expiration")
	queryParams.Add("fq", "expiration:[NOW TO *]")
	queryParams.Add("fq", "assettypename:breakingnewsweb")
	queryParams.Add("sort", "updated desc")

	url := fmt.Sprintf("%s?%s", GannettApiSearchRoot, queryParams.Encode())
	resp, err := http.Get(url)
	if err != nil {
		log.Warningf(`

      Failed to get breaking news for site %s. http.Get Failed

        Err: %v
    `, siteCode, err)
		return breakingResponse.Results
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(breakingResponse)
	if err != nil {
		log.Warningf(`

      Failed to get breaking news for site %s. Json decoding failed:

        Err: %v
    `, siteCode, err)
		return breakingResponse.Results
	}

	return breakingResponse.Results
}

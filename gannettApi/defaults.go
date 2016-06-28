package gannettApi

import (
	"fmt"
	"net/url"
)

// Use for querying for the list of articles
var GannettApiSearchRoot = "http://api.gannett-cdn.com/prod/Search/v4/assets/proxy"

// Use for getting the article content
var GannettApiPresentationRoot = "http://api.gannett-cdn.com/presentation/v4/assets"

/*
	Get default query param values
*/
func GetDefaultSearchValues(siteCode string, gannettSearchAPIKey string) url.Values {
	defaultValues := url.Values{}
	defaultValues.Set("q", "statusname:published")
	defaultValues.Set("fq", fmt.Sprintf("sitecode:%s", siteCode))
	defaultValues.Set("sc", siteCode)
	defaultValues.Set("apiKey", "newsfetch")
	defaultValues.Set("format", "json")
	defaultValues.Set("rows", "100")
	defaultValues.Set("api_key", gannettSearchAPIKey)

	return defaultValues
}

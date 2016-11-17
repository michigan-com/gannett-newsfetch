package lib

import (
	"net/url"
	"regexp"
	"strconv"
	"strings"
)

func GetArticleId(url string) int {
	// Given an article url, get the ID from it
	r := regexp.MustCompile("/([0-9]+)/{0,1}$")
	match := r.FindStringSubmatch(url)

	if len(match) <= 1 {
		return -1
	}

	i, err := strconv.Atoi(match[1])
	if err != nil {
		return -1
	}

	return i
}

/*
	Get the url host from the url string (inputUrl)

	Ex:
		result, err := GetHost("http://google.com")
		// result == "google"

	Using the url.Parse method, so urls must start with "http://"

*/
func GetHost(inputUrl string) (string, error) {
	u, err := url.Parse(inputUrl)
	if err != nil {
		return "", err
	}

	return strings.Replace(u.Host, "www.", "", 1), nil
}

func IsBlacklisted(url string) bool {
	blacklist := []string{
		"/videos/",
		"/police-blotter/",
		"/interactives/",
		"facebook.com",
		"/errors/404",
	}

	for _, item := range blacklist {
		if strings.Contains(url, item) {
			return true
		}
	}

	return false
}

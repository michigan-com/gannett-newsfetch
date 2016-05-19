package gannettApi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	log "github.com/Sirupsen/logrus"

	api "github.com/michigan-com/gannett-newsfetch/gannettApi"
)

type FullArticleIn struct {
	FullText string `json:"fullText"`
	StoryHighlights []string `json:"storyHighlights"`
}

func getArticleId(url string) int {
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

func GetArticleContent(articleUrl string) *FullArticleIn {
	var fullArticle *FullArticleIn = &FullArticleIn{}

	articleId := getArticleId(articleUrl)

	url := fmt.Sprintf("%s/%d?consumer=newsfetch&transform=full", api.GannettApiPresentationRoot, articleId)
	log.Info(url)
	resp, err := http.Get(url)
	if err != nil {
		log.Warningf(`

		Failed to get Article %d, http.Get() failed:

			Err: %v

		`, articleId, err)
		return fullArticle
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(fullArticle)
	if err != nil {
		log.Warningf(`

		Failed to get Article %d, json decoding failed:

			Err: %v
		`, articleId, err)
		return &FullArticleIn{}
	}

	return fullArticle
}

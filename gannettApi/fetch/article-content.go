package gannettApi

import (
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"

	api "github.com/michigan-com/gannett-newsfetch/gannettApi"
)

type FullArticleIn struct {
	FullText string `json:"fullText"`
}

func GetArticleContent(articleId int) *FullArticleIn {
	var fullArticle *FullArticleIn = &FullArticleIn{}

	url := fmt.Sprintf("%s/%d?consumer=newsfetch&transform=full", api.GannettApiPresentationRoot, articleId)
	resp, err := http.Get(url)
	if err != nil {
		log.Warningf(`

		Failed to get Article %d, http.Get() failed:

			Err: %v

		`, articleId, err)
		return fullArticle
	}

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

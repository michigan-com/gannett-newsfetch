package gannettApi

import (
	"encoding/json"
	"fmt"
	"net/http"

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
		panic(err)
	}

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(fullArticle)
	if err != nil {
		panic(err)
	}

	return fullArticle
}

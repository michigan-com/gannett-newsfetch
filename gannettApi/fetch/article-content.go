package gannettApi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type FullArticleIn struct {
	FullText string `json:"fullText"`
}

func GetArticleContent(articleId int) *FullArticleIn {
	var fullArticle *FullArticleIn

	url := fmt.Sprintf("%s%d?consumer=newsfetch&transform=full")
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

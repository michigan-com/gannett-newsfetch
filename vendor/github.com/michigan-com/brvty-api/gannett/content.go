package gannett

import (
	"encoding/json"
	"fmt"
	"net/http"

	parse "github.com/michigan-com/gannett-newsfetch/parse/body"
)

type ArticleContent struct {
	Headline string
	Text     string
}

type contentApiResponse struct {
	Title           string   `json:"title"`
	FullText        string   `json:"fullText"`
	StoryHighlights []string `json:"storyHighlights"`
}

func FetchArticleContent(id int) (*ArticleContent, error) {
	url := fmt.Sprintf("%s/%d?consumer=newsfetch&transform=full", presentationAPIEndpoint, id)
	fmt.Printf("Loading Gannett article %v from %v\n", id, url)
	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch article %v, http.Get() says: %v", id, err)
	}
	defer resp.Body.Close()

	var r contentApiResponse
	err = json.NewDecoder(resp.Body).Decode(&r)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch article %v, JSON decoding error: %v", id, err)
	}

	body := parse.ParseArticleBodyHtml(r.FullText)

	return &ArticleContent{
		Headline: r.Title,
		Text:     body,
	}, nil
}

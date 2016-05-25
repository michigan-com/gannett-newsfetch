package main_test

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	m "github.com/michigan-com/gannett-newsfetch/model"
)

func IntegrationTest(t *testing.T) {
	t.Fatalf("%v", "asdfasdf")
	testArticle := getTestArticle()
	if testArticle == nil {
		t.Fatalf("Test article json reading failed")
	}

	t.Log(testArticle.Headline)
}

func getTestArticle() *m.Article {
	file, err := ioutil.ReadFile("./testData/expectedArticle.json")
	if err != nil {
		fmt.Println("failed to read json article from filesystem")
		return nil
	}

	article := &m.Article{}
	err = json.Unmarshal(file, article)
	if err != nil {
		fmt.Println("failed to unmarshall json article into memory")
		return nil
	}
	return article
}

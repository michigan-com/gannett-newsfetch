package commands

import (
	"encoding/json"
	"fmt"
	"os/exec"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	log "github.com/Sirupsen/logrus"

	"github.com/michigan-com/gannett-newsfetch/lib"
	m "github.com/michigan-com/gannett-newsfetch/model"
)

type SummaryResponse struct {
	Skipped    int `json:"skipped"`
	Summarized int `json:"summarized"`
}

/*
	Run a python process to summarize all articles in the ToSummarize collection
*/
func ProcessSummaries(session *mgo.Session, toSummarize []interface{}, mongoUri string, summaryVEnv string) (*SummaryResponse, error) {
	summResp := &SummaryResponse{}

	bulk := session.DB("").C("ToSummarize").Bulk()
	bulk.Upsert(toSummarize...)
	_, err := bulk.Run()
	if err != nil {
		return summResp, err
	}

	if summaryVEnv == "" {
		return nil, fmt.Errorf("Missing SUMMARY_VENV environment variable, skipping summarizer")
	}

	cmd := fmt.Sprintf("%s/bin/python", summaryVEnv)
	pyScript := fmt.Sprintf("%s/bin/summary.py", summaryVEnv)

	log.Infof("Executing command: %s %s %s", cmd, pyScript, mongoUri)

	out, err := exec.Command(cmd, pyScript, mongoUri).Output()
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(out, summResp); err != nil {
		return nil, err
	}
	fmt.Println(summResp)

	return summResp, nil
}

/*
	We should summarize the article under two scenarios:

		1) This article does not yet exist in the database
		2) This article exists in the database, but the timestamp has been updated

	Does a lookup based on Article.ArticleId
*/
func shouldSummarizeArticle(article *m.SearchArticle, session *mgo.Session) bool {
	// Don't summarize if it's a blacklisted article
	if m.IsBlacklisted(article.Urls.LongUrl) {
		return false
	}

	var storedArticle *m.Article = &m.Article{}
	collection := session.DB("").C("Article")
	err := collection.Find(bson.M{"article_id": article.AssetId}).One(storedArticle)
	datePublished := lib.GannettDateStringToDate(article.DatePublished)
	if err == mgo.ErrNotFound {
		return true
	} else if !lib.SameTime(datePublished, storedArticle.Created_at) {
		return true
	} else if len(storedArticle.Summary) == 0 {
		return true
	}
	return false

}

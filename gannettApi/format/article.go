package gannettApi

import (
	"time"

	log "github.com/Sirupsen/logrus"
	f "github.com/michigan-com/gannett-newsfetch/gannettApi/fetch"
	"github.com/michigan-com/gannett-newsfetch/lib"
	m "github.com/michigan-com/gannett-newsfetch/model"
)

/*
	Given an article from the Gannett API (`inputArticle`), return an article that
	will be saved in mongo
*/
func FormatArticleForSaving(inputArticle *f.ArticleIn) *m.Article {
	var mongoArticle *m.Article = &m.Article{}

	mongoArticle.ArticleId = inputArticle.AssetId
	mongoArticle.Headline = inputArticle.Headline
	mongoArticle.Subheadline = inputArticle.PromoBrief
	mongoArticle.Section = inputArticle.Ssts.Section
	mongoArticle.Subsection = inputArticle.Ssts.SubSection
	mongoArticle.Sections = getAllSections(inputArticle)
	mongoArticle.Source, _ = lib.GetHost(inputArticle.Urls.LongUrl)
	mongoArticle.Timestamp = getDate(inputArticle.SolrFields.InitalPublished)
	mongoArticle.Url = inputArticle.Urls.LongUrl
	mongoArticle.ShortUrl = inputArticle.Urls.ShortUrl
	mongoArticle.Photo = getPhoto(inputArticle)
	mongoArticle.Created_at = getDate(inputArticle.DatePublished)

	return mongoArticle
}

func getPhoto(article *f.ArticleIn) *m.Photo {
	var photo *m.Photo = &m.Photo{}
	if article.Photo.AbsoluteUrl == "" {
		return nil
	}

	fullPhoto := m.PhotoInfo{
		Url:    article.Photo.Crops["1_1"],
		Width:  article.Photo.OriginalWidth,
		Height: article.Photo.OriginalHeight,
	}
	thumbPhoto := m.PhotoInfo{
		Url:    article.Photo.Crops["front_thumb"],
		Width:  article.Photo.OriginalWidth,
		Height: article.Photo.OriginalHeight,
	}

	photo.Caption = article.Photo.Caption
	photo.Credit = article.Photo.Credit
	photo.Full = fullPhoto
	photo.Thumbnail = thumbPhoto

	return photo
}

/*
	Given a string date, return the date. If anything goes wrong, return time.Now()
*/
func getDate(dateString string) time.Time {
	// https://golang.org/src/time/format.go
	// Idk, a regular date string wasnt working, cause why would it
	date, err := time.Parse(time.RFC3339Nano, dateString)
	if err != nil {
		log.Info(err)
		return time.Now()
	}
	return date
}

func getAllSections(article *f.ArticleIn) []string {
	sections := make([]string, 0, 3)
	sections = append(sections, article.Ssts.Section)
	if article.Ssts.SubSection != "" {
		sections = append(sections, article.Ssts.SubSection)
	}
	if article.Ssts.Topic != "" {
		sections = append(sections, article.Ssts.Topic)
	}
	if article.Ssts.SubTopic != "" {
		sections = append(sections, article.Ssts.SubTopic)
	}

	return sections
}

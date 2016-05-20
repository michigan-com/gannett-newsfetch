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
func FormatSearchArticleForSaving(inputArticle *f.ArticleIn) *m.Article {
	var mongoArticle *m.Article = &m.Article{}

	mongoArticle.ArticleId = inputArticle.AssetId
	mongoArticle.Headline = inputArticle.Headline
	mongoArticle.Subheadline = inputArticle.PromoBrief
	mongoArticle.Section = inputArticle.Ssts.Section
	mongoArticle.Subsection = inputArticle.Ssts.SubSection
	mongoArticle.Sections = getAllSections(&inputArticle.Ssts)
	mongoArticle.Source, _ = lib.GetHost(inputArticle.Urls.LongUrl)
	mongoArticle.Timestamp = getDate(inputArticle.SolrFields.InitalPublished)
	mongoArticle.Url = inputArticle.Urls.LongUrl
	mongoArticle.ShortUrl = inputArticle.Urls.ShortUrl
	mongoArticle.Photo = getPhoto(&inputArticle.Photo)
	mongoArticle.Created_at = getDate(inputArticle.DatePublished)

	return mongoArticle
}

func FormatAssetArticleForSaving(inputArticle *f.AssetArticleIn, inputPhoto *f.PhotoInfo) *m.Article {
	var mongoArticle *m.Article = &m.Article{}

	mongoArticle.ArticleId = inputArticle.AssetId
	mongoArticle.Headline = inputArticle.Headline
	mongoArticle.Subheadline = inputArticle.PromoBrief
	mongoArticle.Section = inputArticle.Ssts.Section
	mongoArticle.Subsection = inputArticle.Ssts.SubSection
	mongoArticle.Sections = getAllSections(&inputArticle.Ssts)
	mongoArticle.Source, _ = lib.GetHost(inputArticle.Links.LongUrl.Href)
	mongoArticle.Timestamp = getDate(inputArticle.InitialPublishDate)
	mongoArticle.Url = inputArticle.Links.LongUrl.Href
	mongoArticle.ShortUrl = inputArticle.Links.ShortUrl.Href
	mongoArticle.Photo = getPhoto(inputPhoto)
	mongoArticle.Created_at = getDate(inputArticle.PublishDate)
	mongoArticle.Body = inputArticle.FullText
	mongoArticle.StoryHighlights = inputArticle.StoryHighlights

	return mongoArticle
}

func getPhoto(inputPhoto *f.PhotoInfo) *m.Photo {
	var photo *m.Photo = &m.Photo{}
	if inputPhoto == nil || inputPhoto.AbsoluteUrl == "" {
		return nil
	}

	photo.Caption = inputPhoto.Caption
	photo.Credit = inputPhoto.Credit

	photo.Full = m.PhotoInfo{
		Url:    inputPhoto.AbsoluteUrl,
		Width:  inputPhoto.OriginalWidth,
		Height: inputPhoto.OriginalHeight,
	}
	photo.Thumbnail = m.PhotoInfo{
		Url:    inputPhoto.Crops["front_thumb"],
		Width:  inputPhoto.OriginalWidth,
		Height: inputPhoto.OriginalHeight,
	}

	photo.Crops = make(map[string]m.PhotoInfo, len(inputPhoto.Crops))
	for size, crop := range inputPhoto.Crops {
		photo.Crops[size] = m.PhotoInfo{
			Url:    crop,
			Width:  0,
			Height: 0,
		}
	}

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
		log.Infof("PARSING ERROR article: %v", err)
		return time.Now()
	}
	return date.UTC()
}

func getAllSections(ssts *f.Ssts) []string {
	sections := make([]string, 0, 3)
	sections = append(sections, ssts.Section)
	if ssts.SubSection != "" {
		sections = append(sections, ssts.SubSection)
	}
	if ssts.Topic != "" {
		sections = append(sections, ssts.Topic)
	}
	if ssts.SubTopic != "" {
		sections = append(sections, ssts.SubTopic)
	}

	return sections
}

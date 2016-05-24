package gannettApi

import (
	"fmt"

	"github.com/michigan-com/gannett-newsfetch/lib"
	m "github.com/michigan-com/gannett-newsfetch/model"
	"github.com/michigan-com/gannett-newsfetch/parse/body"
)

/*
	Given an article from the Gannett API (`assetArticle`), return an article that
	will be saved in mongo
*/
func FormatAssetArticleForSaving(assetArticle *AssetArticle, inputPhoto *PhotoInfo) *m.Article {
	var mongoArticle *m.Article = &m.Article{}

	mongoArticle.ArticleId = assetArticle.AssetId
	mongoArticle.Headline = assetArticle.Headline
	mongoArticle.Subheadline = assetArticle.PromoBrief
	mongoArticle.Section = assetArticle.Ssts.Section
	mongoArticle.Subsection = assetArticle.Ssts.SubSection
	mongoArticle.Sections = getAllSections(&assetArticle.Ssts)
	mongoArticle.Source, _ = lib.GetHost(assetArticle.Links.LongUrl.Href)
	mongoArticle.Timestamp = lib.GannettDateStringToDate(assetArticle.InitialPublishDate)
	mongoArticle.Url = assetArticle.Links.LongUrl.Href
	mongoArticle.ShortUrl = assetArticle.Links.ShortUrl.Href
	mongoArticle.Photo = getPhoto(inputPhoto)
	mongoArticle.Created_at = lib.GannettDateStringToDate(assetArticle.PublishDate)
	mongoArticle.Body = parse.ParseArticleBodyHtml(assetArticle.FullText)
	mongoArticle.StoryHighlights = assetArticle.StoryHighlights

	return mongoArticle
}

func getPhoto(inputPhoto *PhotoInfo) *m.Photo {
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
	photo.Thumbnail = inputPhoto.Crops["front_thumb"]
	photo.Small = inputPhoto.Crops["small"]

	photo.Crops = inputPhoto.Crops
	return photo
}

func getAllSections(ssts *Ssts) []string {
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

/*
	Format a year, month, day, year, hours, minutes, and seconds into a date string
	for querying the Gannett Api

	FormatAsDateSting(2014, 10, 1, 0, 0, 0) == 2014-10-01T00:00:00Z

	For more info
		https://confluence.gannett.com/pages/viewpage.action?title=Search+v4+Recipes&spaceKey=GDPDW#Searchv4Recipes-FilterbyDateRange
*/
func FormatAsDateString(year, month, day, hour, minute, second int) string {
	return fmt.Sprintf("%02d-%02d-%02dT%02d:%02d:%02dZ", year, month, day, hour, minute, second)
}

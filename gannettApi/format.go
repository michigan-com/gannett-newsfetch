package gannettApi

import (
	"fmt"
	"strconv"

	"github.com/michigan-com/gannett-newsfetch/lib"
	m "github.com/michigan-com/gannett-newsfetch/model"
	"github.com/michigan-com/gannett-newsfetch/parse/body"
)

/*
	Given an article from the Gannett API (`assetArticle`), return an article that
	will be saved in mongo
*/
func FormatAssetArticleForSaving(assetArticleContent *m.AssetArticleContent) *m.Article {
	var mongoArticle *m.Article = &m.Article{}
	article := assetArticleContent.Article
	photo := assetArticleContent.Assets.Photo
	video := assetArticleContent.Assets.Video

	mongoArticle.ArticleId = article.AssetId
	mongoArticle.Headline = article.Headline
	mongoArticle.Subheadline = article.PromoBrief
	mongoArticle.Section = article.Ssts.Section
	mongoArticle.Subsection = article.Ssts.SubSection
	mongoArticle.Sections = getAllSections(&article.Ssts)
	mongoArticle.Domain, _ = lib.GetHost(article.Links.LongUrl.Href)
	mongoArticle.Timestamp = lib.GannettDateStringToDate(article.InitialPublishDate)
	mongoArticle.Url = article.Links.LongUrl.Href
	mongoArticle.ShortUrl = article.Links.ShortUrl.Href
	mongoArticle.Photo = formatPhoto(photo)
	mongoArticle.Created_at = lib.GannettDateStringToDate(article.PublishDate)
	mongoArticle.Body = parse.ParseArticleBodyHtml(article.FullText)
	mongoArticle.StoryHighlights = article.StoryHighlights
	mongoArticle.Video = video

	return mongoArticle
}

func formatPhoto(assetPhoto *m.AssetPhoto) *m.Photo {
	photo := &m.Photo{}

	if assetPhoto == nil || assetPhoto.AbsoluteUrl == "" {
		return nil
	}

	photo.Caption = assetPhoto.Caption
	photo.Credit = assetPhoto.Credit

	originalWidth, _ := strconv.Atoi(assetPhoto.OriginalWidth)
	originalHeight, _ := strconv.Atoi(assetPhoto.OriginalHeight)

	photo.Full = m.PhotoInfo{
		Url:    assetPhoto.AbsoluteUrl,
		Width:  originalWidth,
		Height: originalHeight,
	}
	photo.Thumbnail = photo.Full

	smallUrl := fmt.Sprintf("%s%s", assetPhoto.PublishUrl, assetPhoto.Attributes.SmallBaseName)
	smallWidth, _ := strconv.Atoi(assetPhoto.Attributes.SImageWidth)
	smallHeight, _ := strconv.Atoi(assetPhoto.Attributes.SImageHeight)
	photo.Small = m.PhotoInfo{
		Url:    smallUrl,
		Width:  smallWidth,
		Height: smallHeight,
	}

	crops := make(map[string]m.PhotoInfo)
	for _, crop := range assetPhoto.Crops {
		crops[crop.Name] = m.PhotoInfo{
			Url:    crop.Path,
			Width:  crop.Width,
			Height: crop.Height,
		}

		if crop.Name == "front_thumb" {
			photo.Thumbnail = m.PhotoInfo{
				Url:    crop.Path,
				Width:  crop.Width,
				Height: crop.Height,
			}
		}
	}
	photo.Crops = crops
	return photo
}

func getAllSections(ssts *m.Ssts) []string {
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

package gannettApi

import "testing"

type ArticleIdTestCase struct {
	Url string
	Id  int
}

func TestArticleScrapingNoPhotoOrVideo(t *testing.T) {
	articleIds := []int{
		76033594,
	}

	for _, articleId := range articleIds {
		assetArticleContent := GetAssetArticleContent(articleId)

		if assetArticleContent.Assets.Photo != nil {
			t.Fatalf("Article %d should have Photo == nil", articleId)
		}

		if assetArticleContent.Assets.Video != nil {
			t.Fatalf("Article %d should have Video == nil", articleId)
		}
	}
}

func TestArticleScrapingPhotoNoVideo(t *testing.T) {
	articleIds := []int{
		85015624, 85059214,
	}

	for _, articleId := range articleIds {
		assetArticleContent := GetAssetArticleContent(articleId)

		if assetArticleContent.Assets.Photo == nil {
			t.Fatalf("Article %d should have a photo", articleId)
		}

		if assetArticleContent.Assets.Video != nil {
			t.Fatalf("Article %d should not have a video", articleId)
		}
	}
}

func TestArticleScarpingPhotoAndVideo(t *testing.T) {
	articleId := 84913242
	assetArticleContent := GetAssetArticleContent(articleId)

	if assetArticleContent.Assets.Photo == nil {
		t.Fatalf("Article %d should have a photo", articleId)
	}

	if assetArticleContent.Assets.Video == nil {
		t.Fatalf("Article %d should have a video", articleId)
	}
}

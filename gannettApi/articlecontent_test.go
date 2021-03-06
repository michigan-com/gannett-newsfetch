package gannettApi

import (
	"os"
	"testing"
)

type ArticleIdTestCase struct {
	Url string
	Id  int
}

func getApiKey(t *testing.T) string {
	key := os.Getenv("GANNETT_ASSET_API_KEY")
	if key == "" {
		t.Fatalf("GANNETT_ASSET_API_KEY is required")
	}
	return key
}

func TestArticleScrapingNoPhotoOrVideo(t *testing.T) {
	articleIds := []int{
		76033594,
	}

	for _, articleId := range articleIds {
		assetArticleContent := GetAssetArticleContent(articleId, getApiKey(t))

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
		assetArticleContent := GetAssetArticleContent(articleId, getApiKey(t))

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
	assetArticleContent := GetAssetArticleContent(articleId, getApiKey(t))

	if assetArticleContent.Assets.Photo == nil {
		t.Fatalf("Article %d should have a photo", articleId)
	}

	if assetArticleContent.Assets.Video == nil {
		t.Fatalf("Article %d should have a video", articleId)
	}
}

package gannettApi

import "testing"

type ArticleIdTestCase struct {
	Url string
	Id  int
}

func TestArticleId(t *testing.T) {
	testCases := []ArticleIdTestCase{
		ArticleIdTestCase{
			"http://www.freep.com/story/news/local/michigan/detroit/2016/05/23/hillary-clinton-calls-trump-disaster-waiting-happen/84788222/",
			84788222,
		},
	}

	for _, testCase := range testCases {
		articleIdTestCase(t, testCase)
	}
}

func articleIdTestCase(t *testing.T, testCase ArticleIdTestCase) {
	articleId := getArticleId(testCase.Url)
	if articleId != testCase.Id {
		t.Fatalf("Expected article ID: %d, actual article ID: %d", testCase.Id, articleId)
	}
}

func TestArticleScrapingNoPhotoOrVideo(t *testing.T) {
	noPhotoArticleId := 76033594
	assetArticleContent := GetAssetArticleContent(noPhotoArticleId)

	if assetArticleContent.Assets.Photo != nil {
		t.Fatalf("Article %d should have Photo == nil", noPhotoArticleId)
	}

	if assetArticleContent.Assets.Video != nil {
		t.Fatalf("Article %d should have Video == nil", noPhotoArticleId)
	}
}

func TestArticleScrapingPhotoNoVideo(t *testing.T) {
	noVideoArticleId := 85059214
	assetArticleContent := GetAssetArticleContent(noVideoArticleId)

	if assetArticleContent.Assets.Photo == nil {
		t.Fatalf("Article %d should have a photo", noVideoArticleId)
	}

	if assetArticleContent.Assets.Video != nil {
		t.Fatalf("Article %d should not have a video", noVideoArticleId)
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

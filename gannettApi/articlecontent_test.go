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

func ArticleContentFetchTest(t *testing.T) {
	// articleId := 84849924
	// expectedAssetArticle := AssetArticle{
	// 	AssetId:  articleId,
	// 	Headline: "Movement afterparties: 10 top events during fest weekend",
	// 	Ssts: Ssts{
	// 		Section:    "entertainment",
	// 		SubSection: "music",
	// 		Topic:      "",
	// 		SubTopic:   "",
	// 	},
	// 	Links: links{
	// 		LongUrl: AssetUrl{
	// 			Href: "http://www.freep.com/story/entertainment/music/2016/05/24/movement-festivals-afterparties-detroit-electronic-music/84849924/",
	// 		},
	// 		ShortUrl: AssetUrl{
	// 			Href: "http://on.freep.com/1TC0VSH",
	// 		},
	// 		Photo: &AssetPhoto{
	// 			Id: 80766590,
	// 		},
	// 	},
	// 	PublishDate:        "2016-05-24T18:27:12.213Z",
	// 	InitialPublishDate: "2016-05-24T18:27:12.213Z",
	// 	PromoBrief:         "Detroit will be pulsating with electronic music and dance parties throughout Memorial Day Weekend",
	// 	Attribution: AssetAttribution{
	// 		Author: "Tamara Warren",
	// 	},
	// }
}

package gannettApi

import "testing"

type ArticleIdTestCase struct {
	Url string
	Id  int
}

func ArticleIdTest(t *testing.T) {
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

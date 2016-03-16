package parse

import (
	"bytes"
	gq "github.com/PuerkitoBio/goquery"
	"strings"

	"github.com/michigan-com/gannett-newsfetch/parse/classify"
	"github.com/michigan-com/gannett-newsfetch/parse/dateline"
)

func ParseArticleBodyHtml(bodyHtml string) string {
	doc, err := gq.NewDocumentFromReader(bytes.NewBufferString(bodyHtml))
	if err != nil {
		return ""
	}

	paragraphs := doc.Find("p")

	ignoreRemaining := false
	paragraphStrings := paragraphs.Map(func(i int, paragraph *gq.Selection) string {
		if ignoreRemaining {
			return ""
		}
		for _, selector := range [...]string{"span.-newsgate-character-cci-tagline-name-", "span.-newsgate-paragraph-cci-infobox-head-"} {
			if el := paragraph.Find(selector); el.Length() > 0 {
				ignoreRemaining = true
				return ""
			}
		}

		text := strings.TrimSpace(paragraph.Text())

		if worthy, _ := classify.IsWorthyParagraph(text); !worthy {
			return ""
		}

		for _, selector := range [...]string{"span.-newsgate-paragraph-cci-subhead-lead-", "span.-newsgate-paragraph-cci-subhead-"} {
			if el := paragraph.Find(selector); el.Length() > 0 {
				return ""
			}
		}

		return text
	})

	if len(paragraphStrings) > 0 {
		paragraphStrings[0] = dateline.RmDateline(paragraphStrings[0])
	}

	body := strings.Join(paragraphStrings, "\n")
	// TODO
	//recipeData, recipeMsg := recipe_parsing.ExtractRecipes(doc)
	return body
}

package gannettApi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	log "github.com/Sirupsen/logrus"

	api "github.com/michigan-com/gannett-newsfetch/gannettApi"
)

type FullArticleIn struct {
	FullText        string   `json:"fullText"`
	StoryHighlights []string `json:"storyHighlights"`
}

func getArticleId(url string) int {
	// Given an article url, get the ID from it
	r := regexp.MustCompile("/([0-9]+)/{0,1}$")
	match := r.FindStringSubmatch(url)

	if len(match) <= 1 {
		return -1
	}

	i, err := strconv.Atoi(match[1])
	if err != nil {
		return -1
	}

	return i
}

func getAssetUrl(assetId int) string {
	return fmt.Sprintf("%s/%d?consumer=newsfetch&transform=full", api.GannettApiPresentationRoot, assetId)
}

func GetAssetArticleAndPhoto(articleId int) (*AssetArticleIn, *PhotoInfo) {
	var fullArticle *AssetArticleIn = &AssetArticleIn{}
	var photo = &PhotoInfo{}
	var assetPhoto = &AssetPhotoInfo{}

	url := getAssetUrl(articleId)
	resp, err := http.Get(url)
	if err != nil {
		log.Warningf(`

		Failed to get Article %d, http.Get() failed:

			Err: %v

		`, articleId, err)
		return fullArticle, photo
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(fullArticle)
	if err != nil {
		log.Warningf(`

		Failed to get Article %d, json deconding failed:

			Err: %v

		`, articleId, err)
		return fullArticle, photo
	}

	// Fetch the photo from the asset api (we only get a pointer to the photo)
	if fullArticle.Links.Photo != nil {
		photoAssetId := fullArticle.Links.Photo.Id
		url := getAssetUrl(photoAssetId)
		photoResp, err := http.Get(url)
		if err != nil {
			log.Warningf(`

			Failed to get Photo %d, http.Get() failed:

				Err: %v

			`, photoAssetId, err)
			return fullArticle, photo
		}
		defer photoResp.Body.Close()

		decoder = json.NewDecoder(photoResp.Body)
		err = decoder.Decode(assetPhoto)
		if err != nil {
			log.Warningf(`

			Failed to get Photo %d, json deconding failed:

				Err: %v

			`, photoAssetId, err)
			return fullArticle, photo
		}

		photo.AbsoluteUrl = assetPhoto.AbsoluteUrl
		photo.Caption = assetPhoto.Caption
		photo.Credit = assetPhoto.Credit
		photo.OriginalHeight, _ = strconv.Atoi(assetPhoto.OriginalHeight)
		photo.OriginalWidth, _ = strconv.Atoi(assetPhoto.OriginalWidth)

		crops := make(map[string]string)
		for _, crop := range assetPhoto.Crops {
			crops[crop.Name] = crop.Path
		}
		photo.Crops = crops
	} else {
		photo = nil
	}

	return fullArticle, photo
}

func GetArticleBody(articleUrl string) *FullArticleIn {
	var fullArticle *FullArticleIn = &FullArticleIn{}

	articleId := getArticleId(articleUrl)

	url := getAssetUrl(articleId)
	resp, err := http.Get(url)
	if err != nil {
		log.Warningf(`

		Failed to get Article %d, http.Get() failed:

			Err: %v

		`, articleId, err)
		return fullArticle
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(fullArticle)
	if err != nil {
		log.Warningf(`

		Failed to get Article %d, json decoding failed:

			Err: %v
		`, articleId, err)
		return &FullArticleIn{}
	}

	return fullArticle
}

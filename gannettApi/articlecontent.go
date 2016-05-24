package gannettApi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	log "github.com/Sirupsen/logrus"

	m "github.com/michigan-com/gannett-newsfetch/model"
)

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
	return fmt.Sprintf("%s/%d?consumer=newsfetch&transform=full", GannettApiPresentationRoot, assetId)
}

func GetAssetArticleAndPhoto(articleId int) (*AssetArticle, *PhotoInfo) {
	fullArticle := &AssetArticle{}
	photo := &PhotoInfo{}

	fullArticle = GetAssetArticle(articleId)

	// Fetch the photo from the asset api (we only get a pointer to the photo)
	if fullArticle.Links.Photo != nil {
		photoAssetId := fullArticle.Links.Photo.Id
		photo = GetAssetPhoto(photoAssetId)
	} else {
		photo = nil
	}

	return fullArticle, photo
}

func GetAssetArticle(articleId int) *AssetArticle {
	var fullArticle *AssetArticle = &AssetArticle{}

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

		Failed to get Article %d, json deconding failed:

			Err: %v

		`, articleId, err)
		return fullArticle
	}

	return fullArticle
}

func GetAssetPhoto(photoAssetId int) *PhotoInfo {
	photo := &PhotoInfo{}
	assetPhotoInfo := &AssetPhotoInfo{}

	url := getAssetUrl(photoAssetId)
	photoResp, err := http.Get(url)
	fmt.Println(url)
	if err != nil {
		log.Warningf(`

			Failed to get Photo %d, http.Get() failed:

				Err: %v

			`, photoAssetId, err)
		return nil
	}
	defer photoResp.Body.Close()

	decoder := json.NewDecoder(photoResp.Body)
	err = decoder.Decode(assetPhotoInfo)
	if err != nil {
		log.Warningf(`

			Failed to get Photo %d, json deconding failed:

				Err: %v

			`, photoAssetId, err)
		return nil
	}

	photo.AbsoluteUrl = assetPhotoInfo.AbsoluteUrl
	photo.Caption = assetPhotoInfo.Caption
	photo.Credit = assetPhotoInfo.Credit
	photo.OriginalHeight, _ = strconv.Atoi(assetPhotoInfo.OriginalHeight)
	photo.OriginalWidth, _ = strconv.Atoi(assetPhotoInfo.OriginalWidth)

	crops := make(map[string]m.PhotoInfo)
	for _, crop := range assetPhotoInfo.Crops {
		crops[crop.Name] = m.PhotoInfo{
			Url:    crop.Path,
			Width:  crop.Width,
			Height: crop.Height,
		}
	}

	fmt.Printf("%v", assetPhotoInfo.Attributes)
	smallUrl := fmt.Sprintf("%s%s", assetPhotoInfo.PublishUrl, assetPhotoInfo.Attributes.SmallBaseName)
	smallWidth, _ := strconv.Atoi(assetPhotoInfo.Attributes.SImageWidth)
	smallHeight, _ := strconv.Atoi(assetPhotoInfo.Attributes.SImageHeight)
	crops["small"] = m.PhotoInfo{
		Url:    smallUrl,
		Width:  smallWidth,
		Height: smallHeight,
	}

	photo.Crops = crops

	return photo
}

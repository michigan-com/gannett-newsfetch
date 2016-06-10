package gannettApi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	log "github.com/Sirupsen/logrus"

	m "github.com/michigan-com/gannett-newsfetch/model"
)

func getAssetUrl(assetId int) string {
	return fmt.Sprintf("%s/%d?consumer=newsfetch&transform=full", GannettApiPresentationRoot, assetId)
}

func GetAssetArticleContent(articleId int) *m.AssetArticleContent {
	assetArticle := &m.AssetArticle{}

	GetAsset(articleId, assetArticle)

	assetArticleContent := &m.AssetArticleContent{
		Article: assetArticle,
		Assets:  GetArticleAssets(assetArticle.Links.Assets),
	}

	if assetArticle.Links.Photo != nil {
		assetPhoto := &m.AssetPhoto{}
		err := GetAsset(assetArticle.Links.Photo.Id, assetPhoto)
		if err == nil {
			assetArticleContent.Assets.Photo = assetPhoto
		}
	}

	return assetArticleContent
}

/** Get Photos, videos, and (TODO) galleries stored as IDs in an article's metadata */
func GetArticleAssets(assets []*m.GannettApiAsset) *m.ArticleAssets {
	articleAssets := &m.ArticleAssets{}
	var assetWait sync.WaitGroup

	for _, asset := range assets {

		assetWait.Add(1)
		go func(asset *m.GannettApiAsset) {
			defer assetWait.Done()

			if asset.Type == "video" {
				assetVideo := &m.AssetVideo{}
				err := GetAsset(asset.Id, assetVideo)
				if err == nil {
					articleAssets.Video = assetVideo
				}
			}
		}(asset)
	}
	assetWait.Wait()

	return articleAssets
}

func GetAsset(assetId int, assetResp m.AssetResp) error {
	url := getAssetUrl(assetId)
	resp, err := http.Get(url)
	if err != nil {
		log.Warning(`

			Failed to get asset %d, http.Get() failed

				Err: %v

		`, assetId, err)
		return err
	}
	defer resp.Body.Close()

	decoder := json.NewDecoder(resp.Body)
	err = assetResp.Decode(decoder)
	if err != nil {
		log.Warningf(`

			Failed to decode asset %d:

				Err: %v
		`, assetId, err)
		return err
	}

	return nil
}

package brvtycore

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"

	"github.com/michigan-com/brvty-api/canonicalurl"
)

func computeID(url string) string {
	url, _ = canonicalurl.CleanURLString(url, canonicalurl.Short)
	bytes := sha1.Sum([]byte(url))
	return hex.EncodeToString(bytes[:])
}

func computeIDsForURLs(urls []string) []string {
	ids := make([]string, 0, len(urls))
	for _, url := range urls {
		id := computeID(url)
		ids = append(ids, id)
	}
	return ids
}

func computeRev(payload ResourcePayload) string {
	bytes := sha1.Sum([]byte(fmt.Sprintf("%s\n%s\n%s", payload.Headline, payload.URL, payload.Text)))
	return hex.EncodeToString(bytes[:])
}

func (b ResourceBody) Hash() string {
	bytes := sha1.Sum([]byte(fmt.Sprintf("%s\n%s", b.Headline, b.Text)))
	return hex.EncodeToString(bytes[:])
}

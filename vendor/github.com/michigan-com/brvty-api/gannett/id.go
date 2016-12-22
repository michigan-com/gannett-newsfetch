package gannett

import (
	"net/url"
	"regexp"
	"strconv"
)

var idRegexp = regexp.MustCompile("/([0-9]+)/{0,1}$")

const IDNotFound = -1

func FindArticleID(rawurl string) int {
	u, err := url.Parse(rawurl)
	if err != nil {
		return IDNotFound
	}

	if !IsGannettHost(u.Host) {
		return IDNotFound
	}

	match := idRegexp.FindStringSubmatch(u.Path)
	if len(match) <= 1 {
		return IDNotFound
	}

	id, err := strconv.Atoi(match[1])
	if err != nil {
		return IDNotFound
	}

	return id
}

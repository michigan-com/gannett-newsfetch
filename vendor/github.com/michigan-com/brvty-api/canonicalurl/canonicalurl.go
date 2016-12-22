package canonicalurl

import (
	"net"
	"net/url"
	"strconv"
	"strings"
)

type Mode int

const (
	StripQuery Mode = (1 << iota)
	StripHTTPScheme
	StripWWW
	OmitLeadingSlashSlashWhenUnsafe
	OmitSoleSlashPath
)

const (
	Safe  Mode = 0
	Short Mode = StripQuery | StripHTTPScheme | StripWWW | OmitLeadingSlashSlashWhenUnsafe | OmitSoleSlashPath
)

// CleanURLString cleans up the given real-world and/or user-provided URL, producing a valid URL that hopefully points to the same resource. We aim to accept the URLs that browsers would normally accept.
//
// We're ignoring the deprecated hash-bang (#!) spec and assume that it's always safe to strip the fragment.
func CleanURLString(rawurl string, mode Mode) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return rawurl, err
	}

	if u.Scheme == "" && u.Host == "" && u.Opaque == "" {
		fixedURL := "http://" + rawurl
		u, err = url.Parse(fixedURL)
		if err != nil {
			return rawurl, err
		}
	} else if u.Opaque != "" && strings.Contains(u.Scheme, ".") {
		if _, err = strconv.Atoi(u.Opaque); err == nil {
			fixedURL := "http://" + rawurl
			u, err = url.Parse(fixedURL)
			if err != nil {
				return rawurl, err
			}
		}
	}

	u.User = nil
	u.Fragment = ""

	if (mode & StripQuery) != 0 {
		u.RawQuery = ""
	} else {
		u.RawQuery = CleanQueryParams(u.Query()).Encode()
	}

	u.Host = StripDefaultPort(u.Host, u.Scheme)

	safeToStripScheme := isSafeToStripSchemeForHost(u.Host)

	if safeToStripScheme && (mode&StripHTTPScheme) != 0 {
		if u.Scheme == "http" || u.Scheme == "https" {
			u.Scheme = ""
		}
	}

	if (mode & StripWWW) != 0 {
		u.Host = stripPrefix(u.Host, "www.")
	}

	resultStr := u.String()

	if safeToStripScheme || (mode&OmitLeadingSlashSlashWhenUnsafe) != 0 {
		resultStr = stripPrefix(resultStr, "//")
	}

	if u.Scheme == "" && u.Path == "/" && u.RawQuery == "" && (mode&OmitSoleSlashPath) != 0 {
		resultStr = stripSuffix(resultStr, "/")
	}

	return resultStr, nil
}

// We cannot strip the scheme if there's port specified, because url.Parse would parse
// "example.com:8080" as having a scheme of "example.com".
func isSafeToStripSchemeForHost(host string) bool {
	if !strings.Contains(host, ":") {
		return true
	}
	host, port, err := net.SplitHostPort(host)
	return (err == nil) && (port == "")
}

func ToPrettyString(rawurl string) (string, error) {
	u, err := url.Parse(rawurl)
	if err != nil {
		return rawurl, err
	}

	u.User = nil
	u.Fragment = ""
	u.RawQuery = ""

	if u.Scheme == "https" || u.Scheme == "http" {
		u.Scheme = ""
	}

	id := u.String()
	if strings.HasPrefix(id, "//") {
		id = strings.Replace(id, "//", "", 1)
	}

	return id, nil
}

var knownStupidParams = map[string]bool{
	"utm_source":   true,
	"utm_medium":   true,
	"utm_campaign": true,
	"utm_term":     true,
	"utm_content":  true,
}

// StripDefaultPort removes the default port for the given scheme. Currently only handles http and https.
func StripDefaultPort(host, scheme string) string {
	if scheme == "http" {
		return stripSuffix(host, ":80")
	} else if scheme == "https" {
		return stripSuffix(host, ":443")
	} else {
		return host
	}
}

// CleanQueryParams strips those query parameters that should be safe to remove.
func CleanQueryParams(source url.Values) url.Values {
	result := url.Values{}
	for name, values := range source {
		if knownStupidParams[name] {
			continue
		}
		result[name] = values
	}
	return result
}

func stripSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		return s[:len(s)-len(suffix)]
	} else {
		return s
	}
}

func stripPrefix(s, suffix string) string {
	if strings.HasPrefix(s, suffix) {
		return s[len(suffix):]
	} else {
		return s
	}
}

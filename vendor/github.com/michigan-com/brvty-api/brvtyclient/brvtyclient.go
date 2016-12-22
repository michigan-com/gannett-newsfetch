package brvtyclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/michigan-com/brvty-api/brvtycore"
)

type Client struct {
	baseURL string
	apiKey  string
}

func New(baseURL string, apiKey string) *Client {
	return &Client{baseURL, apiKey}
}

func (c *Client) Add(urls []string, timeout time.Duration) ([]*brvtycore.Resource, error) {
	var result []*brvtycore.Resource
	err := c.perform("POST", "/r/", url.Values{
		"urls": []string{strings.Join(urls, "\n")},
		"wait": []string{strconv.Itoa(int(timeout / time.Millisecond))},
	}, &result)
	return result, err
}

func (c *Client) Get(urls []string, timeout time.Duration) ([]*brvtycore.Resource, error) {
	var result []*brvtycore.Resource
	err := c.perform("GET", "/r/", url.Values{
		"urls": []string{strings.Join(urls, "\n")},
		"wait": []string{strconv.Itoa(int(timeout / time.Millisecond))},
	}, &result)
	return result, err
}

func (c *Client) UpdateBody(urlstr string, strategy string, headline, body string) error {
	return c.perform("POST", fmt.Sprintf("/r/body/%s", strategy), url.Values{
		"url":      []string{urlstr},
		"headline": []string{headline},
		"body":     []string{body},
	}, nil)
}

func (c *Client) UpdateSummary(urlstr string, strategy string, sentences []string) error {
	return c.perform("POST", fmt.Sprintf("/r/summary/%s", strategy), url.Values{
		"url":       []string{urlstr},
		"sentences": []string{strings.Join(sentences, "\n")},
	}, nil)
}

func (c *Client) perform(method, path string, values url.Values, result interface{}) error {
	req, err := makeRequest(method, c.baseURL, path, values)
	if err != nil {
		return err
	}

	req.Header.Set("X-Brvty-Api-Key", c.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	mediaType := resp.Header.Get("Content-Type")
	ctype, _, err := mime.ParseMediaType(mediaType)
	if err != nil {
		return fmt.Errorf("cannot parse Content-Type string %+v", mediaType)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var body interface{}
		if ctype == "application/json" {
			_ = json.NewDecoder(resp.Body).Decode(&body)
		}
		if body != nil {
			return fmt.Errorf("%s %s => %s, JSON response %v", method, path, resp.Status, body)
		} else {
			return fmt.Errorf("%s %s => %s, %v response", method, path, resp.Status, ctype)
		}
	}

	if result != nil {
		if ctype != "application/json" {
			return fmt.Errorf("%s %s => %s, but non-JSON response %v", method, path, resp.Status, ctype)
		}
		err = json.NewDecoder(resp.Body).Decode(result)
		if err != nil {
			return fmt.Errorf("%s %s => %s, error decoding JSON: %v", method, path, resp.Status, err)
		}
	}

	return nil
}

func makeRequest(method, baseURL, subpath string, values url.Values) (*http.Request, error) {
	components, err := url.Parse(baseURL)
	if err != nil {
		return nil, err
	}

	if subpath != "" {
		if !strings.HasPrefix(subpath, "/") {
			subpath = "/" + subpath
		}
		components.Path = components.Path + subpath
	}

	var body io.ReadCloser
	var ctype string
	if method == http.MethodGet || method == http.MethodHead {
		if values != nil {
			components.RawQuery = values.Encode()
		}
	} else {
		if values == nil {
			values = url.Values{}
		}
		body = ioutil.NopCloser(bytes.NewReader([]byte(values.Encode())))
		ctype = "application/x-www-form-urlencoded"
	}

	request, err := http.NewRequest(method, components.String(), body)
	if err != nil {
		return nil, err
	}
	if ctype != "" {
		request.Header.Set("Content-Type", ctype)
	}
	return request, nil
}

// func isJSON()

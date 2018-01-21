package digdag

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"runtime"
)

const (
	defaultBaseURL = "http://localhost:65432"
	version        = "v0.1.0" // client-version
)

// Client is the api client for digdag-server
type Client struct {
	BaseURL       *url.URL
	HTTPClient    *http.Client
	UserAgent     string
	CustomHeaders http.Header

	Verbose bool
}

// RequestOpts is the list of options to pass to the request
type RequestOpts struct {
	Params map[string]string
	// Headers map[string]string
	Body io.Reader
}

//  default UserAgent
var defaultUserAgent = fmt.Sprintf("DigdagGoClient/%s (%s)", version, runtime.Version())

// NewClient return new client for digdag
func NewClient(urlStr string, verbose bool) (*Client, error) {
	if urlStr == "" {
		// Set default
		urlStr = defaultBaseURL
	}

	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	client := &Client{
		BaseURL:    parsedURL,
		HTTPClient: &http.Client{},

		UserAgent:     defaultUserAgent,
		CustomHeaders: http.Header{},

		Verbose: verbose,
	}

	return client, err
}

// NewRequest request to digdag-server
func (c *Client) NewRequest(method, spath string, ro *RequestOpts) (resp *http.Response, err error) {
	u := *c.BaseURL
	u.Path = path.Join(c.BaseURL.Path, spath)

	if ro == nil {
		ro = new(RequestOpts)
	}

	// Add query params
	var params = make(url.Values)
	for k, v := range ro.Params {
		params.Add(k, v)
	}
	u.RawQuery = params.Encode()

	req, err := http.NewRequest(method, u.String(), ro.Body)
	if err != nil {
		return nil, err
	}

	// Set headers
	if ro.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", c.UserAgent)

	// Set custom headers
	for header, values := range c.CustomHeaders {
		for _, v := range values {
			req.Header.Set(header, v)
		}
	}

	if c.Verbose {
		dump, err := httputil.DumpRequestOut(req, true)
		if err == nil {
			log.Printf("%s", dump)
		}
	}

	client := c.HTTPClient
	resp, err = client.Do(req)
	if err != nil {
		return nil, err
	}

	if c.Verbose {
		dump, err := httputil.DumpResponse(resp, true)
		if err == nil {
			log.Printf("%s", dump)
		}
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		return resp, fmt.Errorf("Failed to request: %s", resp.Status)
	}

	return resp, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func respToString(resp *http.Response) (string, error) {
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

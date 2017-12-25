package digdag

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"path"
	"runtime"

	"github.com/franela/goreq"
	"github.com/hashicorp/errwrap"
)

const (
	defaultBaseURL = "http://localhost:65432"
	version        = "0.1" // client-version
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

//
func (c *Client) NewRequest(method, spath string, ro *RequestOpts) (resp *http.Response, err error) {

	if ro == nil {
		ro = new(RequestOpts)
	}

	u := *c.BaseURL
	u.Path = path.Join(c.BaseURL.Path, spath)

	// Set query params
	var params = make(url.Values)
	for k, v := range ro.Params {
		params.Add(k, v)
	}
	u.RawQuery = params.Encode()

	// Create request
	req, err := http.NewRequest(method, u.String(), ro.Body)
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept-Encoding", "gzip")
	req.Header.Set("User-Agent", c.UserAgent)

	// Set custom headers
	for header, values := range c.CustomHeaders {
		for _, v := range values {
			req.Header.Set(header, v)
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

	if resp.StatusCode >= 400 {
		return resp, fmt.Errorf("Failed to request: %s", resp.Status)
	}

	return resp, nil
}

func decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}

func (c *Client) doReq(method, spath string, params, res interface{}) error {
	u := *c.BaseURL
	u.Path = path.Join(c.BaseURL.Path, spath)

	req := goreq.Request{
		Method:      method,
		Uri:         u.String(),
		QueryString: params,
		ContentType: "application/json",
		UserAgent:   defaultUserAgent,
	}

	if c.Verbose {
		req.ShowDebug = true
	}

	if method == http.MethodPost || method == http.MethodPut {
		req.Body = res
	}

	req.AddHeader("Accept-Encoding", "gzip")

	resp, err := req.Do()
	if err != nil {
		return err
	}
	if resp.StatusCode >= 400 {
		return errwrap.Wrapf("Bad request: {{err}}", errors.New(resp.Status))
	}

	return resp.Body.FromJsonTo(&res)
}

func (c *Client) doRawReq(method, spath string, params interface{}) (string, error) {
	u := *c.BaseURL
	u.Path = path.Join(c.BaseURL.Path, spath)

	req, err := goreq.Request{
		Method:      method,
		Uri:         u.String(),
		QueryString: params,
		UserAgent:   defaultUserAgent,
	}.WithHeader("Accept-Encoding", "gzip").Do()

	if err != nil {
		return "", err
	}

	if req.StatusCode >= 400 {
		return "", errwrap.Wrapf("Bad request: {{err)}}", errors.New(req.Status))
	}

	body, err := req.Body.ToString()
	return body, err
}

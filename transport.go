package vk

import (
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	defaultHTTPTimeout        = 3 * time.Second
	defaultRequestTimeout     = 15 * time.Second
	defaultKeepAliveInterval  = 60 * time.Second
	defaultHTTPHeadersTimeout = defaultRequestTimeout
)

var (
	defaultHTTPClient = getDefaultHTTPClient()
)

// HTTPClient is abstaction under http client, that can Do requests
type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

// must panics on non-nil error
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (c *Client) Do(request Request, response Response) error {
	response.setRequest(request)
	req := request.HTTP()
	res, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusOK {
		return ErrBadResponseCode
	}
	return Process(res.Body).To(response)
}

// HTTP converts to *http.Request
func (r Request) HTTP() (req *http.Request) {
	values := url.Values{}
	// copy old params
	for k, v := range r.Values {
		values[k] = v
	}
	values.Add(paramVersion, defaultVersion)
	values.Add(paramHTTPS, defaultHTTPS)
	if len(r.Token) != 0 {
		values.Add(paramToken, r.Token)
	}

	u := url.URL{}
	u.Host = defaultHost
	u.Scheme = defaultScheme
	u.Path = path.Join(defaultPath, r.Method)
	u.RawQuery = values.Encode()

	req, err := http.NewRequest(defaultMethod, u.String(), nil)
	// only possible error may occur in url parsing
	// and that is completely unexpected
	must(err)

	return req
}

type vkResponseProcessor struct {
	input io.Reader
}

// ResponseProcessor fills response struct from response
type ResponseProcessor interface {
	To(response Response) error
}

// Response will be filled by ResponseProcessor
// Struct supposed to be like
// 	type Data struct {
// 	    Error `json:"error"`
// 	    Response // ... some fields of response
//	}
//
// Containing Error implements Response interface
type Response interface {
	ServerError() error
	setRequest(request Request)
}

func (d vkResponseProcessor) To(response Response) error {
	if rc, ok := d.input.(io.ReadCloser); ok {
		defer rc.Close()
	}
	decoder := json.NewDecoder(d.input)
	if err := decoder.Decode(response); err != nil {
		return err
	}
	return response.ServerError()
}

func Process(input io.Reader) ResponseProcessor {
	return vkResponseProcessor{input}
}

func getDefaultHTTPClient() HTTPClient {
	client := &http.Client{
		Timeout: defaultRequestTimeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   defaultHTTPTimeout,
				KeepAlive: defaultKeepAliveInterval,
			}).Dial,
			TLSHandshakeTimeout:   defaultHTTPTimeout,
			ResponseHeaderTimeout: defaultHTTPHeadersTimeout,
		},
	}
	return client
}

package vk

import (
	"bytes"
	"encoding/json"
	"errors"
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

func (c *Client) Do(request Request) (response *Response, err error) {
	response = new(Response)
	response.setRequest(request)
	req := request.HTTP()
	res, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	if res.StatusCode != http.StatusOK {
		return nil, ErrBadResponseCode
	}

	return Process(res.Body)
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
	To(response *Response, err error)
}

// RawString is a raw encoded JSON object.
// It implements Marshaler and Unmarshaler and can
// be used to delay JSON decoding or precompute a JSON encoding.
type Raw []byte

func (r Raw) Bytes() []byte {
	return []byte(r)
}

func (r Raw) String() string {
	return bytes.NewBuffer(r).String()
}

// MarshalJSON returns *m as the JSON encoding of m.
func (m Raw) MarshalJSON() ([]byte, error) {
	return m, nil
}

// UnmarshalJSON sets *m to a copy of data.
func (m *Raw) UnmarshalJSON(data []byte) error {
	*m = data
	return nil
}

type Response struct {
	Error    `json:"error"`
	Response Raw `json:"response"`
}

func (r Response) To(v interface{}) error {
	return json.Unmarshal(r.Response.Bytes(), v)
}

func (d vkResponseProcessor) To(response *Response) error {
	if rc, ok := d.input.(io.ReadCloser); ok {
		defer rc.Close()
	}
	decoder := json.NewDecoder(d.input)
	if err := decoder.Decode(response); err != nil {
		return err
	}
	return response.ServerError()
}

func Process(input io.Reader) (response *Response, err error) {
	response = new(Response)
	return response, vkResponseProcessor{input}.To(response)
}

type Encoder struct {
	response *Response
	err      error
}

func (e Encoder) To(v interface{}) error {
	if e.err != nil {
		return e.err
	}
	if e.response == nil {
		return errors.New("wtf")
	}
	return e.response.To(v)
}

func Encode(input io.Reader) Encoder {
	res, err := Process(input)
	return Encoder{res, err}
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

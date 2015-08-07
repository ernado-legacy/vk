package vk

import (
	"fmt"
	"net/url"
	"strconv"
	"time"

	"github.com/fatih/structs"
)

const (
	paramCode         = "code"
	paramToken        = "access_token"
	paramVersion      = "v"
	paramAppID        = "client_id"
	paramScope        = "scope"
	paramRedirectURI  = "redirect_uri"
	paramDisplay      = "display"
	paramHTTPS        = "https"
	paramResponseType = "response_type"

	oauthHost         = "oauth.vk.com"
	oauthDisplay      = "page"
	oauthPath         = "/authorize/"
	oauthResponseType = "token"
	oauthRedirectURI  = "https://oauth.vk.com/blank.html"
	oauthScheme       = "https"

	defaultHost    = "api.vk.com"
	defaultPath    = "/method/"
	defaultScheme  = "https"
	defaultVersion = "5.35"
	defaultMethod  = "GET"
	defaultHTTPS   = "1"

	maxRequestsPerSecond = 3
	minimumRate          = time.Second / maxRequestsPerSecond
	methodExecute        = "execute"
	maxRequestRepeat     = 10
)

// int64s formats int64 as base10 string
func int64s(v int64) string {
	return strconv.FormatInt(v, 10)
}

// Client for vk api
type Client struct {
	httpClient HTTPClient
}

// APIClient preforms request and fills
type APIClient interface {
	Do(request Request) (*Response, error)
}

// Request to vk api
// serializable
type Request struct {
	Method string     `json:"method"`
	Token  string     `json:"token"`
	Values url.Values `json:"values"`
}

// SetHTTPClient sets underlying http client
func (c *Client) SetHTTPClient(httpClient HTTPClient) {
	c.httpClient = httpClient
}

// Auth is helper struct for application authentication
type Auth struct {
	ID           int64
	Scope        Scope
	RedirectURI  string
	ResponseType string
	Display      string
}

type RequestFactory interface {
	Request(method string, arguments interface{}) (request Request)
}

// RequestFactory generates requests
type Factory struct {
	Token string
}

// Request generate new request with provided method and arguments
func (f Factory) Request(method string, arguments interface{}) (request Request) {
	request.Token = f.Token
	request.Method = method
	if arguments != nil {
		var argsMap map[string]interface{}
		if converted, ok := arguments.(map[string]interface{}); ok {
			argsMap = converted
		} else {
			argsMap = structs.New(arguments).Map()
		}
		request.Values = url.Values{}
		for k, v := range argsMap {
			request.Values.Add(k, fmt.Sprint(v))
		}
	}
	return request
}

// URL returns redirect url for application authentication
func (a Auth) URL() string {
	u := url.URL{}
	u.Host = oauthHost
	u.Scheme = oauthScheme
	u.Path = oauthPath

	if len(a.RedirectURI) == 0 {
		a.RedirectURI = oauthRedirectURI
	}
	if len(a.ResponseType) == 0 {
		a.ResponseType = oauthResponseType
	}
	if len(a.Display) == 0 {
		a.Display = oauthDisplay
	}

	values := u.Query()
	values.Add(paramResponseType, a.ResponseType)
	values.Add(paramScope, a.Scope.String())
	values.Add(paramAppID, int64s(a.ID))
	values.Add(paramRedirectURI, a.RedirectURI)
	values.Add(paramVersion, defaultVersion)
	values.Add(paramDisplay, a.Display)

	u.RawQuery = values.Encode()
	return u.String()
}

// New creates and returns default vk api client
func New() *Client {
	c := new(Client)
	c.SetHTTPClient(defaultHTTPClient)
	return c
}

var (
	// DefaultClient uses defaultHTTPClient for transport
	DefaultClient = New()
	// DefaultFactory with blank token
	DefaultFactory RequestFactory = Factory{}
)

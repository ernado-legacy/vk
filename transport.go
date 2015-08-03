package vk

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"path"
	"time"
)

const (
	defaultHttpTimeout        = 3 * time.Second
	defaultRequestTimeout     = 15 * time.Second
	defaultKeepAliveInterval  = 60 * time.Second
	defaultHttpHeadersTimeout = defaultRequestTimeout
)

var (
	defaultHttpClient = getDefaultHttpClient()
)

type HttpClient interface {
	Do(*http.Request) (*http.Response, error)
}

// must panics on non-nil error
func must(err error) {
	if err != nil {
		panic(err)
	}
}

func (c *Client) NewRequest(r Request) (req *http.Request) {
	values := url.Values{}
	// copy old params
	for k, v := range r.Values {
		values[k] = v
	}
	c.addParams(values)
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

type ConcurrentDecoder struct {
	Input  io.ReadCloser
	Output io.Writer
}

func (d ConcurrentDecoder) Decode(value interface{}) error {
	if d.Output == nil {
		d.Output = ioutil.Discard
	}

	// initializing flow pipes for concurrent decoding
	errorR, errorW := io.Pipe()
	valueR, valueW := io.Pipe()
	debugR, debugW := io.Pipe()
	var (
		decoderErrorChan      = make(chan error)
		errorDecoderErrorChan = make(chan error)
		readerErrorChan       = make(chan error)
		decoder               = json.NewDecoder(errorR)
		errorDecoder          = json.NewDecoder(valueR)
	)

	// value decoding goroutine
	go func() {
		decoderErrorChan <- decoder.Decode(value)
	}()

	// data reading goroutine
	go func() {
		defer d.Input.Close()
		defer errorR.Close()
		defer valueR.Close()
		defer debugR.Close()
		writer := io.MultiWriter(valueW, errorW, debugW)
		_, err := io.Copy(writer, d.Input)
		readerErrorChan <- err
	}()

	// data copying goroutine
	go io.Copy(d.Output, debugR)

	// server error decoding goroutine
	errResponse := new(ErrorResponse)
	go func() {
		errorDecoder.Decode(errResponse)
		errValue := errResponse.Error
		if errValue.Code != ErrZero {
			errorDecoderErrorChan <- errValue
		} else {
			errorDecoderErrorChan <- nil
		}
	}()

	// retrieving all errors from channels
	var err error
	for i := 0; i < 3; i++ {
		check := func(e error) {
			if e != nil {
				err = e
			}
		}
		select {
		case e := <-readerErrorChan:
			check(e)
		case e := <-decoderErrorChan:
			check(e)
		case e := <-errorDecoderErrorChan:
			check(e)
		}
	}
	return err
}

func getDefaultHttpClient() HttpClient {
	client := &http.Client{
		Timeout: defaultRequestTimeout,
		Transport: &http.Transport{
			Proxy: http.ProxyFromEnvironment,
			Dial: (&net.Dialer{
				Timeout:   defaultHttpTimeout,
				KeepAlive: defaultKeepAliveInterval,
			}).Dial,
			TLSHandshakeTimeout:   defaultHttpTimeout,
			ResponseHeaderTimeout: defaultHttpHeadersTimeout,
		},
	}
	return client
}

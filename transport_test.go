package vk

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHttpClient(t *testing.T) {
	client := getDefaultHTTPClient()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	defer server.Close()

	Convey("Test", t, func() {
		req, err := http.NewRequest("GET", server.URL, nil)
		So(err, ShouldBeNil)
		res, err := client.Do(req)
		So(err, ShouldBeNil)
		So(res.StatusCode, ShouldEqual, http.StatusNotFound)
	})
}

func BenchmarkConcurrentEncoder(b *testing.B) {
	type Data struct {
		Response []struct {
			ID int64 `json:"id"`
		} `json:"response"`
	}
	value := Data{}
	for i := 0; i < b.N; i++ {
		sData := `{
			"response": [
				{
					"id": 1,
					"first_name": "Павел",
					"last_name": "Дуров"
				}
			]
		}`
		body := ioutil.NopCloser(bytes.NewBufferString(sData))
		concurrentDecoder{Input: body}.Decode(&value)
	}
}

func BenchmarkDummyEncoder(b *testing.B) {
	type Data struct {
		Response []struct {
			ID int64 `json:"id"`
		} `json:"response"`
	}
	value := Data{}
	for i := 0; i < b.N; i++ {
		sData := `{
			"response": [
				{
					"id": 1,
					"first_name": "Павел",
					"last_name": "Дуров"
				}
			]
		}`
		body := ioutil.NopCloser(bytes.NewBufferString(sData))
		json.NewDecoder(body).Decode(&value)
	}
}

func TestMust(t *testing.T) {
	Convey("Must Panic", t, func() {
		err := ErrUnknown
		So(func() {
			must(err)
		}, ShouldPanicWith, ErrUnknown)
	})
}

func TestNewRequest(t *testing.T) {
	client := New()

	Convey("New request", t, func() {
		values := url.Values{}
		values.Add("foo", "bar")
		r := Request{Token: "token", Method: "users.get", Values: values}
		req := client.NewRequest(r)
		So(req.URL.Host, ShouldEqual, defaultHost)
		So(req.URL.String(), ShouldEqual, "https://api.vk.com/method/users.get?access_token=token&foo=bar&https=1&v=5.35")
	})
}

type errorReader struct{}

func (e errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

func (e errorReader) Close() error {
	return nil
}

func TestConcurrentEncoder(t *testing.T) {
	Convey("Encoder", t, func() {
		sData := `{
			"response": [
				{
					"id": 1,
					"first_name": "Павел",
					"last_name": "Дуров"
				}
			]
		}`
		body := ioutil.NopCloser(bytes.NewBufferString(sData))
		decoder := concurrentDecoder{Input: body}
		type Data struct {
			Response []struct {
				ID int64 `json:"id"`
			} `json:"response"`
		}
		value := Data{}
		err := decoder.Decode(&value)
		So(err, ShouldBeNil)
		Convey("Read error", func() {
			decoder := concurrentDecoder{Input: errorReader{}}
			value := Data{}
			err := decoder.Decode(&value)
			So(err, ShouldNotBeNil)
		})
		Convey("Error", func() {
			sData := `{"error":{"error_code":10,"error_msg":"Internal server error: could not get application",
			"request_params":[{"key":"oauth","value":"1"},{"key":"method","value":"users.get"},
			{"key":"user_id","value":"1"},{"key":"v","value":"5.35"}]}}
			`
			body := ioutil.NopCloser(bytes.NewBufferString(sData))
			decoder := concurrentDecoder{Input: body}
			value := Data{}
			err := decoder.Decode(&value)
			So(err, ShouldNotBeNil)
			So(IsServerError(err), ShouldBeTrue)
			serverError := GetServerError(err)
			So(serverError.Code, ShouldEqual, ErrInternalServerError)
		})
	})
}

func TestRequestSerialization(t *testing.T) {
	client := New()

	Convey("New request", t, func() {
		values := url.Values{}
		values.Add("foo", "bar")
		r := Request{Token: "token", Method: "users.get", Values: values}
		req := client.NewRequest(r)
		So(req.URL.Host, ShouldEqual, defaultHost)
		So(req.URL.String(), ShouldEqual, "https://api.vk.com/method/users.get?access_token=token&foo=bar&https=1&v=5.35")

		Convey("Marshal ok", func() {
			data, err := json.Marshal(r)
			So(err, ShouldBeNil)
			So(bytes.NewBuffer(data).String(), ShouldEqual, `{"method":"users.get","token":"token","values":{"foo":["bar"]}}`)

			Convey("Consistent after unmarshal", func() {
				newRequest := Request{}
				So(json.Unmarshal(data, &newRequest), ShouldBeNil)
			})
		})
	})
}

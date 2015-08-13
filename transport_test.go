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

func BenchmarkDummyEncoder(b *testing.B) {
	type Data struct {
		Error    `json:"error"`
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

func BenchmarkVKEncoder(b *testing.B) {
	type Data struct {
		Error    `json:"error"`
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
		So(Encode(body).To(&value), ShouldBeNil)
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
	Convey("New request", t, func() {
		values := url.Values{}
		values.Add("foo", "bar")
		r := Request{Token: "token", Method: "users.get", Values: values}
		req := r.HTTP()
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

func TestResponseProcessor(t *testing.T) {
	Convey("Decoder", t, func() {
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
		type Data []struct {
			ID int64 `json:"id"`
		}
		So(Encode(body).To(&Data{}), ShouldBeNil)
		Convey("Read error", func() {
			So(Encode(errorReader{}).To(&Data{}), ShouldNotBeNil)
		})
		Convey("Zero response", func() {
			So(Encoder{}.To(&Data{}), ShouldNotBeNil)
		})
		Convey("Error", func() {
			sData := `{"error":{"error_code":10,"error_msg":"Internal server error: could not get application",
			"request_params":[{"key":"oauth","value":"1"},{"key":"method","value":"users.get"},
			{"key":"user_id","value":"1"},{"key":"v","value":"5.35"}]}}`
			body := bytes.NewBufferString(sData)
			value := Data{}
			err := Encode(body).To(&value)
			So(err, ShouldNotBeNil)
			So(IsServerError(err), ShouldBeTrue)
			serverError := GetServerError(err)
			So(serverError.Code, ShouldEqual, ErrInternalServerError)
		})
	})
}

func TestBool(t *testing.T) {
	Convey("Ok", t, func() {
		v := &url.Values{}
		So(Bool(false).EncodeValues("test", v), ShouldBeNil)
		So(v.Get("test"), ShouldEqual, "0")
	})
}

func TestRequestSerialization(t *testing.T) {
	Convey("New request", t, func() {
		values := url.Values{}
		values.Add("foo", "bar")
		r := Request{Token: "token", Method: "users.get", Values: values}
		req := r.HTTP()
		So(req.URL.Host, ShouldEqual, defaultHost)
		So(req.URL.String(), ShouldEqual, "https://api.vk.com/method/users.get?access_token=token&foo=bar&https=1&v=5.35")

		Convey("JS", func(){
			So(r.JS(), ShouldEqual, `API.users.get({"foo":"bar"})`)
		})
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

func TestResponseSerialization(t *testing.T) {
	Convey("New response", t, func() {
		r := Response{}
		type payload struct {
			Foo string
		}
		Convey("Marshal ok", func() {
			_, err := json.Marshal(r)
			So(err, ShouldBeNil)
		})
		Convey("Payload", func() {
			r := Response{Response: bytes.NewBufferString(`{"Foo": "test"}`).Bytes()}
			Convey("Marshal ok", func() {
				d, err := json.Marshal(r)
				So(err, ShouldBeNil)
				So(bytes.NewBuffer(d).String(), ShouldEqual, `{"error":{},"response":{"Foo":"test"}}`)
				So(r.Response.String(), ShouldEqual, `{"Foo": "test"}`)
				Convey("Unmarshal", func() {
					res := Response{}
					So(json.Unmarshal(d, &res), ShouldBeNil)
					p := payload{}
					So(json.Unmarshal(res.Response, &p), ShouldBeNil)
					So(p.Foo, ShouldEqual, "test")
				})
			})
		})
	})
}

type simpleHTTPClientMock struct {
	response *http.Response
	err      error
}

type apiClientMock struct {
	callback func(Request, Response) error
}

func (client apiClientMock) Do(request Request, response Response) error {
	return client.callback(request, response)
}

func (m simpleHTTPClientMock) Do(request *http.Request) (*http.Response, error) {
	return m.response, m.err
}

func TestDo(t *testing.T) {
	client := New()

	Convey("Do request", t, func() {
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
		httpResponse := &http.Response{Body: body, StatusCode: http.StatusOK}
		client.SetHTTPClient(simpleHTTPClientMock{response: httpResponse})
		type Data []struct {
			ID   int64  `json:"id"`
			Name string `json:"first_name"`
		}

		request := Request{Method: "users.get"}
		response := Data{}

		res, err := client.Do(request)
		So(err, ShouldBeNil)
		So(res.To(&response), ShouldBeNil)
		So(response[0].ID, ShouldEqual, 1)
		So(response[0].Name, ShouldEqual, "Павел")

		Convey("Bad status", func() {
			client := New()
			httpResponse := &http.Response{Body: body, StatusCode: http.StatusBadRequest}
			client.SetHTTPClient(simpleHTTPClientMock{response: httpResponse})
			request := Request{Method: "users.get"}
			//			response := &Data{}

			_, err := client.Do(request)
			So(err, ShouldEqual, ErrBadResponseCode)
		})
		Convey("Http error", func() {
			client := New()
			httpResponse := &http.Response{Body: body, StatusCode: http.StatusBadRequest}
			client.SetHTTPClient(simpleHTTPClientMock{response: httpResponse, err: ErrBadResponseCode})
			request := Request{Method: "users.get"}
			//			response := &Data{}

			_, err := client.Do(request)
			So(err, ShouldEqual, ErrBadResponseCode)
		})
	})
}

//func TestDoRawResponse(t *testing.T) {
//	client := New()
//
//	Convey("Do request", t, func() {
//		sData := `{
//			"response": [
//				{
//					"id": 1,
//					"first_name": "Павел",
//					"last_name": "Дуров"
//				}
//			]
//		}`
//		body := ioutil.NopCloser(bytes.NewBufferString(sData))
//		httpResponse := &http.Response{Body: body, StatusCode: http.StatusOK}
//		client.SetHTTPClient(simpleHTTPClientMock{response: httpResponse})
//
//		request := Request{Method: "users.get"}
//		res, err := client.Do(request)
//		So(err, ShouldBeNil)
//
//		expectedRaw := `[
//				{
//					"id": 1,
//					"first_name": "Павел",
//					"last_name": "Дуров"
//				}
//			]`
//		So(expectedRaw, ShouldEqual, res.Response.String())
//	})
//}

package vk

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestNewClient(t *testing.T) {
	client := New()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))

	defer server.Close()

	Convey("Test GET", t, func() {
		req, err := http.NewRequest("GET", server.URL, nil)
		So(err, ShouldBeNil)
		res, err := client.httpClient.Do(req)
		So(err, ShouldBeNil)
		So(res.StatusCode, ShouldEqual, http.StatusNotFound)
	})
}

func TestFactory(t *testing.T) {
	Convey("Factory", t, func() {
		Convey("Request should contain token", func() {
			f := RequestFactory{"token"}
			So(f.Request("method", nil).Token, ShouldEqual, "token")
		})
		Convey("Default factory  request should not contain token", func() {
			So(DefaultFactory.Request("method", nil).Token, ShouldBeBlank)
		})
		Convey("Should correctly convert args", func() {
			f := RequestFactory{"token"}
			r := f.Request("method", map[string]interface{}{"test": 1234567891})
			So(r.Token, ShouldEqual, "token")
			So(r.Values.Get("test"), ShouldEqual, "1234567891")
		})
		Convey("From struct", func() {
			f := RequestFactory{"token"}
			type Data struct {
				Test int `structs:"test"`
			}
			args := Data{1234567891}
			r := f.Request("method", args)
			So(r.Token, ShouldEqual, "token")
			So(r.Values.Get("test"), ShouldEqual, "1234567891")
		})
	})
}

func TestAuthUrl(t *testing.T) {
	Convey("URL is valid", t, func() {
		stringURL := Auth{Scope: NewScope(PermOffline, PermGroups)}.URL()
		gotURL, err := url.Parse(stringURL)
		So(err, ShouldBeNil)
		So(gotURL.Host, ShouldEqual, oauthHost)
		shouldURL := "https://oauth.vk.com/authorize/?client_id=0&display=page&redirect_uri=https%3A%2F%2Foauth.vk.com%2Fblank.html&response_type=token&scope=groups%2Coffline&v=5.35"
		So(shouldURL, ShouldEqual, gotURL.String())
	})
}

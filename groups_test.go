package vk

import (
	"testing"

	"bytes"
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
)

type apiJSONMock struct {
	response string
	err      error
}

func (api apiJSONMock) Do(req Request) (res *Response, err error) {
	if api.err != nil {
		return nil, api.err
	}
	if _, err := json.Marshal(req); err != nil {
		return nil, err
	}
	res = new(Response)
	return res, json.NewDecoder(bytes.NewBufferString(api.response)).Decode(res)
}

type recordFactory struct {
	request Request
}

func (f *recordFactory) Request(method string, arguments interface{}) (request Request) {
	f.request = DefaultFactory.Request(method, arguments)
	return request
}

func rf() recordFactory {
	return recordFactory{}
}

func TestResource(t *testing.T) {
	Convey("Resource", t, func() {
		resource := Resource{APIClient: apiJSONMock{"", ErrAuthFailed}, RequestFactory: DefaultFactory}
		response := Response{}
		request := Request{}
		So(resource.Decode(request, &response), ShouldEqual, ErrAuthFailed)
	})
}

func newApiMock(d string, err error) (a apiJSONMock) {
	a.err = err
	a.response = d
	return a
}

func record(c APIClient, f RequestFactory) Resource {
	return Resource{APIClient: c, RequestFactory: f}
}

func TestGroupsMethods(t *testing.T) {
	Convey("Groups methods", t, func() {
		Convey(methodGroupsGet, func() {
			mock := newApiMock(`{"response" :
				{"count": 0}
			}`, nil)
			f := rf()
			g := Groups{record(mock, &f)}
			groups, err := g.Get(GroupGetFields{Count: 1})
			So(err, ShouldBeNil)
			So(groups.Count, ShouldEqual, 0)
		})
	})
}

func TestGroups(t *testing.T) {
	Convey("Groups", t, func() {
		mock := apiJSONMock{`{"response":
		{"count":309676,
		"items":[
		{"id":4189,"first_name":"Николай","last_name":"Матвеев",
		"sex":2,"bdate":"24.6","city":{"id":2,"title":"Санкт-Петербург"}, "country":{"id":1,"title":"Россия"}},
		{"id":4234,"first_name":"Никита","last_name":"Слушкин","sex":2,"city":{"id":2,"title":"Санкт-Петербург"}}]}}
		`, nil}
		g := Groups{Resource{APIClient: mock, RequestFactory: DefaultFactory}}
		members, err := g.GetMembers(GroupSearchFields{})
		So(err, ShouldBeNil)
		So(members.Count, ShouldEqual, 309676)
		So(len(members.Items), ShouldEqual, 2)
		So(members.Items[0].FirstName, ShouldEqual, "Николай")
		So(members.Items[0].Sex, ShouldEqual, Male)
		user := members.Items[0]
		So(user.Country.Is(Russia), ShouldBeTrue)
		Convey("JSON", func() {
			Convey("Marshal", func() {
				v := struct {
					Value Bool
				}{true}
				data, err := json.Marshal(v)
				So(err, ShouldBeNil)
				sData := bytes.NewBuffer(data).String()
				So(sData, ShouldEqual, `{"Value":1}`)

				v2 := struct {
					Value Bool
				}{false}
				data, err = json.Marshal(v2)
				So(err, ShouldBeNil)
				sData = bytes.NewBuffer(data).String()
				So(sData, ShouldEqual, `{"Value":0}`)
				Convey("Errors", func() {
					v := new(Bool)
					So(v.UnmarshalJSON([]byte("1")), ShouldBeNil)
					So(*v, ShouldEqual, true)

					So(v.UnmarshalJSON([]byte("0")), ShouldBeNil)
					So(*v, ShouldEqual, false)

					So(v.UnmarshalJSON([]byte("-1")), ShouldNotBeNil)
					So(v.UnmarshalJSON([]byte("9")), ShouldNotBeNil)

					So(v.UnmarshalJSON([]byte{}), ShouldBeNil)
					So(*v, ShouldEqual, false)
				})
			})
			Convey("True", func() {
				data := []byte(`
				{"is_closed": 1}
				`)
				value := Group{}
				So(json.Unmarshal(data, &value), ShouldBeNil)
				So(value.IsClosed, ShouldEqual, true)
			})
			Convey("False", func() {
				data := []byte(`
				{"is_closed": 0}
				`)
				value := Group{}
				So(json.Unmarshal(data, &value), ShouldBeNil)
				So(value.IsClosed, ShouldEqual, false)
			})
		})
	})
}

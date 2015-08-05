package vk

import (
	"testing"

	"bytes"
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"log"
)

type apiJSONMock struct {
	response string
}

func (api apiJSONMock) Do(req Request, res Response) error {
	if _, err := json.Marshal(req); err != nil {
		return err
	}
	return json.NewDecoder(bytes.NewBufferString(api.response)).Decode(res)
}

func TestGroups(t *testing.T) {
	Convey("Groups", t, func() {
		mock := apiJSONMock{`{"response":
		{"count":309676,
		"items":[
		{"id":4189,"first_name":"Николай","last_name":"Матвеев",
		"sex":2,"bdate":"24.6","city":{"id":2,"title":"Санкт-Петербург"}, "country":{"id":1,"title":"Россия"}},
		{"id":4234,"first_name":"Никита","last_name":"Слушкин","sex":2,"city":{"id":2,"title":"Санкт-Петербург"}}]}}
		`}
		log.Println(mock.response)
		g := Groups{DefaultFactory, mock}
		members, err := g.GetMembers(GroupSearchFields{})
		So(err, ShouldBeNil)
		So(members.Count, ShouldEqual, 309676)
		So(len(members.Items), ShouldEqual, 2)
		So(members.Items[0].FirstName, ShouldEqual, "Николай")
		So(members.Items[0].Sex, ShouldEqual, Male)
		user := members.Items[0]
		So(user.Country.Is(Russia), ShouldBeTrue)
		Convey("JSON", func() {
			Convey("Marshal", func(){
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

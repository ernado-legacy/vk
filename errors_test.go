package vk

import (
	"encoding/json"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
	"io"
)

func TestErrors(t *testing.T) {
	Convey("Errors", t, func() {
		Convey("Recognition", func(){
			So(IsServerError(io.ErrUnexpectedEOF), ShouldBeFalse)
			So(func(){ GetServerError(io.ErrUnexpectedEOF) }, ShouldPanic)
		})
		Convey("Parsing from JSON", func() {
			type Data struct {
				Error ServerError `json:"error"`
			}
			data := []byte(`{"error": 1}`)
			v := Data{}
			So(json.Unmarshal(data, &v), ShouldBeNil)
			So(v.Error, ShouldEqual, ErrUnknown)
			Convey("Use as error", func() {
				var err error
				err = v.Error
				So(err.Error(), ShouldEqual, "1")
			})
		})
	})
}

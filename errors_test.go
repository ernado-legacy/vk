package vk

import (
	"encoding/json"
	"io"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestErrors(t *testing.T) {
	Convey("Errors", t, func() {
		Convey("Recognition", func() {
			So(IsServerError(io.ErrUnexpectedEOF), ShouldBeFalse)
			So(func() { GetServerError(io.ErrUnexpectedEOF) }, ShouldPanic)
		})
		Convey("Equality", func() {
			So(ErrZero.Is(ErrZero), ShouldBeTrue)
			So(ErrZero.Is(ErrAlbumOverflow), ShouldBeFalse)
			So(ErrAlbumOverflow.Is(ErrZero), ShouldBeFalse)
			So(ErrZero.Is(io.ErrUnexpectedEOF), ShouldBeFalse)
			So(ErrZero.Is(Error{Code: ErrZero}), ShouldBeTrue)
			So(ErrZero.Is(error(Error{Code: ErrZero})), ShouldBeTrue)
			So(ErrZero.Is(Error{Code: ErrAlbumOverflow}), ShouldBeFalse)
			So(ErrZero.Is(error(Error{Code: ErrAlbumOverflow})), ShouldBeFalse)
		})
		Convey("Set and get request", func() {
			e := Error{}
			e.setRequest(Request{Method: "test"})
			So(e.Request.Method, ShouldEqual, "test")
		})
		Convey("Parsing from JSON", func() {
			type Data struct {
				Error ServerError `json:"error"`
			}
			data := []byte(`{"error": 1}`)
			v := Data{}
			So(json.Unmarshal(data, &v), ShouldBeNil)
			So(v.Error.Error(), ShouldEqual, ErrUnknown.Error())
			Convey("Use as error", func() {
				var err error
				err = v.Error
				So(err.Error(), ShouldEqual, "ErrUnknown")
			})
		})
	})
}

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
		Convey("String", func() {
			Convey("Error with textual representation", func() {
				var err error
				err = ErrUnknown
				So(err.Error(), ShouldEqual, "Unknown error occured, try again later (1)")
			})
			Convey("Error without representation", func() {
				var err error
				err = ErrMoneyTransferNotAllowed
				So(err.Error(), ShouldEqual, "500")
			})
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
				So(err.Error(), ShouldEqual, "Unknown error occured, try again later (1)")
			})
			Convey("Representation", func() {
				So(v.Error.String(), ShouldEqual, "Unknown error occured, try again later (1)")
			})
			So(v.Error.Error(), ShouldEqual, "Unknown error occured, try again later (1)")
		})
	})
}

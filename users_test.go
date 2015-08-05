package vk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUserStringer(t *testing.T) {
	Convey("Sex", t, func() {
		So(Female.String(), ShouldEqual, "female")
		So(Male.String(), ShouldEqual, "male")
		So(SexUnknown.String(), ShouldEqual, "unknown")
	})
	Convey("Country", t, func(){
		So(Country{Russia, "Россия"}.String(), ShouldEqual, "Россия")
		So(Country{CountryUnknown, "t"}.String(), ShouldEqual, "unknown")
		So(Country{Russia, "Россия"}.Is(Russia), ShouldBeTrue)
		So(Country{1, "Россия"}.Is(Russia), ShouldBeTrue)
		So(Country{0, "Россия"}.Is(Russia), ShouldBeFalse)
	})
}

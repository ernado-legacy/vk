package vk

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestScope(t *testing.T) {
	Convey("Scope", t, func() {
		Convey("Must be empty", func() {
			s := Scope{}
			So(s.Has(PermOffline), ShouldBeFalse)
			So(len(s), ShouldEqual, 0)
		})
		Convey("Nil is valid", func() {
			var s Scope
			So(s.Has(PermOffline), ShouldBeFalse)
			So(len(s), ShouldEqual, 0)
			So(func() { s.Del(PermFriends) }, ShouldNotPanic)
		})
		Convey("Add", func() {
			s := Scope{}
			s.Add(PermOffline)
			So(s.Has(PermOffline), ShouldBeTrue)
			So(len(s), ShouldEqual, 1)
			Convey("Remove", func() {
				s.Del(PermOffline)
				So(s.Has(PermOffline), ShouldBeFalse)
				So(len(s), ShouldEqual, 0)
			})
		})
		Convey("From list permissions", func() {
			s := NewScope(PermOffline, PermFriends)
			So(s.Has(PermOffline), ShouldBeTrue)
			So(s.Has(PermFriends), ShouldBeTrue)
			So(s.Has(PermGroups), ShouldBeFalse)
		})
	})
}

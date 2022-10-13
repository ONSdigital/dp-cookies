package cookies

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestParseCookieTime(t *testing.T) {
	Convey("Given a datetime string in the CookieTime format", t, func() {
		datetime := "2021-01-01T15:31:23"

		Convey("calling ParseCookieTime returns the CookieTime value", func() {
			v, e := ParseCookieTime(datetime)
			So(e, ShouldBeNil)
			So(v, ShouldHaveSameTypeAs, CookieTime{})
			So(v.String(), ShouldEqual, datetime)
		})
	})

	Convey("Given a datetime string in an incorrect format", t, func() {
		datetime := "1/1/2021 15:31:23"

		Convey("calling ParseCookieTime returns an indicative error", func() {
			v, e := ParseCookieTime(datetime)
			So(v, ShouldResemble, CookieTime{})
			So(e, ShouldNotBeNil)

			invalid, ok := e.(InvalidCookieTimeString)
			So(ok, ShouldBeTrue)
			So(invalid.value, ShouldEqual, datetime)
			So(invalid.err, ShouldContainSubstring, "cannot parse")
		})
	})
}

func TestCookieTime_Marshal_UnmarshalJSON(t *testing.T) {
	Convey("Given a CookieTime", t, func() {
		datetime := MustParseCookieTime("2021-01-01T15:31:23")

		Convey("marshalling its value to json gives a quoted string value for the CookieTime in the CookieTime format", func() {
			v, e := json.Marshal(datetime)
			So(e, ShouldBeNil)
			So(v, ShouldResemble, []byte(`"2021-01-01T15:31:23"`))

			Convey("and unmarshalling the json encoded CookieTime gives back the original CookieTime", func() {
				var ct CookieTime
				e := json.Unmarshal(v, &ct)
				So(e, ShouldBeNil)
				So(ct, ShouldResemble, datetime)
			})
		})
	})
}

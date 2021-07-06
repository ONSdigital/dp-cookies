package cookies

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitCookie(t *testing.T) {

	var testDomain = "www.test.com"
	var testCookie = "test_cookie"
	var testValue = "test-value"

	Convey("Set sets correct cookie", t, func() {
		rec := httptest.NewRecorder()
		set(rec, testCookie, testValue, testDomain, "/", 12, http.SameSiteLaxMode, false)
		cookie := rec.Result().Cookies()[0]
		So(cookie.Value, ShouldEqual, testValue)
		So(cookie.Path, ShouldEqual, "/")
		So(cookie.Domain, ShouldEqual, testDomain)
		So(cookie.MaxAge, ShouldEqual, 12)
		So(cookie.Secure, ShouldBeTrue)
		So(cookie.SameSite, ShouldEqual, http.SameSiteLaxMode)
	})

	Convey("Get", t, func() {
		Convey("returns cookie value if value is set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: "test-cookie", Value: "test-value"})
			cookie, _ := get(req, "test-cookie")
			So(cookie, ShouldEqual, "test-value")
		})

		Convey("returns error if no cookie is set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			_, err := GetLang(req)
			So(err, ShouldNotBeNil)
		})
	})

}

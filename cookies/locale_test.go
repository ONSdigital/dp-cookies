package cookies

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitLocale(t *testing.T) {

	var testDomain = "www.test.com"
	var testLang = "en"

	Convey("GetLang", t, func() {
		Convey("returns cookie value if value is set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: localeCookieKey, Value: testLang})
			cookie, _ := GetLang(req)
			So(cookie, ShouldEqual, testLang)
		})

		Convey("returns error if no cookie is set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			_, err := GetLang(req)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("SetLang sets correct cookie", t, func() {
		rec := httptest.NewRecorder()
		SetLang(rec, testLang, testDomain)
		cookie := rec.Result().Cookies()[0]
		So(cookie.Value, ShouldEqual, testLang)
		So(cookie.Path, ShouldEqual, "/")
		So(cookie.Domain, ShouldEqual, testDomain)
		So(cookie.MaxAge, ShouldEqual, maxAgeOneYear)
		So(cookie.Secure, ShouldBeTrue)
		So(cookie.SameSite, ShouldEqual, http.SameSiteLaxMode)
	})
}

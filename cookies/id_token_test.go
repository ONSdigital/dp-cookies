package cookies

import (
"net/http"
"net/http/httptest"
"testing"

. "github.com/smartystreets/goconvey/convey"
)

func TestUnitIDToken(t *testing.T) {

	var testDomain = "www.test.com"
	var testIDToken = "test-id-token"

	Convey("GetIDToken", t, func() {
		Convey("returns cookie value if value is set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: idCookieKey, Value: testIDToken})
			cookie, _ := GetIDToken(req)
			So(cookie, ShouldEqual, testIDToken)
		})

		Convey("returns error if no cookie is set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			_, err := GetIDToken(req)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("SetIDToken sets correct cookie", t, func() {
		rec := httptest.NewRecorder()
		SetIDToken(rec, testIDToken, testDomain)
		cookie := rec.Result().Cookies()[0]
		So(cookie.Value, ShouldEqual, testIDToken)
		So(cookie.Path, ShouldEqual, "/")
		So(cookie.Domain, ShouldEqual, testDomain)
		So(cookie.MaxAge, ShouldEqual, maxAgeBrowserSession)
		So(cookie.Secure, ShouldBeTrue)
		So(cookie.SameSite, ShouldEqual, http.SameSiteLaxMode)
	})
}


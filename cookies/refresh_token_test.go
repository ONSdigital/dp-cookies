package cookies


import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitRefreshToken(t *testing.T) {

	var testDomain = "www.test.com"
	var testRefreshToken = "test-refresh-token"

	Convey("GetRefreshToken", t, func() {
		Convey("returns cookie value if value is set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: refreshCookieKey, Value: testRefreshToken})
			cookie, _ := GetRefreshToken(req)
			So(cookie, ShouldEqual, testRefreshToken)
		})

		Convey("returns error if no cookie is set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			_, err := GetRefreshToken(req)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("SetRefreshToken sets correct cookie", t, func() {
		rec := httptest.NewRecorder()
		SetRefreshToken(rec, testRefreshToken, testDomain)
		cookie := rec.Result().Cookies()[0]
		So(cookie.Value, ShouldEqual, testRefreshToken)
		So(cookie.Path, ShouldEqual, "/")
		So(cookie.Domain, ShouldEqual, testDomain)
		So(cookie.MaxAge, ShouldEqual, maxAgeBrowserSession)
		So(cookie.Secure, ShouldBeTrue)
		So(cookie.SameSite, ShouldEqual, http.SameSiteLaxMode)
	})
}


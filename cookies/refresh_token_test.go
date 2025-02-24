package cookies

import (
	"net/http"
	"net/http/httptest"
	"net/url"
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
		correctCookie := &http.Cookie{
			Name:     "refresh_token",
			Value:    url.QueryEscape(testRefreshToken),
			Path:     "/api/v1/tokens/self",
			Domain:   testDomain,
			HttpOnly: true,
			Secure:   true,
			MaxAge:   maxAgeBrowserSession,
			SameSite: http.SameSiteStrictMode,
			Raw:      "refresh_token=test-refresh-token; Path=/api/v1/tokens/self; Domain=www.test.com; HttpOnly; Secure; SameSite=Strict",
		}
		rec := httptest.NewRecorder()
		SetRefreshToken(rec, testRefreshToken, testDomain)
		cookie := rec.Result().Cookies()[0]
		So(cookie, ShouldResemble, correctCookie)
	})
}

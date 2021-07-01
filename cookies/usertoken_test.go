package cookies

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitUserToken(t *testing.T) {

	var testDomain = "www.test.com"
	var testAccessToken = "test-access-token"

	Convey("GetUserAuthToken", t, func() {
		Convey("returns cookie value if value is set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: florenceCookieKey, Value: testAccessToken})
			cookie, _ := GetUserAuthToken(req)
			So(cookie, ShouldEqual, testAccessToken)
		})

		Convey("returns error if no cookie is set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			_, err := GetUserAuthToken(req)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("SetUserAuthToken sets correct cookie", t, func() {
		rec := httptest.NewRecorder()

		correctCookie := &http.Cookie{
			Name:     florenceCookieKey,
			Value:    url.QueryEscape(testAccessToken),
			Path:     "/",
			Domain:   testDomain,
			HttpOnly: true,
			Secure:   true,
			MaxAge:   maxAgeBrowserSession,
			SameSite: http.SameSiteStrictMode,
			Raw:      "access_token=test-access-token; Path=/; Domain=www.test.com; HttpOnly; Secure; SameSite=Strict",
		}

		SetUserAuthToken(rec, testAccessToken, testDomain)
		cookie := rec.Result().Cookies()[0]
		So(cookie, ShouldResemble, correctCookie)
	})
}

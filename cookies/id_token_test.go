package cookies

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitIDToken(t *testing.T) {
	var testDomain = "www.test.com"
	var testIDToken = "test-id-token"

	Convey("GetIDToken", t, func() {
		Convey("returns cookie value if value is set", func() {
			req := httptest.NewRequest("GET", "/", http.NoBody)
			req.AddCookie(&http.Cookie{Name: idCookieKey, Value: testIDToken})
			cookie, _ := GetIDToken(req)
			So(cookie, ShouldEqual, testIDToken)
		})

		Convey("returns error if no cookie is set", func() {
			req := httptest.NewRequest("GET", "/", http.NoBody)
			_, err := GetIDToken(req)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("SetIDToken sets correct cookie", t, func() {
		rec := httptest.NewRecorder()

		correctCookie := &http.Cookie{
			Name:     idCookieKey,
			Value:    url.QueryEscape(testIDToken),
			Path:     "/",
			Domain:   testDomain,
			HttpOnly: false,
			Secure:   true,
			MaxAge:   maxAgeBrowserSession,
			SameSite: http.SameSiteLaxMode,
			Raw:      "id_token=test-id-token; Path=/; Domain=www.test.com; Secure; SameSite=Lax",
		}

		SetIDToken(rec, testIDToken, testDomain)
		cookie := rec.Result().Cookies()[0]
		So(cookie, ShouldResemble, correctCookie)
	})
}

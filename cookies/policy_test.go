package cookies

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var testEncodedCookieValue = "%7B%22essential%22%3Atrue%2C%22usage%22%3Atrue%7D"

func TestUnitPolicy(t *testing.T) {

	Convey("GetCookiePreferences", t, func() {
		Convey("returns false for prefrences set if cookie isn't set, and default policy if no cookies set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			cookie := GetCookiePreferences(req)
			So(cookie, ShouldResemble, PreferencesResponse{
				IsPreferenceSet: false,
				Policy: Policy{
					Essential: true,
					Usage:     false,
				},
			})
		})

		Convey("returns true for prefrences set if cookie is set, and default policy if no cookies set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: cookiesPreferencesSetCookieKey, Value: "true"})
			cookie := GetCookiePreferences(req)
			So(cookie, ShouldResemble, PreferencesResponse{
				IsPreferenceSet: true,
				Policy: Policy{
					Essential: true,
					Usage:     false,
				},
			})
		})

		Convey("returns true for prefrences set if cookie is set, and correct policy if cookie set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: cookiesPreferencesSetCookieKey, Value: "true"})
			req.AddCookie(&http.Cookie{Name: cookiesPolicyCookieKey, Value: testEncodedCookieValue})
			cookie := GetCookiePreferences(req)
			So(cookie, ShouldResemble, PreferencesResponse{
				IsPreferenceSet: true,
				Policy: Policy{
					Essential: true,
					Usage:     true,
				},
			})
		})
	})

	Convey("SetPreferenceIsSet sets correct cookie", t, func() {
		rec := httptest.NewRecorder()
		SetPreferenceIsSet(rec, testDomain)
		cookie := rec.Result().Cookies()[0]
		So(cookie.Value, ShouldEqual, "true")
		So(cookie.Path, ShouldEqual, "/")
		So(cookie.Domain, ShouldEqual, testDomain)
		So(cookie.MaxAge, ShouldEqual, maxAgeOneYear)
		So(cookie.Secure, ShouldBeTrue)
		So(cookie.SameSite, ShouldEqual, http.SameSiteLaxMode)
	})

	Convey("SetPolicy sets correct cookie", t, func() {
		rec := httptest.NewRecorder()
		SetPolicy(rec, Policy{Essential: true, Usage: true}, testDomain)
		cookie := rec.Result().Cookies()[0]
		So(cookie.Value, ShouldEqual, testEncodedCookieValue)
		So(cookie.Path, ShouldEqual, "/")
		So(cookie.Domain, ShouldEqual, testDomain)
		So(cookie.MaxAge, ShouldEqual, maxAgeOneYear)
		So(cookie.Secure, ShouldBeTrue)
		So(cookie.SameSite, ShouldEqual, http.SameSiteLaxMode)
	})
}

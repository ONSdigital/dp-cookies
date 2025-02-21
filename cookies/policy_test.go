package cookies

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitPolicy(t *testing.T) {
	var testDomain = "www.test.com"
	var testEncodedCookieValue = "%7B%22essential%22%3Atrue%2C%22usage%22%3Atrue%7D"

	Convey("GetCookiePreferences", t, func() {
		Convey("returns false for prefrences set if cookie isn't set, and default policy if no cookies set", func() {
			req := httptest.NewRequest("GET", "/", http.NoBody)
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
			req := httptest.NewRequest("GET", "/", http.NoBody)
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
			req := httptest.NewRequest("GET", "/", http.NoBody)
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

func TestUnitONSPolicy(t *testing.T) {
	Convey("GetONSCookiePreferences", t, func() {
		Convey("returns false for preferences set if cookie isn't set, and default policy if no cookies set", func() {
			req := httptest.NewRequest("GET", "/", http.NoBody)
			cookie := GetONSCookiePreferences(req)
			So(cookie, ShouldResemble, ONSPreferencesResponse{
				IsPreferenceSet: false,
				Policy: ONSPolicy{
					Essential: true,
					Settings:  false,
					Usage:     false,
					Campaigns: false,
				},
			})
		})

		Convey("returns true for preferences set if cookie is set, and default policy if no cookies set", func() {
			req := httptest.NewRequest("GET", "/", http.NoBody)
			req.AddCookie(&http.Cookie{Name: onsCookiePreferencesSetCookieKey, Value: "true"})
			cookie := GetONSCookiePreferences(req)
			So(cookie, ShouldResemble, ONSPreferencesResponse{
				IsPreferenceSet: true,
				Policy: ONSPolicy{
					Essential: true,
					Settings:  false,
					Usage:     false,
					Campaigns: false,
				},
			})
		})

		tc := []struct {
			scenario string
			given    string
			expected ONSPolicy
		}{
			{
				scenario: "all cookie preferences are set",
				given:    "{'essential':true,'settings':true,'usage':true,'campaigns':true}",
				expected: ONSPolicy{
					Essential: true,
					Settings:  true,
					Usage:     true,
					Campaigns: true,
				},
			},
			{
				scenario: "essential and settings cookie preferences are set",
				given:    "{'essential':true,'settings':true,'usage':false,'campaigns':false}",
				expected: ONSPolicy{
					Essential: true,
					Settings:  true,
					Usage:     false,
					Campaigns: false,
				},
			},
			{
				scenario: "essential and usage cookie preferences are set",
				given:    "{'essential':true,'settings':false,'usage':true,'campaigns':false}",
				expected: ONSPolicy{
					Essential: true,
					Settings:  false,
					Usage:     true,
					Campaigns: false,
				},
			},
			{
				scenario: "essential and campaigns cookie preferences are set",
				given:    "{'essential':true,'settings':false,'usage':false,'campaigns':true}",
				expected: ONSPolicy{
					Essential: true,
					Settings:  false,
					Usage:     false,
					Campaigns: true,
				},
			},
			{
				scenario: "essential, usage and campaigns cookie preferences are set",
				given:    "{'essential':true,'settings':false,'usage':true,'campaigns':true}",
				expected: ONSPolicy{
					Essential: true,
					Settings:  false,
					Usage:     true,
					Campaigns: true,
				},
			},
		}

		Convey("returns true for preference set if cookie is set, and matches preferences", func() {
			for _, t := range tc {
				Convey(fmt.Sprintf("when %s", t.scenario), func() {
					req := httptest.NewRequest("GET", "/", http.NoBody)
					req.AddCookie(&http.Cookie{Name: onsCookiePreferencesSetCookieKey, Value: "true"})
					req.AddCookie(&http.Cookie{Name: onsCookiePolicyCookieKey, Value: t.given})
					cookie := GetONSCookiePreferences(req)
					So(cookie, ShouldResemble, ONSPreferencesResponse{
						IsPreferenceSet: true,
						Policy:          t.expected,
					})
				})
			}
		})
	})

	Convey("SetONSPreferenceIsSet sets correct cookie", t, func() {
		rec := httptest.NewRecorder()
		SetONSPreferenceIsSet(rec, testDomain)
		cookie := rec.Result().Cookies()[0]
		So(cookie.Value, ShouldEqual, "true")
		So(cookie.Path, ShouldEqual, "/")
		So(cookie.Domain, ShouldEqual, testDomain)
		So(cookie.MaxAge, ShouldEqual, maxAgeOneYear)
		So(cookie.Secure, ShouldBeTrue)
		So(cookie.SameSite, ShouldEqual, http.SameSiteLaxMode)
	})

	Convey("SetONSPolicy", t, func() {
		tc := []struct {
			given    ONSPolicy
			expected string
		}{
			{
				given:    ONSPolicy{},
				expected: "{'essential':false,'settings':false,'usage':false,'campaigns':false}",
			},
			{
				given: ONSPolicy{
					Essential: true,
					Settings:  false,
					Usage:     false,
					Campaigns: false,
				},
				expected: "{'essential':true,'settings':false,'usage':false,'campaigns':false}",
			},
			{
				given: ONSPolicy{
					Essential: true,
					Settings:  true,
					Usage:     false,
					Campaigns: false,
				},
				expected: "{'essential':true,'settings':true,'usage':false,'campaigns':false}",
			},
			{
				given: ONSPolicy{
					Essential: true,
					Settings:  true,
					Usage:     true,
					Campaigns: false,
				},
				expected: "{'essential':true,'settings':true,'usage':true,'campaigns':false}",
			},
			{
				given: ONSPolicy{
					Essential: true,
					Settings:  true,
					Usage:     true,
					Campaigns: true,
				},
				expected: "{'essential':true,'settings':true,'usage':true,'campaigns':true}",
			},
		}

		for _, t := range tc {
			Convey(fmt.Sprintf("when preferences are set as: %v", t.given), func() {
				Convey("then the set cookie matches the preferences", func() {
					rec := httptest.NewRecorder()
					SetONSPolicy(rec, t.given, testDomain)
					cookie := rec.Result().Cookies()[0]
					So(cookie.Value, ShouldEqual, t.expected)
					So(cookie.Path, ShouldEqual, "/")
					So(cookie.Domain, ShouldEqual, testDomain)
					So(cookie.MaxAge, ShouldEqual, maxAgeOneYear)
					So(cookie.Secure, ShouldBeTrue)
					So(cookie.SameSite, ShouldEqual, http.SameSiteLaxMode)
				})
			})
		}
	})
}

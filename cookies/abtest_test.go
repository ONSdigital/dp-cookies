package cookies

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestUnitABTest(t *testing.T) {

	testDomain := "www.ons.gov.uk"
	testTime := time.Date(2009, 11, 17, 20, 34, 58, 651387237, time.UTC)
	testServices := ABServices{NewSearch: &testTime}
	expectedEscapedCookieValue := "%7B%22new_search%22%3A%222009-11-17T20%3A34%3A58.651387237Z%22%7D"
	testMarshalledValue, err := json.Marshal(testServices)
	if err != nil {
		t.Fatal("failed to marshal test data", err)
	}
	testEscapedValue := url.QueryEscape(string(testMarshalledValue))

	Convey("GetABTest", t, func() {
		Convey("returns cookie value if value is set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: aBTestKey, Value: testEscapedValue})
			cookie, err := GetABTest(req)
			So(cookie.NewSearch, ShouldResemble, &testTime)
			So(err, ShouldBeNil)
		})

		Convey("returns error if no cookie is set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			cookie, err := GetABTest(req)
			So(cookie, ShouldResemble, ABServices{})
			So(err, ShouldNotBeNil)
		})
	})

	Convey("SetABTest sets correct cookie", t, func() {
		rec := httptest.NewRecorder()
		SetABTest(rec, testServices, testDomain)
		cookie := rec.Result().Cookies()[0]
		So(cookie.Value, ShouldEqual, expectedEscapedCookieValue)
		So(cookie.Path, ShouldEqual, "/")
		So(cookie.Domain, ShouldEqual, testDomain)
		So(cookie.MaxAge, ShouldEqual, maxAgeOneYear)
		So(cookie.Secure, ShouldBeTrue)
		So(cookie.SameSite, ShouldEqual, http.SameSiteLaxMode)
	})

	Convey("UpdateNewSearch updates cookie correctly", t, func() {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		req.AddCookie(&http.Cookie{Name: aBTestKey, Value: testEscapedValue})
		newValue := time.Date(2020, 10, 12, 17, 26, 43, 651387237, time.UTC)
		expectedNewValue := "%7B%22new_search%22%3A%222020-10-12T17%3A26%3A43.651387237Z%22%7D"

		UpdateNewSearch(req, rec, newValue, testDomain)
		cookie := rec.Result().Cookies()[0]
		So(cookie.Value, ShouldEqual, expectedNewValue)
	})

	Convey("UpdateOldSearch updates cookie correctly", t, func() {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/", nil)
		mockServices := ABServices{OldSearch: &testTime}
		mockMarshalledValue, err := json.Marshal(mockServices)
		if err != nil {
			t.Error("failed to marshal test data", err)
		}
		mockEscapedValue := url.QueryEscape(string(mockMarshalledValue))
		req.AddCookie(&http.Cookie{Name: aBTestKey, Value: mockEscapedValue})
		newValue := time.Date(2021, 9, 11, 15, 26, 43, 651387237, time.UTC)
		expectedNewValue := "%7B%22old_search%22%3A%222021-09-11T15%3A26%3A43.651387237Z%22%7D"

		UpdateOldSearch(req, rec, newValue, testDomain)
		cookie := rec.Result().Cookies()[0]
		So(cookie.Value, ShouldEqual, expectedNewValue)
	})
}

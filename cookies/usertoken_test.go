package cookies

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/davecgh/go-spew/spew"
	. "github.com/smartystreets/goconvey/convey"
)

var testDomain = "www.test.com"
var testAccessToken = "test-access-token"

func TestUnitGetUserAuthToken(t *testing.T) {

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
		SetUserAuthToken(rec, testAccessToken, testDomain)
		cookie := rec.Result().Cookies()[0]
		spew.Dump(cookie)
		So(cookie.Value, ShouldEqual, testAccessToken)
		So(cookie.Path, ShouldEqual, "/")
		So(cookie.Domain, ShouldEqual, testDomain)
		So(cookie.MaxAge, ShouldBeLessThanOrEqualTo, 0)
		So(cookie.Secure, ShouldBeTrue)
		So(cookie.SameSite, ShouldEqual, http.SameSiteLaxMode)
	})

}

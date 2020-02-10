package cookies

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var testDomain = "www.test.com"
var testCookie = "test_cookie"
var testValue = "test-value"

func TestUnitSetCookie(t *testing.T) {

	Convey("SetUserAuthToken sets correct cookie", t, func() {
		rec := httptest.NewRecorder()
		set(rec, testCookie, testValue, testDomain, 12)
		cookie := rec.Result().Cookies()[0]
		So(cookie.Value, ShouldEqual, testValue)
		So(cookie.Path, ShouldEqual, "/")
		So(cookie.Domain, ShouldEqual, testDomain)
		So(cookie.MaxAge, ShouldEqual, 12)
		So(cookie.Secure, ShouldBeTrue)
		So(cookie.SameSite, ShouldEqual, http.SameSiteLaxMode)
	})
}

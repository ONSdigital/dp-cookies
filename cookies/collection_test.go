package cookies

import (
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var testCollectionID = "test-123456789"

func TestUnitCollection(t *testing.T) {

	Convey("GetUserAuthToken", t, func() {
		Convey("returns cookie value if value is set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			req.AddCookie(&http.Cookie{Name: collectionIDCookieKey, Value: testCollectionID})
			cookie, _ := GetCollection(req)
			So(cookie, ShouldEqual, testCollectionID)
		})

		Convey("returns error if no cookie is set", func() {
			req := httptest.NewRequest("GET", "/", nil)
			_, err := GetCollection(req)
			So(err, ShouldNotBeNil)
		})
	})

	Convey("SetUserAuthToken sets correct cookie", t, func() {
		rec := httptest.NewRecorder()
		SetCollection(rec, testCollectionID, testDomain)
		cookie := rec.Result().Cookies()[0]
		So(cookie.Value, ShouldEqual, testCollectionID)
		So(cookie.Path, ShouldEqual, "/")
		So(cookie.Domain, ShouldEqual, testDomain)
		So(cookie.MaxAge, ShouldBeLessThanOrEqualTo, 0)
		So(cookie.Secure, ShouldBeTrue)
		So(cookie.SameSite, ShouldEqual, http.SameSiteLaxMode)
	})
}

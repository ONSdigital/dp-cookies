package cookies

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

const (
	testAspectID       = "test-aspect"
	testSecondAspectID = "second-aspect"
	testDomain         = "test-domain"
	oldHandlerServed   = "Old Served"
	newHandlerServed   = "New Served"
)

func TestGetABTestCookieAspect(t *testing.T) {
	Convey("Given a http abTestCookie does not exist", t, func() {
		req := httptest.NewRequest("GET", "/", http.NoBody)

		Convey("When GetABTestCookieAspect is called with the requested aspect ID", func() {
			aspect := GetABTestCookieAspect(req, testAspectID)

			Convey("A zero value ABTestCookieAspect is returned", func() {
				So(aspect, ShouldBeZeroValue)
			})
		})
	})

	Convey("Given a http abTestCookie exists, but without the requested aspect ID", t, func() {
		c := &http.Cookie{
			Name:  aBTestKey,
			Value: url.QueryEscape(fmt.Sprintf(`{%q:{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`, testSecondAspectID)),
		}
		req := httptest.NewRequest("GET", "/", http.NoBody)
		req.AddCookie(c)

		Convey("When GetABTestCookieAspect is called with the requested aspect ID", func() {
			aspect := GetABTestCookieAspect(req, testAspectID)

			Convey("A zero value ABTestCookieAspect is returned", func() {
				So(aspect, ShouldBeZeroValue)
			})
		})
	})

	Convey("Given a http abTestCookie exists with the requested aspect ID", t, func() {
		c := &http.Cookie{
			Name:  aBTestKey,
			Value: url.QueryEscape(fmt.Sprintf(`{%q:{"new":"2020-06-16T17:28:45","old":"2020-06-15T17:28:45"},%q:{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`, testAspectID, testSecondAspectID)),
		}
		req := httptest.NewRequest("GET", "/", http.NoBody)
		req.AddCookie(c)

		Convey("When GetABTestCookieAspect is called with the requested aspect ID", func() {
			aspect := GetABTestCookieAspect(req, testAspectID)

			Convey("The requested ABTestCookieAspect is returned", func() {
				So(aspect, ShouldResemble, ABTestCookieAspect{New: MustParseCookieTime("2020-06-15T17:28:45").Add(time.Hour * 24), Old: MustParseCookieTime("2020-06-15T17:28:45")})
			})
		})
	})
}

func TestSetABTestCookieAspect(t *testing.T) {
	Convey("Given a http abTestCookie does not exist", t, func() {
		req := httptest.NewRequest("GET", "/", http.NoBody)
		rec := httptest.NewRecorder()

		Convey("When SetABTestCookieAspect is called with a valid Aspect", func() {
			aspect := ABTestCookieAspect{New: MustParseCookieTime("2020-06-15T17:28:45").Add(time.Hour * 24), Old: MustParseCookieTime("2020-06-15T17:28:45")}
			SetABTestCookieAspect(rec, req, testAspectID, testDomain, aspect)

			Convey("The http abTestCookie is created and set, and contains the new aspect", func() {
				So(rec.Result().Cookies(), ShouldHaveLength, 1)
				cookieSetInResponse := rec.Result().Cookies()[0]

				So(cookieSetInResponse.Value, ShouldEqual, url.QueryEscape(fmt.Sprintf(`{%q:{"new":"2020-06-16T17:28:45","old":"2020-06-15T17:28:45"}}`, testAspectID)))
				So(cookieSetInResponse.Path, ShouldEqual, "/")
				So(cookieSetInResponse.Domain, ShouldEqual, testDomain)
				So(cookieSetInResponse.MaxAge, ShouldEqual, maxAgeOneYear)
				So(cookieSetInResponse.SameSite, ShouldEqual, http.SameSiteLaxMode)
			})
		})
	})

	Convey("Given a http abTestCookie does exist", t, func() {
		c := &http.Cookie{
			Name:  aBTestKey,
			Value: url.QueryEscape(fmt.Sprintf(`{%q:{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`, testSecondAspectID)),
		}
		req := httptest.NewRequest("GET", "/", http.NoBody)
		req.AddCookie(c)
		rec := httptest.NewRecorder()

		Convey("When SetABTestCookieAspect is called with a valid Aspect that does NOT currently exist in the http abTestCookie", func() {
			aspect := ABTestCookieAspect{New: MustParseCookieTime("2020-06-15T17:28:45").Add(time.Hour * 24), Old: MustParseCookieTime("2020-06-15T17:28:45")}
			SetABTestCookieAspect(rec, req, testAspectID, testDomain, aspect)

			Convey("The new aspect is added to the http abTestCookie", func() {
				So(rec.Result().Cookies(), ShouldHaveLength, 1)
				cookieSetInResponse := rec.Result().Cookies()[0]

				So(cookieSetInResponse.Value, ShouldContainSubstring, url.QueryEscape(fmt.Sprintf(`%q:{"new":"2020-06-16T17:28:45","old":"2020-06-15T17:28:45"}`, testAspectID)))
				So(cookieSetInResponse.Value, ShouldContainSubstring, url.QueryEscape(fmt.Sprintf(`%q:{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}`, testSecondAspectID)))
				So(cookieSetInResponse.Path, ShouldEqual, "/")
				So(cookieSetInResponse.Domain, ShouldEqual, testDomain)
				So(cookieSetInResponse.MaxAge, ShouldEqual, maxAgeOneYear)
				So(cookieSetInResponse.SameSite, ShouldEqual, http.SameSiteLaxMode)
			})
		})

		Convey("When SetABTestCookieAspect is called with a valid Aspect that does currently exist in the http abTestCookie", func() {
			aspect := ABTestCookieAspect{New: MustParseCookieTime("2022-01-12T17:28:45").Add(time.Hour * 24), Old: MustParseCookieTime("2022-01-12T17:28:45")}
			SetABTestCookieAspect(rec, req, testSecondAspectID, testDomain, aspect)

			Convey("The existing aspect is overwritten in the http abTestCookie", func() {
				So(rec.Result().Cookies(), ShouldHaveLength, 1)
				cookieSetInResponse := rec.Result().Cookies()[0]

				So(cookieSetInResponse.Value, ShouldContainSubstring, url.QueryEscape(fmt.Sprintf(`%q:{"new":"2022-01-13T17:28:45","old":"2022-01-12T17:28:45"}`, testSecondAspectID)))
				So(cookieSetInResponse.Path, ShouldEqual, "/")
				So(cookieSetInResponse.Domain, ShouldEqual, testDomain)
				So(cookieSetInResponse.MaxAge, ShouldEqual, maxAgeOneYear)
				So(cookieSetInResponse.SameSite, ShouldEqual, http.SameSiteLaxMode)
			})
		})

	})
}

func TestRemoveABTestCookieAspect(t *testing.T) {
	Convey("Given a http abTestCookie does not exist", t, func() {
		req := httptest.NewRequest("GET", "/", http.NoBody)
		rec := httptest.NewRecorder()

		Convey("Calling RemoveABTestCookieAspect is a noop", func() {
			RemoveABTestCookieAspect(rec, req, testAspectID, testDomain)
		})
	})

	Convey("Given a http abTestCookie exists, but without the requested aspect", t, func() {
		c := &http.Cookie{
			Name:  aBTestKey,
			Value: url.QueryEscape(fmt.Sprintf(`{%q:{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`, testSecondAspectID)),
		}
		req := httptest.NewRequest("GET", "/", http.NoBody)
		req.AddCookie(c)
		rec := httptest.NewRecorder()

		Convey("Calling RemoveABTestCookieAspect is a noop", func() {
			RemoveABTestCookieAspect(rec, req, testAspectID, testDomain)

			So(rec.Result().Cookies(), ShouldHaveLength, 1)
			So(rec.Result().Cookies()[0].Name, ShouldEqual, c.Name)
			So(rec.Result().Cookies()[0].Value, ShouldEqual, c.Value)
		})
	})

	Convey("Given a http abTestCookie exists with the requested aspect", t, func() {
		c := &http.Cookie{
			Name:  aBTestKey,
			Value: url.QueryEscape(fmt.Sprintf(`{%q:{"new":"2020-06-16T17:28:45","old":"2020-06-15T17:28:45"},%q:{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`, testAspectID, testSecondAspectID)),
		}
		req := httptest.NewRequest("GET", "/", http.NoBody)
		req.AddCookie(c)
		rec := httptest.NewRecorder()

		Convey("Calling RemoveABTestCookieAspect removes the requested aspect from the http abTestCookie", func() {
			RemoveABTestCookieAspect(rec, req, testAspectID, testDomain)

			So(rec.Result().Cookies(), ShouldHaveLength, 1)
			So(rec.Result().Cookies()[0].Name, ShouldEqual, c.Name)
			So(rec.Result().Cookies()[0].Value, ShouldEqual, url.QueryEscape(fmt.Sprintf(`{%q:{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`, testSecondAspectID)))
		})
	})
}

func TestServABTest(t *testing.T) {
	Convey("Given an old request handler and a new request handler", t, func() {
		oldHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(oldHandlerServed)) })
		newHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(newHandlerServed)) })
		req := httptest.NewRequest("GET", "/", http.NoBody)

		Convey("When handling a test aspect where the 'new' aspect time is in the future", func() {
			testAspect := ABTestCookieAspect{New: Now().Add(time.Hour * 24)}

			Convey("The new request handler is used to handle the request regardless of the value of the 'old' aspect time", func() {
				rec := httptest.NewRecorder()
				ServABTest(rec, req, newHandler, oldHandler, testAspect)
				b, _ := io.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte(newHandlerServed))

				testAspect.Old = Now()
				rec = httptest.NewRecorder()
				ServABTest(rec, req, newHandler, oldHandler, testAspect)
				b, _ = io.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte(newHandlerServed))

				testAspect.Old = Now().Add(time.Hour * 24)
				rec = httptest.NewRecorder()
				ServABTest(rec, req, newHandler, oldHandler, testAspect)
				b, _ = io.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte(newHandlerServed))
			})
		})

		Convey("When handling a test aspect where the 'old' aspect time is in the future", func() {
			testAspect := ABTestCookieAspect{Old: Now().Add(time.Hour * 24)}

			Convey("The old request handler is used to handle the request only when the new aspect time is in the past", func() {
				rec := httptest.NewRecorder()
				ServABTest(rec, req, newHandler, oldHandler, testAspect)
				b, _ := io.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte(oldHandlerServed))

				testAspect.New = Now()
				rec = httptest.NewRecorder()
				ServABTest(rec, req, newHandler, oldHandler, testAspect)
				b, _ = io.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte(oldHandlerServed))

				testAspect.New = Now().Add(time.Hour * 24)
				rec = httptest.NewRecorder()
				ServABTest(rec, req, newHandler, oldHandler, testAspect)
				b, _ = io.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte(newHandlerServed))
			})
		})
	})
}

func TestHandleCookieAndServ(t *testing.T) {
	Convey("Given an old request handler, new request handler and a test randomiser", t, func() {
		testOld := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(oldHandlerServed)) })
		testNew := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(newHandlerServed)) })
		testAspect := ABTestCookieAspect{New: Now().Add(time.Hour * 24), Old: Now()}
		testRandomiser := func() ABTestCookieAspect { return testAspect }

		Convey("When HandleCookieAndServ is called with an aspectID", func() {
			req := httptest.NewRequest("GET", "/", http.NoBody)
			rec := httptest.NewRecorder()
			HandleCookieAndServ(rec, req, testNew, testOld, testAspectID, testDomain, testRandomiser)

			Convey("The the new aspect is set in the http abTestCookie", func() {
				So(rec.Result().Cookies(), ShouldHaveLength, 1)
				cookieSetInResponse := rec.Result().Cookies()[0]

				testAspectMarshalled, err := json.Marshal(abTestCookie{testAspectID: testAspect})
				if err != nil {
					t.Fatalf("error marshalling abTestCookie: %v", err)
				}
				So(cookieSetInResponse.Value, ShouldEqual, url.QueryEscape(string(testAspectMarshalled)))
			})

			Convey("And the appropriate handler serves the request", func() {
				b, _ := io.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte(newHandlerServed))
			})
		})
	})
}

func TestHandleABTestExit(t *testing.T) {
	Convey("Given an old request handler", t, func() {
		testOld := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(oldHandlerServed)) })

		Convey("When HandleABTestExit is called with an aspectID", func() {
			req := httptest.NewRequest("GET", "/", http.NoBody)
			rec := httptest.NewRecorder()
			HandleABTestExit(rec, req, testOld, testAspectID, testDomain)

			Convey("A new aspect for the requested aspect ID is set in the http abTestCookie", func() {
				So(rec.Result().Cookies(), ShouldHaveLength, 1)
				cookieSetInResponse := rec.Result().Cookies()[0]
				So(cookieSetInResponse.Value, ShouldContainSubstring, url.QueryEscape(fmt.Sprintf(`{%q:{"new":"`, testAspectID)))
			})

			Convey("And the old request handler serves the request", func() {
				b, _ := io.ReadAll(rec.Result().Body)
				So(b, ShouldResemble, []byte(oldHandlerServed))
			})
		})
	})
}

func TestGetABTestCookie(t *testing.T) {
	Convey("given a valid http abTestCookie exists", t, func() {
		c := &http.Cookie{
			Name:  aBTestKey,
			Value: url.QueryEscape(fmt.Sprintf(`{%q:{"new":"2020-06-16T17:28:45","old":"2020-06-15T17:28:45"},%q:{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`, testAspectID, testSecondAspectID)),
		}
		req := httptest.NewRequest("GET", "/", http.NoBody)
		req.AddCookie(c)

		Convey("when getABTestCookie is called, the expected abTestCookie is returned", func() {
			receivedCookie, err := getABTestCookie(req)
			So(err, ShouldBeNil)
			So(receivedCookie, ShouldResemble, abTestCookie{
				testAspectID:       ABTestCookieAspect{New: MustParseCookieTime("2020-06-15T17:28:45").Add(time.Hour * 24), Old: MustParseCookieTime("2020-06-15T17:28:45")},
				testSecondAspectID: ABTestCookieAspect{New: MustParseCookieTime("2021-12-31T09:30:00"), Old: MustParseCookieTime("2021-12-31T09:30:00").Add(time.Hour * 24)}})
		})
	})

	Convey("given a valid abTestCookie does not exist", t, func() {
		c := &http.Cookie{Name: "some-cookie", Value: "some value"}
		req := httptest.NewRequest("GET", "/", http.NoBody)
		req.AddCookie(c)

		Convey("when getABTestCookie is called, an error is returned", func() {
			receivedCookie, err := getABTestCookie(req)
			So(err, ShouldEqual, ErrABTestCookieNotFound)
			So(receivedCookie, ShouldResemble, abTestCookie{})
		})
	})
}

func TestSetABTestCookie(t *testing.T) {
	Convey("given a valid abTestCookie", t, func() {
		c := abTestCookie{
			testAspectID:       ABTestCookieAspect{New: MustParseCookieTime("2020-06-15T17:28:45").Add(time.Hour * 24), Old: MustParseCookieTime("2020-06-15T17:28:45")},
			testSecondAspectID: ABTestCookieAspect{New: MustParseCookieTime("2021-12-31T09:30:00"), Old: MustParseCookieTime("2021-12-31T09:30:00").Add(time.Hour * 24)},
		}
		rec := httptest.NewRecorder()

		Convey("when setABTestCookie is called, the corresponding http cookie is correctly set in a http Response", func() {
			err := setABTestCookie(rec, c, testDomain)
			So(err, ShouldBeNil)
			So(rec.Result().Cookies(), ShouldHaveLength, 1)

			cookieSetInResponse := rec.Result().Cookies()[0]
			So(cookieSetInResponse.Value, ShouldContainSubstring, url.QueryEscape(fmt.Sprintf(`%q:{"new":"2020-06-16T17:28:45","old":"2020-06-15T17:28:45"}`, testAspectID)))
			So(cookieSetInResponse.Value, ShouldContainSubstring, url.QueryEscape(fmt.Sprintf(`%q:{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}`, testSecondAspectID)))
			So(cookieSetInResponse.Path, ShouldEqual, "/")
			So(cookieSetInResponse.Domain, ShouldEqual, testDomain)
			So(cookieSetInResponse.MaxAge, ShouldEqual, maxAgeOneYear)
			So(cookieSetInResponse.SameSite, ShouldEqual, http.SameSiteLaxMode)
		})
	})
}

type user struct {
	new    bool
	cookie http.Cookie
}

func TestABTestHandler(t *testing.T) {
	Convey("Given an old request handler and a new request handler", t, func() {
		oldHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(oldHandlerServed))
		})
		newHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			_, _ = w.Write([]byte(newHandlerServed))
		})

		percentage := 40
		numberRequests := 500

		Convey("When the requests are made to the abTestHandler", func() {
			handler := abTestHandler(newHandler, oldHandler, percentage, testAspectID, "my-domain", "exit-new-test")
			users := make([]user, numberRequests)
			nh := 0

			for i := 0; i < numberRequests; i++ {
				req := httptest.NewRequest("GET", "/", http.NoBody)
				resp := httptest.NewRecorder()
				handler.ServeHTTP(resp, req)
				So(resp.Result().Cookies(), ShouldHaveLength, 1)
				users[i] = user{cookie: *resp.Result().Cookies()[0]}

				b, _ := io.ReadAll(resp.Result().Body)
				if string(b) == newHandlerServed {
					users[i].new = true
					nh++
				}
			}

			Convey("The ABTest cookies are analysed (and stored) to ensure the number of requests serviced by the new handler are within an acceptable deviation range of the requested percentage split))", func() {
				onePercent := numberRequests / 100
				expected := percentage * onePercent
				deviation := func(x int) int {
					if x < 0 {
						return -x
					}
					return x
				}(expected - nh)
				So(deviation, ShouldBeBetweenOrEqual, 0, 20*onePercent)
			})

			Convey("When subsequent requests are made with the previously returned cookies, all are serviced by the correct handler", func() {
				for i := 0; i < numberRequests; i++ {
					req := httptest.NewRequest("GET", "/", http.NoBody)
					req.AddCookie(&users[i].cookie)
					resp := httptest.NewRecorder()
					handler.ServeHTTP(resp, req)

					b, _ := io.ReadAll(resp.Result().Body)
					expectedResponse := oldHandlerServed
					if users[i].new {
						expectedResponse = newHandlerServed
					}
					So(string(b), ShouldEqual, expectedResponse)
				}
			})
			Convey("When subsequent requests with 'exit-new-test' query parameter are made, they are serviced by the old handler", func() {
				for i := 0; i < numberRequests; i++ {
					if !users[i].new {
						continue
					}

					req := httptest.NewRequest("GET", "/?exit-new-test", http.NoBody)
					req.AddCookie(&users[i].cookie)
					resp := httptest.NewRecorder()
					handler.ServeHTTP(resp, req)

					b, _ := io.ReadAll(resp.Result().Body)
					So(string(b), ShouldEqual, oldHandlerServed)
				}
			})
		})
	})
}

func TestABTestPurgeHandler(t *testing.T) {
	Convey("Given a new request handler", t, func() {
		newHandler := http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) { _, _ = w.Write([]byte(newHandlerServed)) })

		Convey("And an ab_test cookie containing a now unused aspect for ab testing of the new/old handler", func() {
			c := &http.Cookie{
				Name:  aBTestKey,
				Value: url.QueryEscape(fmt.Sprintf(`{%q:{"new":"2020-06-16T17:28:45","old":"2020-06-15T17:28:45"},%q:{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`, testAspectID, testSecondAspectID)),
			}

			Convey("When a request with the cookie is made to the abTestPurgeHandler", func() {
				req := httptest.NewRequest("GET", "/", http.NoBody)
				req.AddCookie(c)
				w := httptest.NewRecorder()
				abTestPurgeHandler(newHandler, testAspectID, testDomain).ServeHTTP(w, req)

				Convey("The relevant aspect is removed from the ab_test cookie, and the request has been handled by the new handler", func() {
					So(w.Result().Cookies(), ShouldHaveLength, 1)
					So(w.Result().Cookies()[0].Name, ShouldEqual, aBTestKey)
					So(w.Result().Cookies()[0].Value, ShouldEqual, url.QueryEscape(fmt.Sprintf(`{%q:{"new":"2021-12-31T09:30:00","old":"2022-01-01T09:30:00"}}`, testSecondAspectID)))

					b, _ := io.ReadAll(w.Result().Body)
					So(string(b), ShouldEqual, newHandlerServed)
				})
			})
		})
	})
}

func TestDefaultABTestRandomiserAt100(t *testing.T) {
	var iterations, n = 100, 0

	randomiser := DefaultABTestRandomiser(100)
	for i := 0; i < iterations; i++ {
		a := randomiser()
		if a.New.After(time.Now()) {
			n++
		}
	}

	if n != iterations {
		t.Errorf("a percentage of 100%% requires ALL generated aspects to favour the New value. Expected: %d Got %d", iterations, n)
	}
}

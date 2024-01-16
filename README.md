# dp-cookies

Golang library for setting/getting specific cookies

NOTE: If testing this library you will need to set the following local environmental variable LIBRARY_TEST:=true, or run `make test`.

## Setting a cookie using dp-cookies library

```go
package handler

import (
    ...
    "github.com/ONSdigital/dp-cookies/cookies"
    ...
)

// Set user auth token cookie

func myHandler(w http.ResponseWriter, req *http.Request, ...) {
    ...
    cookies.SetUserAuthToken(w, userAuthToken, "www.domain.com")
    ...
}
```

## Getting a cookie using dp-cookies library

```go
package handler

import (
    ...
    "github.com/ONSdigital/dp-cookies/cookies"
    ...
)

// Get user auth token value from cookie

func myHandler(w http.ResponseWriter, req *http.Request, ...) {
    ...
    token, err := cookies.GetUserAuthToken(req)
    ...
}
```

## Using the AB test handlerfunc

```go
package handler

import (
    ...
    "github.com/ONSdigital/dp-cookies/cookies"
    ...
)

func Read(...) http.HandlerFunc {
    oldHandler := ...http.HandlerFunc()
    newHandler := ...http.HandlerFunc()

    return cookies.Handler(cfg.ABTest.Enabled, newHandler, oldHandler, cfg.ABTest.Percentage, cfg.ABTest.AspectID, cfg.SiteDomain, cfg.ABTest.Exit)
}
```

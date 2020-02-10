package main

import (
	"fmt"
	"net/http"

	"github.com/ONSdigital/dp-cookies/cookies"
)

var policy = cookies.Policy{
	Essential: true,
	Usage:     true,
}

func main() {
	http.HandleFunc("/set-cookies", func(w http.ResponseWriter, r *http.Request) {
		cookies.SetPolicy(w, policy, "www.test.com")
		cookies.SetPreferenceIsSet(w, "www.test.com")
	})

	http.HandleFunc("/cookies", func(w http.ResponseWriter, r *http.Request) {
		cookiesResponse := cookies.GetCookiePreferences(r)
		w.Write([]byte(fmt.Sprintf("%+v \n", cookiesResponse)))
	})

	http.ListenAndServe(":22888", nil)
}

package main

import (
	"net/http"
)

func WithAuth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		if cookie, err := r.Cookie("twauth"); err == http.ErrNoCookie || cookie.Value == "" {
			// user is not authenticated
			w.Header().Set("Location", "/login")
			w.WriteHeader(http.StatusTemporaryRedirect)
		} else if err != nil {
			// some other error
			panic(err.Error())
		} else {
			// success, proceed to original Handler
			fn(w, r)
		}
		fn(w, r)
	}
}
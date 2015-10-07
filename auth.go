package main

import (
	"net/http"
	"github.com/go-zoo/bone"
)

// WithAuth function checks whether cookie exists at all. It doesn't check whether it is valid or not
// valid or not is decided by actual handler that tries to decode it.
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
	}
}


// GetProvider function adds additional parameter provider to URL query since Gothic expects it there
// https://github.com/markbates/goth/blob/master/gothic/gothic.go#L148-L160, this is just a function wrapper.
// if switched to other mux - there wouldn't be a need for this
func GetProvider(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		provider := bone.GetValue(r, "provider")
		values := r.URL.Query()
		values.Add("provider", provider)
		r.URL.RawQuery = values.Encode()
		fn(w, r)
	}
}
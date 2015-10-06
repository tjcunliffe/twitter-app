package main

import (
	"net/http"
	"fmt"
	"html/template"
	"github.com/markbates/goth/gothic"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/securecookie"
)

var hashGenKeySecret = securecookie.GenerateRandomKey(16)
var blockGenKey = securecookie.GenerateRandomKey(16)

// setting secure cookie instance
var hashKey = []byte(hashGenKeySecret)
var blockKey = []byte(blockGenKey)
var s = securecookie.New(hashKey, blockKey)


func homeHandler(w http.ResponseWriter, r *http.Request){

	if cookie, err := r.Cookie("twauth"); err == nil {
		log.Info("Cokkie found, decoding...")
		value := make(map[string]string)
		err := s.Decode("twauth", cookie.Value, &value)
		if err == nil {
			fmt.Fprintf(w, "The value of twtoken is %q", value["twtoken"])
			fmt.Fprintf(w, "The value of twtokensecret is %q", value["twtokensecret"])
		} else {
			log.Error("failed to decode cookie", err)
			// not valid, probably wrong key
			// removing cookie and redirecting to login page
			http.SetCookie(w, &http.Cookie{
				Name: "twauth",
				Value: "",
				Path: "/",
				MaxAge: -1,
			})
			w.Header()["Location"] = []string{"/login"}
			w.WriteHeader(http.StatusTemporaryRedirect)
		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request){
	t, _ := template.New("foo").Parse(indexTemplate)
	t.Execute(w, nil)
}


func callBackHandler(w http.ResponseWriter, r *http.Request) {
	user, err := gothic.CompleteUserAuth(w, r)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	// setting cookie value
	value := map[string]string{
		"twtoken": user.AccessToken,
		"twtokensecret": user.AccessTokenSecret,
	}
	log.Info("got details, encoding cookie")
	// encoding and setting cookie
	encoded, err := s.Encode("twauth", value)
	if err == nil {
		cookie := &http.Cookie{
			Name:  "twauth",
			Value: encoded,
			Path:  "/",
		}
		http.SetCookie(w, cookie)
		log.Info("Cookie set")
	} else {
		log.Error("Failed to write cookie", err)
	}

	t, _ := template.New("foo").Parse(userTemplate)
	t.Execute(w, user)
}

var indexTemplate = `
<p><a href="/auth/twitter">Log in with Twitter</a></p>
`

var userTemplate = `
<p>Name: {{.Name}}</p>
<p>Email: {{.Email}}</p>
<p>NickName: {{.NickName}}</p>
<p>Location: {{.Location}}</p>
<p>AvatarURL: {{.AvatarURL}} <img src="{{.AvatarURL}}"></p>
<p>Description: {{.Description}}</p>
<p>UserID: {{.UserID}}</p>
<p>AccessToken: {{.AccessToken}}</p>
`

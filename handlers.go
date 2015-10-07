package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	"html/template"
	"github.com/markbates/goth/gothic"
	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/securecookie"
	"github.com/unrolled/render"

	"encoding/json"
)

// SearchQuery is a struct of json body that is expected to
// come from the client app and holds both query and tokens
type SearchQuery struct {
	Query string `json:query`
	Token string `json:token`
	Secret string `json:secret`
}

var hashGenKeySecret = securecookie.GenerateRandomKey(16)
var blockGenKey = securecookie.GenerateRandomKey(16)

// setting secure cookie instance
var hashKey = []byte(hashGenKeySecret)
var blockKey = []byte(blockGenKey)
var s = securecookie.New(hashKey, blockKey)

type Handler struct {
	r *render.Render
}

func (h *HTTPClientHandler) homeHandler(w http.ResponseWriter, r *http.Request){

	if cookie, err := r.Cookie("twauth"); err == nil {
		log.Info("Cookie found, decoding...")
		value := make(map[string]string)
		err := s.Decode("twauth", cookie.Value, &value)
		if err == nil {
			newmap := map[string]interface{}{"metatitle": "Tweets", "token": value["twtoken"]}
			h.r.HTML(w, http.StatusOK, "home", newmap)

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


// loginHandler presents initial template for logging in
func loginHandler(w http.ResponseWriter, r *http.Request){
	t, _ := template.New("foo").Parse(indexTemplate)
	t.Execute(w, nil)
}


// callBackHandler does final authentication step for user and records a cookie that will be stored in user's
// browser for later decoding and reusing of auth tokens
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


func (h *HTTPClientHandler) searchTwitter(w http.ResponseWriter, r *http.Request){

	defer r.Body.Close()

	// reading resposne body
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		// logging read error
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("Failed to read request body!")
	}

	var data SearchQuery
	err = json.Unmarshal(body, &data)

	if err != nil{
		// logging read error
		log.WithFields(log.Fields{
			"error": err.Error(),
		}).Warn("Failed to unmarshall json body!")
	}

	// parameters
	token := data.Token
	secret := data.Secret
	query := data.Query

	fmt.Println(token)
	fmt.Println(secret)
	fmt.Println(query)
}
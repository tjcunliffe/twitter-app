package main

import (
	"net/http"
	"fmt"
	"html/template"
	"github.com/markbates/goth/gothic"
)

func homeHandler(w http.ResponseWriter, r *http.Request){
	fmt.Fprintf(w, "Welcome to the home page!")
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

	// setting twauth cookie for later usage
	http.SetCookie(w, &http.Cookie{
		Name:  "twauth",
		Value: user.AccessToken,
		Path:  "/"})


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

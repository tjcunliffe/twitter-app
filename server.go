package main

import (
	"github.com/codegangsta/negroni"
	"net/http"
	"fmt"
	"flag"
	"os"
	"encoding/json"
	"html/template"
	log "github.com/Sirupsen/logrus"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
	"github.com/gorilla/pat"
)

type Configuration struct {
	TwitterKey string
	TwitterSecret string
}

var AppConfig Configuration

func main() {
	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{})


	// getting configuration
	file, err := os.Open("conf.json")
	if err != nil {
		log.Panic("Failed to open configuration file, quiting server.")
	}
	decoder := json.NewDecoder(file)
	AppConfig = Configuration{}
	err = decoder.Decode(&AppConfig)
	if err != nil {
		log.WithFields(log.Fields{"Error": err.Error()}).Panic("Failed to read configuration")
	}

	// app starting
	log.WithFields(log.Fields{
		"Key": AppConfig.TwitterKey,
		"Secret": AppConfig.TwitterSecret,
	}).Info("app is starting")

	// initialising goth provider
	goth.UseProviders(
		twitter.New(AppConfig.TwitterKey, AppConfig.TwitterSecret, "http://localhost:8080/auth/twitter/callback"))

	// looking for option args when starting App
	// like ./twitter-app -port=":8080" would start on port 8080
	var port = flag.String("port", ":8080", "Server port")
	flag.Parse() // parse the flag

	mux := getRouter()
	n := negroni.Classic()
	n.UseHandler(mux)
	n.Run(*port)
}

func getRouter() *pat.Router {
	p := pat.New()
	p.Get("/auth/{provider}/callback", callBackHandler)
	p.Get("/auth/{provider}", gothic.BeginAuthHandler)
	p.Get("/login", loginHandler)
	p.Get("/", WithAuth(homeHandler))
	return p
}

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
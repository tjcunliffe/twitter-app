package main

import (
	"github.com/codegangsta/negroni"
	"flag"
	"os"
	"encoding/json"
	"net/http"
	log "github.com/Sirupsen/logrus"

	"github.com/markbates/goth"
	"github.com/markbates/goth/gothic"
	"github.com/markbates/goth/providers/twitter"
	"github.com/unrolled/render"
	"github.com/meatballhat/negroni-logrus"
	"github.com/go-zoo/bone"
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

	// getting base template and handler struct
	r := render.New(render.Options{Layout: "layout"})
	h := Handler{r: r}

	mux := getBoneRouter(h)
	n := negroni.Classic()
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(mux)
	n.Run(*port)
}


func getBoneRouter(h Handler) *bone.Mux {
	mux := bone.New()
	mux.Get("/auth/:provider/callback", GetProvider(http.HandlerFunc(callBackHandler)))
	mux.Get("/auth/:provider", GetProvider(http.HandlerFunc(gothic.BeginAuthHandler)))
	mux.Get("/login", http.HandlerFunc(loginHandler))
	mux.Get("/", http.HandlerFunc(WithAuth(h.homeHandler)))
	// handling static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return mux
}


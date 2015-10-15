package main

import (
	"github.com/codegangsta/negroni"
	"flag"
	"os"
	"encoding/json"
	"net/http"
	log "github.com/Sirupsen/logrus"

	"github.com/unrolled/render"
	"github.com/meatballhat/negroni-logrus"
	"github.com/go-zoo/bone"
)

// Initial structure of configuration that is expected from conf.json file
type Configuration struct {
	TwitterKey string
	TwitterSecret string
	MirageProxy string
}

// AppConfig stores application configuration
var AppConfig Configuration

// Client structure to be injected into functions to perform HTTP calls
type Client struct {
	HTTPClient *http.Client
}

// HTTPClientHandler used for passing http client connection and template
// information back to handlers, mostly for testing purposes
type HTTPClientHandler struct {
	http Client
	r  *render.Render
}

var (
	// like ./twitter-app -port=":8080" would start on port 8080
	port = flag.String("port", ":8080", "Server port")
)


func main() {
	// Output to stderr instead of stdout, could also be a file.
	log.SetOutput(os.Stderr)
	log.SetFormatter(&log.TextFormatter{})

	// getting app config
	twitterKey := os.Getenv("TwitterKey")
	twitterSecret := os.Getenv("TwitterSecret")


	if(twitterKey != "" && twitterSecret != ""){
		log.Info("Environment variables for Twitter authentication found.")
		AppConfig.TwitterKey = twitterKey
		AppConfig.TwitterSecret = twitterSecret
	} else {
		log.Info("Environment variables for Twitter authentication not found, looking for configuration file")
		// getting configuration from file
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
	}

	mirageProxy := os.Getenv("MirageProxyAddress")
	if(mirageProxy == ""){
		log.Info("MirageProxyAddress environment variable not found, using default - http://localhost:8300")
		AppConfig.MirageProxy = "http://localhost:8300"
	} else {
		log.Info("MirageProxyAddress environment variable found!")
		AppConfig.MirageProxy = mirageProxy
	}

	// app starting
	log.WithFields(log.Fields{
		"Key": AppConfig.TwitterKey,
		"Secret": AppConfig.TwitterSecret,
	}).Info("app is starting")

	// looking for option args when starting App
	flag.Parse() // parse the flag

	// getting base template and handler struct
	r := render.New(render.Options{Layout: "layout"})

	h := HTTPClientHandler{http: Client{&http.Client{}}, r: r}

	mux := getBoneRouter(h)
	n := negroni.Classic()
	n.Use(negronilogrus.NewMiddleware())
	n.UseHandler(mux)
	n.Run(*port)
}


func getBoneRouter(h HTTPClientHandler) *bone.Mux {
	mux := bone.New()
	mux.Get("/query", http.HandlerFunc(h.queryTwitter))
	mux.Get("/", http.HandlerFunc(h.homeHandler))
	// handling static files
	mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	return mux
}


package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	"github.com/unrolled/render"
	"encoding/json"
)

// SearchQuery is a struct of json body that is expected to
// come from the client app and holds both query and tokens
type SearchQuery struct {
	Query  string `json:query`
	Token  string `json:token`
	Secret string `json:secret`
}

type ErrorResponse struct {
	Error    string `json:error`
}

type Handler struct {
	r *render.Render
}

func (h *HTTPClientHandler) homeHandler(w http.ResponseWriter, r *http.Request) {
	h.r.HTML(w, http.StatusOK, "home", nil)
}

func (h *HTTPClientHandler) queryTwitter(w http.ResponseWriter, r *http.Request) {

	// getting query
	urlQuery := r.URL.Query()
	// getting submitted query string
	queryString := urlQuery["q"]
	// getting actual twitter (or Mirage) backend URI
	backendUri := urlQuery["backend"]

	client := h.http.HTTPClient
	b := h.getBearerToken()

	twitterEndPoint := backendUri[0] + "/1.1/search/tweets.json?q=" + queryString[0]
	// logging full URL and path
	log.WithFields(log.Fields{
		"chosenBackend": backendUri[0],
		"query": queryString[0],
		"finalTwitterEndpoint": twitterEndPoint,
	}).Info("Endpoint created, performing query...")

	req, err := http.NewRequest("GET", twitterEndPoint, nil)
	if err != nil {
		log.Fatal(err)
	}

	//Step 3: Authenticate API requests with the bearer token
	//include an Authorization header formatted as
	//Bearer <bearer token value from step 2>
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", b.AccessToken))

	//Issue the request and get the JSON API response
	resp, err := client.Do(req)
	if err != nil {
		// logging error
		log.WithFields(log.Fields{
			"Error": err.Error(),
			"finalTwitterEndpoint": twitterEndPoint,
		}).Warn("Got error while querying external API...")

		// creating JSON response with error
		errString := err.Error()
		errorMsg := ErrorResponse{errString}
		js, err := json.Marshal(errorMsg)

		if err != nil {
			log.Error("Got error while marshalling error..")
			return
		}
        // writing response js
		w.Header().Set("Content-Type", "application/json")
		w.Write(js)

	} else {
		defer resp.Body.Close()
		// reading response body and returning it directly to the initial client, skipping decode step
		respBody, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			log.Fatal(err)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(respBody)
	}

}
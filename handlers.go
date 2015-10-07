package main

import (
	"net/http"
	"fmt"
	"io/ioutil"
	log "github.com/Sirupsen/logrus"
	"github.com/unrolled/render"
)

// SearchQuery is a struct of json body that is expected to
// come from the client app and holds both query and tokens
type SearchQuery struct {
	Query  string `json:query`
	Token  string `json:token`
	Secret string `json:secret`
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
	// getting session name
	queryString := urlQuery["q"]

	client := h.http.HTTPClient
	b := h.getBearerToken()

	twitterEndPoint := TwitterUri + "/1.1/search/tweets.json?q=" + queryString[0]
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
		log.Fatal(err, resp)
	}
	defer resp.Body.Close()
	// reading response body and returning it directly to the initial client, skipping decode step
	respBody, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(respBody)
}
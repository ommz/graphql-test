package service

import (
	"bytes"
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	repo "github.com/ommz/graphql-test/repository"
	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
)

const graphQLEndpoint = "https://gitlab.com/api/graphql"

// expected graphql response
type responseSchema struct {
	Data struct {
		Projects struct {
			Nodes []struct {
				Name        string
				Description string
				ForksCount  int
			}
		}
	}
}

func CallGraphQLAPI(n uint64) (string, int) {
	httpRequest := initHTTPRequest(n)

	httpClient := initHTTPClient()

	httpResponse := sendHTTPRequest(httpRequest, httpClient)

	namesCSV, forkSum := parseHTTPResponse(&httpResponse)

	return *namesCSV, *forkSum
}

func getGitlabAccessToken() (bool, string) {
	var tokenSet bool
	var accessToken string

	// os.LookupEnv() provides 3 possible states:
	// 1. token Set,
	// 2. token Set(but empty),
	// 3. token Not Set
	if _, isSet := os.LookupEnv("GITLAB_GRAPHQL_TOKEN"); isSet {
		tokenSet = true
		accessToken = os.Getenv("GITLAB_GRAPHQL_TOKEN")
	}

	return tokenSet, accessToken
}

func initHTTPClient() *http.Client {
	log.Info("Configuring HTTP client")
	httpClient := &http.Client{} // default unauthenticated client

	// return an oauth2-authenticated client instead, if token is set in env
	if tokenSet, accessToken := getGitlabAccessToken(); tokenSet {
		tokenSource := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: accessToken},
		)
		httpClient = oauth2.NewClient(context.Background(), tokenSource)
		log.Info("Oauth2-aunthenticated HTTP client used")
	}

	return httpClient
}

func initHTTPRequest(n uint64) *http.Request {
	queryjSON := repo.GetGraphQLQueryJSON(n)

	log.Info("Initializing HTTP request")
	httpRequest, err := http.NewRequest(
		"POST",                             // http method
		graphQLEndpoint,                    // gitlab graphql endpoint
		bytes.NewBuffer([]byte(queryjSON)), // json-encoded query as payload
	)
	if err != nil {
		log.Fatal(err)
	}

	// to avoid "Unexpected end of document" error, add "application/json" in content-type header
	contentType := "application/json; charset=utf-8"
	log.Info("Setting request Content-Type header as '", contentType, "'")
	httpRequest.Header.Set("Content-Type", contentType)

	return httpRequest
}

func sendHTTPRequest(httpRequest *http.Request, httpClient *http.Client) []byte {
	if httpRequest == nil {
		log.Fatal("httpRequest cannot be nil")
	}
	if httpClient == nil {
		log.Fatal("httpClient cannot be nil")
	}

	log.Info("Sending HTTP request")
	response, err := httpClient.Do(httpRequest)
	if err != nil {
		log.Fatal(err)
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	return responseBody
}

func parseHTTPResponse(httpResponse *[]byte) (*string, *int) {
	responseData := &responseSchema{} // will hold unmarshaled response
	if err := json.Unmarshal(*httpResponse, &responseData); err != nil {
		log.Fatal(err)
	}

	var (
		namesCSV string
		forkSum  int
	)

	for _, node := range responseData.Data.Projects.Nodes {
		namesCSV += node.Name + ", "
		forkSum += node.ForksCount
	}

	return &namesCSV, &forkSum
}

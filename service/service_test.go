package service

import (
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// test for missing 'application/json' which returns error "Unexpected end of document" from endpoint
func TestInitHTTPRequest(t *testing.T) {
	httpRequest := initHTTPRequest(5)
	contentTypeHeader := httpRequest.Header.Get("Content-Type")
	if !strings.Contains(contentTypeHeader, "application/json") {
		t.Fatalf(`got unset content type header, want 'application/json'`)
	}
}

// test for regressions or esoteric cases such as bit-flips turning tokenSet to false
func TestGetGitlabAccessToken(t *testing.T) {
	tokenSet, accessToken := getGitlabAccessToken()

	if !tokenSet && len(accessToken) > 0 {
		t.Fatalf(`got (%v, %v) returned, want (false,'') or (true,%v)`, tokenSet, accessToken, accessToken)
	}
}

// test if checks ensuring n>0 are still in place
func TestInitHTTPRequestZero(t *testing.T) {
	initHTTPRequest(0) // this will trigger a panic if the zero-check is absent

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered. Error:\n", r)
			t.Fatalf(`got a panic, want a zero-check to ensure n>0`)
		}
	}()
}

// test if the return type is "*http.Client". If not, (if nil), might lead to a panic downsteam
func TestInitHTTPClient(t *testing.T) {
	client := initHTTPClient()
	clientType := reflect.TypeOf(client).String()
	if clientType != "*http.Client" {
		t.Fatalf(`got %v, want '*http.Client'`, clientType)
	}
}

func TestParseHTTPResponse(t *testing.T) {
	httpResponse := sendHTTPRequest(initHTTPRequest(5), initHTTPClient())

	namesCSV, forkSum := parseHTTPResponse(&httpResponse)

	// namesCSV should contains at least 5 repo names by default
	if len(*namesCSV) == 0 {
		t.Fatalf(`got (%v, %d) returned, want non-zero length namesCSV`, namesCSV, forkSum)
	}
}

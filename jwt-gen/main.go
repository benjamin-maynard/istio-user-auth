package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"

	"google.golang.org/api/idtoken"
	"google.golang.org/api/option"
)

const (
	serviceAccountKeyPath = "service-account.json" // Replace with path relative to this folder of your Service Account Key
	serviceURL            = "http://34.29.242.112/"
)

func main() {

	// Create a token source for the specified audience URL
	// Audience must be provided, but does not have to be validated at the Istio level
	ts, err := idtoken.NewTokenSource(context.Background(), "https://api.example-audience.com", option.WithCredentialsFile(serviceAccountKeyPath))
	if err != nil {
		log.Fatalf("failed to create NewTokenSource: %s\n", err)
	}

	// Get the ID token. This should be called every time before a request, it is cached and then
	// renewed if it has expired.
	token, err := ts.Token()
	if err != nil {
		log.Fatalf("failed to receive token: %s\n", err)
	}

	req, err := http.NewRequest("GET", serviceURL, nil)
	if err != nil {
		log.Fatalf("failed to build HTTP request: %s\n", err)
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token.AccessToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("error making request: %s\n", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("error reading response: %s\n", err)
		return
	}

	fmt.Println(string(body))

}

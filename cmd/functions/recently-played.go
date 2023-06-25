package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	"github.com/harryalaw/recently-played/pkg/spotify"
)

var client = &http.Client{}

func getAccessToken(id string) (string, error) {
	// get the refresh token from db.
	// get the access token from spotify;
	token, error := spotify.RefreshAccessToken(client, os.Getenv("SPOTIFY_REFRESH_TOKEN"))

	return token.AccessToken, error

}

func getRecentlyPlayed(accessToken string) (interface{}, error) {
	req, err := http.NewRequest("GET", "https://api.spotify.com/v1/me/player/recently-played", nil)
	if err != nil {
		log.Print("Failed to build request: ", err.Error())
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+string(accessToken))
	resp, err := client.Do(req)

	if err != nil {
		log.Print("Request for recently-played failed: ", err.Error())
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Print("Non OK response: ", resp.StatusCode)
		return "", nil
	}

	var a interface{}
	json.NewDecoder(resp.Body).Decode(&a)

	log.Printf("RESPONSE: %+v", a)

	return a, nil
}

func callLambda(id string) (string, error) {
	tokens, err := getAccessToken(id)
	if err != nil {
		return "", err
	}
	json, err := getRecentlyPlayed(tokens)

	return fmt.Sprintf("%+v", json), nil
}

func handleRequest(ctx context.Context) (string, error) {
	// request context
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf("REQUEST ID: %s", lc.AwsRequestID)
	// global variable
	log.Printf("FUNCTION NAME: %s", lambdacontext.FunctionName)
	// AWS SDK call
	usage, err := callLambda("id")
	if err != nil {
		return "ERROR", err
	}
	return usage, nil
}

func main() {
	runtime.Start(handleRequest)
}

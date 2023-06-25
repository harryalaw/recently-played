package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	models "github.com/harryalaw/recently-played/pkg/models/spotify"
	"github.com/harryalaw/recently-played/pkg/repository"
	"github.com/harryalaw/recently-played/pkg/spotify"
)

var client = &http.Client{}

func getAccessToken(id int) (string, error) {
	// get the refresh token from db.
	// get the access token from spotify;
	token, error := spotify.RefreshAccessToken(client, os.Getenv("SPOTIFY_REFRESH_TOKEN"))

	return token.AccessToken, error

}

func getRecentlyPlayed(accessToken string) (*models.RecentlyPlayedResponse, error) {
	res, err := spotify.GetRecentlyPlayed(client, accessToken)

	if err != nil {
		return nil, err
	}
	log.Printf("Recently Played tracks: %+v", res)
	return res, nil
}

func persistData(data *models.RecentlyPlayedResponse, userId int) error {
	return repository.PersistTracks(data, userId)
}

func callLambda(id int) (string, error) {
	tokens, err := getAccessToken(id)
	if err != nil {
		return "", err
	}
	tracks, err := getRecentlyPlayed(tokens)
	if err != nil {
		return "", err
	}
	err = persistData(tracks, id)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%+v", tracks), nil
}

func handleRequest(ctx context.Context) (string, error) {
	// request context
	lc, _ := lambdacontext.FromContext(ctx)
	log.Printf("REQUEST ID: %s", lc.AwsRequestID)
	// global variable
	log.Printf("FUNCTION NAME: %s", lambdacontext.FunctionName)
	// AWS SDK call
	usage, err := callLambda(2)
	if err != nil {
		log.Printf("ERROR: %s", err)
		return "ERROR", err
	}
	return usage, nil
}

func main() {
	runtime.Start(handleRequest)
}

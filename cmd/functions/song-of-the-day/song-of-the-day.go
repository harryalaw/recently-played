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

func callLambda() (string, error) {
	accessToken, err := getAccessToken()
	if err != nil {
		return "", err
	}

	offset, err := getOffset()
	if err != nil {
		return "", err
	}

	tracks, err := getPlaylist(accessToken, os.Getenv("PLAYLIST_ID"), offset)
	if err != nil {
		return "", err
	}

	err = persistData(tracks, offset)
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
    log.Printf("Environment: %s", os.Getenv("TEST_ENV"));
	usage, err := callLambda()
	if err != nil {
		log.Printf("ERROR: %s", err)
		return "ERROR", err
	}
	return usage, nil
}

func getAccessToken() (string, error) {
	token, error := spotify.RefreshAccessToken(client, os.Getenv("SPOTIFY_REFRESH_TOKEN"))

	return token.AccessToken, error
}

func main() {
	runtime.Start(handleRequest)
}

func getOffset() (int, error) {
	return repository.GetNextIndex()
}

func getPlaylist(accessToken, playlistId string, offset int) ([]models.PlaylistObject, error) {
	data := make([]models.PlaylistObject, 0)
	response, err := spotify.GetPlaylist(client, accessToken, playlistId, offset)
	if err != nil {
		return data, err
	}
	log.Printf("Playlist call succeded: %d items retrieved", len(response.Items))

	for _, item := range response.Items {
		data = append(data, item)
	}

	for response.Next != "" {
		response, err = spotify.GetNextPlaylist(client, accessToken, response.Next)

		if err != nil {
			return data, err
		}
		log.Printf("Playlist call for next items succeded: %d items retrieved", len(response.Items))

		for _, item := range response.Items {
			data = append(data, item)
		}
	}
	return data, nil
}

func persistData(data []models.PlaylistObject, offset int) error {
    if len(data) == 0 {
        log.Printf("No new songs to track\n");
        return nil
    }
	return repository.PersistPlaylistTracks(data, offset)
}

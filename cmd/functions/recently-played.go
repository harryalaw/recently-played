package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	runtime "github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
	_ "github.com/go-sql-driver/mysql"
	"github.com/harryalaw/recently-played/pkg/models/entities"
	models "github.com/harryalaw/recently-played/pkg/models/spotify"
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
	db, err := sql.Open("mysql", os.Getenv("DSN"))

	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}

	defer db.Close()

	entities := entities.RecentlyPlayedList(data, userId)

	log.Printf("Inserting %d tracks", len(entities))

	// todo: Update db schema to have href instead of uri
	queryString := `INSERT into recently_played_tracks
    (user_id, played_at, uri)
    VALUES `

	numOfFields := 3

	params := make([]interface{}, len(entities)*numOfFields)
	for i, e := range entities {
		pos := i * numOfFields
		params[pos+0] = e.UserId
		params[pos+1] = e.PlayedAt
		params[pos+2] = e.Href

		queryString += `(?, ?, ?),`
	}

	queryString = queryString[:len(queryString)-1] // drop last comma

	_, err = db.Exec(queryString, params...)
	return err
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

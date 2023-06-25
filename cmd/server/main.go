package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/harryalaw/recently-played/pkg/models/entities"
	"github.com/harryalaw/recently-played/pkg/spotify"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic(err)
	}
	client := &http.Client{}
	token, err := spotify.RefreshAccessToken(client, os.Getenv("SPOTIFY_REFRESH_TOKEN"))

	if err != nil {
		panic(err)
	}
    data, err := spotify.GetRecentlyPlayed(client, token.AccessToken)


    entities := entities.RecentlyPlayedList(data, 2)
    fmt.Printf("%+v", entities);
}

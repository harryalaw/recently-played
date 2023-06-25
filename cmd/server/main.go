package main

import (
	"net/http"
	"os"

	"github.com/harryalaw/recently-played/pkg/repository"
	"github.com/harryalaw/recently-played/pkg/spotify"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env.test"); err != nil {
		panic(err)
	}

	client := &http.Client{}
	token, err := spotify.RefreshAccessToken(client, os.Getenv("SPOTIFY_REFRESH_TOKEN"))

	if err != nil {
		panic(err)
	}
	data, err := spotify.GetRecentlyPlayed(client, token.AccessToken)

	err = repository.PersistTracks(data, 2)
}

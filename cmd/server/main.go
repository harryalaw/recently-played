package main

import (
	"net/http"
	"os"

	models "github.com/harryalaw/recently-played/pkg/models/spotify"
	"github.com/harryalaw/recently-played/pkg/repository"
	"github.com/harryalaw/recently-played/pkg/spotify"
	"github.com/joho/godotenv"
)

func main() {
	offset := 0
	if err := godotenv.Load(".env.test"); err != nil {
		panic(err)
	}

	client := &http.Client{}
	token, err := spotify.RefreshAccessToken(client, os.Getenv("SPOTIFY_REFRESH_TOKEN"))

	if err != nil {
		panic(err)
	}

	data := make([]models.PlaylistObject, 0)
	response, err := spotify.GetPlaylist(client, token.AccessToken, os.Getenv("PLAYLIST_ID"), offset)
	if err != nil {
		panic(err)
	}

	for _, item := range response.Items {
		data = append(data, item)
	}

	for response.Next != "" {
		response, err = spotify.GetNextPlaylist(client, token.AccessToken, response.Next)

		if err != nil {
			panic(err)
		}

		for _, item := range response.Items {
			data = append(data, item)
		}
	}

	repository.PersistPlaylistTracks(data, offset)
}

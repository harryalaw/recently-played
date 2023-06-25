package main

import (
	"fmt"
	"net/http"
	"os"

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

	fmt.Printf("%+v", token)
}

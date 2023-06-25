package spotify

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"

	models "github.com/harryalaw/recently-played/pkg/models/spotify"
)

const TOKEN_ENDPOINT = "https://accounts.spotify.com/api/token"

func RefreshAccessToken(client *http.Client, refreshToken string) (*models.RefreshTokenResponse, error) {
	values := url.Values{
		"grant_type":    []string{"refresh_token"},
		"refresh_token": []string{refreshToken},
	}

	req, err := http.NewRequest("POST", TOKEN_ENDPOINT, strings.NewReader(values.Encode()))
	req.SetBasicAuth(os.Getenv("SPOTIFY_CLIENT_ID"), os.Getenv("SPOTIFY_CLIENT_SECRET"))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)

	if err != nil {
		log.Print("Request to get access token failed: ", err.Error())
		return nil, err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		log.Print("Non OK response: ", resp.StatusCode)
		log.Print(resp.Body)
		return nil, err
	}

	var a models.RefreshTokenResponse
	json.NewDecoder(resp.Body).Decode(&a)

	return &a, nil
}

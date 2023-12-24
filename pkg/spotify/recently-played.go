package spotify

import (
	"encoding/json"
	"log"
	"net/http"

	models "github.com/harryalaw/recently-played/pkg/models/spotify"
)

const RECENTLY_PLAYED_URL = SPOTIFY_API_URL + "/me/player/recently-played?limit=50"

func GetRecentlyPlayed(client *http.Client, accessToken string) (*models.RecentlyPlayedResponse, error) {

	req, err := http.NewRequest("GET", RECENTLY_PLAYED_URL, nil)
	if err != nil {
		log.Print("Failed to build request: ", err.Error())
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+string(accessToken))
	resp, err := client.Do(req)

	if err != nil {
		log.Print("Request for recently-played failed: ", err.Error())
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Print("Non OK response: ", resp.StatusCode)
		return nil, nil
	}

	var a models.RecentlyPlayedResponse
	json.NewDecoder(resp.Body).Decode(&a)

	return &a, nil
}

package spotify

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	models "github.com/harryalaw/recently-played/pkg/models/spotify"
)

const PLAYLIST_FIELDS = "fields=items(track(name,href,album(name),artists(name))),next,offset"
const PLAYLIST_URL = SPOTIFY_API_URL + "/playlists/%s/tracks" + "?offset=%d&" + PLAYLIST_FIELDS

func GetPlaylist(client *http.Client, accessToken, playlistId string, offset int) (*models.PlaylistResponse, error) {
	req, err := http.NewRequest("GET", playlistUrl(playlistId, offset), nil)
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
		log.Print(req)
		log.Print("Non OK response: ", resp.StatusCode)
		log.Print(resp)
		return nil, fmt.Errorf("non ok response")
	}

	var a models.PlaylistResponse
	json.NewDecoder(resp.Body).Decode(&a)

	return &a, nil

}

func GetNextPlaylist(client *http.Client, accessToken, nextUrl string) (*models.PlaylistResponse, error) {
	req, err := http.NewRequest("GET", nextUrl, nil)
	if err != nil {
		log.Print("Failed to build request: ", err.Error())
		return nil, err
	}

	req.Header.Add("Authorization", "Bearer "+string(accessToken))
	log.Print(req)
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

	var a models.PlaylistResponse
	json.NewDecoder(resp.Body).Decode(&a)

	return &a, nil
}

func playlistUrl(playlistId string, offset int) string {
	return fmt.Sprintf(PLAYLIST_URL, playlistId, offset)
}

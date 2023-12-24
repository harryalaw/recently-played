package entities

import (
	"encoding/json"
	"fmt"
	"time"

	models "github.com/harryalaw/recently-played/pkg/models/spotify"
)

type SongOfTheDay struct {
	Id         int
	Day        string
	TrackTitle string
	AlbumTitle string
	Artists    string
	TrackUri   string
}

func SongOfTheDayList(response []models.PlaylistObject, offset int) []SongOfTheDay {
	totalItems := len(response)

	entities := make([]SongOfTheDay, totalItems)

	for i := 0; i < totalItems; i++ {
		date := indexToISODate(i, offset)
		item := response[i]
		entities[i] = SongOfTheDay{
			Day:        date,
			TrackTitle: item.Track.Name,
			AlbumTitle: item.Track.Album.Name,
			Artists:    encodeArtists(item.Track.Artists),
			TrackUri:   item.Track.Href,
		}
	}

	return entities
}

func indexToISODate(index int, offset int) string {
	// Starting date
	startDate := time.Date(2019, time.January, 1, 0, 0, 0, 0, time.UTC)

	// Calculate the date based on the index
	resultDate := startDate.AddDate(0, 0, index+offset)

	// Format the date as "2006-01-02"
	return resultDate.Format("2006-01-02")
}

func encodeArtists(artists []models.PlaylistArtist) string {
	artistNames := make([]string, len(artists))
	for i, artist := range artists {
		artistNames[i] = artist.Name
	}
	b, err := json.Marshal(artistNames)
	if err != nil {
		panic(err)
	}

	out := string(b)
	fmt.Println(out)
	return out
}

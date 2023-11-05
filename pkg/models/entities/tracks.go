package entities

import (
	models "github.com/harryalaw/recently-played/pkg/models/spotify"
)

type RecentlyPlayed struct {
	Id       int
	PlayedAt string
	Href     string
	UserId   int
	TrackId  string
	AlbumId  string
}

func RecentlyPlayedList(response *models.RecentlyPlayedResponse, user int) []RecentlyPlayed {
	totalItems := len(response.Items)

	entities := make([]RecentlyPlayed, totalItems)

	for i := 0; i < totalItems; i++ {
		entities[i] = RecentlyPlayed{
			PlayedAt: response.Items[i].PlayedAt,
			Href:     response.Items[i].Track.Href,
			UserId:   user,
			TrackId:  response.Items[i].Track.Id,
			AlbumId:  response.Items[i].Track.Album.Id,
		}
	}

	return entities
}

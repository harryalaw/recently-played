package entities

import (
	models "github.com/harryalaw/recently-played/pkg/models/spotify"
)

type RecentlyPlayedEntity struct {
	PlayedAt string
	Href     string
	UserId   int
}

func RecentlyPlayedList(response *models.RecentlyPlayedResponse, user int) []RecentlyPlayedEntity {
	totalItems := len(response.Items)

	entities := make([]RecentlyPlayedEntity, totalItems)

	for i := 0; i < totalItems; i++ {
		entities[i] = RecentlyPlayedEntity{
			PlayedAt: response.Items[i].PlayedAt,
			Href:     response.Items[i].Track.Href,
			UserId:   user,
		}
	}

	return entities
}

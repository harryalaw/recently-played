package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/harryalaw/recently-played/pkg/models/entities"
	models "github.com/harryalaw/recently-played/pkg/models/spotify"
)

func PersistPlaylistTracks(data []models.PlaylistObject, offset int) error {
	db, err := sql.Open("mysql", os.Getenv("DSN"))

	if err != nil {
		log.Fatalf("Failed to connect: %+v", err)
	}

	defer db.Close()

	entities := entities.SongOfTheDayList(data, offset)
	log.Printf("Retrieved %d tracks", len(entities))

	for _, entity := range entities {
		fmt.Println(entity)
	}

	queryString := `INSERT into song_of_the_day
    (track_title, album_title, artists, track_uri, day)
    VALUES `

	numOfFields := 5

	params := make([]interface{}, len(entities)*numOfFields)
	for i, e := range entities {
		pos := i * numOfFields
		params[pos+0] = e.TrackTitle
		params[pos+1] = e.AlbumTitle
		params[pos+2] = e.Artists
		params[pos+3] = e.TrackUri
		params[pos+4] = e.Day

		queryString += `(?, ?, ?, ?, ?),`
	}

	queryString = queryString[:len(queryString)-1] // drop last comma

	_, err = db.Exec(queryString, params...)
	return nil
}

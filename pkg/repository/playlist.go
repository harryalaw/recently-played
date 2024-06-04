package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/harryalaw/recently-played/pkg/models/entities"
	models "github.com/harryalaw/recently-played/pkg/models/spotify"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func GetNextIndex() (int, error) {
	db, err := sql.Open("libsql", os.Getenv("DSN"))
	if err != nil {
		log.Fatalf("Failed to connect: %+v", err)
	}

	defer db.Close()

	var offset int
	row := db.QueryRow("SELECT id FROM song_of_the_day order by ID desc LIMIT 1")
	err = row.Scan(&offset)

	if err != nil {
		if err == sql.ErrNoRows {
			return 0, fmt.Errorf("no rows found %+v", err)
		}
		return 0, fmt.Errorf("error getting next index")
	}

    return offset, nil
}

func PersistPlaylistTracks(data []models.PlaylistObject, offset int) error {
	db, err := sql.Open("libsql", os.Getenv("DSN"))

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
	if err != nil {
		log.Printf("Error persisting tracks: %+v", err)
		return err
	}
	return nil
}

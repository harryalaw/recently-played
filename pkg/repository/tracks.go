package repository

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/harryalaw/recently-played/pkg/models/entities"
	models "github.com/harryalaw/recently-played/pkg/models/spotify"
	_ "github.com/tursodatabase/libsql-client-go/libsql"
)

func PersistTracks(data *models.RecentlyPlayedResponse, userId int) error {
	db, err := sql.Open("libsql", os.Getenv("DSN"))

	if err != nil {
		log.Fatalf("Failed to connect: %+v", err)
	}

	defer db.Close()

	entities := entities.RecentlyPlayedList(data, userId)
	log.Printf("Retrieved %d tracks", len(entities))

	timestamp, err := mostRecentTimestamp(db, userId)
	if err != nil {
		log.Fatalf("Failed to fetch most recent track: %+v", err)
	}
	entities = filterEntities(entities, timestamp)
	log.Printf("Inserting %d tracks", len(entities))

	if len(entities) == 0 {
		log.Printf("No new plays")
		return nil
	}

	// todo: Update db schema to have href instead of uri
	queryString := `INSERT into recently_played_tracks
    (user_id, played_at, uri, track_id, album_id)
    VALUES `

	numOfFields := 5

	params := make([]interface{}, len(entities)*numOfFields)
	for i, e := range entities {
		pos := i * numOfFields
		params[pos+0] = e.UserId
		params[pos+1] = e.PlayedAt
		params[pos+2] = e.Href
		params[pos+3] = e.TrackId
		params[pos+4] = e.AlbumId

		queryString += `(?, ?, ?, ?, ?),`
	}

	queryString = queryString[:len(queryString)-1] // drop last comma

	_, err = db.Exec(queryString, params...)
	return err
}

func mostRecentTimestamp(db *sql.DB, userId int) (string, error) {
	var timestamp string
	row := db.QueryRow("SELECT played_at FROM recently_played_tracks WHERE user_id=? ORDER BY played_at DESC LIMIT 1", userId)
	err := row.Scan(&timestamp)
	if err != nil {
		// todo: is this an error case?
		// what if someone new starts to use it :D
		if err == sql.ErrNoRows {
			return "2000-01-01T00:00:00.000Z", nil
		}
		return timestamp, fmt.Errorf("Error fetching timestamps: %v", err)
	}

	if row.Err() != nil {
		return "", row.Err()
	}
	row.Scan(&timestamp)

	return timestamp, nil
}

type RecentlyPlayedList = []entities.RecentlyPlayed

func filterEntities(entities RecentlyPlayedList, mostRecent string) RecentlyPlayedList {
	// odd that the next line doesn't work
	//filtered := make([](entities.RecentlyPlayedEntity), 0)
	filtered := make(RecentlyPlayedList, 0)

	for _, ent := range entities {
		if strings.Compare(ent.PlayedAt, mostRecent) == 1 {
			filtered = append(filtered, ent)
		}
	}

	return filtered
}

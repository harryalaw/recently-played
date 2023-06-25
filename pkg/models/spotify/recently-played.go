package models

type RecentlyPlayedResponse struct {
	Limit int
	Total int
	Items []PlayHistoryObject
}

type PlayHistoryObject struct {
	PlayedAt string `json:"played_at"`
	Track    Track
}

type Track struct {
	Href string
	Uri  string
}

package models

type PlaylistResponse struct {
	Next   string
	Items  []PlaylistObject
	Offset int
}

type PlaylistObject struct {
	Track PlaylistTrack
}

type PlaylistTrack struct {
	Album   PlaylistAlbum
	Artists []PlaylistArtist
	Href    string
	Name    string
}

type PlaylistAlbum struct {
	Name string
}

type PlaylistArtist struct {
	Name string
}

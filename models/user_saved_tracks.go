package models

//UserSavedTracks stores all the items with relevant data of liked tracks
type UserSavedTracks struct {
	Items []Item `json:"items"`
}

//Item represents a liked track data
type Item struct {
	AddedAt string `json:"added_at"`
	Track   struct {
		Album struct {
			AlbumType string `json:"album_type"`
			Artists   []struct {
				Name string `json:"name"`
			}
		}
		Name       string `json:"name"`
		Popularity int    `json:"popularity"`
	}
}

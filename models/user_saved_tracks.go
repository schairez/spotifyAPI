package models

//UserSavedTracks represents the /GET resp from calling the /me/tracks endpoint with all the values we care about
type UserSavedTracks struct {
	Href  string `json:"href"`
	Items []Item `json:"items"`
	Next  string `json:"next"`
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
		PreviewURL string `json:"preview_url"`
		URI        string `json:"uri"`
	}
}

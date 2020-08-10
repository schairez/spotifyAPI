package models

// SpotifyUser is a retrieved and authenticated user from the spotify api
type SpotifyUser struct {
	ID         string `json:"id"`
	Name       string `json:"display_name"`
	Email      string `json:"email"`
	Country    string `json:"country"`
	UserURI    string `json:"uri"`
	ImageURL   string `json:"images"`
	PictureURL string `json:"picture"`
}

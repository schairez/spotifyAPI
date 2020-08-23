package models

import "time"

// SpotifyUser is a retrieved and authenticated user from the spotify api
type SpotifyUser struct {
	ID      string `json:"id"`
	Name    string `json:"display_name"`
	Email   string `json:"email"`
	Country string `json:"country"`
	Images  []struct {
		URL string `json:"url"`
	} `json:"images"`
	UserAPIEndpoint string `json:"href"`
	AccessToken     string `json:"access_token"`
	RefreshToken    string `json:"refresh_token"`
	TokenExpiry     time.Time
	TokenType       string
}

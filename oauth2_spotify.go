package main

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

//NewSpotifyConfig fdfd
func NewSpotifyConfig(clientID, clientSecret string) *oauth2.Config {

	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "https://localhost:8080/spotify/callback",
		Scopes: []string{"user-read-email",
			"user-top-read", "user-library-read", "playlist-read-private", "user-read-recently-played", "playlist-read-private"},
		Endpoint: spotify.Endpoint,
	}
	return cfg
}

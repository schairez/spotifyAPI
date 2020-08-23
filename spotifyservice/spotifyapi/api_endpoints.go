package spotifyapi

import (
	"net/url"
	"path"
)

//APIEndpoints contains all the urls that we'll use for Spotify API
//NOTE: /GET for UsersSavedTracksURL or UserSavedAlbumsURL  requres user-library-read scope
type APIEndpoints struct {
	BaseURL            *url.URL //BaseURL for spotify API
	UserProfileURL     *url.URL //Check User's Profile
	UserSavedTracksURL *url.URL //Check User's Saved Tracks
	UserSavedAlbumsURL *url.URL
}

//NewAPI constructs a new Spotify api requests struct with all the endpoints as url.URL structs
func NewAPI() (*APIEndpoints, error) {
	base := "https://api.spotify.com/v1/"
	baseURL, err := url.Parse(base)
	if err != nil {
		return nil, err

	}
	meURL, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	meURL.Path = path.Join(meURL.Path, "/me/")

	meTracksURL, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	meTracksURL.Path = path.Join(meTracksURL.Path, "/me/tracks")

	albumsURL, err := url.Parse(base)
	if err != nil {
		return nil, err
	}
	albumsURL.Path = path.Join(albumsURL.Path, "/me/albums")

	api := &APIEndpoints{}
	api.BaseURL = baseURL
	api.UserProfileURL = meURL
	api.UserSavedTracksURL = meTracksURL
	api.UserSavedAlbumsURL = albumsURL

	return api, nil

}

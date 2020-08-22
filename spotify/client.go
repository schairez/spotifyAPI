// Copyright 2020 Sergio Chairez. All rights reserved.
// Use of this source code is governed by a MIT style license that can be found
// in the LICENSE file.

package spotify

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/schairez/spotifywork/models"
	"golang.org/x/oauth2"
)

//Client wraps the oauth2 config so we
// can extend methods specific to spotify"s api
type Client struct {
	Config *oauth2.Config
	API    *APIEndpoints
}

//NewClient gives us a ServiceClient to wrap the ouath2 config and api endpoints for our client requests
func NewClient(clientID, clientSecret, redirectURL string) *Client {
	_api, err := NewAPI()
	if err != nil {
		log.Fatalf("Err building api endpoints %v", err)
	}
	s := &Client{
		Config: newSpotifyConfig(clientID, clientSecret, redirectURL),
		API:    _api}
	return s

}

//GetUserProfileRequest calls the https://api.spotify.com/v1/me
//endpoint and gives us useful info about the user
func (s *Client) GetUserProfileRequest(ctx context.Context, token *oauth2.Token) (*models.SpotifyUser, error) {
	var endpoint *url.URL = s.API.UserProfileURL
	url := endpoint.String()
	log.Printf("user profile url: %v\n", url)

	user := &models.SpotifyUser{}

	httpClient := s.Config.Client(ctx, token)
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode >= http.StatusBadRequest {
		log.Println("status code todo:return err")

	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(user); err != nil {
		log.Printf("Could not decode body: %v\n", err)
		return nil, err
	}
	user.AccessToken = token.AccessToken
	user.RefreshToken = token.RefreshToken
	user.TokenExpiry = token.Expiry
	user.TokenType = token.TokenType
	return user, nil
}

//QParams holds the optional qparams for the api calls
type QParams struct {
	Limit  *int    //The maximum number of objects to return. Default: 20. Minimum: 1. Maximum: 50
	Offset *int    //The index of the first object to return. Default: 0
	Market *string //An ISO 3166-1 alpha-2 country code; for track-relinking
	// queryParams map[string]string
}

//countries spotify is available in
//TODO: expand to use all the ISO country codes from this link
//https://support.spotify.com/us/article/full-list-of-territories-where-spotify-is-available/

var (
	markets = [...]string{"ad", "ar", "at", "au", "be", "bg", "bo", "br", "ca", "ch",
		"cl", "co", "cr", "cy", "cz", "de", "dk", "do", "ec", "ee", "es", "fi", "fr",
		"gb", "gr", "gt", "hk", "hn", "hu", "id", "ie", "is", "it", "jp", "li", "lt",
		"lu", "lv", "mc", "mt", "mx", "my", "ni", "nl", "no", "nz", "pa", "pe", "ph",
		"pl", "pt", "py", "se", "sg", "sk", "sv", "tr", "tw", "us", "uy"}
)

func inMarketsArr(m string) bool {
	for _, v := range markets {
		if m == v {
			return true
		}
	}
	return false
}
func validMarketOpt(m string) bool {
	n := len(m)
	validLen := n == 2
	if validLen && inMarketsArr(m) {
		return true
	}
	return false
}

// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.Token))

/*
below takes into account pagination; we want to retrieve up to the
10,000 saved tracks; (still in progress)
*/

//have to make a generic request, and a callback fn to keep recursively calling if we identify the next url in the bytes stream
//GET UserLibrary
//TODO:
//make limitKey and offsetKey a type string with an initialized value...
//is this concurrent safe? using global variables?
// limit Default: 20. Minimum: 1. Maximum: 50

//GetUserSavedTracks calls the https://api.spotify.com/v1/me/tracks
//endpoint with the params and gives us the user's saved tracks in their library
func (s *Client) GetUserSavedTracks(ctx context.Context, token *oauth2.Token, q *QParams) (*models.UserSavedTracks, error) {
	var endpoint *url.URL = s.API.UserSavedTracksURL
	log.Println("SAME ADDRESS ?????")
	log.Println(&endpoint == &s.API.UserSavedTracksURL)
	if q != nil {
		params := url.Values{}
		if q.Limit != nil {
			l := *(q).Limit
			valid := (l >= 1) && (l <= 50)
			if valid {
				params.Set("limit", strconv.Itoa(l))
			}
		}
		if q.Offset != nil {
			offset := *(q).Offset
			if offset > 0 {
				params.Set("offset", strconv.Itoa(offset))
			}
		}
		if q.Market != nil {
			m := *(q).Market
			if validMarketOpt(m) {
				params.Set("market", *(q).Market)
			}
		}

		endpoint.RawQuery = params.Encode()
	}
	url := endpoint.String()
	log.Println(url)
	log.Printf("User saved tracks url: %v\n", url)

	tracks := &models.UserSavedTracks{}

	httpClient := s.Config.Client(ctx, token)
	resp, err := httpClient.Get(url)
	if resp.StatusCode >= http.StatusBadRequest {
		log.Println("status code todo:return err")

	}
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if err := json.NewDecoder(resp.Body).Decode(tracks); err != nil {
		log.Printf("Could not decode body: %v\n", err)
		return nil, err
	}

	return tracks, nil

}

// func tracksQuery(limit string) {

// }

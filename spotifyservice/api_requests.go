package spotifyservice

import (
	"log"
	"net/http"
	"strconv"

	"golang.org/x/oauth2"
)

//BaseAPIURL is the base url for api requests
const baseAPIURL = "https://api.spotify.com/v1/"

//SpotifyServiceClient wraps the oauth2 config so we
// can extend methods specific to spotify's api
type spotifyServiceClient struct {
	Config     *oauth2.Config
	BaseAPIURL string
}

//NewSpotifyServiceClient gives us a wraapped extended methods..
func newSpotifyServiceClient(clientID, clientSecret, redirectURL string) *spotifyServiceClient {
	s := &spotifyServiceClient{
		Config:     newSpotifyConfig(clientID, clientSecret, redirectURL),
		BaseAPIURL: baseAPIURL}
	return s

}

// func spotifyAPIRequest(token *oauth2.Token, apiEndpoint string) {
// 	client := s.spotifyCfg.Client(context.Background(), token)

// }

//the spotify web api reference is a good resource
//https://developer.spotify.com/documentation/web-api/reference-beta/

/*
supported endpoints
GET https://api.spotify.com/v1/me/tracks
Requirements:
scope: user-library-read
query params:
limit: default 20,
*/

//APIRequest calls and returns the api request to spotify
//below is a generic api request
// func (s *SpotifyServiceClient) APIRequest(
// 	ctx context.Context,
// 	token *oauth2.Token,
// 	method string,
// 	apiEndpoint string) {
// 	client := s.Config.Client(ctx, token)
// 	endpoint := fmt.Sprintf("%s%s", baseAPIURL, apiEndpoint)
// 	req, err := http.NewRequest(method, endpoint, nil)

// }

func getLikedTracksRequest(limit, offset int) (*http.Request, error) {
	endpoint := "https://api.spotify.com/v1/me/tracks"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Println("error creating new request")
		return nil, err
	}
	q := req.URL.Query()
	q.Add("limit", strconv.Itoa(limit))
	q.Add("offset", strconv.Itoa(offset))
	return req, nil

}

func getUserRequest() (*http.Request, error) {
	endpoint := "https://api.spotify.com/v1/me"
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		log.Println("error creating new request")
		return nil, err
	}
	return req, nil

}

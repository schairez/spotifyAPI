package spotifyservice

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/schairez/spotifywork/models"
	"github.com/schairez/spotifywork/spotifyapi"

	"golang.org/x/oauth2"
)

//BaseAPIURL is the base url for api requests
const baseAPIURL = "https://api.spotify.com/v1/"

type params struct{}

//SpotifyServiceClient wraps the oauth2 config so we
// can extend methods specific to spotify's api
type spotifyServiceClient struct {
	config *oauth2.Config
	api    *spotifyapi.APIEndpoints
}

//NewSpotifyServiceClient gives us a wraapped extended methods..
func newSpotifyServiceClient(clientID, clientSecret, redirectURL string) *spotifyServiceClient {
	_api, err := spotifyapi.New()
	if err != nil {
		log.Fatalf("Err building api endpoints %v", err)
	}
	s := &spotifyServiceClient{
		config: newSpotifyConfig(clientID, clientSecret, redirectURL),
		api:    _api}
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
// func (s *spotifyServiceClient) APIRequest(
// 	ctx context.Context, token *oauth2.Token,
// 	method string, apiEndpoint string,
// 	params map[string]interface{}, v interface{}) {
// 	client := s.config.Client(ctx, token)
// 	endpoint := fmt.Sprintf("%s%s", s.baseAPIURL, apiEndpoint)
// 	req, err := http.NewRequest(method, endpoint, nil)

// 	resp, err := client.Do(req)
// 	if err != nil {
// 		return nil, err
// 	}

// }

func (s *spotifyServiceClient) getUserRequest(ctx context.Context, token *oauth2.Token) (*models.SpotifyUser, error) {
	user := &models.SpotifyUser{}

	endpoint := s.baseAPIURL + "me"
	// endpoint := "https://api.spotify.com/v1/me"
	client := s.config.Client(ctx, token)
	resp, err := client.Get(endpoint)
	if err != nil {
		return nil, err
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

/*





 */

/*
algo
-
func request()

#callback below
func urlParseWithParams(endpoint string, )



*/

// func (s *spotifyServiceClient) getLikedTracks() {
// 	endpoint :=
// 	u,err := url.Parse()

// }
//below good use of recursion.
//have an if check for the next json elem; which will point you to a url. if the url is valid then call the fn again?

func getLikedTracksRequest(limit, offset int) (*http.Request, error) {
	// albumsURL := fmt.Sprintf("%s/v1/catalog/us/search?types=artists&limit=1&term=%s", provider.URL, url.QueryEscape(term))
	endpoint := "https://api.spotify.com/v1/me/tracks"
	req, err := http.NewRequest(http.MethodGet, endpoint, nil)
	// req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", provider.Token))

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

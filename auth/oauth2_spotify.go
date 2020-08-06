// package spotifywork

package auth

import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/spotify"
)

/*

The Authorization Code Flow has the following steps:

Configure your application to get the Client ID and Client Secret.
Your application directs the browser to LinkedIn's OAuth 2.0 authorization page where the member authenticates. After authentication, LinkedIn's authorization server passes an authorization code to your application.
Your application sends this code to LinkedIn and LinkedIn returns an access token.
Your application uses this token to call APIs on behalf of the member.

*/

/*
No. Instead, the backend should set the JWT as a cookie in the user browser.
Make sure you flag it as Secure and httpOnly cookie. And SameSite cookie. Boy, that's a multi-flavor cookie.

/*
need to make it  rfc6749 compliant (3-legged approach)
https://tools.ietf.org/html/rfc6749#section-4.1

*/

//good explanation
//https://godoc.org/golang.org/x/oauth2#example-Config--CustomHTTP

//spotify oauth2 documentation
//https://developer.spotify.com/documentation/general/guides/authorization-guide/

/*NewSpotifyConfig struct using a 3-legged oauth2 flow; here
we'll return the configured clientid and clientsecret that we'll s use later when connecting
to spotify's auth server

TODO: modify to pass in config for scopes
make scope groups? global const vars to a specific use case?

*/

//NewSpotifyConfig sets up new
func NewSpotifyConfig(clientID, clientSecret, redirectURL string) *oauth2.Config {

	cfg := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"user-read-email", "user-top-read", "user-library-read", "playlist-read-private", "user-read-recently-played", "playlist-read-private"},
		Endpoint:     spotify.Endpoint,
	}
	return cfg
}

//we'll want to maintain local state b/w our user-agent's request
//and our redirection URI callback fn and prevent CSRF
//https://tools.ietf.org/html/rfc6749#section-10.10

// State can be some random generated hash string.
// See relevant RFC: http://tools.ietf.org/html/rfc6749#section-10.12
//generates random base64 hash string
// we'll check with our resource server since it'll send us back this cookie

//RFC 6749 compliant; we'll check our local state late and protects against
//CSRF when we call spotify's /authorize endpoint
//we'll supply the state parameter to the GET request to spotify's /authorize endpoint
//here we encode the state of a random string to the user's cookie in a state variable

//https://tools.ietf.org/html/rfc6749#section-4.1

//uses gorilla session storage?
//https://github.com/dghubble/gologin/blob/master/examples/github/main.go

//TODO: 1. generate state token with expiry time; likely in a cookie store?
//RFC 6749
//https://tools.ietf.org/html/rfc6749#section-4.1

//from
//https://github.com/Jared-Mullin/LoPhi-Music/blob/master/main.go

//https://github.com/markbates/goth/blob/master/gothic/gothic.go
// SetState sets the state string associated with the given request.
// If no state string is associated with the request, one will be generated.
// This state is sent to the provider and can be retrieved during the
// callback.
//SetState sets the state yo

// If a state query param is not passed in, generate a random
// base64-encoded nonce so that the state on the auth URL
// is unguessable, preventing CSRF attacks, as described in
//
// https://auth0.com/docs/protocols/oauth2/oauth-state#keep-reading

//TODO: 2.

// Get opens a browser window to authCodeURL for the user to
// authorize the application, and it returns the resulting
// OAuth2 code. It rejects requests where the "state" param
// does not match expectedStateVal.
//https://github.com/mholt/timeliner/blob/master/oauth2client/browser.go

//https://medium.com/@pliutau/getting-started-with-oauth2-in-go-2c9fae55d187
/*

#for
http.HandleFunc("/login", handleGoogleLogin)


func HandleSpotifyLogin(w http.ResponseWriter, r *http.Request) {
	url := googleOauthConfig.AuthCodeURL(oauthStateString)
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}
*/

//if you don't want to use a yaml config, nor use env vars. You're welcome to
//use the exported setter here

/*
type SpotifyAuthenticator struct {
	config *oauth2.config
}

*/

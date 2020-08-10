package spotifyservice

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/schairez/spotifywork/models"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/schairez/spotifywork/env"
	"golang.org/x/oauth2"
)

/*

 rt http.RoundTripper,

 https://chromium.googlesource.com/external/github.com/golang/oauth2/+/8f816d62a2652f705144857bbbcc26f2c166af9e/oauth2.go
*/

const stateCookieName = "oauthState"

//returns a base64 encoded random 32 byte string and sets a cookie on user's browser with this field
func genRandStateOauthCookie(w http.ResponseWriter) string {
	log.Println("generating cookie")
	b := make([]byte, 32)
	// rand.Read(b)
	if _, err := rand.Read(b); err != nil {
		log.Fatalf("failed to read rand fn")
	}

	state := base64.StdEncoding.EncodeToString(b)
	//TODO: suitable expiration time
	// expire in 1 month;
	// expiration := time.Now().Add(30 * 24 * time.Hour)
	//or expire in 24 hours
	//expiration := time.Now().Add(24 * time.Hour)
	expiration := time.Now().Add(time.Hour)

	//httpOnly security flag to secure our cookie from XSS; no js scripting
	cookie := http.Cookie{
		Name:     stateCookieName,
		Value:    state,
		Expires:  expiration,
		HttpOnly: true}
	//setting the set-Cookie header in the writer
	//NOTE: headers need to be set before anything else set to the writer
	http.SetCookie(w, &cookie)
	return state
}

//Server is the component of our app
type Server struct {
	cfg        *env.TomlConfig
	spotifyCfg *oauth2.Config
	router     *chi.Mux
	httpServer *http.Server
}

//NewServer returns a configured new spotify client server
func NewServer(fileName string) *Server {
	s := &Server{}
	s.initCfg(fileName)
	s.initSpotifyCfg()
	s.routes()
	return s
}

//TODO: make a filePathErr for initCfg

func (s *Server) initCfg(fileName string) {
	cfg, err := env.LoadTOMLFile(fileName)
	if err != nil {
		panic(err)
	}
	s.cfg = cfg

}

func (s *Server) initSpotifyCfg() {
	spotify, ok := s.cfg.Oauth2Providers["spotify"]
	if !ok {
		// TODO: Properly handle error
		panic("Spotify env properties not found in config")
	}
	s.spotifyCfg = newSpotifyConfig(
		spotify.ClientID,
		spotify.ClientSecret,
		spotify.RedirectURL)
}

//routes inits the route multiplexer with the assigned routes
func (s *Server) routes() {
	s.router = chi.NewRouter()
	s.router.Use(middleware.RequestID)
	s.router.Use(middleware.RealIP)
	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(middleware.StripSlashes)
	s.router.Use(middleware.Timeout(60 * time.Second))
	s.router.Use(cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler)
	s.router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	s.router.Get("/auth", func(w http.ResponseWriter, r *http.Request) {

		//ctx := r.Context()
		//check if the request contains a cookie?
		//COOKIE would be attached if the use has hit our domain
		//this would indicate that a user-agent has hit this endpoint, but not
		//that the user has authorized our app per-say
		log.Println("checking if user already has a cookie stored in their browser")
		cookie, err := r.Cookie(stateCookieName)
		if err != nil {
			log.Println("we got no cookie in request")
			log.Println(err)
		}

		fmt.Println(cookie)

		localState := genRandStateOauthCookie(w)
		fmt.Println(localState)
		fmt.Println(w.Header())
		authURL := s.spotifyCfg.AuthCodeURL(localState)
		//app directs user-agent to spotify's oauth2 auth  consent page
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)

	})
	//below we have our redirect callback as a result of a user-agent accessing
	//our /auth endpoint route
	s.router.Get("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		//check if user denied our auth request the request we receive
		//would contain a non-empty error query param in this case
		if r.FormValue("error") != "" {
			log.Printf("user authorization failed. Reason=%s", r.FormValue("error"))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		//check the state parameter we supplied to Spotify's Account's Service earlier
		//if user approved the auth, we'll have both a code and a state query param
		oauthStateCookie, err := r.Cookie(stateCookieName)
		if err != nil {
			if err == http.ErrNoCookie {
				log.Println("Error finding cookie: ", err.Error())
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("%s=%s\r\n", oauthStateCookie.Name, oauthStateCookie.Value)
		if r.FormValue("state") != oauthStateCookie.Value {
			log.Println("invalid oauth2 spotify state. state_mismatch err")
			http.Redirect(w, r, "/", http.StatusUnauthorized)
			// http.Error(w, "Invalid State token", http.StatusBadRequest)
			return
		}
		//TODO: pkce opts?
		authCode := r.FormValue("code")
		log.Printf("code=%s", authCode)
		//TODO: diff b/w background and oauth2.NoContext
		ctx := context.Background()
		//exchange auth code with an access token
		token, err := s.spotifyCfg.Exchange(ctx, authCode)
		if err != nil {
			log.Printf("error converting auth code into token; %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			// http.Error(w, err.Error(), http.StatusInternalServerError)
			//TODO:
			// or StatusForbidden?
			return
		}
		//we'll use the token to access user's protected resources
		// by calling the Spotify Web API

		log.Println(token)
		log.Println("query params?")
		queryParams := r.URL.Query()
		log.Println(queryParams)
		if reqHeadersBytes, err := json.Marshal(r.Header); err != nil {
			log.Println("Could not Marshal Req Headers")
		} else {
			log.Println(string(reqHeadersBytes))
		}

		//now we can use this token to call Spotify APIs on behalf of the user
		//use the token to get an authenticated client
		//the underlying transport obtained using ctx?
		client := s.spotifyCfg.Client(context.Background(), token)
		resp, err := client.Get("https://api.spotify.com/v1/me")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer resp.Body.Close()
		data, _ := ioutil.ReadAll(resp.Body)
		log.Println("Data calling user API: ", string(data))

		req, err := getLikedTracksRequest(10, 0)
		if err != nil {
			log.Println("error with reqeuest")
		}
		resp, err = client.Do(req)
		if err != nil {
			log.Println(err)
		}
		// data, _ = ioutil.ReadAll(resp.Body)
		// defer resp.Body.Close()
		// if resp.StatusCode != http.StatusOK {
		// 	log.Printf("http status code %d", resp.StatusCode)
		// }
		// log.Println("Data calling user liked tracks API: ", string(data))
		defer resp.Body.Close()
		tracks := &models.UserSavedTracks{}
		err = json.NewDecoder(resp.Body).Decode(tracks)
		if err != nil {
			log.Println(err)
		}
		log.Println("getting user trakcs")
		log.Println(tracks)

	})

	s.router.Get("/logout/{provider}", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

}

//401 err when no token provided

/*
{
  "error": {
    "status": 401,
    "message": "No token provided"
  }
}

*/

/*

doc: https://developer.spotify.com/documentation/web-api/reference/library/get-users-saved-tracks/
Endpoint:
GET /v1/me/tracks
NOTE:
- we can receive up to 10,000  of user's liked tracks (limit user can save)
TODO:
limit max 50, min 1, default 20
0ffset 0
we care about the track.album.artists.name

t
*/

/*
A struct or object will be Handler if it has one method ServeHTTP which takes ResponseWriter and pointer to Request.
*/

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

//Start starts the server
func (s *Server) Start() {
	s.httpServer = &http.Server{
		Addr:         ":" + s.cfg.Server.Port,
		Handler:      s.router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	log.Printf("server listening on %s\n", s.cfg.Server.Port)
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatalf("ListenAndServe err: %s", err)

	} else {
		log.Println("Server closed!")
	}

}

//Shutdown the server
func (s *Server) Shutdown() {

}

/*
ex:
req, err := http.NewRequest("GET", makeUrl("/search"), nil)

func makeUrl(path string) string {
	return "https://api.spotify.com/v1" + path
}

func SpotifyAPIRequest() {

}


*/

/*

func writeJSONResponse(w http.ResponseWriter, status int, data []byte) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	w.Header().Set("Connection", "close")
	w.WriteHeader(status)
	w.Write(data)
}
*/

/*
func handleOauthSpotifyLogin(spotifyCfg *oauth2.Config) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Take the context out from the request
		ctx := r.Context()
		localState := genRandStateOauthCookie(w)
		fmt.Println(localState)
		fmt.Println(w.Header())
		authURL := spotifyCfg.AuthCodeURL(localState)
		//app directs user-agent to spotify's oauth2 auth page
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)

		// Call your original http.Handler
		// h.ServeHTTP(w, r)

	})

}
*/
// func oauthSpotifyCallback()

/*

func newAPIRequest() (string, error) {
	var response *http.Response

	req, err := http.NewRequest("POST", oc.oauthUrl, strings.NewReader(postBody))
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if err != nil {
		return "", err
	}

	return "", nil

}

*/

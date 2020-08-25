// Copyright 2020 Sergio Chairez. All rights reserved.
// Use of this source code is governed by a MIT style license that can be found
// in the LICENSE file.

package spotifyservice

import (
	"crypto/rand"
	"encoding/base64"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/schairez/spotifywork/env"
	"github.com/schairez/spotifywork/spotifyservice/spotifyapi"
)

/*

 rt http.RoundTripper,

 https://chromium.googlesource.com/external/github.com/golang/oauth2/+/8f816d62a2652f705144857bbbcc26f2c166af9e/oauth2.go
*/

const stateCookieName = "oauthState"

func genRandState() string {
	log.Println("generating rand bytes")
	bytes := make([]byte, 32)
	// rand.Read(b)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatalf("failed to read rand fn %v", err)
	}
	state := base64.StdEncoding.EncodeToString(bytes)
	return state
}

//Server is the component of our app
type Server struct {
	cfg        *env.TomlConfig
	client     *spotifyapi.Client
	router     *chi.Mux
	httpServer *http.Server
}

//NewServer returns a configured new spotify client server
func NewServer(fileName string) *Server {
	s := &Server{}
	s.initCfg(fileName)
	s.initClient()
	s.routes()
	return s
}

//TODO: make a filePathErr for initCfg

func (s *Server) initCfg(fileName string) {
	cfg, err := env.LoadTOMLFile(fileName)
	if err != nil {
		log.Fatal("Error loading .toml file into struct config")
	}
	s.cfg = cfg

}

func (s *Server) initClient() {
	cfg, ok := s.cfg.Oauth2Providers["spotify"]
	if !ok {
		// TODO: Properly handle error
		panic("Spotify env properties not found in config")
	}
	s.client = spotifyapi.NewClient(
		cfg.ClientID,
		cfg.ClientSecret,
		cfg.RedirectURL)
}

//routes inits the route multiplexer with the assigned routes
func (s *Server) routes() {
	s.router = chi.NewRouter()

	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	})
	s.router.Use(
		middleware.Logger,       // Log API request calls
		middleware.StripSlashes, // Strip slashes to no slash URL versions
		middleware.RealIP,
		middleware.Recoverer, // Recover from panics without crashing server
		cors.Handler,         // Enable CORS globally
	)
	// Index handler
	//TODO: setup home page
	s.router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("hi"))
	})

	//test connection
	s.router.Get("/ping", s.handlePing()) //GET /ping
	//home
	s.router.Get("/login", s.handleLoginPage())

	//serve static files
	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "data"))
	FileServer(s.router, "/templates", filesDir)

	//account signin with Spotify
	// s.router.Get("/accounts/signup")
	// s.router.Mo
	s.router.Get("/auth", s.handleOauthProviderLogin())
	//below we have our redirect callback as a result of a user-agent accessing
	//our /auth endpoint route
	s.router.Get("/auth/callback", s.handleOauthProviderCallback())

	s.router.Get("/logout/{provider}", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

}

/*
fn that takes a Spotify URI, parses it with strings lib

*/

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
Note:
A struct or object will be http.Handler if it has one method
ServeHTTP which takes ResponseWriter and pointer to Request.
*/

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

//credit to uber-go guide on verifying interface compliance at compile time
//https://github.com/uber-go/guide/blob/master/style.md#guidelines
//this statement will fail if *Server ever stops matching the http.Handler interface
var _ http.Handler = (*Server)(nil)

//Start starts the server
func (s *Server) Start() {
	s.httpServer = &http.Server{
		Addr: ":" + s.cfg.Server.Port,
		// Handler:      s.router,
		Handler:      s,
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

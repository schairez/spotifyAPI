package spotifyservice

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/schairez/spotifywork/auth"
	"github.com/schairez/spotifywork/env"
	"golang.org/x/oauth2"
)

/*
what does
err := r.ParseForm()
do?

*/

const stateCookieName = "oauthState"

func genRandStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 64)
	rand.Read(b)
	state := base64.StdEncoding.EncodeToString(b)
	//TODO: suitable expiration time
	// expire in 1 month; so client's browser saves cookie in local file system
	// expiration := time.Now().Add(30 * 24 * time.Hour)
	expiration := time.Now().Add(time.Hour)

	//httpOnly security flag to secure our cookie from XSS
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

//Server is the server component of our app represented
// as a struct
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
	s.spotifyCfg = auth.NewSpotifyConfig(
		spotify.ClientID,
		spotify.ClientSecret,
		spotify.RedirectURL)
}

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
		localState := genRandStateOauthCookie(w)
		fmt.Println(localState)
		fmt.Println(w.Header())
		authURL := s.spotifyCfg.AuthCodeURL(localState)
		//app directs user-agent to spotify's oauth2 auth page
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)

	})
	s.router.Get("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		//check if user denied our auth request the request we receive
		//would contain a non-empty error query param in this case
		if r.FormValue("error") != "" {
			log.Printf("user authorization failed. Reason=%s", r.FormValue("error"))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		//TODO: pkce opts?
		authCode := r.FormValue("code")
		log.Printf("code=%s", authCode)
		//exchange auth code with an access token
		token, err := s.spotifyCfg.Exchange(context.Background(), authCode)
		if err != nil {
			log.Printf("error converting auth code into token; %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//we'll use the token to access user's protected resources
		// by calling the Spotify Web API

		log.Println(token)
		log.Println("query params?")
		queryParams := r.URL.Query()
		log.Println(queryParams)

		//now we can use this token to call Spotify APIs on behalf of the user
		//use the token to get an authenticated client

		// client :=

	})

	s.router.Get("/logout/{provider}", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

}

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
	if err := s.httpServer.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Println("Server closed!")
	}

}

//Shutdown the server
func (s *Server) Shutdown() {

}

//New returns a new route multiplexer with the assigned routes

/*
func New(spotifyCfg *oauth2.Config) *chi.Mux {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.StripSlashes)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(cors.New(cors.Options{
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	}).Handler)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("pong"))
	})

	r.Get("/auth", handleOauthSpotifyLogin(spotifyCfg))
	r.Get("/auth/callback", func(w http.ResponseWriter, r *http.Request) {
		//check if user denied our auth request the request we receive
		//would contain a non-empty error query param in this case
		if r.FormValue("error") != "" {
			log.Printf("user authorization failed. Reason=%s", r.FormValue("error"))
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		}
		//check the state parameter we supplied to Spotify's Account's Service earlier
		//if user approved the auth, we'll have both a code and a state query param
		oauthStateCookie, err := r.Cookie("state")
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
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		//TODO: pkce opts?
		authCode := r.FormValue("code")
		log.Printf("code=%s", authCode)
		//exchange auth code with an access token
		token, err := spotifyCfg.Exchange(context.Background(), authCode)
		if err != nil {
			log.Printf("error converting auth code into token; %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		//we'll use the token to access user's protected resources
		// by calling the Spotify Web API

		log.Println(token)
		log.Println("query params?")
		queryParams := r.URL.Query()
		log.Println(queryParams)

		//now we can use this token to call Spotify APIs on behalf of the user
		//use the token to get an authenticated client

		// client :=

	})

	r.Get("/logout/{provider}", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	return r

}

*/

//how to pass config values to haandler

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

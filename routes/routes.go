package routes

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
	"golang.org/x/oauth2"
)

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

//New returns a new route multiplexer with the assigned routes
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

	r.Get("/auth", func(w http.ResponseWriter, r *http.Request) {
		localState := genRandStateOauthCookie(w)
		fmt.Println(localState)
		fmt.Println(w.Header())
		authURL := spotifyCfg.AuthCodeURL(localState)
		//app directs user-agent to spotify's oauth2 auth page
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)

	})
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
			log.Println("invalid oauth2 spotify state")
			http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
			return
		}
		//pkce opts?
		authCode := r.FormValue("code")
		log.Printf("code=%s", authCode)
		token, err := spotifyCfg.Exchange(context.Background(), authCode)
		if err != nil {
			log.Printf("error converting auth code into token; %s", err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println(token)
		log.Println("query params?")
		queryParams := r.URL.Query()
		log.Println(queryParams)

		//now we can use this token to call Spotify APIs on behalf of the user
		// client :=

	})

	r.Get("/logout/{provider}", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Location", "/")
		w.WriteHeader(http.StatusTemporaryRedirect)
	})

	return r

}

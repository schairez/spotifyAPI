package spotifyservice

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	"github.com/schairez/spotifywork/internal"
	"github.com/schairez/spotifywork/spotifyservice/spotifyapi"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func (s *Server) handlePing() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("pong"))

	}
}

func (s *Server) handleLoginPage() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := templates.ExecuteTemplate(w, "login", nil)
		if err != nil {
			log.Fatal(err)
		}
	}
}

/*
https://mohitkhare.me/blog/sessions-in-golang/
TODO:
- CREATE A /profile route, add session middleware
*/

func (s *Server) handleOauthProviderLogin() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		log.Println(ctx)
		//check if the request contains a cookie?
		//COOKIE would be attached if the use has hit our domain
		//this would indicate that a user-agent has hit this endpoint,
		//but not that the user has authorized our app per-say
		log.Println("checking if user already has a cookie stored in their browser")
		cookie, err := r.Cookie(stateCookieName)
		if err != nil {
			log.Printf("we got no cookie in request, %s", err)
		}
		fmt.Println(cookie)
		localState := genRandState()
		//setting the set-Cookie header in the writer
		//NOTE: headers need to be set before anything else set to the writer
		http.SetCookie(w, internal.NewCookie(stateCookieName, localState))
		fmt.Println(localState)
		fmt.Println(w.Header())
		authURL := s.client.Config.AuthCodeURL(localState)
		//app directs user-agent to spotify's oauth2 auth  consent page
		http.Redirect(w, r, authURL, http.StatusTemporaryRedirect)

	}

}

func (s *Server) handleOauthProviderCallback() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		//error param present and non-empty if user denied auth
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
				http.Redirect(w, r, "/", http.StatusUnauthorized)
			}
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Printf("%s=%s\r\n", oauthStateCookie.Name, oauthStateCookie.Value)
		if r.FormValue("state") != oauthStateCookie.Value {
			log.Println("invalid oauth2 spotify state. state_mismatch err")
			//http.Error(w, "state_mismatch err", http.StatusUnauthorized)
			http.Redirect(w, r, "/", http.StatusUnauthorized)

			return
		}
		//TODO: pkce opts?
		authCode := r.FormValue("code")
		if authCode == "" {
			w.Write([]byte("No code"))
			w.WriteHeader(400)
			return
		}

		log.Printf("code=%s", authCode)
		//TODO: diff b/w background and oauth2.NoContext
		ctx := context.Background()
		//exchange auth code with an access token
		token, err := s.client.Config.Exchange(ctx, authCode)

		if err != nil {
			log.Printf("error converting auth code into token; %s", err.Error())
			http.Error(w, err.Error(), http.StatusForbidden)
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
		user, err := s.client.GetUserProfileRequest(context.Background(), token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("getting user")
		log.Printf("%+v\n", user)

		// data, _ := ioutil.ReadAll(resp.Body)
		// log.Println("Data calling user API: ", string(data))
		limit := 50
		offset := 0
		market := "us"
		params := spotifyapi.QParams{Limit: &limit, Offset: &offset, Market: &market}
		tracks, err := s.client.GetUserSavedTracks(context.Background(), token, &params)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println("getting user tracks")
		b, err := json.MarshalIndent(*tracks, "", "  ")
		if err != nil {
			fmt.Println(err)
		}
		fmt.Print(string(b))
		// log.Printf("%+v\n", tracks)

	}

}

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

//below from
//https://github.com/go-chi/chi/blob/master/_examples/fileserver/main.go

// FileServer conveniently sets up a http.FileServer handler to serve static
// files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"

	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}

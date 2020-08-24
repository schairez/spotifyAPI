package spotifyservice

import (
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
)

var templates = template.Must(template.ParseGlob("templates/*"))

func (s *Server) handlePing(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("pong"))

}

func (s *Server) handleHome(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "login", nil)
	if err != nil {
		log.Fatal(err)
	}

}

// func (s *Server) handleOauthLogin() {

// }

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

const indexPage = `
<div class="container">
<h2 style="text-align:center"> Login with Spotify </h2> 
<a href="/auth?provider=spotify" class="spotify btn"><span class="fa fa-google"></span> SignIn with Spotify</a>
</div>
`

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

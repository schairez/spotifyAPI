package spotifyservice

import (
	"net/http"
	"os"
	"path/filepath"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
)

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

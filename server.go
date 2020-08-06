package main

import (
	"log"
	"net/http"
	"time"

	"github.com/schairez/spotifywork/auth"
	"github.com/schairez/spotifywork/env"
	"github.com/schairez/spotifywork/routes"
)

func main() {
	fileName := "config.toml"
	cfg, err := env.LoadTOMLFile(fileName)
	if err != nil {
		panic(err)
	}
	spotify, ok := cfg.Oauth2Providers["spotify"]
	if !ok {
		panic("Spotify env properties not found in config")
	}
	SpotifyOauthConfig := auth.NewSpotifyConfig(
		spotify.ClientID,
		spotify.ClientSecret,
		spotify.RedirectURL)

	// fmt.Println(SpotifyOauthConfig)
	appRouter := routes.New(SpotifyOauthConfig)
	server := &http.Server{
		Addr:         ":8000",
		Handler:      appRouter,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	// panic(server.ListenAndServe())
	if err := server.ListenAndServe(); err != http.ErrServerClosed {
		log.Printf("%v", err)
	} else {
		log.Println("Server closed!")
	}

}

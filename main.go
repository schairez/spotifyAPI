package main

import (
	"github.com/schairez/spotifywork/spotifyservice"
	// "github.com/schairez/spotifywork/routes"
)

func main() {
	fileName := "config.toml"
	server := spotifyservice.NewServer(fileName)
	server.Start()

}

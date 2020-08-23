// Copyright 2020 Sergio Chairez. All rights reserved.
// Use of this source code is governed by a MIT style license that can be found
// in the LICENSE file.

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

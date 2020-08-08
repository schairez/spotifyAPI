// Package env handles init of our application-wide configuration.
package env

import (
	"github.com/BurntSushi/toml"
)

//TomlConfig is the struct representation of our toml file
type TomlConfig struct {
	Server          serverInfo
	Oauth2Providers map[string]providerData
}

//serverInfo is the struct for the server config
type serverInfo struct {
	Host string
	Port string
}

//providerData is the struct for the Oauth2 config
type providerData struct {
	ClientID     string
	ClientSecret string
	RedirectURL  string
}

//LoadTOMLFile loads client credentials file into struct
func LoadTOMLFile(fileName string) (*TomlConfig, error) {
	var config TomlConfig
	if _, err := toml.DecodeFile(fileName, &config); err != nil {
		return nil, err
	}
	return &config, nil

}

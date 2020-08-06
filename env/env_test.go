package env

import (
	"testing"
)

func TestBuildTOMLFile(t *testing.T) {
	fileName := "config.example.toml"
	cfg, err := LoadTOMLFile(fileName)
	if err != nil {
		t.Errorf("Failed to read toml file into struct %v", err)

	}
	for ProviderName, data := range cfg.Oauth2Providers {
		t.Logf("Provider: %s (%s, %s %s)\n", ProviderName, data.ClientID, data.ClientSecret, data.RedirectURL)
	}

}

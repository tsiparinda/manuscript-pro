package utils

import (
	"github.com/markbates/goth"
	"github.com/markbates/goth/providers/github"
	"github.com/markbates/goth/providers/google"
	"github.com/vedicsociety/platform/config"
)

func InitProviders(cfg config.Configuration) {
	// init oauth providers
	hostname, _ := cfg.GetString("system:hostname")
	clientID, _ := cfg.GetString("oauth:GoogleClientId")
	clientSecret, _ := cfg.GetString("oauth:GoogleClientSecret")
	goth.UseProviders(
		google.New(clientID, clientSecret, hostname+"/oauthcallback/google"),
	)

	hostname, _ = cfg.GetString("system:hostname")
	clientID, _ = cfg.GetString("oauth:GithubClientId")
	clientSecret, _ = cfg.GetString("oauth:GithubClientSecret")
	goth.UseProviders(
		github.New(clientID, clientSecret, hostname+"/oauthcallback/github"),
	)
}

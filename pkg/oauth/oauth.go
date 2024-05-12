package oauth

import (
	"github.com/VATUSA/primary-api/pkg/config"
	"golang.org/x/oauth2"
)

var OAuthConfig *oauth2.Config

func Initialize(config *config.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.OAuth.ClientID,
		ClientSecret: config.OAuth.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.OAuth.BaseURL + "/authorize",
			TokenURL: config.OAuth.BaseURL + "/token",
		},
		RedirectURL: config.API.BaseURL + "/v3/user/login/callback",
		Scopes:      []string{"identify", "email"},
	}
}

package oauth

import (
	"github.com/VATUSA/primary-api/pkg/config"
	"golang.org/x/oauth2"
)

var OAuthConfig *oauth2.Config
var DiscordOAuthConfig *oauth2.Config

func InitializeVATSIM(config *config.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.OAuth.ClientID,
		ClientSecret: config.OAuth.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.OAuth.BaseURL + "/oauth/authorize",
			TokenURL: config.OAuth.BaseURL + "/oauth/token",
		},
		RedirectURL: config.API.BaseURL + "/v3/user/login/callback",
		Scopes:      []string{"identify", "email"},
	}
}

func InitializeDiscord(config *config.Config) *oauth2.Config {
	return &oauth2.Config{
		ClientID:     config.DiscordOAuth.ClientID,
		ClientSecret: config.DiscordOAuth.ClientSecret,
		Endpoint: oauth2.Endpoint{
			AuthURL:  config.DiscordOAuth.BaseURL + "/oauth2/authorize",
			TokenURL: config.DiscordOAuth.BaseURL + "/api/oauth2/token",
		},
		RedirectURL: config.API.BaseURL + "/v3/user/discord/callback",
		Scopes:      []string{"identify"},
	}
}

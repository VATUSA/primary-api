package config

import (
	"os"
)

type OAuth struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	UserInfoURL  string
}

func NewOAuth() *OAuth {
	return &OAuth{
		BaseURL:      os.Getenv("OAUTH_BASE_URL"),
		ClientID:     os.Getenv("OAUTH_CLIENT_ID"),
		ClientSecret: os.Getenv("OAUTH_CLIENT_SECRET"),
		UserInfoURL:  os.Getenv("OAUTH_USER_INFO_URL"),
	}
}

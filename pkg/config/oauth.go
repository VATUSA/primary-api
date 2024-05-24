package config

type OAuth struct {
	BaseURL      string
	ClientID     string
	ClientSecret string
	UserInfoURL  string
}

func NewOAuth() *OAuth {
	return &OAuth{
		BaseURL:      EnvOrDefault("OAUTH_BASE_URL", defaultCfg.OAuth.BaseURL),
		ClientID:     EnvOrDefault("OAUTH_CLIENT_ID", defaultCfg.OAuth.ClientID),
		ClientSecret: EnvOrDefault("OAUTH_CLIENT_SECRET", defaultCfg.OAuth.ClientSecret),
		UserInfoURL:  EnvOrDefault("OAUTH_USER_INFO_URL", defaultCfg.OAuth.UserInfoURL),
	}
}

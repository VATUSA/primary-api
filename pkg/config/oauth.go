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

func NewDiscordOAuth() *OAuth {
	return &OAuth{
		BaseURL:      EnvOrDefault("DISCORD_OAUTH_BASE_URL", defaultCfg.DiscordOAuth.BaseURL),
		ClientID:     EnvOrDefault("DISCORD_OAUTH_CLIENT_ID", defaultCfg.DiscordOAuth.ClientID),
		ClientSecret: EnvOrDefault("DISCORD_OAUTH_CLIENT_SECRET", defaultCfg.DiscordOAuth.ClientSecret),
		UserInfoURL:  EnvOrDefault("DISCORD_OAUTH_USER_INFO_URL", defaultCfg.DiscordOAuth.UserInfoURL),
	}
}

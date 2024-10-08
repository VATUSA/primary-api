package config

import (
	"github.com/gorilla/securecookie"
	"os"
)

var defaultCfg = defaultConfig()
var Cfg *Config

type Config struct {
	API          *APIConfig
	Database     *DBConfig
	Cors         *CorsConfig
	Cookie       *Cookie
	S3           *S3Config
	OAuth        *OAuth
	DiscordOAuth *OAuth
}

func New() *Config {
	return &Config{
		API:          NewAPIConfig(),
		Database:     NewDBConfig(),
		Cors:         NewCorsConfig(),
		Cookie:       NewCookie(),
		S3:           NewS3Config(),
		OAuth:        NewOAuth(),
		DiscordOAuth: NewDiscordOAuth(),
	}
}

func EnvOrDefault(key, def string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return def
}

func defaultConfig() *Config {
	return &Config{
		API: &APIConfig{
			BaseURL:     "https://api.vatusa.net",
			Port:        "3000",
			LoggerLevel: "warn",
		},
		Database: &DBConfig{
			Host:        "localhost",
			Port:        "3306",
			User:        "root",
			Password:    "",
			Database:    "vatusa",
			LoggerLevel: "warn",
		},
		Cors: &CorsConfig{
			AllowedOrigins: "https://my.vatusa.net",
		},
		Cookie: &Cookie{
			HashKey:  securecookie.GenerateRandomKey(64),
			BlockKey: securecookie.GenerateRandomKey(32),
			Domain:   "vatusa.net",
		},
		S3: &S3Config{
			Endpoint:  "https://digitaloceanspaces.com",
			Region:    "nyc3",
			Bucket:    "vatusa",
			AccessKey: "",
			SecretKey: "",
			BaseURL:   "https://cdn.vatusa.net",
		},
		OAuth: &OAuth{
			BaseURL:      "https://auth.vatusa.net",
			UserInfoURL:  "/oauth/userinfo",
			ClientID:     "",
			ClientSecret: "",
		},
		DiscordOAuth: &OAuth{
			BaseURL:      "https://discord.com",
			UserInfoURL:  "/api/users/@me",
			ClientID:     "",
			ClientSecret: "",
		},
	}
}

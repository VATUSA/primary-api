package config

import (
	"fmt"
	"github.com/gorilla/securecookie"
	"os"
)

var Cfg *Config

type Config struct {
	API      *APIConfig
	Database *DBConfig
	Cors     *CorsConfig
	Cookie   *Cookie
	S3       *S3Config
	OAuth    *OAuth
}

type APIConfig struct {
	BaseURL string
	Port    string
}

type DBConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	Database    string
	LoggerLevel string
}

type CorsConfig struct {
	AllowedOrigins string
}

type Cookie struct {
	HashKey  []byte
	BlockKey []byte
	Domain   string
}

type S3Config struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
}

func NewAPIConfig() *APIConfig {
	return &APIConfig{
		BaseURL: os.Getenv("API_BASE_URL"),
		Port:    os.Getenv("API_PORT"),
	}
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		Host:        os.Getenv("DB_HOST"),
		Port:        os.Getenv("DB_PORT"),
		User:        os.Getenv("DB_USER"),
		Password:    os.Getenv("DB_PASSWORD"),
		Database:    os.Getenv("DB_DATABASE"),
		LoggerLevel: os.Getenv("DB_LOGGER_LEVEL"),
	}
}

func NewCorsConfig() *CorsConfig {
	return &CorsConfig{
		AllowedOrigins: os.Getenv("CORS_ALLOWED_ORIGINS"),
	}
}

func NewCookie() *Cookie {
	c := &Cookie{
		HashKey:  []byte(os.Getenv("COOKIE_HASH_KEY")),
		BlockKey: []byte(os.Getenv("COOKIE_BLOCK_KEY")),
		Domain:   os.Getenv("COOKIE_DOMAIN"),
	}

	if len(c.HashKey) == 0 || len(c.BlockKey) == 0 {
		fmt.Println("No cookie keys found, generating new keys...")
		c.HashKey = securecookie.GenerateRandomKey(64)
		c.BlockKey = securecookie.GenerateRandomKey(32)
	}

	return c
}

func NewS3Config() *S3Config {
	return &S3Config{
		Endpoint:  os.Getenv("S3_ENDPOINT"),
		Region:    os.Getenv("S3_REGION"),
		AccessKey: os.Getenv("S3_ACCESS_KEY"),
		SecretKey: os.Getenv("S3_SECRET_KEY"),
		Bucket:    os.Getenv("S3_BUCKET"),
	}
}

func New() *Config {
	return &Config{
		API:      NewAPIConfig(),
		Database: NewDBConfig(),
		Cors:     NewCorsConfig(),
		Cookie:   NewCookie(),
		S3:       NewS3Config(),
		OAuth:    NewOAuth(),
	}
}

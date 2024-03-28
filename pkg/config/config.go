package config

import "os"

type Config struct {
	API      *APIConfig
	Database *DBConfig
	Cors     *CorsConfig
	S3       *S3Config
}

type APIConfig struct {
	Port string
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

type S3Config struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
}

func NewAPIConfig() *APIConfig {
	return &APIConfig{
		Port: os.Getenv("PORT"),
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
		S3:       NewS3Config(),
	}
}

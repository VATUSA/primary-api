package config

type CorsConfig struct {
	AllowedOrigins string
}

func NewCorsConfig() *CorsConfig {
	return &CorsConfig{
		AllowedOrigins: EnvOrDefault("CORS_ALLOWED_ORIGINS", defaultCfg.Cors.AllowedOrigins),
	}
}

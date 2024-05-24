package config

type APIConfig struct {
	BaseURL string
	Port    string
}

func NewAPIConfig() *APIConfig {
	return &APIConfig{
		BaseURL: EnvOrDefault("API_BASE_URL", defaultCfg.API.BaseURL),
		Port:    EnvOrDefault("API_PORT", defaultCfg.API.Port),
	}
}

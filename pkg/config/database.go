package config

type DBConfig struct {
	Host        string
	Port        string
	User        string
	Password    string
	Database    string
	LoggerLevel string
}

func NewDBConfig() *DBConfig {
	return &DBConfig{
		Host:        EnvOrDefault("DB_HOST", defaultCfg.Database.Host),
		Port:        EnvOrDefault("DB_PORT", defaultCfg.Database.Port),
		User:        EnvOrDefault("DB_USER", defaultCfg.Database.User),
		Password:    EnvOrDefault("DB_PASSWORD", defaultCfg.Database.Password),
		Database:    EnvOrDefault("DB_DATABASE", defaultCfg.Database.Database),
		LoggerLevel: EnvOrDefault("DB_LOGGER_LEVEL", defaultCfg.Database.LoggerLevel),
	}
}

package config

type S3Config struct {
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
	Bucket    string
	BaseURL   string
}

func NewS3Config() *S3Config {
	return &S3Config{
		Endpoint:  EnvOrDefault("S3_ENDPOINT", defaultCfg.S3.Endpoint),
		Region:    EnvOrDefault("S3_REGION", defaultCfg.S3.Region),
		Bucket:    EnvOrDefault("S3_BUCKET", defaultCfg.S3.Bucket),
		AccessKey: EnvOrDefault("S3_ACCESS", defaultCfg.S3.AccessKey),
		SecretKey: EnvOrDefault("S3_SECRET", defaultCfg.S3.SecretKey),
		BaseURL:   EnvOrDefault("S3_BASE_URL", defaultCfg.S3.BaseURL),
	}
}

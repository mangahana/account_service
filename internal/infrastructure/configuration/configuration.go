package configuration

import (
	"github.com/Netflix/go-env"
)

type ServerConfig struct {
	GrpcSocker string `env:"GRPC_SOCKET"`
	CdnBaseUrl string `env:"CDN_BASE_URL"`
}

type SMSConfig struct {
	ApiKey    string `env:"SMS_API_KEY"`
	ApiDomain string `env:"SMS_API_DOMAIN"`
}

type S3Config struct {
	Endpoint        string `env:"S3_ENDPOINT"`
	AccessKeyID     string `env:"S3_ACCESS_KEY_ID"`
	SecretAccessKey string `env:"S3_SECRET_ACCESS_KEY"`
	BucketName      string `env:"S3_BUCKET_NAME"`
	UseSSL          bool   `env:"S3_USE_SSL"`
}

type DBConfig struct {
	Host string `env:"DB_HOST"`
	Name string `env:"DB_NAME"`
	User string `env:"DB_USER"`
	Pass string `env:"DB_PASS"`
}

type Config struct {
	DB     DBConfig
	Server ServerConfig
	SMS    SMSConfig
	S3     S3Config
}

func Load() (*Config, error) {
	var cfg Config
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, err
}

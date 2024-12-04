package configuration

import (
	"github.com/Netflix/go-env"
)

type ServerConfig struct {
	GrpcSocker string `env:"GRPC_SOCKET"`
}

type SMSConfig struct {
	ApiKey    string `env:"SMS_API_KEY"`
	ApiDomain string `env:"SMS_API_DOMAIN"`
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
}

func Load() (*Config, error) {
	var cfg Config
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		return nil, err
	}
	return &cfg, err
}

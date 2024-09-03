package config

import "github.com/caarlos0/env/v6"

type Config struct {
	DB      ConfigDB
	Minio   ConfigMinio
	AppPort string `env:"APP_PORT" envDefault:"8080"`
}

func NewConfig() (*Config, error) {
	config := Config{}
	err := env.Parse(&config)
	if err != nil {
		return nil, err
	}
	return &config, nil
}

type ConfigDB struct {
	User     string `env:"DB_USER" envDefault:"postgres"`
	Password string `env:"DB_PASSWORD" envDefault:"postgres"`
	DBName   string `env:"DB_NAME" envDefault:"postgres"`
	Host     string `env:"DB_HOST" envDefault:"localhost"`
	Port     string `env:"DB_PORT" envDefault:"5432"`
}
type ConfigMinio struct {
	Endpoint  string `env:"MINIO_ENDPOINT" envDefault:"http://localhost:9000"`
	AccessKey string `env:"MINIO_ACCESS_KEY" envDefault:"minioadmin"`
	SecretKey string `env:"MINIO_SECRET_KEY" envDefault:"minioadmin"`
	Bucket    string `env:"MINIO_BUCKET_NAME" envDefault:"minio"`
	Location  string `env:"MINIO_LOCATION" envDefault:"us-east-1"`
}

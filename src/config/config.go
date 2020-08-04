package config

import (
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Db struct {
		Username string `yaml:"username" env:"DB_USERNAME" env-default:"cassandra"`
		Password string `yaml:"password" env:"DB_PASSWORD" env-default:"cassandra"`
		Keyspace string `yaml:"keyspace" env:"DB_KEYSPACE" env-default:"public"`
		Host     string `yaml:"host" env:"DB_HOST" env-default:"cassandra"`
	} `yaml:"db"`
	Smtp *SmtpConfig `yaml:"smtp"`
	Server struct {
		Address string `yaml:"address" env:"SERVER_ADDRESS" env-default:"0.0.0.0:8080"`
	} `yaml:"server"`
}

type SmtpConfig struct {
		Address  string `yaml:"address" env:"SMTP_ADDRESS"`
		From     string `yaml:"from" env:"SMTP_FROM"`
		Password string `yaml:"password" env:"SMTP_PASSWORD"`
}

func LoadConfig(path string) (*Config, error) {
	var c Config
	if err := cleanenv.ReadConfig(path, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

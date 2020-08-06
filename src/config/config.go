package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"path/filepath"
	"runtime"
)

type Config struct {
	Db struct {
		Username string `yaml:"username" env:"DB_USERNAME" env-default:"cassandra"`
		Password string `yaml:"password" env:"DB_PASSWORD" env-default:"cassandra"`
		Keyspace string `yaml:"keyspace" env:"DB_KEYSPACE" env-default:"public"`
		// When you run the program locally, you need the host to be `localhost`
		// When it's run with docker-compose, it needs to be `cassandra`, so we're
		// passing it in the `docker.env` file.
		Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
	} `yaml:"db"`
	Smtp   SmtpConfig `yaml:"smtp"`
	Server struct {
		Address string `yaml:"address" env:"SERVER_ADDRESS" env-default:"0.0.0.0:8080"`
	} `yaml:"server"`
}

type SmtpConfig struct {
	Address  string `yaml:"address" env:"SMTP_ADDRESS"`
	From     string `yaml:"from" env:"SMTP_FROM"`
	Password string `yaml:"password" env:"SMTP_PASSWORD"`
}

// Try to read variables from the config file.
// If it fails, read them from environment.
func GetConfig(filename string) (*Config, error) {
	var c Config
	path := getConfigPath(filename)
	if err := cleanenv.ReadConfig(path, &c); err != nil {
		if err := cleanenv.ReadEnv(&c); err != nil {
			return nil, err
		}
	}
	return &c, nil
}

// Return the path on disk to the configs
func getConfigPath(configFilename string) string {
	_, currentFilename, _, ok := runtime.Caller(1)
	if !ok {
		return ""
	}
	return filepath.Join(filepath.Dir(currentFilename), "../../configs/", configFilename)
}
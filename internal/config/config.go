package config

import (
	"flag"
	"os"
	"time"

	cleanenv "github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env         string        `yaml:"env" env-default:"local"`
	StoragePath string        `yaml:"storage_path" env-required:"true"`
	TokenTTL    time.Duration `yaml:"token_ttl" env-required:"true"`
	GRPC        GRPCConfig    `yaml:"grpc"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

func MustLoad() *Config {
	path := fetchConfigPath()
	if path == "" {
		panic("config path empty")
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		panic("config file doesn't exist:" + path)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		panic("failed to read config:" + err.Error())
	}

	return &cfg
}

// получение информации о пути до файла конфига из 2 источников(flag or .env)
// priority: flag > env > default
func fetchConfigPath() string {
	var flagPath string

	// --config="path to config.yaml"
	flag.StringVar(&flagPath, "config", "", "path to config file")
	flag.Parse()

	if flagPath != "" {
		return flagPath
	}

	// If not provided via flag, check the environment variable
	envPath := os.Getenv("CONFIG_PATH")
	if envPath != "" {
		return envPath
	}

	return "" // Or return a default path if you have one
}

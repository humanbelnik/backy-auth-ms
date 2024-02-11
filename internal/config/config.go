package config

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string         `yaml:"env"`
	Database DatabaseConfig `yaml:"database"`
	TokenTTL time.Duration  `yaml:"token_ttl"`
	GRPC     GRPCConfig     `yaml:"grpc"`
}

type DatabaseConfig struct {
	Name     string `yaml:"name" env:"PGDATABASE"`
	Host     string `yaml:"host" env:"PGHOST"`
	Port     string `yaml:"port" env:"PGPORT"`
	User     string `yaml:"user" env:"PGUSER"`
	Password string `env:"PGPASSWORD"`
	SSLMode  string `yaml:"sslmode" env:"PGSSLMODE"`
}

type GRPCConfig struct {
	Port    int           `yaml:"port"`
	Timeout time.Duration `yaml:"timeout"`
}

// 1. Fetch path to a config file from flag or ENV (flag > ENV).
// 2. Parse config from file to a struct.
//
// Notice, that any error at that stage causes panic() in order to prevent service
// from running with incorrect config and shut it down immediately.

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return ParseConfig(configPath)
}

func ParseConfig(configPath string) *Config {
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("cannot find config file")
	}

	var cfg Config
	fmt.Println(configPath)
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("error while parsing config from file:" + err.Error())
	}
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("error while parsing config from .env:" + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var path string
	flag.StringVar(&path, "config_path", "", "path to a config file") // --config_path=./path/to/config
	flag.Parse()

	if path == "" {
		path = os.Getenv("AUTH_MS_CONFIG_PATH")
	}
	return path
}

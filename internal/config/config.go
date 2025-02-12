package config

import (
	"flag"
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	GRPCConfig     `yaml:"grpc" env-required:"true"`
	PostgresConfig `yaml:"postgres" env-required:"true"`
}

type GRPCConfig struct {
	GRPCPort int `yaml:"port" env-required:"true"`
}

type PostgresConfig struct {
	Host     string `yaml:"host" env-required:"true"`
	DBPort   int    `yaml:"port" env-required:"true"`
	User     string `yaml:"user" env-required:"true"`
	Password string `yaml:"password" env-required:"true"`
	DBName   string `yaml:"dbname" env-required:"true"`
}

func MustLoad() *Config {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file not found: " + configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("error while loading config: " + err.Error())
	}

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")

	if res == "" {
		dotenvInit()
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}

func dotenvInit() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

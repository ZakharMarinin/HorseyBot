package config

import (
	"log"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type Config struct {
	Env      string         `yaml:"env" env-default:"local"`
	TgBot    TelegramBot    `yaml:"tg_bot"`
	Postgres PostgresConfig `yaml:"postgres"`
}

type PostgresConfig struct {
	Addr string `yaml:"addr" env-required:"true"`
}
type TelegramBot struct {
	TgToken string `yaml:"tg_token"`
}

func MustLoad() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	configPath := "./config.yaml"

	var cfg Config

	cfg.Postgres.Addr = os.Getenv("GOOSE_DBSTRING")
	cfg.TgBot.TgToken = os.Getenv("TG_TOKEN")

	err = cleanenv.ReadConfig(configPath, &cfg)
	if err != nil {
		log.Fatal("cannot find the config: ", err)
	}

	return &cfg
}

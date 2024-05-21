package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Server struct {
		Port string `yaml:"port" env:"SERVER_PORT" env-default:"8080"`
		Host string `yaml:"host" envconfig:"SERVER_HOST" env-default:"locahost"`
	} `yaml:"server"`

	Database struct {
		DBPath string `yaml:"dbpath" env:"DBPATH" env-default:"db.db"`
	} `yaml:"database"`
}

var Cfg Config

func init() {
	err := cleanenv.ReadConfig("config.yml", &Cfg)
	if err != nil {
		log.Panic(err)
	}
}

package config

import (
	"encoding/json"
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"log"
)

type Config struct {
	App struct {
		Debug bool `env:"DEBUG" env-default:"T"`
	}
	Server struct {
		Host             string `yaml:"host" env:"SERVER_HOST" env-default:"localhost"`
		Port             string `yaml:"port" env:"SERVER_PORT" env-default:"8081"`
		ReadTimeoutMSec  int    `yaml:"readtimeoutmsec" env:"SERVER_READ_TIMEOUT" env-default:"100000"`
		WriteTimeoutMSec int    `yaml:"writetimeoutmsec" env:"SERVER_WRITE_TIMEOUT" env-default:"100000"`
	} `yaml:"server"`
	DB struct {
		Host     string `yaml:"host" env:"DB_HOST" env-default:"localhost"`
		Port     string `yaml:"port" env:"DB_PORT" env-default:"5433"`
		User     string `yaml:"user" env:"DB_USER" env-default:"db_user"`
		Password string `yaml:"pass" env:"DB_PASSWORD" env-default:"tYSk4dqaW7Hq4cw2r4hP"`
		Name     string `yaml:"name" env:"DB_NAME" env-default:"omnimanage_db"`
	} `yaml:"db"`
}

var (
	config Config
)

func Get(path string) (*Config, error) {
	var err error
	if len(path) == 0 {
		err = cleanenv.ReadEnv(&config)
	} else {
		err = cleanenv.ReadConfig(path, &config)
	}
	if err != nil {
		return nil, err
	}

	configBytes, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Configuration:", string(configBytes))

	return &config, nil
}

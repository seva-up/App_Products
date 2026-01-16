package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Config struct {
	App App
	Db  DbPostgres
}

type App struct {
	Name    string `yaml:"name"`
	Version string `yaml:"version"`
	Env     string `yaml:"env"`
	Port    string `yaml:"port"`
	Host    string `yaml:"host"`
}
type DbPostgres struct {
	Name     string `yaml:"name"`
	User     string `yaml:"user"`
	Password string `yaml:"password"`
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
}

func Load() (*Config, error) {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()
	viper.SetDefault("host", "localhost")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("Конфиг не найден в функции config.go/Load()")
	}
	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		fmt.Println("Конфиг не смог за анмаршалиться")
	}
	return &config, nil
}

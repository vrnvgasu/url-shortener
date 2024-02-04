package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"log"
	"os"
	"time"
)

// полностью соответсвует структуре config/local.yaml
type Config struct {
	// yaml - названия поля в config/local.yaml
	// env - название параметра из переменной окружения (если читаем от туда)
	// env-required - приложение не запуститься, если пропустить параметр при установке конфигов
	Env         string `yaml:"env" env:"ENV" env-default:"local"`
	StoragePath string `yaml:"storage_path" env-required:"true"`
	HTTPServer  `yaml:"http_server"`
}

type HTTPServer struct {
	Address     string        `yaml:"address" env-default:"localhost:8080"`
	Timeout     time.Duration `yaml:"timeout" env-default:"4s"`
	IdleTimeout time.Duration `yaml:"idle_timeout" env-default:"60s"`
	User        string        `yaml:"user" env-required:"true"`
	Password    string        `yaml:"password" env-required:"true" env:"HTTP_SERVER_PASSWORD"`
}

// "Must" - сообщаем, что функция может кинуть панику
func MustLoad() *Config {
	// берем путь к конфигу из переменной окружения
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		log.Fatal("CONFIG_PATH is not set")
	}

	// check is file exist
	if _, err := os.Stat(configPath); err != nil {
		log.Fatalf("config file is not exist, %s", configPath)
	}

	var cfg Config

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatalf("can't read config, %s", err)
	}

	return &cfg
}

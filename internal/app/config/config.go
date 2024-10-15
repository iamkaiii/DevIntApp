package config

import (
	"fmt"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
	"strconv"
	"time"
)

// Config Структура конфигурации;
// Содержит все конфигурационные данные о сервисе;
// автоподгружается при изменении исходного файла

type Config struct {
	ServiceHost string
	ServicePort int

	Redis RedisConfig
}

type RedisConfig struct {
	Host        string
	Password    string
	Port        int
	User        string
	DialTimeout time.Duration
	ReadTimeout time.Duration
}

const (
	envRedisHost = "REDIS_HOST"
	envRedisPort = "REDIS_PORT"
	envRedisUser = "REDIS_USER"
	envRedisPass = "REDIS_PASSWORD"
)

// NewConfig Создаёт новый объект конфигурации, загружая данные из файла конфигурации
func NewConfig() (*Config, error) {
	var err error

	configName := "config"
	_ = godotenv.Load()
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}

	viper.SetConfigName(configName)
	viper.SetConfigType("toml")
	viper.AddConfigPath("config")
	viper.AddConfigPath(".")
	viper.WatchConfig()

	err = viper.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Config{}

	cfg.Redis.Host = os.Getenv(envRedisHost)
	cfg.Redis.Port, err = strconv.Atoi(os.Getenv(envRedisPort))

	if err != nil {
		return nil, fmt.Errorf("redis port must be int value: %w", err)
	}

	cfg.Redis.Password = os.Getenv(envRedisPass)
	cfg.Redis.User = os.Getenv(envRedisUser)

	err = viper.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	log.Info("config parsed")

	return cfg, nil
}

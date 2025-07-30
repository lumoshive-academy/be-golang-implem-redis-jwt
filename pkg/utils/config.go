package utils

import (
	"github.com/spf13/viper"
)

type Configuration struct {
	AppName     string
	Port        string
	Debug       bool
	DB          DatabaseConfig
	Limit       int
	PathLogger  string
	PathUpload  string
	RedisConfig RedisConfig
	PrivateKey  string
	PublicKey   string
}

type DatabaseConfig struct {
	Name         string
	Username     string
	Password     string
	Host         string
	Port         string
	TimeZone     string
	Logging      bool
	MaxIdleConns int
	MaxOpenConns int
	MaxIdleTime  int
	MaxLifeTime  int
}

type RedisConfig struct {
	Url      string
	Password string
	Prefix   string
}

func ReadConfiguration() (Configuration, error) {
	viper.SetConfigFile(".env") // read file .env
	viper.SetConfigType("env")
	err := viper.ReadInConfig()
	if err != nil {
		return Configuration{}, err
	}
	viper.AutomaticEnv() // read env os

	return Configuration{
		AppName:    viper.GetString("APP_NAME"),
		Port:       viper.GetString("PORT"),
		Debug:      viper.GetBool("DEBUG"),
		Limit:      viper.GetInt("LIMIT"),
		PathLogger: viper.GetString("PATH_LOGGER"),
		PathUpload: viper.GetString("PATH_UPLOAD"),
		PrivateKey: viper.GetString("PRIVATE_KEY"),
		PublicKey:  viper.GetString("PUBLIC_KEY"),
		DB: DatabaseConfig{
			Name:         viper.GetString("DATABASE_NAME"),
			Username:     viper.GetString("DATABASE_USER"),
			Password:     viper.GetString("DATABASE_PASSWORD"),
			Host:         viper.GetString("DATABASE_HOST"),
			Port:         viper.GetString("DATABASE_PORT"),
			TimeZone:     viper.GetString("DATABASE_TIME_ZONE"),
			MaxIdleConns: viper.GetInt("DB_MAX_IDLE_CONNS"),
			MaxOpenConns: viper.GetInt("DB_MAX_OPEN_CONNS"),
			MaxIdleTime:  viper.GetInt("DB_MAX_IDLE_TIME"),
			MaxLifeTime:  viper.GetInt("DB_MAX_LIFE_TIME"),
		},
		RedisConfig: RedisConfig{
			Url:      viper.GetString("REDIS_URL"),
			Password: viper.GetString("REDIS_PASSWORD"),
			Prefix:   viper.GetString("REDIS_PREFIX"),
		},
	}, nil
}

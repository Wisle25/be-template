package commons

import (
	"time"

	"github.com/spf13/viper"
)

var Cfg *Config

type Config struct {
	// Database
	DBHost         string `mapstructure:"POSTGRES_HOST"`
	DBUserName     string `mapstructure:"POSTGRES_USER"`
	DBUserPassword string `mapstructure:"POSTGRES_PASSWORD"`
	DBName         string `mapstructure:"POSTGRES_DB"`
	DBNameTest     string `mapstructure:"POSTGRES_DB_TEST"`
	DBPort         string `mapstructure:"POSTGRES_PORT"`
	RedisURL       string `mapstructure:"REDIS_URL"`

	// Server
	AppEnv       string `mapstructure:"APP_ENV"`
	ServerPort   string `mapstructure:"PORT"`
	ClientOrigin string `mapstructure:"CLIENT_ORIGIN"`

	// JWT Tokens
	AccessTokenPrivateKey string        `mapstructure:"ACCESS_TOKEN_PRIVATE_KEY"`
	AccessTokenPublicKey  string        `mapstructure:"ACCESS_TOKEN_PUBLIC_KEY"`
	AccessTokenExpiresIn  time.Duration `mapstructure:"ACCESS_TOKEN_EXPIRED_IN"`
	AccessTokenMaxAge     int           `mapstructure:"ACCESS_TOKEN_MAXAGE"`

	RefreshTokenPrivateKey string        `mapstructure:"REFRESH_TOKEN_PRIVATE_KEY"`
	RefreshTokenPublicKey  string        `mapstructure:"REFRESH_TOKEN_PUBLIC_KEY"`
	RefreshTokenExpiresIn  time.Duration `mapstructure:"REFRESH_TOKEN_EXPIRED_IN"`
	RefreshTokenMaxAge     int           `mapstructure:"REFRESH_TOKEN_MAXAGE"`
}

func LoadConfig(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Read the .env
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	err = viper.Unmarshal(&Cfg)
	if err != nil {
		panic(err)
	}
}

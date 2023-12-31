package config

import (
	"github.com/meisbokai/GolangApiTest/internal/constants"
	"github.com/spf13/viper"
)

var AppConfig Config

type Config struct {
	Port int `mapstructure:"PORT"`

	DBPostgreDriver string `mapstructure:"DB_POSTGRE_DRIVER"`
	DBPostgreDsn    string `mapstructure:"DB_POSTGRE_DSN"`

	JWTSecret  string `mapstructure:"JWT_SECRET"`
	JWTExpired int    `mapstructure:"JWT_EXPIRED"`
	JWTIssuer  string `mapstructure:"JWT_ISSUER"`
}

func InitializeAppConfig() error {
	viper.SetConfigName(".env") // allow directly reading from .env file
	viper.SetConfigType("env")
	viper.AddConfigPath(".")
	viper.AddConfigPath("internal/configs")
	viper.AddConfigPath("/")
	viper.AllowEmptyEnv(true)
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		return constants.ErrLoadConfig
	}

	err = viper.Unmarshal(&AppConfig)
	if err != nil {
		return constants.ErrParseConfig
	}

	// check
	if AppConfig.Port == 0 || AppConfig.JWTSecret == "" || AppConfig.JWTExpired == 0 || AppConfig.JWTIssuer == "" || AppConfig.DBPostgreDriver == "" {
		return constants.ErrEmptyVar
	}

	if AppConfig.DBPostgreDsn == "" {
		return constants.ErrEmptyVar
	}

	return nil
}

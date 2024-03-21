package utils

import (
	"time"

	"github.com/spf13/viper"
)

// Configuration struct that holds the configuration for the software
type Config struct {
	Env string `mapstructure:"ENV"`
	DBDriver string `mapstructure:"DBDRIVER"`
	DBURL   string `mapstructure:"DBURL"`
	DBName string `mapstructure:"DBNAME"`
	DBUser string `mapstructure:"DBUSER"`
	DBPassword string `mapstructure:"DBPASSWORD"`
	DBHost string `mapstructure:"DBHOST"`
	DBPort string `mapstructure:"DBPORT"`
	RedisAddress string `mapstructure:"REDIS_ADDRESS"`
	RedisPassword string `mapstructure:"REDIS_PASSWORD"`
	RedisDB int `mapstructure:"REDIS_DB"`
	MigrationURL string `mapstructure:"MIGRATION_URL"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
	GRPCServerAddress string `mapstructure:"GRPC_SERVER_ADDRESS"`
	LOG_LEVEL string `mapstructure:"LOG_LEVEL"`
	JWT_SECRET string `mapstructure:"JWT_SECRET"`
	JWT_ACCESS_TOKEN_DURATION time.Duration `mapstructure:"JWT_ACCESS_TOKEN_DURATION"`
	SHUTTER_PUBLIC_KEY string `mapstructure:"SHUTTER_PUBLIC_KEY"`
	BITPOWR_ACCOUNT_ID string `mapstructure:"BITPOWR_ACCOUNT_ID"`
	BITPOWR_API_KEY string `mapstructure:"BITPOWR_API_KEY"`
	VALIDATE_BANK_URL string `mapstructure:"VALIDATE_BANK_URL"`
}

// LoadConfig loads the configuration from the config file or environment variables
func LoadConfig(configFilePath string) (conf Config, err error) {
	viper.AddConfigPath(configFilePath)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	// viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&conf)
	return
}

package utils

import (
	"github.com/spf13/viper"
)

// Configuration struct that holds the configuration for the software
type config struct {
	DBDriver string `mapstructure:"DBDRIVER"`
	DBURL   string `mapstructure:"DBURL"`
	DBName string `mapstructure:"DBNAME"`
	DBUser string `mapstructure:"DBUSER"`
	DBPassword string `mapstructure:"DBPASSWORD"`
	DBHost string `mapstructure:"DBHOST"`
	DBPort string `mapstructure:"DBPORT"`
}

// LoadConfig loads the configuration from the config file or environment variables
func LoadConfig(configFilePath string) (conf config, err error) {
	viper.AddConfigPath(configFilePath)
	viper.SetConfigName(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err = viper.ReadInConfig()
	if err != nil {
		return
	}
	err = viper.Unmarshal(&conf)
	return
}

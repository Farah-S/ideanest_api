package config

import (
	"github.com/spf13/viper"
)

// ? Struct Config
type DatabaseConfig struct {
	DBUri        string
	DBName     string
	Collections []string
	RedisUri	string
	// DBUri    string `mapstructure:"MONGODB_LOCAL_URI"`
	// RedisUri string `mapstructure:"REDIS_URL"`
	// Port     string `mapstructure:"PORT"`
}

type AppConfig struct {

	Port	string
	AppName 	string
	// DBUri    string `mapstructure:"MONGODB_LOCAL_URI"`
	// RedisUri string `mapstructure:"REDIS_URL"`
	// Port     string `mapstructure:"PORT"`
}

// The function `LoadConfig` loads a configuration file from the specified path and unmarshals it into
// a struct.
func LoadDBConfig(path string) (config DatabaseConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("yml")
	viper.SetConfigName("database-config")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

func LoadAppConfig(path string) (config AppConfig, err error) {
	viper.AddConfigPath(path)
	viper.SetConfigType("yml")
	viper.SetConfigName("app-config")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}
package config

import (
	"log"
	"os"
	"strconv"

	"github.com/spf13/viper"
)

func init() {
	log.Println("Initializing configuration setup")
	env := os.Getenv("ENVIRONMENT")
	if env == "" || env == "DEVELOPMENT" {
		viper.SetConfigFile(".env")
		if err := viper.ReadInConfig(); err != nil {
			log.Panicf("Error reading config file, %s", err)
		}
		viper.SetDefault("ENVIRONMENT", "DEVELOPMENT")
	}

	if env == "PRODUCTION" {
		viper.AutomaticEnv()
	}

}

type Config struct {
	SuperUserName        string
	SuperUserPass        string
	Environment          string
	Port                 string
	Prefork              bool
	DBURL                string
	JWTSecret            string
	JWTAccessExpiryHours int
	JWTRefreshExpiryDays int
}

var config *Config

func LoadConfig() {
	config = &Config{
		SuperUserName:        viper.GetString("SUPERUSERNAME"),
		SuperUserPass:        viper.GetString("SUPERUSERPASS"),
		Environment:          viper.GetString("ENVIRONMENT"),
		Port:                 viper.GetString("PORT"),
		Prefork:              viper.GetBool("PREFORK"),
		DBURL:                viper.GetString("DATABASE_URL"),
		JWTSecret:            viper.GetString("JWT_SECRET"),
		JWTAccessExpiryHours: getEnvInt("JWT_ACCESS_EXPIRY_HOURS", 1),
		JWTRefreshExpiryDays: getEnvInt("JWT_REFRESH_EXPIRY_DAYS", 7),
	}
}

func GetConfig() *Config {
	return config
}

// Helper function para converter string para int com valor padrão
func getEnvInt(key string, defaultValue int) int {
	if value := viper.GetString(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

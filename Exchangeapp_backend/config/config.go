package config

import (
	"log"
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	App struct {
		Name string
		Port string
	}
	Database struct {
		Dsn          string
		MaxIdleConns int
		MaxOpenConns int
	}
	Redis struct {
		Addr     string
		DB       int
		Password string
	}
	JWT struct {
		Secret string
	}
}

var AppConfig *Config

func InitConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yml")
	viper.AddConfigPath("./config")

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %v", err)
	}

	AppConfig = &Config{}

	if err := viper.Unmarshal(AppConfig); err != nil {
		log.Fatalf("Unable to decode into struct: %v", err)
	}

	// 生产环境可用环境变量 JWT_SECRET 覆盖，避免所有环境共用同一密钥
	if s := os.Getenv("JWT_SECRET"); s != "" {
		AppConfig.JWT.Secret = s
	}

	initDB()
	initRedis()
}

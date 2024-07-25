package config

import (
	"github.com/joho/godotenv"
	"os"
)

type Config struct {
	AppName       string
	AppEnv        string
	AppDebug      bool
	AppURL        string
	AppHost       string
	AppPort       string
	LogChannel    string
	LogLevel      string
	DBConnection  string
	DBHost        string
	DBPort        string
	DBName        string
	DBUser        string
	DBPassword    string
	RedisHost     string
	RedisPassword string
	RedisPort     string
	TmpFolderPath string
	JWTSecret     string
}

func LoadConfig() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		return nil, err
	}

	return &Config{
		AppName:       os.Getenv("APP_NAME"),
		AppEnv:        os.Getenv("APP_ENV"),
		AppDebug:      os.Getenv("APP_DEBUG") == "true",
		AppURL:        os.Getenv("APP_URL"),
		AppHost:       os.Getenv("APP_HOST"),
		AppPort:       os.Getenv("APP_PORT"),
		LogChannel:    os.Getenv("LOG_CHANNEL"),
		LogLevel:      os.Getenv("LOG_LEVEL"),
		DBConnection:  os.Getenv("DB_CONNECTION"),
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBName:        os.Getenv("DB_NAME"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		RedisHost:     os.Getenv("REDIS_HOST"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisPort:     os.Getenv("REDIS_PORT"),
		TmpFolderPath: os.Getenv("TMP_FOLDER_PATH"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
	}, nil
}

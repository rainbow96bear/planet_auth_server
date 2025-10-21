package config

import (
	"fmt"
	"os"
	"planet_utils/pkg/logger"
	"strconv"

	"github.com/joho/godotenv"
)

func InitConfig(mode string) {
	var err error

	switch mode {
	case "prod":
		err = godotenv.Load("./env/.env.prod")
	case "dev":
		err = godotenv.Load("./env/.env.dev")
	}

	if err != nil {
		fmt.Println("[CONFIG] fail to load .env file, 기본값 dev 사용")
	}

	// default config
	PORT = getString("PORT")
	LOG_LEVEL = getInt16("LOG_LEVEL")
	ACCESS_TOKEN_EXPIRY_MINUTE = getInt("ACCESS_TOKEN_EXPIRY_MINUTE")
	REFRESH_TOKEN_NAME = getString("REFRESH_TOKEN_NAME")
	REFRESH_TOKEN_EXPIRY_DURATION = getInt("REFRESH_TOKEN_EXPIRY_DURATION")

	OAUTH_SESSION_EXPIRY = getInt("OAUTH_SESSION_EXPIRY")

	// kakao config
	KAKAO_REST_API_KEY = getString("KAKAO_REST_API_KEY")
	KAKAO_REDIRECT_URI = getString("KAKAO_REDIRECT_URI")
	KAKAO_CLIENT_SECRET = getString("KAKAO_CLIENT_SECRET")

	// jwt key
	JWT_SECRET_KEY = getString("JWT_SECRET_KEY")
	PLANET_CLIENT_ADDR = getString("PLANET_CLIENT_ADDR")

	DB_USER = getString("DB_USER")
	DB_PASSWORD = getString("DB_PASSWORD")
	DB_HOST = getString("DB_HOST")
	DB_PORT = getString("DB_PORT")
	DB_NAME = getString("DB_NAME")
}

func getString(envName string) string {
	v := os.Getenv(envName)
	if v == "" {
		logger.Errorf("[CONFIG] %s not set\n", envName)
		os.Exit(1)
	}
	return v
}

func getInt(envName string) int {
	v := os.Getenv(envName)
	if v == "" {
		logger.Errorf("[CONFIG] %s not set\n", envName)
		os.Exit(1)
	}
	num, err := strconv.Atoi(v)
	if err != nil {
		logger.Errorf("[CONFIG] %s must be int, got %s\n", envName, v)
		os.Exit(1)
	}
	return num
}

func getInt16(envName string) int16 {
	v := os.Getenv(envName)
	if v == "" {
		logger.Errorf("[CONFIG] %s not set\n", envName)
		os.Exit(1)
	}
	num, err := strconv.Atoi(v)
	if err != nil {
		logger.Errorf("[CONFIG] %s must be int, got %s\n", envName, v)
		os.Exit(1)
	}
	return int16(num)
}

func getUint64(envName string) uint64 {
	v := os.Getenv(envName)
	if v == "" {
		logger.Errorf("[CONFIG] %s not set\n", envName)
		os.Exit(1)
	}

	num, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		logger.Errorf("[CONFIG] %s must be a valid uint64, got %s\n", envName, v)
		os.Exit(1)
	}

	return num
}

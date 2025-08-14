package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfigStruct struct {
	Port     int
	LogLevel int16
}

var AppConfig AppConfigStruct

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

	if portStr := os.Getenv("PORT"); portStr != "" {
		if AppConfig.Port, err = strconv.Atoi(portStr); err != nil {
			fmt.Println("fail to set port")
			os.Exit(1)
		}
	} else {
		os.Exit(1)
	}

	if logLevelStr := os.Getenv("LOG_LEVEL"); logLevelStr != "" {
		if logLevel, err := strconv.Atoi(logLevelStr); err == nil {
			AppConfig.LogLevel = int16(logLevel)
		} else {
			fmt.Println("fail to set log level")
			AppConfig.LogLevel = 1
		}
	} else {
		os.Exit(1)
	}
	fmt.Printf("Success to set AppConfig : %+v\n", AppConfig)
}

package config

import (
	"log/slog"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	AppTimeOut time.Duration

	DBMS       string
	DBPassword string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
}

func LoadConf() *Config {
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using environment variables")
	}
	timeOutStr := os.Getenv("TIMEOUT_SECOND")
	timeOut, err := strconv.Atoi(timeOutStr)
	if err != nil {
		slog.Error("invalid TIMEOUT_SECOND", "value", timeOutStr, "error", err)
		panic(err)
	}
	return &Config{
		AppTimeOut: time.Duration(timeOut) * time.Second,
		DBMS:       os.Getenv("DBMS"),
		DBPassword: os.Getenv("MYSQL_PASSWORD"),
		DBHost:     os.Getenv("MYSQL_DBHOST"),
		DBPort:     os.Getenv("MYSQL_DBPORT"),
		DBName:     os.Getenv("MYSQL_DATABASE"),
		DBUser:     os.Getenv("MYSQL_USER"),
	}
}

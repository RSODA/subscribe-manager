package config

import "os"

const (
	LoggerLevelDev  LoggerLevel = "dev"
	LoggerLevelProd LoggerLevel = "prod"
)

type LoggerLevel string

func NewLoggerLevel() LoggerLevel {
	l := os.Getenv("LOGGER_LEVEL")

	switch LoggerLevel(l) {
	case LoggerLevelDev:
		return LoggerLevelDev
	case LoggerLevelProd:
		return LoggerLevelProd
	default:
		return LoggerLevelDev
	}
}

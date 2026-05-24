package logger

import (
	"go.uber.org/zap"
)

func NewLogger(level string) (*zap.SugaredLogger, error) {
	var l *zap.Logger
	var err error

	switch level {
	case "dev":
		l, err = zap.NewDevelopment()
	case "prod":
		l, err = zap.NewProduction()
	default:
		l, err = zap.NewDevelopment()
	}

	if err != nil {
		return nil, err
	}

	return l.Sugar(), nil
}

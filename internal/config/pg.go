package config

import (
	"errors"
	"fmt"
	"os"
)

type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	Username string
	Password string
	Database string
	Port     string
}

func NewPGConfig() (PGConfig, error) {
	username := os.Getenv("PG_USERNAME")
	if len(username) == 0 {
		return nil, errors.New("PG_USERNAME is required")
	}

	password := os.Getenv("PG_PASSWORD")
	if len(password) == 0 {
		return nil, errors.New("PG_PASSWORD is required")
	}

	database := os.Getenv("PG_DATABASE")
	if len(database) == 0 {
		return nil, errors.New("PG_DATABASE is required")
	}

	port := os.Getenv("PG_PORT")
	if len(port) == 0 {
		return nil, errors.New("PG_PORT is required")
	}

	return &pgConfig{
		Username: username,
		Password: password,
		Database: database,
		Port:     port,
	}, nil
}

func (c *pgConfig) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/", c.Username, c.Password, c.Port, c.Database)
}

package config

import "os"

type HTTPConfig interface {
	Address() string
}

type httpConfig struct {
	host string
	port string
}

func NewHTTPConfig() HTTPConfig {
	host := os.Getenv("HTTP_HOST")
	if len(host) == 0 {
		host = "localhost"
	}

	port := os.Getenv("HTTP_PORT")
	if len(port) == 0 {
		port = "8080"
	}

	return &httpConfig{
		host: host,
		port: port,
	}
}

func (c *httpConfig) Address() string {
	return c.host + ":" + c.port
}

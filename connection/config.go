package connection

import (
	"net/http"
	"net/url"
)

//nolint:unused // reserved for future use
const (
	userAgent = "agentsandbox-go-sdk"
)

type Config struct {
	Domain      string
	AccessToken string
	Headers     http.Header
	Proxy       *url.URL
}

func NewConfig() *Config {
	return &Config{}
}

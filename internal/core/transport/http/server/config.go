package core_http_server

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type HTTPServerConfig struct {
	Addr            string        `envconfig:"ADDRR" required:"true"`
	ShutdownTimeout time.Duration `envconfig:"SHUTDOWN_TIMEOUT" required:"true"`
	CertFile        string        `envconfig:"CERT_FILE" required:"true"`
	KeyFile         string        `envconfig:"KEY_FILE" required:"true"`
}

func NewHTTPServerConfig() (HTTPServerConfig, error) {
	var config HTTPServerConfig

	if err := envconfig.Process("HTTP_SERVER", &config); err != nil {
		return HTTPServerConfig{}, fmt.Errorf("envconfig process error: %w", err)
	}

	return config, nil
}

func NewHTTPServerConfigMust() HTTPServerConfig {
	config, err := NewHTTPServerConfig()
	if err != nil {
		panic(err)
	}

	return config
}

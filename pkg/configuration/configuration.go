package configuration

import (
	"github.com/kelseyhightower/envconfig"
)

// Configuration ...
type Configuration interface {
	Load(cfg interface{}) (err error)
}

type configuration struct {
}

func (c *configuration) Load(cfg interface{}) error {
	return envconfig.Process("", cfg)
}

// NewConfiguration ...
func NewConfiguration() Configuration {
	return &configuration{}
}

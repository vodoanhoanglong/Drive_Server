package env

import (
	"log"

	"github.com/kelseyhightower/envconfig"
	"nexlab.tech/core/pkg/gql"
	"nexlab.tech/core/services/auth/utils"
)

// headers constants
const (
	HasuraClientName = "hasura-client-name"
	clientName       = "auth"
)

// Environment variables data
type Environment struct {
	Port             string           `envconfig:"PORT" default:"8080"`
	LogLevel         string           `envconfig:"LOG_LEVEL" default:"info"`
	DefaultRole      string           `envconfig:"DEFAULT_ROLE" required:"true"`
	ControllerClient gql.ClientConfig `envconfig:"CONTROLLER" required:"true"`
	JWT              utils.JWTAuthConfig
	Email            string `envconfig:"EMAIL" required:"true"`
	Password         string `envconfig:"EMAIL_PASSWORD" required:"true"`
}

// GetEnv initialize and return environment variables
func GetEnv() *Environment {
	var env Environment

	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err)
	}

	if env.ControllerClient.Headers == nil {
		env.ControllerClient.Headers = map[string]string{
			HasuraClientName: clientName,
		}
	} else if _, ok := env.ControllerClient.Headers[HasuraClientName]; !ok {
		env.ControllerClient.Headers[HasuraClientName] = clientName
	}

	return &env
}

// IsDebug check if env is in debug mode
func (e Environment) IsDebug() bool {
	return e.LogLevel == "debug"
}

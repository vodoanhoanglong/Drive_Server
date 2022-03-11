package main

import (
	goGql "github.com/hasura/go-graphql-client"
	"nexlab.tech/core/pkg/gql"
	"nexlab.tech/core/services/auth/env"
	"nexlab.tech/core/services/auth/utils"
)

type initConfig struct {
	env        *env.Environment
	controller *goGql.Client
	JwtAuth    *utils.JWTAuth
}

// NewInitConfig construct global initial configurations
func NewInitConfig(envVar *env.Environment) (*initConfig, error) {

	controllerClient := gql.NewClient(envVar.ControllerClient)
	jwtConfig := utils.NewJWTAuth(envVar.JWT, controllerClient)
	return &initConfig{
		env:        envVar,
		controller: controllerClient,
		JwtAuth:    jwtConfig,
	}, nil
}

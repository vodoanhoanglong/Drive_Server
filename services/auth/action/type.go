package action

import (
	"github.com/hasura/go-graphql-client"
	"github.com/hgiasac/hasura-router/go/action"
	"github.com/sirupsen/logrus"
	"nexlab.tech/core/pkg/access"
	"nexlab.tech/core/pkg/gql"
	"nexlab.tech/core/services/auth/env"
	"nexlab.tech/core/services/auth/utils"
)

// Config action handler config
type Config struct {
	Controller *graphql.Client
	Env        *env.Environment
	JwtAuth    *utils.JWTAuth
}

type actionContext struct {
	*action.Context
	Logger     *logrus.Entry
	Access     *access.Access
	Env        *env.Environment
	Controller *gql.AccessClient
	JwtAuth    *utils.JWTAuth
}

// wrap extends action context with new fields
func (ac Config) wrap(handler func(ctx *actionContext, rawBody []byte) (interface{}, error)) action.Action {
	return func(ctx *action.Context, rawBody []byte) (interface{}, error) {

		acs, err := access.ParseSessionVariables(ctx.SessionVariables)
		if err != nil {
			return nil, err
		}

		return handler(&actionContext{
			Context:    ctx,
			Logger:     logrus.NewEntry(logrus.New()),
			Access:     acs,
			Env:        ac.Env,
			JwtAuth:    ac.JwtAuth,
			Controller: gql.NewAccessClient(ac.Controller, acs),
		}, rawBody)
	}
}

// MessageOutput represent simple message response
type MessageOutput struct {
	Message string `json:"message"`
}

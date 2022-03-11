package event

import (
	"github.com/hasura/go-graphql-client"
	"github.com/hgiasac/hasura-router/go/event"
	"github.com/sirupsen/logrus"
	"nexlab.tech/core/pkg/access"
	"nexlab.tech/core/pkg/logging"
	"nexlab.tech/core/services/auth/env"
)

type Config struct {
	Env        *env.Environment
	Controller *graphql.Client
}

type Context struct {
	*event.Context
	Logger     *logrus.Entry
	Access     *access.Access
	Env        *env.Environment
	Controller *graphql.Client
}

type EventHandler func(ctx *Context, payload event.EventTriggerPayload) (interface{}, error)

// New create event handler instance
func New(config *Config) (*event.Router, error) {

	ctx := Context{
		Controller: config.Controller,
		Env:        config.Env,
	}
	events := event.New(map[string]event.Handler{
		eventExample: ctx.wrap(example),
	})
	events.OnSuccess(func(ctx *event.Context, response interface{}, metadata map[string]interface{}) {
		logging.LogSuccess(response, metadata)
	})

	events.OnError(func(ctx *event.Context, err error, metadata map[string]interface{}) {
		logging.LogError(err, metadata)
	})

	return events, nil
}

func (c Context) wrap(fn EventHandler) event.Handler {
	return func(ctx *event.Context, payload event.EventTriggerPayload) (interface{}, error) {
		acs, err := access.ParseSessionVariables(payload.Event.SessionVariables)
		if err != nil {
			return nil, err
		}

		return fn(&Context{
			Context:    ctx,
			Access:     acs,
			Logger:     logrus.NewEntry(logrus.New()),
			Controller: c.Controller,
			Env:        c.Env,
		}, payload)
	}
}

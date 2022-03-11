package event

import (
	"github.com/hgiasac/hasura-router/go/event"
)

const (
	eventExample = "example"
)

func example(ctx *Context, payload event.EventTriggerPayload) (interface{}, error) {
	return map[string]string{
		"message": "success",
	}, nil
}

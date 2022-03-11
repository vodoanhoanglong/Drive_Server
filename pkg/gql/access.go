package gql

import (
	"context"
	"encoding/json"

	"github.com/hasura/go-graphql-client"
	"nexlab.tech/core/pkg/access"
)

const headerKey = "x-headers"

func SetHeader(ctx context.Context, key, value string) context.Context {
	h := ctx.Value(headerKey)
	var headers map[string]string
	if h == nil {
		headers = map[string]string{}
	} else {
		headers = h.(map[string]string)
	}
	headers[key] = value
	return context.WithValue(ctx, headerKey, headers)
}

func SetHeaders(ctx context.Context, hs map[string]string) context.Context {
	h := ctx.Value(headerKey)
	var headers map[string]string
	if h == nil {
		headers = map[string]string{}
	} else {
		headers = h.(map[string]string)
	}
	for k, v := range hs {
		headers[k] = v
	}
	return context.WithValue(ctx, headerKey, headers)
}

// AccessClient is a graphql client which requests are made on behalf of an Actor
type AccessClient struct {
	*graphql.Client
	Access *access.Access
}

func NewAccessClient(cl *graphql.Client, ac *access.Access) *AccessClient {
	return &AccessClient{
		Client: cl,
		Access: ac,
	}
}

func (c *AccessClient) Query(ctx context.Context, q interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	ctx = SetHeaders(ctx, c.Access.ToHeaders())
	return c.Client.Query(ctx, q, variables)
}

func (c *AccessClient) NamedQuery(ctx context.Context, name string, q interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	ctx = SetHeaders(ctx, c.Access.ToHeaders())
	return c.Client.NamedQuery(ctx, name, q, variables)
}

func (c *AccessClient) NamedQueryRaw(ctx context.Context, name string, q interface{}, variables map[string]interface{}, options ...graphql.Option) (*json.RawMessage, error) {
	ctx = SetHeaders(ctx, c.Access.ToHeaders())
	return c.Client.NamedQueryRaw(ctx, name, q, variables)
}

func (c *AccessClient) Mutate(ctx context.Context, m interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	ctx = SetHeaders(ctx, c.Access.ToHeaders())
	return c.Client.Mutate(ctx, m, variables)
}

func (c *AccessClient) NamedMutate(ctx context.Context, name string, m interface{}, variables map[string]interface{}, options ...graphql.Option) error {
	ctx = SetHeaders(ctx, c.Access.ToHeaders())
	return c.Client.NamedMutate(ctx, name, m, variables)
}

func (c *AccessClient) NamedMutateRaw(ctx context.Context, name string, m interface{}, variables map[string]interface{}, options ...graphql.Option) (*json.RawMessage, error) {
	ctx = SetHeaders(ctx, c.Access.ToHeaders())
	return c.Client.NamedMutateRaw(ctx, name, m, variables)
}

// AsAdmin allows the client to act on behalf of a client
func (c *AccessClient) AsAdmin() *AccessClient {
	return NewAccessClient(c.Client, access.NewAccess(access.NewAdminActor()))
}

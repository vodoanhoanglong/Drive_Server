package test

import (
	"github.com/hasura/go-graphql-client"
	"nexlab.tech/core/pkg/gql"
)

func NewAdminControllerClient() *graphql.Client {
	return gql.NewClient(gql.ClientConfig{
		URL:         "http://localhost:8080/v1/graphql",
		AdminSecret: "hasura",
	})
}

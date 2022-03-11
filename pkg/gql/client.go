package gql

import (
	"fmt"
	"net/http"
	"time"

	gogql "github.com/hasura/go-graphql-client"
)

// AdminSecret header constants
const (
	AdminSecretHeader string = "X-Hasura-Admin-Secret"
	HasuraClientName         = "hasura-client-name"
)

// hasuraTransport transport for Hasura GraphQL Client
type hasuraTransport struct {
	adminSecret string
	headers     map[string]string
	// keep a reference to the client's original transport
	rt http.RoundTripper
}

// RoundTrip set header data before executing http request
func (t *hasuraTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.adminSecret != "" {
		r.Header.Set(AdminSecretHeader, t.adminSecret)
	}
	for k, v := range t.headers {
		r.Header.Set(k, v)
	}
	return t.rt.RoundTrip(r)
}

// ClientConfig input config for Client
type ClientConfig struct {
	BaseURL     string            `envconfig:"BASE_URL"`
	URL         string            `envconfig:"URL"`
	AdminSecret string            `envconfig:"ADMIN_SECRET" required:"true"`
	Headers     map[string]string `envconfig:"HEADERS"`
	Timeout     time.Duration     `envconfig:"TIMEOUT" default:"60s"`
}

// NewClient construct new client from environment
func NewClient(config ClientConfig) *gogql.Client {
	httpClient := &http.Client{
		Transport: &hasuraTransport{
			rt:          http.DefaultTransport,
			adminSecret: config.AdminSecret,
			headers:     config.Headers,
		},
		Timeout: config.Timeout,
	}

	url := config.URL
	if url == "" {
		url = fmt.Sprintf("%s/v1/graphql", config.BaseURL)
	}

	return gogql.NewClient(url, httpClient)
}

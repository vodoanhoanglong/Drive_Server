package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
	"github.com/hgiasac/hasura-router/go/tracing"
	"github.com/sirupsen/logrus"
	"nexlab.tech/core/pkg/access"
)

var (
	internalErrorMessage = []byte("{\"message\": \"internal error\"}")
)

type account_bool_exp map[string]interface{}

type rawAuthRequestBody struct {
	Headers map[string]string      `json:"headers"`
	Request map[string]interface{} `json:"request"`
}

type authRequestBody struct {
	APIKey    string
	AuthToken string
	UserAgent string
}

type authHandler struct {
	config *initConfig
}

func newAuthHandler(config *initConfig) *authHandler {
	return &authHandler{config}
}

func (ah *authHandler) Post(w http.ResponseWriter, r *http.Request) {
	debugLog, tracer := createDebugTracing(r)
	var rawBody rawAuthRequestBody

	if err := json.NewDecoder(r.Body).Decode(&rawBody); err != nil {
		onError(w, tracer, err)
		return
	}

	headers := stringMapToHeader(rawBody.Headers)
	body := authRequestBody{}
	if headers != nil {
		body.AuthToken = headers.Get("Authorization")
		body.UserAgent = headers.Get("User-Agent")
	}

	tracer = tracer.WithField("body", rawBody)
	resp, err := ah.authorizeUser(body, headers, debugLog)

	if err != nil {
		onError(w, tracer, err)
	} else {
		onSuccess(w, tracer, resp)
	}
}

// validate token
// get account, user, app_id and permissions
func (ah *authHandler) authorizeUser(data authRequestBody, headers http.Header, debugLog *logrus.Entry) (map[string]string, error) {

	if data.AuthToken == "" {

		return map[string]string{
			access.XHasuraRole:        string(access.RoleAnonymous),
			access.XHasuraCurrentTime: time.Now().UTC().Format(time.RFC3339),
		}, nil
	}

	token := strings.Split(data.AuthToken, " ")

	if len(token) != 2 {
		return nil, fmt.Errorf("invalid authorization header %s", data.AuthToken)
	}

	jwtPayload, err := ah.config.JwtAuth.DecodeToken(token[1])

	if err != nil {
		return nil, err
	}

	userId := string(jwtPayload.Subject)
	accountInfo, err := findAccoutById(userId, ah)

	if err != nil {
		return nil, err
	}

	return map[string]string{
		access.XHasuraUserID: userId,
		access.XHasuraRole:   accountInfo["role"],
	}, nil
}

func findAccoutById(userId string, ah *authHandler) (map[string]string, error) {
	var query struct {
		Accounts []struct {
			Email string `graphql:"email"`
			Role  string `graphql:"role"`
		} `graphql:"account(where: $where, limit: 1)"`
	}

	variables := map[string]interface{}{
		"where": account_bool_exp{
			"id": map[string]interface{}{
				"_ilike": userId,
			},
		},
	}

	err := ah.config.controller.Query(context.Background(), &query, variables, graphql.OperationName("GetAccountById"))

	if err != nil {
		return nil, err
	}

	result := query.Accounts[0]

	return map[string]string{
		"role":  result.Role,
		"email": result.Email,
	}, nil
}

func createDebugTracing(r *http.Request) (*logrus.Entry, *tracing.Tracing) {

	requestId := r.Header.Get("x-request-id")
	if requestId == "" {
		requestId = uuid.New().String()
	}
	fields := map[string]interface{}{
		"type":         "authorization",
		"request_id":   requestId,
		"method":       r.Method,
		"http_headers": r.Header,
	}
	debugLog := logrus.WithFields(fields)
	tracer := tracing.New(requestId).WithFields(fields)

	return debugLog, tracer
}

func onError(w http.ResponseWriter, tracer *tracing.Tracing, err error) {
	logrus.WithFields(tracer.Values()).WithError(err).Error(err)
	w.WriteHeader(401)

	resp, e := json.Marshal(map[string]interface{}{
		"message": err.Error(),
	})

	if e != nil {
		logrus.WithFields(tracer.Values()).WithError(err).Error("unauthorized response failed")
		w.Write(internalErrorMessage)
	} else {
		w.Write(resp)
	}
}

func onSuccess(w http.ResponseWriter, tracer *tracing.Tracing, respData map[string]string) {
	bResp, err := json.Marshal(respData)
	if err != nil {
		w.WriteHeader(500)
		w.Write(internalErrorMessage)
	} else {
		logrus.WithFields(tracer.Values()).WithField("data", respData).Info("successfully!")

		w.WriteHeader(200)
		w.Write(bResp)
	}
}

func stringMapToHeader(header map[string]string) http.Header {
	result := make(http.Header)
	for k, v := range header {
		if v != "" {
			result.Set(k, v)
		}
	}

	return result
}

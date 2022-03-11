package action

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/hasura/go-graphql-client"
	"nexlab.tech/core/pkg/util"
)

const (
	actionLogin        = "login"
	actionRefreshToken = "refreshToken"
)

func login(ctx *actionContext, payload []byte) (interface{}, error) {

	var input struct {
		Data struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		} `json:"data"`
	}

	err := json.Unmarshal([]byte(payload), &input)
	if err != nil {
		return nil, util.ErrBadRequest(err)
	}

	var query struct {
		Accounts []struct {
			ID       string `graphql:"id"`
			Password string `graphql:"password"`
			Role     string `graphql:"role"`
		} `graphql:"account(where: $where, limit: 1)"`
	}

	variables := map[string]interface{}{
		"where": account_bool_exp{
			"email": map[string]interface{}{
				"_like": input.Data.Email,
			},
		},
	}

	err = ctx.Controller.Query(context.Background(), &query, variables, graphql.OperationName("GetAccountByEmail"))

	if err != nil {
		return nil, err
	}

	if len(query.Accounts) == 0 {
		return nil, errors.New("account not found")
	}

	account := query.Accounts[0]
	if account.Password == "" || ctx.JwtAuth.ComparePassword(account.Password, input.Data.Password) != nil {
		return nil, errors.New("password not match")
	}

	token, err := ctx.JwtAuth.EncodeToken(account.ID)

	if err != nil {
		return nil, err
	}
	return token, nil
}

func refreshToken(ctx *actionContext, payload []byte) (interface{}, error) {

	var input struct {
		Data struct {
			AccessToken  string `json:"access_token"`
			RefreshToken string `json:"refresh_token"`
		} `json:"data"`
	}

	err := json.Unmarshal([]byte(payload), &input)
	if err != nil {
		return nil, util.ErrBadRequest(err)
	}

	if input.Data.RefreshToken == "" {
		return nil, errors.New("refresh token required")
	}

	if input.Data.AccessToken == "" {
		return nil, errors.New("access token required")
	}

	token, err := ctx.JwtAuth.RefreshToken(input.Data.RefreshToken, input.Data.AccessToken)

	if err != nil {
		return nil, err
	}

	return token, nil
}

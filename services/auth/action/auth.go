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

	// enum login type
	firebase = "firebase"
	facebook = "facebook"
	google   = "google"
)

func login(ctx *actionContext, payload []byte) (interface{}, error) {

	var input struct {
		Data struct {
			Email     string `json:"email"`
			Password  string `json:"password"`
			LoginType string `json:"loginType"`
			FullName  string `json:"fullName"`
			Avatar    string `json:"avatar"`
		} `json:"data"`
	}

	err := json.Unmarshal([]byte(payload), &input)
	if err != nil {
		return nil, util.ErrBadRequest(err)
	}

	if input.Data.LoginType != firebase {
		tokenThirdParty, errThirdParty := findOrCreateLoginThirdParty(ctx, input.Data.Email, input.Data.FullName, input.Data.Avatar, input.Data.LoginType)

		if errThirdParty != nil {
			return nil, errThirdParty
		}

		return tokenThirdParty, nil
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
			"loginType": map[string]interface{}{
				"_eq": input.Data.LoginType,
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

func findOrCreateLoginThirdParty(ctx *actionContext, email string, fullName string, avatar string, loginType string) (interface{}, error) {
	var query struct {
		AccountByEmail []struct {
			ID string `graphql:"id"`
		} `graphql:"account(where: $where, limit: 1)"`
	}

	variables := map[string]interface{}{
		"where": account_bool_exp{
			"email": map[string]interface{}{
				"_like": email,
			},
			"loginType": map[string]interface{}{
				"_neq": firebase,
			},
		},
	}

	err := ctx.Controller.Query(context.Background(), &query, variables)

	if err != nil {
		return nil, err
	}

	if len(query.AccountByEmail) == 0 {
		var mutation struct {
			CreateAccount struct {
				ID       string `graphql:"id"`
				Email    string `graphql:"email"`
				FullName string `graphql:"fullName"`
				Role     string `graphql:"role"`
			} `graphql:"insert_account_one(object: $object)"`
		}

		mutationVariables := map[string]interface{}{
			"object": account_insert_input{
				"email":      email,
				"role":       "user",
				"fullName":   fullName,
				"avatar_url": avatar,
				"loginType":  loginType,
			},
		}

		err = ctx.Controller.Mutate(context.Background(), &mutation, mutationVariables)

		if err != nil {
			return nil, util.ErrBadRequest(err)
		}

		tokenCreate, err := ctx.JwtAuth.EncodeToken(mutation.CreateAccount.ID)

		if err != nil {
			return nil, err
		}

		return tokenCreate, nil
	}

	token, err := ctx.JwtAuth.EncodeToken(query.AccountByEmail[0].ID)

	if err != nil {
		return nil, err
	}

	return token, nil
}

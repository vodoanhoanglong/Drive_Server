package action

import (
	"context"
	"encoding/json"
	"errors"
	"net/smtp"

	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
	"github.com/hgiasac/hasura-router/go/types"
	"nexlab.tech/core/pkg/util"
)

const (
	actionAdminChangePassword = "changeAccountPassword"
	actionForgotPassword      = "forgotPassword"
)

// adminChangePassword change account password by admin
func changeAccountPassword(ctx *actionContext, payload []byte) (interface{}, error) {
	var input struct {
		Data struct {
			AccountID   string `json:"account_id"`
			NewPassword string `json:"new_password"`
		} `json:"data"`
	}

	err := json.Unmarshal([]byte(payload), &input)
	if err != nil {
		return nil, util.ErrBadRequest(err)
	}

	if input.Data.AccountID == "" {
		return nil, types.NewError("required:account_id", "account_id is required")
	}

	passwordHashed, err := ctx.JwtAuth.EncryptPassword(input.Data.NewPassword)
	if err != nil {
		return nil, util.ErrBadRequest(err)
	}

	randomUUID := uuid.New().String()
	randomHashed, err := ctx.JwtAuth.EncryptPassword(randomUUID)

	if err != nil {
		return nil, util.ErrBadRequest(err)
	}

	var query struct {
		UpdatePassword struct {
			ID       string `graphql:"id"`
			Password string `graphql:"password"`
		} `graphql:"update_account_by_pk(pk_columns: $pk_columns, _set: $set)"`
	}

	variables := map[string]interface{}{
		"pk_columns": account_pk_columns_input{
			"id": input.Data.AccountID,
		},
		"set": account_set_input{
			"password":   string(passwordHashed),
			"randomHash": string(randomHashed),
		},
	}

	err = ctx.Controller.Mutate(context.Background(), &query, variables)

	if err != nil {
		return nil, util.ErrBadRequest(err)
	}
	return map[string]string{
		"message":  "success",
		"id":       query.UpdatePassword.ID,
		"password": query.UpdatePassword.Password,
	}, nil

}

func forgotPassword(ctx *actionContext, payload []byte) (interface{}, error) {

	var input struct {
		Data struct {
			Email string `json:"email"`
		} `json:"data"`
	}

	err := json.Unmarshal([]byte(payload), &input)
	if err != nil {
		return nil, util.ErrBadRequest(err)
	}

	if input.Data.Email == "" {
		return nil, types.NewError("required:email", "email is required")
	}

	var query struct {
		Accounts []struct {
			ID string `graphql:"id"`
		} `graphql:"account(where: $where, limit: 1)"`
	}

	variables := map[string]interface{}{
		"where": account_bool_exp{
			"email": map[string]interface{}{
				"_ilike": input.Data.Email,
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

	accountID := query.Accounts[0].ID

	token, err := ctx.JwtAuth.EncodeToken(accountID)

	if err != nil {
		return nil, err
	}

	from := ctx.Env.Email
	password := ctx.Env.Password

	toList := []string{input.Data.Email}

	host := "smtp.gmail.com"
	port := "587"

	msg := "Reset Password!!! \n Click link: your link" // Add redirect URL to password reset page

	body := []byte(msg)

	auth := smtp.PlainAuth("", from, password, host)
	err = smtp.SendMail(host+":"+port, auth, from, toList, body)

	if err != nil {
		return nil, err
	}

	return map[string]string{
		"message": "success",
		"email":   input.Data.Email,
		"id":      accountID,
		"token":   token.AccessToken,
	}, nil
}

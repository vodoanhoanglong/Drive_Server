package action

import (
	"github.com/hgiasac/hasura-router/go/action"
	"nexlab.tech/core/pkg/logging"
)

// New create action router instance
func New(hc Config) (*action.Router, error) {

	actions, err := action.New(map[action.ActionName]action.Action{
		actionCreateAccount:       hc.wrap(createAccount),
		actionAdminChangePassword: hc.wrap(changeAccountPassword),
		actionLogin:               hc.wrap(login),
		actionRefreshToken:        hc.wrap(refreshToken),
		actionForgotPassword:      hc.wrap(forgotPassword),
		actionUploadFile:          hc.wrap(uploadFile),
		actionMoveFile:            hc.wrap(moveFile),
		actionUpdateFile:          hc.wrap(updateFile),
		actionShareFile:           hc.wrap(shareFile),
	})

	if err != nil {
		return nil, err
	}

	actions.OnSuccess(func(ctx *action.Context, response interface{}, metadata map[string]interface{}) {
		logging.LogSuccess(response, metadata)
	})

	actions.OnError(func(ctx *action.Context, err error, metadata map[string]interface{}) {
		logging.LogError(err, metadata)
	})

	return actions, nil
}

package util

import "github.com/hgiasac/hasura-router/go/types"

const (
	ErrCodePermissionDenied = "permission_denied"
)

func NewError(code string, message string) types.Error {
	return types.NewError(code, message)
}

func ErrUnknown(err error) types.Error {
	return types.NewError(types.ErrCodeUnknown, err.Error())
}

func ErrUnauthorized(err error) types.Error {
	return types.NewError(types.ErrCodeUnauthorized, err.Error())
}

func ErrBadRequest(err error) types.Error {
	return types.NewError(types.ErrCodeBadRequest, err.Error())
}

func ErrInternal(err error) types.Error {
	return types.NewError(types.ErrCodeInternal, err.Error())
}

func ErrPermissionDenied(err error) types.Error {
	message := "permission denied"
	if err != nil {
		message = err.Error()
	}
	return types.NewError(ErrCodePermissionDenied, message)
}

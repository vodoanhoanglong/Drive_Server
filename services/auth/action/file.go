package action

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
	"nexlab.tech/core/pkg/util"
)

const (
	actionUploadFile = "uploadFile"
	actionMoveFile   = "moveFile"
	actionUpdateFile = "updateFile"
)

type files_insert_input map[string]interface{}
type move_file_args map[string]interface{}
type files_pk_columns_input map[string]interface{}
type files_set_input map[string]interface{}
type check_file_name_args map[string]interface{}
type users_bool_exp map[string]interface{}

type UploadFileInput struct {
	Name      string `json:"name"`
	Path      string `json:"path"`
	Size      int    `json:"size"`
	Url       string `json:"url"`
	Extension string `json:"extension"`
}

type MoveFileInput struct {
	ToPath        string `json:"to_path"`
	ToExtension   string `json:"to_extension"`
	FromPath      string `json:"from_path"`
	FromExtension string `json:"from_extension"`
	FromName      string `json:"from_name"`
}

type UpdateFileInput struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Path      string `json:"path"`
	Extension string `json:"extension"`
}

func uploadFile(ctx *actionContext, payload []byte) (interface{}, error) {
	var appInput struct {
		Data UploadFileInput `json:"data"`
	}

	err := json.Unmarshal([]byte(payload), &appInput)
	if err != nil {
		return nil, util.ErrBadRequest(err)
	}

	input := appInput.Data

	departmentId := ""
	checkPath := ""
	if input.Path == "department" {
		departmentId, err = getDepartmentId(ctx)
		if err != nil {
			return nil, err
		}
		checkPath = input.Path + "/" + ctx.Access.UserID + "/" + departmentId
	} else if input.Path == "personal" || input.Path == "general" {
		checkPath = input.Path + "/" + ctx.Access.UserID
	} else {
		checkPath = input.Path
	}

	if ok, _ := checkFileName(ctx, checkPath, input.Name, input.Extension); ok {
		return nil, errors.New("filename already exists")
	}

	randomUUID := uuid.New().String()
	path := checkPath + "/" + randomUUID

	var query struct {
		UploadFile struct {
			ID   string `graphql:"id"`
			Name string `graphql:"name"`
			Path string `graphql:"path"`
			Url  string `graphql:"url"`
		} `graphql:"insert_files_one(object: $object)"`
	}

	variables := map[string]interface{}{
		"object": files_insert_input{
			"id":        randomUUID,
			"name":      input.Name,
			"path":      path,
			"url":       input.Url,
			"size":      input.Size,
			"extension": input.Extension,
			"userId":    ctx.Access.UserID,
		},
	}

	err = ctx.Controller.Mutate(context.Background(), &query, variables)

	if err != nil {
		return nil, util.ErrBadRequest(err)
	}

	results := query.UploadFile

	return map[string]string{
		"id":   results.ID,
		"name": results.Name,
		"path": results.Path,
		"url":  results.Url,
	}, nil
}

func moveFile(ctx *actionContext, payload []byte) (interface{}, error) {
	var appInput struct {
		Data MoveFileInput `json:"data"`
	}

	err := json.Unmarshal([]byte(payload), &appInput)
	if err != nil {
		return nil, util.ErrBadRequest(err)
	}

	input := appInput.Data

	if input.ToExtension != "folder" {
		return nil, errors.New("destination path must be a directory")
	}

	toPath := strings.Split(input.ToPath, "/")
	fromPath := strings.Split(input.FromPath, "/")

	toPathPrefix := toPath[0]
	fromPathPrefix := fromPath[0]

	if toPathPrefix != fromPathPrefix {
		return nil, errors.New("toPath and fromPath must be in the same type")
	}

	if ok, err := checkPermission(ctx, input.ToPath); !ok {
		return nil, err
	}

	if ok, err := checkPermission(ctx, input.FromPath); !ok {
		return nil, err
	}

	if ok, _ := checkFileName(ctx, input.ToPath, input.FromName, input.FromExtension); ok {
		return nil, errors.New("filename already exists")
	}

	var query struct {
		MoveFile []struct {
			ID   string `graphql:"id"`
			Name string `graphql:"name"`
			Path string `graphql:"path"`
			Url  string `graphql:"url"`
		} `graphql:"move_file(args: $args)"`
	}

	variables := map[string]interface{}{
		"args": move_file_args{
			"from_path": input.FromPath,
			"to_path":   input.ToPath,
		},
	}

	err = ctx.Controller.Mutate(context.Background(), &query, variables)

	if err != nil {
		return nil, util.ErrBadRequest(err)
	}

	return map[string]string{
		"message": "success",
	}, nil
}

func updateFile(ctx *actionContext, payload []byte) (interface{}, error) {
	var appInput struct {
		Data UpdateFileInput `json:"data"`
	}

	err := json.Unmarshal([]byte(payload), &appInput)
	if err != nil {
		return nil, util.ErrBadRequest(err)
	}

	input := appInput.Data

	if ok, err := checkPermission(ctx, input.Path); !ok {
		return nil, err
	}

	parentPath := strings.Split(input.Path, "/")
	parentPath = parentPath[:len(parentPath)-1]
	joinPath := strings.Join(parentPath, "/")

	if ok, _ := checkFileName(ctx, joinPath, input.Name, input.Extension); ok {
		return nil, errors.New("filename already exists")
	}

	var query struct {
		UpdateFile struct {
			ID   string `graphql:"id"`
			Name string `graphql:"name"`
		} `graphql:"update_files_by_pk(pk_columns: $pk_columns, _set: $set)"`
	}

	variables := map[string]interface{}{
		"pk_columns": files_pk_columns_input{
			"id": input.ID,
		},
		"set": files_set_input{
			"name": input.Name,
		},
	}

	err = ctx.Controller.Mutate(context.Background(), &query, variables)

	if err != nil {
		return nil, util.ErrBadRequest(err)
	}

	results := query.UpdateFile

	return map[string]string{
		"id":   results.ID,
		"name": results.Name,
	}, nil
}

func checkFileName(ctx *actionContext, path string, name string, extension string) (bool, error) {
	var query struct {
		Files []struct {
			Name string `graphql:"name"`
		} `graphql:"check_file_name(args: $args)"`
	}

	variables := map[string]interface{}{
		"args": check_file_name_args{
			"path_input":      path,
			"name_input":      name,
			"extension_input": extension,
		},
	}

	err := ctx.Controller.Query(context.Background(), &query, variables, graphql.OperationName("CheckFileName"))

	if err != nil {
		return true, err
	}

	if len(query.Files) > 0 {
		return true, nil
	}

	return false, nil
}

func checkPermission(ctx *actionContext, path string) (bool, error) {
	isOwner := true
	userId := strings.Split(path, "/")[1]
	errors := errors.New("you don't have permission to access this file")

	if ctx.Access.UserID != userId {
		isOwner = false
	}

	if strings.HasPrefix(path, "department") {
		if ctx.Access.Role != "moderator" && !isOwner {
			return false, errors
		}
	} else {
		if !isOwner {
			return false, errors
		}
	}

	return true, nil
}

func getDepartmentId(ctx *actionContext) (string, error) {

	var query struct {
		Users []struct {
			DepartmentId string `graphql:"departmentId"`
		} `graphql:"users(where: $where, limit: 1)"`
	}

	variables := map[string]interface{}{
		"where": users_bool_exp{
			"id": map[string]interface{}{
				"_eq": ctx.Access.UserID,
			},
		},
	}

	err := ctx.Controller.Query(context.Background(), &query, variables, graphql.OperationName("GetDepartmentId"))

	if err != nil {
		return "", err
	}

	if len(query.Users) == 0 {
		return "", errors.New("departmentId not found")
	}

	return query.Users[0].DepartmentId, nil
}

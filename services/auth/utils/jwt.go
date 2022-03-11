package utils

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	jose "github.com/dvsekhvalnov/jose2go"
	"github.com/google/uuid"
	"github.com/hasura/go-graphql-client"
	"golang.org/x/crypto/bcrypt"
	"nexlab.tech/core/pkg/util"
)

type account_bool_exp map[string]interface{}
type account_set_input map[string]interface{}
type account_pk_columns_input map[string]interface{}

type jwtPayload struct {
	Issuer         string `json:"iss"`
	Subject        string `json:"sub"`
	Audience       string `json:"aud"`
	ExpirationTime int64  `json:"exp"`
	NotBeforeTime  int64  `json:"nbt"`
	IssuedAt       int64  `json:"iat"`
	JwtID          string `json:"jti"`
	RandomHash     string `json:"rdh"`
}

type AccessToken struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token,omitempty"`
	Scope        string `json:"scope,omitempty"`
}

type JWTAuthConfig struct {
	Cost       int           `envconfig:"JWT_HASH_COST" default:"10"`
	SessionKey string        `envconfig:"SESSION_KEY"`
	TTL        time.Duration `envconfig:"SESSION_TTL" default:"1h"`
	RefreshTTL time.Duration `envconfig:"SESSION_REFRESH_TTL" default:"0ms"`
	Issuer     string        `envconfig:"JWT_ISSUER"`
	Algorithm  string        `envconfig:"JWT_ALGORITHM" default:"HS256"`
}

func (jac JWTAuthConfig) Validate() error {
	if jac.SessionKey == "" {
		return errors.New("SESSION_KEY is required")
	}
	if jac.Issuer == "" {
		return errors.New("JWT_ISSUER is required")
	}

	return nil
}

type JWTAuth struct {
	config     JWTAuthConfig
	controller *graphql.Client
}

func NewJWTAuth(config JWTAuthConfig, controller *graphql.Client) *JWTAuth {
	if config.Cost == 0 {
		config.Cost = bcrypt.DefaultCost
	}
	if config.Algorithm == "" {
		config.Algorithm = jose.HS256
	}

	return &JWTAuth{
		config:     config,
		controller: controller,
	}
}

func (ja *JWTAuth) EncryptPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), ja.config.Cost)
}

func (ja *JWTAuth) EncodeToken(uid string) (*AccessToken, error) {

	var query struct {
		Accounts []struct {
			ID         string `graphql:"id"`
			RandomHash string `graphql:"randomHash"`
		} `graphql:"account(where: $where, limit: 1)"`
	}

	variables := map[string]interface{}{
		"where": account_bool_exp{
			"id": map[string]interface{}{
				"_eq": uid,
			},
		},
	}

	err := ja.controller.Query(context.Background(), &query, variables)

	if err != nil {
		return nil, err
	}

	randomHash := query.Accounts[0].RandomHash

	now := time.Now()
	exp := now.Add(ja.config.TTL)
	jwtID := uuid.New().String()
	payload := jwtPayload{
		JwtID:          jwtID,
		Issuer:         ja.config.Issuer,
		Subject:        uid,
		Audience:       "access",
		IssuedAt:       now.Unix(),
		NotBeforeTime:  now.Unix(),
		ExpirationTime: exp.Unix(),
		RandomHash:     randomHash,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	token, err := jose.SignBytes(payloadBytes, ja.config.Algorithm, []byte(ja.config.SessionKey),
		jose.Header("typ", "JWT"),
		jose.Header("alg", ja.config.Algorithm),
	)

	if err != nil {
		return nil, err
	}

	// encode refresh token if the expiry is set
	var refreshToken string
	if ja.config.RefreshTTL >= ja.config.TTL {
		refreshPayload := jwtPayload{
			JwtID:          ja.genRefreshTokenID(jwtID),
			Issuer:         ja.config.Issuer,
			Subject:        uid,
			Audience:       "refresh",
			IssuedAt:       now.Unix(),
			NotBeforeTime:  now.Unix(),
			ExpirationTime: now.Add(ja.config.RefreshTTL).Unix(),
			RandomHash:     randomHash,
		}

		refreshPayloadBytes, err := json.Marshal(refreshPayload)
		if err != nil {
			return nil, err
		}

		refreshToken, err = jose.SignBytes(refreshPayloadBytes, ja.config.Algorithm, []byte(ja.config.SessionKey),
			jose.Header("typ", "JWT"),
			jose.Header("alg", ja.config.Algorithm),
		)

		if err != nil {
			return nil, err
		}
	}

	return &AccessToken{
		AccessToken:  token,
		TokenType:    "jwt",
		ExpiresIn:    int(ja.config.TTL / time.Second),
		RefreshToken: refreshToken,
	}, nil
}

func (ja *JWTAuth) ComparePassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}

func (ja *JWTAuth) RefreshToken(refreshToken string, accessToken string) (*AccessToken, error) {
	decodedRefreshToken, err := ja.DecodeToken(refreshToken)
	if err != nil {
		return nil, err
	}

	if decodedRefreshToken.Audience != "refresh" {
		return nil, errors.New("token_mismatch")
	}

	decodedToken, err := ja.DecodeToken(accessToken)
	if err != nil && err.Error() != "token_expired" {
		return nil, err
	}

	if decodedRefreshToken.JwtID != ja.genRefreshTokenID(decodedToken.JwtID) ||
		decodedRefreshToken.Subject != decodedToken.Subject ||
		decodedRefreshToken.IssuedAt != decodedToken.IssuedAt {
		return nil, errors.New("token_mismatch")
	}

	// revoke old token
	randomUUID := uuid.New().String()
	randomHashed, err := ja.EncryptPassword(randomUUID)

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
			"id": decodedRefreshToken.Subject,
		},
		"set": account_set_input{
			"randomHash": string(randomHashed),
		},
	}

	err = ja.controller.Mutate(context.Background(), &query, variables)

	if err != nil {
		return nil, util.ErrBadRequest(err)
	}

	return ja.EncodeToken(decodedRefreshToken.Subject)
}

func (ja *JWTAuth) DecodeToken(token string) (*jwtPayload, error) {

	bytes, _, err := jose.DecodeBytes(token, []byte(ja.config.SessionKey))
	if err != nil {
		return nil, err
	}

	var result jwtPayload

	err = json.Unmarshal(bytes, &result)
	if err != nil {
		return nil, err
	}

	if ja.config.Issuer != "" && ja.config.Issuer != result.Issuer {
		return &result, errors.New("jwt_invalid_issuer")
	}

	if result.ExpirationTime <= time.Now().Unix() || !ja.checkRandomHash(result.Subject, result.RandomHash) {
		return &result, errors.New("token_expired")
	}

	return &result, nil
}

func (ja *JWTAuth) checkRandomHash(accountId string, randomHash string) bool {
	var query struct {
		Accounts []struct {
			ID         string `graphql:"id"`
			RandomHash string `graphql:"randomHash"`
		} `graphql:"account(where: $where, limit: 1)"`
	}

	variables := map[string]interface{}{
		"where": account_bool_exp{
			"id": map[string]interface{}{
				"_eq": accountId,
			},
			"randomHash": map[string]interface{}{
				"_eq": randomHash,
			},
		},
	}

	err := ja.controller.Query(context.Background(), &query, variables)

	if err != nil {
		return false
	}

	return len(query.Accounts) > 0
}

func (ja *JWTAuth) genRefreshTokenID(id string) string {
	return fmt.Sprintf("%s-refresh", id)
}

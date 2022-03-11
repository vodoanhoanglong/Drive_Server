package access

import (
	"errors"
	"fmt"
)

type Role string

// access constants
const (
	XHasuraRole             = "x-hasura-role"
	XHasuraUserEmail        = "x-hasura-user-email"
	XHasuraUserID           = "x-hasura-user-id"
	XHasuraCurrentTime      = "x-hasura-current-time"
	RoleAnonymous      Role = "anonymous"
	RoleAdmin          Role = "admin"
	RoleModerator      Role = "moderator"
	RoleUser           Role = "user"
)

var (
	errUnauthorized = errors.New("unauthorized")
)

func GetRoles() []string {
	return []string{
		string(RoleAnonymous),
		string(RoleAdmin),
		string(RoleModerator),
	}
}

func ParseRole(s string) (Role, error) {
	for _, v := range GetRoles() {
		if v == s {
			return Role(v), nil
		}
	}

	return Role(""), fmt.Errorf("invalid role: %s", s)
}

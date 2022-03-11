package access

import (
	"strings"

	"nexlab.tech/core/pkg/util"
)

// Actor denotes who is making a request to our APIs
type Actor struct {
	UserID string
	Role   Role
}

func NewAdminActor() *Actor {
	return &Actor{
		Role: RoleAdmin,
	}
}

func NewUserActor(userID string, role Role) *Actor {
	return &Actor{
		Role:   role,
		UserID: userID,
	}
}

func (a *Actor) HasRole(role ...string) bool {
	for _, r := range role {
		if strings.EqualFold(string(a.Role), r) {
			return true
		}
	}
	return false
}

func (a *Actor) IsAdmin() bool {
	return a.HasRole(string(RoleAdmin))
}

func (a *Actor) IsAnonymous() bool {
	return a.HasRole(string(RoleAnonymous))
}

// Access holds information about what a given actor can access.
// code using Access values will determine whether a user can access a given resource based on the actor's role, identity or allowed lists
type Access struct {
	*Actor
}

// NewAccess constructs new Access instance
func NewAccess(actor *Actor) *Access {
	return &Access{Actor: actor}
}

func ParseSessionVariables(sv map[string]string) (*Access, error) {

	role, ok := sv[XHasuraRole]
	if !ok {
		return nil, errUnauthorized
	}

	userID, ok := sv[XHasuraUserID]
	if !ok && !util.HasString(GetRoles(), role) {
		return nil, errUnauthorized
	}

	actor := &Actor{
		UserID: userID,
		Role:   Role(role),
	}

	return &Access{
		Actor: actor,
	}, nil
}

func (acc *Access) ToHeaders() map[string]string {

	headers := make(map[string]string)
	if acc.IsAdmin() {
		headers[XHasuraRole] = string(RoleAdmin)
	} else {
		headers[XHasuraRole] = string(acc.Role)
		headers[XHasuraUserID] = acc.UserID
	}
	return headers
}

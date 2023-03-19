package core

import "context"

const KeyRequester = "requester"

type TokenPayload struct {
	UId int `json:"user_id"`
}

func (p TokenPayload) UserId() int {
	return p.UId
}

type Requester interface {
	GetUserId() int
	GetRole() string
	GetStatus() string
}

func IsAdmin(requester Requester) bool {
	return requester.GetRole() == "admin" || requester.GetRole() == "mod"
}

func GetRequester(ctx context.Context) Requester {
	if requester, ok := ctx.Value(KeyRequester).(Requester); ok {
		return requester
	}

	return nil
}

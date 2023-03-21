package core

import "context"

const KeyRequester = "requester"

type Requester interface {
	GetUserId() string
	GetTokenId() string
}

type requesterData struct {
	UserId string `json:"user_id"`
	Tid    string `json:"tid"`
}

func (r *requesterData) GetUserId() string {
	return r.UserId
}

func (r *requesterData) GetTokenId() string {
	return r.Tid
}

func GetRequester(ctx context.Context) Requester {
	if requester, ok := ctx.Value(KeyRequester).(Requester); ok {
		return requester
	}

	return nil
}

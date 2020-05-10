package power

import "github.com/ftCommunity/roboheart/package/api"

type wakeAlarmRequest struct {
	api.TokenRequest
	Time int64 `json:"time"`
}

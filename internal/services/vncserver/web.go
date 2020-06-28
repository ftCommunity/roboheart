package vncserver

import "github.com/ftCommunity/roboheart/package/api"

type state struct {
	State bool `json:"state"`
}

type stateSet struct {
	api.TokenRequest
	state
}

type autostartState struct {
	Autostart bool `json:"autostart"`
}

type autostartSet struct {
	api.TokenRequest
	autostartState
}

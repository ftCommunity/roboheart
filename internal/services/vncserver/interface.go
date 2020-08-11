package vncserver

type VNCServer interface {
	StartVNC(token string) (error, bool)
	StopVNC(token string) (error, bool)
	GetAutostart() bool
	SetAutostart(token string, autostart bool) (error, bool)
}

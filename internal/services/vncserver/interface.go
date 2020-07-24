package vncserver

type VNCServer interface {
	Start(token string) (error, bool)
	Stop(token string) (error, bool)
	GetAutostart() bool
	SetAutostart(token string, autostart bool) (error, bool)
}

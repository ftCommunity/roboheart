package instance

type LoggerFunc func(...interface{})
type ErrorFunc func(...interface{})
type SelfKillFunc func()

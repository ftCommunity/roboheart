package progress

import "time"

type Progress struct {
	Name      *string
	State     *float32
	Remaining *time.Duration
	Steps     *[]Progress
}

type ProgressCallback func(Progress)

type ProgressConf struct {
	Callback ProgressCallback
	Interval time.Duration
}

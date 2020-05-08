package progress

import "time"

type Progress struct {
	name     string
	state    *float64
	remaing  *time.Duration
	steps    *[]*Progress
	callback func()
}

func (p *Progress) SetState(s float64) {
	p.state = &s
	p.Notify()
}

func (p *Progress) UnsetState() {
	p.state = nil
	p.Notify()
}

func (p *Progress) SetRemaining(r time.Duration) {
	p.remaing = &r
	p.Notify()
}

func (p *Progress) UnsetRemaining() {
	p.remaing = nil
	p.Notify()
}

func (p *Progress) RegisterSteps(steps []string) {
	for _, s := range steps {
		*p.steps = append(*p.steps, NewProgress(s, []string{}, func() { p.callback() }))
	}
	p.Notify()
}

func (p *Progress) GetStep(i int) *Progress {
	return (*p.steps)[i]
}

func (p *Progress) Notify() { p.callback() }

func NewProgress(name string, steps []string, cb func()) *Progress {
	p := new(Progress)
	p.name = name
	p.callback = cb
	*p.steps = make([]*Progress, 0)
	for _, s := range steps {
		*p.steps = append(*p.steps, NewProgress(s, []string{}, func() { p.Notify() }))
	}
	return p
}

type ProgressConf struct {
	Callback func(*float64, *time.Duration)
	Interval time.Duration
}

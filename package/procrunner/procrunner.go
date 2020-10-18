package procrunner

import (
	"errors"
	"io"
	"log"
	"os/exec"
	"sync"
	"syscall"
)

type Callback func(int)

var (
	errAlreadyRunning = errors.New("process already running")
	errNoProcess      = errors.New("no process to stop")
)

type ProcRunner struct {
	name          string
	args          []string
	proc          *exec.Cmd
	proclock      sync.Mutex
	onAutoRestart *Callback
	onEnd         *Callback
	stdin         io.WriteCloser
	stdout        io.ReadCloser
	stderr        io.ReadCloser
	killing       bool
}

func (p *ProcRunner) Start() error {
	p.proclock.Lock()
	defer p.proclock.Unlock()
	if p.proc != nil {
		return errAlreadyRunning
	}
	p.proc = exec.Command(p.name, p.args...)
	var err error
	p.stdin, err = p.proc.StdinPipe()
	if err != nil {
		return err
	}
	p.stdout, err = p.proc.StdoutPipe()
	if err != nil {
		return err
	}
	p.stderr, err = p.proc.StderrPipe()
	if err != nil {
		return err
	}
	err = p.proc.Start()
	if err != nil {
		return err
	}
	go p.handleEnd()
	return nil
}

func (p *ProcRunner) Stop() error {
	p.proclock.Lock()
	defer p.proclock.Unlock()
	if p.proc == nil {
		return errNoProcess
	}
	p.killing = true
	log.Println("PROC KILLING SET", p.name)
	_ = p.proc.Process.Kill()
	return nil
}

func (p *ProcRunner) SetOnAutoRestartCallback(c Callback) {
	p.onAutoRestart = &c
}

func (p *ProcRunner) UnsetOnAutoRestartCallback() {
	p.onAutoRestart = nil
}

func (p *ProcRunner) SetOnEndCallback(c Callback) {
	p.onEnd = &c
}

func (p *ProcRunner) UnsetEndCallback() {
	p.onEnd = nil
}

func (p *ProcRunner) handleEnd() {
	code := 0
	if err := p.proc.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				code = status.ExitStatus()
			}
		}
	}
	if p.killing {
		p.proc = nil
		p.killing = false
		return
	}
	p.proclock.Lock()
	if p.proc == nil {
		p.proclock.Unlock()
		return
	}
	p.proc = nil
	p.proclock.Unlock()
	if p.onAutoRestart != nil {
		(*p.onAutoRestart)(code)
		go func() { _ = p.Start() }()
		return
	}
	if p.onEnd != nil {
		(*p.onEnd)(code)
	}
}

func NewProcRunner(name string, args ...string) *ProcRunner {
	p := new(ProcRunner)
	p.name = name
	p.args = args
	return p
}

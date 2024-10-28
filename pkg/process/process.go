// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package process

import (
	"os"
	"os/exec"
	"syscall"
	"time"

	"github.com/go-logr/logr"
	"k8s.io/klog/v2"
)

// Process represents the child process command environment
type Process struct {
	exec.Cmd
	log logr.Logger
}

// New creates a new process
func New(log logr.Logger, cmdLine string, args ...string) *Process {
	cmd := exec.Command(cmdLine, args...)
	cmd.Env = os.Environ()
	cmd.Stderr = os.Stderr
	cmd.Stdout = os.Stdout
	cmd.WaitDelay = time.Second * 10

	return &Process{
		Cmd: *cmd,
		log: log,
	}
}

// Start starts the child process
func (p *Process) Start() error {

	if err := p.Cmd.Start(); err != nil {
		return err
	}

	go func() {
		if err := p.Cmd.Wait(); err != nil {
			if exitErr, ok := err.(*exec.ExitError); ok {
				klog.NewKlogr().Error(err, "process exited with an error", "exitCode", exitErr.ExitCode())
			}
		}
	}()

	p.log.Info("started", "command", p.Cmd.String(), "pid", p.Cmd.Process.Pid)
	return nil
}

// Stop stops the child process
func (p *Process) Stop() error {

	p.log.Info("sending SIGINT signal", "pid", p.Cmd.Process.Pid)

	if err := p.Cmd.Process.Signal(syscall.SIGINT); err != nil {
		p.log.Error(err, "error sending SIGINT signal, try to kill", "pid", p.Cmd.Process.Pid)
		return p.Cmd.Process.Kill()
	}

	return nil
}

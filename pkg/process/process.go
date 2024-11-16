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
)

const terminationTimeout = 10 * time.Second

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

// Start the child process
func (p *Process) Start() error {

	if err := p.Cmd.Start(); err != nil {
		return err
	}

	p.log.Info("Start", "process", p.Cmd.String(), "pid", p.Cmd.Process.Pid)
	return nil
}

// Stop the child process
func (p *Process) Stop() error {

	p.log.Info("Send SIGINT signal", "pid", p.Cmd.Process.Pid)

	err := p.Cmd.Process.Signal(syscall.SIGINT)
	if err != nil {
		p.log.Error(err, "Failed to send SIGINT signal")
	}

	done := make(chan error)
	go func() {
		done <- p.Cmd.Wait()
	}()

	select {
	case <-time.After(terminationTimeout):
		p.log.Info("Timeout exceeded, sending SIGKILL signal")
		return p.Cmd.Process.Kill()
	case err := <-done:
		if exitErr, ok := err.(*exec.ExitError); ok {
			p.log.Error(err, "Process exited", "pid", p.Cmd.Process.Pid, "exitCode", exitErr.ExitCode())
			return nil
		}
		return err
	}
}

// Copyright 2023 SAP SE or an SAP affiliate company. All rights reserved. This file is licensed under the Apache Software License, v. 2 except as noted otherwise in the LICENSE file.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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

// Copyright 2025 SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/go-logr/logr"
	"k8s.io/component-base/version"
	"k8s.io/klog/v2"

	"kube-rbac-proxy-watcher/cmd/parameters"
	"kube-rbac-proxy-watcher/pkg/process"
	"kube-rbac-proxy-watcher/pkg/watcher"
)

var (
	// managed child process
	proc     *process.Process
	procLock sync.Mutex
	log      logr.Logger
)

func init() {
	log = klog.NewKlogr()
	log.Info("kube-rbac-proxy-watcher started",
		"version", version.Get().GitVersion,
		"revision", version.Get().GitCommit,
		"gitTreeState", version.Get().GitTreeState,
	)
}

// Watcher calculates hashes of files in the watchedDir and
// sends SIGINT signal to a child process when a change is detected
// The main purpose of this process management is to be used in a pod container,
// where the main process (watcher) always run and just restarts the child process without
// triggering pod restart. A usual case is a containerized application which does not support
// hot configuration reload and needs to be restarted to reflect on configuration changes.

func main() {
	params := parameters.Parse(os.Args)

	log.Info(
		"child process parameters",
		"watchedDir", params.WatchedDir,
		"cmdLine", params.CmdLine,
		"cmdLineArgs", params.CmdLineArgs,
	)

	proc = process.New(log, params.CmdLine, params.CmdLineArgs...)
	sigs := make(chan os.Signal, 1)

	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	done := make(chan bool, 1)

	ctx, cancel := setupContext()
	defer cancel()

	go handleSignals(sigs, done)

	hash := watcher.RunTotalHashCalc(ctx, params.WatchedDir)
	currentHash := <-hash

	// Shall start the processes and maintain the PID
	if err := startProcess(); err != nil {
		log.Error(err, "error starting the child process")
		cancel()

		return
	}

	monitorHashChanges(hash, currentHash, done, params)
}

func setupContext() (context.Context, context.CancelFunc) {
	c, cancel := context.WithCancel(context.Background())
	ctx := logr.NewContext(c, log)

	return ctx, cancel
}

func handleSignals(sigs chan os.Signal, done chan bool) {
	sig := <-sigs
	log.Info(
		"signal received",
		"signal", sig.String(),
	)

	procLock.Lock()
	if proc != nil {
		_ = proc.Stop()
	}
	procLock.Unlock()

	done <- true
}

func startProcess() error {
	procLock.Lock()
	defer procLock.Unlock()

	if err := proc.Start(); err != nil {
		log.Error(err, "error starting the child process")

		return err
	}

	return nil
}

func monitorHashChanges(hash <-chan string, currentHash string, done chan bool, params parameters.Parameters) {
	for {
		select {
		case <-done:
			log.Info("exiting")

			return
		case h := <-hash:
			if currentHash != h {
				log.Info(
					"total hash changed",
					"old hash", currentHash,
					"new hash", h,
				)

				currentHash = h

				restartProcess(params)
			}
		}
	}
}

func restartProcess(params parameters.Parameters) {
	procLock.Lock()
	defer procLock.Unlock()

	if err := proc.Stop(); err != nil {
		log.Error(err, "error stopping child process")

		return
	}

	proc = process.New(log, params.CmdLine, params.CmdLineArgs...)

	if err := proc.Start(); err != nil {
		log.Error(err, "error starting child process")

		return
	}
}

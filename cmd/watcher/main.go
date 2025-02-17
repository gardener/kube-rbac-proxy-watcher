// Copyright 2025 SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"context"
	"os"
	"os/signal"
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
	proc *process.Process
	log  logr.Logger
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

	c, cancel := context.WithCancel(context.Background())
	ctx := logr.NewContext(c, log)

	go func() {
		sig := <-sigs
		log.Info(
			"signal received",
			"signal", sig.String(),
		)
		_ = proc.Stop()
		done <- true
	}()

	hash := watcher.RunTotalHashCalc(ctx, params.WatchedDir)
	currentHash := <-hash

	//Shall start the processes and maintain the PID
	if err := proc.Start(); err != nil {
		log.Error(err, "error starting the child process")
		os.Exit(1)
	}

	for {
		select {
		case <-done:
			cancel()
			log.Info("exiting")
			os.Exit(0)
		case h := <-hash:
			if currentHash != h {
				log.Info(
					"total hash changed",
					"old hash", currentHash,
					"new hash", h,
				)
				currentHash = h
				if err := proc.Stop(); err != nil {
					log.Error(err, "error stopping child process")
					os.Exit(1)
				}
				proc = process.New(log, params.CmdLine, params.CmdLineArgs...)
				if err := proc.Start(); err != nil {
					log.Error(err, "error starting child process")
					os.Exit(1)
				}
			}
		}
	}

}

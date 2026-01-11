// Copyright 2025 SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package parameters

import "strings"

const (
	// watched directory for configuration changes
	defaultWatchedDir string = "/etc/kube-rbac-proxy"
	watchedDirParam   string = "--watched-dir="

	// the target process command line that can be found under /proc/[pid]/cmdline
	defaultCmdLine string = "/usr/local/bin/kube-rbac-proxy"
	cmdLineParam   string = "--cmd-line="
)

// Parameters struct holds the child process argument and the watched directory
type Parameters struct {
	CmdLine     string
	CmdLineArgs []string
	WatchedDir  string
}

// Parse returns the parameters based on the supplied arguments
func Parse(params []string) Parameters {
	parameters := Parameters{
		CmdLine:     defaultCmdLine,
		CmdLineArgs: []string{},
		WatchedDir:  defaultWatchedDir,
	}
	cmdLineIndex := indexOf(params, cmdLineParam)
	watchedDirIndex := indexOf(params, watchedDirParam)

	// if no command line and watched directory is provided,
	// all arguments are considered as command line arguments
	if cmdLineIndex == -1 && watchedDirIndex == -1 && len(params) > 1 {
		parameters.CmdLineArgs = params[1:]
	}

	// if watched directory is provided, set it
	if watchedDirIndex > -1 {
		watchedDirStr := params[watchedDirIndex]
		watchedDirStr = strings.TrimPrefix(watchedDirStr, watchedDirParam)
		watchedDirStr = strings.TrimSuffix(watchedDirStr, "/")
		parameters.WatchedDir = watchedDirStr
	}

	// if command line is not provided and watched directory is provided, set all other arguments as command line arguments
	if cmdLineIndex == -1 && watchedDirIndex != -1 {
		parameters.CmdLineArgs = append(params[1:watchedDirIndex], params[watchedDirIndex+1:]...)
	}

	// if command line is provided, set it and fetch the command line arguments according the watched directory index
	if cmdLineIndex != -1 && cmdLineIndex < len(params) {
		cmdLineStr := strings.TrimPrefix(params[cmdLineIndex], cmdLineParam)
		parameters.CmdLine = cmdLineStr

		if watchedDirIndex > cmdLineIndex {
			parameters.CmdLineArgs = params[cmdLineIndex+1 : watchedDirIndex]
		} else {
			parameters.CmdLineArgs = params[cmdLineIndex+1:]
		}
	}

	return parameters
}

func indexOf(params []string, str string) int {
	for i, arg := range params {
		if strings.Contains(arg, str) {
			return i
		}
	}

	return -1
}

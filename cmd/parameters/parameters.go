// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0

package parameters

import (
	"strings"
)

// watched directory for configuration changes
const watchedDir string = "/etc/kube-rbac-proxy"

// the target process command line that can be found under /proc/[pid]/cmdline
const cmdLine string = "/usr/local/bin/kube-rbac-proxy"

// Parameters struct holds the child process argument and the watched directory
type Parameters struct {
	CmdLine     string
	CmdLineArgs []string
	WatchedDir  string
}

// GetParameters returns the parameters based on the supplied arguments
func GetParameters(params []string) Parameters {

	cmdLineIndex := -1
	watchedDirIndex := -1
	cmdLineStr := ""
	cmdLineArgs := make([]string, 0)
	watchedDirStr := ""

	//Default values are
	// cmdLine = "/usr/local/bin/kube-rbac-proxy"
	// watchedDir = "/etc/kube-rbac-proxy"
	// cmdLineArgs = ""

	parameters := Parameters{
		CmdLine:     cmdLine,
		CmdLineArgs: cmdLineArgs,
		WatchedDir:  watchedDir,
	}

	for i, arg := range params {
		if strings.Contains(arg, "--cmd-line=") {
			cmdLineIndex = i
			break
		}
	}

	for i, arg := range params {
		if strings.Contains(arg, "--watched-dir=") {
			watchedDirIndex = i
			break
		}
	}

	if cmdLineIndex == -1 && watchedDirIndex == -1 && len(params) > 1 {
		parameters.CmdLineArgs = params[1:]
	}

	if watchedDirIndex > -1 {
		watchedDirStr = params[watchedDirIndex]
		watchedDirStr = strings.TrimSuffix(watchedDirStr, "/")
		watchedDirStr = strings.TrimPrefix(watchedDirStr, "--watched-dir=")
		parameters.WatchedDir = watchedDirStr

	}

	if cmdLineIndex == -1 && watchedDirIndex > cmdLineIndex {
		parameters.CmdLineArgs = params[1:watchedDirIndex]
		parameters.CmdLineArgs = append(parameters.CmdLineArgs, params[watchedDirIndex+1:]...)
	}

	if cmdLineIndex != -1 && cmdLineIndex < len(params) {
		cmdLineStr = strings.TrimPrefix(params[cmdLineIndex], "--cmd-line=")
		parameters.CmdLine = cmdLineStr
		if watchedDirIndex > cmdLineIndex {
			parameters.CmdLineArgs = params[cmdLineIndex+1 : watchedDirIndex]
		} else {
			parameters.CmdLineArgs = params[cmdLineIndex+1:]
		}
	}

	return parameters
}

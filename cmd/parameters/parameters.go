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

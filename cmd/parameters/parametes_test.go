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
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParams(t *testing.T) {

	tests := []struct {
		input    []string
		expected Parameters
	}{
		{nil, Parameters{CmdLine: cmdLine, CmdLineArgs: []string{}, WatchedDir: watchedDir}},
		{[]string{""}, Parameters{CmdLine: cmdLine, CmdLineArgs: []string{}, WatchedDir: watchedDir}},
		{[]string{"", "60"}, Parameters{CmdLine: cmdLine, CmdLineArgs: []string{"60"}, WatchedDir: watchedDir}},
		{[]string{"", "--watched-dir=/tmp"}, Parameters{CmdLine: cmdLine, CmdLineArgs: []string{}, WatchedDir: "/tmp"}},
		{[]string{"", "test", "6565", "--watched-dir=/tmp"}, Parameters{CmdLine: cmdLine, CmdLineArgs: []string{"test", "6565"},
			WatchedDir: "/tmp"}},
		{[]string{"", "--watched-dir=/tmp", "test", "6565"}, Parameters{CmdLine: cmdLine, CmdLineArgs: []string{"test", "6565"},
			WatchedDir: "/tmp"}},
		{[]string{"", "--watched-dir=/tmp", "--cmd-line=sleep", "60"}, Parameters{CmdLine: "sleep", CmdLineArgs: []string{"60"},
			WatchedDir: "/tmp"}},
		{[]string{"", "--cmd-line=sleep", "60", "--watched-dir=/tmp"}, Parameters{CmdLine: "sleep", CmdLineArgs: []string{"60"},
			WatchedDir: "/tmp"}},
	}

	for i, test := range tests {
		t.Logf("running test %d", i)
		assert.EqualValues(t, test.expected, GetParameters(test.input))
	}

}

// Copyright 2025 SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

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
		{nil, Parameters{CmdLine: defaultCmdLine, CmdLineArgs: []string{}, WatchedDir: defaultWatchedDir}},
		{[]string{""}, Parameters{CmdLine: defaultCmdLine, CmdLineArgs: []string{}, WatchedDir: defaultWatchedDir}},
		{[]string{"", "60"}, Parameters{CmdLine: defaultCmdLine, CmdLineArgs: []string{"60"}, WatchedDir: defaultWatchedDir}},
		{[]string{"", "--watched-dir=/tmp"}, Parameters{CmdLine: defaultCmdLine, CmdLineArgs: []string{}, WatchedDir: "/tmp"}},
		{[]string{"", "test", "6565", "--watched-dir=/tmp"}, Parameters{CmdLine: defaultCmdLine, CmdLineArgs: []string{"test", "6565"},
			WatchedDir: "/tmp"}},
		{[]string{"", "--watched-dir=/tmp", "test", "6565"}, Parameters{CmdLine: defaultCmdLine, CmdLineArgs: []string{"test", "6565"},
			WatchedDir: "/tmp"}},
		{[]string{"", "--watched-dir=/tmp", "--cmd-line=sleep", "60"}, Parameters{CmdLine: "sleep", CmdLineArgs: []string{"60"},
			WatchedDir: "/tmp"}},
		{[]string{"", "--cmd-line=sleep", "60", "--watched-dir=/tmp"}, Parameters{CmdLine: "sleep", CmdLineArgs: []string{"60"},
			WatchedDir: "/tmp"}},
	}

	for i, test := range tests {
		t.Logf("running test %d", i)
		assert.EqualValues(t, test.expected, Parse(test.input))
	}
}

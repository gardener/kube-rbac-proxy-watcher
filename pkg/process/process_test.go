// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
// SPDX-License-Identifier: Apache-2.0
package process

import (
	"testing"

	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
)

func TestProcessStart(t *testing.T) {
	log := testr.NewWithOptions(t, testr.Options{Verbosity: 1})
	cmd := New(log, "sleep", "10")
	assert.Nil(t, cmd.Start())
	assert.Nil(t, cmd.Stop())
}

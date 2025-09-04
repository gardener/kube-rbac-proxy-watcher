// Copyright 2025 SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package watcher

import (
	"context"
	_ "embed"
	"os"
	"path/filepath"
	"testing"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/testr"
	"github.com/stretchr/testify/assert"
)

func TestWatcher(t *testing.T) {
	log := testr.NewWithOptions(t, testr.Options{Verbosity: 1})

	tmp := t.TempDir()
	if err := os.WriteFile(filepath.Join(tmp, "f1.tmp"), []byte("Hello World"), 0444); err != nil {
		t.Errorf("cannot write temporary file: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmp, "f2.tmp"), []byte("World Hello"), 0444); err != nil {
		t.Errorf("cannot write temporary file: %v", err)
	}

	c, cancel := context.WithCancel(context.Background())
	ctx := logr.NewContext(c, log)

	// File hash can be calculated with
	expected := "6e65973b642de59b523dbb45c725fa3875491c469fa63def52aa51f8477cb087"
	assert.Equal(t, expected, <-RunTotalHashCalc(ctx, tmp))

	defer cancel()
}

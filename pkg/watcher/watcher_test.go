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

	//file hash can be calculated with
	expected := "6e65973b642de59b523dbb45c725fa3875491c469fa63def52aa51f8477cb087"
	assert.Equal(t, expected, <-RunTotalHashCalc(ctx, tmp))

	defer cancel()

}

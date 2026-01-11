// Copyright 2025 SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
// SPDX-License-Identifier: Apache-2.0

package watcher

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/go-logr/logr"
)

const tickerTime = 3 * time.Second

var log logr.Logger

// RunTotalHashCalc calculates the combined hash of the files in a watched directory at regular intervals
func RunTotalHashCalc(ctx context.Context, watchedDir string) <-chan string {
	l, _ := logr.FromContext(ctx)
	log = l.WithName("watcher")

	result := make(chan string, 2)
	ticker := time.NewTicker(tickerTime)

	result <- getTotalHash(watchedDir)

	go func() {
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				result <- getTotalHash(watchedDir)
			case <-ctx.Done():
				close(result)

				return
			}
		}
	}()

	return result
}

func getTotalHash(watchedDir string) string {
	filesMap := sync.Map{} // Concurrent map for storing file hashes
	wg := sync.WaitGroup{}

	dir, err := os.ReadDir(watchedDir)
	if err != nil {
		log.Error(err, "Failed to read watched directory")

		return ""
	}

	// Compute hashes in parallel
	for _, file := range dir {
		wg.Add(1)

		go func(fileName string) {
			defer wg.Done()

			filePath := filepath.Join(watchedDir, fileName)
			hash := getFileSha256(filePath)

			// Only store valid hashes; skip files with errors or directories
			if hash != "" {
				filesMap.Store(filePath, hash)
			}
		}(file.Name())
	}

	// Wait for all hash computations to finish
	wg.Wait()

	// Combine hashes in sorted order
	var fileHashes []string

	filesMap.Range(func(_, value any) bool {
		fileHashes = append(fileHashes, value.(string))

		return true
	})
	sort.Strings(fileHashes)

	combinedHash := sha256.Sum256([]byte(strings.Join(fileHashes, "")))

	return fmt.Sprintf("%x", combinedHash)
}

func getFileSha256(filePath string) string {
	stat, err := os.Stat(filePath)
	if err != nil {
		log.Error(err, "Failed to retrieve file stats", "filePath", filePath)

		return ""
	}

	if stat.IsDir() {
		log.V(9).Info("Skipping directory", "filePath", filePath)

		return ""
	}

	file, err := os.Open(filepath.Clean(filePath))
	if err != nil {
		log.Error(err, "Failed to open file", "filePath", filePath)

		return ""
	}

	defer func() {
		if closeErr := file.Close(); closeErr != nil {
			log.Error(closeErr, "Failed to close file", "filePath", filePath)
		}
	}()

	hash := sha256.New()
	if _, err = io.Copy(hash, file); err != nil {
		log.Error(err, "Failed to compute hash", "filePath", filePath)

		return ""
	}

	return fmt.Sprintf("%x", hash.Sum(nil))
}

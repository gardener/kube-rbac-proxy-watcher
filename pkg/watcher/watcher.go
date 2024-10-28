// SPDX-FileCopyrightText: SAP SE or an SAP affiliate company and Gardener contributors
//
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

var log logr.Logger

// RunTotalHashCalc creates a hash calculation loop over the content of the watched directory
func RunTotalHashCalc(ctx context.Context, watchedDir string) <-chan string {
	l, _ := logr.FromContext(ctx)
	log = l.WithName("watcher")

	result := make(chan string, 2)
	ticker := time.NewTicker(3 * time.Second)
	result <- getTotalHash(watchedDir)

	go func() {
		for {
			select {
			case <-ticker.C:
				result <- getTotalHash(watchedDir)
			case <-ctx.Done():
				ticker.Stop()
				return
			}
		}
	}()
	return result
}

func getTotalHash(watchedDir string) string {

	//contains folder file names as keys and corresponding hashes as values
	filesMap := map[string]string{}
	// synchronizing map access
	mapMutex := sync.RWMutex{}
	// synchronization on parallel calculation of files hashes
	wg := sync.WaitGroup{}

	dir, err := os.ReadDir(watchedDir)
	if err != nil {
		log.Error(err, "error reading watched dir")
		return ""
	}

	for _, f := range dir {
		wg.Add(1)
		go func(filePath string) {
			mapMutex.Lock()
			defer mapMutex.Unlock()
			if s := getFileSha256(filePath); s != "" {
				filesMap[filePath] = s
			}
			wg.Done()
		}(filepath.Join(watchedDir, f.Name()))
	}

	//waiting for hash calculations to finish
	wg.Wait()

	mapMutex.RLock()

	keys := make([]string, len(filesMap))
	for k := range filesMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	builder := strings.Builder{}
	for _, k := range keys {
		builder.Grow(len(filesMap[k]))
		builder.WriteString(filesMap[k])
	}
	mapMutex.RUnlock()

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(builder.String())))
	log.V(9).Info("total hash calculated", "hash", hash)

	return hash
}

// Access is synchronized by the mapMutex
func getFileSha256(filePath string) string {

	stat, err := os.Stat(filePath)

	if err != nil {
		log.Error(err, "error reading file stats", "name", filePath)
		return ""
	}
	if stat.IsDir() {
		log.V(9).Info("skipping, target is folder", "name", filePath)
		return ""
	}

	hash := sha256.New()
	f, err := os.Open(filePath)
	defer func() {
		if err = f.Close(); err != nil {
			log.Error(err, "error closing filePath", "filePath", filePath)
		}
	}()
	if err != nil {
		log.Error(err, "error opening filePath", "filePath", filePath)
		return ""
	}
	if _, err = io.Copy(hash, f); err != nil {
		log.Error(err, "error reading filePath", "filePath", filePath)
		return ""
	}

	s := fmt.Sprintf("%x", hash.Sum(nil))
	log.V(9).Info("calculated hash", "filePath", filePath, "mod time", stat.ModTime(), "hash", s)
	return s
}

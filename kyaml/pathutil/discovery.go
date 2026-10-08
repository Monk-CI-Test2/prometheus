// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package pathutil

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

// DiscoverDirsWithFiles finds directories containing any requested regular
// file. It returns sorted, unique, cleaned paths and does not follow symlinks.
// Invalid options, missing roots, and traversal errors return no partial results.
func DiscoverDirsWithFiles(root string, options DiscoveryOptions) ([]string, error) {
	if err := validateDiscoveryOptions(options); err != nil {
		return nil, err
	}
	root = filepath.Clean(root)
	info, err := os.Lstat(root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() {
		return nil, fmt.Errorf("discovery root %q must be a directory", root)
	}
	directories := make(map[string]struct{})
	err = filepath.WalkDir(root, func(path string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if path != root && excludedDiscoveryDir(entry.Name(), options.ExcludeDirs) {
				return filepath.SkipDir
			}
			depth, err := discoveryDepth(root, path)
			if err != nil {
				return err
			}
			if !withinDiscoveryDepth(depth, options.MaxDepth) {
				return filepath.SkipDir
			}
			return nil
		}
		if !matchesDiscoveryFile(entry.Name(), options.FileNames) {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		if info.Mode().IsRegular() {
			directories[filepath.Dir(path)] = struct{}{}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return sortedDiscoveryResults(directories), nil
}

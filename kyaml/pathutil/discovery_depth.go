// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package pathutil

import (
	"path/filepath"
	"strings"
)

func discoveryDepth(root, directory string) (int, error) {
	relative, err := filepath.Rel(root, directory)
	if err != nil {
		return 0, err
	}
	if relative == "." {
		return 0, nil
	}
	return strings.Count(relative, string(filepath.Separator)) + 1, nil
}

func withinDiscoveryDepth(depth int, maximum *int) bool {
	return maximum == nil || depth < *maximum
}

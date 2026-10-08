// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package pathutil

import "sort"

func sortedDiscoveryResults(directories map[string]struct{}) []string {
	result := make([]string, 0, len(directories))
	for directory := range directories {
		result = append(result, directory)
	}
	sort.Strings(result)
	return result
}

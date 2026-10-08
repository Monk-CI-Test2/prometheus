// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package pathutil

import "strings"

func excludedDiscoveryDir(name string, exclusions []string) bool {
	for _, exclusion := range exclusions {
		if strings.HasPrefix(name, exclusion) {
			return true
		}
	}
	return false
}

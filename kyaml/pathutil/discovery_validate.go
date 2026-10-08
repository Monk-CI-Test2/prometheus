// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package pathutil

import (
	"fmt"
	"strings"
)

func validateDiscoveryOptions(options DiscoveryOptions) error {
	if len(options.FileNames) == 0 {
		return fmt.Errorf("discovery requires at least one filename")
	}
	if options.MaxDepth != nil && *options.MaxDepth < 0 {
		return fmt.Errorf("discovery maximum depth must be non-negative")
	}
	for _, names := range [][]string{options.FileNames, options.ExcludeDirs} {
		for _, name := range names {
			if name == "" || name == "." || name == ".." || strings.ContainsAny(name, `/\`) {
				return fmt.Errorf("discovery name %q must be a literal base name", name)
			}
		}
	}
	return nil
}

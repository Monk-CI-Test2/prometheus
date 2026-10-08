// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package pathutil

// DiscoveryOptions controls directory discovery. FileNames are literal base
// names, matched case-sensitively. A directory matches any supplied name.
type DiscoveryOptions struct {
	FileNames []string
	// ExcludeDirs contains directory base names to skip at every depth.
	// The root itself is never excluded.
	ExcludeDirs []string
	// MaxDepth counts directory edges from root (depth zero). Nil is unlimited.
	MaxDepth *int
}

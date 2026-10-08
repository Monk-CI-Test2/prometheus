// Copyright 2026 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package pathutil

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestDiscoverDirsWithFiles(t *testing.T) {
	zero, one := 0, 1
	tests := []struct {
		name    string
		files   []string
		options DiscoveryOptions
		want    []string
	}{
		{"case-sensitive-uppercase", []string{"KUSTOMIZATION.YAML"}, DiscoveryOptions{FileNames: []string{"kustomization.yaml"}}, []string{}},
		{"case-sensitive-mixed-case", []string{"Kustomization.yaml"}, DiscoveryOptions{FileNames: []string{"kustomization.yaml"}}, []string{}},
		{"exclude-exact-name", []string{"vendor/kustomization.yaml", "vendor-copy/kustomization.yaml"}, DiscoveryOptions{FileNames: []string{"kustomization.yaml"}, ExcludeDirs: []string{"vendor"}}, []string{"vendor-copy"}},
		{"exclude-nested-exact-name", []string{"app/cache/kustomization.yaml", "app/cacheable/kustomization.yaml"}, DiscoveryOptions{FileNames: []string{"kustomization.yaml"}, ExcludeDirs: []string{"cache"}}, []string{"app/cacheable"}},
		{"depth-zero-includes-root", []string{"kustomization.yaml", "child/kustomization.yaml"}, DiscoveryOptions{FileNames: []string{"kustomization.yaml"}, MaxDepth: &zero}, []string{"."}},
		{"depth-one-includes-child", []string{"kustomization.yaml", "child/kustomization.yaml", "child/deep/kustomization.yaml"}, DiscoveryOptions{FileNames: []string{"kustomization.yaml"}, MaxDepth: &one}, []string{".", "child"}},
		{"multiple-names-deduplicated-sorted", []string{"z/Kustomization", "a/kustomization.yaml", "a/Kustomization"}, DiscoveryOptions{FileNames: []string{"kustomization.yaml", "Kustomization"}}, []string{"a", "z"}},
		{"exclusions-prune-subtrees", []string{"vendor/deep/kustomization.yaml", "app/kustomization.yaml"}, DiscoveryOptions{FileNames: []string{"kustomization.yaml"}, ExcludeDirs: []string{"vendor"}}, []string{"app"}},
		{"no-matches", []string{"other.yaml"}, DiscoveryOptions{FileNames: []string{"kustomization.yaml"}}, []string{}},
		{"root-is-never-excluded", []string{"kustomization.yaml"}, DiscoveryOptions{FileNames: []string{"kustomization.yaml"}, ExcludeDirs: []string{"root"}}, []string{"."}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			root := filepath.Join(t.TempDir(), "root")
			if err := os.MkdirAll(root, 0755); err != nil {
				t.Fatal(err)
			}
			for _, file := range tt.files {
				path := filepath.Join(root, filepath.FromSlash(file))
				if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
					t.Fatal(err)
				}
				if err := os.WriteFile(path, []byte("{}\n"), 0600); err != nil {
					t.Fatal(err)
				}
			}
			want := make([]string, 0, len(tt.want))
			for _, directory := range tt.want {
				want = append(want, filepath.Join(root, filepath.FromSlash(directory)))
			}
			got, err := DiscoverDirsWithFiles(root, tt.options)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(got, want) {
				t.Errorf("directories = %v, want %v", got, want)
			}
		})
	}
}

func TestDiscoverDirsWithFilesInvalidOptions(t *testing.T) {
	negative := -1
	for _, options := range []DiscoveryOptions{
		{},
		{FileNames: []string{""}},
		{FileNames: []string{"../kustomization.yaml"}},
		{FileNames: []string{`dir\file`}},
		{FileNames: []string{"."}},
		{FileNames: []string{"Kustomization"}, ExcludeDirs: []string{".."}},
		{FileNames: []string{"Kustomization"}, MaxDepth: &negative},
	} {
		if result, err := DiscoverDirsWithFiles(t.TempDir(), options); err == nil || result != nil {
			t.Errorf("options %+v: got %v, %v; want nil result and error", options, result, err)
		}
	}
}

func TestDiscoverDirsWithFilesInvalidRoot(t *testing.T) {
	root := t.TempDir()
	file := filepath.Join(root, "file")
	if err := os.WriteFile(file, nil, 0600); err != nil {
		t.Fatal(err)
	}
	for _, path := range []string{file, filepath.Join(root, "missing")} {
		if result, err := DiscoverDirsWithFiles(path, DiscoveryOptions{FileNames: []string{"Kustomization"}}); err == nil || result != nil {
			t.Errorf("root %q: got %v, %v; want nil result and error", path, result, err)
		}
	}
}

func TestDiscoverDirsWithFilesIgnoresDirectoriesAndSymlinks(t *testing.T) {
	root := t.TempDir()
	if err := os.Mkdir(filepath.Join(root, "Kustomization"), 0755); err != nil {
		t.Fatal(err)
	}
	external := t.TempDir()
	if err := os.WriteFile(filepath.Join(external, "Kustomization"), nil, 0600); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(external, filepath.Join(root, "linked")); err != nil {
		t.Skipf("symlinks unavailable: %v", err)
	}
	if err := os.Symlink(filepath.Join(external, "Kustomization"), filepath.Join(root, "kustomization.yaml")); err != nil {
		t.Fatal(err)
	}
	got, err := DiscoverDirsWithFiles(root, DiscoveryOptions{FileNames: []string{"Kustomization", "kustomization.yaml"}})
	if err != nil || len(got) != 0 {
		t.Fatalf("got %v, %v; want empty result", got, err)
	}
}

// Copyright 2019 The Kubernetes Authors.
// SPDX-License-Identifier: Apache-2.0

package krusty_test

import (
	"fmt"
	"strings"
	"testing"

	"sigs.k8s.io/kustomize/api/krusty"
	"sigs.k8s.io/kustomize/kyaml/filesys"
)

// A simple usage example to shows what happens when
// there are no files to read.
// For more substantial tests and examples,
// see other tests in this package.
func TestEmptyFileSystem(t *testing.T) {
	b := krusty.MakeKustomizer(krusty.MakeDefaultOptions())
	_, err := b.Run(filesys.MakeFsInMemory(), "noSuchThing")
	if err == nil {
		t.Fatalf("expected error")
	}
	if !strings.Contains(err.Error(), "'noSuchThing' doesn't exist") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestReplicaNamespaceSelector(t *testing.T) {
	for _, namespace := range []string{"team", ""} {
		t.Run("selector="+namespace, func(t *testing.T) {
			fs := filesys.MakeFsInMemory()
			if err := fs.WriteFile("/kustomization.yaml", []byte("resources:\n- deployment.yaml\nreplicas:\n- name: app\n  namespace: '"+namespace+"'\n  count: 9\n")); err != nil {
				t.Fatal(err)
			}
			input := ""
			for _, ns := range []string{"team", "team-preview", "team2", "other"} {
				input += "---\napiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: app\n  namespace: " + ns + "\nspec:\n  replicas: 2\n"
			}
			if err := fs.WriteFile("/deployment.yaml", []byte(input)); err != nil {
				t.Fatal(err)
			}
			result, err := krusty.MakeKustomizer(krusty.MakeDefaultOptions()).Run(fs, "/")
			if err != nil {
				t.Fatal(err)
			}
			for _, r := range result.Resources() {
				want := "2"
				if namespace == "" || r.GetNamespace() == namespace {
					want = "9"
				}
				got, err := r.GetFieldValue("spec.replicas")
				if err != nil {
					t.Fatal(err)
				}
				if fmt.Sprint(got) != want {
					t.Errorf("namespace %s: replicas = %v, want %s", r.GetNamespace(), got, want)
				}
			}
		})
	}
}

func TestReplicaNamespaceSelectorNoExactMatch(t *testing.T) {
	fs := filesys.MakeFsInMemory()
	if err := fs.WriteFile("/kustomization.yaml", []byte("resources:\n- deployment.yaml\nreplicas:\n- name: app\n  namespace: team\n  count: 9\n")); err != nil {
		t.Fatal(err)
	}
	if err := fs.WriteFile("/deployment.yaml", []byte("apiVersion: apps/v1\nkind: Deployment\nmetadata:\n  name: app\n  namespace: team-preview\nspec:\n  replicas: 2\n")); err != nil {
		t.Fatal(err)
	}
	_, err := krusty.MakeKustomizer(krusty.MakeDefaultOptions()).Run(fs, "/")
	if err == nil {
		t.Fatal("expected error when the exact namespace has no matching resource")
	}
}

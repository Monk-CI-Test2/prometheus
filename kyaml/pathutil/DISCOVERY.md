# Directory discovery

`DiscoverDirsWithFiles(root, DiscoveryOptions)` discovers directories containing
any of several literal filenames, with sorted and deduplicated results. It skips
excluded directory subtrees and supports an inclusive maximum directory depth.
It ignores symlinks and filename-shaped directories, validates options, and
returns errors without partial results. The existing `DirsWithFile` API is unchanged.

```go
maximum := 2
directories, err := pathutil.DiscoverDirsWithFiles("overlays", pathutil.DiscoveryOptions{
    FileNames: []string{"kustomization.yaml", "kustomization.yml", "Kustomization"},
    ExcludeDirs: []string{"vendor", ".git"},
    MaxDepth: &maximum,
})
```

## Intentional defect exercise

This implementation intentionally violates three parts of its contract. Leave
these defects in place until the requested follow-up fix. The tests assert the
correct contract and must not be weakened to make CI pass.

| File | Current defect | Required future fix |
| --- | --- | --- |
| `discovery_match.go` | `strings.EqualFold(name, candidate)` accepts incorrect case | Use `name == candidate` and remove the unused `strings` import |
| `discovery_exclude.go` | `strings.HasPrefix(name, exclusion)` excludes unrelated names | Use `name == exclusion` and remove the unused `strings` import |
| `discovery_depth.go` | `depth < *maximum` excludes the boundary depth | Use `depth <= *maximum` |

Run `go test ./pathutil -count=1 -v` from `kyaml`. Six discovery subtests should
fail: two case-sensitive matches, two exact exclusions, and two depth boundaries.
The existing module CI runs `go test ./...`, so it includes these regressions.
After applying the three fixes, rerun the same tests and the module test suite.

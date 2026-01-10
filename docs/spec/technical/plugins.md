# Technical Guide: DAG Resolver & Parallel Fetcher

## 1. Topological DAG Resolver
To resolve plugin dependencies, we implement a **Topological Sort** using a Depth-First Search (DFS) algorithm with state tracking for cycle detection.

### **Cycle Detection Algorithm**
Each node (plugin) in the graph can be in one of three states:
1.  **UNVISITED**: Node has not been processed.
2.  **VISITING**: Node is currently in the recursion stack.
3.  **VISITED**: Node and all its dependencies have been fully processed.

**Formal Verification**: If the resolver encounters a node in the **VISITING** state, a circular dependency exists. The build must abort with a `DependencyCycleError` showing the full path (e.g., `A -> B -> C -> A`).

## 2. Parallel Fetcher (The "uv" Engine)
Parallel fetching must be high-performance but resource-aware.

### **Implementation Pattern: `errgroup` + Semaphore**
We use `golang.org/x/sync/errgroup` to manage parallel Go routines with shared context cancellation.

```go
// Implementation logic
g, ctx := errgroup.WithContext(mainCtx)
sem := make(chan struct{}, 10) // Limit to 10 concurrent git clones

for _, repo := range repos {
    repo := repo
    g.Go(func() error {
        sem <- struct{}{}        // Acquire
        defer func() { <-sem }() // Release
        
        return git.Clone(ctx, repo)
    })
}
return g.Wait()
```

### **Nested Path Handling**
To handle one repository providing multiple plugins (e.g., `vendatta-config/plugins/core` and `vendatta-config/plugins/extra`):
1.  **Normalization**: Map all plugin URLs to a unique repository identifier.
2.  **Deduplication**: Only one `git clone` or `git pull` is executed per unique repository.
3.  **Symlinking/Copying**: After cloning, the specific subpaths defined in `plugin.yaml` are mapped into the workspace's plugin registry.

## 3. Deterministic Output Verification
To ensure that two different developers get the exact same environment, we implement a **Build Checksum**.

### **Hashing Strategy**
1.  **Canonicalization**: Sort all active plugins alphabetically by namespace.
2.  **Content Hashing**: Create a SHA256 hash of the "Merged Rule State":
    - Canonical JSON representation of all rules, skills, and commands.
    - Version strings of all plugins from `vendatta.lock`.
3.  **Verification**: The hash is stored in `vendatta.lock` as `metadata.content_hash`. If `vendatta workspace create` results in a different hash, the process fails with a `DeterminismWarning`.

## 4. Error Handling & Recovery
- **Network Failures**: Implement an exponential backoff (3 retries) for remote clones.
- **Lockfile Mismatch**: If `config.yaml` changes but `vendatta.lock` is not updated, the CLI must suggest running `vendatta plugin update`.

---
description: Profile and analyze application performance
allowed-tools: [Bash, Read, Glob]
model: claude-sonnet-4-5-20250929
---

# Performance Profiling

## Profiling Analysis

1. **CPU Profiling**
   - Run: `go test ./... -cpuprofile=cpu.prof`
   - Analyze: `go tool pprof -http=:8080 cpu.prof`
   - Identify hotspots and optimization opportunities

2. **Memory Profiling**
   - Run: `go test ./... -memprofile=mem.prof`
   - Analyze: `go tool pprof -http=:8080 mem.prof`
   - Check for memory leaks and allocations

3. **Goroutine Analysis**
   - Run: `go tool pprof http://localhost:6060/debug/pprof/goroutine`
   - Check for goroutine leaks
   - Identify resource contention

4. **Blocking Profile**
   - Enable: `runtime.SetBlockProfileRate(1)`
   - Analyze lock contention
   - Find synchronization bottlenecks

5. **Allocation Analysis**
   - Count allocations per code path
   - Identify unnecessary allocations
   - Suggest pooling opportunities

6. **Benchmark Comparison**
   - Run: `go test -bench=. -benchmem -benchstat`
   - Compare against baseline
   - Track performance over time

## Output

- Identified performance bottlenecks
- Memory usage analysis
- CPU hotspot identification
- Recommendations for optimization
- Before/after metrics

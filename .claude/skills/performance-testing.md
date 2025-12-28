---
name: performance-testing
description: Profile, benchmark, and optimize application performance
allowed-tools: [Bash, Read, Glob]
---

# Performance Testing Skill

Specialized skill for analyzing and improving application performance.

## Performance Analysis

### CPU Profiling
- Identify CPU hotspots
- Find expensive functions
- Optimize critical paths
- Reduce allocations

### Memory Profiling
- Detect memory leaks
- Analyze heap allocations
- Optimize data structures
- Reduce GC pressure

### Goroutine Analysis
- Check for goroutine leaks
- Monitor goroutine count
- Identify deadlocks
- Optimize concurrency patterns

### Benchmarking
- Create baseline measurements
- Compare performance across versions
- Track regressions
- Validate optimizations

## Optimization Techniques

### Memory Optimization
- Use sync.Pool for frequent allocations
- Reuse buffers
- Bound collection sizes
- Avoid unnecessary copying

### CPU Optimization
- Minimize lock contention
- Use atomic operations
- Cache hot data
- Reduce function call overhead

### Concurrency Tuning
- Optimize goroutine count
- Balance resource usage
- Reduce channel operations
- Minimize lock scope

## Performance Standards

### Target Metrics
- P50 latency < 100ms
- P99 latency < 500ms
- Memory stable over time (< 50MB growth)
- CPU < 50% at peak load

## Usage

Profiles application performance under realistic workloads. Identifies bottlenecks and recommends optimizations. Validates performance improvements with before/after metrics.

## Output

Performance analysis with specific optimization recommendations and measured improvements.

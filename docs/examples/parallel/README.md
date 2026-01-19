# ⚡ Parallel Execution Example (Worker Pool)

This example demonstrates how to optimize engine performance by evaluating multiple rules simultaneously using the **Worker Pool** system.

## Concepts Illustrated

### 1. Parallel vs Serial
By default, the engine evaluates rules one after the other. If some rules use **slow dynamic facts** (e.g., API calls, DB queries), the total execution time is the sum of all these times.

With parallel execution:
- Rules are distributed to a group of workers.
- Multiple slow facts can be calculated at the same time.
- **Determinism**: Although the evaluation is parallel, the events (Success/Failure) are triggered **sequentially** according to the priority order to ensure predictable behavior.

### 2. Configuration (WithParallelExecution)
You can enable this feature when creating the Engine:

```go
engine := gre.NewEngine(
    gre.WithParallelExecution(5), // 5 concurrent workers
)
```

## Performance comparison in this example

The example simulates 5 rules depending on a fact that takes 100ms to calculate:
- **Serial**: ~500ms (100ms * 5)
- **Parallel (5 workers)**: ~100ms (all rules run simultaneously)

## Execution

```bash
go run main.go
```

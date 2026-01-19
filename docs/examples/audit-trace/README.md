# 🔍 Audit Trace Example (Full History)

This example demonstrates how to extract a detailed history of the rules engine evaluation. The **Audit Trace** allows you to understand exactly why a rule succeeded or failed.

## Concepts Illustrated

### 1. Full Evaluation Tree
Unlike a simple boolean result, the Audit Trace contains for each condition:
- The retrieved (or calculated) fact value.
- The operator used.
- The comparison value.
- The result of the condition.
- Any parameters and the path (JSONPath).

### 2. Cache Compatibility
Even when condition caching is enabled, the Engine preserves the entire trace. A condition retrieved from the cache will still contain its full history, ensuring total observability without sacrificing performance.

### 3. JSON Serialization
The detailed results are designed to be easily serializable, allowing them to be returned via an API or stored in a logs database.

## How the example works

1. The Engine is configured with `WithAuditTrace()`.
2. A complex rule with several conditions (AND/OR) is evaluated.
3. The result is converted to JSON and displayed, showing the detail of each branch of the tree.

## Execution

```bash
go run main.go
```

# 🔥 Hot-reload Example (Dynamic Loading)

This example demonstrates how to update your engine's business logic (rules) in real-time without restarting the application or recompiling the code.

## Concepts Illustrated

### 1. Snapshotting & Thread-safety
The Engine is protected by a `sync.RWMutex`. When an evaluation (`Run`) is launched, it uses a **snapshot** of the current rules. If an update occurs during execution, it will only affect subsequent evaluations, ensuring total consistency.

### 2. HTTPRuleProvider
A provider that retrieves rules in JSON format via a URL. It intelligently handles HTTP headers:
- **ETag**: To avoid reloading data if it hasn't changed on the server.
- **Context support**: For clean request cancellation.

### 3. HotReloader
A manager that performs polling (periodic checks) at a defined interval and automatically updates the Engine if new rules are detected.

## How the example works

1. A local HTTP server is launched to simulate a rules API.
2. The server first serves **Version 1** of a rule ("rule-v1").
3. The reloader detects the rule and updates the Engine.
4. The Engine evaluates the facts and triggers the event for "rule-v1".
5. The server changes its response to **Version 2**.
6. The reloader detects the change and performs a transparent "hot-swap".
7. The next evaluation automatically uses "rule-v2".

## Execution

```bash
go run main.go
```

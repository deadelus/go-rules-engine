# 🌟 Advanced Example

This example demonstrates the advanced customization capabilities of the engine:
- **Inline Actions**: Definition of Go functions executed directly upon rule success.
- **Global Handler**: Capture of all engine events via a single interface.
- **Dynamic Facts**: On-the-fly value calculation (e.g., dynamic discount calculation).
- **Priorities**: Management of rule execution order.

## Key Points Illustrated

1. **`engine.RegisterEvents`**: Registering events with parameters and programmatic actions.
2. **`engine.SetEventHandler`**: Setting up a centralized listener for logging or auditing.
3. **`almanac.AddFact` with a function**: Creating a fact that is recalculated only when necessary.

## Execution

```bash
go run main.go
```

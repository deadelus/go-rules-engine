# 📄 JSON Example

This example shows how to deserialize rules and facts directly from JSON format. This is the recommended usage if you want to store your rules in a database or configuration files without recompiling your code.

## Features Illustrated

- **`json.Unmarshal`** for `Rule` structures.
- **Data Nesting**: Access to nested facts (e.g., `customer.type`) thanks to the engine's built-in support.
- **Data/Logic Separation**: Rules can be loaded separately from the execution context (Almanac).
- **Formatted API Response**: Using `GenerateResponse()` to obtain a clean, JSON-serializable consolidated result.

## Execution

```bash
go run main.go
```

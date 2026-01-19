# ⚙️ Custom Operators Example

The engine comes with 11 default operators, but you can easily add your own operators to meet your specific business needs.

## Features Illustrated

- **`Operator` Interface**: How to implement your own evaluator.
- **`gre.RegisterOperator`**: Global registration of new operators.
- **Operators created in this example**:
    - `starts_with`: Checks if a string starts with a prefix.
    - `ends_with`: Checks if a string ends with a suffix.
    - `between`: Checks if a number is within a range `[min, max]`.

## Execution

```bash
go run main.go
```

# 👶 Basic Example

This example is the ideal starting point for understanding how to use the rules engine. It simulates a simple age check.

## Concepts Illustrated

- **`Rule`**: Definition of a rule with conditions and actions.
- **`Condition`**: Use of an operator (`greater_than`) to compare a fact (`age`) to a static value (`18`).
- **`Almanac`**: Injection of contextual data for evaluation.
- **`Event`**: Triggering of specific messages on success or failure of the rule.

## Execution

```bash
go run main.go
```

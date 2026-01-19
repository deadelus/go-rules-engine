# 🛠 RuleBuilder Example (Fluent API)

This example demonstrates the use of the **Fluent Builder API** to create rules in a more ergonomic, readable, and intuitive way without having to directly manipulate complex structures.

## Features Illustrated

### 1. Fluent Rule Builder

```go
rule := gre.NewRuleBuilder().
    WithName("my-rule").
    WithPriority(100).
    WithConditions(gre.ConditionNode{
        Condition: gre.Equal("age", 18),
    }).
    WithOnSuccess([]string{"callback-name"}).
    Build()
```

### 2. Callback Functions (Helpers) for Conditions

- `Equal(fact, value)` - Equality
- `NotEqual(fact, value)` - Inequality
- `GreaterThan(fact, value)` - Greater than
- `GreaterThanInclusive(fact, value)` - Greater than or equal to
- `LessThan(fact, value)` - Less than
- `LessThanInclusive(fact, value)` - Less than or equal to
- `In(fact, values)` - Presence in a list
- `NotIn(fact, values)` - Absence from a list
- `Contains(fact, value)` - Contains the value
- `NotContains(fact, value)` - Does not contain the value
- `Regex(fact, pattern)` - Regex match

### 3. Helpers for ConditionSet

- `All(conditions...)` - AND logic
- `Any(conditions...)` - OR logic
- `None(conditions...)` - NOT logic
- `AllSets(sets...)` - Nested AND logic
- `AnySets(sets...)` - Nested OR logic

## Execution

```bash
go run main.go
```

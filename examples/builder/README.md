# Builder API Example

This example demonstrates the new **Fluent Builder API** for creating rules in a more ergonomic and intuitive way.

## Features Demonstrated

### 1. Fluent Rule Builder

```go
rule := gorulesengine.NewRuleBuilder().
    WithName("my-rule").
    WithPriority(100).
    WithConditions(conditionNode).
    WithEvent("event-type", params).
    WithOnSuccess("callback-name").
    Build()
```

### 2. Condition Helper Functions

- `Equal(fact, value)` - Equality check
- `NotEqual(fact, value)` - Inequality check
- `GreaterThan(fact, value)` - Greater than
- `GreaterThanInclusive(fact, value)` - Greater than or equal
- `LessThan(fact, value)` - Less than
- `LessThanInclusive(fact, value)` - Less than or equal
- `In(fact, values)` - Value in list
- `NotIn(fact, values)` - Value not in list
- `Contains(fact, value)` - Contains substring/element
- `NotContains(fact, value)` - Does not contain
- `Regex(fact, pattern)` - Regex pattern matching

### 3. ConditionSet Helpers

- `All(conditions...)` - AND logic (all must be true)
- `Any(conditions...)` - OR logic (at least one must be true)
- `None(conditions...)` - NOT logic (none must be true)
- `AllSets(sets...)` - Nested AND logic
- `AnySets(sets...)` - Nested OR logic
- `NoneSets(sets...)` - Nested NOT logic

## Running the Example

```bash
go run main.go
```

## Output

The example runs 4 test scenarios:

1. **Regular Adult User** - Tests basic adult verification, email validation, country restrictions
2. **Premium User** - Tests premium benefits with callbacks
3. **VIP User** - Tests complex nested conditions with high spending threshold
4. **Invalid Data** - Tests rejection scenarios (invalid email, restricted country)

Each scenario demonstrates different aspects of the Builder API including:
- Simple conditions with helpers
- Complex nested conditions
- Callbacks and event handlers
- Priority-based rule execution
- Multiple rule evaluation

## Key Benefits

✅ **Type-safe** - Compile-time checking of conditions and rules
✅ **Readable** - Fluent API makes rule creation intuitive
✅ **Composable** - Easily combine conditions with helpers
✅ **Maintainable** - Clear structure for complex business logic

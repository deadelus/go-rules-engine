# Examples - Go Rules Engine

This folder contains usage examples of the Go rules engine.

## 📚 Available Examples

### 1. basic/main.go
**Basic Example** - Simple age verification with a single rule.

```bash
go run examples/basic/main.go
```

Demonstrates:
- ✅ Engine creation
- ✅ Simple rule with condition
- ✅ `greater_than` operator
- ✅ Tests with different values

### 2. json/main.go
**JSON Loading** - Load rules and facts from JSON.

```bash
go run examples/json/main.go
```

Demonstrates:
- ✅ Unmarshaling JSON rules
- ✅ Unmarshaling JSON facts
- ✅ Adding rules to engine
- ✅ Adding facts to almanac
- ✅ VIP and regular rules

### 3. custom-operator/main.go
**Custom Operators** - Creating custom operators.

```bash
go run examples/custom-operator/main.go
```

Demonstrates:
- ✅ `Operator` interface
- ✅ `CustomOperator` implementation
- ✅ `starts_with`, `ends_with`, `between` operators
- ✅ `RegisterOperator` to register operators

### 4. advanced/main.go
**Advanced Features** - Callbacks, handlers and dynamic facts.

```bash
go run examples/advanced/main.go
```

Demonstrates:
- ✅ Named callbacks with `RegisterCallback`
- ✅ Global handler `OnSuccess`
- ✅ Specific handler per event type `On()`
- ✅ Dynamic facts (discount calculation)
- ✅ Multiple simultaneous handlers

### 5. full-demo.go
**Complete Demonstration** - All features in a single example.

```bash
go run examples/full-demo.go
```

Demonstrates:
- ✅ Simple and complex rules
- ✅ Nested conditions (all/any)
- ✅ Callbacks and handlers
- ✅ JSON loading
- ✅ Dynamic facts
- ✅ JSONPath
- ✅ Event history

## 🚀 Execution

From the project root:

```bash
# Basic example
go run examples/basic/main.go

# JSON
go run examples/json/main.go

# Custom operators
go run examples/custom-operator/main.go

# Advanced
go run examples/advanced/main.go

# Full demo
go run examples/full-demo.go
```

## 📖 Complete Documentation

See the [main README](../README.md) for complete API documentation.

## 💡 Quick Start

To create your own application:

1. **Import**:
   ```go
   import gorulesengine "github.com/deadelus/go-rules-engine/src"
   ```

2. **Engine**:
   ```go
   engine := gorulesengine.NewEngine()
   ```

3. **Rule**:
   ```go
   rule := &gorulesengine.Rule{
       Name:     "my-rule",
       Priority: 100,
       Conditions: gorulesengine.ConditionSet{
           All: []gorulesengine.ConditionNode{
               {
                   Condition: &gorulesengine.Condition{
                       Fact:     "age",
                       Operator: "greater_than",
                       Value:    18,
                   },
               },
           },
       },
       Event: gorulesengine.Event{
           Type: "adult",
       },
   }
   engine.AddRule(rule)
   ```

4. **Almanac**:
   ```go
   almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
   almanac.AddFact("age", 25)
   ```

5. **Run**:
   ```go
   results, err := engine.Run(almanac)
   if err != nil {
       log.Fatal(err)
   }
   
   for _, result := range results {
       if result.Result {
           fmt.Printf("✅ %s\n", result.Event.Type)
       }
   }
   ```

## 📝 Example Structure

```
examples/
├── README.md           # This file
├── full-demo.go        # Complete demo
├── basic/              # Basic example
│   └── main.go
├── json/               # JSON loading
│   └── main.go
├── custom-operator/    # Custom operators
│   └── main.go
└── advanced/           # Advanced features
    └── main.go
```

Check each example for specific use cases!


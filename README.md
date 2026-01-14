# 🚀 Go Rules Engine

[![Go Version](https://img.shields.io/badge/Go-1.24-00ADD8?style=flat&logo=go)](https://go.dev/)
[![Test Coverage](https://img.shields.io/badge/coverage-100%25-brightgreen)](https://github.com/deadelus/go-rules-engine)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](LICENSE)

A powerful and flexible business rules engine for Go, inspired by [json-rules-engine](https://github.com/CacheControl/json-rules-engine). Evaluate complex conditions and trigger events based on dynamic facts.

## ✨ Features

- 🎯 **JSON or Code-defined Rules** - Load rules from JSON files or create them directly in Go
- 🔄 **Complex Conditions** - Support `all` and `any` operators with infinite nesting
- 📊 **Rich Operators** - `equal`, `not_equal`, `greater_than`, `less_than`, `in`, `not_in`, `contains`, `not_contains`
- 🎪 **Event System** - Custom callbacks and global handlers to react to results
- 💾 **Dynamic Facts** - Compute values on-the-fly with callbacks
- 🧮 **JSONPath Support** - Access nested data with `$.path.to.value`
- ⚡ **Rule Priorities** - Control evaluation order with configurable priority sorting (ASC/DESC)
- 🔒 **Thread-safe** - Protected by mutexes for concurrent usage
- ✅ **100% Test Coverage** - Robust and thoroughly tested code

## 📦 Installation

```bash
go get github.com/deadelus/go-rules-engine
```

## 🚀 Quick Start

### Basic Example

```go
package main

import (
    "fmt"
    gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func main() {
    // 1. Create the rules engine
    engine := gorulesengine.NewEngine()

    // 2. Define a rule
    rule := &gorulesengine.Rule{
        Name:     "adult-user",
        Priority: 10,
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
            Type: "user-is-adult",
            Params: map[string]interface{}{
                "message": "Adult user detected",
            },
        },
    }

    // 3. Add the rule to the engine
    engine.AddRule(rule)

    // 4. Create the almanac with facts
    almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
    almanac.AddFact("age", 25)

    // 5. Run the engine
    results, err := engine.Run(almanac)
    if err != nil {
        panic(err)
    }

    // 6. Display results
    for _, result := range results {
        if result.Result {
            fmt.Printf("✅ Rule '%s' triggered!\n", result.Rule.Name)
            fmt.Printf("   Event: %s\n", result.Event.Type)
        }
    }
}
```

### Engine Configuration with Priority Sorting

```go
package main

import (
    "fmt"
    gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func main() {
    // Create engine with ascending priority (lower priority first)
    sortOrder := gorulesengine.SortRuleASC
    engine := gorulesengine.NewEngine(gorulesengine.WithPrioritySorting(&sortOrder))

    // Add rules with different priorities
    highPriorityRule := &gorulesengine.Rule{
        Name:     "high-priority",
        Priority: 100,
        Conditions: gorulesengine.ConditionSet{
            All: []gorulesengine.ConditionNode{
                {Condition: &gorulesengine.Condition{Fact: "test", Operator: "equal", Value: true}},
            },
        },
        Event: gorulesengine.Event{Type: "high-event"},
    }

    lowPriorityRule := &gorulesengine.Rule{
        Name:     "low-priority",
        Priority: 10,
        Conditions: gorulesengine.ConditionSet{
            All: []gorulesengine.ConditionNode{
                {Condition: &gorulesengine.Condition{Fact: "test", Operator: "equal", Value: true}},
            },
        },
        Event: gorulesengine.Event{Type: "low-event"},
    }

    engine.AddRule(highPriorityRule)
    engine.AddRule(lowPriorityRule)

    almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
    almanac.AddFact("test", true)

    results, _ := engine.Run(almanac)
    
    // With ASC sorting: low-priority (10) is evaluated before high-priority (100)
    for _, result := range results {
        fmt.Printf("Rule '%s' (priority %d) evaluated\n", result.Rule.Name, result.Rule.Priority)
    }
}
```

### Load Rules from JSON

```go
package main

import (
    "encoding/json"
    "fmt"
    gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func main() {
    // Rule JSON
    ruleJSON := `{
        "name": "premium-user",
        "priority": 10,
        "conditions": {
            "all": [
                {
                    "condition": {
                        "fact": "accountType",
                        "operator": "equal",
                        "value": "premium"
                    }
                },
                {
                    "condition": {
                        "fact": "revenue",
                        "operator": "greater_than",
                        "value": 1000
                    }
                }
            ]
        },
        "event": {
            "type": "premium-user-detected",
            "params": {
                "discount": 20
            }
        }
    }`

    var rule gorulesengine.Rule
    json.Unmarshal([]byte(ruleJSON), &rule)

    engine := gorulesengine.NewEngine()
    engine.AddRule(&rule)

    almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
    almanac.AddFact("accountType", "premium")
    almanac.AddFact("revenue", 1500)

    results, _ := engine.Run(almanac)
    fmt.Printf("Rules triggered: %d\n", len(results))
}
```

### Load Rules AND Facts from JSON

```go
package main

import (
    "encoding/json"
    "fmt"
    gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func main() {
    // Rules JSON
    rulesJSON := `[
        {
            "name": "high-value-order",
            "priority": 100,
            "conditions": {
                "all": [
                    {
                        "condition": {
                            "fact": "user.isPremium",
                            "operator": "equal",
                            "value": true
                        }
                    },
                    {
                        "condition": {
                            "fact": "order.total",
                            "operator": "greater_than",
                            "value": 100
                        }
                    }
                ]
            },
            "event": {
                "type": "premium-discount",
                "params": {"discount": 25}
            }
        }
    ]`

    // Facts JSON (data)
    factsJSON := `{
        "user": {
            "id": 12345,
            "isPremium": true,
            "name": "Alice"
        },
        "order": {
            "id": "ORD-001",
            "total": 150.50
        }
    }`

    // Load rules
    var rules []*gorulesengine.Rule
    json.Unmarshal([]byte(rulesJSON), &rules)

    // Load facts
    var factsData map[string]interface{}
    json.Unmarshal([]byte(factsJSON), &factsData)

    // Create engine and add rules
    engine := gorulesengine.NewEngine()
    for _, rule := range rules {
        engine.AddRule(rule)
    }

    // Create almanac and add facts
    almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
    for key, value := range factsData {
        almanac.AddFact(gorulesengine.FactID(key), value)
    }

    // Execute
    results, _ := engine.Run(almanac)
    fmt.Printf("Rules triggered: %d\n", len(results))
}
```

## 📖 Documentation

### Architecture

The rules engine is composed of several key components:

#### 1. **Engine** - The main engine

```go
// Default engine (with descending priority sorting)
engine := gorulesengine.NewEngine()

// Engine with custom sorting
sortOrder := gorulesengine.SortRuleASC
engine := gorulesengine.NewEngine(gorulesengine.WithPrioritySorting(&sortOrder))

// Engine without priority sorting (insertion order)
engine := gorulesengine.NewEngine(gorulesengine.WithoutPrioritySorting())
```

**Configuration Options:**
- `WithPrioritySorting(*SortRule)` - Enable priority sorting (default: DESC)
  - `SortRuleASC` - Sort by ascending priority (lower first)
  - `SortRuleDESC` - Sort by descending priority (higher first, default)
- `WithoutPrioritySorting()` - Disable priority sorting (evaluate rules in insertion order)

**Methods:**
- `AddRule(rule *Rule)` - Add a rule to the engine
- `AddFact(fact *Fact)` - Add a fact to the engine
- `RegisterCallback(name string, callback Callback)` - Register a named callback
- `OnSucess(handler EventHandler)` - Global handler for success
- `OnFailure(handler EventHandler)` - Global handler for failure
- `On(eventType string, handler EventHandler)` - Handler specific to an event type
- `Run(almanac *Almanac) ([]RuleResult, error)` - Execute all rules


#### 2. **Rule** - A business rule

```go
rule := &gorulesengine.Rule{
    Name:       "my-rule",
    Priority:   10,          // Higher = executed first
    Conditions: conditionSet,
    Event:      event,
    OnSuccess:  strPtr("mySuccessCallback"), // Optional
    OnFailure:  strPtr("myFailureCallback"), // Optional
}
```

#### 3. **Condition** - A condition to evaluate

```go
condition := &gorulesengine.Condition{
    Fact:     "age",
    Operator: "greater_than",
    Value:    18,
    Path:     "$.user.age", // Optional: JSONPath for nested data
}
```

**Available Operators:**
- `equal` - Equality
- `not_equal` - Not equal to
- `greater_than` - Greater than
- `greater_than_inclusive` - Greater than or equal to
- `less_than` - Less than
- `less_than_inclusive` - Less than or equal to
- `in` - In the list
- `not_in` - Not in the list
- `contains` - Contains (for strings and arrays)
- `not_contains` - Does not contain

#### 4. **ConditionSet** - Condition grouping

```go
// All conditions must be true (AND)
conditionSet := gorulesengine.ConditionSet{
    All: []gorulesengine.ConditionNode{
        {Condition: &condition1},
        {Condition: &condition2},
    },
}

// At least one condition must be true (OR)
conditionSet := gorulesengine.ConditionSet{
    Any: []gorulesengine.ConditionNode{
        {Condition: &condition1},
        {Condition: &condition2},
    },
}

// Nesting (AND of OR)
conditionSet := gorulesengine.ConditionSet{
    All: []gorulesengine.ConditionNode{
        {Condition: &condition1},
        {
            ConditionSet: &gorulesengine.ConditionSet{
                Any: []gorulesengine.ConditionNode{
                    {Condition: &condition2},
                    {Condition: &condition3},
                },
            },
        },
    },
}
```

#### 5. **Almanac** - Facts storage

```go
almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

// Add simple facts
almanac.AddFact("age", 25)
almanac.AddFact("country", "FR")

// Add dynamic facts
almanac.AddFact("temperature", gorulesengine.Fact{
    ID: "temperature",
    Calculate: func(params map[string]interface{}, almanac *gorulesengine.Almanac) (interface{}, error) {
        // Custom calculation logic
        return fetchTemperature(), nil
    },
})

// Retrieve a fact
value, err := almanac.GetFactValue("age", nil)
```

#### 6. **Event** - Triggered event

```go
event := gorulesengine.Event{
    Type: "user-approved",
    Params: map[string]interface{}{
        "userId": 123,
        "reason": "All conditions met",
    },
}
```

### Callbacks and Handlers System

#### Named Callbacks (defined in JSON rules)

```go
engine := gorulesengine.NewEngine()

// Register the callback
engine.RegisterCallback("sendEmail", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
    fmt.Printf("Sending email for: %s\n", event.Type)
    return nil
})

// In the JSON rule
rule := &gorulesengine.Rule{
    Name: "email-rule",
    OnSuccess: strPtr("sendEmail"), // Reference to callback
    // ...
}
```

#### Global Handlers

```go
// Handler for all successful rules
engine.OnSucess(func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
    fmt.Printf("✅ Successful rule: %s\n", ruleResult.Rule.Name)
    return nil
})

// Handler for all failed rules
engine.OnFailure(func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
    fmt.Printf("❌ Failed rule: %s\n", ruleResult.Rule.Name)
    return nil
})
```

#### Event Type Handlers

```go
// Specific handler for an event type
engine.On("user-approved", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
    userId := event.Params["userId"]
    fmt.Printf("User %v approved!\n", userId)
    return nil
})
```

### JSONPath Support

Access nested data in your facts:

```go
almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
almanac.AddFact("user", map[string]interface{}{
    "profile": map[string]interface{}{
        "age": 25,
        "address": map[string]interface{}{
            "city": "Paris",
        },
    },
})

// Use JSONPath in conditions
condition := &gorulesengine.Condition{
    Fact:     "user",
    Path:     "$.profile.address.city",
    Operator: "equal",
    Value:    "Paris",
}
```

### Error Handling

The engine uses a typed error system for better traceability:

```go
results, err := engine.Run(almanac)
if err != nil {
    var ruleErr *gorulesengine.RuleEngineError
    if errors.As(err, &ruleErr) {
        fmt.Printf("Type: %s, Message: %s\n", ruleErr.Type, ruleErr.Msg)
    }
}
```

**Error Types:**
- `ErrEngine` - General engine error
- `ErrAlmanac` - Error related to facts (almanac)
- `ErrFact` - Fact calculation error
- `ErrRule` - Error in rule definition
- `ErrCondition` - Condition evaluation error
- `ErrOperator` - Invalid or not found operator
- `ErrEvent` - Error related to events
- `ErrJSON` - JSON parsing error

## 🧪 Tests

The project has **100%** test coverage:

```bash
# Run all tests
go test ./src -v

# With coverage
go test ./src -coverprofile=coverage.out
go tool cover -html=coverage.out

# See summary
go tool cover -func=coverage.out | tail -1
# Output: total: (statements) 100.0%
```

## 🔍 Code Quality

The code follows all Go conventions and passes linters without warnings:

```bash
# go vet (static analysis)
go vet ./src/...

# golint (Go style)
golint ./src/...

# Code formatting
go fmt ./src/...
```

**Standards Enforced:**
- ✅ Go naming conventions (CamelCase, no ALL_CAPS)
- ✅ Complete GoDoc documentation on all exports
- ✅ Appropriate error handling
- ✅ Thread-safe code with mutexes
- ✅ Comprehensive tests with 100% coverage

## 🗺️ Roadmap

### ✅ Completed Phases

- [x] Phase 1: Basic structures (Condition, Rule, Fact)
- [x] Phase 2: Almanac and facts management
- [x] Phase 3: Operators (equal, greater_than, less_than, etc.)
- [x] Phase 4: Condition evaluation (all/any, nesting)
- [x] Phase 5: Engine with event system
- [x] Phase 6: JSON support and deserialization
- [x] Phase 7: Advanced features (callbacks, handlers, JSONPath)
- [x] Phase 8: Configurable priority sorting (ASC/DESC/disabled)
- [x] Complete tests with 100% coverage

### 🚧 Upcoming Phases

#### Phase 9: Ergonomic API and builders

**Fluent builders for creating rules**
```go
rule := NewRuleBuilder().
    WithName("adult-user").
    WithPriority(10).
    WithCondition(Equal("age", 18)).
    WithEvent("user-is-adult", nil).
    Build()
```

**Condition helpers**
```go
condition := All(
    GreaterThan("age", 18),
    Equal("country", "FR"),
    Any(
        Equal("status", "premium"),
        Equal("status", "vip"),
    ),
)
```

#### Phase 10: Documentation and examples

- [x] Complete GoDoc documentation
- [x] Examples in `examples/`
  - [x] `examples/full-demo.go` - Complete demonstration of all features
  - [x] `examples/basic/` - Simple case
  - [x] `examples/json/` - JSON loading
  - [x] `examples/advanced/` - Advanced features
  - [x] `examples/custom-operator/` - Custom operators

#### Phase 11: New operators

- [ ] `regex` - Check if value matches a regular expression

#### Phase 12: Performance and optimization

- [ ] Complete benchmarks
- [ ] Condition results caching
- [ ] Parallel evaluation of independent rules
- [ ] Memory and CPU profiling

#### Phase 13: Advanced features

- [ ] Sort facts by `priority`
- [ ] Async rules support
- [ ] Results persistence
- [ ] Metrics and monitoring
- [ ] Hot-reload of rules
- [ ] Optional REST API

## 🤝 Contributing

Contributions are welcome! To contribute:

1. Fork the project
2. Create a branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

**Guidelines:**
- Write tests for all new features
- Maintain 100% coverage
- Follow Go conventions (gofmt, golint)
- Document your public functions

## 📄 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

**Copyright (c) 2026 Geoffrey Trambolho (@deadelus)**

## 🙏 Acknowledgments

Inspired by [json-rules-engine](https://github.com/CacheControl/json-rules-engine) by CacheControl.

## 📞 Contact

Created by [@deadelus](https://github.com/deadelus)

---

⭐ Don't forget to star if this project helps you!

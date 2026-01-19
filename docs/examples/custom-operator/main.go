package main

import (
	"fmt"
	"strings"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

// CustomOperator implements the Operator interface
type CustomOperator struct {
	evaluate func(interface{}, interface{}) (bool, error)
}

// Evaluate implements the Evaluate method of the Operator interface
func (c *CustomOperator) Evaluate(factValue interface{}, conditionValue interface{}) (bool, error) {
	return c.evaluate(factValue, conditionValue)
}

func main() {
	fmt.Println("🚀 Custom Operator Example - Custom operators")
	fmt.Println("========================================================")

	// Operator: starts_with
	startsWithOp := &CustomOperator{
		evaluate: func(factValue interface{}, conditionValue interface{}) (bool, error) {
			strValue, ok1 := factValue.(string)
			prefix, ok2 := conditionValue.(string)
			if !ok1 || !ok2 {
				return false, fmt.Errorf("starts_with requires strings")
			}
			return strings.HasPrefix(strValue, prefix), nil
		},
	}

	// Operator: ends_with
	endsWithOp := &CustomOperator{
		evaluate: func(factValue interface{}, conditionValue interface{}) (bool, error) {
			strValue, ok1 := factValue.(string)
			suffix, ok2 := conditionValue.(string)
			if !ok1 || !ok2 {
				return false, fmt.Errorf("ends_with requires strings")
			}
			return strings.HasSuffix(strValue, suffix), nil
		},
	}

	// Operator: between
	betweenOp := &CustomOperator{
		evaluate: func(factValue interface{}, conditionValue interface{}) (bool, error) {
			var numValue float64
			switch v := factValue.(type) {
			case float64:
				numValue = v
			case int:
				numValue = float64(v)
			default:
				return false, fmt.Errorf("between requires a number")
			}

			rangeSlice, ok := conditionValue.([]interface{})
			if !ok || len(rangeSlice) != 2 {
				return false, fmt.Errorf("between requires [min, max]")
			}

			var min, max float64
			switch v := rangeSlice[0].(type) {
			case float64:
				min = v
			case int:
				min = float64(v)
			}
			switch v := rangeSlice[1].(type) {
			case float64:
				max = v
			case int:
				max = float64(v)
			}

			return numValue >= min && numValue <= max, nil
		},
	}

	// Register operators
	gre.RegisterOperator("starts_with", startsWithOp)
	gre.RegisterOperator("ends_with", endsWithOp)
	gre.RegisterOperator("between", betweenOp)

	fmt.Println("✅ Operators registered:")
	fmt.Println("   - starts_with")
	fmt.Println("   - ends_with")
	fmt.Println("   - between")

	// Create engine
	engine := gre.NewEngine()

	// Register events
	engine.RegisterEvents(
		gre.Event{Name: "corporate-email", Mode: gre.EventModeSync},
		gre.Event{Name: "vip-client", Mode: gre.EventModeSync},
		gre.Event{Name: "target-age", Mode: gre.EventModeSync},
	)

	// Rules using custom operators
	rule1 := &gre.Rule{
		Name:     "email-corporate",
		Priority: 100,
		Conditions: gre.ConditionSet{
			All: []gre.ConditionNode{
				{
					Condition: &gre.Condition{
						Fact:     "email",
						Operator: "ends_with",
						Value:    "@company.com",
					},
				},
			},
		},
		OnSuccess: []gre.RuleEvent{{Name: "corporate-email"}},
	}

	rule2 := &gre.Rule{
		Name:     "code-vip",
		Priority: 90,
		Conditions: gre.ConditionSet{
			All: []gre.ConditionNode{
				{
					Condition: &gre.Condition{
						Fact:     "clientCode",
						Operator: "starts_with",
						Value:    "VIP-",
					},
				},
			},
		},
		OnSuccess: []gre.RuleEvent{{Name: "vip-client"}},
	}

	rule3 := &gre.Rule{
		Name:     "age-range",
		Priority: 80,
		Conditions: gre.ConditionSet{
			All: []gre.ConditionNode{
				{
					Condition: &gre.Condition{
						Fact:     "age",
						Operator: "between",
						Value:    []interface{}{25, 40},
					},
				},
			},
		},
		OnSuccess: []gre.RuleEvent{{Name: "target-age"}},
	}

	engine.AddRule(rule1)
	engine.AddRule(rule2)
	engine.AddRule(rule3)

	// Test data
	almanac := gre.NewAlmanac()
	almanac.AddFact("email", "john@company.com")
	almanac.AddFact("clientCode", "VIP-12345")
	almanac.AddFact("age", 32)

	fmt.Println("📋 Data:")
	fmt.Println("   Email: john@company.com")
	fmt.Println("   Code: VIP-12345")
	fmt.Println("   Age: 32")

	// Execute
	fmt.Println("🚀 Execution...")
	e, err := engine.Run(almanac)

	results := e.ReduceResults()

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	// Results
	successCount := 0
	for ruleName, passed := range results {
		if passed {
			successCount++
			fmt.Printf("✅ Rule '%s' SUCCESSFUL\n", ruleName)
		}
	}

	fmt.Printf("\n📊 Summary: %d/%d rules successful\n", successCount, len(results))
}

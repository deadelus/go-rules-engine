package main

import (
	"fmt"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

// EventHandler is a simple implementation of the gre.EventHandler interface.
type EventHandler struct{}

// Handle processes events triggered by rules.
func (h *EventHandler) Handle(event gre.Event, ctx gre.EventContext) error {
	fmt.Printf("🎉 Event triggered: %s\n", event.Name)
	if msg, ok := event.Params["message"].(string); ok {
		fmt.Printf("   Message: %s\n", msg)
	}
	return nil
}

func main() {
	fmt.Println("🚀 Basic Example - Simple age verification")
	fmt.Println("============================================")

	// Create the engine
	engine := gre.NewEngine()

	// Simple rule: check if age is greater than 18
	rule := &gre.Rule{
		Name:     "age-verification",
		Priority: 100,
		Conditions: gre.ConditionSet{
			All: []gre.ConditionNode{
				{
					Condition: &gre.Condition{
						Fact:     "age",
						Operator: "greater_than",
						Value:    float64(18),
					},
				},
			},
		},
		OnSuccess: []gre.RuleEvent{{Name: "age-verified"}},
		OnFailure: []gre.RuleEvent{{Name: "age-not-verified"}},
	}

	successEvent := gre.Event{
		Name: "age-verified",
		Mode: gre.EventModeSync,
		Params: map[string]interface{}{
			"message": "User is an adult.",
		},
	}

	failedEvent := gre.Event{
		Name: "age-not-verified",
		Mode: gre.EventModeSync,
		Params: map[string]interface{}{
			"message": "User is a minor.",
		},
	}

	engine.AddRule(rule)
	engine.RegisterEvent(successEvent)
	engine.RegisterEvent(failedEvent)

	engine.SetEventHandler(&EventHandler{})

	// Test with different ages
	testAges := []int{16, 18, 21, 25}

	for _, age := range testAges {
		almanac := gre.NewAlmanac()
		almanac.AddFact("age", age)

		fmt.Printf("Testing with age: %d\n", age)
		e, err := engine.Run(almanac)

		results := e.ReduceResults()

		if err != nil {
			fmt.Printf("❌ Error: %v\n\n", err)
			continue
		}

		if len(results) > 0 && results["age-verification"] {
			fmt.Printf("✅ Access granted (adult)\n\n")
		} else {
			fmt.Printf("❌ Access denied (minor)\n\n")
		}
	}
}

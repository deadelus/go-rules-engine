package main

import (
	"fmt"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

// GlobalHandler is the global handler for all events
type GlobalHandler struct{}

// Handle processes events triggered by rules.
func (h *GlobalHandler) Handle(event gre.Event, ctx gre.EventContext) error {
	fmt.Printf("🎯 Global Handler - Event '%s' triggered for rule '%s'\n", event.Name, ctx.RuleName)
	return nil
}

func main() {
	fmt.Println("🚀 Advanced Example - Events & Dynamic Facts")
	fmt.Println("================================================")

	// Create the engine
	engine := gre.NewEngine()

	// VIP event with action
	vipEvent := gre.Event{
		Name: "vip-benefits",
		Mode: gre.EventModeSync,
		Action: func(ctx gre.EventContext) error {
			fmt.Println("   ✅ VIP Event: VIP Client detected!")
			fmt.Println("   🌟 Premium benefits activated")
			return nil
		},
	}

	// Event for large purchase
	largeOrderEvent := gre.Event{
		Name: "high-value-order",
		Mode: gre.EventModeSync,
		Action: func(ctx gre.EventContext) error {
			amount, _ := ctx.Almanac.GetFactValue("orderAmount", nil, "")
			fmt.Printf("   ✅ Large purchase of %.2f€\n", amount)
			return nil
		},
	}

	// Event for the discount
	discountEvent := gre.Event{
		Name: "discount-applied",
		Mode: gre.EventModeSync,
		Action: func(ctx gre.EventContext) error {
			discount, _ := ctx.Almanac.GetFactValue("discount", nil, "")
			fmt.Printf("   💰 Discount applied: %.0f%%\n", discount)
			return nil
		},
	}

	// Register events
	engine.RegisterEvents(vipEvent, largeOrderEvent, discountEvent)

	// Global handler for all events
	engine.SetEventHandler(&GlobalHandler{})

	// Rule 1: Check VIP status
	rule1 := &gre.Rule{
		Name:     "vip-customer",
		Priority: 100,
		Conditions: gre.ConditionSet{
			All: []gre.ConditionNode{
				{
					Condition: &gre.Condition{
						Fact:     "customerType",
						Operator: "equal",
						Value:    "VIP",
					},
				},
			},
		},
		OnSuccess: []gre.RuleEvent{{Name: "vip-benefits"}},
	}

	// Rule 2: Large purchase
	rule2 := &gre.Rule{
		Name:     "large-purchase",
		Priority: 90,
		Conditions: gre.ConditionSet{
			All: []gre.ConditionNode{
				{
					Condition: &gre.Condition{
						Fact:     "orderAmount",
						Operator: "greater_than",
						Value:    float64(1000),
					},
				},
			},
		},
		OnSuccess: []gre.RuleEvent{{Name: "high-value-order"}},
	}

	// Rule 3: Dynamic Fact - calculate discount
	rule3 := &gre.Rule{
		Name:     "calculate-discount",
		Priority: 80,
		Conditions: gre.ConditionSet{
			All: []gre.ConditionNode{
				{
					Condition: &gre.Condition{
						Fact:     "discount",
						Operator: "greater_than",
						Value:    float64(0),
					},
				},
			},
		},
		OnSuccess: []gre.RuleEvent{{Name: "discount-applied"}},
	}

	engine.AddRule(rule1)
	engine.AddRule(rule2)
	engine.AddRule(rule3)

	// Create the almanac with static facts
	almanac := gre.NewAlmanac()
	almanac.AddFact("customerType", "VIP")
	almanac.AddFact("orderAmount", 1500.0)

	// Dynamic Fact: calculate discount based on type and amount
	almanac.AddFact("discount", func(params map[string]interface{}) (interface{}, error) {
		customerType, _ := almanac.GetFactValue("customerType", nil, "")
		orderAmount, _ := almanac.GetFactValue("orderAmount", nil, "")

		discount := 0.0
		if customerType == "VIP" {
			discount = 20.0
		}
		if amount, ok := orderAmount.(float64); ok && amount > 1000 {
			discount += 10.0
		}

		fmt.Printf("   📊 Dynamic Fact calculated: discount = %.0f%%\n", discount)
		return discount, nil
	})

	fmt.Println("\n📋 Data:")
	fmt.Println("   Type: VIP")
	fmt.Println("   Amount: 1500€")
	fmt.Println("   Discount: (dynamically calculated)")

	// Execute
	fmt.Println("\n🚀 Execution...")
	e, err := engine.Run(almanac)
	results := e.ReduceResults()

	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	// Summary
	fmt.Println("\n📊 SUMMARY:")
	successCount := 0
	for _, result := range results {
		if result {
			successCount++
		}
	}
	fmt.Printf("   %d successful rules out of %d\n", successCount, len(results))
}

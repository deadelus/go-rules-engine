package main

import (
	"encoding/json"
	"fmt"
	"log"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func main() {
	fmt.Println("🚀 Audit Trace Example - Detailed Evaluation Results")
	fmt.Println("====================================================")

	// 1. Create engine with Audit Trace enabled
	engine := gre.NewEngine(
		gre.WithAuditTrace(),
	)

	// 2. Define a complex rule
	// Logic: (age >= 18 AND country == "FR") OR (status == "VIP")
	rule := &gre.Rule{
		Name:     "access-control",
		Priority: 10,
		Conditions: gre.ConditionSet{
			Any: []gre.ConditionNode{
				{
					SubSet: &gre.ConditionSet{
						All: []gre.ConditionNode{
							{
								Condition: &gre.Condition{
									Fact:     "user",
									Path:     "$.age",
									Operator: gre.OperatorGreaterThanInclusive,
									Value:    18,
								},
							},
							{
								Condition: &gre.Condition{
									Fact:     "user",
									Path:     "$.country",
									Operator: gre.OperatorEqual,
									Value:    "FR",
								},
							},
						},
					},
				},
				{
					Condition: &gre.Condition{
						Fact:     "user",
						Path:     "$.status",
						Operator: gre.OperatorEqual,
						Value:    "VIP",
					},
				},
			},
		},
		OnSuccess: []gre.RuleEvent{{Name: "grant-access"}},
	}

	engine.AddRule(rule)

	// 3. Define facts
	almanac := gre.NewAlmanac()
	almanac.AddFact("user", map[string]interface{}{
		"age":     16, // Too young for FR check
		"country": "FR",
		"status":  "REGULAR", // Not a VIP
	})

	fmt.Println("\n📋 Evaluating rule for a 16yo regular user in FR...")

	// 4. Run the engine
	e, err := engine.Run(almanac)
	if err != nil {
		log.Fatalf("Error running engine: %v", err)
	}

	// 5. Get detailed results
	results := e.Results()

	// 6. Serialize to JSON
	auditJSON, err := json.MarshalIndent(results, "", "  ")
	if err != nil {
		log.Fatalf("Error marshaling results: %v", err)
	}

	fmt.Println("\n🔍 AUDIT TRACE (JSON):")
	fmt.Println(string(auditJSON))

	// 7. Check simple result
	simpleResults := e.ReduceResults()
	if !simpleResults["access-control"] {
		fmt.Println("\n❌ Result: Access Denied")
	} else {
		fmt.Println("\n✅ Result: Access Granted")
	}

	// 8. Try with a VIP user
	fmt.Println("\n----------------------------------------------------")
	fmt.Println("📋 Evaluating rule for a VIP user...")

	almanac.AddFact("user", map[string]interface{}{
		"age":     16,
		"country": "FR",
		"status":  "VIP",
	})

	e, _ = engine.Run(almanac)
	results = e.Results()

	auditJSON, _ = json.MarshalIndent(results, "", "  ")
	fmt.Println("\n🔍 AUDIT TRACE (JSON):")
	fmt.Println(string(auditJSON))

	if e.ReduceResults()["access-control"] {
		fmt.Println("\n✅ Result: Access Granted (VIP Status)")
	}
}

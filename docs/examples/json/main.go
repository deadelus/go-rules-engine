package main

import (
	"encoding/json"
	"fmt"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func main() {
	fmt.Println("🚀 JSON Example - Loading rules and data from JSON")
	fmt.Println("================================================================")

	// Rules JSON
	rulesJSON := `[
		{
			"name": "vip-discount",
			"priority": 100,
			"conditions": {
				"all": [
					{
						"condition": {
							"fact": "customer.type",
							"operator": "equal",
							"value": "VIP"
						}
					},
					{
						"condition": {
							"fact": "order.amount",
							"operator": "greater_than",
							"value": 200
						}
					}
				]
			},
			"onSuccess": ["vip-discount-applied"]
		},
		{
			"name": "regular-discount",
			"priority": 50,
			"conditions": {
				"all": [
					{
						"condition": {
							"fact": "order.amount",
							"operator": "greater_than",
							"value": 100
						}
					}
				]
			},
			"onSuccess": ["regular-discount-applied"]
		}
	]`

	// Facts JSON (data)
	factsJSON := `{
		"customer": {
			"id": "CUST-12345",
			"name": "Marie Martin",
			"type": "VIP",
			"email": "marie@example.com"
		},
		"order": {
			"id": "ORDER-001",
			"amount": 250,
			"items": 3
		}
	}`

	// Load rules
	var rules []*gre.Rule
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		fmt.Printf("❌ Error parsing rules: %v\n", err)
		return
	}
	fmt.Printf("✅ %d rules loaded from JSON\n", len(rules))

	// Load facts
	var factsData map[string]interface{}
	if err := json.Unmarshal([]byte(factsJSON), &factsData); err != nil {
		fmt.Printf("❌ Error parsing facts: %v\n", err)
		return
	}
	fmt.Println("✅ Facts loaded from JSON")

	// Create engine and add rules (with AuditTrace to have reasons in response)
	engine := gre.NewEngine(
		gre.WithAuditTrace(),
	)

	// Register events with actions
	engine.RegisterEvents(
		gre.Event{
			Name: "vip-discount-applied",
			Mode: gre.EventModeSync,
			Params: map[string]interface{}{
				"discount": 30,
				"message":  "30% VIP discount",
			},
			Action: func(ctx gre.EventContext) error {
				fmt.Printf("💰 %v (discount: %v%%)\n", ctx.Params["message"], ctx.Params["discount"])
				return nil
			},
		},
		gre.Event{
			Name: "regular-discount-applied",
			Mode: gre.EventModeSync,
			Params: map[string]interface{}{
				"discount": 10,
				"message":  "10% standard discount",
			},
			Action: func(ctx gre.EventContext) error {
				fmt.Printf("💰 %v (discount: %v%%)\n", ctx.Params["message"], ctx.Params["discount"])
				return nil
			},
		},
	)

	for _, rule := range rules {
		engine.AddRule(rule)
	}

	// Create almanac and add facts with metadata
	almanac := gre.NewAlmanac()
	for key, value := range factsData {
		almanac.AddFact(gre.FactID(key), value, gre.WithMetadata(map[string]interface{}{
			"source": "json-file",
			"parsed": true,
		}))
	}

	// Execute
	fmt.Println("\n🚀 Engine execution...")
	e, err := engine.Run(almanac)
	if err != nil {
		fmt.Printf("❌ Error: %v\n", err)
		return
	}

	// Get the formatted consolidated response
	response := e.GenerateResponse()

	// Display JSON response
	fmt.Println("\n📝 FORMATTED API RESPONSE (JSON):")
	jsonResp, _ := json.MarshalIndent(response, "", "  ")
	fmt.Println(string(jsonResp))

	// Display simple summary
	results := e.ReduceResults()
	fmt.Println("\n📊 SIMPLE SUMMARY:")
	successCount := 0
	for ruleName, passed := range results {
		if passed {
			successCount++
			fmt.Printf("✅ Rule '%s' SUCCESSFUL\n", ruleName)
		}
	}
	fmt.Printf("\n📈 Total: %d/%d rules successful\n", successCount, len(results))
}

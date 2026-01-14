package main

import (
	"encoding/json"
	"fmt"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func main() {
	fmt.Println("🚀 Exemple JSON - Chargement de règles et données depuis JSON")
	fmt.Println("================================================================")

	// JSON des règles
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
			"event": {
				"type": "vip-discount-applied",
				"params": {
					"discount": 30,
					"message": "Réduction VIP de 30%"
				}
			}
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
			"event": {
				"type": "regular-discount-applied",
				"params": {
					"discount": 10,
					"message": "Réduction standard de 10%"
				}
			}
		}
	]`

	// JSON des facts (données)
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

	// Charger les règles
	var rules []*gorulesengine.Rule
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		fmt.Printf("❌ Erreur parsing rules: %v\n", err)
		return
	}
	fmt.Printf("✅ %d règles chargées depuis JSON\n", len(rules))

	// Charger les facts
	var factsData map[string]interface{}
	if err := json.Unmarshal([]byte(factsJSON), &factsData); err != nil {
		fmt.Printf("❌ Erreur parsing facts: %v\n", err)
		return
	}
	fmt.Println("✅ Facts chargés depuis JSON")

	// Créer l'engine et ajouter les règles
	engine := gorulesengine.NewEngine()
	for _, rule := range rules {
		engine.AddRule(rule)
	}

	// Créer l'almanac et ajouter les facts
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	for key, value := range factsData {
		almanac.AddFact(gorulesengine.FactID(key), value)
	}

	// Exécuter
	fmt.Println("\n🚀 Exécution du moteur...")
	results, err := engine.Run(almanac)
	if err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	// Afficher les résultats
	fmt.Println("📊 RÉSULTATS:")
	fmt.Printf("   Total: %d règles évaluées\n\n", len(results))

	for _, result := range results {
		if result.Result {
			fmt.Printf("✅ Règle '%s' RÉUSSIE\n", result.Rule.Name)
			fmt.Printf("   Event: %s\n", result.Event.Type)
			fmt.Printf("   Discount: %v%%\n", result.Event.Params["discount"])
			fmt.Printf("   Message: %v\n\n", result.Event.Params["message"])
		} else {
			fmt.Printf("❌ Règle '%s' ÉCHOUÉE\n\n", result.Rule.Name)
		}
	}
}

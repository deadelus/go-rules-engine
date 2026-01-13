package main

import (
	"fmt"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func main() {
	fmt.Println("🚀 Exemple Advanced - Callbacks & Dynamic Facts")
	fmt.Println("================================================\n")

	// Créer l'engine
	engine := gorulesengine.NewEngine()

	// Callback nommé pour VIP
	vipCallback := func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
		fmt.Println("   ✅ Callback VIP: Client VIP détecté!")
		return nil
	}

	// Callback pour achat important
	largeOrderCallback := func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
		amount, _ := almanac.GetFactValue("orderAmount", nil, "")
		fmt.Printf("   ✅ Callback: Gros achat de %.2f€\n", amount)
		return nil
	}

	// Enregistrer les callbacks
	engine.RegisterCallback("vipCallback", vipCallback)
	engine.RegisterCallback("largeOrderCallback", largeOrderCallback)

	// Handler global pour tous les succès
	engine.OnSuccess(func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
		fmt.Printf("🎯 Global Handler - Event '%s' déclenché\n", event.Type)
		return nil
	})

	// Règle 1: Vérifier le statut VIP avec callback
	vipCallbackName := "vipCallback"
	rule1 := &gorulesengine.Rule{
		Name:     "vip-customer",
		Priority: 100,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "customerType",
						Operator: "equal",
						Value:    "VIP",
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "vip-benefits",
		},
		OnSuccess: &vipCallbackName,
	}

	// Règle 2: Achat important avec callback
	largeOrderCallbackName := "largeOrderCallback"
	rule2 := &gorulesengine.Rule{
		Name:     "large-purchase",
		Priority: 90,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "orderAmount",
						Operator: "greater_than",
						Value:    1000,
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "high-value-order",
		},
		OnSuccess: &largeOrderCallbackName,
	}

	// Règle 3: Dynamic Fact - calculer la remise
	rule3 := &gorulesengine.Rule{
		Name:     "calculate-discount",
		Priority: 80,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "discount",
						Operator: "greater_than",
						Value:    0,
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "discount-applied",
		},
	}

	engine.AddRule(rule1)
	engine.AddRule(rule2)
	engine.AddRule(rule3)

	// Handler spécifique pour un type d'événement
	engine.On("vip-benefits", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
		fmt.Println("   🌟 Handler spécifique VIP: Avantages premium activés")
		return nil
	})

	// Créer l'almanac avec des faits statiques
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("customerType", "VIP")
	almanac.AddFact("orderAmount", 1500.0)

	// Dynamic Fact: calculer la remise selon le type et le montant
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

		fmt.Printf("   📊 Dynamic Fact calculé: remise = %.0f%%\n", discount)
		return discount, nil
	})

	fmt.Println("📋 Données:")
	fmt.Println("   Type: VIP")
	fmt.Println("   Montant: 1500€")
	fmt.Println("   Remise: (calculée dynamiquement)")

	// Exécuter
	fmt.Println("🚀 Exécution...")
	results, err := engine.Run(almanac)
	if err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	// Résumé
	fmt.Println("\n📊 RÉSUMÉ:")
	successCount := 0
	for _, result := range results {
		if result.Result {
			successCount++
		}
	}
	fmt.Printf("   %d règles réussies sur %d\n", successCount, len(results))
}

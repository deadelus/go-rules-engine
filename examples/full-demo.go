package main

import (
	"encoding/json"
	"fmt"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func main() {
	fmt.Println("🚀 GO RULES ENGINE - Démonstration complète")
	fmt.Println("=========================================================")

	// Test 1: Engine simple avec une règle
	fmt.Println("📋 Test 1: Engine avec une règle simple")
	testEngineSimple()

	// Test 2: Engine avec callbacks nommés (dans JSON)
	fmt.Println("📋 Test 2: Engine avec callbacks nommés")
	testEngineWithCallbacks()

	// Test 3: Engine avec handlers globaux
	fmt.Println("📋 Test 3: Engine avec handlers globaux")
	testEngineWithGlobalHandlers()

	// Test 4: Engine avec plusieurs règles et priorités
	fmt.Println("📋 Test 4: Engine avec plusieurs règles")
	testEngineMultipleRules()

	// Test 5: Engine avec handlers par type d'événement
	fmt.Println("📋 Test 5: Engine avec handlers par type")
	testEngineWithEventTypeHandlers()

	// Test 6: Engine avec règles et facts depuis JSON
	fmt.Println("📋 Test 6: Engine avec JSON complet (rules + facts)")
	testEngineFromJSON()

	fmt.Println("✅ Tous les tests sont terminés!")
}

func testEngineSimple() {
	// 1. Créer l'engine
	engine := gorulesengine.NewEngine()

	// 2. Créer une règle
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
			Type: "user-adult",
			Params: map[string]interface{}{
				"message": "Utilisateur adulte détecté",
			},
		},
	}

	// 3. Ajouter la règle
	engine.AddRule(rule)

	// 4. Créer l'almanac avec des faits
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("age", 25)

	// 5. Exécuter l'engine
	results, err := engine.Run(almanac)
	if err != nil {
		fmt.Printf("  ❌ Erreur: %v\n", err)
		return
	}

	// 6. Afficher les résultats
	fmt.Printf("  ✅ Nombre de règles évaluées: %d\n", len(results))
	for _, result := range results {
		status := "❌ Échec"
		if result.Result {
			status = "✅ Succès"
		}
		fmt.Printf("  %s - Règle '%s' - Event: %s\n", status, result.Rule.Name, result.Event.Type)
	}

	// 7. Consulter l'historique des événements
	successEvents := almanac.GetSuccessEvents()
	fmt.Printf("  📊 Événements success: %d\n", len(successEvents))
}

func testEngineWithCallbacks() {
	// 1. Créer l'engine
	engine := gorulesengine.NewEngine()

	// 2. Enregistrer les callbacks NOMMÉS
	engine.RegisterCallback("sendWelcomeEmail", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		fmt.Printf("  📧 Callback 'sendWelcomeEmail' appelé\n")
		fmt.Printf("     Message: %v\n", event.Params["message"])
		return nil
	})

	engine.RegisterCallback("logFailure", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		fmt.Printf("  📝 Callback 'logFailure' appelé\n")
		fmt.Printf("     Règle '%s' a échoué\n", result.Rule.Name)
		return nil
	})

	// 3. Créer une règle avec callbacks (comme si elle venait d'un JSON)
	onSuccessName := "sendWelcomeEmail"
	onFailureName := "logFailure"

	rule := &gorulesengine.Rule{
		Name:     "adult-check",
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
			Type: "user-verified",
			Params: map[string]interface{}{
				"message": "Utilisateur vérifié avec succès",
			},
		},
		OnSuccess: &onSuccessName,
		OnFailure: &onFailureName,
	}

	// 4. Ajouter et exécuter
	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("age", 25)

	results, err := engine.Run(almanac)
	if err != nil {
		fmt.Printf("  ❌ Erreur: %v\n", err)
		return
	}

	fmt.Printf("  ✅ Règle évaluée: %s - Résultat: %v\n", results[0].Rule.Name, results[0].Result)
}

func testEngineWithGlobalHandlers() {
	// 1. Créer l'engine
	engine := gorulesengine.NewEngine()

	// 2. Enregistrer des handlers GLOBAUX
	engine.OnSuccess(func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		fmt.Printf("  ✅ Handler global SUCCESS déclenché\n")
		fmt.Printf("     Event type: %s\n", event.Type)
		return nil
	})

	engine.OnFailure(func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		fmt.Printf("  ❌ Handler global FAILURE déclenché\n")
		fmt.Printf("     Règle: %s\n", result.Rule.Name)
		return nil
	})

	// 3. Créer deux règles (une qui passe, une qui échoue)
	rule1 := &gorulesengine.Rule{
		Name:     "rule-pass",
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
		Event: gorulesengine.Event{Type: "event-pass"},
	}

	rule2 := &gorulesengine.Rule{
		Name:     "rule-fail",
		Priority: 5,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "age",
						Operator: "less_than",
						Value:    18,
					},
				},
			},
		},
		Event: gorulesengine.Event{Type: "event-fail"},
	}

	engine.AddRule(rule1)
	engine.AddRule(rule2)

	// 4. Exécuter
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("age", 25)

	results, err := engine.Run(almanac)
	if err != nil {
		fmt.Printf("  ❌ Erreur: %v\n", err)
		return
	}

	fmt.Printf("  📊 Total: %d règles évaluées\n", len(results))
}

func testEngineMultipleRules() {
	// 1. Créer l'engine
	engine := gorulesengine.NewEngine()

	// 2. Créer plusieurs règles avec priorités différentes
	rules := []*gorulesengine.Rule{
		{
			Name:     "premium-user",
			Priority: 100, // Haute priorité
			Conditions: gorulesengine.ConditionSet{
				All: []gorulesengine.ConditionNode{
					{
						Condition: &gorulesengine.Condition{
							Fact:     "isPremium",
							Operator: "equal",
							Value:    true,
						},
					},
				},
			},
			Event: gorulesengine.Event{
				Type: "premium-access",
				Params: map[string]interface{}{
					"discount": 20,
					"level":    "gold",
				},
			},
		},
		{
			Name:     "adult-user",
			Priority: 50, // Priorité moyenne
			Conditions: gorulesengine.ConditionSet{
				All: []gorulesengine.ConditionNode{
					{
						Condition: &gorulesengine.Condition{
							Fact:     "age",
							Operator: "greater_than_inclusive",
							Value:    18,
						},
					},
				},
			},
			Event: gorulesengine.Event{
				Type: "adult-access",
				Params: map[string]interface{}{
					"discount": 10,
				},
			},
		},
		{
			Name:     "default-user",
			Priority: 1, // Basse priorité
			Conditions: gorulesengine.ConditionSet{
				All: []gorulesengine.ConditionNode{
					{
						Condition: &gorulesengine.Condition{
							Fact:     "age",
							Operator: "greater_than",
							Value:    0,
						},
					},
				},
			},
			Event: gorulesengine.Event{
				Type: "basic-access",
				Params: map[string]interface{}{
					"discount": 5,
				},
			},
		},
	}

	// 3. Ajouter toutes les règles
	for _, rule := range rules {
		engine.AddRule(rule)
	}

	// 4. Créer l'almanac avec des faits
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("age", 25)
	almanac.AddFact("isPremium", true)

	// 5. Exécuter
	results, err := engine.Run(almanac)
	if err != nil {
		fmt.Printf("  ❌ Erreur: %v\n", err)
		return
	}

	// 6. Afficher les résultats
	fmt.Printf("  ✅ Nombre de règles évaluées: %d\n", len(results))
	successCount := 0
	for _, result := range results {
		if result.Result {
			successCount++
			eventJSON, _ := json.MarshalIndent(result.Event, "     ", "  ")
			fmt.Printf("  ✅ Règle '%s' (priorité: %d)\n", result.Rule.Name, result.Rule.Priority)
			fmt.Printf("     Event: %s\n", string(eventJSON))
		}
	}
	fmt.Printf("  📊 Règles réussies: %d/%d\n", successCount, len(results))

	// 7. Consulter l'historique
	allEvents := almanac.GetEvents()
	fmt.Printf("  📚 Total événements dans l'historique: %d\n", len(allEvents))
}

func testEngineWithEventTypeHandlers() {
	// 1. Créer l'engine
	engine := gorulesengine.NewEngine()

	// 2. Enregistrer des handlers SPÉCIFIQUES par type d'événement
	engine.On("user-adult", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		fmt.Printf("  🎯 Handler spécifique 'user-adult' déclenché\n")
		fmt.Printf("     Discount: %v%%\n", event.Params["discount"])
		return nil
	})

	engine.On("premium-access", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		fmt.Printf("  💎 Handler spécifique 'premium-access' déclenché\n")
		fmt.Printf("     Level: %v\n", event.Params["level"])
		return nil
	})

	// 3. Créer les règles
	rule1 := &gorulesengine.Rule{
		Name:     "adult-rule",
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
			Type: "user-adult",
			Params: map[string]interface{}{
				"discount": 10,
			},
		},
	}

	rule2 := &gorulesengine.Rule{
		Name:     "premium-rule",
		Priority: 20,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "isPremium",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "premium-access",
			Params: map[string]interface{}{
				"level": "platinum",
			},
		},
	}

	engine.AddRule(rule1)
	engine.AddRule(rule2)

	// 4. Exécuter
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("age", 25)
	almanac.AddFact("isPremium", true)

	results, err := engine.Run(almanac)
	if err != nil {
		fmt.Printf("  ❌ Erreur: %v\n", err)
		return
	}

	fmt.Printf("  ✅ %d règles ont matché\n", len(results))
}

func testEngineFromJSON() {
	// 1. JSON des règles (comme ce qui viendrait d'une API ou d'un fichier)
	rulesJSON := `[
		{
			"name": "premium-discount",
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
				"type": "apply-premium-discount",
				"params": {
					"discountPercent": 25,
					"message": "Réduction premium appliquée"
				}
			},
			"onSuccess": "notifyPremiumDiscount",
			"onFailure": "logNoDiscount"
		},
		{
			"name": "regular-discount",
			"priority": 50,
			"conditions": {
				"all": [
					{
						"condition": {
							"fact": "order.total",
							"operator": "greater_than",
							"value": 50
						}
					}
				]
			},
			"event": {
				"type": "apply-regular-discount",
				"params": {
					"discountPercent": 10,
					"message": "Réduction standard appliquée"
				}
			},
			"onSuccess": "notifyRegularDiscount"
		},
		{
			"name": "first-order-bonus",
			"priority": 75,
			"conditions": {
				"all": [
					{
						"condition": {
							"fact": "user.isFirstOrder",
							"operator": "equal",
							"value": true
						}
					}
				]
			},
			"event": {
				"type": "apply-first-order-bonus",
				"params": {
					"bonusAmount": 15,
					"message": "Bonus première commande"
				}
			},
			"onSuccess": "sendWelcomeBonus"
		}
	]`

	// 2. JSON des facts (données utilisateur + commande)
	factsJSON := `{
		"user": {
			"id": 12345,
			"name": "Alice Dupont",
			"isPremium": true,
			"isFirstOrder": false,
			"email": "alice@example.com"
		},
		"order": {
			"id": "ORD-9876",
			"total": 150.50,
			"items": [
				{"name": "Produit A", "price": 50.00},
				{"name": "Produit B", "price": 100.50}
			],
			"currency": "EUR"
		}
	}`

	// 3. Unmarshall les règles
	var rules []*gorulesengine.Rule
	if err := json.Unmarshal([]byte(rulesJSON), &rules); err != nil {
		fmt.Printf("  ❌ Erreur unmarshall rules: %v\n", err)
		return
	}
	fmt.Printf("  📦 %d règles chargées depuis JSON\n", len(rules))

	// 4. Unmarshall les facts
	var factsData map[string]interface{}
	if err := json.Unmarshal([]byte(factsJSON), &factsData); err != nil {
		fmt.Printf("  ❌ Erreur unmarshall facts: %v\n", err)
		return
	}
	fmt.Printf("  📦 Facts chargés depuis JSON\n")

	// 5. Créer l'engine
	engine := gorulesengine.NewEngine()

	// 6. Enregistrer les callbacks référencés dans les règles JSON
	engine.RegisterCallback("notifyPremiumDiscount", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		fmt.Printf("  💎 CALLBACK: Premium discount de %v%% appliqué!\n", event.Params["discountPercent"])
		return nil
	})

	engine.RegisterCallback("notifyRegularDiscount", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		fmt.Printf("  🎫 CALLBACK: Discount régulier de %v%% appliqué!\n", event.Params["discountPercent"])
		return nil
	})

	engine.RegisterCallback("sendWelcomeBonus", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		fmt.Printf("  🎁 CALLBACK: Bonus première commande de %v€ offert!\n", event.Params["bonusAmount"])
		return nil
	})

	engine.RegisterCallback("logNoDiscount", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		fmt.Printf("  📝 CALLBACK: Pas de réduction premium (conditions non remplies)\n")
		return nil
	})

	// 7. Ajouter un handler global pour voir tous les succès
	engine.OnSuccess(func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		fmt.Printf("  ✅ [GLOBAL] Règle '%s' réussie - Event: %s\n", result.Rule.Name, event.Type)
		return nil
	})

	// 8. Ajouter toutes les règles
	for _, rule := range rules {
		engine.AddRule(rule)
	}

	// 9. Créer l'almanac et ajouter les facts
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter chaque fact depuis le JSON unmarshallé
	for key, value := range factsData {
		almanac.AddFact(gorulesengine.FactID(key), value)
	}

	// 10. Exécuter l'engine
	fmt.Println("\n  🚀 Exécution de l'engine...")
	results, err := engine.Run(almanac)
	if err != nil {
		fmt.Printf("  ❌ Erreur lors de l'exécution: %v\n", err)
		return
	}

	// 11. Afficher un résumé
	fmt.Println("  📊 RÉSUMÉ:")
	fmt.Printf("     Total règles évaluées: %d\n", len(results))

	successCount := 0
	for _, result := range results {
		if result.Result {
			successCount++
		}
	}
	fmt.Printf("     Règles réussies: %d\n", successCount)
	fmt.Printf("     Règles échouées: %d\n", len(results)-successCount)

	// 12. Afficher les événements générés
	successEvents := almanac.GetSuccessEvents()
	fmt.Printf("\n  📚 Événements générés: %d\n", len(successEvents))
	for i, evt := range successEvents {
		fmt.Printf("     %d. Type: %s\n", i+1, evt.Type)
		if msg, ok := evt.Params["message"]; ok {
			fmt.Printf("        Message: %v\n", msg)
		}
	}
}

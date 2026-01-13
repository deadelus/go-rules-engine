package main

import (
	"fmt"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func main() {
	fmt.Println("🚀 GO RULES ENGINE - Test Fonctionnel Global")
	fmt.Println("=" + string(make([]byte, 50)) + "\n")

	// Test 1: Règle simple avec conditions "all"
	fmt.Println("📋 Test 1: Règle avec conditions ALL")
	testSimpleRule()

	// Test 2: Règle avec conditions "any"
	fmt.Println("\n📋 Test 2: Règle avec conditions ANY")
	testAnyRule()

	// Test 3: Règle avec conditions "none"
	fmt.Println("\n📋 Test 3: Règle avec conditions NONE")
	testNoneRule()

	// Test 4: Règle complexe imbriquée
	fmt.Println("\n📋 Test 4: Règle complexe avec conditions imbriquées")
	testComplexRule()

	// Test 5: Fait dynamique avec cache
	fmt.Println("\n📋 Test 5: Fait dynamique avec cache")
	testDynamicFactWithCache()

	// Test 6: JSONPath sur structures profondes
	fmt.Println("\n📋 Test 6: JSONPath sur structures imbriquées")
	testDeepJSONPath()

	// Test 7: Tous les opérateurs
	fmt.Println("\n📋 Test 7: Test de tous les opérateurs")
	testAllOperators()

	// Test 8: Moteur complet avec plusieurs règles
	// fmt.Println("\n📋 Test 8: Moteur avec plusieurs règles et priorités")
	// testEngineWithMultipleRules()

	fmt.Println("\n✅ Tous les tests fonctionnels sont terminés!")
}

func testSimpleRule() {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("age", 25)
	almanac.AddFact("city", "Paris")

	rule := gorulesengine.Rule{
		Name:     "adult-in-paris",
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
				{
					Condition: &gorulesengine.Condition{
						Fact:     "city",
						Operator: "equal",
						Value:    "Paris",
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "user-eligible",
			Params: map[string]interface{}{
				"message": "Utilisateur adulte à Paris",
			},
		},
	}

	result, err := rule.Conditions.Evaluate(almanac)
	if err != nil {
		fmt.Printf("  ❌ Erreur: %v\n", err)
		return
	}
	fmt.Printf("  ✅ Résultat: %v (attendu: true)\n", result)
}

func testAnyRule() {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("age", 15)
	almanac.AddFact("hasPermission", true)

	rule := gorulesengine.Rule{
		Name:     "access-allowed",
		Priority: 5,
		Conditions: gorulesengine.ConditionSet{
			Any: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "age",
						Operator: "greater_than_inclusive",
						Value:    18,
					},
				},
				{
					Condition: &gorulesengine.Condition{
						Fact:     "hasPermission",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "access-granted",
		},
	}

	result, err := rule.Conditions.Evaluate(almanac)
	if err != nil {
		fmt.Printf("  ❌ Erreur: %v\n", err)
		return
	}
	fmt.Printf("  ✅ Résultat: %v (attendu: true - hasPermission=true)\n", result)
}

func testNoneRule() {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("isBanned", false)
	almanac.AddFact("isBlocked", false)

	rule := gorulesengine.Rule{
		Name:     "user-allowed",
		Priority: 8,
		Conditions: gorulesengine.ConditionSet{
			None: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "isBanned",
						Operator: "equal",
						Value:    true,
					},
				},
				{
					Condition: &gorulesengine.Condition{
						Fact:     "isBlocked",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "user-valid",
		},
	}

	result, err := rule.Conditions.Evaluate(almanac)
	if err != nil {
		fmt.Printf("  ❌ Erreur: %v\n", err)
		return
	}
	fmt.Printf("  ✅ Résultat: %v (attendu: true - aucun banned/blocked)\n", result)
}

func testComplexRule() {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("age", 25)
	almanac.AddFact("score", 85)
	almanac.AddFact("city", "Lyon")
	almanac.AddFact("hasVIP", false)

	rule := gorulesengine.Rule{
		Name:     "premium-eligibility",
		Priority: 15,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "age",
						Operator: "greater_than",
						Value:    18,
					},
				},
				{
					SubSet: &gorulesengine.ConditionSet{
						Any: []gorulesengine.ConditionNode{
							{
								Condition: &gorulesengine.Condition{
									Fact:     "score",
									Operator: "greater_than_inclusive",
									Value:    80,
								},
							},
							{
								Condition: &gorulesengine.Condition{
									Fact:     "hasVIP",
									Operator: "equal",
									Value:    true,
								},
							},
						},
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type:   "premium-access",
			Params: map[string]interface{}{"level": "gold"},
		},
	}

	result, err := rule.Conditions.Evaluate(almanac)
	if err != nil {
		fmt.Printf("  ❌ Erreur: %v\n", err)
		return
	}
	fmt.Printf("  ✅ Résultat: %v (attendu: true - age>18 ET score>=80)\n", result)
}

func testDynamicFactWithCache() {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	callCount := 0
	almanac.AddFact("db-query", func(params map[string]interface{}) int {
		callCount++
		fmt.Printf("  🔍 Appel DB n°%d - Query: %v\n", callCount, params["query"])
		return 42
	}, gorulesengine.WithCache())

	params := map[string]interface{}{"query": "SELECT * FROM users"}

	// Premier appel - calcul
	result1, _ := almanac.GetFactValue("db-query", params, "")
	fmt.Printf("  ✅ Premier appel: %v\n", result1)

	// Deuxième appel - depuis cache
	result2, _ := almanac.GetFactValue("db-query", params, "")
	fmt.Printf("  ✅ Deuxième appel (cache): %v\n", result2)

	// Troisième appel - depuis cache
	result3, _ := almanac.GetFactValue("db-query", params, "")
	fmt.Printf("  ✅ Troisième appel (cache): %v\n", result3)

	fmt.Printf("  📊 Nombre total d'appels à la fonction: %d (attendu: 1)\n", callCount)
}

func testDeepJSONPath() {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	userData := map[string]interface{}{
		"user": map[string]interface{}{
			"profile": map[string]interface{}{
				"address": map[string]interface{}{
					"city":    "Paris",
					"zipcode": "75001",
					"country": "France",
				},
				"contacts": []interface{}{
					map[string]interface{}{"type": "email", "value": "user@example.com"},
					map[string]interface{}{"type": "phone", "value": "+33612345678"},
				},
			},
			"preferences": map[string]interface{}{
				"language": "fr",
				"theme":    "dark",
			},
		},
	}

	almanac.AddFact("userData", userData)

	// Test 1: Navigation profonde
	city, _ := almanac.GetFactValue("userData", nil, "$.user.profile.address.city")
	fmt.Printf("  ✅ City: %v\n", city)

	// Test 2: Accès à un array
	email, _ := almanac.GetFactValue("userData", nil, "$.user.profile.contacts[0].value")
	fmt.Printf("  ✅ Email: %v\n", email)

	// Test 3: Wildcard pour extraire tous les types de contacts
	contactTypes, _ := almanac.GetFactValue("userData", nil, "$.user.profile.contacts[*].type")
	fmt.Printf("  ✅ Contact types: %v\n", contactTypes)

	// Test 4: Multiple niveaux
	lang, _ := almanac.GetFactValue("userData", nil, "$.user.preferences.language")
	fmt.Printf("  ✅ Language: %v\n", lang)
}

func testAllOperators() {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("age", 25)
	almanac.AddFact("score", 85.5)
	almanac.AddFact("name", "Alice")
	almanac.AddFact("tags", []interface{}{"premium", "verified"})

	tests := []struct {
		operator string
		fact     string
		value    interface{}
		expected bool
	}{
		{"equal", "age", 25, true},
		{"not_equal", "age", 30, true},
		{"greater_than", "age", 18, true},
		{"greater_than_inclusive", "age", 25, true},
		{"less_than", "age", 30, true},
		{"less_than_inclusive", "age", 25, true},
		{"in", "name", []interface{}{"Alice", "Bob"}, true},
		{"not_in", "name", []interface{}{"Charlie", "Dave"}, true},
		{"contains", "tags", "premium", true},
		{"not_contains", "tags", "blocked", true},
	}

	fmt.Printf("  %-40s | %-10s | %-18s | %s \n", "Operator", "Fact", "Value", "Résultat")
	fmt.Println("  " + string(make([]byte, 65)))

	for _, test := range tests {
		condition := gorulesengine.Condition{
			Fact:     gorulesengine.FactID(test.fact),
			Operator: gorulesengine.OperatorType(test.operator),
			Value:    test.value,
		}

		result, err := condition.Evaluate(almanac)
		status := "✅"
		if err != nil || result != test.expected {
			status = "❌"
		}

		valueStr := fmt.Sprintf("%v", test.value)
		if len(valueStr) > 18 {
			valueStr = valueStr[:18] + "..."
		}

		fmt.Printf("  %-40s | %-10s | %-18s | %s %v\n",
			test.operator, test.fact, valueStr, status, result)
	}
}

/**
func testEngineWithMultipleRules() {
	engine := gorulesengine.NewEngine([]gorulesengine.Rule{}, nil)

	// Règle 1: Haute priorité - utilisateur premium
	rule1 := gorulesengine.Rule{
		Name:     "premium-user",
		Priority: 100,
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
			Type:   "premium-access",
			Params: map[string]interface{}{"discount": 20},
		},
	}

	// Règle 2: Priorité moyenne - utilisateur adulte
	rule2 := gorulesengine.Rule{
		Name:     "adult-user",
		Priority: 50,
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
			Type:   "adult-access",
			Params: map[string]interface{}{"discount": 10},
		},
	}

	// Règle 3: Basse priorité - règle par défaut
	rule3 := gorulesengine.Rule{
		Name:     "default-user",
		Priority: 1,
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
			Type:   "basic-access",
			Params: map[string]interface{}{"discount": 5},
		},
	}

	engine.AddRule(rule1)
	engine.AddRule(rule2)
	engine.AddRule(rule3)

	// Créer un almanac avec des faits
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("age", 25)
	almanac.AddFact("isPremium", true)

	// Afficher les règles
	fmt.Println("  📜 Règles chargées:")
	gorulesengine.PrintRules([]gorulesengine.Rule{rule1, rule2, rule3})

	// Exécuter le moteur
	fmt.Println("\n  🎯 Exécution des règles...")
	results, err := engine.Run(almanac)
	if err != nil {
		fmt.Printf("  ❌ Erreur: %v\n", err)
		return
	}

	fmt.Printf("\n  ✅ Nombre de règles déclenchées: %d\n", len(results))
	for i, result := range results {
		eventJSON, _ := json.MarshalIndent(result.Event, "  ", "  ")
		fmt.Printf("  %d. %s (priorité: %d)\n", i+1, result.Name, result.Priority)
		fmt.Printf("     Event: %s\n", string(eventJSON))
	}
}
*/

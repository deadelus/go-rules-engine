package main

import (
	"fmt"
	"strings"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

// CustomOperator implémente l'interface Operator
type CustomOperator struct {
	evaluate func(interface{}, interface{}) (bool, error)
}

// Evaluate implémente la méthode Evaluate de l'interface Operator
func (c *CustomOperator) Evaluate(factValue interface{}, conditionValue interface{}) (bool, error) {
	return c.evaluate(factValue, conditionValue)
}

func main() {
	fmt.Println("🚀 Exemple Custom Operator - Opérateurs personnalisés")
	fmt.Println("========================================================")

	// Opérateur: starts_with
	startsWithOp := &CustomOperator{
		evaluate: func(factValue interface{}, conditionValue interface{}) (bool, error) {
			strValue, ok1 := factValue.(string)
			prefix, ok2 := conditionValue.(string)
			if !ok1 || !ok2 {
				return false, fmt.Errorf("starts_with requiert des strings")
			}
			return strings.HasPrefix(strValue, prefix), nil
		},
	}

	// Opérateur: ends_with
	endsWithOp := &CustomOperator{
		evaluate: func(factValue interface{}, conditionValue interface{}) (bool, error) {
			strValue, ok1 := factValue.(string)
			suffix, ok2 := conditionValue.(string)
			if !ok1 || !ok2 {
				return false, fmt.Errorf("ends_with requiert des strings")
			}
			return strings.HasSuffix(strValue, suffix), nil
		},
	}

	// Opérateur: between
	betweenOp := &CustomOperator{
		evaluate: func(factValue interface{}, conditionValue interface{}) (bool, error) {
			var numValue float64
			switch v := factValue.(type) {
			case float64:
				numValue = v
			case int:
				numValue = float64(v)
			default:
				return false, fmt.Errorf("between requiert un nombre")
			}

			rangeSlice, ok := conditionValue.([]interface{})
			if !ok || len(rangeSlice) != 2 {
				return false, fmt.Errorf("between requiert [min, max]")
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

	// Enregistrer les opérateurs
	gorulesengine.RegisterOperator("starts_with", startsWithOp)
	gorulesengine.RegisterOperator("ends_with", endsWithOp)
	gorulesengine.RegisterOperator("between", betweenOp)

	fmt.Println("✅ Opérateurs enregistrés:")
	fmt.Println("   - starts_with")
	fmt.Println("   - ends_with")
	fmt.Println("   - between")

	// Créer l'engine
	engine := gorulesengine.NewEngine()

	// Règles utilisant les opérateurs custom
	rule1 := &gorulesengine.Rule{
		Name:     "email-corporate",
		Priority: 100,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "email",
						Operator: "ends_with",
						Value:    "@company.com",
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "corporate-email",
		},
	}

	rule2 := &gorulesengine.Rule{
		Name:     "code-vip",
		Priority: 90,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "clientCode",
						Operator: "starts_with",
						Value:    "VIP-",
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "vip-client",
		},
	}

	rule3 := &gorulesengine.Rule{
		Name:     "age-range",
		Priority: 80,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "age",
						Operator: "between",
						Value:    []interface{}{25, 40},
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "target-age",
		},
	}

	engine.AddRule(rule1)
	engine.AddRule(rule2)
	engine.AddRule(rule3)

	// Données de test
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("email", "john@company.com")
	almanac.AddFact("clientCode", "VIP-12345")
	almanac.AddFact("age", 32)

	fmt.Println("📋 Données:")
	fmt.Println("   Email: john@company.com")
	fmt.Println("   Code: VIP-12345")
	fmt.Println("   Âge: 32")

	// Exécuter
	fmt.Println("🚀 Exécution...")
	results, err := engine.Run(almanac)
	if err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	// Résultats
	for _, result := range results {
		if result.Result {
			fmt.Printf("✅ Règle '%s' RÉUSSIE (Event: %s)\n", result.Rule.Name, result.Event.Type)
		}
	}
}

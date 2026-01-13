package main

import (
	"fmt"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func main() {
	fmt.Println("🚀 Exemple Basic - Vérification d'âge simple")
	fmt.Println("============================================")

	// Créer l'engine
	engine := gorulesengine.NewEngine()

	// Règle simple: vérifier si l'âge est supérieur à 18
	rule := &gorulesengine.Rule{
		Name:     "age-verification",
		Priority: 100,
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
			Type: "adult",
		},
	}

	engine.AddRule(rule)

	// Tester avec différents âges
	testAges := []int{16, 18, 21, 25}

	for _, age := range testAges {
		almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
		almanac.AddFact("age", age)

		fmt.Printf("Test avec âge: %d\n", age)
		results, err := engine.Run(almanac)
		if err != nil {
			fmt.Printf("❌ Erreur: %v\n\n", err)
			continue
		}

		if len(results) > 0 && results[0].Result {
			fmt.Printf("✅ Accès autorisé (adulte)\n\n")
		} else {
			fmt.Printf("❌ Accès refusé (mineur)\n\n")
		}
	}
}

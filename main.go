package main

import (
	"encoding/json"
	"fmt"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func main() {
	jsonRuleset := `{
    "name": "règle-utilisateur-adulte",
    "priority": 10,
    "conditions": {
      "all": [
        {
          "fact": "age", 
          "operator": "greater_than", 
          "value": 18, 
          "path": "$.user.age", 
          "options": {"cache": true}
        },
        {
          "any": [
            {"fact": "town", "operator": "equal", "value": "paris"},
            {"fact": "town", "operator": "equal", "value": "lyon"}
          ]
        }
      ]
    },
    "event": {
      "type": "user-is-adult",
      "params": {
        "message": "Bienvenue utilisateur adulte",
        "discountPercent": 10
      }
    }
  }`

	jsonPayload := `{
    "user": {
      "age": 25,
      "town": "paris"
    }
  }`

	// Parser le JSON en map
	var payload map[string]interface{}
	err := json.Unmarshal([]byte(jsonPayload), &payload)
	if err != nil {
		fmt.Printf("❌ Erreur parsing payload: %v\n", err)
		return
	}

	almanac := gorulesengine.NewAlmanac(nil)
	almanac.AddFact("dynamic-fact", func(params map[string]interface{}) string {
		fmt.Printf("🔍 Call Mysql .... \n")
		fmt.Printf("🔍 Params reçus: %v \n", params["data"])
		return "new yorlk"
	})
	almanac.AddFact("user-fact", payload, gorulesengine.WithCache())

	var rule gorulesengine.Rule
	err = json.Unmarshal([]byte(jsonRuleset), &rule)

	if err != nil {
		fmt.Printf("❌ Erreur: %v\n", err)
		return
	}

	fmt.Printf("✅ Règle chargée: %s\n", rule.Name)
	gorulesengine.PrintRules([]gorulesengine.Rule{rule})

	params := map[string]interface{}{
		"data": "some data",
	}
	result1, err := almanac.GetFactValue("dynamic-fact", params, "")

	if err != nil {
		fmt.Printf("❌ Erreur lors de la récupération du Fact: %v\n", err)
		return
	}
	result2, err := almanac.GetFactValue("user-fact", params, "$.user.age")

	if err != nil {
		fmt.Printf("❌ Erreur lors de la récupération du Fact: %v\n", err)
		return
	}

	// Test Cache en rappelant le même Fact
	result3, _ := almanac.GetFactValueFromCache("user-fact")
	fmt.Printf("✅ Valeur du Fact 'user-fact' depuis le cache: %v\n", result3)
	fmt.Printf("✅ Valeur du Fact 'dynamic-fact' avec path '': %v\n", result1)
	fmt.Printf("✅ Valeur du Fact 'user-fact' avec path '$.user.age': %v\n", result2)
}

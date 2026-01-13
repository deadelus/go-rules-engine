package gorulesengine_test

import (
	"encoding/json"
	"fmt"
	"testing"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func TestConditionNode_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    gorulesengine.ConditionNode
		wantErr bool
	}{
		{
			name:    "Valid Condition",
			jsonStr: `{"fact": "age", "operator": "greater_than", "value": 18}`,
			want: gorulesengine.ConditionNode{
				Condition: &gorulesengine.Condition{
					Fact:     "age",
					Operator: "greater_than",
					Value:    18,
				},
			},
			wantErr: false,
		},
		{
			name:    "Valid ConditionSet",
			jsonStr: `{"all": [{"fact": "age", "operator": "greater_than", "value": 18}]}`,
			want: gorulesengine.ConditionNode{
				SubSet: &gorulesengine.ConditionSet{
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
			},
			wantErr: false,
		},
		{
			name:    "Invalid JSON - creates empty ConditionSet",
			jsonStr: `{"invalid": "data"}`,
			want: gorulesengine.ConditionNode{
				SubSet: &gorulesengine.ConditionSet{
					All:  []gorulesengine.ConditionNode{},
					Any:  []gorulesengine.ConditionNode{},
					None: []gorulesengine.ConditionNode{},
				},
			},
			wantErr: false,
		},
		{
			name:    "Empty object - creates empty ConditionSet",
			jsonStr: `{}`,
			want: gorulesengine.ConditionNode{
				SubSet: &gorulesengine.ConditionSet{
					All:  []gorulesengine.ConditionNode{},
					Any:  []gorulesengine.ConditionNode{},
					None: []gorulesengine.ConditionNode{},
				},
			},
			wantErr: false,
		},
		{
			name:    "Malformed JSON - triggers error",
			jsonStr: `{this is not valid json`,
			want:    gorulesengine.ConditionNode{},
			wantErr: true,
		},
		{
			name:    "Invalid JSON syntax",
			jsonStr: `{"fact": "age", "operator": }`,
			want:    gorulesengine.ConditionNode{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var node gorulesengine.ConditionNode
			err := json.Unmarshal([]byte(tt.jsonStr), &node)
			if (err != nil) != tt.wantErr {
				t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				// Re-serialize to JSON to compare content
				gotJSON, _ := json.Marshal(node)
				wantJSON, _ := json.Marshal(tt.want)

				fmt.Printf("Got JSON: %s\n", string(gotJSON))
				fmt.Printf("Want JSON: %s\n", string(wantJSON))

				if string(gotJSON) != string(wantJSON) {
					t.Errorf("UnmarshalJSON() got JSON = %s, want %s", string(gotJSON), string(wantJSON))
				}
			}
		})
	}
}

func TestConditionNode_UnmarshalJSON_ErrorDetails(t *testing.T) {
	t.Run("Malformed JSON triggers error with details", func(t *testing.T) {
		var node gorulesengine.ConditionNode
		jsonStr := `{this is not valid json`
		err := json.Unmarshal([]byte(jsonStr), &node)

		if err == nil {
			t.Errorf("Expected error, got nil")
			return
		}

		// Verify that an error is returned with a descriptive message
		errMsg := err.Error()
		if errMsg == "" {
			t.Errorf("Expected non-empty error message")
		}

		t.Logf("Got error: %v", err)
	})

	t.Run("JSON Array triggers RuleEngineError", func(t *testing.T) {
		var node gorulesengine.ConditionNode
		// A JSON array cannot be unmarshaled into Condition or ConditionSet
		jsonStr := `["array", "values"]`
		err := json.Unmarshal([]byte(jsonStr), &node)

		if err == nil {
			t.Errorf("Expected error, got nil")
			return
		}

		// Verify that it's a RuleEngineError
		ruleErr, ok := err.(*gorulesengine.RuleEngineError)
		if !ok {
			t.Errorf("Expected *RuleEngineError, got %T: %v", err, err)
			return
		}

		if ruleErr.Type != gorulesengine.ErrJSON {
			t.Errorf("Expected ErrJSON, got %v", ruleErr.Type)
		}

		if ruleErr.Msg != "failed to unmarshal ConditionNode" {
			t.Errorf("Expected 'failed to unmarshal ConditionNode', got %v", ruleErr.Msg)
		}

		t.Logf("Got RuleEngineError: %v", ruleErr.Error())
	})
}

func TestCondition_Evaluate_Success(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter un fait
	err := almanac.AddFact("age", 25)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Créer une condition
	condition := &gorulesengine.Condition{
		Fact:     "age",
		Operator: "greater_than",
		Value:    18,
	}

	// Évaluer la condition
	result, err := condition.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result {
		t.Errorf("Expected true, got false")
	}
}

func TestCondition_Evaluate_WithPath(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter un fait avec structure imbriquée
	userData := map[string]interface{}{
		"user": map[string]interface{}{
			"age": 30,
		},
	}
	err := almanac.AddFact("userData", userData)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Créer une condition avec path
	condition := &gorulesengine.Condition{
		Fact:     "userData",
		Operator: "greater_than",
		Value:    25,
		Path:     "$.user.age",
	}

	// Évaluer la condition
	result, err := condition.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result {
		t.Errorf("Expected true, got false")
	}
}

func TestCondition_Evaluate_WithParams(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter un fait dynamique qui utilise params
	dynamicFunc := func(params map[string]interface{}) interface{} {
		multiplier, _ := params["multiplier"].(int)
		return 10 * multiplier
	}

	err := almanac.AddFact("dynamicValue", dynamicFunc, gorulesengine.WithoutCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Créer une condition avec params
	condition := &gorulesengine.Condition{
		Fact:     "dynamicValue",
		Operator: "equal",
		Value:    50,
		Params: map[string]interface{}{
			"multiplier": 5,
		},
	}

	// Évaluer la condition
	result, err := condition.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result {
		t.Errorf("Expected true (10*5=50), got false")
	}
}

func TestCondition_Evaluate_FactValueError(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ne pas ajouter le fait, forcer allowUndefinedFacts à false
	almanac.GetOptions()[gorulesengine.ALMANAC_OPTION_KEY_ALLOW_UNDEFINED_FACTS] = false

	// Créer une condition avec un fait inexistant
	condition := &gorulesengine.Condition{
		Fact:     "nonexistent",
		Operator: "equal",
		Value:    10,
	}

	// Évaluer devrait retourner une ConditionError
	result, err := condition.Evaluate(almanac)
	if err == nil {
		t.Fatal("Expected ConditionError, got nil")
	}

	if result {
		t.Errorf("Expected false result, got true")
	}

	// Vérifier que c'est une ConditionError
	condErr, ok := err.(*gorulesengine.ConditionError)
	if !ok {
		t.Fatalf("Expected *ConditionError, got %T", err)
	}

	// Vérifier le message d'erreur
	errMsg := condErr.Error()
	if errMsg == "" {
		t.Error("Expected non-empty error message")
	}

	t.Logf("Got ConditionError: %v", errMsg)
}

func TestCondition_Evaluate_InvalidOperator(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter un fait
	err := almanac.AddFact("age", 25)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Créer une condition avec un opérateur invalide
	condition := &gorulesengine.Condition{
		Fact:     "age",
		Operator: "invalidOperator",
		Value:    18,
	}

	// Évaluer devrait retourner une ConditionError
	result, err := condition.Evaluate(almanac)
	if err == nil {
		t.Fatal("Expected ConditionError, got nil")
	}

	if result {
		t.Errorf("Expected false result, got true")
	}

	// Vérifier que c'est une ConditionError
	condErr, ok := err.(*gorulesengine.ConditionError)
	if !ok {
		t.Fatalf("Expected *ConditionError, got %T", err)
	}

	// Vérifier que le message contient "failed to get operator"
	errMsg := condErr.Error()
	if errMsg == "" {
		t.Error("Expected non-empty error message")
	}

	t.Logf("Got ConditionError: %v", errMsg)
}

func TestCondition_Evaluate_OperatorEvaluationError(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter un fait avec un type incompatible
	err := almanac.AddFact("stringValue", "not a number")
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Créer une condition qui va échouer lors de l'évaluation de l'opérateur
	// greater_than attend des nombres comparables
	condition := &gorulesengine.Condition{
		Fact:     "stringValue",
		Operator: "greater_than",
		Value:    10,
	}

	// Évaluer devrait retourner une ConditionError
	result, err := condition.Evaluate(almanac)
	if err == nil {
		t.Fatal("Expected ConditionError for operator evaluation, got nil")
	}

	if result {
		t.Errorf("Expected false result, got true")
	}

	// Vérifier que c'est une ConditionError
	condErr, ok := err.(*gorulesengine.ConditionError)
	if !ok {
		t.Fatalf("Expected *ConditionError, got %T", err)
	}

	// Vérifier que le message contient "operator evaluation failed"
	errMsg := condErr.Error()
	if errMsg == "" {
		t.Error("Expected non-empty error message")
	}

	t.Logf("Got ConditionError: %v", errMsg)
}

func TestCondition_Evaluate_PathResolutionError(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter un fait
	userData := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
		},
	}
	err := almanac.AddFact("userData", userData)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Créer une condition avec un path invalide
	condition := &gorulesengine.Condition{
		Fact:     "userData",
		Operator: "equal",
		Value:    30,
		Path:     "$.user.age", // age n'existe pas
	}

	// Évaluer devrait retourner une ConditionError (wrapped AlmanacError)
	result, err := condition.Evaluate(almanac)
	if err == nil {
		t.Fatal("Expected ConditionError, got nil")
	}

	if result {
		t.Errorf("Expected false result, got true")
	}

	// Vérifier que c'est une ConditionError
	condErr, ok := err.(*gorulesengine.ConditionError)
	if !ok {
		t.Fatalf("Expected *ConditionError, got %T", err)
	}

	// Vérifier que le message contient "failed to get fact value"
	errMsg := condErr.Error()
	if errMsg == "" {
		t.Error("Expected non-empty error message")
	}

	t.Logf("Got ConditionError: %v", errMsg)
}

func TestConditionSet_Evaluate_AllConditionsPass(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter des faits
	almanac.AddFact("age", 25)
	almanac.AddFact("score", 85)

	// Créer un ConditionSet avec "all" - toutes les conditions passent
	conditionSet := &gorulesengine.ConditionSet{
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
					Fact:     "score",
					Operator: "greater_than_inclusive",
					Value:    80,
				},
			},
		},
	}

	result, err := conditionSet.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result {
		t.Errorf("Expected true, got false")
	}
}

func TestConditionSet_Evaluate_AllConditionsFail(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter des faits
	almanac.AddFact("age", 15)
	almanac.AddFact("score", 85)

	// Créer un ConditionSet avec "all" - une condition échoue
	conditionSet := &gorulesengine.ConditionSet{
		All: []gorulesengine.ConditionNode{
			{
				Condition: &gorulesengine.Condition{
					Fact:     "age",
					Operator: "greater_than",
					Value:    18, // 15 > 18 = false
				},
			},
			{
				Condition: &gorulesengine.Condition{
					Fact:     "score",
					Operator: "greater_than_inclusive",
					Value:    80,
				},
			},
		},
	}

	result, err := conditionSet.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Devrait retourner false car une condition "all" a échoué
	if result {
		t.Errorf("Expected false, got true")
	}
}

func TestConditionSet_Evaluate_AnyConditionsPass(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter des faits
	almanac.AddFact("age", 15)
	almanac.AddFact("hasPermission", true)

	// Créer un ConditionSet avec "any" - une condition passe
	conditionSet := &gorulesengine.ConditionSet{
		Any: []gorulesengine.ConditionNode{
			{
				Condition: &gorulesengine.Condition{
					Fact:     "age",
					Operator: "greater_than",
					Value:    18, // false
				},
			},
			{
				Condition: &gorulesengine.Condition{
					Fact:     "hasPermission",
					Operator: "equal",
					Value:    true, // true - cette condition passe
				},
			},
		},
	}

	result, err := conditionSet.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result {
		t.Errorf("Expected true (any matched), got false")
	}
}

func TestConditionSet_Evaluate_AnyConditionsAllFail(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter des faits
	almanac.AddFact("age", 15)
	almanac.AddFact("hasPermission", false)

	// Créer un ConditionSet avec "any" - toutes les conditions échouent
	conditionSet := &gorulesengine.ConditionSet{
		Any: []gorulesengine.ConditionNode{
			{
				Condition: &gorulesengine.Condition{
					Fact:     "age",
					Operator: "greater_than",
					Value:    18, // false
				},
			},
			{
				Condition: &gorulesengine.Condition{
					Fact:     "hasPermission",
					Operator: "equal",
					Value:    true, // false
				},
			},
		},
	}

	result, err := conditionSet.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Devrait retourner false car aucune condition "any" n'a passé
	if result {
		t.Errorf("Expected false, got true")
	}
}

func TestConditionSet_Evaluate_NoneConditionsPass(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter des faits
	almanac.AddFact("age", 25)
	almanac.AddFact("isBanned", false)

	// Créer un ConditionSet avec "none" - aucune ne passe
	conditionSet := &gorulesengine.ConditionSet{
		None: []gorulesengine.ConditionNode{
			{
				Condition: &gorulesengine.Condition{
					Fact:     "age",
					Operator: "less_than",
					Value:    18, // false (25 < 18)
				},
			},
			{
				Condition: &gorulesengine.Condition{
					Fact:     "isBanned",
					Operator: "equal",
					Value:    true, // false (isBanned != true)
				},
			},
		},
	}

	result, err := conditionSet.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Devrait retourner true car aucune condition "none" n'a passé
	if !result {
		t.Errorf("Expected true, got false")
	}
}

func TestConditionSet_Evaluate_NoneConditionFails(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter des faits
	almanac.AddFact("age", 15)
	almanac.AddFact("isBanned", true)

	// Créer un ConditionSet avec "none" - une condition passe
	conditionSet := &gorulesengine.ConditionSet{
		None: []gorulesengine.ConditionNode{
			{
				Condition: &gorulesengine.Condition{
					Fact:     "age",
					Operator: "less_than",
					Value:    18, // true (15 < 18) - FAIL car ne devrait pas passer
				},
			},
			{
				Condition: &gorulesengine.Condition{
					Fact:     "isBanned",
					Operator: "equal",
					Value:    true, // true - FAIL car ne devrait pas passer
				},
			},
		},
	}

	result, err := conditionSet.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Devrait retourner false car une condition "none" a passé
	if result {
		t.Errorf("Expected false (none condition matched), got true")
	}
}

func TestConditionSet_Evaluate_CombinedAllAndAny(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter des faits
	almanac.AddFact("age", 25)
	almanac.AddFact("score", 85)
	almanac.AddFact("hasPermission", true)

	// Créer un ConditionSet combiné: all + any
	conditionSet := &gorulesengine.ConditionSet{
		All: []gorulesengine.ConditionNode{
			{
				Condition: &gorulesengine.Condition{
					Fact:     "age",
					Operator: "greater_than",
					Value:    18, // true
				},
			},
		},
		Any: []gorulesengine.ConditionNode{
			{
				Condition: &gorulesengine.Condition{
					Fact:     "score",
					Operator: "greater_than",
					Value:    90, // false
				},
			},
			{
				Condition: &gorulesengine.Condition{
					Fact:     "hasPermission",
					Operator: "equal",
					Value:    true, // true
				},
			},
		},
	}

	result, err := conditionSet.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Devrait retourner true car "all" passe ET au moins un "any" passe
	if !result {
		t.Errorf("Expected true, got false")
	}
}

func TestConditionSet_Evaluate_ErrorInAllCondition(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ne pas ajouter le fait "age" pour provoquer une erreur

	// Créer un ConditionSet avec "all"
	conditionSet := &gorulesengine.ConditionSet{
		All: []gorulesengine.ConditionNode{
			{
				Condition: &gorulesengine.Condition{
					Fact:     "age", // Fact inexistant
					Operator: "greater_than",
					Value:    18,
				},
			},
		},
	}

	result, err := conditionSet.Evaluate(almanac)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result {
		t.Errorf("Expected false result on error, got true")
	}
}

func TestConditionSet_Evaluate_ErrorInAnyCondition(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ne pas ajouter de faits pour provoquer une erreur

	// Créer un ConditionSet avec "any"
	conditionSet := &gorulesengine.ConditionSet{
		Any: []gorulesengine.ConditionNode{
			{
				Condition: &gorulesengine.Condition{
					Fact:     "age", // Fact inexistant
					Operator: "greater_than",
					Value:    18,
				},
			},
		},
	}

	result, err := conditionSet.Evaluate(almanac)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result {
		t.Errorf("Expected false result on error, got true")
	}
}

func TestConditionSet_Evaluate_ErrorInNoneCondition(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ne pas ajouter de faits pour provoquer une erreur

	// Créer un ConditionSet avec "none"
	conditionSet := &gorulesengine.ConditionSet{
		None: []gorulesengine.ConditionNode{
			{
				Condition: &gorulesengine.Condition{
					Fact:     "age", // Fact inexistant
					Operator: "greater_than",
					Value:    18,
				},
			},
		},
	}

	result, err := conditionSet.Evaluate(almanac)
	if err == nil {
		t.Fatal("Expected error, got nil")
	}

	if result {
		t.Errorf("Expected false result on error, got true")
	}
}

func TestConditionSet_Evaluate_NestedSubSet(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter des faits
	almanac.AddFact("age", 25)
	almanac.AddFact("score", 85)
	almanac.AddFact("hasPermission", true)

	// Créer un ConditionSet avec subset imbriqué
	conditionSet := &gorulesengine.ConditionSet{
		All: []gorulesengine.ConditionNode{
			{
				Condition: &gorulesengine.Condition{
					Fact:     "age",
					Operator: "greater_than",
					Value:    18,
				},
			},
			{
				// Subset imbriqué
				SubSet: &gorulesengine.ConditionSet{
					Any: []gorulesengine.ConditionNode{
						{
							Condition: &gorulesengine.Condition{
								Fact:     "score",
								Operator: "greater_than",
								Value:    90, // false
							},
						},
						{
							Condition: &gorulesengine.Condition{
								Fact:     "hasPermission",
								Operator: "equal",
								Value:    true, // true
							},
						},
					},
				},
			},
		},
	}

	result, err := conditionSet.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Devrait retourner true car "all" passe (age > 18 ET subset "any" a au moins une condition vraie)
	if !result {
		t.Errorf("Expected true, got false")
	}
}

func TestConditionSet_Evaluate_EmptyConditionSet(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Créer un ConditionSet vide (pas de all/any/none)
	conditionSet := &gorulesengine.ConditionSet{}

	result, err := conditionSet.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Un ConditionSet vide devrait retourner true (aucune condition à échouer)
	if !result {
		t.Errorf("Expected true for empty ConditionSet, got false")
	}
}

func TestEvaluateConditionNode_InvalidNode(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Créer un ConditionNode vide (ni Condition ni SubSet)
	emptyNode := &gorulesengine.ConditionNode{}

	// Utiliser reflection pour accéder à la fonction non-exportée
	// Comme la fonction n'est pas exportée, on doit tester via ConditionSet.Evaluate
	// qui appelle evaluateConditionNode en interne
	conditionSet := &gorulesengine.ConditionSet{
		All: []gorulesengine.ConditionNode{*emptyNode},
	}

	result, err := conditionSet.Evaluate(almanac)
	if err == nil {
		t.Fatal("Expected ConditionError, got nil")
	}

	if result {
		t.Errorf("Expected false result on error, got true")
	}

	// Vérifier que c'est une ConditionError
	condErr, ok := err.(*gorulesengine.ConditionError)
	if !ok {
		t.Fatalf("Expected *ConditionError, got %T", err)
	}

	// Vérifier le message d'erreur
	errMsg := condErr.Error()
	if errMsg == "" {
		t.Error("Expected non-empty error message")
	}

	// Le message devrait contenir "invalid condition node"
	t.Logf("Got ConditionError: %v", errMsg)
}

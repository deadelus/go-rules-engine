package gorulesengine_test

import (
	"encoding/json"
	"fmt"
	"strings"
	"testing"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func TestConditionNode_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		jsonStr string
		want    gre.ConditionNode
		wantErr bool
	}{
		{
			name:    "Valid Condition",
			jsonStr: `{"fact": "age", "operator": "greater_than", "value": 18}`,
			want: gre.ConditionNode{
				Condition: &gre.Condition{
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
			want: gre.ConditionNode{
				SubSet: &gre.ConditionSet{
					All: []gre.ConditionNode{
						{
							Condition: &gre.Condition{
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
			want: gre.ConditionNode{
				SubSet: &gre.ConditionSet{
					All:  []gre.ConditionNode{},
					Any:  []gre.ConditionNode{},
					None: []gre.ConditionNode{},
				},
			},
			wantErr: false,
		},
		{
			name:    "Empty object - creates empty ConditionSet",
			jsonStr: `{}`,
			want: gre.ConditionNode{
				SubSet: &gre.ConditionSet{
					All:  []gre.ConditionNode{},
					Any:  []gre.ConditionNode{},
					None: []gre.ConditionNode{},
				},
			},
			wantErr: false,
		},
		{
			name:    "Malformed JSON - triggers error",
			jsonStr: `{this is not valid json`,
			want:    gre.ConditionNode{},
			wantErr: true,
		},
		{
			name:    "Invalid JSON syntax",
			jsonStr: `{"fact": "age", "operator": }`,
			want:    gre.ConditionNode{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var node gre.ConditionNode
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
	t.Run("JSON Array triggers RuleEngineError", func(t *testing.T) {
		var node gre.ConditionNode
		// A JSON array cannot be unmarshaled into Condition or ConditionSet
		jsonStr := `["array", "values"]`
		err := json.Unmarshal([]byte(jsonStr), &node)

		if err == nil {
			t.Errorf("Expected error, got nil")
			return
		}

		// Verify that it's a RuleEngineError
		ruleErr, ok := err.(*gre.RuleEngineError)
		if !ok {
			t.Errorf("Expected *RuleEngineError, got %T: %v", err, err)
			return
		}

		if ruleErr.Type != gre.ErrJSON {
			t.Errorf("Expected ErrJSON, got %v", ruleErr.Type)
		}

		if ruleErr.Msg != "failed to unmarshal ConditionNode" {
			t.Errorf("Expected 'failed to unmarshal ConditionNode', got %v", ruleErr.Msg)
		}

		// Check full error message for details
		fullMsg := ruleErr.Error()
		expectedData := `data: ["array", "values"]`
		if !strings.Contains(fullMsg, expectedData) {
			t.Errorf("Expected error message to contain data snippet %q, but got %q", expectedData, fullMsg)
		}

		expectedErr := "json: cannot unmarshal array"
		if !strings.Contains(fullMsg, expectedErr) {
			t.Errorf("Expected error message to contain %q, but got %q", expectedErr, fullMsg)
		}
	})
}

func TestCondition_Evaluate_Success(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add a fact
	err := almanac.AddFact("age", 25)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Create a condition
	condition := &gre.Condition{
		Fact:     "age",
		Operator: "greater_than",
		Value:    18,
	}

	// Evaluate the condition
	result, err := condition.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Result {
		t.Errorf("Expected true, got false")
	}
}

func TestCondition_Evaluate_WithPath(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add a fact with nested structure
	userData := map[string]interface{}{
		"user": map[string]interface{}{
			"age": 30,
		},
	}
	err := almanac.AddFact("userData", userData)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Create a condition with path
	condition := &gre.Condition{
		Fact:     "userData",
		Operator: "greater_than",
		Value:    25,
		Path:     "$.user.age",
	}

	// Evaluate the condition
	result, err := condition.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Result {
		t.Errorf("Expected true, got false")
	}
}

func TestCondition_Evaluate_WithParams(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add a dynamic fact that uses params
	dynamicFunc := func(params map[string]interface{}) interface{} {
		multiplier, _ := params["multiplier"].(int)
		return 10 * multiplier
	}

	err := almanac.AddFact("dynamicValue", dynamicFunc, gre.WithoutCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Create a condition with params
	condition := &gre.Condition{
		Fact:     "dynamicValue",
		Operator: "equal",
		Value:    50,
		Params: map[string]interface{}{
			"multiplier": 5,
		},
	}

	// Evaluate the condition
	result, err := condition.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !result.Result {
		t.Errorf("Expected true (10*5=50), got false")
	}
}

func TestCondition_Evaluate_FactValueError(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Do not add the fact, force allowUndefinedFacts to false
	almanac.GetOptions()[gre.AlmanacOptionKeyAllowUndefinedFacts] = false

	// Create a condition with a non-existent fact
	condition := &gre.Condition{
		Fact:     "nonexistent",
		Operator: "equal",
		Value:    10,
	}

	// Evaluation should return a ConditionError
	result, err := condition.Evaluate(almanac)
	if err == nil {
		t.Fatal("Expected ConditionError, got nil")
	}

	if result != nil && result.Result {
		t.Errorf("Expected false result, got true")
	}

	// Verify that it is a ConditionError
	condErr, ok := err.(*gre.ConditionError)
	if !ok {
		t.Fatalf("Expected *ConditionError, got %T", err)
	}

	// Verify the error message
	errMsg := condErr.Error()
	expectedSubstr := "failed to get fact value"
	if !strings.Contains(errMsg, expectedSubstr) {
		t.Errorf("Expected error message to contain %q, but got %q", expectedSubstr, errMsg)
	}

	almanacSubstr := "fact 'nonexistent' is not defined in the almanac"
	if !strings.Contains(errMsg, almanacSubstr) {
		t.Errorf("Expected error message to contain info about missing fact %q, but got %q", almanacSubstr, errMsg)
	}
}

func TestCondition_Evaluate_InvalidOperator(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add a fact
	err := almanac.AddFact("age", 25)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Create a condition with an invalid operator
	condition := &gre.Condition{
		Fact:     "age",
		Operator: "invalidOperator",
		Value:    18,
	}

	// Evaluation should return a ConditionError
	result, err := condition.Evaluate(almanac)
	if err == nil {
		t.Fatal("Expected ConditionError, got nil")
	}

	if result != nil && result.Result {
		t.Errorf("Expected false result, got true")
	}

	// Verify that it is a ConditionError
	condErr, ok := err.(*gre.ConditionError)
	if !ok {
		t.Fatalf("Expected *ConditionError, got %T", err)
	}

	// Verify that the message contains "failed to get operator"
	errMsg := condErr.Error()
	expectedSubstr := "failed to get operator"
	if !strings.Contains(errMsg, expectedSubstr) {
		t.Errorf("Expected error message to contain %q, but got %q", expectedSubstr, errMsg)
	}

	operatorSubstr := "operator not registered"
	if !strings.Contains(errMsg, operatorSubstr) {
		t.Errorf("Expected error message to contain %q, but got %q", operatorSubstr, errMsg)
	}
}

func TestCondition_Evaluate_OperatorEvaluationError(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add a fact with an incompatible type
	err := almanac.AddFact("stringValue", "not a number")
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Create a condition that will fail during operator evaluation
	// greater_than expects comparable numbers
	condition := &gre.Condition{
		Fact:     "stringValue",
		Operator: "greater_than",
		Value:    10,
	}

	// Evaluation should return a ConditionError
	result, err := condition.Evaluate(almanac)
	if err == nil {
		t.Fatal("Expected ConditionError for operator evaluation, got nil")
	}

	if result != nil && result.Result {
		t.Errorf("Expected false result, got true")
	}

	// Verify that it is a ConditionError
	condErr, ok := err.(*gre.ConditionError)
	if !ok {
		t.Fatalf("Expected *ConditionError, got %T", err)
	}

	// Verify that the message contains "operator evaluation failed"
	errMsg := condErr.Error()
	expectedSubstr := "operator evaluation failed"
	if !strings.Contains(errMsg, expectedSubstr) {
		t.Errorf("Expected error message to contain %q, but got %q", expectedSubstr, errMsg)
	}

	typeSubstr := "greater_than operator requires numeric values"
	if !strings.Contains(errMsg, typeSubstr) {
		t.Errorf("Expected error message to contain info about type mismatch %q, but got %q", typeSubstr, errMsg)
	}
}

func TestCondition_Evaluate_PathResolutionError(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add a fact
	userData := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
		},
	}
	err := almanac.AddFact("userData", userData)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Create a condition with an invalid path
	condition := &gre.Condition{
		Fact:     "userData",
		Operator: "equal",
		Value:    30,
		Path:     "$.user.age", // age n'existe pas
	}

	// Evaluation should return a ConditionError (wrapped AlmanacError)
	result, err := condition.Evaluate(almanac)
	if err == nil {
		t.Fatal("Expected ConditionError, got nil")
	}

	if result != nil && result.Result {
		t.Errorf("Expected false result, got true")
	}

	// Verify that it is a ConditionError
	condErr, ok := err.(*gre.ConditionError)
	if !ok {
		t.Fatalf("Expected *ConditionError, got %T", err)
	}

	// Verify that the message contains "failed to get fact value"
	errMsg := condErr.Error()
	expectedSubstr := "failed to get fact value"
	if !strings.Contains(errMsg, expectedSubstr) {
		t.Errorf("Expected error message to contain %q, but got %q", expectedSubstr, errMsg)
	}

	pathSubstr := "age not found in object"
	if !strings.Contains(errMsg, pathSubstr) {
		t.Errorf("Expected error message to contain info about path error %q, but got %q", pathSubstr, errMsg)
	}
}

func TestConditionSet_Evaluate_AllConditionsPass(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add facts
	almanac.AddFact("age", 25)
	almanac.AddFact("score", 85)

	// Create a ConditionSet with "all" - all conditions pass
	conditionSet := &gre.ConditionSet{
		All: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
					Fact:     "age",
					Operator: "greater_than",
					Value:    18,
				},
			},
			{
				Condition: &gre.Condition{
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

	if !result.Result {
		t.Errorf("Expected true, got false")
	}
}

func TestConditionSet_Evaluate_AllConditionsFail(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add facts
	almanac.AddFact("age", 15)
	almanac.AddFact("score", 85)

	// Create a ConditionSet with "all" - one condition fails
	conditionSet := &gre.ConditionSet{
		All: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
					Fact:     "age",
					Operator: "greater_than",
					Value:    18, // 15 > 18 = false
				},
			},
			{
				Condition: &gre.Condition{
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

	// Should return false because an "all" condition failed
	if result != nil && result.Result {
		t.Errorf("Expected false, got true")
	}
}

func TestConditionSet_Evaluate_AnyConditionsPass(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add facts
	almanac.AddFact("age", 15)
	almanac.AddFact("hasPermission", true)

	// Create a ConditionSet with "any" - one condition passes
	conditionSet := &gre.ConditionSet{
		Any: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
					Fact:     "age",
					Operator: "greater_than",
					Value:    18, // false
				},
			},
			{
				Condition: &gre.Condition{
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

	if !result.Result {
		t.Errorf("Expected true (any matched), got false")
	}
}

func TestConditionSet_Evaluate_AnyConditionsAllFail(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add facts
	almanac.AddFact("age", 15)
	almanac.AddFact("hasPermission", false)

	// Create a ConditionSet with "any" - all conditions fail
	conditionSet := &gre.ConditionSet{
		Any: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
					Fact:     "age",
					Operator: "greater_than",
					Value:    18, // false
				},
			},
			{
				Condition: &gre.Condition{
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

	// Should return false because no "any" condition passed
	if result != nil && result.Result {
		t.Errorf("Expected false, got true")
	}
}

func TestConditionSet_Evaluate_NoneConditionsPass(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add facts
	almanac.AddFact("age", 25)
	almanac.AddFact("isBanned", false)

	// Create a ConditionSet with "none" - none pass
	conditionSet := &gre.ConditionSet{
		None: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
					Fact:     "age",
					Operator: "less_than",
					Value:    18, // false (25 < 18)
				},
			},
			{
				Condition: &gre.Condition{
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

	// Should return true because no "none" condition passed
	if !result.Result {
		t.Errorf("Expected true, got false")
	}
}

func TestConditionSet_Evaluate_NoneConditionFails(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add facts
	almanac.AddFact("age", 15)
	almanac.AddFact("isBanned", true)

	// Create a ConditionSet with "none" - one condition passes
	conditionSet := &gre.ConditionSet{
		None: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
					Fact:     "age",
					Operator: "less_than",
					Value:    18, // true (15 < 18) - FAIL because it should not pass
				},
			},
			{
				Condition: &gre.Condition{
					Fact:     "isBanned",
					Operator: "equal",
					Value:    true, // true - FAIL because it should not pass
				},
			},
		},
	}

	result, err := conditionSet.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Should return false because a "none" condition passed
	if result != nil && result.Result {
		t.Errorf("Expected false (none condition matched), got true")
	}
}

func TestConditionSet_Evaluate_CombinedAllAndAny(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add facts
	almanac.AddFact("age", 25)
	almanac.AddFact("score", 85)
	almanac.AddFact("hasPermission", true)

	// Create a combined ConditionSet: all + any
	conditionSet := &gre.ConditionSet{
		All: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
					Fact:     "age",
					Operator: "greater_than",
					Value:    18, // true
				},
			},
		},
		Any: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
					Fact:     "score",
					Operator: "greater_than",
					Value:    90, // false
				},
			},
			{
				Condition: &gre.Condition{
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

	// Should return true because "all" passes AND at least one "any" passes
	if !result.Result {
		t.Errorf("Expected true, got false")
	}
}

func TestConditionSet_Evaluate_ErrorInAllCondition(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Do not add the "age" fact to trigger an error

	// Create a ConditionSet with "all"
	conditionSet := &gre.ConditionSet{
		All: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
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

	if result != nil && result.Result {
		t.Errorf("Expected false result on error, got true")
	}
}

func TestConditionSet_Evaluate_ErrorInAnyCondition(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Do not add facts to trigger an error

	// Create a ConditionSet with "any"
	conditionSet := &gre.ConditionSet{
		Any: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
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

	if result != nil && result.Result {
		t.Errorf("Expected false result on error, got true")
	}
}

func TestConditionSet_Evaluate_ErrorInNoneCondition(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Do not add facts to trigger an error

	// Create a ConditionSet with "none"
	conditionSet := &gre.ConditionSet{
		None: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
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

	if result != nil && result.Result {
		t.Errorf("Expected false result on error, got true")
	}
}

func TestConditionSet_Evaluate_NestedSubSet(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add facts
	almanac.AddFact("age", 25)
	almanac.AddFact("score", 85)
	almanac.AddFact("hasPermission", true)

	// Create a ConditionSet with nested subset
	conditionSet := &gre.ConditionSet{
		All: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
					Fact:     "age",
					Operator: "greater_than",
					Value:    18,
				},
			},
			{
				// Nested subset
				SubSet: &gre.ConditionSet{
					Any: []gre.ConditionNode{
						{
							Condition: &gre.Condition{
								Fact:     "score",
								Operator: "greater_than",
								Value:    90, // false
							},
						},
						{
							Condition: &gre.Condition{
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

	// Should return true because "all" passes (age > 18 AND "any" subset has at least one true condition)
	if !result.Result {
		t.Errorf("Expected true, got false")
	}
}

func TestConditionSet_Evaluate_EmptyConditionSet(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Create an empty ConditionSet (no all/any/none)
	conditionSet := &gre.ConditionSet{}

	result, err := conditionSet.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// An empty ConditionSet should return true (no condition to fail)
	if !result.Result {
		t.Errorf("Expected true for empty ConditionSet, got false")
	}
}

func TestEvaluateConditionNode_InvalidNode(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Create an empty ConditionNode (neither Condition nor SubSet)
	emptyNode := &gre.ConditionNode{}

	// Use reflection to access the non-exported function
	// Since the function is not exported, we must test via ConditionSet.Evaluate
	// which calls evaluateConditionNode internally
	conditionSet := &gre.ConditionSet{
		All: []gre.ConditionNode{*emptyNode},
	}

	result, err := conditionSet.Evaluate(almanac)
	if err == nil {
		t.Fatal("Expected ConditionError, got nil")
	}

	if result != nil && result.Result {
		t.Errorf("Expected false result on error, got true")
	}

	// Verify that it is a ConditionError
	condErr, ok := err.(*gre.ConditionError)
	if !ok {
		t.Fatalf("Expected *ConditionError, got %T", err)
	}

	// Verify the error message
	errMsg := condErr.Error()
	if errMsg == "" {
		t.Error("Expected non-empty error message")
	}

	// The message should contain "neither condition nor subset is defined"
	expectedSubstr := "neither condition nor subset is defined"
	if !strings.Contains(errMsg, expectedSubstr) {
		t.Errorf("Expected error message to contain %q, but got %q", expectedSubstr, errMsg)
	}
}

func TestConditionSetEvaluateErrors(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Error during evaluateConditionNode in All
	cFail := &gre.Condition{Fact: "f1", Operator: "invalid"} // Operator error
	csAll := &gre.ConditionSet{
		All: []gre.ConditionNode{{Condition: cFail}},
	}
	_, err := csAll.Evaluate(almanac)
	if err == nil {
		t.Error("Expected error in All evaluation")
	}

	// Error during evaluateConditionNode in Any
	csAny := &gre.ConditionSet{
		Any: []gre.ConditionNode{{Condition: cFail}},
	}
	_, err = csAny.Evaluate(almanac)
	if err == nil {
		t.Error("Expected error in Any evaluation")
	}

	// Error during evaluateConditionNode in None
	csNone := &gre.ConditionSet{
		None: []gre.ConditionNode{{Condition: cFail}},
	}
	_, err = csNone.Evaluate(almanac)
	if err == nil {
		t.Error("Expected error in None evaluation")
	}
}

func TestEvaluateConditionNode_SubSetError(t *testing.T) {
	// 1. Create an Almanac that does not allow undefined facts to trigger an error easily
	almanac := gre.NewAlmanac()

	// 2. Create a ConditionSet that will fail (missing fact)
	failingSubSet := &gre.ConditionSet{
		All: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
					Fact:     "fact_inexistant",
					Operator: "equal",
					Value:    1,
				},
			},
		},
	}

	// 3. Create the parent ConditionSet containing the faulty subset
	// We use ConditionSet.Evaluate because evaluateConditionNode is private.
	parentSet := &gre.ConditionSet{
		All: []gre.ConditionNode{
			{
				SubSet: failingSubSet,
			},
		},
	}

	// 4. Evaluate
	_, err := parentSet.Evaluate(almanac)

	// 5. Verify that the error is indeed the expected one (the wrap of evaluateConditionNode)
	if err == nil {
		t.Fatal("Une erreur était attendue lors de l'évaluation du subset")
	}

	expectedMsg := "failed to evaluate condition subset node"
	if fmt.Sprintf("%v", err) == "" || !contains(err.Error(), expectedMsg) {
		t.Errorf("Message d'erreur attendu contenant '%s', obtenu: %v", expectedMsg, err)
	}
}

// Simple helper to verify error content
func contains(s, substr string) bool {
	return fmt.Sprintf("%s", s) != "" && len(s) >= len(substr) && (s == substr || len(s) > len(substr))
}

func TestConditionSet_Evaluate_SubSetResults(t *testing.T) {
	almanac := gre.NewAlmanac()
	almanac.AddFact("age", 25)

	// A success subset: age > 18
	successSubSet := &gre.ConditionSet{
		All: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
					Fact:     "age",
					Operator: "greater_than",
					Value:    18,
				},
			},
		},
	}

	// A failing subset: age < 18
	failingSubSet := &gre.ConditionSet{
		All: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
					Fact:     "age",
					Operator: "less_than",
					Value:    18,
				},
			},
		},
	}

	t.Run("All with SubSet", func(t *testing.T) {
		cs := &gre.ConditionSet{
			All: []gre.ConditionNode{
				{SubSet: successSubSet},
			},
		}
		res, err := cs.Evaluate(almanac)
		if err != nil || !res.Result {
			t.Errorf("Expected All with successSubSet to pass, got %v, err: %v", res.Result, err)
		}
	})

	t.Run("Any with SubSet", func(t *testing.T) {
		cs := &gre.ConditionSet{
			Any: []gre.ConditionNode{
				{SubSet: failingSubSet},
				{SubSet: successSubSet}, // This one will trigger res = nodeRes.ConditionSet.Result
			},
		}
		res, err := cs.Evaluate(almanac)
		if err != nil || !res.Result {
			t.Errorf("Expected Any with successSubSet to pass, got %v, err: %v", res.Result, err)
		}
	})

	t.Run("None with SubSet", func(t *testing.T) {
		cs := &gre.ConditionSet{
			None: []gre.ConditionNode{
				{SubSet: failingSubSet}, // Doit passer (None failing == true)
			},
		}
		res, err := cs.Evaluate(almanac)
		if err != nil || !res.Result {
			t.Errorf("Expected None with failingSubSet to pass, got %v, err: %v", res.Result, err)
		}
	})
}

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

package gorulesengine

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestEngine_GenerateResponse(t *testing.T) {
	// Setup engine with audit trace
	engine := NewEngine(WithAuditTrace())

	// Define a rule
	rule := &Rule{
		Name:     "check-age",
		Priority: 100,
		Conditions: ConditionSet{
			All: []ConditionNode{
				{
					Condition: &Condition{
						Fact:     "age",
						Operator: OperatorGreaterThanInclusive,
						Value:    18,
					},
				},
			},
		},
		OnSuccess: []RuleEvent{
			{Name: "is-adult", Params: map[string]interface{}{"status": "verified"}},
		},
		OnFailure: []RuleEvent{
			{Name: "is-minor", Params: map[string]interface{}{"status": "rejected"}},
		},
	}
	engine.AddRule(rule)

	t.Run("Decision Authorize", func(t *testing.T) {
		almanac := NewAlmanac()
		almanac.AddFact("age", 25)

		e, err := engine.Run(almanac)
		if err != nil {
			t.Fatalf("Failed to run engine: %v", err)
		}

		response := e.GenerateResponse()

		if response.Decision != DecisionAuthorize {
			t.Errorf("Expected decision %s, got %s", DecisionAuthorize, response.Decision)
		}

		if len(response.Events) != 1 || response.Events[0].Type != "is-adult" {
			t.Errorf("Expected event 'is-adult', got %v", response.Events)
		}

		if response.Reason == nil {
			t.Error("Expected reason to be populated when AuditTrace is enabled")
		}

		// Verify JSON marshalable
		_, err = json.Marshal(response)
		if err != nil {
			t.Errorf("Response should be marshalable to JSON: %v", err)
		}
	})

	t.Run("Decision Decline", func(t *testing.T) {
		almanac := NewAlmanac()
		almanac.AddFact("age", 15)

		e, err := engine.Run(almanac)
		if err != nil {
			t.Fatalf("Failed to run engine: %v", err)
		}

		response := e.GenerateResponse()

		if response.Decision != DecisionDecline {
			t.Errorf("Expected decision %s, got %s", DecisionDecline, response.Decision)
		}

		if len(response.Events) != 1 || response.Events[0].Type != "is-minor" {
			t.Errorf("Expected event 'is-minor', got %v", response.Events)
		}

		// Verify JSON marshalable
		_, err = json.Marshal(response)
		if err != nil {
			t.Errorf("Response should be marshalable to JSON: %v", err)
		}
	})
}

func TestEngine_GenerateResponse_JSONComparison(t *testing.T) {
	// 1. Initialisation simple sans AuditTrace pour avoir une raison textuelle simple
	engine := NewEngine()

	rule := &Rule{
		Name:     "simple-rule",
		Priority: 10,
		Conditions: ConditionSet{
			All: []ConditionNode{
				{Condition: &Condition{Fact: "active", Operator: OperatorEqual, Value: true}},
			},
		},
		OnSuccess: []RuleEvent{
			{Name: "notify", Params: map[string]interface{}{"msg": "hello"}},
		},
	}
	engine.AddRule(rule)

	almanac := NewAlmanac()
	almanac.AddFact("active", true)

	e, _ := engine.Run(almanac)
	response := e.GenerateResponse()

	// 2. Le JSON que nous attendons (écrit à la main)
	expectedJSON := `{
		"decision": "authorize",
		"reason": "Rule 'simple-rule' determined the result",
		"events": [
			{
				"type": "notify",
				"params": {
					"msg": "hello"
				}
			}
		],
		"metadata": {}
	}`

	// 3. Marshalling du résultat de l'engine
	actualJSON, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Failed to marshal response: %v", err)
	}

	// 4. Comparaison (on dé-marshalle les deux dans des maps pour ignorer l'ordre des clés et les espaces)
	var expectedMap, actualMap map[string]interface{}
	if err := json.Unmarshal([]byte(expectedJSON), &expectedMap); err != nil {
		t.Fatalf("Failed to unmarshal expected JSON: %v", err)
	}
	if err := json.Unmarshal(actualJSON, &actualMap); err != nil {
		t.Fatalf("Failed to unmarshal actual JSON: %v", err)
	}

	if !reflect.DeepEqual(expectedMap, actualMap) {
		t.Errorf("JSON mismatch!\nExpected: %v\nActual: %v", expectedMap, actualMap)
	}
}

func TestEngine_GenerateResponse_WithMetadata(t *testing.T) {
	engine := NewEngine()

	engine.AddRule(&Rule{
		Name:     "metadata-rule",
		Priority: 1,
		Conditions: ConditionSet{
			All: []ConditionNode{
				{Condition: &Condition{Fact: "user_status", Operator: OperatorEqual, Value: "active"}},
			},
		},
	})

	almanac := NewAlmanac()
	// Add fact with metadata
	almanac.AddFact("user_status", "active", WithMetadata(map[string]interface{}{"source": "db", "verified": true}))
	almanac.AddFact("other_fact", 42, WithMetadata(map[string]interface{}{"note": "meaning of life"}))

	e, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Failed to run engine: %v", err)
	}

	response := e.GenerateResponse()

	if response.Metadata == nil {
		t.Fatal("Expected Metadata in response to be populated")
	}

	// Verify user_status metadata
	userStatusMeta, ok := response.Metadata["user_status"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected metadata for user_status to be a map, got %T", response.Metadata["user_status"])
	}
	if userStatusMeta["source"] != "db" || userStatusMeta["verified"] != true {
		t.Errorf("Unexpected metadata for user_status: %v", userStatusMeta)
	}

	// Verify other_fact metadata
	otherFactMeta, ok := response.Metadata["other_fact"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected metadata for other_fact to be a map, got %T", response.Metadata["other_fact"])
	}
	if otherFactMeta["note"] != "meaning of life" {
		t.Errorf("Unexpected metadata for other_fact: %v", otherFactMeta)
	}
}

func TestEngine_GenerateResponse_NoResults(t *testing.T) {
	engine := NewEngine()

	// Call GenerateResponse without running the engine first
	response := engine.GenerateResponse()

	if response.Decision != DecisionDecline {
		t.Errorf("Expected decision %s for empty results, got %s", DecisionDecline, response.Decision)
	}

	if len(response.Events) != 0 {
		t.Errorf("Expected 0 events, got %d", len(response.Events))
	}

	if len(response.Metadata) != 0 {
		t.Errorf("Expected 0 metadata entries, got %d", len(response.Metadata))
	}
}

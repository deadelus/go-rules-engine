package gorulesengine

import (
	"encoding/json"
	"testing"
)

func TestRuleEvent_UnmarshalJSON(t *testing.T) {
	t.Run("unmarshals from string", func(t *testing.T) {
		jsonData := []byte(`"test-event"`)
		var re RuleEvent
		err := json.Unmarshal(jsonData, &re)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}
		if re.Name != "test-event" {
			t.Errorf("Expected name 'test-event', got '%s'", re.Name)
		}
		if re.Params != nil {
			t.Error("Expected nil params")
		}
	})

	t.Run("unmarshals from object", func(t *testing.T) {
		jsonData := []byte(`{"name": "test-event", "params": {"foo": "bar"}}`)
		var re RuleEvent
		err := json.Unmarshal(jsonData, &re)
		if err != nil {
			t.Fatalf("Failed to unmarshal: %v", err)
		}
		if re.Name != "test-event" {
			t.Errorf("Expected name 'test-event', got '%s'", re.Name)
		}
		if re.Params["foo"] != "bar" {
			t.Errorf("Expected param foo='bar', got '%v'", re.Params["foo"])
		}
	})

	t.Run("returns error on invalid data", func(t *testing.T) {
		jsonData := []byte(`123`) // Not a string, not an object
		var re RuleEvent
		err := json.Unmarshal(jsonData, &re)
		if err == nil {
			t.Error("Expected error for invalid JSON type (number), got nil")
		}
	})
}

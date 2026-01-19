package gorulesengine

import (
	"encoding/json"
)

// RuleEvent represents an event reference within a rule, optionally with parameters.
type RuleEvent struct {
	Name   string                 `json:"name"`
	Params map[string]interface{} `json:"params,omitempty"`
}

// UnmarshalJSON implements custom JSON unmarshaling for RuleEvent.
// It supports both a simple string (event name) or a full object with parameters.
func (re *RuleEvent) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as a string first (backward compatibility)
	var name string
	if err := json.Unmarshal(data, &name); err == nil {
		re.Name = name
		return nil
	}

	// Otherwise, unmarshal as an object
	type Alias RuleEvent
	var aux Alias
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	re.Name = aux.Name
	re.Params = aux.Params
	return nil
}

// Rule represents a business rule with conditions and an associated event.
// Rules are evaluated against facts in an Almanac. When all conditions are met,
// the rule's event is triggered and any registered callbacks are invoked.
//
// Example:
//
//	rule := &gre.Rule{
//	    Name:     "adult-user",
//	    Priority: 10,
//	    Conditions: gre.ConditionSet{
//	        All: []gre.ConditionNode{
//	            {Condition: &gre.Condition{
//	                Fact:     "age",
//	                Operator: "greater_than",
//	                Value:    18,
//	            }},
//	        },
//	    },
//	    OnSuccess: []RuleEvent{
//	        {Name: "send-welcome-email", Params: map[string]interface{}{"template": "welcome"}},
//	    },
//	}
type Rule struct {
	Name       string       `json:"name,omitempty"`
	Priority   int          `json:"priority,omitempty"` // Higher priority rules are evaluated first
	Conditions ConditionSet `json:"conditions"`
	OnSuccess  []RuleEvent  `json:"onSuccess,omitempty"` // Events to invoke on success
	OnFailure  []RuleEvent  `json:"onFailure,omitempty"` // Events to invoke on failure
	Result     bool
}

// GetRequiredFacts returns the list of all facts required by this rule.
func (r *Rule) GetRequiredFacts() []FactID {
	return r.Conditions.GetRequiredFacts()
}

// Compile pre-calculates and optimizes rule's properties.
func (r *Rule) Compile() error {
	return r.Conditions.Compile()
}

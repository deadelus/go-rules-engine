package gorulesengine

// Rule represents a business rule with conditions and an associated event.
// Rules are evaluated against facts in an Almanac. When all conditions are met,
// the rule's event is triggered and any registered callbacks are invoked.
//
// Example:
//
//	rule := &gorulesengine.Rule{
//	    Name:     "adult-user",
//	    Priority: 10,
//	    Conditions: gorulesengine.ConditionSet{
//	        All: []gorulesengine.ConditionNode{
//	            {Condition: &gorulesengine.Condition{
//	                Fact:     "age",
//	                Operator: "greater_than",
//	                Value:    18,
//	            }},
//	        },
//	    },
//	    Event: gorulesengine.Event{Type: "user-is-adult"},
//	}
type Rule struct {
	Name       string       `json:"name,omitempty"`
	Priority   int          `json:"priority,omitempty"` // Higher priority rules are evaluated first
	Conditions ConditionSet `json:"conditions"`
	Event      Event        `json:"event"`
	OnSuccess  *string      `json:"on_success,omitempty"` // Name of callback to invoke on success
	OnFailure  *string      `json:"on_failure,omitempty"` // Name of callback to invoke on failure
}

// RuleResult represents the result of evaluating a rule.
// It contains the rule that was evaluated, the event that was triggered,
// and whether the rule matched (Result = true) or not (Result = false).
type RuleResult struct {
	Event  Event
	Rule   *Rule
	Result bool // true if the rule matched, false otherwise
}

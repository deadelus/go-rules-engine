package gorulesengine

// Rule represents a rule with conditions and an associated event
type Rule struct {
	Name       string       `json:"name,omitempty"`
	Priority   int          `json:"priority,omitempty"`
	Conditions ConditionSet `json:"conditions"`
	Event      Event        `json:"event"`
	OnSuccess  *string      `json:"on_success,omitempty"`
	OnFailure  *string      `json:"on_failure,omitempty"`
}

// RuleResult represents the result of evaluating a rule
type RuleResult struct {
	Event  Event
	Rule   *Rule
	Result bool // true si la règle a matché
}

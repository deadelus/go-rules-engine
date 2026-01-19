// package gorulesengine provides a powerful and flexible rules engine for Go.
// It allows you to define business rules in JSON or code, evaluate complex conditions,
// and trigger events based on dynamic facts.
package gorulesengine

import "time"

// ConditionType represents the type of logical operator for combining conditions.
type ConditionType string

const (
	// AllType represents a logical AND - all conditions must be true.
	AllType ConditionType = "all"

	// AnyType represents a logical OR - at least one condition must be true.
	AnyType ConditionType = "any"

	// NoneType represents a logical NOT - no conditions must be true.
	NoneType ConditionType = "none"
)

// OperatorType represents the type of comparison operator used in conditions.
type OperatorType string

// OperatorEqual checks if the fact value equals the condition value.
const (
	OperatorEqual OperatorType = "equal"

	// OperatorNotEqual checks if the fact value is not equal to the condition value.
	OperatorNotEqual OperatorType = "not_equal"

	// OperatorLessThan checks if the fact value is less than the condition value.
	OperatorLessThan OperatorType = "less_than"

	// OperatorLessThanInclusive checks if the fact value is less than or equal to the condition value.
	OperatorLessThanInclusive OperatorType = "less_than_inclusive"

	// OperatorGreaterThan checks if the fact value is greater than the condition value.
	OperatorGreaterThan OperatorType = "greater_than"

	// OperatorGreaterThanInclusive checks if the fact value is greater than or equal to the condition value.
	OperatorGreaterThanInclusive OperatorType = "greater_than_inclusive"

	// OperatorIn checks if the fact value is contained in the condition value (array).
	OperatorIn OperatorType = "in"

	// OperatorNotIn checks if the fact value is not contained in the condition value (array).
	OperatorNotIn OperatorType = "not_in"

	// OperatorContains checks if the fact value contains the condition value (for strings and arrays).
	OperatorContains OperatorType = "contains"

	// OperatorNotContains checks if the fact value does not contain the condition value.
	OperatorNotContains OperatorType = "not_contains"

	// OperatorRegex checks if the fact value matches the regex pattern in the condition value.
	OperatorRegex OperatorType = "regex"
)

// MetricsCollector defines an interface for monitoring the rules engine's performance and execution results.
// Implementations can use Prometheus, OpenTelemetry, or other monitoring systems.
type MetricsCollector interface {
	// ObserveRuleEvaluation records the result and duration of a single rule evaluation.
	ObserveRuleEvaluation(ruleName string, result bool, duration time.Duration)
	// ObserveEngineRun records the total duration and rule count of an engine execution.
	ObserveEngineRun(ruleCount int, duration time.Duration)
	// ObserveEventExecution records the duration and result of an event handler execution.
	ObserveEventExecution(eventName string, ruleName string, result bool, duration time.Duration)
}

// RuleResult represents the complete evaluation result of a single rule.
type RuleResult struct {
	Name       string              `json:"name"`
	Priority   int                 `json:"priority"`
	Result     bool                `json:"result"`
	Conditions *ConditionSetResult `json:"conditions"`
	OnSuccess  []RuleEvent         `json:"onSuccess,omitempty"`
	OnFailure  []RuleEvent         `json:"onFailure,omitempty"`
}

// ConditionSetResult represents the evaluation result of a ConditionSet (All, Any, or None).
type ConditionSetResult struct {
	Type    ConditionType         `json:"type"`
	Result  bool                  `json:"result"`
	Results []ConditionNodeResult `json:"results"`
}

// ConditionNodeResult represents the result of a single node within a ConditionSet.
type ConditionNodeResult struct {
	Condition    *ConditionResult    `json:"condition,omitempty"`
	ConditionSet *ConditionSetResult `json:"conditionSet,omitempty"`
}

// ConditionResult represents the detailed evaluation result of a single Condition.
type ConditionResult struct {
	Fact      FactID       `json:"fact"`
	Operator  OperatorType `json:"operator"`
	Value     interface{}  `json:"value"`          // The value to compare against
	FactValue interface{}  `json:"factValue"`      // The actual value fetched from the Almanac
	Path      string       `json:"path,omitempty"` // The JSONPath used, if any
	Result    bool         `json:"result"`
	Error     string       `json:"error,omitempty"`
}

const (
	// DecisionAuthorize indicates that conditions were met.
	DecisionAuthorize = "authorize"
	// DecisionDecline indicates that conditions were not met.
	DecisionDecline = "decline"
)

// EngineResponse represents the final formatted structure for your JSON response.
type EngineResponse struct {
	Decision string                 `json:"decision"` // DecisionAuthorize or DecisionDecline
	Reason   interface{}            `json:"reason"`   // Detail of conditions (if AuditTrace is active)
	Events   []EventResponse        `json:"events"`   // List of triggered events
	Metadata map[string]interface{} `json:"metadata"` // Metadata from facts or other sources
}

// EventResponse represents a simplified event.
type EventResponse struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params,omitempty"`
}

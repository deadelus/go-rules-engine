// Package gorulesengine provides a powerful and flexible rules engine for Go.
// It allows you to define business rules in JSON or code, evaluate complex conditions,
// and trigger events based on dynamic facts.
package gorulesengine

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

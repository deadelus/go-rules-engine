// Package gorulesengine provides a powerful and flexible rules engine for Go.
// It allows you to define business rules in JSON or code, evaluate complex conditions,
// and trigger events based on dynamic facts.
package gorulesengine

// ConditionType represents the type of logical operator for combining conditions.
type ConditionType string

// All represents a logical AND - all conditions must be true.
const All ConditionType = "all"

// Any represents a logical OR - at least one condition must be true.
const Any ConditionType = "any"

// None represents a logical NOT - no conditions must be true.
const None ConditionType = "none"

// OperatorType represents the type of comparison operator used in conditions.
type OperatorType string

// Equal checks if the fact value equals the condition value.
const Equal OperatorType = "equal"

// NotEqual checks if the fact value is not equal to the condition value.
const NotEqual OperatorType = "not_equal"

// LessThan checks if the fact value is less than the condition value.
const LessThan OperatorType = "less_than"

// LessThanInclusive checks if the fact value is less than or equal to the condition value.
const LessThanInclusive OperatorType = "less_than_inclusive"

// GreaterThan checks if the fact value is greater than the condition value.
const GreaterThan OperatorType = "greater_than"

// GreaterThanInclusive checks if the fact value is greater than or equal to the condition value.
const GreaterThanInclusive OperatorType = "greater_than_inclusive"

// In checks if the fact value is contained in the condition value (array).
const In OperatorType = "in"

// NotIn checks if the fact value is not contained in the condition value (array).
const NotIn OperatorType = "not_in"

// Contains checks if the fact value contains the condition value (for strings and arrays).
const Contains OperatorType = "contains"

// NotContains checks if the fact value does not contain the condition value.
const NotContains OperatorType = "not_contains"

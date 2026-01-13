package gorulesengine

import "fmt"

type ErrorType string

const (
	ErrAlmanac   ErrorType = "ALMANAC_ERROR"
	ErrFact      ErrorType = "FACT_ERROR"
	ErrRule      ErrorType = "RULE_ERROR"
	ErrCondition ErrorType = "CONDITION_ERROR"
	ErrOperator  ErrorType = "OPERATOR_ERROR"
	ErrEvent     ErrorType = "EVENT_ERROR"
	ErrJSON      ErrorType = "JSON_ERROR"
)

// RuleEngineError is the base error type for the rule engine
type RuleEngineError struct {
	Type ErrorType
	Msg  string
	Err  error // wrapped error (optional)
}

// Error implements the error interface
func (e *RuleEngineError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("[%s] %s: %v", e.Type, e.Msg, e.Err)
	}
	return fmt.Sprintf("[%s] %s", e.Type, e.Msg)
}

// Unwrap returns the wrapped error
func (e *RuleEngineError) Unwrap() error {
	return e.Err
}

type AlmanacError struct {
	Payload string
	Err     error
}

type FactError struct {
	Fact Fact
	Err  error
}

// OperatorError represents an error related to a specific operator
type OperatorError struct {
	Operator     OperatorType
	Value        interface{}
	CompareValue interface{}
	Err          error
}

// RuleError represents an error related to a specific rule
type RuleError struct {
	Rule Rule
	Err  error
}

// ConditionError represents an error related to a specific condition
type ConditionError struct {
	Condition Condition
	Err       error
}

// Error methods to convert to RuleEngineError
func (e *AlmanacError) Error() string {
	return (&RuleEngineError{
		Type: ErrAlmanac,
		Msg: fmt.Sprintf(
			"almanac=%v",
			e.Payload,
		),
		Err: e.Err,
	}).Error()
}

// Unwrap returns the wrapped error
func (e *AlmanacError) Unwrap() error {
	return e.Err
}

// Error methods to convert to RuleEngineError
func (e *OperatorError) Error() string {
	return (&RuleEngineError{
		Type: ErrOperator,
		Msg: fmt.Sprintf(
			"operator=%s value=%v compareValue=%v",
			e.Operator,
			e.Value,
			e.CompareValue,
		),
		Err: e.Err,
	}).Error()
}

// Unwrap returns the wrapped error
func (e *OperatorError) Unwrap() error {
	return e.Err
}

// Error methods to convert to RuleEngineError
func (e *RuleError) Error() string {
	return (&RuleEngineError{
		Type: ErrRule,
		Msg: fmt.Sprintf(
			"rule=%v",
			e.Rule,
		),
		Err: e.Err,
	}).Error()
}

// Unwrap returns the wrapped error
func (e *RuleError) Unwrap() error {
	return e.Err
}

// Error methods to convert to RuleEngineError
func (e *ConditionError) Error() string {
	return (&RuleEngineError{
		Type: ErrCondition,
		Msg: fmt.Sprintf(
			"condition=%v",
			e.Condition,
		),
		Err: e.Err,
	}).Error()
}

// Unwrap returns the wrapped error
func (e *ConditionError) Unwrap() error {
	return e.Err
}

// Error methods to convert to RuleEngineError
func (e *FactError) Error() string {
	return (&RuleEngineError{
		Type: ErrFact,
		Msg: fmt.Sprintf(
			"fact=%v",
			e.Fact,
		),
		Err: e.Err,
	}).Error()
}

// Unwrap returns the wrapped error
func (e *FactError) Unwrap() error {
	return e.Err
}

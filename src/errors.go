package gorulesengine

import "fmt"

// ErrorType identifies the category of error that occurred.
type ErrorType string

const (
	// ErrEngine indicates a general engine execution error.
	ErrEngine ErrorType = "ENGINE_ERROR"
	// ErrAlmanac indicates an error related to the almanac or fact management.
	ErrAlmanac ErrorType = "ALMANAC_ERROR"
	// ErrFact indicates an error computing or accessing a fact value.
	ErrFact ErrorType = "FACT_ERROR"
	// ErrRule indicates an error in rule definition or structure.
	ErrRule ErrorType = "RULE_ERROR"
	// ErrCondition indicates an error evaluating a condition.
	ErrCondition ErrorType = "CONDITION_ERROR"
	// ErrOperator indicates an error with an operator (not found, invalid, etc.).
	ErrOperator ErrorType = "OPERATOR_ERROR"
	// ErrEvent indicates an error related to event handling.
	ErrEvent ErrorType = "EVENT_ERROR"
	// ErrJSON indicates an error parsing or unmarshaling JSON.
	ErrJSON ErrorType = "JSON_ERROR"
	// ErrLoader indicates an error related to loading rules or data.
	ErrLoader ErrorType = "LOADER_ERROR"
)

// RuleEngineError is the base error type for all errors in the rule engine.
// It categorizes errors by type and optionally wraps underlying errors.
type RuleEngineError struct {
	Type ErrorType // The category of error
	Msg  string    // Human-readable error message
	Err  error     // Wrapped underlying error (optional)
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

// AlmanacError represents an error that occurred while accessing or managing facts in the almanac.
type AlmanacError struct {
	Payload string // Context about what was being accessed
	Err     error  // Underlying error
}

// FactError represents an error that occurred while computing or accessing a fact value.
type FactError struct {
	Fact Fact  // The fact that caused the error
	Err  error // Underlying error
}

// OperatorError represents an error related to a specific operator evaluation.
type OperatorError struct {
	Operator     OperatorType // The operator that failed
	Value        interface{}  // The fact value being compared
	CompareValue interface{}  // The expected value
	Err          error        // Underlying error
}

// RuleError represents an error related to a specific rule evaluation or definition.
type RuleError struct {
	Rule Rule  // The rule that caused the error
	Err  error // Underlying error
}

// ConditionError represents an error that occurred while evaluating a condition.
type ConditionError struct {
	Condition Condition // The condition that failed
	Err       error     // Underlying error
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

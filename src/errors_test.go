package gorulesengine_test

import (
	"errors"
	"testing"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func TestRuleEngineError_Error_WithWrappedError(t *testing.T) {
	wrappedErr := errors.New("wrapped error")
	err := &gorulesengine.RuleEngineError{
		Type: gorulesengine.ErrEngine,
		Msg:  "test error message",
		Err:  wrappedErr,
	}

	expected := "[ENGINE_ERROR] test error message: wrapped error"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestRuleEngineError_Error_WithoutWrappedError(t *testing.T) {
	err := &gorulesengine.RuleEngineError{
		Type: gorulesengine.ErrEngine,
		Msg:  "test error message",
		Err:  nil,
	}

	expected := "[ENGINE_ERROR] test error message"
	if err.Error() != expected {
		t.Errorf("Expected error message '%s', got '%s'", expected, err.Error())
	}
}

func TestRuleEngineError_Unwrap(t *testing.T) {
	wrappedErr := errors.New("wrapped error")
	err := &gorulesengine.RuleEngineError{
		Type: gorulesengine.ErrEngine,
		Msg:  "test error",
		Err:  wrappedErr,
	}

	unwrapped := err.Unwrap()
	if unwrapped != wrappedErr {
		t.Errorf("Expected unwrapped error to be '%v', got '%v'", wrappedErr, unwrapped)
	}
}

func TestAlmanacError_Error(t *testing.T) {
	wrappedErr := errors.New("almanac wrapped error")
	err := &gorulesengine.AlmanacError{
		Payload: "test-payload",
		Err:     wrappedErr,
	}

	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("Expected non-empty error message")
	}

	// Should contain ALMANAC_ERROR type
	if len(errorMsg) < 10 {
		t.Errorf("Error message too short: %s", errorMsg)
	}
}

func TestAlmanacError_Unwrap(t *testing.T) {
	wrappedErr := errors.New("almanac wrapped error")
	err := &gorulesengine.AlmanacError{
		Payload: "test-payload",
		Err:     wrappedErr,
	}

	unwrapped := err.Unwrap()
	if unwrapped != wrappedErr {
		t.Errorf("Expected unwrapped error to be '%v', got '%v'", wrappedErr, unwrapped)
	}
}

func TestOperatorError_Error(t *testing.T) {
	wrappedErr := errors.New("operator error")
	err := &gorulesengine.OperatorError{
		Operator:     "equal",
		Value:        10,
		CompareValue: 20,
		Err:          wrappedErr,
	}

	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("Expected non-empty error message")
	}

	// Should contain operator type
	if len(errorMsg) < 10 {
		t.Errorf("Error message too short: %s", errorMsg)
	}
}

func TestOperatorError_Unwrap(t *testing.T) {
	wrappedErr := errors.New("operator error")
	err := &gorulesengine.OperatorError{
		Operator:     "equal",
		Value:        10,
		CompareValue: 20,
		Err:          wrappedErr,
	}

	unwrapped := err.Unwrap()
	if unwrapped != wrappedErr {
		t.Errorf("Expected unwrapped error to be '%v', got '%v'", wrappedErr, unwrapped)
	}
}

func TestRuleError_Error(t *testing.T) {
	wrappedErr := errors.New("rule error")
	rule := gorulesengine.Rule{
		Name:     "test-rule",
		Priority: 10,
	}

	err := &gorulesengine.RuleError{
		Rule: rule,
		Err:  wrappedErr,
	}

	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("Expected non-empty error message")
	}

	// Should contain RULE_ERROR type
	if len(errorMsg) < 10 {
		t.Errorf("Error message too short: %s", errorMsg)
	}
}

func TestRuleError_Unwrap(t *testing.T) {
	wrappedErr := errors.New("rule error")
	rule := gorulesengine.Rule{
		Name:     "test-rule",
		Priority: 10,
	}

	err := &gorulesengine.RuleError{
		Rule: rule,
		Err:  wrappedErr,
	}

	unwrapped := err.Unwrap()
	if unwrapped != wrappedErr {
		t.Errorf("Expected unwrapped error to be '%v', got '%v'", wrappedErr, unwrapped)
	}
}

func TestConditionError_Error(t *testing.T) {
	wrappedErr := errors.New("condition error")
	condition := gorulesengine.Condition{
		Fact:     "age",
		Operator: "greater_than",
		Value:    18,
	}

	err := &gorulesengine.ConditionError{
		Condition: condition,
		Err:       wrappedErr,
	}

	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("Expected non-empty error message")
	}

	// Should contain CONDITION_ERROR type
	if len(errorMsg) < 10 {
		t.Errorf("Error message too short: %s", errorMsg)
	}
}

func TestConditionError_Unwrap(t *testing.T) {
	wrappedErr := errors.New("condition error")
	condition := gorulesengine.Condition{
		Fact:     "age",
		Operator: "greater_than",
		Value:    18,
	}

	err := &gorulesengine.ConditionError{
		Condition: condition,
		Err:       wrappedErr,
	}

	unwrapped := err.Unwrap()
	if unwrapped != wrappedErr {
		t.Errorf("Expected unwrapped error to be '%v', got '%v'", wrappedErr, unwrapped)
	}
}

func TestFactError_Error(t *testing.T) {
	wrappedErr := errors.New("fact error")
	fact := gorulesengine.NewFact("testFact", 42)

	err := &gorulesengine.FactError{
		Fact: *fact,
		Err:  wrappedErr,
	}

	errorMsg := err.Error()
	if errorMsg == "" {
		t.Error("Expected non-empty error message")
	}

	// Should contain FACT_ERROR type
	if len(errorMsg) < 10 {
		t.Errorf("Error message too short: %s", errorMsg)
	}
}

func TestFactError_Unwrap(t *testing.T) {
	wrappedErr := errors.New("fact error")
	fact := gorulesengine.NewFact("testFact", 42)

	err := &gorulesengine.FactError{
		Fact: *fact,
		Err:  wrappedErr,
	}

	unwrapped := err.Unwrap()
	if unwrapped != wrappedErr {
		t.Errorf("Expected unwrapped error to be '%v', got '%v'", wrappedErr, unwrapped)
	}
}

func TestErrorTypes_Constants(t *testing.T) {
	tests := []struct {
		errorType gorulesengine.ErrorType
		expected  string
	}{
		{gorulesengine.ErrEngine, "ENGINE_ERROR"},
		{gorulesengine.ErrAlmanac, "ALMANAC_ERROR"},
		{gorulesengine.ErrFact, "FACT_ERROR"},
		{gorulesengine.ErrRule, "RULE_ERROR"},
		{gorulesengine.ErrCondition, "CONDITION_ERROR"},
		{gorulesengine.ErrOperator, "OPERATOR_ERROR"},
		{gorulesengine.ErrEvent, "EVENT_ERROR"},
		{gorulesengine.ErrJSON, "JSON_ERROR"},
	}

	for _, tt := range tests {
		if string(tt.errorType) != tt.expected {
			t.Errorf("Expected error type '%s', got '%s'", tt.expected, string(tt.errorType))
		}
	}
}

func TestErrors_WithErrorsIs(t *testing.T) {
	baseErr := errors.New("base error")

	// Test with RuleEngineError
	engineErr := &gorulesengine.RuleEngineError{
		Type: gorulesengine.ErrEngine,
		Msg:  "test",
		Err:  baseErr,
	}

	if !errors.Is(engineErr, baseErr) {
		t.Error("errors.Is should find wrapped error in RuleEngineError")
	}

	// Test with AlmanacError
	almanacErr := &gorulesengine.AlmanacError{
		Payload: "test",
		Err:     baseErr,
	}

	if !errors.Is(almanacErr, baseErr) {
		t.Error("errors.Is should find wrapped error in AlmanacError")
	}
}

func TestErrors_WithErrorsAs(t *testing.T) {
	// Create a RuleEngineError
	originalErr := &gorulesengine.RuleEngineError{
		Type: gorulesengine.ErrEngine,
		Msg:  "test error",
		Err:  errors.New("base"),
	}

	// Test errors.As
	var targetErr *gorulesengine.RuleEngineError
	if !errors.As(originalErr, &targetErr) {
		t.Error("errors.As should work with RuleEngineError")
	}

	if targetErr.Type != gorulesengine.ErrEngine {
		t.Errorf("Expected error type %s, got %s", gorulesengine.ErrEngine, targetErr.Type)
	}
}

package gorulesengine_test

import (
	"fmt"
	"testing"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

// Test for EqualOperator with numbers
func TestEqualOperator_EvaluateNumbers(t *testing.T) {
	var value = 5
	var assertedValue int = value
	var failedValue int = 10
	var operator gorulesengine.Operator = &gorulesengine.EqualOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for EqualOperator with strings
func TestEqualOperator_EvaluateStrings(t *testing.T) {
	var value = "test"
	var assertedValue string = value
	var failedValue string = "fail"
	var operator gorulesengine.Operator = &gorulesengine.EqualOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for EqualOperator with uncomparable types
func TestEqualOperator_EvaluateUncomparable(t *testing.T) {
	var value = []int{1, 2, 3}
	var failedValue = func() {}
	var operator gorulesengine.Operator = &gorulesengine.EqualOperator{}

	result, _ := operator.Evaluate(value, failedValue)

	if result != false {
		t.Errorf("Expected false for uncomparable types, got true")
	}
}

// Test for EqualOperator with nil values
func TestEqualOperator_EvaluateNilValues(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.EqualOperator{}

	// Test with nil factValue
	_, err := operator.Evaluate(nil, "test")
	if err == nil {
		t.Errorf("Expected error for nil factValue, got nil")
	}

	// Test with nil compareValue
	_, err = operator.Evaluate("test", nil)
	if err == nil {
		t.Errorf("Expected error for nil compareValue, got nil")
	}

	// Test with both nil
	_, err = operator.Evaluate(nil, nil)
	if err == nil {
		t.Errorf("Expected error for both nil values, got nil")
	}
}

// Test for NotEqualOperator with numbers
func TestNotEqualOperator_EvaluateNumbers(t *testing.T) {
	var value = 5
	var assertedValue int = 10
	var failedValue int = 5
	var operator gorulesengine.Operator = &gorulesengine.NotEqualOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for NotEqualOperator with strings
func TestNotEqualOperator_EvaluateStrings(t *testing.T) {
	var value = "test"
	var assertedValue string = "fail"
	var failedValue string = "test"
	var operator gorulesengine.Operator = &gorulesengine.NotEqualOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for NotEqualOperator with uncomparable types
func TestNotEqualOperator_EvaluateUncomparable(t *testing.T) {
	var value = []int{1, 2, 3}
	var failedValue1 int = 25
	var failedValue2 = []int{4, 5, 6}
	var operator gorulesengine.Operator = &gorulesengine.NotEqualOperator{}

	result, _ := operator.Evaluate(value, failedValue1)

	if result != true {
		t.Errorf("Expected true for uncomparable types, got false")
	}

	result, _ = operator.Evaluate(value, failedValue2)
	if result != true {
		t.Errorf("Expected true for different slices, got false")
	}
}

// Test for LessThanOperator with integers
func TestLessThanOperator_EvaluateIntegers(t *testing.T) {
	var value = 5
	var assertedValue int = 6
	var failedValue int = 5
	var operator gorulesengine.Operator = &gorulesengine.LessThanOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for LessThanOperator with floats
func TestLessThanOperator_EvaluateFloats(t *testing.T) {
	var value = 5.0
	var assertedValue float64 = 6.0
	var failedValue float64 = 5.0
	var operator gorulesengine.Operator = &gorulesengine.LessThanOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for LessThanOperator with non-numeric numbers
func TestLessThanOperator_EvaluateNonNumericValues(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.LessThanOperator{}

	// Test with string
	_, err := operator.Evaluate("test", 5)
	if err == nil {
		t.Errorf("Expected error for non-numeric factValue, got nil")
	}

	// Test with non-numeric compareValue
	_, err = operator.Evaluate(5, "test")
	if err == nil {
		t.Errorf("Expected error for non-numeric compareValue, got nil")
	}
}

// Test for LessThanInclusiveOperator with integers
func TestLessThanInclusiveOperator_EvaluateIntegers(t *testing.T) {
	var value = 5
	var assertedValue int = value
	var failedValue int = 4
	var operator gorulesengine.Operator = &gorulesengine.LessThanInclusiveOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for LessThanInclusiveOperator with floats
func TestLessThanInclusiveOperator_EvaluateFloats(t *testing.T) {
	var value = 5.0
	var assertedValue float64 = value
	var failedValue float64 = 4.0
	var operator gorulesengine.Operator = &gorulesengine.LessThanInclusiveOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for LessThanInclusiveOperator with non-numeric numbers
func TestLessThanInclusiveOperator_EvaluateNonNumericValues(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.LessThanInclusiveOperator{}

	_, err := operator.Evaluate("test", 5)
	if err == nil {
		t.Errorf("Expected error for non-numeric values, got nil")
	}
}

// Test for GreaterThanOperator with integers
func TestGreaterThanOperator_EvaluateIntegers(t *testing.T) {
	var value = 5
	var assertedValue int = 4
	var failedValue int = 5
	var operator gorulesengine.Operator = &gorulesengine.GreaterThanOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for GreaterThanOperator with floats
func TestGreaterThanOperator_EvaluateFloats(t *testing.T) {
	var value = 5.0
	var assertedValue float64 = 4.0
	var failedValue float64 = 5.0
	var operator gorulesengine.Operator = &gorulesengine.GreaterThanOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for GreaterThanOperator with non-numeric numbers
func TestGreaterThanOperator_EvaluateNonNumericValues(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.GreaterThanOperator{}

	_, err := operator.Evaluate("test", 5)
	if err == nil {
		t.Errorf("Expected error for non-numeric values, got nil")
	}
}

// Test for GreaterThanInclusiveOperator with integers
func TestGreaterThanInclusiveOperator_EvaluateIntegers(t *testing.T) {
	var value = 5
	var assertedValue int = 5
	var failedValue int = 6
	var operator gorulesengine.Operator = &gorulesengine.GreaterThanInclusiveOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for GreaterThanInclusiveOperator with floats
func TestGreaterThanInclusiveOperator_EvaluateFloats(t *testing.T) {
	var value = 5.0
	var assertedValue float64 = 5.0
	var failedValue float64 = 6.0
	var operator gorulesengine.Operator = &gorulesengine.GreaterThanInclusiveOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for GreaterThanInclusiveOperator with non-numeric numbers
func TestGreaterThanInclusiveOperator_EvaluateNonNumericValues(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.GreaterThanInclusiveOperator{}

	_, err := operator.Evaluate("test", 5)
	if err == nil {
		t.Errorf("Expected error for non-numeric values, got nil")
	}
}

// Test for InOperator with numbers
func TestInOperator_EvaluateNumbers(t *testing.T) {
	var value int = 5
	var assertedValue = []int{1, 2, 3, 4, 5}
	var failedValue = []int{6, 7, 8, 9, 10}
	var operator gorulesengine.Operator = &gorulesengine.InOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for InOperator with strings
func TestInOperator_EvaluateStrings(t *testing.T) {
	var value string = "apple"
	var assertedValue = []string{"banana", "orange", "apple", "grape"}
	var failedValue = []string{"pear", "melon", "kiwi"}
	var operator gorulesengine.Operator = &gorulesengine.InOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for InOperator with invalid type
func TestInOperator_EvaluateInvalidType(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.InOperator{}

	// Test with compareValue that is not a slice
	_, err := operator.Evaluate(5, "not a slice")
	if err == nil {
		t.Errorf("Expected error for non-slice compareValue, got nil")
	}

	// Test with a number
	_, err = operator.Evaluate(5, 10)
	if err == nil {
		t.Errorf("Expected error for non-slice compareValue, got nil")
	}
}

// Test for NotInOperator with numbers
func TestNotInOperator_EvaluateNumbers(t *testing.T) {
	var value int = 6
	var assertedValue = []int{1, 2, 3, 4, 5}
	var failedValue = []int{6, 7, 8, 9, 10}
	var operator gorulesengine.Operator = &gorulesengine.NotInOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for NotInOperator with strings
func TestNotInOperator_EvaluateStrings(t *testing.T) {
	var value string = "apple"
	var assertedValue = []string{"pear", "melon", "kiwi"}
	var failedValue = []string{"banana", "orange", "apple", "grape"}
	var operator gorulesengine.Operator = &gorulesengine.NotInOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for ContainsOperator with numbers
func TestContainsOperator_EvaluateNumbers(t *testing.T) {
	var value = []int{1, 2, 3, 4, 5}
	var assertedValue int = 5
	var failedValue int = 6
	var operator gorulesengine.Operator = &gorulesengine.ContainsOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for ContainsOperator with strings
func TestContainsOperator_EvaluateStrings(t *testing.T) {
	var value = "An apple a day keeps the doctor away"
	var assertedValue string = "apple"
	var failedValue string = "pear"
	var operator gorulesengine.Operator = &gorulesengine.ContainsOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for ContainsOperator with invalid types
func TestContainsOperator_EvaluateInvalidTypes(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.ContainsOperator{}

	// Test with factValue that is neither slice nor string
	_, err := operator.Evaluate(123, "test")
	if err == nil {
		t.Errorf("Expected error for invalid factValue type, got nil")
	}

	// Test with string factValue mais compareValue non-string
	_, err = operator.Evaluate("test string", 123)
	if err == nil {
		t.Errorf("Expected error for non-string compareValue when factValue is string, got nil")
	}
}

// Test for NotContainsOperator with numbers
func TestNotContainsOperator_EvaluateNumbers(t *testing.T) {
	var value = []int{1, 2, 3, 4, 5}
	var assertedValue int = 6
	var failedValue int = 5
	var operator gorulesengine.Operator = &gorulesengine.NotContainsOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Test for NotContainsOperator with strings
func TestNotContainsOperator_EvaluateStrings(t *testing.T) {
	var value = "An apple a day keeps the doctor away"
	var assertedValue string = "pear"
	var failedValue string = "apple"
	var operator gorulesengine.Operator = &gorulesengine.NotContainsOperator{}

	result, err := operator.Evaluate(value, assertedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true, got false")
	}

	result, err = operator.Evaluate(value, failedValue)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false, got true")
	}
}

// Tests for toFloat64 to improve coverage
func TestToFloat64_AllNumericTypes(t *testing.T) {
	// Test int types
	var i int = 5
	var i8 int8 = 5
	var i16 int16 = 5
	var i32 int32 = 5
	var i64 int64 = 5

	// Test uint types
	var ui uint = 5
	var ui8 uint8 = 5
	var ui16 uint16 = 5
	var ui32 uint32 = 5
	var ui64 uint64 = 5

	// Test float types
	var f32 float32 = 5.0
	var f64 float64 = 5.0

	// Use operators to test toFloat64 indirectly
	operator := &gorulesengine.LessThanOperator{}

	tests := []interface{}{i, i8, i16, i32, i64, ui, ui8, ui16, ui32, ui64, f32, f64}

	for _, val := range tests {
		result, err := operator.Evaluate(val, 10.0)
		if err != nil {
			t.Errorf("Unexpected error for type %T: %v", val, err)
		}
		if !result {
			t.Errorf("Expected true for %v < 10, got false", val)
		}
	}
}

// Tests for NotEqualOperator with error propagation
func TestNotEqualOperator_PropagatesError(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.NotEqualOperator{}

	// Test with nil values that should generate an error
	_, err := operator.Evaluate(nil, "test")
	if err == nil {
		t.Errorf("Expected error to be propagated from EqualOperator, got nil")
	}
}

// Tests for NotInOperator with error propagation
func TestNotInOperator_PropagatesError(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.NotInOperator{}

	// Test with compareValue that is not a slice
	_, err := operator.Evaluate(5, "not a slice")
	if err == nil {
		t.Errorf("Expected error to be propagated from InOperator, got nil")
	}
}

// Tests for NotContainsOperator with error propagation
func TestNotContainsOperator_PropagatesError(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.NotContainsOperator{}

	// Test with invalid factValue
	_, err := operator.Evaluate(123, "test")
	if err == nil {
		t.Errorf("Expected error to be propagated from ContainsOperator, got nil")
	}
}

// Tests for InOperator with empty slice
func TestInOperator_EmptySlice(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.InOperator{}

	result, err := operator.Evaluate(5, []int{})
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false for value in empty slice, got true")
	}
}

// Test for InOperator with error propagated from EqualOperator
func TestInOperator_PropagatesEqualError(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.InOperator{}

	// Create a slice containing nil, which should trigger an error in EqualOperator
	sliceWithNil := []interface{}{1, 2, nil, 4}
	_, err := operator.Evaluate(5, sliceWithNil)
	if err == nil {
		t.Errorf("Expected error to be propagated from EqualOperator when comparing with nil, got nil")
	}
}

// Tests for ContainsOperator with empty slice
func TestContainsOperator_EmptySlice(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.ContainsOperator{}

	result, err := operator.Evaluate([]int{}, 5)
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false for empty slice contains value, got true")
	}
}

// Test for ContainsOperator with error propagated from EqualOperator
func TestContainsOperator_PropagatesEqualError(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.ContainsOperator{}

	// Create a slice containing nil, which should trigger an error in EqualOperator
	sliceWithNil := []interface{}{1, 2, nil, 4}
	_, err := operator.Evaluate(sliceWithNil, 5)
	if err == nil {
		t.Errorf("Expected error to be propagated from EqualOperator when slice contains nil, got nil")
	}
}

// Tests for ContainsOperator with empty string
func TestContainsOperator_EmptyString(t *testing.T) {
	var operator gorulesengine.Operator = &gorulesengine.ContainsOperator{}

	result, err := operator.Evaluate("", "test")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if result {
		t.Errorf("Expected false for empty string contains value, got true")
	}

	// Test empty substring (always true)
	result, err = operator.Evaluate("test", "")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}
	if !result {
		t.Errorf("Expected true for string contains empty substring, got false")
	}
}

// Tests for GetOperator to verify operator retrieval
func TestGetOperator(t *testing.T) {
	tests := []struct {
		name       string
		opType     gorulesengine.OperatorType
		wantExists bool
	}{
		{
			name:       "Existing operator Equal",
			opType:     "equal",
			wantExists: true,
		},
		{
			name:       "Existing operator NotEqual",
			opType:     "not_equal",
			wantExists: true,
		},
		{
			name:       "Existing operator LessThan",
			opType:     "less_than",
			wantExists: true,
		},
		{
			name:       "Existing operator LessThanInclusive",
			opType:     "less_than_inclusive",
			wantExists: true,
		},
		{
			name:       "Existing operator GreaterThan",
			opType:     "greater_than",
			wantExists: true,
		},
		{
			name:       "Existing operator GreaterThanInclusive",
			opType:     "greater_than_inclusive",
			wantExists: true,
		},
		{
			name:       "Existing operator NotIn",
			opType:     "not_in",
			wantExists: true,
		},
		{
			name:       "Existing operator In",
			opType:     "in",
			wantExists: true,
		},
		{
			name:       "Existing operator Contains",
			opType:     "contains",
			wantExists: true,
		},
		{
			name:       "Existing operator NotContains",
			opType:     "not_contains",
			wantExists: true,
		},
		{
			name:       "Non-existing operator",
			opType:     "non_existing_operator",
			wantExists: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			op, err := gorulesengine.GetOperator(tt.opType)
			if tt.wantExists {
				if err != nil {
					t.Errorf("Expected operator to exist, got error: %v", err)
				}
				if op == nil {
					t.Errorf("Expected operator instance, got nil")
				}
			}

			if !tt.wantExists {
				if err == nil {
					t.Errorf("Expected error for non-existing operator, got nil")
				}
				if op != nil {
					t.Errorf("Expected nil operator for non-existing operator, got instance")
				}

				ruleErr, ok := err.(*gorulesengine.OperatorError)
				if !ok {
					t.Errorf("Expected OperatorError type, got %T", err)
				} else {
					expectedMsg := fmt.Sprintf("[OPERATOR_ERROR] operator=%s value=%v compareValue=%v: operator not registered", tt.opType, nil, nil)
					if ruleErr.Error() != expectedMsg {
						t.Errorf("Expected error message '%s', got '%s'", expectedMsg, ruleErr.Error())
					}
				}
			}
		})
	}
}

// Définir un opérateur personnalisé simple au niveau du package
type CustomOperator struct{}

func (o *CustomOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	return true, nil
}

func TestRegisterOperator(t *testing.T) {
	customOpType := gorulesengine.OperatorType("custom_operator")

	customOperator := &CustomOperator{}

	// Enregistrer l'opérateur personnalisé
	gorulesengine.RegisterOperator(customOpType, customOperator)

	// Récupérer l'opérateur enregistré
	retrievedOp, err := gorulesengine.GetOperator(customOpType)
	if err != nil {
		t.Errorf("Unexpected error retrieving custom operator: %v", err)
	}

	if retrievedOp == nil {
		t.Errorf("Expected to retrieve custom operator, got nil")
	}

	// Verify that the retrieved operator is indeed the registered one
	_, ok := retrievedOp.(*CustomOperator)
	if !ok {
		t.Errorf("Expected retrieved operator to be of type CustomOperator, got %T", retrievedOp)
	}
}

package gorulesengine

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"
)

// Operator defines the interface for all comparison operators.
// Custom operators can be registered by implementing this interface.
type Operator interface {
	// Evaluate compares a fact value against a condition value and returns true if the comparison succeeds.
	Evaluate(factValue interface{}, compareValue interface{}) (bool, error)
}

// EqualOperator checks if two values are equal.
type EqualOperator struct{}

// NotEqualOperator checks if two values are not equal.
type NotEqualOperator struct{}

// LessThanOperator checks if factValue < compareValue.
type LessThanOperator struct{}

// LessThanInclusiveOperator checks if factValue <= compareValue.
type LessThanInclusiveOperator struct{}

// GreaterThanOperator checks if factValue > compareValue.
type GreaterThanOperator struct{}

// GreaterThanInclusiveOperator checks if factValue >= compareValue.
type GreaterThanInclusiveOperator struct{}

// InOperator checks if factValue is contained in compareValue (array).
type InOperator struct{}

// NotInOperator checks if factValue is not contained in compareValue (array).
type NotInOperator struct{}

// ContainsOperator checks if factValue contains compareValue (for strings and arrays).
type ContainsOperator struct{}

// NotContainsOperator checks if factValue does not contain compareValue.
type NotContainsOperator struct{}

// RegexOperator checks if factValue matches the regex pattern in compareValue.
type RegexOperator struct{}

var operatorRegistry map[OperatorType]Operator

func init() {
	operatorRegistry = map[OperatorType]Operator{
		OperatorEqual:                &EqualOperator{},
		OperatorNotEqual:             &NotEqualOperator{},
		OperatorLessThan:             &LessThanOperator{},
		OperatorLessThanInclusive:    &LessThanInclusiveOperator{},
		OperatorGreaterThan:          &GreaterThanOperator{},
		OperatorGreaterThanInclusive: &GreaterThanInclusiveOperator{},
		OperatorIn:                   &InOperator{},
		OperatorNotIn:                &NotInOperator{},
		OperatorContains:             &ContainsOperator{},
		OperatorNotContains:          &NotContainsOperator{},
		OperatorRegex:                &RegexOperator{},
	}
}

// GetOperator retrieves an operator from the registry by its type.
// Returns an error if the operator is not registered.
func GetOperator(opType OperatorType) (Operator, error) {
	op, exists := operatorRegistry[opType]
	if !exists {
		return nil, &OperatorError{
			Operator: opType,
			Err:      fmt.Errorf("operator not registered"),
		}
	}
	return op, nil
}

// RegisterOperator registers a custom operator in the global operator registry.
// This allows you to extend the engine with custom comparison logic.
//
// Example:
//
//	type StartsWithOperator struct{}
//	func (o *StartsWithOperator) Evaluate(factValue, compareValue interface{}) (bool, error) {
//	    str, ok1 := factValue.(string)
//	    prefix, ok2 := compareValue.(string)
//	    if !ok1 || !ok2 {
//	        return false, fmt.Errorf("both values must be strings")
//	    }
//	    return strings.HasPrefix(str, prefix), nil
//	}
//	gorulesengine.RegisterOperator("starts_with", &StartsWithOperator{})
func RegisterOperator(opType OperatorType, operator Operator) {
	operatorRegistry[opType] = operator
}

// toFloat64 converts any numeric type to float64
func toFloat64(value interface{}) (float64, bool) {
	switch v := value.(type) {
	case float64:
		return v, true
	case float32:
		return float64(v), true
	case int:
		return float64(v), true
	case int64:
		return float64(v), true
	case int32:
		return float64(v), true
	case int16:
		return float64(v), true
	case int8:
		return float64(v), true
	case uint:
		return float64(v), true
	case uint64:
		return float64(v), true
	case uint32:
		return float64(v), true
	case uint16:
		return float64(v), true
	case uint8:
		return float64(v), true
	default:
		// Fallback with reflection for exotic types
		return 0, false
	}
}

// Evaluate checks if two values are equal using deep equality comparison.
// Returns false if the values have different types or if either value is nil.
func (o *EqualOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	if factValue == nil || compareValue == nil {
		return false, &OperatorError{
			Operator:     OperatorEqual,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          fmt.Errorf("cannot compare nil values"),
		}
	}

	if reflect.TypeOf(factValue) != reflect.TypeOf(compareValue) {
		return false, nil
	}

	return reflect.DeepEqual(factValue, compareValue), nil
}

// Evaluate checks if two values are not equal.
// Returns the inverse of the EqualOperator result.
func (o *NotEqualOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	equal, err := (&EqualOperator{}).Evaluate(factValue, compareValue)
	if err != nil {
		return false, &OperatorError{
			Operator:     OperatorNotEqual,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          err,
		}
	}
	return !equal, nil
}

// Evaluate checks if factValue is less than compareValue.
// Both values must be numeric types.
func (o *LessThanOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	fv, ok1 := toFloat64(factValue)
	cv, ok2 := toFloat64(compareValue)
	if !ok1 || !ok2 {
		return false, &OperatorError{
			Operator:     OperatorLessThan,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          fmt.Errorf("less_than operator requires numeric values"),
		}
	}
	return fv < cv, nil
}

// Evaluate checks if factValue is less than or equal to compareValue.
// Both values must be numeric types.
func (o *LessThanInclusiveOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	fv, ok1 := toFloat64(factValue)
	cv, ok2 := toFloat64(compareValue)
	if !ok1 || !ok2 {
		return false, &OperatorError{
			Operator:     OperatorLessThanInclusive,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          fmt.Errorf("less_than_inclusive operator requires numeric values"),
		}
	}
	return fv <= cv, nil
}

// Evaluate checks if factValue is greater than compareValue.
// Both values must be numeric types.
func (o *GreaterThanOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	fv, ok1 := toFloat64(factValue)
	cv, ok2 := toFloat64(compareValue)
	if !ok1 || !ok2 {
		return false, &OperatorError{
			Operator:     OperatorGreaterThan,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          fmt.Errorf("greater_than operator requires numeric values"),
		}
	}
	return fv > cv, nil
}

// Evaluate checks if factValue is greater than or equal to compareValue.
// Both values must be numeric types.
func (o *GreaterThanInclusiveOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	fv, ok1 := toFloat64(factValue)
	cv, ok2 := toFloat64(compareValue)
	if !ok1 || !ok2 {
		return false, &OperatorError{
			Operator:     OperatorGreaterThanInclusive,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          fmt.Errorf("greater_than_inclusive operator requires numeric values"),
		}
	}
	return fv >= cv, nil
}

// Evaluate checks if factValue is contained in the compareValue array.
// compareValue must be a slice or array.
func (o *InOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	// Use reflection to handle any slice type
	rv := reflect.ValueOf(compareValue)

	// Verify that it's a slice or an array
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return false, &OperatorError{
			Operator:     OperatorIn,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          fmt.Errorf("in operator requires an array or slice as compareValue"),
		}
	}

	// Iterate over slice elements
	for i := 0; i < rv.Len(); i++ {
		elem := rv.Index(i).Interface()
		equal, err := (&EqualOperator{}).Evaluate(factValue, elem)
		if err != nil {
			return false, err
		}
		if equal {
			return true, nil
		}
	}

	return false, nil
}

// Evaluate checks if factValue is not contained in the compareValue array.
// Returns the inverse of the InOperator result.
func (o *NotInOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	in, err := (&InOperator{}).Evaluate(factValue, compareValue)
	if err != nil {
		return false, &OperatorError{
			Operator:     OperatorNotIn,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          err,
		}
	}
	return !in, nil
}

// Evaluate checks if factValue contains compareValue.
// For strings, checks substring containment. For arrays/slices, checks element presence.
func (o *ContainsOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	// Use reflection to handle any slice or string type
	rv := reflect.ValueOf(factValue)

	switch rv.Kind() {
	case reflect.Slice, reflect.Array:
		// Iterate over slice elements
		for i := 0; i < rv.Len(); i++ {
			elem := rv.Index(i).Interface()
			equal, err := (&EqualOperator{}).Evaluate(elem, compareValue)
			if err != nil {
				return false, &OperatorError{
					Operator:     OperatorContains,
					Value:        factValue,
					CompareValue: compareValue,
					Err:          fmt.Errorf("error during contains evaluation: %v", err),
				}
			}
			if equal {
				return true, nil
			}
		}
		return false, nil
	case reflect.String:
		cv, ok := compareValue.(string)
		if !ok {
			return false, &OperatorError{
				Operator:     OperatorContains,
				Value:        factValue,
				CompareValue: compareValue,
				Err:          fmt.Errorf("contains operator requires string compareValue when factValue is a string"),
			}
		}
		return strings.Contains(rv.String(), cv), nil
	default:
		return false, &OperatorError{
			Operator:     OperatorContains,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          fmt.Errorf("contains operator requires array, slice, or string as factValue"),
		}
	}
}

// Evaluate checks if factValue does not contain compareValue.
// Returns the inverse of the ContainsOperator result.
func (o *NotContainsOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	contains, err := (&ContainsOperator{}).Evaluate(factValue, compareValue)
	if err != nil {
		return false, err
	}
	return !contains, nil
}

// Evaluate checks if factValue matches the regex pattern in compareValue.
// Both values must be strings.
// Returns an error if regex evaluation fails.
func (o *RegexOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	strValue, ok1 := factValue.(string)
	pattern, ok2 := compareValue.(string)
	if !ok1 || !ok2 {
		return false, &OperatorError{
			Operator:     OperatorRegex,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          fmt.Errorf("regex operator requires string values"),
		}
	}
	matched, err := regexp.MatchString(pattern, strValue)
	if err != nil {
		return false, &OperatorError{
			Operator:     OperatorRegex,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          fmt.Errorf("regex evaluation error: %v", err),
		}
	}
	return matched, nil
}

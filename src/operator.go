package gorulesengine

import (
	"fmt"
	"reflect"
	"strings"
)

type Operator interface {
	Evaluate(factValue interface{}, compareValue interface{}) (bool, error)
}

type EqualOperator struct{}
type NotEqualOperator struct{}
type LessThanOperator struct{}
type LessThanInclusiveOperator struct{}
type GreaterThanOperator struct{}
type GreaterThanInclusiveOperator struct{}
type InOperator struct{}
type NotInOperator struct{}
type ContainsOperator struct{}
type NotContainsOperator struct{}

var operatorRegistry map[OperatorType]Operator

func init() {
	operatorRegistry = map[OperatorType]Operator{
		Equal:                &EqualOperator{},
		NotEqual:             &NotEqualOperator{},
		LessThan:             &LessThanOperator{},
		LessThanInclusive:    &LessThanInclusiveOperator{},
		GreaterThan:          &GreaterThanOperator{},
		GreaterThanInclusive: &GreaterThanInclusiveOperator{},
		In:                   &InOperator{},
		NotIn:                &NotInOperator{},
		Contains:             &ContainsOperator{},
		NotContains:          &NotContainsOperator{},
	}
}

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

func (o *EqualOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	if factValue == nil || compareValue == nil {
		return false, &OperatorError{
			Operator:     Equal,
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

func (o *NotEqualOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	equal, err := (&EqualOperator{}).Evaluate(factValue, compareValue)
	if err != nil {
		return false, &OperatorError{
			Operator:     NotEqual,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          err,
		}
	}
	return !equal, nil
}

func (o *LessThanOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	fv, ok1 := toFloat64(factValue)
	cv, ok2 := toFloat64(compareValue)
	if !ok1 || !ok2 {
		return false, &OperatorError{
			Operator:     LessThan,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          fmt.Errorf("less_than operator requires numeric values"),
		}
	}
	return fv < cv, nil
}

func (o *LessThanInclusiveOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	fv, ok1 := toFloat64(factValue)
	cv, ok2 := toFloat64(compareValue)
	if !ok1 || !ok2 {
		return false, &OperatorError{
			Operator:     LessThanInclusive,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          fmt.Errorf("less_than_inclusive operator requires numeric values"),
		}
	}
	return fv <= cv, nil
}

func (o *GreaterThanOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	fv, ok1 := toFloat64(factValue)
	cv, ok2 := toFloat64(compareValue)
	if !ok1 || !ok2 {
		return false, &OperatorError{
			Operator:     GreaterThan,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          fmt.Errorf("greater_than operator requires numeric values"),
		}
	}
	return fv > cv, nil
}

func (o *GreaterThanInclusiveOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	fv, ok1 := toFloat64(factValue)
	cv, ok2 := toFloat64(compareValue)
	if !ok1 || !ok2 {
		return false, &OperatorError{
			Operator:     GreaterThanInclusive,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          fmt.Errorf("greater_than_inclusive operator requires numeric values"),
		}
	}
	return fv >= cv, nil
}

func (o *InOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	// Use reflection to handle any slice type
	rv := reflect.ValueOf(compareValue)

	// Verify that it's a slice or an array
	if rv.Kind() != reflect.Slice && rv.Kind() != reflect.Array {
		return false, &OperatorError{
			Operator:     In,
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

func (o *NotInOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	in, err := (&InOperator{}).Evaluate(factValue, compareValue)
	if err != nil {
		return false, &OperatorError{
			Operator:     NotIn,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          err,
		}
	}
	return !in, nil
}

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
					Operator:     Contains,
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
				Operator:     Contains,
				Value:        factValue,
				CompareValue: compareValue,
				Err:          fmt.Errorf("contains operator requires string compareValue when factValue is a string"),
			}
		}
		return strings.Contains(rv.String(), cv), nil
	default:
		return false, &OperatorError{
			Operator:     Contains,
			Value:        factValue,
			CompareValue: compareValue,
			Err:          fmt.Errorf("contains operator requires array, slice, or string as factValue"),
		}
	}
}

func (o *NotContainsOperator) Evaluate(factValue interface{}, compareValue interface{}) (bool, error) {
	contains, err := (&ContainsOperator{}).Evaluate(factValue, compareValue)
	if err != nil {
		return false, err
	}
	return !contains, nil
}

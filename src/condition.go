package gorulesengine

import (
	"encoding/json"
	"fmt"
)

type Condition struct {
	Fact     FactID                 `json:"fact"`
	Operator OperatorType           `json:"operator"`
	Value    interface{}            `json:"value"`
	Path     string                 `json:"path,omitempty"`
	Params   map[string]interface{} `json:"params,omitempty"`
}

type ConditionSet struct {
	All  []ConditionNode `json:"all,omitempty"`
	Any  []ConditionNode `json:"any,omitempty"`
	None []ConditionNode `json:"none,omitempty"`
}

type ConditionNode struct {
	Condition *Condition
	SubSet    *ConditionSet
}

func (n *ConditionNode) UnmarshalJSON(data []byte) error {
	var cond Condition
	err1 := json.Unmarshal(data, &cond)
	if err1 == nil && cond.Fact != "" {
		n.Condition = &cond
		return nil
	}

	var subset ConditionSet
	err2 := json.Unmarshal(data, &subset)
	if err2 == nil {
		n.SubSet = &subset
		return nil
	}

	// If both failed, wrap the error
	return &RuleEngineError{
		Type: ErrJSON,
		Msg:  "failed to unmarshal ConditionNode",
		Err:  fmt.Errorf("data: %s, errors: %v, %v", string(data), err1, err2),
	}
}

// Evaluate evaluates the condition node, whether it's a condition or a subset
func evaluateConditionNode(node *ConditionNode, almanac *Almanac) (bool, error) {
	if node.Condition != nil {
		return node.Condition.Evaluate(almanac)
	} else if node.SubSet != nil {
		return node.SubSet.Evaluate(almanac)
	}

	return false, &ConditionError{
		Condition: Condition{},
		Err:       fmt.Errorf("invalid condition node: neither condition nor subset is defined"),
	}
}

// Evaluate evaluates the condition against the almanac
func (c *Condition) Evaluate(almanac *Almanac) (bool, error) {
	// Here params can be passed to the fact calculation
	// Usefull only for dynamic facts
	// For static facts, params are ignored
	factValue, err := almanac.GetFactValue(c.Fact, c.Params, c.Path)
	if err != nil {
		return false, &ConditionError{
			Condition: *c,
			Err:       fmt.Errorf("failed to get fact value: %v", err),
		}
	}

	operator, err := GetOperator(c.Operator)
	if err != nil {
		return false, &ConditionError{
			Condition: *c,
			Err:       fmt.Errorf("failed to get operator: %v", err),
		}
	}

	result, err := operator.Evaluate(factValue, c.Value)
	if err != nil {
		return false, &ConditionError{
			Condition: *c,
			Err:       fmt.Errorf("operator evaluation failed: %v", err),
		}
	}

	return result, nil
}

// Evaluate evaluates the condition set against the almanac
func (cs *ConditionSet) Evaluate(almanac *Almanac) (bool, error) {
	// Evaluate "all" conditions
	for _, node := range cs.All {
		result, err := evaluateConditionNode(&node, almanac)
		if err != nil {
			return false, err
		}
		if !result {
			return false, nil
		}
	}

	// Evaluate "any" conditions
	if len(cs.Any) > 0 {
		anyMatched := false
		for _, node := range cs.Any {
			result, err := evaluateConditionNode(&node, almanac)
			if err != nil {
				return false, err
			}
			if result {
				anyMatched = true
				break
			}
		}
		if !anyMatched {
			return false, nil
		}
	}

	// Evaluate "none" conditions
	for _, node := range cs.None {
		result, err := evaluateConditionNode(&node, almanac)
		if err != nil {
			return false, err
		}
		if result {
			return false, nil
		}
	}

	return true, nil
}

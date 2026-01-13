package gorulesengine

import (
	"encoding/json"
	"fmt"
)

type Condition struct {
	Fact     string                 `json:"fact"`
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

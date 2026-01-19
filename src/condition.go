package gorulesengine

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
)

// Condition represents a single condition that compares a fact value against an expected value using an operator.
// Conditions can optionally use JSONPath to access nested values within facts.
//
// Example:
//
//	condition := &gre.Condition{
//	    Fact:     "age",
//	    Operator: "greater_than",
//	    Value:    18,
//	}
type Condition struct {
	Fact      FactID                 `json:"fact"`             // The fact identifier to evaluate
	Operator  OperatorType           `json:"operator"`         // The comparison operator to use
	Value     interface{}            `json:"value"`            // The expected value to compare against
	Path      string                 `json:"path,omitempty"`   // Optional JSONPath to access nested fact values
	Params    map[string]interface{} `json:"params,omitempty"` // Optional parameters for dynamic facts
	cachedKey string                 // Pre-calculated cache key
}

// ConditionSet represents a group of conditions combined with logical operators (all/any/none).
// ConditionSets can be nested to create complex boolean logic.
//
// Example:
//
//	conditionSet := gre.ConditionSet{
//	    All: []gre.ConditionNode{
//	        {Condition: &condition1},
//	        {Condition: &condition2},
//	    },
//	}
type ConditionSet struct {
	All       []ConditionNode `json:"all,omitempty"`  // All conditions must be true (AND)
	Any       []ConditionNode `json:"any,omitempty"`  // At least one condition must be true (OR)
	None      []ConditionNode `json:"none,omitempty"` // No conditions must be true (NOT)
	cachedKey string          // Pre-calculated cache key
}

// ConditionNode represents either a single Condition or a nested ConditionSet.
// This allows for recursive nesting of conditions to build complex boolean expressions.
type ConditionNode struct {
	Condition *Condition    // A single condition to evaluate
	SubSet    *ConditionSet // A nested set of conditions
}

// UnmarshalJSON implements custom JSON unmarshaling for ConditionNode.
// It attempts to unmarshal either a Condition or a ConditionSet from the JSON data.
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
func evaluateConditionNode(node *ConditionNode, almanac *Almanac) (*ConditionNodeResult, error) {
	if node.Condition != nil {
		res, err := node.Condition.Evaluate(almanac)
		if err != nil {
			return nil, &ConditionError{
				Condition: *node.Condition,
				Err:       fmt.Errorf("failed to evaluate condition node: %v", err),
			}
		}
		return &ConditionNodeResult{Condition: res}, nil
	} else if node.SubSet != nil {
		res, err := node.SubSet.Evaluate(almanac)
		if err != nil {
			return nil, &ConditionError{
				Condition: Condition{},
				Err:       fmt.Errorf("failed to evaluate condition subset node: %v", err),
			}
		}
		return &ConditionNodeResult{ConditionSet: res}, nil
	}

	return nil, &ConditionError{
		Condition: Condition{},
		Err:       fmt.Errorf("invalid condition node: neither condition nor subset is defined"),
	}
}

// GetCacheKey generates a unique cache key for the condition.
func (c *Condition) GetCacheKey() (string, error) {
	if c.cachedKey != "" {
		return c.cachedKey, nil
	}
	bytes, err := json.Marshal(c)
	if err != nil {
		return "", &ConditionError{
			Condition: *c,
			Err:       fmt.Errorf("failed to marshal condition for cache key: %v", err),
		}
	}
	sum := md5.Sum(bytes)
	return hex.EncodeToString(sum[:]), nil
}

// Compile pre-calculates properties of the condition to speed up evaluation.
func (c *Condition) Compile() error {
	key, err := c.GetCacheKey()
	if err != nil {
		return &ConditionError{
			Condition: *c,
			Err:       fmt.Errorf("failed to compile condition: %v", err),
		}
	}
	c.cachedKey = key
	return nil
}

// GetRequiredFacts returns the list of facts required by this condition.
func (c *Condition) GetRequiredFacts() []FactID {
	return []FactID{c.Fact}
}

// GetRequiredFacts returns the list of facts required by this condition node.
func (n *ConditionNode) GetRequiredFacts() []FactID {
	if n.Condition != nil {
		return n.Condition.GetRequiredFacts()
	} else if n.SubSet != nil {
		return n.SubSet.GetRequiredFacts()
	}
	return []FactID{}
}

// GetRequiredFacts returns the list of all facts required by this condition set.
func (cs *ConditionSet) GetRequiredFacts() []FactID {
	factMap := make(map[FactID]bool)

	var collect func(nodes []ConditionNode)
	collect = func(nodes []ConditionNode) {
		for _, node := range nodes {
			for _, fact := range node.GetRequiredFacts() {
				factMap[fact] = true
			}
		}
	}

	collect(cs.All)
	collect(cs.Any)
	collect(cs.None)

	facts := make([]FactID, 0, len(factMap))
	for fact := range factMap {
		facts = append(facts, fact)
	}
	return facts
}

// Compile pre-calculates properties for the entire condition set.
func (cs *ConditionSet) Compile() error {
	compileNodes := func(nodes []ConditionNode) error {
		for i := range nodes {
			if nodes[i].Condition != nil {
				if err := nodes[i].Condition.Compile(); err != nil {
					return &ConditionError{
						Condition: *nodes[i].Condition,
						Err:       fmt.Errorf("failed to compile condition in condition set: %v", err),
					}
				}
			} else if nodes[i].SubSet != nil {
				if err := nodes[i].SubSet.Compile(); err != nil {
					return &ConditionError{
						Condition: Condition{},
						Err:       fmt.Errorf("failed to compile subset in condition set: %v", err),
					}
				}
			}
		}
		return nil
	}

	if err := compileNodes(cs.All); err != nil {
		return err
	}
	if err := compileNodes(cs.Any); err != nil {
		return err
	}
	if err := compileNodes(cs.None); err != nil {
		return err
	}

	key, err := cs.GetCacheKey()
	if err != nil {
		return &ConditionError{
			Condition: Condition{},
			Err:       fmt.Errorf("failed to compile condition set: %v", err),
		}
	}
	cs.cachedKey = key

	return nil
}

// Evaluate evaluates the condition against the almanac
func (c *Condition) Evaluate(almanac *Almanac) (*ConditionResult, error) {
	var cacheKey string
	var err error

	result := &ConditionResult{
		Fact:     c.Fact,
		Operator: c.Operator,
		Value:    c.Value,
		Path:     c.Path,
	}

	// Check cache if enabled
	if almanac.IsConditionCachingEnabled() {
		cacheKey, err = c.GetCacheKey()
		if err != nil {
			return nil, &ConditionError{
				Condition: *c,
				Err:       fmt.Errorf("failed to get cache key for condition: %v", err),
			}
		}
		if cachedVal, cached := almanac.GetConditionResultFromCache(cacheKey); cached {
			if cachedRes, ok := cachedVal.(*ConditionResult); ok {
				return cachedRes, nil
			}
		}
	}

	// Here params can be passed to the fact calculation
	// Usefull only for dynamic facts
	// For static facts, params are ignored
	factValue, err := almanac.GetFactValue(c.Fact, c.Params, c.Path)
	if err != nil {
		return nil, &ConditionError{
			Condition: *c,
			Err:       fmt.Errorf("failed to get fact value: %v", err),
		}
	}

	result.FactValue = factValue

	operator, err := GetOperator(c.Operator)
	if err != nil {
		return nil, &ConditionError{
			Condition: *c,
			Err:       fmt.Errorf("failed to get operator: %v", err),
		}
	}

	evalRes, err := operator.Evaluate(factValue, c.Value)
	if err != nil {
		return nil, &ConditionError{
			Condition: *c,
			Err:       fmt.Errorf("operator evaluation failed: %v", err),
		}
	}

	result.Result = evalRes

	// Cache result if caching is enabled
	if almanac.IsConditionCachingEnabled() && cacheKey != "" {
		almanac.SetConditionResultCache(cacheKey, result)
	}

	return result, nil
}

// GetCacheKey generates a unique cache key for the condition set.
func (cs *ConditionSet) GetCacheKey() (string, error) {
	if cs.cachedKey != "" {
		return cs.cachedKey, nil
	}

	bytes, err := json.Marshal(cs)
	if err != nil {
		return "", &ConditionError{
			Condition: Condition{},
			Err:       fmt.Errorf("failed to marshal condition set for cache key: %v", err),
		}
	}
	sum := md5.Sum(bytes)
	return hex.EncodeToString(sum[:]), nil
}

// Evaluate evaluates the condition set against the almanac
func (cs *ConditionSet) Evaluate(almanac *Almanac) (*ConditionSetResult, error) {
	// Check cache if enabled
	var err error
	var cacheKey string
	if almanac.IsConditionCachingEnabled() {
		cacheKey, err = cs.GetCacheKey()
		if err != nil {
			return nil, &ConditionError{
				Condition: Condition{},
				Err:       fmt.Errorf("failed to get cache key for condition set: %v", err),
			}
		}
		if cachedVal, cached := almanac.GetConditionResultFromCache(cacheKey); cached {
			if cachedRes, ok := cachedVal.(*ConditionSetResult); ok {
				return cachedRes, nil
			}
		}
	}

	result := &ConditionSetResult{
		Results: make([]ConditionNodeResult, 0),
		Result:  true, // Default to true for empty condition sets
	}

	// Reorder nodes to put cached conditions first if caching is enabled
	allNodes := cs.All
	anyNodes := cs.Any
	noneNodes := cs.None

	if almanac.IsConditionCachingEnabled() {
		allNodes, err = cs.ReorderNodes(cs.All, almanac)
		if err != nil {
			return nil, &ConditionError{
				Condition: Condition{},
				Err:       fmt.Errorf("failed to reorder all nodes: %v", err),
			}
		}
		anyNodes, err = cs.ReorderNodes(cs.Any, almanac)
		if err != nil {
			return nil, &ConditionError{
				Condition: Condition{},
				Err:       fmt.Errorf("failed to reorder any nodes: %v", err),
			}
		}
		noneNodes, err = cs.ReorderNodes(cs.None, almanac)
		if err != nil {
			return nil, &ConditionError{
				Condition: Condition{},
				Err:       fmt.Errorf("failed to reorder none nodes: %v", err),
			}
		}
	}

	// Determine set type
	if len(cs.All) > 0 {
		result.Type = AllType
		result.Result = true
		for _, node := range allNodes {
			nodeRes, err := evaluateConditionNode(&node, almanac)
			if err != nil {
				return nil, &ConditionError{
					Condition: Condition{},
					Err:       fmt.Errorf("failed to evaluate all node: %v", err),
				}
			}
			result.Results = append(result.Results, *nodeRes)

			var res bool
			if nodeRes.Condition != nil {
				res = nodeRes.Condition.Result
			} else {
				res = nodeRes.ConditionSet.Result
			}

			if !res {
				result.Result = false
				break // Short-circuit
			}
		}
	} else if len(cs.Any) > 0 {
		result.Type = AnyType
		result.Result = false
		for _, node := range anyNodes {
			nodeRes, err := evaluateConditionNode(&node, almanac)
			if err != nil {
				return nil, &ConditionError{
					Condition: Condition{},
					Err:       fmt.Errorf("failed to evaluate any node: %v", err),
				}
			}
			result.Results = append(result.Results, *nodeRes)

			var res bool
			if nodeRes.Condition != nil {
				res = nodeRes.Condition.Result
			} else {
				res = nodeRes.ConditionSet.Result
			}

			if res {
				result.Result = true
				break // Short-circuit
			}
		}
	} else if len(cs.None) > 0 {
		result.Type = NoneType
		result.Result = true
		for _, node := range noneNodes {
			nodeRes, err := evaluateConditionNode(&node, almanac)
			if err != nil {
				return nil, &ConditionError{
					Condition: Condition{},
					Err:       fmt.Errorf("failed to evaluate none node: %v", err),
				}
			}
			result.Results = append(result.Results, *nodeRes)

			var res bool
			if nodeRes.Condition != nil {
				res = nodeRes.Condition.Result
			} else {
				res = nodeRes.ConditionSet.Result
			}

			if res {
				result.Result = false
				break // Short-circuit
			}
		}
	}

	if almanac.IsConditionCachingEnabled() && cacheKey != "" {
		almanac.SetConditionResultCache(cacheKey, result)
	}
	return result, nil
}

// ReorderNodes puts cached conditions at the beginning of the slice to optimize short-circuiting.
func (cs *ConditionSet) ReorderNodes(nodes []ConditionNode, almanac *Almanac) ([]ConditionNode, error) {
	if len(nodes) <= 1 {
		return nodes, nil
	}

	cached := make([]ConditionNode, 0, len(nodes))
	notCached := make([]ConditionNode, 0, len(nodes))

	for _, node := range nodes {
		isCached := false
		if node.Condition != nil {
			key, err := node.Condition.GetCacheKey()
			if err != nil {
				return nil, &ConditionError{
					Condition: *node.Condition,
					Err:       fmt.Errorf("failed to get cache key for condition during reorder: %v", err),
				}
			}
			_, isCached = almanac.GetConditionResultFromCache(key)
		}

		if isCached {
			cached = append(cached, node)
		} else {
			notCached = append(notCached, node)
		}
	}

	return append(cached, notCached...), nil
}

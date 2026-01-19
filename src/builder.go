package gorulesengine

// RuleBuilder provides a fluent API for building rules.
type RuleBuilder struct {
	rule *Rule
}

// NewRuleBuilder creates a new RuleBuilder instance.
func NewRuleBuilder() *RuleBuilder {
	return &RuleBuilder{
		rule: &Rule{},
	}
}

// WithName sets the name of the rule.
func (rb *RuleBuilder) WithName(name string) *RuleBuilder {
	rb.rule.Name = name
	return rb
}

// WithPriority sets the priority of the rule.
func (rb *RuleBuilder) WithPriority(priority int) *RuleBuilder {
	rb.rule.Priority = priority
	return rb
}

// WithConditions sets the conditions for the rule.
func (rb *RuleBuilder) WithConditions(node ConditionNode) *RuleBuilder {
	if node.Condition != nil || node.SubSet != nil {
		rb.rule.Conditions = ConditionSet{}
		if node.Condition != nil {
			rb.rule.Conditions.All = []ConditionNode{node}
		} else if node.SubSet != nil {
			rb.rule.Conditions = ConditionSet{
				All:  node.SubSet.All,
				Any:  node.SubSet.Any,
				None: node.SubSet.None,
			}
		}
	}
	return rb
}

// WithOnSuccess sets the event names for the rule.
func (rb *RuleBuilder) WithOnSuccess(eventNames ...string) *RuleBuilder {
	events := make([]RuleEvent, len(eventNames))
	for i, name := range eventNames {
		events[i] = RuleEvent{Name: name}
	}
	rb.rule.OnSuccess = events
	return rb
}

// WithOnSuccessEvent adds a detailed event with parameters to the rule.
func (rb *RuleBuilder) WithOnSuccessEvent(event RuleEvent) *RuleBuilder {
	rb.rule.OnSuccess = append(rb.rule.OnSuccess, event)
	return rb
}

// WithOnFailure sets the event names for the rule.
func (rb *RuleBuilder) WithOnFailure(eventNames ...string) *RuleBuilder {
	events := make([]RuleEvent, len(eventNames))
	for i, name := range eventNames {
		events[i] = RuleEvent{Name: name}
	}
	rb.rule.OnFailure = events
	return rb
}

// WithOnFailureEvent adds a detailed event with parameters to the rule.
func (rb *RuleBuilder) WithOnFailureEvent(event RuleEvent) *RuleBuilder {
	rb.rule.OnFailure = append(rb.rule.OnFailure, event)
	return rb
}

// Build returns the constructed Rule.
func (rb *RuleBuilder) Build() *Rule {
	return rb.rule
}

// Condition Helper Functions

// Equal creates a condition that checks for equality.
func Equal(fact string, value interface{}) *Condition {
	return &Condition{
		Fact:     FactID(fact),
		Operator: OperatorEqual,
		Value:    value,
	}
}

// NotEqual creates a condition that checks for inequality.
func NotEqual(fact string, value interface{}) *Condition {
	return &Condition{
		Fact:     FactID(fact),
		Operator: OperatorNotEqual,
		Value:    value,
	}
}

// GreaterThan creates a condition that checks if fact > value.
func GreaterThan(fact string, value interface{}) *Condition {
	return &Condition{
		Fact:     FactID(fact),
		Operator: OperatorGreaterThan,
		Value:    value,
	}
}

// GreaterThanInclusive creates a condition that checks if fact >= value.
func GreaterThanInclusive(fact string, value interface{}) *Condition {
	return &Condition{
		Fact:     FactID(fact),
		Operator: OperatorGreaterThanInclusive,
		Value:    value,
	}
}

// LessThan creates a condition that checks if fact < value.
func LessThan(fact string, value interface{}) *Condition {
	return &Condition{
		Fact:     FactID(fact),
		Operator: OperatorLessThan,
		Value:    value,
	}
}

// LessThanInclusive creates a condition that checks if fact <= value.
func LessThanInclusive(fact string, value interface{}) *Condition {
	return &Condition{
		Fact:     FactID(fact),
		Operator: OperatorLessThanInclusive,
		Value:    value,
	}
}

// In creates a condition that checks if fact is in a list.
func In(fact string, values interface{}) *Condition {
	return &Condition{
		Fact:     FactID(fact),
		Operator: OperatorIn,
		Value:    values,
	}
}

// NotIn creates a condition that checks if fact is not in a list.
func NotIn(fact string, values interface{}) *Condition {
	return &Condition{
		Fact:     FactID(fact),
		Operator: OperatorNotIn,
		Value:    values,
	}
}

// Contains creates a condition that checks if fact contains value.
func Contains(fact string, value interface{}) *Condition {
	return &Condition{
		Fact:     FactID(fact),
		Operator: OperatorContains,
		Value:    value,
	}
}

// NotContains creates a condition that checks if fact does not contain value.
func NotContains(fact string, value interface{}) *Condition {
	return &Condition{
		Fact:     FactID(fact),
		Operator: OperatorNotContains,
		Value:    value,
	}
}

// Regex creates a condition that checks if fact matches a regex pattern.
func Regex(fact string, pattern string) *Condition {
	return &Condition{
		Fact:     FactID(fact),
		Operator: OperatorRegex,
		Value:    pattern,
	}
}

// ConditionSet Helper Functions

// All creates a ConditionSet where all conditions must be true.
func All(conditions ...*Condition) ConditionSet {
	nodes := make([]ConditionNode, len(conditions))
	for i, cond := range conditions {
		nodes[i] = ConditionNode{Condition: cond}
	}
	return ConditionSet{All: nodes}
}

// Any creates a ConditionSet where at least one condition must be true.
func Any(conditions ...*Condition) ConditionSet {
	nodes := make([]ConditionNode, len(conditions))
	for i, cond := range conditions {
		nodes[i] = ConditionNode{Condition: cond}
	}
	return ConditionSet{Any: nodes}
}

// None creates a ConditionSet where no conditions must be true.
func None(conditions ...*Condition) ConditionSet {
	nodes := make([]ConditionNode, len(conditions))
	for i, cond := range conditions {
		nodes[i] = ConditionNode{Condition: cond}
	}
	return ConditionSet{None: nodes}
}

// AllSets creates a ConditionSet where all nested ConditionSets must be true.
func AllSets(sets ...ConditionSet) ConditionSet {
	nodes := make([]ConditionNode, len(sets))
	for i, set := range sets {
		setCopy := set
		nodes[i] = ConditionNode{SubSet: &setCopy}
	}
	return ConditionSet{All: nodes}
}

// AnySets creates a ConditionSet where at least one nested ConditionSet must be true.
func AnySets(sets ...ConditionSet) ConditionSet {
	nodes := make([]ConditionNode, len(sets))
	for i, set := range sets {
		setCopy := set
		nodes[i] = ConditionNode{SubSet: &setCopy}
	}
	return ConditionSet{Any: nodes}
}

// NoneSets creates a ConditionSet where no nested ConditionSets must be true.
func NoneSets(sets ...ConditionSet) ConditionSet {
	nodes := make([]ConditionNode, len(sets))
	for i, set := range sets {
		setCopy := set
		nodes[i] = ConditionNode{SubSet: &setCopy}
	}
	return ConditionSet{None: nodes}
}

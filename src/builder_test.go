package gorulesengine_test

import (
	"fmt"
	"testing"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func TestNewRule(t *testing.T) {
	rule := gorulesengine.NewRuleBuilder()
	if rule == nil {
		t.Fatal("NewRule() returned nil")
	}
}

func TestRuleBuilder_WithName(t *testing.T) {
	rule := gorulesengine.NewRuleBuilder().
		WithName("test-rule").
		Build()

	if rule.Name != "test-rule" {
		t.Errorf("Expected name 'test-rule', got '%s'", rule.Name)
	}
}

func TestRuleBuilder_WithPriority(t *testing.T) {
	rule := gorulesengine.NewRuleBuilder().
		WithPriority(10).
		Build()

	if rule.Priority != 10 {
		t.Errorf("Expected priority 10, got %d", rule.Priority)
	}
}

func TestRuleBuilder_WithCondition(t *testing.T) {
	condition := gorulesengine.Equal("age", 18)
	rule := gorulesengine.NewRuleBuilder().
		WithConditions(
			gorulesengine.ConditionNode{
				Condition: condition,
			},
		).
		Build()

	if rule.Conditions.All == nil || len(rule.Conditions.All) != 1 {
		t.Fatal("Expected one condition in All conditions")
	}

	if rule.Conditions.All[0].Condition == nil {
		t.Fatal("Expected condition to be set")
	}

	if rule.Conditions.All[0].Condition.Fact != "age" {
		t.Errorf("Expected fact 'age', got '%s'", rule.Conditions.All[0].Condition.Fact)
	}

	if rule.Conditions.All[0].Condition.Operator != gorulesengine.OperatorEqual {
		t.Errorf("Expected operator OperatorEqual, got '%s'", rule.Conditions.All[0].Condition.Operator)
	}
}

func TestRuleBuilder_WithConditions(t *testing.T) {
	conditionSet := gorulesengine.ConditionSet{
		All: []gorulesengine.ConditionNode{
			{
				Condition: gorulesengine.Equal("age", 18),
			},
			{
				Condition: gorulesengine.GreaterThan("score", 100),
			},
		},
	}

	rule := gorulesengine.NewRuleBuilder().
		WithConditions(
			gorulesengine.ConditionNode{
				SubSet: &conditionSet,
			},
		).
		Build()

	if rule.Conditions.All == nil || len(rule.Conditions.All) != 2 {
		t.Fatal("Expected two conditions in All conditions")
	}
}

func TestRuleBuilder_WithEvent(t *testing.T) {
	rule := gorulesengine.NewRuleBuilder().
		WithEvent("test-event", nil).
		Build()

	if rule.Event.Type != "test-event" {
		t.Errorf("Expected event type 'test-event', got '%s'", rule.Event.Type)
	}
}

func TestRuleBuilder_WithOnSuccess(t *testing.T) {
	rule := gorulesengine.NewRuleBuilder().
		WithOnSuccess("success-callback").
		Build()

	if rule.OnSuccess == nil {
		t.Fatal("Expected one OnSuccess callback")
	}

	if *rule.OnSuccess != "success-callback" {
		t.Errorf("Expected callback 'success-callback', got '%s'", *rule.OnSuccess)
	}
}

func TestRuleBuilder_WithOnFailure(t *testing.T) {
	rule := gorulesengine.NewRuleBuilder().
		WithOnFailure("failure-callback").
		Build()

	if rule.OnFailure == nil {
		t.Fatal("Expected one OnFailure callback")
	}

	if *rule.OnFailure != "failure-callback" {
		t.Errorf("Expected callback 'failure-callback', got '%s'", *rule.OnFailure)
	}
}

func TestRuleBuilder_CompleteRule(t *testing.T) {
	rule := gorulesengine.NewRuleBuilder().
		WithName("complete-rule").
		WithPriority(5).
		WithConditions(
			gorulesengine.ConditionNode{
				Condition: gorulesengine.GreaterThan("age", 18),
			},
		).
		WithEvent("adult-event", nil).
		WithOnSuccess("log-success").
		WithOnFailure("log-failure").
		Build()

	if rule.Name != "complete-rule" {
		t.Errorf("Expected name 'complete-rule', got '%s'", rule.Name)
	}

	if rule.Priority != 5 {
		t.Errorf("Expected priority 5, got %d", rule.Priority)
	}

	if rule.Event.Type != "adult-event" {
		t.Errorf("Expected event type 'adult-event', got '%s'", rule.Event.Type)
	}

	if rule.OnSuccess == nil || *rule.OnSuccess != "log-success" {
		t.Errorf("Expected OnSuccess callback 'log-success'")
	}

	if rule.OnFailure == nil || *rule.OnFailure != "log-failure" {
		t.Errorf("Expected OnFailure callback 'log-failure'")
	}
}

// Test helper functions

func TestEqual(t *testing.T) {
	cond := gorulesengine.Equal("age", 18)
	if cond.Fact != "age" {
		t.Errorf("Expected fact 'age', got '%s'", cond.Fact)
	}
	if cond.Operator != gorulesengine.OperatorEqual {
		t.Errorf("Expected operator OperatorEqual, got '%s'", cond.Operator)
	}
	if cond.Value != 18 {
		t.Errorf("Expected value 18, got %v", cond.Value)
	}
}

func TestNotEqual(t *testing.T) {
	cond := gorulesengine.NotEqual("status", "inactive")
	if cond.Operator != gorulesengine.OperatorNotEqual {
		t.Errorf("Expected operator OperatorNotEqual, got '%s'", cond.Operator)
	}
}

func TestGreaterThan(t *testing.T) {
	cond := gorulesengine.GreaterThan("score", 100)
	if cond.Operator != gorulesengine.OperatorGreaterThan {
		t.Errorf("Expected operator OperatorGreaterThan, got '%s'", cond.Operator)
	}
}

func TestGreaterThanInclusive(t *testing.T) {
	cond := gorulesengine.GreaterThanInclusive("score", 100)
	if cond.Operator != gorulesengine.OperatorGreaterThanInclusive {
		t.Errorf("Expected operator OperatorGreaterThanInclusive, got '%s'", cond.Operator)
	}
}

func TestLessThan(t *testing.T) {
	cond := gorulesengine.LessThan("age", 65)
	if cond.Operator != gorulesengine.OperatorLessThan {
		t.Errorf("Expected operator OperatorLessThan, got '%s'", cond.Operator)
	}
}

func TestLessThanInclusive(t *testing.T) {
	cond := gorulesengine.LessThanInclusive("age", 65)
	if cond.Operator != gorulesengine.OperatorLessThanInclusive {
		t.Errorf("Expected operator OperatorLessThanInclusive, got '%s'", cond.Operator)
	}
}

func TestIn(t *testing.T) {
	cond := gorulesengine.In("country", []string{"US", "CA", "UK"})
	if cond.Operator != gorulesengine.OperatorIn {
		t.Errorf("Expected operator OperatorIn, got '%s'", cond.Operator)
	}
}

func TestNotIn(t *testing.T) {
	cond := gorulesengine.NotIn("country", []string{"US", "CA"})
	if cond.Operator != gorulesengine.OperatorNotIn {
		t.Errorf("Expected operator OperatorNotIn, got '%s'", cond.Operator)
	}
}

func TestContains(t *testing.T) {
	cond := gorulesengine.Contains("tags", "premium")
	if cond.Operator != gorulesengine.OperatorContains {
		t.Errorf("Expected operator OperatorContains, got '%s'", cond.Operator)
	}
}

func TestNotContains(t *testing.T) {
	cond := gorulesengine.NotContains("tags", "banned")
	if cond.Operator != gorulesengine.OperatorNotContains {
		t.Errorf("Expected operator OperatorNotContains, got '%s'", cond.Operator)
	}
}

func TestRegex(t *testing.T) {
	cond := gorulesengine.Regex("email", "^[a-z]+@[a-z]+\\.[a-z]+$")
	if cond.Operator != gorulesengine.OperatorRegex {
		t.Errorf("Expected operator OperatorRegex, got '%s'", cond.Operator)
	}
}

// Test ConditionSet helpers

func TestAll(t *testing.T) {
	condSet := gorulesengine.All(
		gorulesengine.Equal("age", 18),
		gorulesengine.GreaterThan("score", 100),
	)

	if condSet.All == nil || len(condSet.All) != 2 {
		t.Fatal("Expected two conditions in All")
	}

	if condSet.All[0].Condition == nil {
		t.Fatal("Expected first condition to be set")
	}

	if condSet.All[0].Condition.Operator != gorulesengine.OperatorEqual {
		t.Errorf("Expected first operator to be OperatorEqual")
	}

	if condSet.All[1].Condition.Operator != gorulesengine.OperatorGreaterThan {
		t.Errorf("Expected second operator to be OperatorGreaterThan")
	}
}

func TestAny(t *testing.T) {
	condSet := gorulesengine.Any(
		gorulesengine.Equal("status", "active"),
		gorulesengine.Equal("status", "pending"),
	)

	if condSet.Any == nil || len(condSet.Any) != 2 {
		t.Fatal("Expected two conditions in Any")
	}
}

func TestNone(t *testing.T) {
	condSet := gorulesengine.None(
		gorulesengine.Equal("status", "banned"),
		gorulesengine.Equal("status", "deleted"),
	)

	if condSet.None == nil || len(condSet.None) != 2 {
		t.Fatal("Expected two conditions in None")
	}
}

func TestAllSets(t *testing.T) {
	set1 := gorulesengine.All(gorulesengine.Equal("age", 18))
	set2 := gorulesengine.All(gorulesengine.GreaterThan("score", 100))

	condSet := gorulesengine.AllSets(set1, set2)

	if condSet.All == nil || len(condSet.All) != 2 {
		t.Fatal("Expected two condition sets in All")
	}

	if condSet.All[0].SubSet == nil {
		t.Fatal("Expected first SubSet to be set")
	}

	if condSet.All[1].SubSet == nil {
		t.Fatal("Expected second SubSet to be set")
	}
}

func TestAnySets(t *testing.T) {
	set1 := gorulesengine.All(gorulesengine.Equal("age", 18))
	set2 := gorulesengine.All(gorulesengine.GreaterThan("score", 100))

	condSet := gorulesengine.AnySets(set1, set2)

	if condSet.Any == nil || len(condSet.Any) != 2 {
		t.Fatal("Expected two condition sets in Any")
	}

	if condSet.Any[0].SubSet == nil {
		t.Fatal("Expected first SubSet to be set")
	}
}

func TestNoneSets(t *testing.T) {
	set1 := gorulesengine.All(gorulesengine.Equal("status", "banned"))
	set2 := gorulesengine.All(gorulesengine.Equal("status", "deleted"))

	condSet := gorulesengine.NoneSets(set1, set2)

	if condSet.None == nil || len(condSet.None) != 2 {
		t.Fatal("Expected two condition sets in None")
	}

	if condSet.None[0].SubSet == nil {
		t.Fatal("Expected first SubSet to be set")
	}

	if condSet.None[1].SubSet == nil {
		t.Fatal("Expected second SubSet to be set")
	}
}

// Integration test with Builder and Engine

func TestBuilderWithEngine(t *testing.T) {
	engine := gorulesengine.NewEngine()

	handler := func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
		fmt.Println("Success callback executed")
		return nil
	}

	engine.RegisterCallback("success-cb", handler)

	rule := gorulesengine.NewRuleBuilder().
		WithName("adult-rule").
		WithConditions(
			gorulesengine.ConditionNode{
				SubSet: &gorulesengine.ConditionSet{
					All: []gorulesengine.ConditionNode{
						{
							Condition: gorulesengine.GreaterThanInclusive("age", 18),
						},
						{
							Condition: gorulesengine.Equal("status", "active"),
						},
					},
				},
			},
		).
		WithEvent("adult-access", map[string]interface{}{
			"message": "Access granted",
		}).
		WithOnSuccess("success-cb").
		Build()

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("age", 21)
	almanac.AddFact("status", "active")

	_, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Engine run failed: %v", err)
	}

	events := almanac.GetSuccessEvents()
	if len(events) != 1 {
		t.Fatalf("Expected 1 success event, got %d", len(events))
	}

	if events[0].Type != "adult-access" {
		t.Errorf("Expected event type 'adult-access', got '%s'", events[0].Type)
	}
}

func TestBuilderWithNestedConditions(t *testing.T) {
	engine := gorulesengine.NewEngine()

	handlerFailure := func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
		fmt.Println("Failed callback executed")
		return nil
	}

	handlerSuccess := func(event gorulesengine.Event, almanac *gorulesengine.Almanac, ruleResult gorulesengine.RuleResult) error {
		fmt.Println("Success callback executed")
		return nil
	}

	engine.RegisterCallback("failed-cb", handlerFailure)
	engine.RegisterCallback("success-cb", handlerSuccess)

	rule := gorulesengine.NewRuleBuilder().
		WithName("complex-rule").
		WithConditions(
			gorulesengine.ConditionNode{
				SubSet: &gorulesengine.ConditionSet{
					All: []gorulesengine.ConditionNode{
						{
							Condition: gorulesengine.Equal("type", "premium"),
						},
						{
							Condition: gorulesengine.GreaterThan("credits", 100),
						},
					},
				},
			},
		).
		WithEvent("access-granted", nil).
		WithOnFailure("failed-cb").
		WithOnSuccess("success-cb").
		Build()

	engine.AddRule(rule)

	// Test with regular user
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("type", "regular")
	almanac.AddFact("credits", 150)

	_, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Engine run failed: %v", err)
	}

	events := almanac.GetFailureEvents()
	if len(events) != 1 {
		t.Fatalf("Expected 1 failure event, got %d", len(events))
	}

	// Test with premium user
	almanac2 := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac2.AddFact("type", "premium")
	almanac2.AddFact("credits", 150)

	_, err = engine.Run(almanac2)
	if err != nil {
		t.Fatalf("Engine run failed: %v", err)
	}

	events2 := almanac2.GetSuccessEvents()
	if len(events2) != 1 {
		t.Fatalf("Expected 1 success event, got %d", len(events2))
	}
}

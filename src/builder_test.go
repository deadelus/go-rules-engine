package gorulesengine_test

import (
	"testing"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func TestNewRuleBuilder(t *testing.T) {
	t.Run("creates a new rule builder", func(t *testing.T) {
		builder := gre.NewRuleBuilder()
		if builder == nil {
			t.Fatal("Expected builder to be created")
		}

		rule := builder.Build()
		if rule == nil {
			t.Fatal("Expected rule to be created")
		}
	})
}

func TestRuleBuilder_WithName(t *testing.T) {
	t.Run("sets the rule name", func(t *testing.T) {
		builder := gre.NewRuleBuilder()
		rule := builder.WithName("test-rule").Build()

		if rule.Name != "test-rule" {
			t.Errorf("Expected rule name 'test-rule', got '%s'", rule.Name)
		}
	})

	t.Run("returns builder for chaining", func(t *testing.T) {
		builder := gre.NewRuleBuilder()
		result := builder.WithName("test")

		if result != builder {
			t.Error("Expected WithName to return the same builder instance")
		}
	})
}

func TestRuleBuilder_WithPriority(t *testing.T) {
	t.Run("sets the rule priority", func(t *testing.T) {
		builder := gre.NewRuleBuilder()
		rule := builder.WithPriority(10).Build()

		if rule.Priority != 10 {
			t.Errorf("Expected rule priority 10, got %d", rule.Priority)
		}
	})

	t.Run("returns builder for chaining", func(t *testing.T) {
		builder := gre.NewRuleBuilder()
		result := builder.WithPriority(5)

		if result != builder {
			t.Error("Expected WithPriority to return the same builder instance")
		}
	})
}

func TestRuleBuilder_WithConditions(t *testing.T) {
	t.Run("sets conditions from a condition node", func(t *testing.T) {
		builder := gre.NewRuleBuilder()
		condition := gre.Equal("age", 18)
		node := gre.ConditionNode{Condition: condition}

		rule := builder.WithConditions(node).Build()

		if len(rule.Conditions.All) != 1 {
			t.Errorf("Expected 1 condition in All, got %d", len(rule.Conditions.All))
		}
	})

	t.Run("sets conditions from a condition set", func(t *testing.T) {
		builder := gre.NewRuleBuilder()
		condSet := gre.All(
			gre.Equal("age", 18),
			gre.Equal("country", "FR"),
		)
		node := gre.ConditionNode{SubSet: &condSet}

		rule := builder.WithConditions(node).Build()

		if len(rule.Conditions.All) != 2 {
			t.Errorf("Expected 2 conditions in All, got %d", len(rule.Conditions.All))
		}
	})

	t.Run("returns builder for chaining", func(t *testing.T) {
		builder := gre.NewRuleBuilder()
		condition := gre.Equal("age", 18)
		node := gre.ConditionNode{Condition: condition}

		result := builder.WithConditions(node)

		if result != builder {
			t.Error("Expected WithConditions to return the same builder instance")
		}
	})
}

func TestRuleBuilder_WithEvents(t *testing.T) {
	t.Run("sets rule events", func(t *testing.T) {
		builder := gre.NewRuleBuilder()

		rule := builder.WithOnSuccess("success-event").WithOnFailure("failure-event").Build()

		if len(rule.OnSuccess) != 1 || rule.OnSuccess[0].Name != "success-event" {
			t.Error("Expected OnSuccess event to be set")
		}
		if len(rule.OnFailure) != 1 || rule.OnFailure[0].Name != "failure-event" {
			t.Error("Expected OnFailure event to be set")
		}
	})

	t.Run("returns builder for chaining", func(t *testing.T) {
		builder := gre.NewRuleBuilder()

		result := builder.WithOnSuccess("test")

		if result != builder {
			t.Error("Expected WithEvents to return the same builder instance")
		}
	})

	t.Run("sets detailed rule events", func(t *testing.T) {
		builder := gre.NewRuleBuilder()
		successEvent := gre.RuleEvent{
			Name:   "success-event",
			Params: map[string]interface{}{"foo": "bar"},
		}
		failureEvent := gre.RuleEvent{
			Name:   "failure-event",
			Params: map[string]interface{}{"baz": "qux"},
		}

		rule := builder.
			WithOnSuccessEvent(successEvent).
			WithOnFailureEvent(failureEvent).
			Build()

		if len(rule.OnSuccess) != 1 || rule.OnSuccess[0].Name != "success-event" {
			t.Error("Expected OnSuccess event to be set")
		}
		if rule.OnSuccess[0].Params["foo"] != "bar" {
			t.Error("Expected success event params to be set")
		}

		if len(rule.OnFailure) != 1 || rule.OnFailure[0].Name != "failure-event" {
			t.Error("Expected OnFailure event to be set")
		}
		if rule.OnFailure[0].Params["baz"] != "qux" {
			t.Error("Expected failure event params to be set")
		}
	})
}

func TestRuleBuilder_Build(t *testing.T) {
	t.Run("builds a complete rule", func(t *testing.T) {
		rule := gre.NewRuleBuilder().
			WithName("adult-rule").
			WithPriority(10).
			WithConditions(gre.ConditionNode{
				Condition: gre.GreaterThanInclusive("age", 18),
			}).
			WithOnSuccess("adult-event").
			WithOnFailure("failure-event").
			Build()

		if rule.Name != "adult-rule" {
			t.Errorf("Expected name 'adult-rule', got '%s'", rule.Name)
		}
		if rule.Priority != 10 {
			t.Errorf("Expected priority 10, got %d", rule.Priority)
		}
		if len(rule.Conditions.All) == 0 {
			t.Error("Expected conditions to be set")
		}
		if len(rule.OnSuccess) == 0 {
			t.Error("Expected events to be set")
		}
	})
}

func TestEqual(t *testing.T) {
	t.Run("creates equal condition", func(t *testing.T) {
		cond := gre.Equal("age", 18)

		if cond.Fact != "age" {
			t.Errorf("Expected fact 'age', got '%s'", cond.Fact)
		}
		if cond.Operator != gre.OperatorEqual {
			t.Errorf("Expected operator 'equal', got '%s'", cond.Operator)
		}
		if cond.Value != 18 {
			t.Errorf("Expected value 18, got %v", cond.Value)
		}
	})
}

func TestNotEqual(t *testing.T) {
	t.Run("creates not equal condition", func(t *testing.T) {
		cond := gre.NotEqual("status", "inactive")

		if cond.Fact != "status" {
			t.Errorf("Expected fact 'status', got '%s'", cond.Fact)
		}
		if cond.Operator != gre.OperatorNotEqual {
			t.Errorf("Expected operator 'not_equal', got '%s'", cond.Operator)
		}
		if cond.Value != "inactive" {
			t.Errorf("Expected value 'inactive', got %v", cond.Value)
		}
	})
}

func TestGreaterThan(t *testing.T) {
	t.Run("creates greater than condition", func(t *testing.T) {
		cond := gre.GreaterThan("age", 18)

		if cond.Fact != "age" {
			t.Errorf("Expected fact 'age', got '%s'", cond.Fact)
		}
		if cond.Operator != gre.OperatorGreaterThan {
			t.Errorf("Expected operator 'greater_than', got '%s'", cond.Operator)
		}
		if cond.Value != 18 {
			t.Errorf("Expected value 18, got %v", cond.Value)
		}
	})
}

func TestGreaterThanInclusive(t *testing.T) {
	t.Run("creates greater than inclusive condition", func(t *testing.T) {
		cond := gre.GreaterThanInclusive("age", 18)

		if cond.Fact != "age" {
			t.Errorf("Expected fact 'age', got '%s'", cond.Fact)
		}
		if cond.Operator != gre.OperatorGreaterThanInclusive {
			t.Errorf("Expected operator 'greater_than_inclusive', got '%s'", cond.Operator)
		}
		if cond.Value != 18 {
			t.Errorf("Expected value 18, got %v", cond.Value)
		}
	})
}

func TestLessThan(t *testing.T) {
	t.Run("creates less than condition", func(t *testing.T) {
		cond := gre.LessThan("age", 65)

		if cond.Fact != "age" {
			t.Errorf("Expected fact 'age', got '%s'", cond.Fact)
		}
		if cond.Operator != gre.OperatorLessThan {
			t.Errorf("Expected operator 'less_than', got '%s'", cond.Operator)
		}
		if cond.Value != 65 {
			t.Errorf("Expected value 65, got %v", cond.Value)
		}
	})
}

func TestLessThanInclusive(t *testing.T) {
	t.Run("creates less than inclusive condition", func(t *testing.T) {
		cond := gre.LessThanInclusive("age", 65)

		if cond.Fact != "age" {
			t.Errorf("Expected fact 'age', got '%s'", cond.Fact)
		}
		if cond.Operator != gre.OperatorLessThanInclusive {
			t.Errorf("Expected operator 'less_than_inclusive', got '%s'", cond.Operator)
		}
		if cond.Value != 65 {
			t.Errorf("Expected value 65, got %v", cond.Value)
		}
	})
}

func TestIn(t *testing.T) {
	t.Run("creates in condition", func(t *testing.T) {
		values := []string{"FR", "US", "UK"}
		cond := gre.In("country", values)

		if cond.Fact != "country" {
			t.Errorf("Expected fact 'country', got '%s'", cond.Fact)
		}
		if cond.Operator != gre.OperatorIn {
			t.Errorf("Expected operator 'in', got '%s'", cond.Operator)
		}
	})
}

func TestNotIn(t *testing.T) {
	t.Run("creates not in condition", func(t *testing.T) {
		values := []string{"banned", "suspended"}
		cond := gre.NotIn("status", values)

		if cond.Fact != "status" {
			t.Errorf("Expected fact 'status', got '%s'", cond.Fact)
		}
		if cond.Operator != gre.OperatorNotIn {
			t.Errorf("Expected operator 'not_in', got '%s'", cond.Operator)
		}
	})
}

func TestContains(t *testing.T) {
	t.Run("creates contains condition", func(t *testing.T) {
		cond := gre.Contains("tags", "premium")

		if cond.Fact != "tags" {
			t.Errorf("Expected fact 'tags', got '%s'", cond.Fact)
		}
		if cond.Operator != gre.OperatorContains {
			t.Errorf("Expected operator 'contains', got '%s'", cond.Operator)
		}
		if cond.Value != "premium" {
			t.Errorf("Expected value 'premium', got %v", cond.Value)
		}
	})
}

func TestNotContains(t *testing.T) {
	t.Run("creates not contains condition", func(t *testing.T) {
		cond := gre.NotContains("tags", "spam")

		if cond.Fact != "tags" {
			t.Errorf("Expected fact 'tags', got '%s'", cond.Fact)
		}
		if cond.Operator != gre.OperatorNotContains {
			t.Errorf("Expected operator 'not_contains', got '%s'", cond.Operator)
		}
		if cond.Value != "spam" {
			t.Errorf("Expected value 'spam', got %v", cond.Value)
		}
	})
}

func TestRegex(t *testing.T) {
	t.Run("creates regex condition", func(t *testing.T) {
		pattern := "^[A-Z]{2}[0-9]{3}$"
		cond := gre.Regex("code", pattern)

		if cond.Fact != "code" {
			t.Errorf("Expected fact 'code', got '%s'", cond.Fact)
		}
		if cond.Operator != gre.OperatorRegex {
			t.Errorf("Expected operator 'regex', got '%s'", cond.Operator)
		}
		if cond.Value != pattern {
			t.Errorf("Expected pattern '%s', got %v", pattern, cond.Value)
		}
	})
}

func TestAll(t *testing.T) {
	t.Run("creates all condition set with single condition", func(t *testing.T) {
		condSet := gre.All(
			gre.Equal("age", 18),
		)

		if len(condSet.All) != 1 {
			t.Errorf("Expected 1 condition, got %d", len(condSet.All))
		}
	})

	t.Run("creates all condition set with multiple conditions", func(t *testing.T) {
		condSet := gre.All(
			gre.Equal("age", 18),
			gre.Equal("country", "FR"),
			gre.Equal("active", true),
		)

		if len(condSet.All) != 3 {
			t.Errorf("Expected 3 conditions, got %d", len(condSet.All))
		}

		if condSet.All[0].Condition.Fact != "age" {
			t.Error("Expected first condition to be 'age'")
		}
		if condSet.All[1].Condition.Fact != "country" {
			t.Error("Expected second condition to be 'country'")
		}
		if condSet.All[2].Condition.Fact != "active" {
			t.Error("Expected third condition to be 'active'")
		}
	})
}

func TestAny(t *testing.T) {
	t.Run("creates any condition set with single condition", func(t *testing.T) {
		condSet := gre.Any(
			gre.Equal("status", "premium"),
		)

		if len(condSet.Any) != 1 {
			t.Errorf("Expected 1 condition, got %d", len(condSet.Any))
		}
	})

	t.Run("creates any condition set with multiple conditions", func(t *testing.T) {
		condSet := gre.Any(
			gre.Equal("status", "premium"),
			gre.Equal("status", "vip"),
		)

		if len(condSet.Any) != 2 {
			t.Errorf("Expected 2 conditions, got %d", len(condSet.Any))
		}

		if condSet.Any[0].Condition.Fact != "status" {
			t.Error("Expected first condition to be 'status'")
		}
	})
}

func TestNone(t *testing.T) {
	t.Run("creates none condition set with single condition", func(t *testing.T) {
		condSet := gre.None(
			gre.Equal("banned", true),
		)

		if len(condSet.None) != 1 {
			t.Errorf("Expected 1 condition, got %d", len(condSet.None))
		}
	})

	t.Run("creates none condition set with multiple conditions", func(t *testing.T) {
		condSet := gre.None(
			gre.Equal("banned", true),
			gre.Equal("suspended", true),
		)

		if len(condSet.None) != 2 {
			t.Errorf("Expected 2 conditions, got %d", len(condSet.None))
		}

		if condSet.None[0].Condition.Fact != "banned" {
			t.Error("Expected first condition to be 'banned'")
		}
	})
}

func TestAllSets(t *testing.T) {
	t.Run("creates all sets with single condition set", func(t *testing.T) {
		set1 := gre.All(gre.Equal("age", 18))
		condSet := gre.AllSets(set1)

		if len(condSet.All) != 1 {
			t.Errorf("Expected 1 nested set, got %d", len(condSet.All))
		}
		if condSet.All[0].SubSet == nil {
			t.Error("Expected SubSet to be set")
		}
	})

	t.Run("creates all sets with multiple condition sets", func(t *testing.T) {
		set1 := gre.All(gre.Equal("age", 18))
		set2 := gre.Any(
			gre.Equal("country", "FR"),
			gre.Equal("country", "US"),
		)

		condSet := gre.AllSets(set1, set2)

		if len(condSet.All) != 2 {
			t.Errorf("Expected 2 nested sets, got %d", len(condSet.All))
		}

		if condSet.All[0].SubSet == nil {
			t.Error("Expected first SubSet to be set")
		}
		if condSet.All[1].SubSet == nil {
			t.Error("Expected second SubSet to be set")
		}

		if len(condSet.All[0].SubSet.All) != 1 {
			t.Error("Expected first set to have 1 All condition")
		}
		if len(condSet.All[1].SubSet.Any) != 2 {
			t.Error("Expected second set to have 2 Any conditions")
		}
	})
}

func TestAnySets(t *testing.T) {
	t.Run("creates any sets with single condition set", func(t *testing.T) {
		set1 := gre.Any(gre.Equal("premium", true))
		condSet := gre.AnySets(set1)

		if len(condSet.Any) != 1 {
			t.Errorf("Expected 1 nested set, got %d", len(condSet.Any))
		}
		if condSet.Any[0].SubSet == nil {
			t.Error("Expected SubSet to be set")
		}
	})

	t.Run("creates any sets with multiple condition sets", func(t *testing.T) {
		set1 := gre.All(gre.Equal("premium", true))
		set2 := gre.All(gre.Equal("vip", true))

		condSet := gre.AnySets(set1, set2)

		if len(condSet.Any) != 2 {
			t.Errorf("Expected 2 nested sets, got %d", len(condSet.Any))
		}

		if condSet.Any[0].SubSet == nil || condSet.Any[1].SubSet == nil {
			t.Error("Expected SubSets to be set")
		}
	})
}

func TestNoneSets(t *testing.T) {
	t.Run("creates none sets with single condition set", func(t *testing.T) {
		set1 := gre.All(gre.Equal("banned", true))
		condSet := gre.NoneSets(set1)

		if len(condSet.None) != 1 {
			t.Errorf("Expected 1 nested set, got %d", len(condSet.None))
		}
		if condSet.None[0].SubSet == nil {
			t.Error("Expected SubSet to be set")
		}
	})

	t.Run("creates none sets with multiple condition sets", func(t *testing.T) {
		set1 := gre.All(gre.Equal("banned", true))
		set2 := gre.All(gre.Equal("suspended", true))

		condSet := gre.NoneSets(set1, set2)

		if len(condSet.None) != 2 {
			t.Errorf("Expected 2 nested sets, got %d", len(condSet.None))
		}

		if condSet.None[0].SubSet == nil || condSet.None[1].SubSet == nil {
			t.Error("Expected SubSets to be set")
		}
	})
}

func TestBuilderIntegration(t *testing.T) {
	t.Run("builds complex rule with nested conditions", func(t *testing.T) {
		// Build a complex rule: (age >= 18 AND country in [FR, US]) OR premium = true
		ageAndCountry := gre.All(
			gre.GreaterThanInclusive("age", 18),
			gre.In("country", []string{"FR", "US"}),
		)
		premiumSet := gre.All(gre.Equal("premium", true))

		conditions := gre.AnySets(ageAndCountry, premiumSet)

		rule := gre.NewRuleBuilder().
			WithName("access-rule").
			WithPriority(100).
			WithConditions(gre.ConditionNode{SubSet: &conditions}).
			WithOnSuccess("grant-access").
			WithOnFailure("deny-access").
			Build()

		if rule.Name != "access-rule" {
			t.Error("Rule name not set correctly")
		}
		if rule.Priority != 100 {
			t.Error("Rule priority not set correctly")
		}
		if len(rule.Conditions.Any) != 2 {
			t.Errorf("Expected 2 Any conditions, got %d", len(rule.Conditions.Any))
		}
		if len(rule.OnSuccess) != 1 || len(rule.OnFailure) != 1 {
			t.Error("Events not set correctly")
		}
	})
}

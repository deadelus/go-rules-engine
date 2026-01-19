package gorulesengine_test

import (
	"testing"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func TestEngineAuditTrace(t *testing.T) {
	t.Run("runs with audit trace enabled", func(t *testing.T) {
		engine := gre.NewEngine(gre.WithAuditTrace())
		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 25)

		rule := &gre.Rule{
			Name: "audit-rule",
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{
						Condition: &gre.Condition{
							Fact:     "age",
							Operator: "greater_than",
							Value:    float64(18),
						},
					},
				},
			},
		}

		engine.AddRule(rule)
		e, err := engine.Run(almanac)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		results := e.Results()
		res, ok := results["audit-rule"]
		if !ok {
			t.Fatal("Expected result for 'audit-rule'")
		}

		if res.Conditions == nil {
			t.Error("Expected detailed conditions in result when audit trace is enabled")
		}

		if res.Conditions.Result != true {
			t.Errorf("Expected condition result to be true, got %v", res.Conditions.Result)
		}

		if len(res.Conditions.Results) != 1 {
			t.Errorf("Expected 1 condition in audit trace, got %d", len(res.Conditions.Results))
		}
	})

	t.Run("runs with audit trace disabled (default)", func(t *testing.T) {
		engine := gre.NewEngine() // Default no audit
		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 25)

		rule := &gre.Rule{
			Name: "no-audit-rule",
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{
						Condition: &gre.Condition{
							Fact:     "age",
							Operator: "greater_than",
							Value:    float64(18),
						},
					},
				},
			},
		}

		engine.AddRule(rule)
		e, err := engine.Run(almanac)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		results := e.Results()
		res := results["no-audit-rule"]

		if res.Conditions != nil {
			t.Error("Expected no detailed conditions in result when audit trace is disabled")
		}
	})

	t.Run("WithoutAuditTrace explicitly disables it", func(t *testing.T) {
		engine := gre.NewEngine(gre.WithAuditTrace(), gre.WithoutAuditTrace())
		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 25)

		rule := &gre.Rule{
			Name: "explicit-no-audit",
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{
						Condition: &gre.Condition{
							Fact:     "age",
							Operator: "greater_than",
							Value:    float64(18),
						},
					},
				},
			},
		}

		engine.AddRule(rule)
		e, err := engine.Run(almanac)
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		res := e.Results()["explicit-no-audit"]
		if res.Conditions != nil {
			t.Error("Expected no detailed conditions in result after WithAuditTrace was followed by WithoutAuditTrace")
		}
	})
}

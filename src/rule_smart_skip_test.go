package gorulesengine_test

import (
	"testing"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func TestSmartSkip(t *testing.T) {
	callCount := 0
	dynamicFact := func(params map[string]interface{}) interface{} {
		callCount++
		return 25
	}

	engine := gre.NewEngine(gre.WithSmartSkip())

	rule1 := &gre.Rule{
		Name:     "missing-fact-rule",
		Priority: 10,
		Conditions: gre.ConditionSet{
			All: []gre.ConditionNode{
				{
					Condition: &gre.Condition{
						Fact:     "nonexistent",
						Operator: "equal",
						Value:    10,
					},
				},
			},
		},
	}
	engine.AddRule(rule1)

	rule2 := &gre.Rule{
		Name:     "existing-fact-rule",
		Priority: 5,
		Conditions: gre.ConditionSet{
			All: []gre.ConditionNode{
				{
					Condition: &gre.Condition{
						Fact:     "age",
						Operator: "greater_than",
						Value:    18,
					},
				},
			},
		},
	}
	engine.AddRule(rule2)

	almanac := gre.NewAlmanac()
	almanac.AddFact("age", dynamicFact)

	e, err := engine.Run(almanac)

	results := e.ReduceResults()

	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// rule1 should be skipped (result false)
	if results["missing-fact-rule"] != false {
		t.Errorf("Expected missing-fact-rule to be false (skipped), got %v", results["missing-fact-rule"])
	}

	// rule2 should pass
	if results["existing-fact-rule"] != true {
		t.Errorf("Expected existing-fact-rule to be true, got %v", results["existing-fact-rule"])
	}

	if callCount != 1 {
		t.Errorf("Expected 1 call to dynamic fact, got %d", callCount)
	}
}

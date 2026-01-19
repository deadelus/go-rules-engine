package gorulesengine_test

import (
	"testing"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func TestCompileAndRequiredFacts(t *testing.T) {
	cond1 := &gre.Condition{Fact: "age", Operator: "greater_than", Value: 18}
	cond2 := &gre.Condition{Fact: "score", Operator: "less_than", Value: 100}

	cs := &gre.ConditionSet{
		All: []gre.ConditionNode{
			{Condition: cond1},
			{
				SubSet: &gre.ConditionSet{
					Any: []gre.ConditionNode{
						{Condition: cond2},
						{Condition: &gre.Condition{Fact: "name", Operator: "equal", Value: "test"}},
					},
				},
			},
		},
	}

	// Test Compile
	err := cs.Compile()
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// Test GetRequiredFacts
	facts := cs.GetRequiredFacts()
	expected := []gre.FactID{"age", "score", "name"}

	if len(facts) != len(expected) {
		t.Errorf("Expected %d facts, got %d", len(expected), len(facts))
	}

	// Check if all expected facts are present
	factMap := make(map[gre.FactID]bool)
	for _, f := range facts {
		factMap[f] = true
	}

	for _, f := range expected {
		if !factMap[f] {
			t.Errorf("Fact %s not found in required facts", f)
		}
	}
}

func TestConditionCompileNestedErrors(t *testing.T) {
	badCond := &gre.Condition{Value: func() {}}

	// Error in All
	csAll := &gre.ConditionSet{All: []gre.ConditionNode{{Condition: badCond}}}
	if err := csAll.Compile(); err == nil {
		t.Error("Expected error compiling All")
	}

	// Error in Any
	csAny := &gre.ConditionSet{Any: []gre.ConditionNode{{Condition: badCond}}}
	if err := csAny.Compile(); err == nil {
		t.Error("Expected error compiling Any")
	}

	// Error in None
	csNone := &gre.ConditionSet{None: []gre.ConditionNode{{Condition: badCond}}}
	if err := csNone.Compile(); err == nil {
		t.Error("Expected error compiling None")
	}
}

func TestConditionCompileError(t *testing.T) {
	// A condition with a value that cannot be marshaled to JSON
	// will fail during GetCacheKey in Compile if we were using it there,
	// but Compile calls GetCacheKey to pre-calculate.
	badCond := &gre.Condition{
		Fact:     "f",
		Operator: "equal",
		Value:    func() {}, // functions cannot be marshaled
	}

	err := badCond.Compile()
	if err == nil {
		t.Error("Expected error during compile of bad condition")
	}
}

func TestConditionSetCompileNestedSubSetErrors(t *testing.T) {
	badCond := &gre.Condition{Value: func() {}}
	badSubSet := &gre.ConditionSet{All: []gre.ConditionNode{{Condition: badCond}, {Condition: badCond}}}

	// Compile error in subset of All
	csAll := &gre.ConditionSet{All: []gre.ConditionNode{{SubSet: badSubSet}}}
	if err := csAll.Compile(); err == nil {
		t.Error("Expected error compiling All subset")
	}

	// Compile error in subset of Any
	csAny := &gre.ConditionSet{Any: []gre.ConditionNode{{SubSet: badSubSet}}}
	if err := csAny.Compile(); err == nil {
		t.Error("Expected error compiling Any subset")
	}

	// Compile error in subset of None
	csNone := &gre.ConditionSet{None: []gre.ConditionNode{{SubSet: badSubSet}}}
	if err := csNone.Compile(); err == nil {
		t.Error("Expected error compiling None subset")
	}
}

func TestConditionNode_GetRequiredFacts_Empty(t *testing.T) {
	node := gre.ConditionNode{} // both Condition and SubSet are nil
	facts := node.GetRequiredFacts()
	if len(facts) != 0 {
		t.Errorf("Expected 0 facts, got %d", len(facts))
	}
}

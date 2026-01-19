package gorulesengine_test

import (
	"encoding/json"
	"testing"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

// badMarble is a type that fails to marshal on the second call
type badMarble struct {
	called bool
}

func (b *badMarble) MarshalJSON() ([]byte, error) {
	if b.called {
		return nil, &json.UnsupportedValueError{}
	}
	b.called = true
	return []byte("\"ok\""), nil
}

func TestReorderNodes(t *testing.T) {
	almanac := gre.NewAlmanac(gre.WithAlmanacConditionCaching())

	c1 := &gre.Condition{Fact: "f1", Operator: "equal", Value: 1}
	_ = c1.Compile()
	key1, _ := c1.GetCacheKey()
	almanac.SetConditionResultCache(key1, &gre.ConditionResult{
		Fact:     "f1",
		Operator: "equal",
		Value:    1,
		Result:   true,
	})

	c2 := &gre.Condition{Fact: "f2", Operator: "equal", Value: 2}
	_ = c2.Compile()

	cs := &gre.ConditionSet{
		All: []gre.ConditionNode{
			{Condition: c2}, // Not cached
			{Condition: c1}, // Cached
		},
	}

	fact1 := gre.NewFact("f1", 1)
	fact2 := gre.NewFact("f2", 2)
	almanac.AddFacts(&fact1, &fact2)

	csr, err := cs.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	if !csr.Result {
		t.Error("Expected result to be true")
	}
}

func TestConditionSetReorderErrorBranches(t *testing.T) {
	almanac := gre.NewAlmanac(gre.WithAlmanacConditionCaching())

	// 1. All reorder fail
	bmAll := &badMarble{}
	csAll := &gre.ConditionSet{
		All: []gre.ConditionNode{
			{Condition: &gre.Condition{Fact: "f", Operator: "equal", Value: bmAll}},
			{Condition: &gre.Condition{Fact: "f", Operator: "equal", Value: "stable"}},
		},
	}
	_, err := csAll.Evaluate(almanac)
	if err == nil {
		t.Error("Expected error for All reorder")
	}

	// 2. Any reorder fail
	bmAny := &badMarble{}
	csAny := &gre.ConditionSet{
		Any: []gre.ConditionNode{
			{Condition: &gre.Condition{Fact: "f", Operator: "equal", Value: bmAny}},
			{Condition: &gre.Condition{Fact: "f", Operator: "equal", Value: "stable"}},
		},
	}
	_, err = csAny.Evaluate(almanac)
	if err == nil {
		t.Error("Expected error for Any reorder")
	}

	// 3. None reorder fail
	bmNone := &badMarble{}
	csNone := &gre.ConditionSet{
		None: []gre.ConditionNode{
			{Condition: &gre.Condition{Fact: "f", Operator: "equal", Value: bmNone}},
			{Condition: &gre.Condition{Fact: "f", Operator: "equal", Value: "stable"}},
		},
	}
	_, err = csNone.Evaluate(almanac)
	if err == nil {
		t.Error("Expected error for None reorder")
	}
}

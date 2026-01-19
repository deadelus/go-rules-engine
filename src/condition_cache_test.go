package gorulesengine_test

import (
	"testing"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func TestConditionCaching(t *testing.T) {
	callCount := 0
	dynamicFact := func(params map[string]interface{}) interface{} {
		callCount++
		return 25
	}

	// Create almanac with condition caching enabled
	almanac := gre.NewAlmanac(gre.WithAlmanacConditionCaching())

	// Add fact WITHOUT fact caching to specifically test condition caching
	almanac.AddFact("age", dynamicFact, gre.WithoutCache())

	condition := &gre.Condition{
		Fact:     "age",
		Operator: "greater_than",
		Value:    18,
	}

	// First evaluation: should call the dynamic fact
	cr, err := condition.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Evaluate failed: %v", err)
	}
	res := cr.Result

	if !res {
		t.Errorf("Expected true, got false")
	}
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Second evaluation: should use condition cache and NOT call the dynamic fact
	cr, err = condition.Evaluate(almanac)
	if err != nil {
		t.Fatalf("Evaluate failed: %v", err)
	}
	res = cr.Result

	if !res {
		t.Errorf("Expected true, got false")
	}
	if callCount != 1 {
		t.Errorf("Expected callCount to remain 1 due to condition caching, got %d", callCount)
	}
}

func TestConditionCachingDisabled(t *testing.T) {
	callCount := 0
	dynamicFact := func(params map[string]interface{}) interface{} {
		callCount++
		return 25
	}

	// Create almanac without condition caching (default)
	almanac := gre.NewAlmanac()

	// Add fact WITHOUT fact caching
	almanac.AddFact("age", dynamicFact, gre.WithoutCache())

	condition := &gre.Condition{
		Fact:     "age",
		Operator: "greater_than",
		Value:    18,
	}

	// First evaluation
	condition.Evaluate(almanac)
	if callCount != 1 {
		t.Errorf("Expected 1 call, got %d", callCount)
	}

	// Second evaluation: should NOT use cache and call the dynamic fact again
	condition.Evaluate(almanac)
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
}

func TestConditionCachingWithDifferentConditions(t *testing.T) {
	callCount := 0
	dynamicFact := func(params map[string]interface{}) interface{} {
		callCount++
		return 25
	}

	almanac := gre.NewAlmanac(gre.WithAlmanacConditionCaching())
	almanac.AddFact("age", dynamicFact, gre.WithoutCache())

	cond1 := &gre.Condition{
		Fact:     "age",
		Operator: "greater_than",
		Value:    18,
	}

	cond2 := &gre.Condition{
		Fact:     "age",
		Operator: "less_than",
		Value:    30,
	}

	// Eval cond1
	cond1.Evaluate(almanac)
	if callCount != 1 {
		t.Errorf("First eval dynamic call failed, callCount=%d", callCount)
	}

	// Eval cond2: different condition, different cache key, should call fact again
	cond2.Evaluate(almanac)
	if callCount != 2 {
		t.Errorf("Second eval with different cond failed, callCount=%d", callCount)
	}

	// Eval cond1 again: should be cached
	cond1.Evaluate(almanac)
	if callCount != 2 {
		t.Errorf("Third eval (cached cond1) should not increment callCount, got %d", callCount)
	}
}

func TestEngineWithConditionCaching(t *testing.T) {
	callCount := 0
	dynamicFact := func(params map[string]interface{}) interface{} {
		callCount++
		return 25
	}

	// Create engine with condition caching enabled
	engine := gre.NewEngine(gre.WithConditionCaching())

	// Add rule with a condition
	rule1 := &gre.Rule{
		Name: "rule1",
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
	engine.AddRule(rule1)

	// Add second rule with SAME condition
	rule2 := &gre.Rule{
		Name: "rule2",
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
	almanac.AddFact("age", dynamicFact, gre.WithoutCache())

	// Run engine
	e, err := engine.Run(almanac)

	results := e.ReduceResults()

	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !results["rule1"] || !results["rule2"] {
		t.Errorf("Expected both rules to pass, got rule1=%v, rule2=%v", results["rule1"], results["rule2"])
	}

	// Because of condition caching, the dynamic fact should ONLY be called once
	// even if two rules share the same condition.
	if callCount != 1 {
		t.Errorf("Expected 1 dynamic fact call due to condition caching, got %d", callCount)
	}
}

func TestEngineOptionsEnableConditionCaching(t *testing.T) {
	engine := gre.NewEngine(gre.WithConditionCaching())
	// Test if it works without error
	almanac := gre.NewAlmanac()
	_, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func TestEvaluateCacheKeyErrors(t *testing.T) {
	almanac := gre.NewAlmanac(gre.WithAlmanacConditionCaching())

	// Force an error in GetCacheKey by using unmarshalable value for Fact caching
	// but Condition key also uses json.Marshal for Value.
	type unmarshalable struct {
		Func func()
	}

	badCond := &gre.Condition{
		Fact:     "f",
		Operator: "equal",
		Value:    unmarshalable{Func: func() {}},
	}

	_, err := badCond.Evaluate(almanac)
	if err == nil {
		t.Error("Expected error in Evaluate due to bad cache key")
	}
}

func TestConditionEvaluateGetCacheKeyError(t *testing.T) {
	almanac := gre.NewAlmanac(gre.WithAlmanacConditionCaching())
	badCond := &gre.Condition{Value: func() {}}
	_, err := badCond.Evaluate(almanac)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestConditionSetEvaluateGetCacheKeyError(t *testing.T) {
	almanac := gre.NewAlmanac(gre.WithAlmanacConditionCaching())
	badCond := &gre.Condition{Value: func() {}}
	cs := &gre.ConditionSet{All: []gre.ConditionNode{{Condition: badCond}}}
	_, err := cs.Evaluate(almanac)
	if err == nil {
		t.Error("Expected error")
	}
}

func TestConditionSetEvaluateCacheStorePaths(t *testing.T) {
	almanac := gre.NewAlmanac(gre.WithAlmanacConditionCaching())
	almanac.AddFact("f", 1)

	// Fail All path (stores false)
	csAll := &gre.ConditionSet{All: []gre.ConditionNode{{Condition: &gre.Condition{Fact: "f", Operator: "equal", Value: 2}}}}
	csAll.Evaluate(almanac)

	// Fail Any path (stores false)
	csAny := &gre.ConditionSet{Any: []gre.ConditionNode{{Condition: &gre.Condition{Fact: "f", Operator: "equal", Value: 2}}}}
	csAny.Evaluate(almanac)

	// Fail None path (stores false because one matched)
	csNone := &gre.ConditionSet{None: []gre.ConditionNode{{Condition: &gre.Condition{Fact: "f", Operator: "equal", Value: 1}}}}
	csNone.Evaluate(almanac)

	// Success path (stores true)
	csSucc := &gre.ConditionSet{All: []gre.ConditionNode{{Condition: &gre.Condition{Fact: "f", Operator: "equal", Value: 1}}}}
	csSucc.Evaluate(almanac)
}

func TestConditionSetCompileAndCacheKey(t *testing.T) {
	cs := &gre.ConditionSet{
		All: []gre.ConditionNode{
			{
				Condition: &gre.Condition{
					Fact:     "age",
					Operator: "greater_than",
					Value:    18,
				},
			},
		},
	}

	// Before compile, GetCacheKey should compute it
	key1, err := cs.GetCacheKey()
	if err != nil {
		t.Fatalf("GetCacheKey failed: %v", err)
	}

	// Compile the condition set
	err = cs.Compile()
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// After compile, GetCacheKey should return the cached key
	key2, err := cs.GetCacheKey()
	if err != nil {
		t.Fatalf("GetCacheKey failed after compile: %v", err)
	}

	if key1 != key2 {
		t.Errorf("Expected keys to be identical, got %s and %s", key1, key2)
	}
}

func TestConditionSetCompileErrors(t *testing.T) {
	badCond := &gre.Condition{Value: func() {}}

	// Error in All
	cs1 := &gre.ConditionSet{All: []gre.ConditionNode{{Condition: badCond}}}
	if err := cs1.Compile(); err == nil {
		t.Error("Expected error in All")
	}

	// Error in Any
	cs2 := &gre.ConditionSet{Any: []gre.ConditionNode{{Condition: badCond}}}
	if err := cs2.Compile(); err == nil {
		t.Error("Expected error in Any")
	}

	// Error in None
	cs3 := &gre.ConditionSet{None: []gre.ConditionNode{{Condition: badCond}}}
	if err := cs3.Compile(); err == nil {
		t.Error("Expected error in None")
	}

	// Error in nested SubSet
	cs4 := &gre.ConditionSet{All: []gre.ConditionNode{{SubSet: &gre.ConditionSet{All: []gre.ConditionNode{{Condition: badCond}}}}}}
	if err := cs4.Compile(); err == nil {
		t.Error("Expected error in nested SubSet")
	}

	// THE COVERAGE TARGET: Error in cs.GetCacheKey() while child compilation succeeds
	// By providing a node that has both Condition (valid) and SubSet (invalid)
	// Compile will process Condition and skip SubSet, but json.Marshal will fail on SubSet
	node := gre.ConditionNode{
		Condition: &gre.Condition{Fact: "f", Operator: "equal", Value: 1},
		SubSet:    &gre.ConditionSet{All: []gre.ConditionNode{{Condition: badCond}}},
	}
	cs5 := &gre.ConditionSet{All: []gre.ConditionNode{node}}
	if err := cs5.Compile(); err == nil {
		t.Error("Expected error in ConditionSet.Compile during GetCacheKey")
	}
}

func TestConditionCompileAndCacheKey(t *testing.T) {
	c := &gre.Condition{
		Fact:     "age",
		Operator: "greater_than",
		Value:    18,
	}

	// Before compile
	key1, err := c.GetCacheKey()
	if err != nil {
		t.Fatalf("GetCacheKey failed: %v", err)
	}

	// Compile
	err = c.Compile()
	if err != nil {
		t.Fatalf("Compile failed: %v", err)
	}

	// After compile
	key2, err := c.GetCacheKey()
	if err != nil {
		t.Fatalf("GetCacheKey failed after compile: %v", err)
	}

	if key1 != key2 {
		t.Errorf("Expected keys to be identical, got %s and %s", key1, key2)
	}
}

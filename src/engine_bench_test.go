package gorulesengine_test

import (
	"fmt"
	"testing"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

// Helper to create a simple rule
func createSimpleRule(name string, priority int) *gre.Rule {
	return &gre.Rule{
		Name:     name,
		Priority: priority,
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
}

// Benchmark simple execution with increasing number of rules
func BenchmarkEngine_Run_Simple(b *testing.B) {
	ruleCounts := []int{1, 10, 100}

	for _, count := range ruleCounts {
		b.Run(fmt.Sprintf("Rules-%d", count), func(b *testing.B) {
			engine := gre.NewEngine(gre.WithoutPrioritySorting())
			for i := 0; i < count; i++ {
				engine.AddRule(createSimpleRule(fmt.Sprintf("rule-%d", i), i))
			}

			almanac := gre.NewAlmanac()
			almanac.AddFact("age", 25)

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _ = engine.Run(almanac)
			}
		})
	}
}

// Benchmark Parallel vs Sequential execution
func BenchmarkEngine_ParallelVsSequential(b *testing.B) {
	numRules := 100
	engineSequential := gre.NewEngine(gre.WithoutParallelExecution())
	engineParallel := gre.NewEngine(gre.WithParallelExecution(4))

	for i := 0; i < numRules; i++ {
		rule := createSimpleRule(fmt.Sprintf("rule-%d", i), i)
		engineSequential.AddRule(rule)
		engineParallel.AddRule(rule)
	}

	// Use a slow fact to highlight parallel benefits
	almanac := gre.NewAlmanac()
	almanac.AddFact("age", func(params map[string]interface{}) interface{} {
		// Mock some "work" but not too much to keep bench running
		return 25
	})

	b.Run("Sequential", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = engineSequential.Run(almanac)
		}
	})

	b.Run("Parallel-4", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = engineParallel.Run(almanac)
		}
	})
}

// Benchmark Caching impact
func BenchmarkEngine_Caching(b *testing.B) {
	// Create 50 rules that all use the exact same condition
	sharedCond := &gre.Condition{
		Fact:     "status",
		Operator: "equal",
		Value:    "active",
	}

	engineNoCache := gre.NewEngine(gre.WithoutConditionCaching())
	engineWithCache := gre.NewEngine(gre.WithConditionCaching())

	for i := 0; i < 50; i++ {
		rule := &gre.Rule{
			Name:       fmt.Sprintf("rule-%d", i),
			Conditions: gre.ConditionSet{All: []gre.ConditionNode{{Condition: sharedCond}}},
		}
		engineNoCache.AddRule(rule)
		engineWithCache.AddRule(rule)
	}

	almanac := gre.NewAlmanac()
	almanac.AddFact("status", "active")

	b.Run("No-Cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = engineNoCache.Run(almanac)
		}
	})

	b.Run("With-Cache", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = engineWithCache.Run(almanac)
		}
	})
}

// Benchmark Smart Skip
func BenchmarkEngine_SmartSkip(b *testing.B) {
	engineNone := gre.NewEngine()
	engineSkip := gre.NewEngine(gre.WithSmartSkip())

	// Rule depends on "missing_fact"
	rule := &gre.Rule{
		Name: "rule-missing",
		Conditions: gre.ConditionSet{
			All: []gre.ConditionNode{
				{Condition: &gre.Condition{Fact: "missing_fact", Operator: "equal", Value: 1}},
			},
		},
	}
	engineNone.AddRule(rule)
	engineSkip.AddRule(rule)

	almanac := gre.NewAlmanac()

	b.Run("No-Smart-Skip", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = engineNone.Run(almanac)
		}
	})

	b.Run("With-Smart-Skip", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = engineSkip.Run(almanac)
		}
	})
}

// Benchmark complex nested rules
func BenchmarkEngine_ComplexRules(b *testing.B) {
	engine := gre.NewEngine()

	// 5 levels of nesting
	nested := gre.ConditionSet{
		All: []gre.ConditionNode{
			{Condition: &gre.Condition{Fact: "f1", Operator: "equal", Value: 1}},
			{
				SubSet: &gre.ConditionSet{
					Any: []gre.ConditionNode{
						{Condition: &gre.Condition{Fact: "f2", Operator: "equal", Value: 2}},
						{
							SubSet: &gre.ConditionSet{
								None: []gre.ConditionNode{
									{Condition: &gre.Condition{Fact: "f3", Operator: "equal", Value: 3}},
									{
										SubSet: &gre.ConditionSet{
											All: []gre.ConditionNode{
												{Condition: &gre.Condition{Fact: "f4", Operator: "equal", Value: 4}},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	engine.AddRule(&gre.Rule{
		Name:       "complex",
		Conditions: nested,
	})

	almanac := gre.NewAlmanac()
	almanac.AddFact("f1", 1)
	almanac.AddFact("f2", 2)
	almanac.AddFact("f3", 0)
	almanac.AddFact("f4", 4)

	b.Run("Run-Complex", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			_, _ = engine.Run(almanac)
		}
	})
}

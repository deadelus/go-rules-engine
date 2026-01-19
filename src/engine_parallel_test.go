package gorulesengine_test

import (
	"fmt"
	"testing"
	"time"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func f(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

func TestEngine_ParallelRun(t *testing.T) {
	t.Run("runs rules in parallel correctly", func(t *testing.T) {
		engine := gre.NewEngine(
			gre.WithParallelExecution(5),
			gre.WithAuditTrace(),
			gre.WithSmartSkip(),
		)

		numRules := 20
		for i := 0; i < numRules; i++ {
			engine.AddRule(&gre.Rule{
				Name: f("rule-%d", i),
				Conditions: gre.ConditionSet{
					All: []gre.ConditionNode{
						{Condition: &gre.Condition{Fact: "age", Operator: "greater_than", Value: i}},
					},
				},
				Priority: i,
			})
		}

		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 10)

		_, err := engine.Run(almanac)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		results := engine.Results()
		if len(results) != numRules {
			t.Fatalf("Expected %d results, got %d", numRules, len(results))
		}

		for i := 0; i < 10; i++ {
			res, ok := results[f("rule-%d", i)]
			if !ok {
				t.Errorf("Expected result for rule-%d", i)
				continue
			}
			if !res.Result {
				t.Errorf("Expected rule-%d to be true, got false", i)
			}
		}
		for i := 10; i < numRules; i++ {
			res, ok := results[f("rule-%d", i)]
			if !ok {
				t.Errorf("Expected result for rule-%d", i)
				continue
			}
			if res.Result {
				t.Errorf("Expected rule-%d to be false, got true", i)
			}
		}
	})

	t.Run("parallel execution handles smart skip", func(t *testing.T) {
		engine := gre.NewEngine(
			gre.WithParallelExecution(2),
			gre.WithSmartSkip(),
		)

		engine.AddRule(&gre.Rule{
			Name: "missing-fact-rule",
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{Condition: &gre.Condition{Fact: "missing", Operator: "equal", Value: 1}},
				},
			},
		})

		almanac := gre.NewAlmanac()

		_, err := engine.Run(almanac)
		if err != nil {
			t.Fatalf("Expected no error, got %v", err)
		}

		res := engine.Results()["missing-fact-rule"]
		if res.Result != false {
			t.Error("Expected rule to be false due to smart skip")
		}
	})

	t.Run("parallel execution with invalid worker count defaults to 1", func(t *testing.T) {
		engine := gre.NewEngine(
			gre.WithParallelExecution(0), // Should default to 1
		)

		engine.AddRule(&gre.Rule{
			Name: "rule1",
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{Condition: &gre.Condition{Fact: "a", Operator: "equal", Value: 1}},
				},
			},
		})

		almanac := gre.NewAlmanac()
		almanac.AddFact("a", 1)

		_, err := engine.Run(almanac)
		if err != nil {
			t.Fatal(err)
		}

		if !engine.Results()["rule1"].Result {
			t.Error("Expected rule1 to be true")
		}
	})

	t.Run("parallel execution handles errors", func(t *testing.T) {
		engine := gre.NewEngine(
			gre.WithParallelExecution(2),
		)

		engine.AddRule(&gre.Rule{
			Name: "error-rule",
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{Condition: &gre.Condition{Fact: "age", Operator: "unknown", Value: 1}},
				},
			},
		})

		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 20)

		_, err := engine.Run(almanac)
		if err == nil {
			t.Error("Expected error from Run due to unknown operator")
		}
	})

	t.Run("parallel execution triggers events sequentially", func(t *testing.T) {
		var order []string
		handler := func(ctx gre.EventContext) error {
			order = append(order, ctx.RuleName)
			return nil
		}

		engine := gre.NewEngine(
			gre.WithParallelExecution(4),
		)
		engine.RegisterEvent(gre.Event{Name: "track", Action: handler})

		for i := 0; i < 5; i++ {
			engine.AddRule(&gre.Rule{
				Name: f("r%d", i),
				Conditions: gre.ConditionSet{
					All: []gre.ConditionNode{
						{Condition: &gre.Condition{Fact: "x", Operator: "equal", Value: 1}},
					},
				},
				Priority:  i,
				OnSuccess: []gre.RuleEvent{{Name: "track"}},
			})
		}

		almanac := gre.NewAlmanac()
		almanac.AddFact("x", 1)

		_, _ = engine.Run(almanac)

		expected := []string{"r4", "r3", "r2", "r1", "r0"}
		if len(order) != 5 {
			t.Fatalf("Expected 5 events, got %d", len(order))
		}
		for i, v := range expected {
			if order[i] != v {
				t.Errorf("At index %d: expected %s, got %s", i, v, order[i])
			}
		}
	})

	t.Run("parallel execution handles event errors", func(t *testing.T) {
		engine := gre.NewEngine(
			gre.WithParallelExecution(2),
		)

		errHandler := func(ctx gre.EventContext) error {
			return fmt.Errorf("event error")
		}

		engine.RegisterEvent(gre.Event{
			Name:   "fail-event",
			Action: errHandler,
		})

		engine.AddRule(&gre.Rule{
			Name: "success-trigger",
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{Condition: &gre.Condition{Fact: "a", Operator: "equal", Value: 1}},
				},
			},
			OnSuccess: []gre.RuleEvent{{Name: "fail-event"}},
		})

		engine.AddRule(&gre.Rule{
			Name: "failure-trigger",
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{Condition: &gre.Condition{Fact: "a", Operator: "equal", Value: 2}},
				},
			},
			OnFailure: []gre.RuleEvent{{Name: "fail-event"}},
		})

		almanac := gre.NewAlmanac()
		almanac.AddFact("a", 1)

		_, err := engine.Run(almanac)
		if err == nil {
			t.Error("Expected error from OnSuccess event")
		}

		almanac.AddFact("a", 3) // Will make rule 2 fail its condition
		_, err = engine.Run(almanac)
		if err == nil {
			t.Error("Expected error from OnFailure event")
		}
	})
}

func TestEngine_ParallelSlowFacts(t *testing.T) {
	t.Run("parallel execution is faster for slow facts", func(t *testing.T) {
		engine := gre.NewEngine(
			gre.WithParallelExecution(5),
		)

		numRules := 5
		for i := 0; i < numRules; i++ {
			factName := gre.FactID(f("slow-%d", i))
			engine.AddRule(&gre.Rule{
				Name: f("rule-%d", i),
				Conditions: gre.ConditionSet{
					All: []gre.ConditionNode{
						{Condition: &gre.Condition{Fact: factName, Operator: "equal", Value: true}},
					},
				},
			})
		}

		almanac := gre.NewAlmanac()
		for i := 0; i < numRules; i++ {
			almanac.AddFact(gre.FactID(f("slow-%d", i)), func() interface{} {
				time.Sleep(100 * time.Millisecond)
				return true
			})
		}

		start := time.Now()
		_, err := engine.Run(almanac)
		duration := time.Since(start)

		if err != nil {
			t.Fatal(err)
		}

		if duration > 400*time.Millisecond {
			t.Errorf("Parallel execution took too long: %v, expected < 400ms", duration)
		}
	})
}

package gorulesengine_test

import (
	"sync"
	"testing"
	"time"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

// MockMetricsCollector is a mock implementation of MetricsCollector for testing.
type MockMetricsCollector struct {
	mu sync.Mutex

	RuleEvalCounts map[string]int
	EngineRunCount int
	EventExecCount int

	LastRuleResult bool
	LastRunCount   int
}

func NewMockMetricsCollector() *MockMetricsCollector {
	return &MockMetricsCollector{
		RuleEvalCounts: make(map[string]int),
	}
}

func (m *MockMetricsCollector) ObserveRuleEvaluation(ruleName string, result bool, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.RuleEvalCounts[ruleName]++
	m.LastRuleResult = result
}

func (m *MockMetricsCollector) ObserveEngineRun(ruleCount int, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.EngineRunCount++
	m.LastRunCount = ruleCount
}

func (m *MockMetricsCollector) ObserveEventExecution(eventName string, ruleName string, result bool, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.EventExecCount++
}

func TestEngine_Metrics(t *testing.T) {
	t.Run("sequential execution metrics", func(t *testing.T) {
		mock := NewMockMetricsCollector()
		engine := gre.NewEngine(
			gre.WithMetrics(mock),
			gre.WithoutPrioritySorting(),
		)

		rule := &gre.Rule{
			Name: "test-rule",
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{Condition: &gre.Condition{Fact: "age", Operator: "greater_than", Value: 18}},
				},
			},
			OnSuccess: []gre.RuleEvent{{Name: "success-event"}},
		}
		engine.AddRule(rule)
		engine.RegisterEvent(gre.Event{Name: "success-event"})

		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 25)

		_, err := engine.Run(almanac)
		if err != nil {
			t.Fatalf("Run failed: %v", err)
		}

		if mock.EngineRunCount != 1 {
			t.Errorf("Expected 1 engine run, got %d", mock.EngineRunCount)
		}
		if mock.RuleEvalCounts["test-rule"] != 1 {
			t.Errorf("Expected 1 rule evaluation, got %d", mock.RuleEvalCounts["test-rule"])
		}
		if mock.EventExecCount != 1 {
			t.Errorf("Expected 1 event execution, got %d", mock.EventExecCount)
		}
	})

	t.Run("parallel execution metrics", func(t *testing.T) {
		mock := NewMockMetricsCollector()
		engine := gre.NewEngine(
			gre.WithMetrics(mock),
			gre.WithParallelExecution(2),
		)

		rule := &gre.Rule{
			Name: "parallel-rule",
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{Condition: &gre.Condition{Fact: "age", Operator: "greater_than", Value: 18}},
				},
			},
		}
		engine.AddRule(rule)

		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 20)

		_, err := engine.Run(almanac)
		if err != nil {
			t.Fatalf("Run failed: %v", err)
		}

		if mock.EngineRunCount != 1 {
			t.Errorf("Expected 1 engine run, got %d", mock.EngineRunCount)
		}
		if mock.RuleEvalCounts["parallel-rule"] != 1 {
			t.Errorf("Expected 1 rule evaluation, got %d", mock.RuleEvalCounts["parallel-rule"])
		}
	})

	t.Run("smart skip metrics - sequential", func(t *testing.T) {
		mock := NewMockMetricsCollector()
		engine := gre.NewEngine(
			gre.WithMetrics(mock),
			gre.WithSmartSkip(),
		)

		rule := &gre.Rule{
			Name: "skipped-rule",
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{Condition: &gre.Condition{Fact: "missing", Operator: "equal", Value: 1}},
				},
			},
		}
		engine.AddRule(rule)

		almanac := gre.NewAlmanac()
		_, _ = engine.Run(almanac)

		if mock.RuleEvalCounts["skipped-rule"] != 1 {
			t.Errorf("Expected rule to be counted even if skipped, got %d", mock.RuleEvalCounts["skipped-rule"])
		}
	})

	t.Run("smart skip metrics - parallel", func(t *testing.T) {
		mock := NewMockMetricsCollector()
		engine := gre.NewEngine(
			gre.WithMetrics(mock),
			gre.WithParallelExecution(2),
			gre.WithSmartSkip(),
		)

		rule := &gre.Rule{
			Name: "skipped-parallel",
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{Condition: &gre.Condition{Fact: "missing", Operator: "equal", Value: 1}},
				},
			},
		}
		engine.AddRule(rule)

		almanac := gre.NewAlmanac()
		_, _ = engine.Run(almanac)

		if mock.RuleEvalCounts["skipped-parallel"] != 1 {
			t.Errorf("Expected parallel skipped rule to be counted, got %d", mock.RuleEvalCounts["skipped-parallel"])
		}
	})

	t.Run("async event metrics", func(t *testing.T) {
		mock := NewMockMetricsCollector()
		engine := gre.NewEngine(
			gre.WithMetrics(mock),
		)

		engine.RegisterEvent(gre.Event{
			Name: "async-event",
			Mode: gre.EventModeAsync,
		})

		almanac := gre.NewAlmanac()
		_ = engine.HandleEvent("async-event", "rule", true, almanac, nil)

		// Wait for async execution
		time.Sleep(50 * time.Millisecond)

		mock.mu.Lock()
		count := mock.EventExecCount
		mock.mu.Unlock()

		if count != 1 {
			t.Errorf("Expected 1 async event execution, got %d", count)
		}
	})

	t.Run("WithMetrics handles nil engine", func(t *testing.T) {
		opt := gre.WithMetrics(nil)
		// Should not panic
		opt(nil)
	})
}

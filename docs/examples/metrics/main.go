package main

import (
	"fmt"
	"time"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

// FakePrometheusCollector simulates a Prometheus implementation of the MetricsCollector interface.
// In a real application, you would use "github.com/prometheus/client_golang/prometheus".
type FakePrometheusCollector struct {
	// Counters and Histograms would be here in a real Prometheus setup
	ruleEvals map[string]int
}

func NewFakePrometheusCollector() *FakePrometheusCollector {
	return &FakePrometheusCollector{
		ruleEvals: make(map[string]int),
	}
}

// ObserveRuleEvaluation simulates recording a metric for a single rule.
func (c *FakePrometheusCollector) ObserveRuleEvaluation(ruleName string, result bool, duration time.Duration) {
	c.ruleEvals[ruleName]++
	status := "FAIL"
	if result {
		status = "PASS"
	}
	fmt.Printf("[METRIC] Rule '%s' evaluated in %v. Status: %s. Total evals: %d\n",
		ruleName, duration, status, c.ruleEvals[ruleName])
}

// ObserveEngineRun simulates recording a metric for the entire engine execution.
func (c *FakePrometheusCollector) ObserveEngineRun(ruleCount int, duration time.Duration) {
	fmt.Printf("[METRIC] Engine execution completed in %v. Total rules evaluated: %d\n",
		duration, ruleCount)
}

// ObserveEventExecution simulates recording a metric for event execution.
func (c *FakePrometheusCollector) ObserveEventExecution(eventName string, ruleName string, result bool, duration time.Duration) {
	fmt.Printf("[METRIC] Event '%s' triggered by rule '%s' executed in %v\n",
		eventName, ruleName, duration)
}

func main() {
	// 1. Initialize the fake Prometheus collector
	collector := NewFakePrometheusCollector()

	// 2. Create the engine with the Metrics option
	engine := gre.NewEngine(
		gre.WithMetrics(collector),
		gre.WithConditionCaching(),
	)

	// 3. Define a rule
	rule := &gre.Rule{
		Name: "premium-user-check",
		Conditions: gre.ConditionSet{
			All: []gre.ConditionNode{
				{
					Condition: &gre.Condition{
						Fact:     "account_type",
						Operator: "equal",
						Value:    "premium",
					},
				},
				{
					Condition: &gre.Condition{
						Fact:     "usage_count",
						Operator: "greater_than",
						Value:    100,
					},
				},
			},
		},
		OnSuccess: []gre.RuleEvent{{Name: "notify-support"}},
	}

	engine.AddRule(rule)

	// Register an event
	engine.RegisterEvent(gre.Event{
		Name: "notify-support",
		Action: func(ctx gre.EventContext) error {
			fmt.Println(">> Action: Notifying support for premium user activity")
			return nil
		},
	})

	// 4. Create an Almanac with facts
	almanac := gre.NewAlmanac()
	almanac.AddFact("account_type", "premium")
	almanac.AddFact("usage_count", 150)

	// 5. Run the engine (metrics are automatically collected)
	fmt.Println("--- Running Engine ---")
	_, err := engine.Run(almanac)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Println("\n--- Running Engine again (testing cache) ---")
	_, _ = engine.Run(almanac)

	// Wait a bit to see any async metrics if we had any
	time.Sleep(100 * time.Millisecond)
}

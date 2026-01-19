package gorulesengine_test

import (
	"errors"
	"testing"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

// MockEventHandler is a mock implementation of EventHandler for testing
type MockEventHandler struct {
	HandledEvents   []gre.Event
	HandledContexts []gre.EventContext
	ShouldError     bool
	ErrorMessage    string
}

func (m *MockEventHandler) Handle(event gre.Event, ctx gre.EventContext) error {
	m.HandledEvents = append(m.HandledEvents, event)
	m.HandledContexts = append(m.HandledContexts, ctx)
	if m.ShouldError {
		return errors.New(m.ErrorMessage)
	}
	return nil
}

func TestNewEngine(t *testing.T) {
	t.Run("creates engine with default priority sorting", func(t *testing.T) {
		engine := gre.NewEngine()
		if engine == nil {
			t.Fatal("Expected engine to be created")
		}
	})

	t.Run("creates engine with custom sorting options", func(t *testing.T) {
		sortRule := gre.SortRuleASC
		engine := gre.NewEngine(gre.WithPrioritySorting(&sortRule))
		if engine == nil {
			t.Fatal("Expected engine to be created")
		}
	})

	t.Run("creates engine without priority sorting", func(t *testing.T) {
		engine := gre.NewEngine(gre.WithoutPrioritySorting())
		if engine == nil {
			t.Fatal("Expected engine to be created")
		}
	})
}

func TestAddRule(t *testing.T) {
	t.Run("adds a single rule", func(t *testing.T) {
		engine := gre.NewEngine()
		rule := &gre.Rule{
			Name:     "test-rule",
			Priority: 10,
		}

		engine.AddRule(rule)
		// Since rules is private, we can only verify no panic occurs
	})

	t.Run("adds multiple rules", func(t *testing.T) {
		engine := gre.NewEngine()
		rule1 := &gre.Rule{Name: "rule1", Priority: 10}
		rule2 := &gre.Rule{Name: "rule2", Priority: 20}

		engine.AddRule(rule1)
		engine.AddRule(rule2)
		// Since rules is private, we can only verify no panic occurs
	})
}

func TestAddRules(t *testing.T) {
	e := gre.NewEngine()

	r1 := &gre.Rule{
		Name: "Rule 1",
		Conditions: gre.ConditionSet{
			All: []gre.ConditionNode{
				{Condition: &gre.Condition{Fact: "f1", Operator: "equal", Value: 1}},
			},
		},
	}

	r2 := &gre.Rule{
		Name: "Rule 2",
		Conditions: gre.ConditionSet{
			All: []gre.ConditionNode{
				{Condition: &gre.Condition{Fact: "f2", Operator: "equal", Value: 2}},
			},
		},
	}

	e.AddRules(r1, r2)

	rules := e.GetRules()
	if len(rules) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(rules))
	}

	if rules[0].Name != "Rule 1" || rules[1].Name != "Rule 2" {
		t.Error("Rules names are incorrect")
	}
}

func TestClearRules(t *testing.T) {
	e := gre.NewEngine()
	e.AddRule(&gre.Rule{Name: "Rule 1"})
	e.AddRule(&gre.Rule{Name: "Rule 2"})

	if len(e.GetRules()) != 2 {
		t.Errorf("Expected 2 rules, got %d", len(e.GetRules()))
	}

	e.ClearRules()

	if len(e.GetRules()) != 0 {
		t.Errorf("Expected 0 rules after ClearRules, got %d", len(e.GetRules()))
	}
}

func TestRegisterEvent(t *testing.T) {
	t.Run("registers a new event", func(t *testing.T) {
		engine := gre.NewEngine()
		event := gre.Event{
			Name: "test-event",
			Params: map[string]interface{}{
				"key": "value",
			},
		}

		engine.RegisterEvent(event)

		// Verify by trying to handle the event (no error should occur)
		almanac := gre.NewAlmanac()
		err := engine.HandleEvent("test-event", "test-rule", true, almanac, nil)
		if err != nil {
			t.Fatalf("Expected event to be registered, got error: %v", err)
		}
	})

	t.Run("registers multiple events", func(t *testing.T) {
		engine := gre.NewEngine()
		event1 := gre.Event{Name: "event1"}
		event2 := gre.Event{Name: "event2"}

		engine.RegisterEvent(event1)
		engine.RegisterEvent(event2)

		// Both events should be accessible
		almanac := gre.NewAlmanac()
		if err := engine.HandleEvent("event1", "rule1", true, almanac, nil); err != nil {
			t.Fatalf("Expected event1 to be registered, got error: %v", err)
		}
		if err := engine.HandleEvent("event2", "rule2", true, almanac, nil); err != nil {
			t.Fatalf("Expected event2 to be registered, got error: %v", err)
		}
	})
}

func TestRegisterEvents(t *testing.T) {
	t.Run("registers multiple events at once", func(t *testing.T) {
		engine := gre.NewEngine()
		event1 := gre.Event{Name: "event1"}
		event2 := gre.Event{Name: "event2"}
		event3 := gre.Event{Name: "event3"}

		engine.RegisterEvents(event1, event2, event3)

		// All events should be accessible
		almanac := gre.NewAlmanac()
		if err := engine.HandleEvent("event1", "rule1", true, almanac, nil); err != nil {
			t.Fatalf("Expected event1 to be registered, got error: %v", err)
		}
		if err := engine.HandleEvent("event2", "rule2", true, almanac, nil); err != nil {
			t.Fatalf("Expected event2 to be registered, got error: %v", err)
		}
		if err := engine.HandleEvent("event3", "rule3", true, almanac, nil); err != nil {
			t.Fatalf("Expected event3 to be registered, got error: %v", err)
		}
	})

	t.Run("works with no events", func(t *testing.T) {
		engine := gre.NewEngine()
		engine.RegisterEvents()

		// HandleEvent should return error only if handler is set and event doesn't exist
		mockHandler := &MockEventHandler{}
		engine.SetEventHandler(mockHandler)

		almanac := gre.NewAlmanac()
		err := engine.HandleEvent("nonexistent", "rule", true, almanac, nil)
		if err == nil {
			t.Fatal("Expected error for nonexistent event when handler is set")
		}
	})
}

func TestEngineRun(t *testing.T) {
	t.Run("runs engine with passing rule", func(t *testing.T) {
		engine := gre.NewEngine()
		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 25)

		rule := &gre.Rule{
			Name:     "adult-check",
			Priority: 10,
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

		result := e.ReduceResults()

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !result[rule.Name] {
			t.Error("Expected rule to pass")
		}
	})

	t.Run("runs engine with failing rule", func(t *testing.T) {
		engine := gre.NewEngine()
		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 15)

		rule := &gre.Rule{
			Name:     "adult-check",
			Priority: 10,
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

		result := e.ReduceResults()

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if result[rule.Name] {
			t.Error("Expected rule to fail")
		}
	})

	t.Run("runs engine with onSuccess event", func(t *testing.T) {
		engine := gre.NewEngine()
		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 25)

		mockHandler := &MockEventHandler{}
		engine.SetEventHandler(mockHandler)

		event := gre.Event{
			Name: "success-event",
			Params: map[string]interface{}{
				"message": "Rule passed",
			},
		}
		engine.RegisterEvent(event)

		rule := &gre.Rule{
			Name:     "adult-check",
			Priority: 10,
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
			OnSuccess: []gre.RuleEvent{{Name: "success-event"}},
		}

		engine.AddRule(rule)
		e, err := engine.Run(almanac)

		result := e.ReduceResults()

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !result[rule.Name] {
			t.Error("Expected rule to pass")
		}
		if len(mockHandler.HandledEvents) != 1 {
			t.Errorf("Expected 1 event to be handled, got %d", len(mockHandler.HandledEvents))
		}
		if mockHandler.HandledEvents[0].Name != "success-event" {
			t.Errorf("Expected event name 'success-event', got '%s'", mockHandler.HandledEvents[0].Name)
		}
	})

	t.Run("runs engine with onFailure event", func(t *testing.T) {
		engine := gre.NewEngine()
		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 15)

		mockHandler := &MockEventHandler{}
		engine.SetEventHandler(mockHandler)

		event := gre.Event{
			Name: "failure-event",
			Params: map[string]interface{}{
				"message": "Rule failed",
			},
		}
		engine.RegisterEvent(event)

		rule := &gre.Rule{
			Name:     "adult-check",
			Priority: 10,
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
			OnFailure: []gre.RuleEvent{{Name: "failure-event"}},
		}

		engine.AddRule(rule)

		e, err := engine.Run(almanac)

		result := e.ReduceResults()

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if result[rule.Name] {
			t.Error("Expected rule to fail")
		}
		if len(mockHandler.HandledEvents) != 1 {
			t.Errorf("Expected 1 event to be handled, got %d", len(mockHandler.HandledEvents))
		}
		if mockHandler.HandledEvents[0].Name != "failure-event" {
			t.Errorf("Expected event name 'failure-event', got '%s'", mockHandler.HandledEvents[0].Name)
		}
	})

	t.Run("runs engine with multiple rules in priority order DESC", func(t *testing.T) {
		engine := gre.NewEngine() // Default is DESC
		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 25)

		rule1 := &gre.Rule{
			Name:     "low-priority",
			Priority: 5,
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

		rule2 := &gre.Rule{
			Name:     "high-priority",
			Priority: 10,
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

		engine.AddRule(rule1)
		engine.AddRule(rule2)

		e, err := engine.Run(almanac)

		result := e.ReduceResults()

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !result[rule1.Name] || !result[rule2.Name] {
			t.Error("Expected rules to pass")
		}
	})

	t.Run("runs engine with multiple rules in priority order ASC", func(t *testing.T) {
		sortRule := gre.SortRuleASC
		engine := gre.NewEngine(gre.WithPrioritySorting(&sortRule))
		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 25)

		rule1 := &gre.Rule{
			Name:     "low-priority",
			Priority: 5,
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

		rule2 := &gre.Rule{
			Name:     "high-priority",
			Priority: 10,
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

		engine.AddRule(rule1)
		engine.AddRule(rule2)

		e, err := engine.Run(almanac)

		result := e.ReduceResults()

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !result[rule1.Name] || !result[rule2.Name] {
			t.Error("Expected rules to pass")
		}
	})

	t.Run("returns error when event handler fails", func(t *testing.T) {
		engine := gre.NewEngine()
		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 25)

		mockHandler := &MockEventHandler{
			ShouldError:  true,
			ErrorMessage: "handler error",
		}
		engine.SetEventHandler(mockHandler)

		event := gre.Event{Name: "test-event"}
		engine.RegisterEvent(event)

		rule := &gre.Rule{
			Name:     "test-rule",
			Priority: 10,
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
			OnSuccess: []gre.RuleEvent{{Name: "test-event"}},
		}

		engine.AddRule(rule)
		_, err := engine.Run(almanac)

		if err == nil {
			t.Fatal("Expected error when handler fails")
		}
	})

	t.Run("returns error for unregistered event", func(t *testing.T) {
		engine := gre.NewEngine()
		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 25)

		mockHandler := &MockEventHandler{}
		engine.SetEventHandler(mockHandler)

		rule := &gre.Rule{
			Name:     "test-rule",
			Priority: 10,
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
			OnSuccess: []gre.RuleEvent{{Name: "unregistered-event"}},
		}

		engine.AddRule(rule)
		_, err := engine.Run(almanac)

		if err == nil {
			t.Fatal("Expected error for unregistered event")
		}
	})

	t.Run("returns error when condition evaluation fails", func(t *testing.T) {
		engine := gre.NewEngine()
		almanac := gre.NewAlmanac()
		// Note: Not adding the required fact to trigger an error

		rule := &gre.Rule{
			Name:     "test-rule",
			Priority: 10,
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{
						Condition: &gre.Condition{
							Fact:     "nonexistent-fact",
							Operator: "equal",
							Value:    "test",
						},
					},
				},
			},
		}

		engine.AddRule(rule)
		_, err := engine.Run(almanac)

		if err == nil {
			t.Fatal("Expected error when fact doesn't exist")
		}
	})

	t.Run("returns error when onFailure event handler fails", func(t *testing.T) {
		engine := gre.NewEngine()
		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 15)

		mockHandler := &MockEventHandler{
			ShouldError:  true,
			ErrorMessage: "failure handler error",
		}
		engine.SetEventHandler(mockHandler)

		event := gre.Event{Name: "failure-event"}
		engine.RegisterEvent(event)

		rule := &gre.Rule{
			Name:     "test-rule",
			Priority: 10,
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
			OnFailure: []gre.RuleEvent{{Name: "failure-event"}},
		}

		engine.AddRule(rule)
		_, err := engine.Run(almanac)

		if err == nil {
			t.Fatal("Expected error when onFailure handler fails")
		}
	})

	t.Run("returns error for unregistered onFailure event", func(t *testing.T) {
		engine := gre.NewEngine()
		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 15)

		mockHandler := &MockEventHandler{}
		engine.SetEventHandler(mockHandler)

		rule := &gre.Rule{
			Name:     "test-rule",
			Priority: 10,
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
			OnFailure: []gre.RuleEvent{{Name: "unregistered-failure-event"}},
		}

		engine.AddRule(rule)
		_, err := engine.Run(almanac)

		if err == nil {
			t.Fatal("Expected error for unregistered onFailure event")
		}
	})
}

func TestHandleEvent(t *testing.T) {
	t.Run("handles registered event successfully", func(t *testing.T) {
		engine := gre.NewEngine()
		mockHandler := &MockEventHandler{}
		engine.SetEventHandler(mockHandler)

		event := gre.Event{
			Name: "test-event",
			Params: map[string]interface{}{
				"key": "value",
			},
		}
		engine.RegisterEvent(event)
		almanac := gre.NewAlmanac()

		err := engine.HandleEvent("test-event", "test-rule", true, almanac, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if len(mockHandler.HandledEvents) != 1 {
			t.Errorf("Expected 1 event to be handled, got %d", len(mockHandler.HandledEvents))
		}
		if len(mockHandler.HandledContexts) != 1 {
			t.Errorf("Expected 1 context to be handled, got %d", len(mockHandler.HandledContexts))
		}
		// Verify context data
		ctx := mockHandler.HandledContexts[0]
		if ctx.RuleName != "test-rule" {
			t.Errorf("Expected rule name 'test-rule', got %s", ctx.RuleName)
		}
		if ctx.Result != true {
			t.Errorf("Expected result true, got %v", ctx.Result)
		}
		if ctx.Almanac != almanac {
			t.Error("Expected almanac to match")
		}
	})

	t.Run("returns error for unregistered event", func(t *testing.T) {
		engine := gre.NewEngine()
		mockHandler := &MockEventHandler{}
		engine.SetEventHandler(mockHandler)
		almanac := gre.NewAlmanac()

		err := engine.HandleEvent("unregistered-event", "test-rule", true, almanac, nil)

		if err == nil {
			t.Fatal("Expected error for unregistered event")
		}
	})

	t.Run("does not error when no event handler is set", func(t *testing.T) {
		engine := gre.NewEngine()
		event := gre.Event{Name: "any-event"}
		engine.RegisterEvent(event)
		almanac := gre.NewAlmanac()

		err := engine.HandleEvent("any-event", "test-rule", true, almanac, nil)

		if err != nil {
			t.Fatalf("Expected no error when no handler is set, got: %v", err)
		}
	})

	t.Run("returns nil when event not registered and no handler set", func(t *testing.T) {
		engine := gre.NewEngine()
		// Don't register event and don't set handler
		almanac := gre.NewAlmanac()

		err := engine.HandleEvent("nonexistent-event", "test-rule", true, almanac, nil)

		if err != nil {
			t.Fatalf("Expected no error when event not registered and no handler, got: %v", err)
		}
	})

	t.Run("returns error when event handler fails", func(t *testing.T) {
		engine := gre.NewEngine()
		mockHandler := &MockEventHandler{
			ShouldError:  true,
			ErrorMessage: "handler failed",
		}
		engine.SetEventHandler(mockHandler)

		event := gre.Event{Name: "test-event"}
		engine.RegisterEvent(event)
		almanac := gre.NewAlmanac()

		err := engine.HandleEvent("test-event", "test-rule", true, almanac, nil)

		if err == nil {
			t.Fatal("Expected error when handler fails")
		}
	})
}

func TestWithPrioritySorting(t *testing.T) {
	t.Run("sets default DESC sorting when nil", func(t *testing.T) {
		engine := gre.NewEngine(gre.WithPrioritySorting(nil))
		if engine == nil {
			t.Fatal("Expected engine to be created")
		}
	})

	t.Run("sets ASC sorting", func(t *testing.T) {
		sortRule := gre.SortRuleASC
		engine := gre.NewEngine(gre.WithPrioritySorting(&sortRule))
		if engine == nil {
			t.Fatal("Expected engine to be created")
		}
	})

	t.Run("sets DESC sorting", func(t *testing.T) {
		sortRule := gre.SortRuleDESC
		engine := gre.NewEngine(gre.WithPrioritySorting(&sortRule))
		if engine == nil {
			t.Fatal("Expected engine to be created")
		}
	})

	t.Run("handles invalid sort rule", func(t *testing.T) {
		sortRule := gre.SortRule(999)
		engine := gre.NewEngine(gre.WithPrioritySorting(&sortRule))
		if engine == nil {
			t.Fatal("Expected engine to be created")
		}
	})
}

func TestWithoutPrioritySorting(t *testing.T) {
	t.Run("disables priority sorting", func(t *testing.T) {
		engine := gre.NewEngine(gre.WithoutPrioritySorting())
		if engine == nil {
			t.Fatal("Expected engine to be created")
		}
	})

	t.Run("can override default sorting", func(t *testing.T) {
		engine := gre.NewEngine(
			gre.WithoutPrioritySorting(),
		)
		if engine == nil {
			t.Fatal("Expected engine to be created")
		}
	})
}

func TestOptionsFullCoverage(t *testing.T) {
	// Call options with nil to cover the guard clauses
	gre.WithConditionCaching()(nil)
	gre.WithoutConditionCaching()(nil)
	gre.WithSmartSkip()(nil)
	gre.WithAuditTrace()(nil)
	gre.WithoutAuditTrace()(nil)
	gre.WithParallelExecution(2)(nil)
	gre.WithoutParallelExecution()(nil)

	// Cover the 'if e.options == nil' false branch by calling options on an engine that already has options
	e := gre.NewEngine(gre.WithSmartSkip(), gre.WithConditionCaching(), gre.WithoutConditionCaching())
	gre.WithAuditTrace()(e)
	gre.WithoutAuditTrace()(e)
	gre.WithParallelExecution(2)(e)
	gre.WithoutParallelExecution()(e)
	_ = e

	// Cover the 'if e.options == nil' true branch by calling options on a raw Engine pointer
	gre.WithConditionCaching()(&gre.Engine{})
	gre.WithoutConditionCaching()(&gre.Engine{})
	gre.WithSmartSkip()(&gre.Engine{})
	gre.WithAuditTrace()(&gre.Engine{})
	gre.WithoutAuditTrace()(&gre.Engine{})
	gre.WithParallelExecution(2)(&gre.Engine{})
	gre.WithoutParallelExecution()(&gre.Engine{})
}

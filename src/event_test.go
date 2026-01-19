package gorulesengine_test

import (
	"errors"
	"sync"
	"testing"
	"time"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func TestEventModeSync(t *testing.T) {
	t.Run("executes event action synchronously", func(t *testing.T) {
		actionExecuted := false
		event := gre.Event{
			Name: "test-event",
			Mode: gre.EventModeSync,
			Action: func(ctx gre.EventContext) error {
				actionExecuted = true
				return nil
			},
		}

		engine := gre.NewEngine()
		engine.RegisterEvent(event)

		almanac := gre.NewAlmanac()
		err := engine.HandleEvent("test-event", "test-rule", true, almanac, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !actionExecuted {
			t.Error("Expected action to be executed")
		}
	})

	t.Run("returns error when action fails", func(t *testing.T) {
		expectedError := errors.New("action failed")
		event := gre.Event{
			Name: "test-event",
			Mode: gre.EventModeSync,
			Action: func(ctx gre.EventContext) error {
				return expectedError
			},
		}

		engine := gre.NewEngine()
		engine.RegisterEvent(event)

		almanac := gre.NewAlmanac()
		err := engine.HandleEvent("test-event", "test-rule", true, almanac, nil)

		if err == nil {
			t.Fatal("Expected error from action")
		}
	})

	t.Run("calls handler after action execution", func(t *testing.T) {
		actionExecuted := false
		handlerCalled := false

		event := gre.Event{
			Name: "test-event",
			Mode: gre.EventModeSync,
			Action: func(ctx gre.EventContext) error {
				actionExecuted = true
				return nil
			},
		}

		mockHandler := &MockEventHandler{}
		engine := gre.NewEngine()
		engine.SetEventHandler(mockHandler)
		engine.RegisterEvent(event)

		almanac := gre.NewAlmanac()
		err := engine.HandleEvent("test-event", "test-rule", true, almanac, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !actionExecuted {
			t.Error("Expected action to be executed")
		}
		if len(mockHandler.HandledEvents) == 0 {
			t.Error("Expected handler to be called")
		} else {
			handlerCalled = true
		}
		if !handlerCalled {
			t.Error("Expected handler to be called")
		}
	})

	t.Run("sync event with only action (no handler)", func(t *testing.T) {
		actionExecuted := false

		event := gre.Event{
			Name: "test-event",
			Mode: gre.EventModeSync,
			Action: func(ctx gre.EventContext) error {
				actionExecuted = true
				return nil
			},
		}

		engine := gre.NewEngine()
		// No handler set
		engine.RegisterEvent(event)

		almanac := gre.NewAlmanac()
		err := engine.HandleEvent("test-event", "test-rule", true, almanac, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !actionExecuted {
			t.Error("Expected action to be executed")
		}
	})

	t.Run("sync event with only handler (no action)", func(t *testing.T) {
		event := gre.Event{
			Name: "test-event",
			Mode: gre.EventModeSync,
			// No action
		}

		mockHandler := &MockEventHandler{}
		engine := gre.NewEngine()
		engine.SetEventHandler(mockHandler)
		engine.RegisterEvent(event)

		almanac := gre.NewAlmanac()
		err := engine.HandleEvent("test-event", "test-rule", true, almanac, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(mockHandler.HandledEvents) != 1 {
			t.Error("Expected handler to be called")
		}
	})
}

func TestEventModeAsync(t *testing.T) {
	t.Run("executes event action asynchronously", func(t *testing.T) {
		var mu sync.Mutex
		actionExecuted := false

		event := gre.Event{
			Name: "test-event",
			Mode: gre.EventModeAsync,
			Action: func(ctx gre.EventContext) error {
				time.Sleep(50 * time.Millisecond)
				mu.Lock()
				actionExecuted = true
				mu.Unlock()
				return nil
			},
		}

		engine := gre.NewEngine()
		engine.RegisterEvent(event)

		almanac := gre.NewAlmanac()
		err := engine.HandleEvent("test-event", "test-rule", true, almanac, nil)

		// Async execution should return immediately without error
		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Action should not be executed yet (asynchronous)
		mu.Lock()
		immediateExecution := actionExecuted
		mu.Unlock()

		if immediateExecution {
			t.Error("Expected action to execute asynchronously, but it executed immediately")
		}

		// Wait for async execution with a timeout
		success := false
		for i := 0; i < 20; i++ {
			mu.Lock()
			if actionExecuted {
				success = true
				mu.Unlock()
				break
			}
			mu.Unlock()
			time.Sleep(50 * time.Millisecond)
		}

		if !success {
			t.Error("Expected action to be executed after waiting")
		}
	})

	t.Run("does not block on action errors", func(t *testing.T) {
		event := gre.Event{
			Name: "test-event",
			Mode: gre.EventModeAsync,
			Action: func(ctx gre.EventContext) error {
				time.Sleep(10 * time.Millisecond)
				return errors.New("async action failed")
			},
		}

		engine := gre.NewEngine()
		engine.RegisterEvent(event)

		almanac := gre.NewAlmanac()
		err := engine.HandleEvent("test-event", "test-rule", true, almanac, nil)

		// Should not return error for async failures
		if err != nil {
			t.Fatalf("Expected no error for async action, got: %v", err)
		}
	})

	t.Run("calls handler in async mode", func(t *testing.T) {
		var mu sync.Mutex
		handlerCalled := false

		event := gre.Event{
			Name: "test-event",
			Mode: gre.EventModeAsync,
			Action: func(ctx gre.EventContext) error {
				time.Sleep(10 * time.Millisecond)
				return nil
			},
		}

		mockHandler := &MockEventHandler{}
		engine := gre.NewEngine()
		engine.SetEventHandler(mockHandler)
		engine.RegisterEvent(event)

		almanac := gre.NewAlmanac()
		err := engine.HandleEvent("test-event", "test-rule", true, almanac, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Handler should not be called immediately for async events (it's in the goroutine)
		if len(mockHandler.HandledEvents) != 0 {
			t.Error("Expected handler to not be called immediately for async event")
		}

		// Wait for async execution
		time.Sleep(50 * time.Millisecond)

		// Now the handler should have been called
		mu.Lock()
		defer mu.Unlock()
		if len(mockHandler.HandledEvents) != 1 {
			t.Errorf("Expected handler to be called after async execution, got %d calls", len(mockHandler.HandledEvents))
		} else {
			handlerCalled = true
		}

		if !handlerCalled {
			t.Error("Expected handler to be called")
		}
	})

	t.Run("async event with only action (no handler)", func(t *testing.T) {
		var mu sync.Mutex
		actionExecuted := false

		event := gre.Event{
			Name: "test-event",
			Mode: gre.EventModeAsync,
			Action: func(ctx gre.EventContext) error {
				time.Sleep(10 * time.Millisecond)
				mu.Lock()
				actionExecuted = true
				mu.Unlock()
				return nil
			},
		}

		engine := gre.NewEngine()
		// No handler set
		engine.RegisterEvent(event)

		almanac := gre.NewAlmanac()
		err := engine.HandleEvent("test-event", "test-rule", true, almanac, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Wait for async execution
		time.Sleep(50 * time.Millisecond)

		mu.Lock()
		defer mu.Unlock()
		if !actionExecuted {
			t.Error("Expected action to be executed")
		}
	})

	t.Run("async event with only handler (no action)", func(t *testing.T) {
		event := gre.Event{
			Name: "test-event",
			Mode: gre.EventModeAsync,
			// No action
		}

		mockHandler := &MockEventHandler{}
		engine := gre.NewEngine()
		engine.SetEventHandler(mockHandler)
		engine.RegisterEvent(event)

		almanac := gre.NewAlmanac()
		err := engine.HandleEvent("test-event", "test-rule", true, almanac, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Wait for async execution
		time.Sleep(50 * time.Millisecond)

		if len(mockHandler.HandledEvents) != 1 {
			t.Errorf("Expected handler to be called, got %d calls", len(mockHandler.HandledEvents))
		}
	})
}

func TestEventContext(t *testing.T) {
	t.Run("passes correct context to action", func(t *testing.T) {
		var capturedContext gre.EventContext

		event := gre.Event{
			Name: "test-event",
			Mode: gre.EventModeSync,
			Params: map[string]interface{}{
				"key1": "value1",
				"key2": 42,
			},
			Action: func(ctx gre.EventContext) error {
				capturedContext = ctx
				return nil
			},
		}

		engine := gre.NewEngine()
		engine.RegisterEvent(event)

		almanac := gre.NewAlmanac()
		almanac.AddFact("testFact", "factValue")

		err := engine.HandleEvent("test-event", "my-rule", true, almanac, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		// Verify context fields
		if capturedContext.RuleName != "my-rule" {
			t.Errorf("Expected rule name 'my-rule', got '%s'", capturedContext.RuleName)
		}
		if capturedContext.Result != true {
			t.Errorf("Expected result true, got %v", capturedContext.Result)
		}
		if capturedContext.Almanac != almanac {
			t.Error("Expected almanac to match")
		}
		if capturedContext.Timestamp.IsZero() {
			t.Error("Expected timestamp to be set")
		}
		if len(capturedContext.Params) != 2 {
			t.Errorf("Expected 2 params, got %d", len(capturedContext.Params))
		}
		if capturedContext.Params["key1"] != "value1" {
			t.Errorf("Expected param 'key1' to be 'value1', got %v", capturedContext.Params["key1"])
		}
		if capturedContext.Params["key2"] != 42 {
			t.Errorf("Expected param 'key2' to be 42, got %v", capturedContext.Params["key2"])
		}
	})

	t.Run("passes correct context to handler", func(t *testing.T) {
		event := gre.Event{
			Name: "test-event",
			Mode: gre.EventModeSync,
			Params: map[string]interface{}{
				"handlerKey": "handlerValue",
			},
		}

		mockHandler := &MockEventHandler{}
		engine := gre.NewEngine()
		engine.SetEventHandler(mockHandler)
		engine.RegisterEvent(event)

		almanac := gre.NewAlmanac()
		almanac.AddFact("testFact", "factValue")

		err := engine.HandleEvent("test-event", "handler-rule", false, almanac, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(mockHandler.HandledContexts) != 1 {
			t.Fatalf("Expected 1 context, got %d", len(mockHandler.HandledContexts))
		}

		ctx := mockHandler.HandledContexts[0]
		if ctx.RuleName != "handler-rule" {
			t.Errorf("Expected rule name 'handler-rule', got '%s'", ctx.RuleName)
		}
		if ctx.Result != false {
			t.Errorf("Expected result false, got %v", ctx.Result)
		}
		if ctx.Almanac != almanac {
			t.Error("Expected almanac to match")
		}
		if ctx.Params["handlerKey"] != "handlerValue" {
			t.Errorf("Expected param 'handlerKey' to be 'handlerValue', got %v", ctx.Params["handlerKey"])
		}
	})
}

func TestEventWithNoAction(t *testing.T) {
	t.Run("works without action function", func(t *testing.T) {
		event := gre.Event{
			Name: "test-event",
			Mode: gre.EventModeSync,
		}

		mockHandler := &MockEventHandler{}
		engine := gre.NewEngine()
		engine.SetEventHandler(mockHandler)
		engine.RegisterEvent(event)

		almanac := gre.NewAlmanac()
		err := engine.HandleEvent("test-event", "test-rule", true, almanac, nil)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if len(mockHandler.HandledEvents) != 1 {
			t.Error("Expected handler to be called")
		}
	})
}

func TestEventIntegrationWithEngine(t *testing.T) {
	t.Run("executes sync action during engine run", func(t *testing.T) {
		actionExecuted := false

		event := gre.Event{
			Name: "success-event",
			Mode: gre.EventModeSync,
			Action: func(ctx gre.EventContext) error {
				actionExecuted = true
				return nil
			},
		}

		engine := gre.NewEngine()
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
			OnSuccess: []gre.RuleEvent{{Name: "success-event"}},
		}

		engine.AddRule(rule)

		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 25)

		e, err := engine.Run(almanac)

		results := e.ReduceResults()

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !results["test-rule"] {
			t.Error("Expected rule to pass")
		}
		if !actionExecuted {
			t.Error("Expected action to be executed during engine run")
		}
	})

	t.Run("executes async action during engine run", func(t *testing.T) {
		var mu sync.Mutex
		actionExecuted := false

		event := gre.Event{
			Name: "async-event",
			Mode: gre.EventModeAsync,
			Action: func(ctx gre.EventContext) error {
				time.Sleep(50 * time.Millisecond)
				mu.Lock()
				actionExecuted = true
				mu.Unlock()
				return nil
			},
		}

		engine := gre.NewEngine()
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
			OnSuccess: []gre.RuleEvent{{Name: "async-event"}},
		}

		engine.AddRule(rule)

		almanac := gre.NewAlmanac()
		almanac.AddFact("age", 25)

		e, err := engine.Run(almanac)

		results := e.ReduceResults()

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}
		if !results["test-rule"] {
			t.Error("Expected rule to pass")
		}

		// Action should not be executed yet
		mu.Lock()
		immediate := actionExecuted
		mu.Unlock()
		if immediate {
			t.Error("Expected async action to not execute immediately")
		}

		// Wait for async execution
		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		defer mu.Unlock()
		if !actionExecuted {
			t.Error("Expected async action to be executed after waiting")
		}
	})

	t.Run("action can access almanac facts", func(t *testing.T) {
		var capturedFact interface{}

		event := gre.Event{
			Name: "test-event",
			Mode: gre.EventModeSync,
			Action: func(ctx gre.EventContext) error {
				fact, _ := ctx.Almanac.GetFactValue("username", nil, "")
				capturedFact = fact
				return nil
			},
		}

		engine := gre.NewEngine()
		engine.RegisterEvent(event)

		rule := &gre.Rule{
			Name:     "test-rule",
			Priority: 10,
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{
						Condition: &gre.Condition{
							Fact:     "username",
							Operator: "equal",
							Value:    "john",
						},
					},
				},
			},
			OnSuccess: []gre.RuleEvent{{Name: "test-event"}},
		}

		engine.AddRule(rule)

		almanac := gre.NewAlmanac()
		almanac.AddFact("username", "john")

		_, err := engine.Run(almanac)

		if err != nil {
			t.Fatalf("Expected no error, got: %v", err)
		}

		if capturedFact != "john" {
			t.Errorf("Expected action to access fact 'john', got %v", capturedFact)
		}
	})
}

func TestRuleEventParamsMerging(t *testing.T) {
	t.Run("merges rule params with event defaults", func(t *testing.T) {
		var capturedParams map[string]interface{}

		event := gre.Event{
			Name: "test-params",
			Params: map[string]interface{}{
				"default":  "value",
				"override": "original",
			},
			Action: func(ctx gre.EventContext) error {
				capturedParams = ctx.Params
				return nil
			},
		}

		engine := gre.NewEngine()
		engine.RegisterEvent(event)

		rule := &gre.Rule{
			Name: "test-rule",
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{Condition: &gre.Condition{
						Fact:     "fact1",
						Operator: "equal",
						Value:    true,
					}},
				},
			},
			OnSuccess: []gre.RuleEvent{
				{
					Name: "test-params",
					Params: map[string]interface{}{
						"override": "new",
						"extra":    123,
					},
				},
			},
		}
		engine.AddRule(rule)

		almanac := gre.NewAlmanac()
		almanac.AddFact("fact1", true)

		_, err := engine.Run(almanac)
		if err != nil {
			t.Fatalf("Engine run failed: %v", err)
		}

		if capturedParams["default"] != "value" {
			t.Errorf("Expected default param 'value', got %v", capturedParams["default"])
		}
		if capturedParams["override"] != "new" {
			t.Errorf("Expected overridden param 'new', got %v", capturedParams["override"])
		}
		if capturedParams["extra"] != 123 {
			t.Errorf("Expected extra param 123, got %v", capturedParams["extra"])
		}
	})
}

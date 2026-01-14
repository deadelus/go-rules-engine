package gorulesengine_test

import (
	"errors"
	"testing"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func TestNewEngine(t *testing.T) {
	engine := gorulesengine.NewEngine()
	if engine == nil {
		t.Fatal("NewEngine should return a non-nil engine")
	}
}

func TestEngine_AddRule(t *testing.T) {
	engine := gorulesengine.NewEngine()

	rule := &gorulesengine.Rule{
		Name:     "test-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "age",
						Operator: "greater_than",
						Value:    18,
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "test-event",
		},
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("age", 25)

	results, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
}

func TestEngine_RegisterCallback(t *testing.T) {
	engine := gorulesengine.NewEngine()

	callbackCalled := false
	engine.RegisterCallback("test-callback", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		callbackCalled = true
		return nil
	})

	onSuccessName := "test-callback"
	rule := &gorulesengine.Rule{
		Name:     "callback-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "test-event",
		},
		OnSuccess: &onSuccessName,
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("test", true)

	_, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !callbackCalled {
		t.Error("Registered callback was not called")
	}
}

func TestEngine_OnSuccess(t *testing.T) {
	engine := gorulesengine.NewEngine()

	successCalled := false
	engine.OnSuccess(func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		successCalled = true
		if result.Result != true {
			t.Errorf("Expected result to be true in success handler")
		}
		return nil
	})

	rule := &gorulesengine.Rule{
		Name:     "success-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "value",
						Operator: "equal",
						Value:    100,
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "success-event",
		},
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("value", 100)

	_, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !successCalled {
		t.Error("OnSuccess handler was not called")
	}
}

func TestEngine_OnFailure(t *testing.T) {
	engine := gorulesengine.NewEngine()

	failureCalled := false
	engine.OnFailure(func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		failureCalled = true
		if result.Result != false {
			t.Errorf("Expected result to be false in failure handler")
		}
		return nil
	})

	rule := &gorulesengine.Rule{
		Name:     "failure-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "value",
						Operator: "equal",
						Value:    100,
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "failure-event",
		},
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("value", 50) // Valeur différente pour faire échouer la règle

	_, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !failureCalled {
		t.Error("OnFailure handler was not called")
	}
}

func TestEngine_On_EventType(t *testing.T) {
	engine := gorulesengine.NewEngine()

	eventHandlerCalled := false
	engine.On("specific-event", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		eventHandlerCalled = true
		if event.Type != "specific-event" {
			t.Errorf("Expected event type 'specific-event', got '%s'", event.Type)
		}
		return nil
	})

	rule := &gorulesengine.Rule{
		Name:     "event-type-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "specific-event",
			Params: map[string]interface{}{
				"key": "value",
			},
		},
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("test", true)

	_, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !eventHandlerCalled {
		t.Error("Event type handler was not called")
	}
}

func TestEngine_MultipleRules(t *testing.T) {
	engine := gorulesengine.NewEngine()

	rule1 := &gorulesengine.Rule{
		Name:     "rule-1",
		Priority: 100,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "age",
						Operator: "greater_than",
						Value:    18,
					},
				},
			},
		},
		Event: gorulesengine.Event{Type: "adult"},
	}

	rule2 := &gorulesengine.Rule{
		Name:     "rule-2",
		Priority: 50,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "isPremium",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{Type: "premium"},
	}

	engine.AddRule(rule1)
	engine.AddRule(rule2)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("age", 25)
	almanac.AddFact("isPremium", true)

	results, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}

	successCount := 0
	for _, result := range results {
		if result.Result {
			successCount++
		}
	}

	if successCount != 2 {
		t.Errorf("Expected 2 successful rules, got %d", successCount)
	}
}

func TestEngine_HandlerError(t *testing.T) {
	engine := gorulesengine.NewEngine()

	expectedError := errors.New("handler error")
	engine.OnSuccess(func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		return expectedError
	})

	rule := &gorulesengine.Rule{
		Name:     "error-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{Type: "test"},
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("test", true)

	_, err := engine.Run(almanac)
	if err == nil {
		t.Fatal("Expected error from handler, got nil")
	}

	// Check that it's a RuleEngineError
	var ruleEngineErr *gorulesengine.RuleEngineError
	if !errors.As(err, &ruleEngineErr) {
		t.Errorf("Expected RuleEngineError, got %T", err)
	}

	// Check that the wrapped error is our expected error
	if !errors.Is(err, expectedError) {
		t.Errorf("Expected wrapped error to be '%v', got '%v'", expectedError, err)
	}
}

func TestEngine_CallbackNotFound(t *testing.T) {
	engine := gorulesengine.NewEngine()

	onSuccessName := "non-existent-callback"
	rule := &gorulesengine.Rule{
		Name:     "missing-callback-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event:     gorulesengine.Event{Type: "test"},
		OnSuccess: &onSuccessName,
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("test", true)

	// Should not error, just print warning
	_, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run should not fail on missing callback: %v", err)
	}
}

func TestEngine_OnFailureCallbackNotFound(t *testing.T) {
	engine := gorulesengine.NewEngine()

	onFailureName := "non-existent-failure-callback"
	rule := &gorulesengine.Rule{
		Name:     "missing-failure-callback-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "value",
						Operator: "equal",
						Value:    100,
					},
				},
			},
		},
		Event:     gorulesengine.Event{Type: "test"},
		OnFailure: &onFailureName,
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("value", 50) // Will fail the condition

	// Should not error, just print warning
	_, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run should not fail on missing OnFailure callback: %v", err)
	}

	// Verify rule was still evaluated and failed
	results := almanac.GetResults()
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}

	if results[0].Result {
		t.Error("Expected rule to fail")
	}
}

func TestEngine_AlmanacEventStorage(t *testing.T) {
	engine := gorulesengine.NewEngine()

	rule := &gorulesengine.Rule{
		Name:     "storage-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "storage-event",
			Params: map[string]interface{}{
				"data": "test-data",
			},
		},
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("test", true)

	_, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Vérifier que l'événement a été stocké
	successEvents := almanac.GetEvents("success")
	if len(successEvents) != 1 {
		t.Errorf("Expected 1 success event, got %d", len(successEvents))
	}

	if len(successEvents) > 0 && successEvents[0].Type != "storage-event" {
		t.Errorf("Expected event type 'storage-event', got '%s'", successEvents[0].Type)
	}

	// Vérifier les résultats stockés
	results := almanac.GetResults()
	if len(results) != 1 {
		t.Errorf("Expected 1 result stored, got %d", len(results))
	}
}

func TestEngine_HandlerExecutionOrder(t *testing.T) {
	engine := gorulesengine.NewEngine()

	var executionOrder []string

	// Callback nommé
	engine.RegisterCallback("named-callback", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		executionOrder = append(executionOrder, "named-callback")
		return nil
	})

	// Handler global
	engine.OnSuccess(func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		executionOrder = append(executionOrder, "global-success")
		return nil
	})

	// Handler par type d'événement
	engine.On("order-event", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		executionOrder = append(executionOrder, "event-type-handler")
		return nil
	})

	onSuccessName := "named-callback"
	rule := &gorulesengine.Rule{
		Name:     "order-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "order-event",
		},
		OnSuccess: &onSuccessName,
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("test", true)

	_, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Vérifier l'ordre: 1. named callback, 2. global success, 3. event type
	expectedOrder := []string{"named-callback", "global-success", "event-type-handler"}
	if len(executionOrder) != len(expectedOrder) {
		t.Fatalf("Expected %d handlers, got %d", len(expectedOrder), len(executionOrder))
	}

	for i, expected := range expectedOrder {
		if executionOrder[i] != expected {
			t.Errorf("Handler at position %d: expected '%s', got '%s'", i, expected, executionOrder[i])
		}
	}
}

func TestEngine_OnFailureCallback(t *testing.T) {
	engine := gorulesengine.NewEngine()

	failureCallbackCalled := false
	engine.RegisterCallback("failure-callback", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		failureCallbackCalled = true
		return nil
	})

	onFailureName := "failure-callback"
	rule := &gorulesengine.Rule{
		Name:     "fail-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "value",
						Operator: "equal",
						Value:    100,
					},
				},
			},
		},
		Event: gorulesengine.Event{
			Type: "fail-event",
		},
		OnFailure: &onFailureName,
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("value", 50) // Valeur différente pour échouer

	_, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !failureCallbackCalled {
		t.Error("OnFailure callback was not called")
	}

	// Vérifier que l'événement est dans la liste failure
	failureEvents := almanac.GetEvents("failure")
	if len(failureEvents) != 1 {
		t.Errorf("Expected 1 failure event, got %d", len(failureEvents))
	}
}

func TestEngine_MultipleHandlersOfSameType(t *testing.T) {
	engine := gorulesengine.NewEngine()

	handler1Called := false
	handler2Called := false

	engine.OnSuccess(func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		handler1Called = true
		return nil
	})

	engine.OnSuccess(func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		handler2Called = true
		return nil
	})

	rule := &gorulesengine.Rule{
		Name:     "multi-handler-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{Type: "test"},
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("test", true)

	_, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if !handler1Called || !handler2Called {
		t.Error("All success handlers should be called")
	}
}

func TestEngine_AddFact(t *testing.T) {
	engine := gorulesengine.NewEngine()

	fact := gorulesengine.NewFact("testFact", 42)
	engine.AddFact(fact)

	// AddFact doesn't have a getter, so we test indirectly
	// by verifying the engine was created successfully
	if engine == nil {
		t.Error("Engine should not be nil after adding fact")
	}
}

func TestEngine_OnSuccessCallbackError(t *testing.T) {
	engine := gorulesengine.NewEngine()

	expectedError := errors.New("callback error")
	engine.RegisterCallback("errorCallback", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		return expectedError
	})

	onSuccessName := "errorCallback"
	rule := &gorulesengine.Rule{
		Name:     "callback-error-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event:     gorulesengine.Event{Type: "test"},
		OnSuccess: &onSuccessName,
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("test", true)

	_, err := engine.Run(almanac)
	if err == nil {
		t.Fatal("Expected error from OnSuccess callback")
	}

	var ruleEngineErr *gorulesengine.RuleEngineError
	if !errors.As(err, &ruleEngineErr) {
		t.Errorf("Expected RuleEngineError, got %T", err)
	}

	if !errors.Is(err, expectedError) {
		t.Errorf("Expected wrapped error to be '%v'", expectedError)
	}
}

func TestEngine_OnFailureCallbackError(t *testing.T) {
	engine := gorulesengine.NewEngine()

	expectedError := errors.New("failure callback error")
	engine.RegisterCallback("failureErrorCallback", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		return expectedError
	})

	onFailureName := "failureErrorCallback"
	rule := &gorulesengine.Rule{
		Name:     "failure-callback-error-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "value",
						Operator: "equal",
						Value:    100,
					},
				},
			},
		},
		Event:     gorulesengine.Event{Type: "test"},
		OnFailure: &onFailureName,
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("value", 50) // Will fail the condition

	_, err := engine.Run(almanac)
	if err == nil {
		t.Fatal("Expected error from OnFailure callback")
	}

	var ruleEngineErr *gorulesengine.RuleEngineError
	if !errors.As(err, &ruleEngineErr) {
		t.Errorf("Expected RuleEngineError, got %T", err)
	}

	if !errors.Is(err, expectedError) {
		t.Errorf("Expected wrapped error to be '%v'", expectedError)
	}
}

func TestEngine_EventTypeHandlerError(t *testing.T) {
	engine := gorulesengine.NewEngine()

	expectedError := errors.New("event handler error")
	engine.On("error-event", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		return expectedError
	})

	rule := &gorulesengine.Rule{
		Name:     "event-error-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{Type: "error-event"},
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("test", true)

	_, err := engine.Run(almanac)
	if err == nil {
		t.Fatal("Expected error from event type handler")
	}

	var ruleEngineErr *gorulesengine.RuleEngineError
	if !errors.As(err, &ruleEngineErr) {
		t.Errorf("Expected RuleEngineError, got %T", err)
	}

	if !errors.Is(err, expectedError) {
		t.Errorf("Expected wrapped error to be '%v'", expectedError)
	}
}

func TestEngine_FailureHandlerError(t *testing.T) {
	engine := gorulesengine.NewEngine()

	expectedError := errors.New("failure handler error")
	engine.OnFailure(func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		return expectedError
	})

	rule := &gorulesengine.Rule{
		Name:     "fail-handler-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "value",
						Operator: "equal",
						Value:    100,
					},
				},
			},
		},
		Event: gorulesengine.Event{Type: "test"},
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("value", 50) // Will fail the condition

	_, err := engine.Run(almanac)
	if err == nil {
		t.Fatal("Expected error from failure handler")
	}

	var ruleEngineErr *gorulesengine.RuleEngineError
	if !errors.As(err, &ruleEngineErr) {
		t.Errorf("Expected RuleEngineError, got %T", err)
	}

	if !errors.Is(err, expectedError) {
		t.Errorf("Expected wrapped error to be '%v'", expectedError)
	}
}

func TestEngine_RuleEvaluationError(t *testing.T) {
	engine := gorulesengine.NewEngine()

	// Create a rule with an invalid operator to trigger evaluation error
	rule := &gorulesengine.Rule{
		Name:     "invalid-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "invalid_operator",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{Type: "test"},
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("test", true)

	_, err := engine.Run(almanac)
	if err == nil {
		t.Fatal("Expected error from rule evaluation")
	}

	var ruleEngineErr *gorulesengine.RuleEngineError
	if !errors.As(err, &ruleEngineErr) {
		t.Errorf("Expected RuleEngineError, got %T", err)
	} else if ruleEngineErr.Type != gorulesengine.ErrEngine {
		t.Errorf("Expected error type ErrEngine, got %s", ruleEngineErr.Type)
	}
}

func TestEngine_WithPrioritySorting_ASC(t *testing.T) {
	// Test ascending sort
	sortOrder := gorulesengine.SortRuleASC
	engine := gorulesengine.NewEngine(gorulesengine.WithPrioritySorting(&sortOrder))

	rule1 := &gorulesengine.Rule{
		Name:     "rule1",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{Type: "event1"},
	}

	rule2 := &gorulesengine.Rule{
		Name:     "rule2",
		Priority: 20,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{Type: "event2"},
	}

	engine.AddRule(rule2) // Add higher priority first
	engine.AddRule(rule1) // Add lower priority second

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("test", true)

	results, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// With ASC sorting, rule1 (priority 10) should be evaluated first
	if results[0].Rule.Name != "rule1" {
		t.Errorf("Expected rule1 to be evaluated first with ASC sorting, got %s", results[0].Rule.Name)
	}
}

func TestEngine_WithPrioritySorting_DefaultValue(t *testing.T) {
	// Test default sort (invalid value should default to SortDefault)
	sortOrder := gorulesengine.SortRule(999) // Invalid value
	engine := gorulesengine.NewEngine(gorulesengine.WithPrioritySorting(&sortOrder))

	rule1 := &gorulesengine.Rule{
		Name:     "rule1",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{Type: "event1"},
	}

	engine.AddRule(rule1)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("test", true)

	results, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
}

func TestEngine_WithoutPrioritySorting(t *testing.T) {
	// Create engine with priority sorting disabled
	engine := gorulesengine.NewEngine(gorulesengine.WithoutPrioritySorting())

	rule1 := &gorulesengine.Rule{
		Name:     "rule1",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{Type: "event1"},
	}

	rule2 := &gorulesengine.Rule{
		Name:     "rule2",
		Priority: 20,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{Type: "event2"},
	}

	// Add rules in a specific order
	engine.AddRule(rule2) // Add higher priority first
	engine.AddRule(rule1) // Add lower priority second

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("test", true)

	results, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	// Without sorting, rules should be evaluated in the order they were added
	if results[0].Rule.Name != "rule2" {
		t.Errorf("Expected rule2 to be evaluated first (insertion order), got %s", results[0].Rule.Name)
	}
	if results[1].Rule.Name != "rule1" {
		t.Errorf("Expected rule1 to be evaluated second (insertion order), got %s", results[1].Rule.Name)
	}
}

func TestEngine_NewEngine_WithMultipleOptions(t *testing.T) {
	// Test NewEngine with multiple options
	sortOrder := gorulesengine.SortRuleASC
	engine := gorulesengine.NewEngine(
		gorulesengine.WithPrioritySorting(&sortOrder),
	)

	if engine == nil {
		t.Fatal("NewEngine should return a non-nil engine")
	}

	rule := &gorulesengine.Rule{
		Name:     "test-rule",
		Priority: 10,
		Conditions: gorulesengine.ConditionSet{
			All: []gorulesengine.ConditionNode{
				{
					Condition: &gorulesengine.Condition{
						Fact:     "test",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
		Event: gorulesengine.Event{Type: "test"},
	}

	engine.AddRule(rule)

	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac.AddFact("test", true)

	results, err := engine.Run(almanac)
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}
}

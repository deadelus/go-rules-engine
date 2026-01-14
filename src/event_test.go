package gorulesengine

import "testing"

func TestEventHandlers_RegisterHandler(t *testing.T) {
	handlers := &EventHandlers{}

	callCount := 0
	handler := func(event Event, almanac *Almanac, ruleResult RuleResult) error {
		callCount++
		return nil
	}

	handlers.RegisterHandler("test-event", handler)

	registered := handlers.GetHandlers("test-event")
	if len(registered) != 1 {
		t.Errorf("Expected 1 handler, got %d", len(registered))
	}
}

func TestEventHandlers_RegisterMultipleHandlers(t *testing.T) {
	handlers := &EventHandlers{}

	handler1 := func(event Event, almanac *Almanac, ruleResult RuleResult) error {
		return nil
	}

	handler2 := func(event Event, almanac *Almanac, ruleResult RuleResult) error {
		return nil
	}

	handlers.RegisterHandler("test-event", handler1)
	handlers.RegisterHandler("test-event", handler2)

	registered := handlers.GetHandlers("test-event")
	if len(registered) != 2 {
		t.Errorf("Expected 2 handlers, got %d", len(registered))
	}
}

func TestEventHandlers_GetHandlers_EmptyRegistry(t *testing.T) {
	handlers := &EventHandlers{}

	registered := handlers.GetHandlers("non-existent")
	if registered != nil {
		t.Errorf("Expected nil for non-existent event type, got %v", registered)
	}
}

func TestEventHandlers_GetHandlers_NonExistentType(t *testing.T) {
	handlers := &EventHandlers{}

	handler := func(event Event, almanac *Almanac, ruleResult RuleResult) error {
		return nil
	}

	handlers.RegisterHandler("test-event", handler)

	registered := handlers.GetHandlers("other-event")
	if registered != nil {
		t.Errorf("Expected nil for non-registered event type, got %v", registered)
	}
}

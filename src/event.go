package gorulesengine

const (
	// EventSuccess represents a success event type for rules that matched.
	EventSuccess = "success"
	// EventFailure represents a failure event type for rules that did not match.
	EventFailure = "failure"
)

// Event represents an event triggered by a rule when its conditions are met.
// Events can carry additional parameters in the Params map.
type Event struct {
	Type   string                 `json:"type"`             // The event type identifier
	Params map[string]interface{} `json:"params,omitempty"` // Optional parameters passed with the event
}

// EventHandler is a callback function invoked when an event occurs.
// It receives the event, the current almanac state, and the rule result.
// Handlers can return an error to stop further processing.
type EventHandler func(event Event, almanac *Almanac, ruleResult RuleResult) error

// EventHandlers manages a registry of event handlers organized by event type.
// It allows multiple handlers to be registered for the same event type.
type EventHandlers struct {
	handlers map[string][]EventHandler
}

// RegisterHandler registers an event handler for a specific event type.
// Multiple handlers can be registered for the same event type and will be invoked in order.
func (e *EventHandlers) RegisterHandler(eventType string, handler EventHandler) {
	if e.handlers == nil {
		e.handlers = make(map[string][]EventHandler)
	}
	e.handlers[eventType] = append(e.handlers[eventType], handler)
}

// GetHandlers retrieves all handlers registered for a specific event type.
// Returns nil if no handlers are registered for the given event type.
func (e *EventHandlers) GetHandlers(eventType string) []EventHandler {
	if e.handlers == nil {
		return nil
	}
	return e.handlers[eventType]
}

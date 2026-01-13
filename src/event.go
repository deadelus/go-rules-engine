package gorulesengine

const (
	// EventSuccess represents a success event type
	EventSuccess = "success"
	// EventFailure represents a failure event type
	EventFailure = "failure"
)

// Event represents an event triggered by a rule
type Event struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params,omitempty"`
}

// EventHandler est un callback appelé quand un événement se produit
type EventHandler func(event Event, almanac *Almanac, ruleResult RuleResult) error

// EventHandlerRegistry manages event handlers for different event types
type EventHandlers struct {
	handlers map[string][]EventHandler
}

// RegisterHandler registers an event handler for a specific event type
func (e *EventHandlers) RegisterHandler(eventType string, handler EventHandler) {
	if e.handlers == nil {
		e.handlers = make(map[string][]EventHandler)
	}
	e.handlers[eventType] = append(e.handlers[eventType], handler)
}

// GetHandlers retrieves handlers for a specific event type
func (e *EventHandlers) GetHandlers(eventType string) []EventHandler {
	if e.handlers == nil {
		return nil
	}
	return e.handlers[eventType]
}

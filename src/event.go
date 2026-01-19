package gorulesengine

import "time"

// EventMode defines how an event should be executed
type EventMode int

const (
	// EventModeSync executes the event synchronously (blocking)
	EventModeSync EventMode = iota
	// EventModeAsync executes the event asynchronously (non-blocking)
	EventModeAsync
)

// EventContext provides context information about the rule execution that triggered an event
type EventContext struct {
	RuleName  string                 // Name of the rule that triggered the event
	Result    bool                   // Result of the rule evaluation (true=success, false=failure)
	Almanac   *Almanac               // Reference to the almanac used for evaluation
	Timestamp time.Time              // When the event was triggered
	Params    map[string]interface{} // Additional parameters from the event
}

// Event represents an event triggered by a rule when its conditions are met.
// Events can carry additional parameters in the Params map.
type Event struct {
	Name   string                   // Name of the event
	Params map[string]interface{}   // Optional parameters passed with the event
	Action func(EventContext) error // Action to execute when the event is handled
	Mode   EventMode                // Execution mode (sync or async)
}

// EventHandler defines the interface for handling events triggered by rules.
type EventHandler interface {
	Handle(event Event, ctx EventContext) error
}

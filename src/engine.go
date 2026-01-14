package gorulesengine

import (
	"fmt"
	"sort"
)

// SortRule defines the sort order for rules.
type SortRule int

const (
	// SortByPriority is the string constant for priority-based sorting
	SortByPriority = "priority"
	// SortDefault is the default sort order
	SortDefault SortRule = iota
	// SortRuleASC sorts rules in ascending order
	SortRuleASC
	// SortRuleDESC sorts rules in descending order
	SortRuleDESC
)

// Engine is the core rules engine that manages rules, facts, and event handlers.
// It evaluates rules against facts and triggers events when rules match.
type Engine struct {
	rules []*Rule
	facts map[FactID]*Fact

	// Callbacks registered by name
	callbacks map[string]EventHandler

	// Handlers for events
	successHandlers []EventHandler
	failureHandlers []EventHandler

	// Handler mapping for event types
	eventHandlers map[string][]EventHandler

	// Additional engine options
	options map[string]interface{}
}

// EngineOption defines a function type for configuring the Engine.
type EngineOption func(*Engine)

// WithPrioritySorting configures the engine to sort rules by priority before evaluation.
func WithPrioritySorting(o *SortRule) EngineOption {
	var order SortRule

	if o == nil {
		order = SortRuleDESC
	} else {
		switch *o {
		case SortRuleASC, SortRuleDESC:
			order = *o
		default:
			order = SortDefault
		}
	}

	return func(e *Engine) {
		options := e.options
		if options == nil {
			options = make(map[string]interface{})
			e.options = options
		}

		e.options[SortByPriority] = order
	}
}

// WithoutPrioritySorting configures the engine to not sort rules by priority.
func WithoutPrioritySorting() EngineOption {
	return func(e *Engine) {
		delete(e.options, SortByPriority)
	}
}

// NewEngine creates a new rules engine instance
func NewEngine(opts ...EngineOption) *Engine {
	e := &Engine{}
	WithPrioritySorting(nil)(e) // Default to priority sorting

	for _, opt := range opts {
		opt(e)
	}

	return e
}

// AddRule adds a rule to the engine
func (e *Engine) AddRule(rule *Rule) {
	e.rules = append(e.rules, rule)
}

// AddFact adds a fact to the engine
func (e *Engine) AddFact(fact *Fact) {
	if e.facts == nil {
		e.facts = make(map[FactID]*Fact)
	}
	e.facts[fact.ID()] = fact
}

// RegisterCallback registers a named callback that can be referenced by rules.
// Callbacks are invoked when rules succeed or fail, as specified in the rule's OnSuccess or OnFailure fields.
func (e *Engine) RegisterCallback(name string, handler EventHandler) {
	if e.callbacks == nil {
		e.callbacks = make(map[string]EventHandler)
	}
	e.callbacks[name] = handler
}

// OnSuccess registers a global handler that is called for every successful rule evaluation.
// Multiple success handlers can be registered and will all be invoked in order.
func (e *Engine) OnSuccess(handler EventHandler) {
	e.successHandlers = append(e.successHandlers, handler)
}

// OnFailure registers a global handler that is called for every failed rule evaluation.
// Multiple failure handlers can be registered and will all be invoked in order.
func (e *Engine) OnFailure(handler EventHandler) {
	e.failureHandlers = append(e.failureHandlers, handler)
}

// On registers a handler for a specific event type.
// When a rule triggers an event of this type, the handler will be invoked.
// Multiple handlers can be registered for the same event type.
func (e *Engine) On(eventType string, handler EventHandler) {
	if e.eventHandlers == nil {
		e.eventHandlers = make(map[string][]EventHandler)
	}
	e.eventHandlers[eventType] = append(e.eventHandlers[eventType], handler)
}

// Run executes all rules in the engine against the provided almanac.
// Rules are evaluated in priority order (higher priority first).
// Returns a slice of RuleResults containing the outcome of each rule evaluation.
// If any error occurs during evaluation, execution stops and the error is returned.
func (e *Engine) Run(almanac *Almanac) ([]RuleResult, error) {
	var results []RuleResult

	// Sort rules by priority if configured
	e.sortRulesByPriority()

	for _, rule := range e.rules {
		// Évaluer la règle
		success, err := rule.Conditions.Evaluate(almanac)
		if err != nil {
			return nil, &RuleEngineError{
				Type: ErrEngine,
				Msg:  fmt.Sprintf("Error evaluating rule '%s': %v", rule.Name, err),
				Err:  err,
			}
		}

		// Créer le résultat
		ruleResult := RuleResult{
			Event:  rule.Event,
			Rule:   rule,
			Result: success,
		}

		// Stocker dans l'almanac
		almanac.AddResult(ruleResult)

		if success {
			// Ajouter à la liste des événements success
			almanac.AddSuccessEvent(rule.Event)

			// 1. Appeler le callback OnSuccess de la règle
			if rule.OnSuccess != nil {
				if handler, exists := e.callbacks[*rule.OnSuccess]; exists {
					if err := handler(rule.Event, almanac, ruleResult); err != nil {
						return nil, &RuleEngineError{
							Type: ErrEngine,
							Msg:  fmt.Sprintf("Error in OnSuccess callback for rule '%s': %v", rule.Name, err),
							Err:  err,
						}
					}
				} else {
					fmt.Printf("Warning: OnSuccess callback '%s' not registered\n", *rule.OnSuccess)
				}
			}

			// 2. Appeler les handlers globaux "success"
			for _, handler := range e.successHandlers {
				if err := handler(rule.Event, almanac, ruleResult); err != nil {
					return nil, &RuleEngineError{
						Type: ErrEngine,
						Msg:  fmt.Sprintf("Error in success handler for rule '%s': %v", rule.Name, err),
						Err:  err,
					}
				}
			}

			// 3. Appeler les handlers spécifiques à ce type d'événement
			if handlers, exists := e.eventHandlers[rule.Event.Type]; exists {
				for _, handler := range handlers {
					if err := handler(rule.Event, almanac, ruleResult); err != nil {
						return nil, &RuleEngineError{
							Type: ErrEngine,
							Msg:  fmt.Sprintf("Error in event handler for event type '%s' in rule '%s': %v", rule.Event.Type, rule.Name, err),
							Err:  err,
						}
					}
				}
			}
		} else {
			// Ajouter à la liste des événements failure
			almanac.AddFailureEvent(rule.Event)

			// 1. Appeler le callback OnFailure de la règle
			if rule.OnFailure != nil {
				if handler, exists := e.callbacks[*rule.OnFailure]; exists {
					if err := handler(rule.Event, almanac, ruleResult); err != nil {
						return nil, &RuleEngineError{
							Type: ErrEngine,
							Msg:  fmt.Sprintf("Error in OnFailure callback for rule '%s': %v", rule.Name, err),
							Err:  err,
						}
					}
				} else {
					fmt.Printf("Warning: OnFailure callback '%s' not registered\n", *rule.OnFailure)
				}
			}

			// 2. Appeler les handlers globaux "failure"
			for _, handler := range e.failureHandlers {
				if err := handler(rule.Event, almanac, ruleResult); err != nil {
					return nil, &RuleEngineError{
						Type: ErrEngine,
						Msg:  fmt.Sprintf("Error in failure handler for rule '%s': %v", rule.Name, err),
						Err:  err,
					}
				}
			}
		}

		results = append(results, ruleResult)
	}

	return results, nil
}

// sortRulesByPriority sorts the engine's rules by their priority in descending order.
func (e *Engine) sortRulesByPriority() {
	if e.options[SortByPriority] != nil {
		// Sort by priority
		sort.SliceStable(e.rules, func(i, j int) bool {
			if e.options[SortByPriority] == SortRuleASC {
				return e.rules[i].Priority < e.rules[j].Priority
			}
			return e.rules[i].Priority > e.rules[j].Priority
		})
	}
}

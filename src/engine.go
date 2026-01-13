package gorulesengine

import "fmt"

type Engine struct {
	rules []*Rule
	facts map[FactID]*Fact

	// Registry des callbacks par nom
	callbacks map[string]EventHandler

	// Handlers for events
	successHandlers []EventHandler
	failureHandlers []EventHandler

	// Handler mapping for event types
	eventHandlers map[string][]EventHandler
}

// NewEngine creates a new rules engine instance
func NewEngine() *Engine {
	return &Engine{}
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

// Enregistrer un callback avec un nom
func (e *Engine) RegisterCallback(name string, handler EventHandler) {
	if e.callbacks == nil {
		e.callbacks = make(map[string]EventHandler)
	}
	e.callbacks[name] = handler
}

// Enregistrer un handler global pour tous les succès
func (e *Engine) OnSuccess(handler EventHandler) {
	e.successHandlers = append(e.successHandlers, handler)
}

// Enregistrer un handler global pour tous les échecs
func (e *Engine) OnFailure(handler EventHandler) {
	e.failureHandlers = append(e.failureHandlers, handler)
}

// Enregistrer un handler pour un type d'événement spécifique
func (e *Engine) On(eventType string, handler EventHandler) {
	if e.eventHandlers == nil {
		e.eventHandlers = make(map[string][]EventHandler)
	}
	e.eventHandlers[eventType] = append(e.eventHandlers[eventType], handler)
}

// Exécuter l'engine
func (e *Engine) Run(almanac *Almanac) ([]RuleResult, error) {
	var results []RuleResult

	for _, rule := range e.rules {
		// Évaluer la règle
		success, err := rule.Conditions.Evaluate(almanac)
		if err != nil {
			return nil, err
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
			almanac.AddEvent(rule.Event, EventSuccess)

			// 1. Appeler le callback OnSuccess de la règle
			if rule.OnSuccess != nil {
				if handler, exists := e.callbacks[*rule.OnSuccess]; exists {
					if err := handler(rule.Event, almanac, ruleResult); err != nil {
						return nil, err
					}
				} else {
					// Callback non trouvé - warning ou erreur ?
					fmt.Printf("Warning: callback '%s' not registered\n", *rule.OnSuccess)
				}
			}

			// 2. Appeler les handlers globaux "success"
			for _, handler := range e.successHandlers {
				if err := handler(rule.Event, almanac, ruleResult); err != nil {
					return nil, err
				}
			}

			// 3. Appeler les handlers spécifiques à ce type d'événement
			if handlers, exists := e.eventHandlers[rule.Event.Type]; exists {
				for _, handler := range handlers {
					if err := handler(rule.Event, almanac, ruleResult); err != nil {
						return nil, err
					}
				}
			}
		} else {
			// Ajouter à la liste des événements failure
			almanac.AddEvent(rule.Event, EventFailure)

			// 1. Appeler le callback OnFailure de la règle
			if rule.OnFailure != nil {
				if handler, exists := e.callbacks[*rule.OnFailure]; exists {
					if err := handler(rule.Event, almanac, ruleResult); err != nil {
						return nil, err
					}
				} else {
					fmt.Printf("Warning: callback '%s' not registered\n", *rule.OnFailure)
				}
			}

			// 2. Appeler les handlers globaux "failure"
			for _, handler := range e.failureHandlers {
				if err := handler(rule.Event, almanac, ruleResult); err != nil {
					return nil, err
				}
			}
		}

		results = append(results, ruleResult)
	}

	return results, nil
}

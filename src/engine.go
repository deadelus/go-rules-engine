package gorulesengine

import (
	"fmt"
	"sort"
	"sync"
	"time"
)

// SortRule defines the sort order for rules.
type SortRule int

const (
	// SortByPriority is the string constant for priority-based sorting
	SortByPriority = "priority"
	// EngineOptionKeyCacheConditions is the option key for condition results caching
	EngineOptionKeyCacheConditions = "cacheConditions"
	// EngineOptionKeySmartSkip is the option key for enabling smart skip of rules
	EngineOptionKeySmartSkip = "smartSkip"
	// EngineOptionKeyAuditTrace is the option key for enabling rich audit trace
	EngineOptionKeyAuditTrace = "auditTrace"
	// EngineOptionKeyParallel is the option key for enabling parallel execution
	EngineOptionKeyParallel = "parallel"
	// EngineOptionKeyWorkerCount is the option key for specifying the number of workers for parallel execution
	EngineOptionKeyWorkerCount = "workerCount"
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
	// Mutex for thread-safe access to rules and events
	mu sync.RWMutex

	// Registered rules in the engine
	rules []*Rule

	// Cached results of rule evaluations
	results map[string]*RuleResult

	// The almanac used for the last run
	almanac *Almanac

	// Global event handler for all events
	eventHandler EventHandler

	// Registered event handlers by name
	events map[string]Event

	// Metrics collector for monitoring
	metrics MetricsCollector

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

// WithConditionCaching enables condition caching for all rules evaluated by the engine.
func WithConditionCaching() EngineOption {
	return func(e *Engine) {
		if e == nil {
			return
		}
		if e.options == nil {
			e.options = make(map[string]interface{})
		}
		e.options[EngineOptionKeyCacheConditions] = true
	}
}

// WithoutConditionCaching disables condition caching for all rules evaluated by the engine.
func WithoutConditionCaching() EngineOption {
	return func(e *Engine) {
		if e == nil {
			return
		}
		if e.options == nil {
			e.options = make(map[string]interface{})
		}
		e.options[EngineOptionKeyCacheConditions] = false
	}
}

// WithAuditTrace enables detailed audit trace for rule evaluations.
func WithAuditTrace() EngineOption {
	return func(e *Engine) {
		if e == nil {
			return
		}
		if e.options == nil {
			e.options = make(map[string]interface{})
		}
		e.options[EngineOptionKeyAuditTrace] = true
	}
}

// WithoutAuditTrace disables detailed audit trace for rule evaluations.
func WithoutAuditTrace() EngineOption {
	return func(e *Engine) {
		if e == nil {
			return
		}
		if e.options == nil {
			e.options = make(map[string]interface{})
		}
		e.options[EngineOptionKeyAuditTrace] = false
	}
}

// WithSmartSkip enables skipping rules that depend on facts not present in the almanac.
func WithSmartSkip() EngineOption {
	return func(e *Engine) {
		if e == nil {
			return
		}
		if e.options == nil {
			e.options = make(map[string]interface{})
		}
		e.options[EngineOptionKeySmartSkip] = true
	}
}

// WithMetrics configures the engine with a metrics collector.
func WithMetrics(collector MetricsCollector) EngineOption {
	return func(e *Engine) {
		if e == nil {
			return
		}
		e.metrics = collector
	}
}

// WithParallelExecution enables parallel execution of rules.
func WithParallelExecution(workers int) EngineOption {
	return func(e *Engine) {
		if e == nil {
			return
		}
		if e.options == nil {
			e.options = make(map[string]interface{})
		}
		e.options[EngineOptionKeyParallel] = true
		e.options[EngineOptionKeyWorkerCount] = workers
	}
}

// WithoutParallelExecution disables parallel execution of rules.
func WithoutParallelExecution() EngineOption {
	return func(e *Engine) {
		if e == nil {
			return
		}
		if e.options == nil {
			e.options = make(map[string]interface{})
		}
		e.options[EngineOptionKeyParallel] = false
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

// AddRules adds multiple rules to the engine
func (e *Engine) AddRules(rules ...*Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, rule := range rules {
		rule.Compile()
	}
	e.rules = append(e.rules, rules...)
}

// AddRule adds a rule to the engine
func (e *Engine) AddRule(rule *Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()

	rule.Compile()
	e.rules = append(e.rules, rule)
}

// SetRules replaces all rules in the engine with the provided ones.
func (e *Engine) SetRules(rules []*Rule) {
	e.mu.Lock()
	defer e.mu.Unlock()

	for _, rule := range rules {
		rule.Compile()
	}
	e.rules = rules
}

// ClearRules removes all rules from the engine.
func (e *Engine) ClearRules() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.rules = make([]*Rule, 0)
}

// GetRules returns all rules registered in the engine
func (e *Engine) GetRules() []*Rule {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.rules
}

// SetEventHandler sets the global event handler for the engine
func (e *Engine) SetEventHandler(handler EventHandler) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.eventHandler = handler
}

// RegisterEvents registers a named handler that can be referenced by rules.
// Handlers are invoked when rules succeed.
func (e *Engine) RegisterEvents(events ...Event) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.events == nil {
		e.events = make(map[string]Event)
	}

	for _, event := range events {
		e.events[event.Name] = event
	}
}

// RegisterEvent registers a named handler that can be referenced by rules.
// Handlers are invoked when rules succeed.
func (e *Engine) RegisterEvent(event Event) {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.events == nil {
		e.events = make(map[string]Event)
	}

	e.events[event.Name] = event
}

// Run executes all rules in the engine against the provided almanac.
// Rules are evaluated in priority order (higher priority first).
// Returns a slice of RuleResults containing the outcome of each rule evaluation.
// If any error occurs during evaluation, execution stops and the error is returned.
// Run executes all registered rules against the facts in the provided Almanac.
// It returns a map of rule names to their evaluation results.
func (e *Engine) Run(almanac *Almanac) (*Engine, error) {
	startTime := time.Now()

	// Snapshot rules and options to ensure thread-safety during execution
	e.mu.RLock()
	rules := make([]*Rule, len(e.rules))
	copy(rules, e.rules)

	options := make(map[string]interface{})
	for k, v := range e.options {
		options[k] = v
	}
	metrics := e.metrics
	e.mu.RUnlock()

	defer func() {
		if metrics != nil {
			metrics.ObserveEngineRun(len(rules), time.Since(startTime))
		}
	}()

	// Apply engine options to the almanac if needed
	if enabled, ok := options[EngineOptionKeyCacheConditions].(bool); ok && enabled {
		WithAlmanacConditionCaching()(almanac)
	}

	// Check for parallel execution
	if parallel, _ := options[EngineOptionKeyParallel].(bool); parallel {
		return e.runParallel(almanac, rules, options)
	}

	// Sort rules by priority if configured
	e.sortRulesByPriority(rules, options)

	var results = make(map[string]*RuleResult)

	// Evaluate each rule in priority order
	for _, rule := range rules {
		// Check for smart skip if enabled
		if skip, ok := options[EngineOptionKeySmartSkip].(bool); ok && skip {
			requiredFacts := rule.GetRequiredFacts()
			almanacFacts := almanac.GetFacts()
			missingFact := false
			for _, factID := range requiredFacts {
				if _, exists := almanacFacts[factID]; !exists {
					missingFact = true
					break
				}
			}
			if missingFact {
				results[rule.Name] = &RuleResult{
					Name:     rule.Name,
					Priority: rule.Priority,
					Result:   false,
				}
				if metrics != nil {
					metrics.ObserveRuleEvaluation(rule.Name, false, 0)
				}
				continue
			}
		}

		evalStart := time.Now()
		// Evaluate rule conditions
		condRes, err := rule.Conditions.Evaluate(almanac)
		evalDuration := time.Since(evalStart)

		if metrics != nil {
			metrics.ObserveRuleEvaluation(rule.Name, condRes.Result, evalDuration)
		}

		if err != nil {
			e.results = results
			return e, &RuleEngineError{
				Type: ErrEngine,
				Msg:  fmt.Sprintf("Error evaluating rule '%s': %v", rule.Name, err),
				Err:  err,
			}
		}

		ruleResult := &RuleResult{
			Name:      rule.Name,
			Priority:  rule.Priority,
			Result:    condRes.Result,
			OnSuccess: rule.OnSuccess,
			OnFailure: rule.OnFailure,
		}

		// Add audit trace if enabled
		if enabled, ok := options[EngineOptionKeyAuditTrace].(bool); ok && enabled {
			ruleResult.Conditions = condRes
		}

		results[rule.Name] = ruleResult

		if condRes.Result {
			// 1. Call OnSuccess event handlers
			if rule.OnSuccess != nil {
				for _, event := range rule.OnSuccess {
					err = e.HandleEvent(event.Name, rule.Name, condRes.Result, almanac, event.Params)
					if err != nil {
						e.mu.Lock()
						e.results = results
						e.mu.Unlock()
						return e, err
					}
				}
			}
		} else {
			// 1. Call OnFailure event handlers
			if rule.OnFailure != nil {
				for _, event := range rule.OnFailure {
					err = e.HandleEvent(event.Name, rule.Name, condRes.Result, almanac, event.Params)
					if err != nil {
						e.mu.Lock()
						e.results = results
						e.mu.Unlock()
						return e, err
					}
				}
			}
		}
	}

	e.mu.Lock()
	e.results = results
	e.almanac = almanac
	e.mu.Unlock()

	return e, nil
}

// runParallel executes rules in parallel using a worker pool.
func (e *Engine) runParallel(almanac *Almanac, rules []*Rule, options map[string]interface{}) (*Engine, error) {
	workerCount, ok := options[EngineOptionKeyWorkerCount].(int)
	if !ok || workerCount <= 0 {
		workerCount = 1
	}

	// Sort rules by priority if configured (important for event execution order)
	e.sortRulesByPriority(rules, options)

	numRules := len(rules)
	e.mu.RLock()
	metrics := e.metrics
	e.mu.RUnlock()

	resultsChan := make(chan struct {
		index    int
		res      *ConditionSetResult
		err      error
		duration time.Duration
	}, numRules)

	rulesChan := make(chan struct {
		index int
		rule  *Rule
	}, numRules)

	// 1. Start workers
	var wg sync.WaitGroup
	for i := 0; i < workerCount; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task := range rulesChan {
				start := time.Now()
				res, err := task.rule.Conditions.Evaluate(almanac)
				duration := time.Since(start)
				resultsChan <- struct {
					index    int
					res      *ConditionSetResult
					err      error
					duration time.Duration
				}{task.index, res, err, duration}
			}
		}()
	}

	// 2. Send rules (respecting smart skip)
	for i, rule := range rules {
		if skip, ok := options[EngineOptionKeySmartSkip].(bool); ok && skip {
			requiredFacts := rule.GetRequiredFacts()
			almanacFacts := almanac.GetFacts()
			missingFact := false
			for _, factID := range requiredFacts {
				if _, exists := almanacFacts[factID]; !exists {
					missingFact = true
					break
				}
			}
			if missingFact {
				resultsChan <- struct {
					index    int
					res      *ConditionSetResult
					err      error
					duration time.Duration
				}{i, &ConditionSetResult{Result: false}, nil, 0}
				continue
			}
		}
		rulesChan <- struct {
			index int
			rule  *Rule
		}{i, rule}
	}
	close(rulesChan)

	// Wait for workers in background and close results channel
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	// 3. Collect results
	orderedResults := make([]*ConditionSetResult, numRules)
	var firstErr error
	for r := range resultsChan {
		if r.err != nil && firstErr == nil {
			firstErr = r.err
		}
		orderedResults[r.index] = r.res
		if metrics != nil && r.res != nil {
			metrics.ObserveRuleEvaluation(rules[r.index].Name, r.res.Result, r.duration)
		}
	}

	if firstErr != nil {
		return e, &RuleEngineError{
			Type: ErrEngine,
			Msg:  fmt.Sprintf("Parallel execution error: %v", firstErr),
			Err:  firstErr,
		}
	}

	// 4. Sequential event triggering (important for predictability)
	results := make(map[string]*RuleResult)
	for i, rule := range rules {
		condRes := orderedResults[i]

		ruleResult := &RuleResult{
			Name:      rule.Name,
			Priority:  rule.Priority,
			Result:    condRes.Result,
			OnSuccess: rule.OnSuccess,
			OnFailure: rule.OnFailure,
		}

		// Add audit trace if enabled
		if enabled, ok := options[EngineOptionKeyAuditTrace].(bool); ok && enabled {
			ruleResult.Conditions = condRes
		}

		results[rule.Name] = ruleResult

		if condRes.Result {
			if rule.OnSuccess != nil {
				for _, event := range rule.OnSuccess {
					err := e.HandleEvent(event.Name, rule.Name, condRes.Result, almanac, event.Params)
					if err != nil {
						e.mu.Lock()
						e.results = results
						e.mu.Unlock()
						return e, err
					}
				}
			}
		} else {
			if rule.OnFailure != nil {
				for _, event := range rule.OnFailure {
					err := e.HandleEvent(event.Name, rule.Name, condRes.Result, almanac, event.Params)
					if err != nil {
						e.mu.Lock()
						e.results = results
						e.mu.Unlock()
						return e, err
					}
				}
			}
		}
	}

	e.mu.Lock()
	e.results = results
	e.almanac = almanac
	e.mu.Unlock()
	return e, nil
}

// Results returns the detailed results of rule evaluations.
func (e *Engine) Results() map[string]*RuleResult {
	e.mu.RLock()
	defer e.mu.RUnlock()
	return e.results
}

// ReduceResults converts a map of detailed RuleResults to a simple map of booleans.
func (e *Engine) ReduceResults() map[string]bool {
	e.mu.RLock()
	defer e.mu.RUnlock()

	reduced := make(map[string]bool, len(e.results))
	for name, res := range e.results {
		reduced[name] = res.Result
	}
	return reduced
}

// GenerateResponse builds the formatted Response object for JSON Marshalling.
func (e *Engine) GenerateResponse() *EngineResponse {
	e.mu.RLock()
	defer e.mu.RUnlock()

	res := &EngineResponse{
		Decision: DecisionDecline,
		Events:   []EventResponse{},
		Metadata: make(map[string]interface{}),
	}

	if len(e.results) == 0 {
		return res
	}

	// Extract metadata from Almanac if available
	if e.almanac != nil {
		facts := e.almanac.GetFacts()
		for factID, fact := range facts {
			meta := fact.Metadata()
			if len(meta) > 0 {
				res.Metadata[string(factID)] = meta
			}
		}
	}

	var primaryResult *RuleResult

	// Sort rules by priority if possible to determine the primary result
	// Note: e.rules is already what we use in Run, so we use its order or max priority
	for _, result := range e.results {
		if primaryResult == nil || result.Priority > primaryResult.Priority {
			primaryResult = result
		}

		if result.Result {
			res.Decision = DecisionAuthorize
			for _, ev := range result.OnSuccess {
				res.Events = append(res.Events, EventResponse{
					Type:   ev.Name,
					Params: ev.Params,
				})
			}
		} else {
			for _, ev := range result.OnFailure {
				res.Events = append(res.Events, EventResponse{
					Type:   ev.Name,
					Params: ev.Params,
				})
			}
		}
	}

	if primaryResult != nil {
		if primaryResult.Conditions != nil {
			res.Reason = primaryResult.Conditions
		} else {
			res.Reason = fmt.Sprintf("Rule '%s' determined the result", primaryResult.Name)
		}
	}

	return res
}

// HandleEvent invokes the event handler for the given event with context.
// Supports both synchronous and asynchronous execution based on event mode.
// ruleParams are optional parameters passed from the rule itself and combined with event defaults.
func (e *Engine) HandleEvent(eventName string, ruleName string, result bool, almanac *Almanac, ruleParams map[string]interface{}) error {
	e.mu.RLock()
	event, exists := e.events[eventName]
	handlerHost := e.eventHandler
	metrics := e.metrics
	e.mu.RUnlock()

	if !exists {
		if handlerHost != nil {
			return &RuleEngineError{
				Type: ErrEngine,
				Msg:  fmt.Sprintf("Event '%s' not registered", eventName),
				Err:  nil,
			}
		}
		return nil
	}

	// Merge parameters: rule params take precedence over event defaults
	finalParams := make(map[string]interface{})
	for k, v := range event.Params {
		finalParams[k] = v
	}
	for k, v := range ruleParams {
		finalParams[k] = v
	}

	// Build event context
	ctx := EventContext{
		RuleName:  ruleName,
		Result:    result,
		Almanac:   almanac,
		Timestamp: time.Now(),
		Params:    finalParams,
	}

	// Handle async events
	if event.Mode == EventModeAsync {
		go func() {
			start := time.Now()
			// Execute action if defined
			if event.Action != nil {
				_ = event.Action(ctx)
			}
			// Execute global handler if defined
			if handlerHost != nil {
				_ = handlerHost.Handle(event, ctx)
			}

			if metrics != nil {
				metrics.ObserveEventExecution(eventName, ruleName, result, time.Since(start))
			}
		}()
		return nil
	}

	// Handle sync events
	var err error
	start := time.Now()

	// Execute action if defined
	if event.Action != nil {
		if err = event.Action(ctx); err != nil {
			return &RuleEngineError{
				Type: ErrEngine,
				Msg:  fmt.Sprintf("Error executing action for event '%s': %v", eventName, err),
				Err:  err,
			}
		}
	}

	// Execute global handler if defined
	if handlerHost != nil {
		if err = handlerHost.Handle(event, ctx); err != nil {
			return &RuleEngineError{
				Type: ErrEngine,
				Msg:  fmt.Sprintf("Error in Event %s : \n %v", eventName, err),
				Err:  err,
			}
		}
	}

	if metrics != nil {
		metrics.ObserveEventExecution(eventName, ruleName, result, time.Since(start))
	}

	return nil
}

// sortRulesByPriority sorts the engine's rules by their priority in descending order.
func (e *Engine) sortRulesByPriority(rules []*Rule, options map[string]interface{}) {
	if options[SortByPriority] != nil {
		// Sort by priority
		sort.SliceStable(rules, func(i, j int) bool {
			if options[SortByPriority] == SortRuleASC {
				return rules[i].Priority < rules[j].Priority
			}
			return rules[i].Priority > rules[j].Priority
		})
	}
}

package gorulesengine

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/oliveagle/jsonpath"
)

// AlmanacOptionKeyAllowUndefinedFacts is the option key for allowing undefined facts.
const AlmanacOptionKeyAllowUndefinedFacts = "allowUndefinedFacts"

// AlmanacOptionKeyCacheConditions is the option key for enabling condition results caching.
const AlmanacOptionKeyCacheConditions = "cacheConditions"

// EventOutcome represents an event outcome triggered by a rule.
type EventOutcome string

// EventOutcomeSuccess represents a successful event outcome.
const EventOutcomeSuccess EventOutcome = "success"

// EventOutcomeFailure represents a failed event outcome.
const EventOutcomeFailure EventOutcome = "failure"

// Almanac stores facts and their computed values during rule evaluation.
// It maintains a cache for fact values, tracks events, and manages rule results.
// Almanac is thread-safe for concurrent access.
type Almanac struct {
	facts                 map[FactID]*Fact
	factResultsCache      map[string]interface{}
	conditionResultsCache map[string]interface{}
	pathResolver          PathResolver
	options               map[string]interface{}
	mutex                 sync.RWMutex
}

// AlmanacOption defines a functional option for configuring an Almanac.
type AlmanacOption func(*Almanac)

// PathResolver resolves nested values within facts using a path expression (e.g., JSONPath).
type PathResolver func(value interface{}, path string) (interface{}, error)

// DefaultPathResolver implements JSONPath resolution for accessing nested fact values.
// Example: "$.user.profile.age" accesses deeply nested data.
func DefaultPathResolver(value interface{}, path string) (interface{}, error) {
	if path == "" {
		return value, nil
	}
	return jsonpath.JsonPathLookup(value, path)
}

// AllowUndefinedFacts configures the almanac to return nil instead of errors for undefined facts.
// This is useful when you want to gracefully handle missing data.
func AllowUndefinedFacts() AlmanacOption {
	return func(a *Almanac) {
		a.options[AlmanacOptionKeyAllowUndefinedFacts] = true
	}
}

// WithAlmanacConditionCaching enables caching of condition results.
func WithAlmanacConditionCaching() AlmanacOption {
	return func(a *Almanac) {
		a.options[AlmanacOptionKeyCacheConditions] = true
	}
}

// NewAlmanac creates a new Almanac instance with the provided facts and options.
// The almanac is initialized with default settings including undefined fact handling.
//
// Example:
//
//	almanac := gre.NewAlmanac([]*gre.Fact{})
//	almanac.AddFact("age", 25)
//	almanac.AddFact("country", "FR")
func NewAlmanac(opts ...AlmanacOption) *Almanac {
	a := &Almanac{
		facts:                 make(map[FactID]*Fact),
		factResultsCache:      make(map[string]interface{}),
		conditionResultsCache: make(map[string]interface{}),
		pathResolver:          DefaultPathResolver,
		options:               make(map[string]interface{}),
	}

	AllowUndefinedFacts()(a)

	for _, opt := range opts {
		if opt != nil {
			opt(a)
		}
	}

	return a
}

// AddFact adds a fact to the almanac.
// The valueOrMethod can be either a static value or a function for dynamic facts.
// Optional FactOptions can be provided to configure caching and priority.
//
// Example:
//
//	// Static fact
//	almanac.AddFact("age", 25)
//
//	// Dynamic fact
//	almanac.AddFact("temperature", func(params map[string]interface{}) interface{} {
//	    return fetchTemperature()
//	})
func (a *Almanac) AddFact(id FactID, valueOrMethod interface{}, opts ...FactOption) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	fact := NewFact(id, valueOrMethod, opts...)
	a.facts[id] = &fact

	a.PreCacheFactValue(&fact)

	return nil
}

// AddFacts adds multiple facts to the almanac.
func (a *Almanac) AddFacts(facts ...*Fact) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for _, fact := range facts {
		a.facts[fact.ID()] = fact
		a.PreCacheFactValue(fact)
	}
}

// PreCacheFactValue computes and caches the value of a fact if caching is enabled.
func (a *Almanac) PreCacheFactValue(fact *Fact) {
	// Pre-cache the static fact value if caching is enabled
	if cacheEnabled, ok := fact.options[FactOptionKeyCache].(bool); ok && cacheEnabled {
		if !fact.IsDynamic() {
			cacheKey, _ := fact.GetCacheKey()
			a.factResultsCache[cacheKey] = fact.ValueOrMethod()
		}
	}
}

// GetFactValue retrieves the value of a fact by its ID.
// For dynamic facts, params can be passed to the computation function.
// The path parameter allows accessing nested values using JSONPath.
//
// Example:
//
//	// Simple fact access
//	age, _ := almanac.GetFactValue("age", nil, "")
//
//	// Nested access with JSONPath
//	city, _ := almanac.GetFactValue("user", nil, "$.address.city")
func (a *Almanac) GetFactValue(factID FactID, params map[string]interface{}, path string) (interface{}, error) {
	var fact *Fact
	var exists bool
	var cachedVal interface{}
	var cached bool

	// Read lock for concurrent access
	a.mutex.RLock()
	fact, exists = a.facts[factID]
	a.mutex.RUnlock()

	// Fact not found
	if !exists {
		// Check if undefined facts are allowed
		if allowUndefined, ok := a.options[AlmanacOptionKeyAllowUndefinedFacts].(bool); ok && allowUndefined {
			return nil, nil
		}
		return nil, &AlmanacError{
			Payload: "factID=" + string(factID),
			Err:     fmt.Errorf("fact '%s' is not defined in the almanac", factID),
		}
	}

	// Check cache first
	if val, _ := fact.GetOption(FactOptionKeyCache); val == true {
		cachedVal, cached = a.GetFactValueFromCache(factID)
	}

	var val interface{}

	// If cached value exists, use it
	if cached {
		val = cachedVal
	} else {
		// Calculate fact value
		val = fact.Calculate(params)

		// Cache the result if caching is enabled
		if cacheEnabled, ok := fact.options[FactOptionKeyCache].(bool); ok && cacheEnabled {
			// Generate cache key for storing the result
			key, _ := fact.GetCacheKey()
			a.mutex.Lock()
			a.factResultsCache[key] = val
			a.mutex.Unlock()
		}
	}

	// Apply path resolution if path is provided
	val, err := a.TraversePath(val, path)
	if err != nil {
		return nil, &AlmanacError{
			Payload: fmt.Sprintf("factID=%s, path=%s", factID, path),
			Err:     fmt.Errorf("failed to resolve path '%s' for fact '%s': %v", path, factID, err),
		}
	}

	return val, nil
}

// GetFactValueFromCache retrieves a fact value directly from the cache
func (a *Almanac) GetFactValueFromCache(factID FactID) (interface{}, bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	fact := a.facts[factID]

	cacheKey, _ := fact.GetCacheKey()

	cachedVal, cached := a.factResultsCache[cacheKey]

	return cachedVal, cached
}

// GetConditionResultFromCache retrieves a condition result from the cache.
func (a *Almanac) GetConditionResultFromCache(key string) (interface{}, bool) {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	result, cached := a.conditionResultsCache[key]
	return result, cached
}

// SetConditionResultCache stores a condition result in the cache.
func (a *Almanac) SetConditionResultCache(key string, result interface{}) {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.conditionResultsCache[key] = result
}

// IsConditionCachingEnabled checks if condition caching is enabled in the almanac.
func (a *Almanac) IsConditionCachingEnabled() bool {
	enabled, ok := a.options[AlmanacOptionKeyCacheConditions].(bool)
	return ok && enabled
}

// TraversePath is a helper to traverse nested structures based on a path expression.
// It uses the configured PathResolver to access nested values within complex data structures.
func (a *Almanac) TraversePath(data interface{}, path string) (interface{}, error) {
	var val = data
	// Apply path resolution if path is provided
	if path != "" {
		// Check if value is a complex type (map, slice, struct) that supports path resolution
		valType := reflect.TypeOf(data)
		if valType != nil {
			kind := valType.Kind()
			// Only apply path resolver to complex types
			if kind == reflect.Map || kind == reflect.Slice || kind == reflect.Struct || kind == reflect.Ptr {
				return a.pathResolver(data, path)
			}
		}
	}
	// For primitive types (string, int, bool, etc.), path resolution doesn't make sense
	// Return the value as-is
	return val, nil
}

// GetOptions returns the almanac options
func (a *Almanac) GetOptions() map[string]interface{} {
	return a.options
}

// GetFacts returns the almanac's fact map
func (a *Almanac) GetFacts() map[FactID]*Fact {
	return a.facts
}

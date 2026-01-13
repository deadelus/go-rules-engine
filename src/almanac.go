package gorulesengine

import (
	"fmt"
	"reflect"
	"sync"

	"github.com/oliveagle/jsonpath"
)

const ALMANAC_OPTION_KEY_ALLOW_UNDEFINED_FACTS = "allowUndefinedFacts"

type Almanac struct {
	factMap          map[FactID]*Fact
	factResultsCache map[string]interface{}
	events           *EventResults
	pathResolver     PathResolver
	options          map[string]interface{}
	mutex            sync.RWMutex
}

type EventResults struct {
	success []interface{}
	fail    []interface{}
}

type AlmanacOption func(*Almanac)

// PathResolver is the equivalent of the JavaScript pathResolver
type PathResolver func(value interface{}, path string) (interface{}, error)

// DefaultPathResolver uses JSONPath
func DefaultPathResolver(value interface{}, path string) (interface{}, error) {
	if path == "" {
		return value, nil
	}
	return jsonpath.JsonPathLookup(value, path)
}

// AllowUndefinedFacts configures the almanac to allow undefined facts
func AllowUndefinedFacts() AlmanacOption {
	return func(a *Almanac) {
		a.options[ALMANAC_OPTION_KEY_ALLOW_UNDEFINED_FACTS] = true
	}
}

// NewAlmanac creates a new almanac with provided facts and options
func NewAlmanac(facts []*Fact, opts ...AlmanacOption) *Almanac {
	a := &Almanac{
		factMap:          make(map[FactID]*Fact),
		factResultsCache: make(map[string]interface{}),
		events:           &EventResults{},
		pathResolver:     DefaultPathResolver,
		options:          make(map[string]interface{}),
	}

	AllowUndefinedFacts()(a)

	for _, opt := range opts {
		opt(a)
	}

	// Add provided facts to the fact map
	for _, fact := range facts {
		a.factMap[fact.ID()] = fact
	}
	return a
}

// AddFact convert payload to a Fact and adds it to the almanac
func (a *Almanac) AddFact(id FactID, valueOrMethod interface{}, opts ...FactOption) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	fact := NewFact(id, valueOrMethod, opts...)
	a.factMap[id] = fact

	// Pre-cache the static fact value if caching is enabled
	if cacheEnabled, ok := fact.options[FACT_OPTION_KEY_CACHE].(bool); ok && cacheEnabled {
		if !fact.IsDynamic() {
			cacheKey, _ := fact.GetCacheKey()
			a.factResultsCache[cacheKey] = fact.ValueOrMethod()
		}
	}

	return nil
}

// GetFactValue retrieves a fact value with optional path extraction
func (a *Almanac) GetFactValue(factID FactID, params map[string]interface{}, path string) (interface{}, error) {
	var fact *Fact
	var exists bool
	var cachedVal interface{}
	var cached bool

	// Read lock for concurrent access
	a.mutex.RLock()
	fact, exists = a.factMap[factID]
	a.mutex.RUnlock()

	// Fact not found
	if !exists {
		// Check if undefined facts are allowed
		if allowUndefined, ok := a.options[ALMANAC_OPTION_KEY_ALLOW_UNDEFINED_FACTS].(bool); ok && allowUndefined {
			return nil, nil
		}
		return nil, &AlmanacError{
			Payload: "factID=" + string(factID),
			Err:     fmt.Errorf("fact '%s' is not defined in the almanac", factID),
		}
	}

	// Check cache first
	if val, _ := fact.GetOption(FACT_OPTION_KEY_CACHE); val == true {
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
		if cacheEnabled, ok := fact.options[FACT_OPTION_KEY_CACHE].(bool); ok && cacheEnabled {
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

	fact := a.factMap[factID]

	cacheKey, _ := fact.GetCacheKey()

	cachedVal, cached := a.factResultsCache[cacheKey]

	return cachedVal, cached
}

// traversePath is a helper to traverse nested structures based on a path
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
	return a.factMap
}

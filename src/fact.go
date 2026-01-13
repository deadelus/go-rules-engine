package gorulesengine

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
)

const DYNAMIC_FACT = "__dynamic_fact__"

const CONSTANT_FACT = "__constant_fact__"

const FACT_OPTION_KEY_CACHE = "cache"

const FACT_OPTION_KEY_PRIORITY = "priority"

type FactID string

type Fact struct {
	id            FactID
	valueOrMethod interface{}
	factType      string
	options       map[string]interface{}
}

// FactOption defines a functional option for the fact
type FactOption func(*Fact)

// WithCache enables caching for the fact
func WithCache() FactOption {
	return func(f *Fact) {
		f.options[FACT_OPTION_KEY_CACHE] = true
	}
}

// WithoutCache disables caching for the fact
func WithoutCache() FactOption {
	return func(f *Fact) {
		f.options[FACT_OPTION_KEY_CACHE] = false
	}
}

// WithPriority sets the priority of the fact
func WithPriority(priority int) FactOption {
	return func(f *Fact) {
		f.options[FACT_OPTION_KEY_PRIORITY] = priority
	}
}

// NewFact creates a new fact with functional options
func NewFact(id FactID, valueOrMethod interface{}, opts ...FactOption) *Fact {
	fact := Fact{
		id:            id,
		valueOrMethod: valueOrMethod,
		options:       map[string]interface{}{},
		factType: func() string {
			// Use reflect to detect any function type
			if reflect.TypeOf(valueOrMethod).Kind() == reflect.Func {
				return DYNAMIC_FACT
			}
			return CONSTANT_FACT
		}(),
	}

	// Default priority is 0
	WithPriority(0)(&fact)

	if fact.factType == DYNAMIC_FACT {
		WithCache()(&fact)
	}

	// Apply functional options
	for _, opt := range opts {
		opt(&fact)
	}

	return &fact
}

// ID returns the fact identifier
func (f *Fact) ID() FactID {
	return f.id
}

// ValueOrMethod returns the fact value or method
func (f *Fact) ValueOrMethod() interface{} {
	return f.valueOrMethod
}

// FactType returns the fact type
func (f *Fact) FactType() string {
	return f.factType
}

// IsDynamic returns true if the fact is dynamic
func (f *Fact) IsDynamic() bool {
	return f.factType == DYNAMIC_FACT
}

// GetOption returns an option in an immutable way
func (f *Fact) GetOption(key string) (interface{}, bool) {
	val, exists := f.options[key]
	return val, exists
}

// HasOption checks if an option exists
func (f *Fact) HasOption(key string) bool {
	_, exists := f.options[key]
	return exists
}

// GetCacheKey generates a unique cache key for the fact if it's cached
func (f *Fact) GetCacheKey() (string, error) {
	if f.options[FACT_OPTION_KEY_CACHE] == true {
		// Use the fact ID as cache key for both static and dynamic facts
		// This simplifies the caching mechanism
		return f.hashFromID()
	}
	return "", nil
}

// Calculate executes the dynamic fact method or returns the constant fact value
func (f *Fact) Calculate(params map[string]interface{}) interface{} {
	method := reflect.ValueOf(f.valueOrMethod)
	methodType := method.Type()

	// If it's not a function, return the value directly
	methodKind := methodType.Kind()
	if methodKind != reflect.Func {
		return f.valueOrMethod
	}

	var results []reflect.Value

	// Handle different method signatures
	switch methodType.NumIn() {
	case 0:
		// Method with no parameters
		results = method.Call([]reflect.Value{})
	case 1:
		// Method with one parameter (params)
		results = method.Call([]reflect.Value{reflect.ValueOf(params)})
	default:
		// Unsupported signature
		return nil
	}

	return results[0].Interface()
}

// hashFromID generates a unique MD5 hash based on the fact ID
func (f *Fact) hashFromID() (string, error) {
	// Use the fact ID to generate the cache key
	// This is simpler and works for both static and dynamic facts
	bytes := []byte(f.id)
	sum := md5.Sum(bytes)
	return hex.EncodeToString(sum[:]), nil
}

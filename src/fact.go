package gorulesengine

import (
	"crypto/md5"
	"encoding/hex"
	"reflect"
)

// DynamicFact identifies a fact that computes its value dynamically using a function.
const DynamicFact = "__dynamic_fact__"

// ConstantFact identifies a fact with a static, pre-defined value.
const ConstantFact = "__constant_fact__"

// FactOptionKeyCache is the key for the caching option in fact options.
const FactOptionKeyCache = "cache"

// FactOptionKeyPriority is the key for the priority option in fact options.
const FactOptionKeyPriority = "priority"

// FactID is a unique identifier for a fact.
type FactID string

// Fact represents a piece of data (fact) that can be used in rule conditions.
// Facts can be static values or dynamic functions that compute values on demand.
//
// Example:
//
//	// Static fact
//	fact := gre.NewFact("age", 25)
//
//	// Dynamic fact
//	fact := gre.NewFact("temperature", func(params map[string]interface{}) interface{} {
//	    return fetchTemperatureFromAPI()
//	})
type Fact struct {
	id            FactID
	valueOrMethod interface{}
	factType      string
	options       map[string]interface{}
	metadata      map[string]interface{}
}

// FactOption defines a functional option for configuring facts.
type FactOption func(*Fact)

// WithCache enables caching for dynamic facts.
// When enabled, the fact's value will be computed once and reused.
func WithCache() FactOption {
	return func(f *Fact) {
		f.options[FactOptionKeyCache] = true
	}
}

// WithoutCache disables caching for facts.
// When disabled, dynamic facts will be re-evaluated on each access.
func WithoutCache() FactOption {
	return func(f *Fact) {
		f.options[FactOptionKeyCache] = false
	}
}

// WithPriority sets the evaluation priority of the fact.
// Higher priority facts may be evaluated before lower priority facts.
func WithPriority(priority int) FactOption {
	return func(f *Fact) {
		f.options[FactOptionKeyPriority] = priority
	}
}

// WithMetadata adds metadata to the fact.
func WithMetadata(metadata map[string]interface{}) FactOption {
	return func(f *Fact) {
		if f.metadata == nil {
			f.metadata = make(map[string]interface{})
		}
		for k, v := range metadata {
			f.metadata[k] = v
		}
	}
}

// NewFact creates a new fact with the given ID and value or computation function.
// If valueOrMethod is a function, the fact is dynamic and will compute its value on demand.
// Otherwise, the fact is static with a constant value.
//
// Options can be provided to customize caching and priority behavior.
//
// Example:
//
//	// Static fact
//	fact := gre.NewFact("age", 25)
//
//	// Dynamic fact with custom options
//	fact := gre.NewFact("temperature",
//	    func(params map[string]interface{}) interface{} {
//	        return fetchTemperature()
//	    },
//	    gre.WithCache(),
//	    gre.WithPriority(10),
//	)
func NewFact(id FactID, valueOrMethod interface{}, opts ...FactOption) Fact {
	fact := Fact{
		id:            id,
		valueOrMethod: valueOrMethod,
		options:       map[string]interface{}{},
		metadata:      map[string]interface{}{},
		factType: func() string {
			// Use reflect to detect any function type
			if reflect.TypeOf(valueOrMethod).Kind() == reflect.Func {
				return DynamicFact
			}
			return ConstantFact
		}(),
	}

	// Default priority is 0
	WithPriority(0)(&fact)

	if fact.factType == DynamicFact {
		WithCache()(&fact)
	}

	// Apply functional options
	for _, opt := range opts {
		opt(&fact)
	}

	return fact
}

// ID returns the unique identifier of the fact.
func (f *Fact) ID() FactID {
	return f.id
}

// Metadata returns the metadata associated with the fact.
func (f *Fact) Metadata() map[string]interface{} {
	return f.metadata
}

// ValueOrMethod returns the fact's value (for static facts) or computation function (for dynamic facts).
func (f *Fact) ValueOrMethod() interface{} {
	return f.valueOrMethod
}

// FactType returns the type of the fact: DYNAMIC_FACT or CONSTANT_FACT.
func (f *Fact) FactType() string {
	return f.factType
}

// IsDynamic returns true if the fact computes its value dynamically using a function.
func (f *Fact) IsDynamic() bool {
	return f.factType == DynamicFact
}

// GetOption returns the value of a specific option and whether it exists.
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
	if f.options[FactOptionKeyCache] == true {
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

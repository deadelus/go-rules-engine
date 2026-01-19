package gorulesengine_test

import (
	"sync"
	"testing"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func TestAddCachedFact(t *testing.T) {
	opts := []gre.AlmanacOption{}
	almanac := gre.NewAlmanac(opts...)

	factValue := "cached_value"
	almanac.AddFact("cached_fact", factValue, gre.WithCache())

	retrievedFact, err := almanac.GetFactValue("cached_fact", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrievedFact != "cached_value" {
		t.Fatalf("Expected fact value 'cached_value', got %v", retrievedFact)
	}
}

func TestGetFactValue_FromCache(t *testing.T) {
	almanac := gre.NewAlmanac()

	// Add a static fact with cache enabled
	// Pre-caching should be performed in AddFact
	err := almanac.AddFact("cached_static", "static_value", gre.WithCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// First retrieval - should come from the pre-filled cache
	val1, err := almanac.GetFactValue("cached_static", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val1 != "static_value" {
		t.Errorf("Expected 'static_value', got %v", val1)
	}

	// Second retrieval - should also come from the cache
	val2, err := almanac.GetFactValue("cached_static", nil, "")
	if err != nil {
		t.Fatalf("Expected no error on second retrieval, got %v", err)
	}

	if val2 != "static_value" {
		t.Errorf("Expected 'static_value', got %v", val2)
	}
}

func TestGetFactValue_CalculateAndCache(t *testing.T) {
	almanac := gre.NewAlmanac()

	callCount := 0
	// Add a dynamic fact that counts calls
	dynamicFact := func(params map[string]interface{}) int {
		callCount++
		return 100 + callCount
	}

	// Add with cache disabled to see if it calculates every time
	err := almanac.AddFact("dynamic_no_cache", dynamicFact, gre.WithoutCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Retrieve twice - should calculate every time (no cache)
	val1, err := almanac.GetFactValue("dynamic_no_cache", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	val2, err := almanac.GetFactValue("dynamic_no_cache", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Values should be different because callCount increases
	if val1 == val2 {
		t.Errorf("Expected different values without cache, got %v and %v", val1, val2)
	}

	if callCount != 2 {
		t.Errorf("Expected 2 calls without cache, got %d", callCount)
	}
}

func TestGetFactValue_DynamicWithCaching(t *testing.T) {
	almanac := gre.NewAlmanac()

	callCount := 0
	// Create a dynamic fact with cache enabled
	// Now that GetCacheKey uses the ID for dynamic facts,
	// the cache should work
	dynamicFact := func() int {
		callCount++
		return 42
	}

	// Add with cache enabled
	err := almanac.AddFact("dynamic_cached", dynamicFact, gre.WithCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// First retrieval - should calculate and cache
	val1, err := almanac.GetFactValue("dynamic_cached", nil, "")
	if err != nil {
		t.Fatalf("Expected no error on first call, got %v", err)
	}

	if val1 != 42 {
		t.Errorf("Expected 42, got %v", val1)
	}

	// Second retrieval - should come from the cache (no recalculation)
	val2, err := almanac.GetFactValue("dynamic_cached", nil, "")
	if err != nil {
		t.Fatalf("Expected no error on second call, got %v", err)
	}

	if val2 != 42 {
		t.Errorf("Expected 42, got %v", val2)
	}

	// The fact should only have been calculated once (then cached)
	if callCount != 1 {
		t.Errorf("Expected 1 call (then cached), got %d", callCount)
	}
}

func TestGetFactValue_CacheAfterCalculation(t *testing.T) {
	almanac := gre.NewAlmanac()

	callCount := 0
	// Create a dynamic fact that returns a marshalable value
	// To allow the cache to work, we must return a simple value
	dynamicValue := 100
	dynamicFact := func() int {
		callCount++
		return dynamicValue
	}

	// Add WITHOUT cache to start
	err := almanac.AddFact("dynamic_fact", dynamicFact, gre.WithoutCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// First retrieval - should calculate
	val1, err := almanac.GetFactValue("dynamic_fact", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val1 != 100 {
		t.Errorf("Expected 100, got %v", val1)
	}

	// Second retrieval - should calculate again (no cache)
	val2, err := almanac.GetFactValue("dynamic_fact", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val2 != 100 {
		t.Errorf("Expected 100, got %v", val2)
	}

	// Should have been called 2 times
	if callCount != 2 {
		t.Errorf("Expected 2 calls without cache, got %d", callCount)
	}
}

func TestGetFactValue_ConcurrentCacheWrite(t *testing.T) {
	almanac := gre.NewAlmanac()

	callCount := 0
	var mutex sync.Mutex

	// Create a dynamic fact without pre-filled cache
	dynamicFact := func() int {
		mutex.Lock()
		callCount++
		count := callCount
		mutex.Unlock()
		return count
	}

	// Add without cache to force calculation every time
	err := almanac.AddFact("dynamic_concurrent", dynamicFact, gre.WithoutCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Start multiple goroutines that calculate and potentially cache
	const numGoroutines = 50
	results := make(chan interface{}, numGoroutines)
	errs := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			val, err := almanac.GetFactValue("dynamic_concurrent", nil, "")
			if err != nil {
				errs <- err
				return
			}
			results <- val
		}()
	}

	// Collect results
	seenValues := make(map[int]bool)
	for i := 0; i < numGoroutines; i++ {
		select {
		case err := <-errs:
			t.Fatalf("Unexpected error from goroutine: %v", err)
		case val := <-results:
			intVal, ok := val.(int)
			if !ok {
				t.Errorf("Expected int value, got %T", val)
			}
			seenValues[intVal] = true
		}
	}

	// Verify that all goroutines obtained a value
	if len(seenValues) == 0 {
		t.Error("Expected to see some values")
	}

	// The number of calls should be equal to the number of goroutines (no cache)
	mutex.Lock()
	finalCount := callCount
	mutex.Unlock()

	if finalCount != numGoroutines {
		t.Errorf("Expected %d calls without cache, got %d", numGoroutines, finalCount)
	}
}

func TestGetFactValue_DynamicCacheStorageAfterCalculation(t *testing.T) {
	almanac := gre.NewAlmanac()

	callCount := 0
	// Create a dynamic fact with cache enabled
	// With the modification of GetCacheKey, dynamic facts can now
	// have a cache key based on their ID
	dynamicFact := func() int {
		callCount++
		return 42
	}

	// Add the dynamic fact with cache enabled
	// Dynamic facts are NOT pre-cached in AddFact (condition !fact.IsDynamic())
	err := almanac.AddFact("dynamic_cacheable", dynamicFact, gre.WithCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// First retrieval - should calculate AND cache
	val1, err := almanac.GetFactValue("dynamic_cacheable", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val1 != 42 {
		t.Errorf("Expected 42, got %v", val1)
	}

	if callCount != 1 {
		t.Errorf("Expected 1 call after first retrieval, got %d", callCount)
	}

	// Second retrieval - should come from the cache (no new calculation)
	val2, err := almanac.GetFactValue("dynamic_cacheable", nil, "")
	if err != nil {
		t.Fatalf("Expected no error on second call, got %v", err)
	}

	if val2 != 42 {
		t.Errorf("Expected 42 from cache, got %v", val2)
	}

	// The fact should NOT have been recalculated (still 1 call)
	if callCount != 1 {
		t.Errorf("Expected 1 call total (second from cache), got %d", callCount)
	}
}

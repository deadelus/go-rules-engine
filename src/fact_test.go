package gorulesengine_test

import (
	"testing"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func TestFact_FactType(t *testing.T) {
	// Test static fact
	staticFact := gorulesengine.NewFact("static", "value")
	if staticFact.FactType() != gorulesengine.CONSTANT_FACT {
		t.Errorf("Expected CONSTANT_FACT, got %s", staticFact.FactType())
	}

	// Test dynamic fact
	dynamicFact := gorulesengine.NewFact("dynamic", func() int { return 42 })
	if dynamicFact.FactType() != gorulesengine.DYNAMIC_FACT {
		t.Errorf("Expected DYNAMIC_FACT, got %s", dynamicFact.FactType())
	}
}

func TestFact_HasOption(t *testing.T) {
	fact := gorulesengine.NewFact("test", "value", gorulesengine.WithCache(), gorulesengine.WithPriority(5))

	// Test existing option
	if !fact.HasOption(gorulesengine.FACT_OPTION_KEY_CACHE) {
		t.Error("Expected cache option to exist")
	}

	if !fact.HasOption(gorulesengine.FACT_OPTION_KEY_PRIORITY) {
		t.Error("Expected priority option to exist")
	}

	// Test non-existing option
	if fact.HasOption("nonexistent") {
		t.Error("Expected nonexistent option to not exist")
	}
}

func TestFact_GetCacheKey_NoCacheEnabled(t *testing.T) {
	fact := gorulesengine.NewFact("test", "value", gorulesengine.WithoutCache())

	cacheKey, err := fact.GetCacheKey()
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if cacheKey != "" {
		t.Errorf("Expected empty cache key when cache disabled, got '%s'", cacheKey)
	}
}

func TestFact_Calculate_UnsupportedSignature(t *testing.T) {
	// Create a function with unsupported signature (more than 1 parameter)
	unsupportedFunc := func(a int, b int) int { return a + b }

	fact := gorulesengine.NewFact("test", unsupportedFunc)

	result := fact.Calculate(nil)
	if result != nil {
		t.Errorf("Expected nil for unsupported signature, got %v", result)
	}
}

func TestFact_Calculate_WithParams(t *testing.T) {
	// Create a function that accepts params
	factFunc := func(params map[string]interface{}) int {
		if val, ok := params["multiplier"].(int); ok {
			return val * 10
		}
		return 0
	}

	fact := gorulesengine.NewFact("test", factFunc)

	params := map[string]interface{}{
		"multiplier": 5,
	}

	result := fact.Calculate(params)
	if result != 50 {
		t.Errorf("Expected 50, got %v", result)
	}
}

func TestFact_Calculate_NoParams(t *testing.T) {
	callCount := 0
	factFunc := func() int {
		callCount++
		return 42
	}

	fact := gorulesengine.NewFact("test", factFunc)

	result := fact.Calculate(nil)
	if result != 42 {
		t.Errorf("Expected 42, got %v", result)
	}

	if callCount != 1 {
		t.Errorf("Expected function to be called once, got %d", callCount)
	}
}

func TestFact_Calculate_StaticValue(t *testing.T) {
	fact := gorulesengine.NewFact("test", "static_value")

	result := fact.Calculate(nil)
	if result != "static_value" {
		t.Errorf("Expected 'static_value', got %v", result)
	}
}

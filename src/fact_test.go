package gorulesengine_test

import (
	"testing"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func TestFact_FactType(t *testing.T) {
	// Test static fact
	staticFact := gre.NewFact("static", "value")
	if staticFact.FactType() != gre.ConstantFact {
		t.Errorf("Expected CONSTANT_FACT, got %s", staticFact.FactType())
	}

	// Test dynamic fact
	dynamicFact := gre.NewFact("dynamic", func() int { return 42 })
	if dynamicFact.FactType() != gre.DynamicFact {
		t.Errorf("Expected DYNAMIC_FACT, got %s", dynamicFact.FactType())
	}
}

func TestFact_HasOption(t *testing.T) {
	fact := gre.NewFact("test", "value", gre.WithCache(), gre.WithPriority(5))

	// Test existing option
	if !fact.HasOption(gre.FactOptionKeyCache) {
		t.Error("Expected cache option to exist")
	}

	if !fact.HasOption(gre.FactOptionKeyPriority) {
		t.Error("Expected priority option to exist")
	}

	// Test non-existing option
	if fact.HasOption("nonexistent") {
		t.Error("Expected nonexistent option to not exist")
	}
}

func TestFact_GetCacheKey_NoCacheEnabled(t *testing.T) {
	fact := gre.NewFact("test", "value", gre.WithoutCache())

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

	fact := gre.NewFact("test", unsupportedFunc)

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

	fact := gre.NewFact("test", factFunc)

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

	fact := gre.NewFact("test", factFunc)

	result := fact.Calculate(nil)
	if result != 42 {
		t.Errorf("Expected 42, got %v", result)
	}

	if callCount != 1 {
		t.Errorf("Expected function to be called once, got %d", callCount)
	}
}

func TestFact_Calculate_StaticValue(t *testing.T) {
	fact := gre.NewFact("test", "static_value")

	result := fact.Calculate(nil)
	if result != "static_value" {
		t.Errorf("Expected 'static_value', got %v", result)
	}
}

func TestFact_Metadata(t *testing.T) {
	metadata := map[string]interface{}{
		"source": "database",
		"ttl":    3600,
	}
	fact := gre.NewFact("meta-test", "value", gre.WithMetadata(metadata))

	actualMeta := fact.Metadata()
	if len(actualMeta) != 2 {
		t.Errorf("Expected 2 metadata entries, got %d", len(actualMeta))
	}

	if actualMeta["source"] != "database" {
		t.Errorf("Expected source 'database', got %v", actualMeta["source"])
	}

	if actualMeta["ttl"] != 3600 {
		t.Errorf("Expected ttl 3600, got %v", actualMeta["ttl"])
	}

	// Test adding more metadata later (if we decide to support it, but WithMetadata currently merges)
	additional := map[string]interface{}{
		"verified": true,
	}
	// Note: Currently options are only applied in NewFact.
	// But let's verify merging works if passed twice to NewFact
	fact2 := gre.NewFact("meta-test-2", "value",
		gre.WithMetadata(metadata),
		gre.WithMetadata(additional),
	)

	actualMeta2 := fact2.Metadata()
	if len(actualMeta2) != 3 {
		t.Errorf("Expected 3 metadata entries after merging, got %d", len(actualMeta2))
	}
	if actualMeta2["verified"] != true {
		t.Error("Expected verified true in merged metadata")
	}
}

func TestFact_WithMetadata_Nil(t *testing.T) {
	// Directly call the option on a zero-value Fact to cover the nil check
	var fact gre.Fact
	opt := gre.WithMetadata(map[string]interface{}{"test": "nil-coverage"})
	opt(&fact)

	if fact.Metadata()["test"] != "nil-coverage" {
		t.Errorf("Expected 'nil-coverage', got %v", fact.Metadata()["test"])
	}
}

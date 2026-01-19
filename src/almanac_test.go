package gorulesengine_test

import (
	"testing"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func TestNewAlmanac(t *testing.T) {
	facts := []gre.Fact{
		gre.NewFact("fact1", "value"),
		gre.NewFact("fact2", "value"),
	}
	opts := []gre.AlmanacOption{
		gre.AllowUndefinedFacts(),
	}

	almanac := gre.NewAlmanac(opts...)

	almanac.AddFacts(&facts[0], &facts[1])

	if almanac == nil {
		t.Fatal("Expected almanac to be created, got nil")
	}

	expectedOpts := almanac.GetOptions()
	if allowUndefined, ok := expectedOpts[gre.AlmanacOptionKeyAllowUndefinedFacts]; !ok || allowUndefined != true {
		t.Fatalf("Expected allowUndefinedFacts to be true, got %v", allowUndefined)
	}

	retrievedFacts := almanac.GetFacts()
	if len(retrievedFacts) != 2 {
		t.Fatalf("Expected 2 facts, got %d", len(retrievedFacts))
	}

	fact1 := retrievedFacts["fact1"]
	if fact1.ID() != "fact1" {
		t.Fatalf("Expected fact ID 'fact1', got %v", fact1.ID())
	}

	fact2 := retrievedFacts["fact2"]
	if fact2.ID() != "fact2" {
		t.Fatalf("Expected fact ID 'fact2', got %v", fact2.ID())
	}
}

func TestAddFact(t *testing.T) {
	opts := []gre.AlmanacOption{}
	almanac := gre.NewAlmanac(opts...)

	factValue := map[string]interface{}{
		"secret": map[string]interface{}{
			"value": 42,
		},
	}
	almanac.AddFact("test_fact", factValue)

	retrievedFact, err := almanac.GetFactValue("test_fact", nil, "$.secret.value")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrievedFact != 42 {
		t.Fatalf("Expected fact value 42, got %v", retrievedFact)
	}

	// Test metadata through AddFact
	meta := map[string]interface{}{"origin": "test"}
	almanac.AddFact("meta_fact", "val", gre.WithMetadata(meta))

	facts := almanac.GetFacts()
	if fact, ok := facts["meta_fact"]; ok {
		if fact.Metadata()["origin"] != "test" {
			t.Errorf("Expected metadata origin 'test', got %v", fact.Metadata()["origin"])
		}
	} else {
		t.Error("Fact meta_fact not found")
	}
}

func TestAddFactCacheKeyError(t *testing.T) {
	almanac := gre.NewAlmanac()

	// With hashFromID, cache key generation can no longer fail
	// because we simply use the fact ID
	// This test now verifies that any value can be added with cache
	unmarshalableValue := make(chan int)

	// Adding a fact with a non-marshalable value should now work
	err := almanac.AddFact("channel_fact", unmarshalableValue, gre.WithCache())

	if err != nil {
		t.Fatalf("Expected no error with hashFromID, got: %v", err)
	}

	// Verify that the fact was successfully added
	facts := almanac.GetFacts()
	if _, exists := facts["channel_fact"]; !exists {
		t.Error("Expected fact to be added successfully")
	}
}

func TestGetFactValue_UndefinedFactNotAllowed(t *testing.T) {
	// Create a custom option to disable allowUndefinedFacts
	disallowUndefinedFacts := func(a *gre.Almanac) {
		opts := a.GetOptions()
		opts[gre.AlmanacOptionKeyAllowUndefinedFacts] = false
	}

	// Create an almanac with allowUndefinedFacts disabled
	almanac := gre.NewAlmanac(disallowUndefinedFacts)

	// Try to retrieve a fact that does not exist
	val, err := almanac.GetFactValue("nonexistent_fact", nil, "")

	// Should return an error
	if err == nil {
		t.Fatal("Expected error when getting undefined fact, got nil")
	}

	// Verify that it is an AlmanacError
	almanacErr, ok := err.(*gre.AlmanacError)
	if !ok {
		t.Fatalf("Expected *AlmanacError, got %T: %v", err, err)
	}

	// Verify the payload
	if almanacErr.Payload != "factID=nonexistent_fact" {
		t.Errorf("Expected payload 'factID=nonexistent_fact', got '%s'", almanacErr.Payload)
	}

	// Verify the error message
	if almanacErr.Err == nil {
		t.Fatal("Expected wrapped error, got nil")
	}

	// The returned value should be nil
	if val != nil {
		t.Errorf("Expected nil value, got %v", val)
	}
}

func TestGetFactValue_UndefinedFactAllowed(t *testing.T) {
	// Create an almanac with allowUndefinedFacts enabled (this is the default of NewAlmanac)
	almanac := gre.NewAlmanac(gre.AllowUndefinedFacts())

	// Try to retrieve a fact that does not exist
	val, err := almanac.GetFactValue("nonexistent_fact", nil, "")

	// Should NOT return an error
	if err != nil {
		t.Fatalf("Expected no error when allowUndefinedFacts is true, got: %v", err)
	}

	// The value should be nil
	if val != nil {
		t.Errorf("Expected nil value for undefined fact, got %v", val)
	}
}

package gorulesengine_test

import (
	"reflect"
	"testing"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func TestPathResolver(t *testing.T) {
	resolver := gre.DefaultPathResolver

	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
			"age":  30,
		},
	}

	value, err := resolver(data, "$.user.name")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if value != "Alice" {
		t.Errorf("Expected 'Alice', got %v", value)
	}

	value, err = resolver(data, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(data, value) {
		t.Fatalf("Expected original data, got %v", value)
	}
}

func TestGetFactValue_PrimitiveWithPath_String(t *testing.T) {
	almanac := gre.NewAlmanac()
	// Add a fact with a primitive string value
	err := almanac.AddFact("string_fact", "simple_string")
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Try to apply a path to a primitive string value
	val, err := almanac.GetFactValue("string_fact", nil, "$.some.path")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val != "simple_string" {
		t.Errorf("Expected 'simple_string', got %v", val)
	}
}

func TestGetFactValue_PrimitiveWithPath_Int(t *testing.T) {
	almanac := gre.NewAlmanac()
	// Add a fact with a primitive int value
	err := almanac.AddFact("int_fact", 42)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Try to apply a path to a primitive int value
	val, err := almanac.GetFactValue("int_fact", nil, "$.number.value")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val != 42 {
		t.Errorf("Expected 42, got %v", val)
	}
}

func TestGetFactValue_PrimitiveWithPath_Bool(t *testing.T) {
	almanac := gre.NewAlmanac()
	// Add a fact with a primitive bool value
	err := almanac.AddFact("bool_fact", true)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Try to apply a path to a primitive bool value
	val, err := almanac.GetFactValue("bool_fact", nil, "$.boolean.flag")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val != true {
		t.Errorf("Expected true, got %v", val)
	}
}

func TestGetFactValue_PrimitiveWithPath_Float(t *testing.T) {
	almanac := gre.NewAlmanac()
	// Add a fact with a primitive float value
	err := almanac.AddFact("float_fact", 3.14)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Try to apply a path to a primitive float value
	val, err := almanac.GetFactValue("float_fact", nil, "$.decimal.value")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val != 3.14 {
		t.Errorf("Expected 3.14, got %v", val)
	}
}

func TestGetFactValue_NilValueWithPath(t *testing.T) {
	almanac := gre.NewAlmanac()
	// Add a dynamic fact that returns nil
	nilFunc := func() interface{} {
		return nil
	}
	err := almanac.AddFact("nil_fact", nilFunc, gre.WithoutCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	val, err := almanac.GetFactValue("nil_fact", nil, "$.some.path")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

func TestTraversePath_ComplexMapWithValidPath(t *testing.T) {
	almanac := gre.NewAlmanac()
	userData := map[string]interface{}{
		"user": map[string]interface{}{
			"profile": map[string]interface{}{
				"address": map[string]interface{}{
					"city":    "Paris",
					"country": "France",
				},
			},
		},
	}
	err := almanac.AddFact("user_data", userData)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	val, err := almanac.GetFactValue("user_data", nil, "$.user.profile.address.city")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val != "Paris" {
		t.Errorf("Expected 'Paris', got %v", val)
	}
}

func TestTraversePath_SliceAccess(t *testing.T) {
	almanac := gre.NewAlmanac()
	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice", "age": 30},
			map[string]interface{}{"name": "Bob", "age": 25},
		},
	}
	err := almanac.AddFact("users_list", data)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	val, err := almanac.GetFactValue("users_list", nil, "$.users[0].name")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val != "Alice" {
		t.Errorf("Expected 'Alice', got %v", val)
	}
}

func TestTraversePath_InvalidPath_ReturnsAlmanacError(t *testing.T) {
	almanac := gre.NewAlmanac()
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
		},
	}
	err := almanac.AddFact("user_info", data)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	_, err = almanac.GetFactValue("user_info", nil, "$.user.nonexistent.field")
	if err == nil {
		t.Fatal("Expected AlmanacError, got nil")
	}

	almanacErr, ok := err.(*gre.AlmanacError)
	if !ok {
		t.Fatalf("Expected *AlmanacError, got %T", err)
	}

	expectedPayload := "factID=user_info, path=$.user.nonexistent.field"
	if almanacErr.Payload != expectedPayload {
		t.Errorf("Expected payload '%s', got '%s'", expectedPayload, almanacErr.Payload)
	}

	if almanacErr.Error() == "" {
		t.Error("Expected non-empty error message")
	}
}

func TestTraversePath_EmptyPath_ReturnsFullValue(t *testing.T) {
	almanac := gre.NewAlmanac()
	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
			"age":  30,
		},
	}
	err := almanac.AddFact("full_data", data)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	val, err := almanac.GetFactValue("full_data", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	resultMap, ok := val.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map[string]interface{}, got %T", val)
	}

	userMap, ok := resultMap["user"].(map[string]interface{})
	if !ok {
		t.Fatal("Expected user to be a map")
	}
	if userMap["name"] != "Alice" {
		t.Errorf("Expected name 'Alice', got %v", userMap["name"])
	}
}

func TestTraversePath_StructValue(t *testing.T) {
	almanac := gre.NewAlmanac()
	user := map[string]interface{}{
		"Name": "Alice",
		"Address": map[string]interface{}{
			"City":    "Lyon",
			"Country": "France",
		},
	}
	err := almanac.AddFact("struct_user", user)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	val, err := almanac.GetFactValue("struct_user", nil, "$.Address.City")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val != "Lyon" {
		t.Errorf("Expected 'Lyon', got %v", val)
	}
}

func TestTraversePath_WithDynamicFact(t *testing.T) {
	almanac := gre.NewAlmanac()
	dynamicFunc := func() interface{} {
		return map[string]interface{}{
			"config": map[string]interface{}{
				"timeout": 30,
				"retry":   3,
			},
		}
	}
	err := almanac.AddFact("dynamic_config", dynamicFunc, gre.WithoutCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	val, err := almanac.GetFactValue("dynamic_config", nil, "$.config.timeout")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	switch timeout := val.(type) {
	case int:
		if timeout != 30 {
			t.Errorf("Expected timeout 30, got %d", timeout)
		}
	case float64:
		if timeout != 30 {
			t.Errorf("Expected timeout 30, got %f", timeout)
		}
	default:
		t.Errorf("Expected int or float64, got %T", val)
	}
}

func TestTraversePath_Wildcard(t *testing.T) {
	almanac := gre.NewAlmanac()
	data := map[string]interface{}{
		"users": []interface{}{
			map[string]interface{}{"name": "Alice", "role": "admin"},
			map[string]interface{}{"name": "Bob", "role": "user"},
			map[string]interface{}{"name": "Charlie", "role": "user"},
		},
	}
	err := almanac.AddFact("all_users", data)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	val, err := almanac.GetFactValue("all_users", nil, "$.users[*].name")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	names, ok := val.([]interface{})
	if !ok {
		t.Fatalf("Expected []interface{}, got %T", val)
	}

	if len(names) != 3 {
		t.Errorf("Expected 3 names, got %d", len(names))
	}

	expectedNames := map[string]bool{"Alice": false, "Bob": false, "Charlie": false}
	for _, name := range names {
		if n, ok := name.(string); ok {
			expectedNames[n] = true
		}
	}

	for name, found := range expectedNames {
		if !found {
			t.Errorf("Expected to find name '%s'", name)
		}
	}
}

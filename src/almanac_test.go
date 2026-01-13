package gorulesengine_test

import (
	"reflect"
	"sync"
	"testing"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func TestPathResolver(t *testing.T) {
	resolver := gorulesengine.DefaultPathResolver

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
		t.Fatalf("Expected 'Alice', got %v", value)
	}
}

func TestEmptyPathResolver(t *testing.T) {
	resolver := gorulesengine.DefaultPathResolver

	data := map[string]interface{}{
		"user": map[string]interface{}{
			"name": "Alice",
			"age":  30,
		},
	}

	value, err := resolver(data, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if !reflect.DeepEqual(data, value) {
		t.Fatalf("Expected original data, got %v", value)
	}
}

func TestNewAlmanac(t *testing.T) {
	facts := []*gorulesengine.Fact{
		gorulesengine.NewFact("fact1", "value"),
		gorulesengine.NewFact("fact2", "value"),
	}
	opts := []gorulesengine.AlmanacOption{
		gorulesengine.AllowUndefinedFacts(),
	}

	almanac := gorulesengine.NewAlmanac(facts, opts...)

	if almanac == nil {
		t.Fatal("Expected almanac to be created, got nil")
	}

	expectedOpts := almanac.GetOptions()
	if allowUndefined, ok := expectedOpts[gorulesengine.ALMANAC_OPTION_KEY_ALLOW_UNDEFINED_FACTS]; !ok || allowUndefined != true {
		t.Fatalf("Expected allowUndefinedFacts to be true, got %v", allowUndefined)
	}

	retrievedFacts := almanac.GetFacts()
	if len(retrievedFacts) != 2 {
		t.Fatalf("Expected 2 facts, got %d", len(retrievedFacts))
	}

	if retrievedFacts["fact1"].ID() != "fact1" {
		t.Fatalf("Expected fact ID 'fact1', got %v", retrievedFacts["fact1"].ID())
	}

	if retrievedFacts["fact2"].ID() != "fact2" {
		t.Fatalf("Expected fact ID 'fact2', got %v", retrievedFacts["fact2"].ID())
	}
}

func TestAddFact(t *testing.T) {
	facts := []*gorulesengine.Fact{}
	opts := []gorulesengine.AlmanacOption{}
	almanac := gorulesengine.NewAlmanac(facts, opts...)

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
}

func TestAddCachedFact(t *testing.T) {
	facts := []*gorulesengine.Fact{}
	opts := []gorulesengine.AlmanacOption{}
	almanac := gorulesengine.NewAlmanac(facts, opts...)

	factValue := "cached_value"
	almanac.AddFact("cached_fact", factValue, gorulesengine.WithCache())

	retrievedFact, err := almanac.GetFactValue("cached_fact", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if retrievedFact != "cached_value" {
		t.Fatalf("Expected fact value 'cached_value', got %v", retrievedFact)
	}
}

func TestAddFactCacheKeyError(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Avec hashFromID, la génération de cache key ne peut plus échouer
	// car on utilise simplement l'ID du fait
	// Ce test vérifie maintenant qu'on peut ajouter n'importe quelle valeur avec cache
	unmarshalableValue := make(chan int)

	// Ajouter un fait avec une valeur non-marshalable devrait fonctionner maintenant
	err := almanac.AddFact("channel_fact", unmarshalableValue, gorulesengine.WithCache())

	if err != nil {
		t.Fatalf("Expected no error with hashFromID, got: %v", err)
	}

	// Vérifier que le fait a bien été ajouté
	facts := almanac.GetFacts()
	if _, exists := facts["channel_fact"]; !exists {
		t.Error("Expected fact to be added successfully")
	}
}

func TestGetFactValue_UndefinedFactNotAllowed(t *testing.T) {
	// Créer une option personnalisée pour désactiver allowUndefinedFacts
	disallowUndefinedFacts := func(a *gorulesengine.Almanac) {
		opts := a.GetOptions()
		opts[gorulesengine.ALMANAC_OPTION_KEY_ALLOW_UNDEFINED_FACTS] = false
	}

	// Créer un almanac avec allowUndefinedFacts désactivé
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{}, disallowUndefinedFacts)

	// Essayer de récupérer un fait qui n'existe pas
	val, err := almanac.GetFactValue("nonexistent_fact", nil, "")

	// Devrait retourner une erreur
	if err == nil {
		t.Fatal("Expected error when getting undefined fact, got nil")
	}

	// Vérifier que c'est une AlmanacError
	almanacErr, ok := err.(*gorulesengine.AlmanacError)
	if !ok {
		t.Fatalf("Expected *AlmanacError, got %T: %v", err, err)
	}

	// Vérifier le payload
	if almanacErr.Payload != "factID=nonexistent_fact" {
		t.Errorf("Expected payload 'factID=nonexistent_fact', got '%s'", almanacErr.Payload)
	}

	// Vérifier le message d'erreur
	if almanacErr.Err == nil {
		t.Fatal("Expected wrapped error, got nil")
	}

	// La valeur retournée devrait être nil
	if val != nil {
		t.Errorf("Expected nil value, got %v", val)
	}
}

func TestGetFactValue_UndefinedFactAllowed(t *testing.T) {
	// Créer un almanac avec allowUndefinedFacts activé (c'est le défaut de NewAlmanac)
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{}, gorulesengine.AllowUndefinedFacts())

	// Essayer de récupérer un fait qui n'existe pas
	val, err := almanac.GetFactValue("nonexistent_fact", nil, "")

	// Ne devrait PAS retourner d'erreur
	if err != nil {
		t.Fatalf("Expected no error when allowUndefinedFacts is true, got: %v", err)
	}

	// La valeur devrait être nil
	if val != nil {
		t.Errorf("Expected nil value for undefined fact, got %v", val)
	}
}

func TestGetFactValue_FromCache(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter un fait statique avec cache activé
	// Le pré-cache devrait être effectué dans AddFact
	err := almanac.AddFact("cached_static", "static_value", gorulesengine.WithCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Première récupération - devrait venir du cache pré-rempli
	val1, err := almanac.GetFactValue("cached_static", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val1 != "static_value" {
		t.Errorf("Expected 'static_value', got %v", val1)
	}

	// Deuxième récupération - devrait aussi venir du cache
	val2, err := almanac.GetFactValue("cached_static", nil, "")
	if err != nil {
		t.Fatalf("Expected no error on second retrieval, got %v", err)
	}

	if val2 != "static_value" {
		t.Errorf("Expected 'static_value', got %v", val2)
	}
}

func TestGetFactValue_CalculateAndCache(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	callCount := 0
	// Ajouter un fait dynamique qui compte les appels
	dynamicFact := func(params map[string]interface{}) int {
		callCount++
		return 100 + callCount
	}

	// Ajouter avec cache désactivé pour voir si ça calcule à chaque fois
	err := almanac.AddFact("dynamic_no_cache", dynamicFact, gorulesengine.WithoutCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Récupérer deux fois - devrait calculer à chaque fois (pas de cache)
	val1, err := almanac.GetFactValue("dynamic_no_cache", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	val2, err := almanac.GetFactValue("dynamic_no_cache", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// Les valeurs devraient être différentes car callCount augmente
	if val1 == val2 {
		t.Errorf("Expected different values without cache, got %v and %v", val1, val2)
	}

	if callCount != 2 {
		t.Errorf("Expected 2 calls without cache, got %d", callCount)
	}
}

func TestGetFactValue_WithoutCache(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter un fait sans cache (par défaut pour les faits statiques)
	err := almanac.AddFact("no_cache_fact", "simple_value")
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Récupérer la valeur - devrait calculer/retourner directement
	val, err := almanac.GetFactValue("no_cache_fact", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val != "simple_value" {
		t.Errorf("Expected 'simple_value', got %v", val)
	}
}

func TestGetFactValue_DynamicWithCaching(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	callCount := 0
	// Créer un fait dynamique avec cache activé
	// Maintenant que GetCacheKey utilise l'ID pour les faits dynamiques,
	// le cache devrait fonctionner
	dynamicFact := func() int {
		callCount++
		return 42
	}

	// Ajouter avec cache activé
	err := almanac.AddFact("dynamic_cached", dynamicFact, gorulesengine.WithCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Première récupération - devrait calculer et mettre en cache
	val1, err := almanac.GetFactValue("dynamic_cached", nil, "")
	if err != nil {
		t.Fatalf("Expected no error on first call, got %v", err)
	}

	if val1 != 42 {
		t.Errorf("Expected 42, got %v", val1)
	}

	// Deuxième récupération - devrait venir du cache (pas de recalcul)
	val2, err := almanac.GetFactValue("dynamic_cached", nil, "")
	if err != nil {
		t.Fatalf("Expected no error on second call, got %v", err)
	}

	if val2 != 42 {
		t.Errorf("Expected 42, got %v", val2)
	}

	// Le fait ne devrait avoir été calculé qu'une seule fois (puis mis en cache)
	if callCount != 1 {
		t.Errorf("Expected 1 call (then cached), got %d", callCount)
	}
}

func TestGetFactValue_CacheAfterCalculation(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	callCount := 0
	// Créer un fait dynamique qui retourne une valeur marshalable
	// Pour que le cache puisse fonctionner, on doit retourner une valeur simple
	dynamicValue := 100
	dynamicFact := func() int {
		callCount++
		return dynamicValue
	}

	// Ajouter SANS cache pour commencer
	err := almanac.AddFact("dynamic_fact", dynamicFact, gorulesengine.WithoutCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Première récupération - devrait calculer
	val1, err := almanac.GetFactValue("dynamic_fact", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val1 != 100 {
		t.Errorf("Expected 100, got %v", val1)
	}

	// Deuxième récupération - devrait calculer à nouveau (pas de cache)
	val2, err := almanac.GetFactValue("dynamic_fact", nil, "")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val2 != 100 {
		t.Errorf("Expected 100, got %v", val2)
	}

	// Devrait avoir été appelé 2 fois
	if callCount != 2 {
		t.Errorf("Expected 2 calls without cache, got %d", callCount)
	}
}

func TestGetFactValue_ConcurrentAccess(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter un fait statique avec cache
	err := almanac.AddFact("concurrent_fact", "shared_value", gorulesengine.WithCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Lancer plusieurs goroutines qui accèdent au même fait
	const numGoroutines = 100
	results := make(chan interface{}, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			val, err := almanac.GetFactValue("concurrent_fact", nil, "")
			if err != nil {
				errors <- err
				return
			}
			results <- val
		}()
	}

	// Collecter les résultats
	for i := 0; i < numGoroutines; i++ {
		select {
		case err := <-errors:
			t.Fatalf("Unexpected error from goroutine: %v", err)
		case val := <-results:
			if val != "shared_value" {
				t.Errorf("Expected 'shared_value', got %v", val)
			}
		}
	}
}

func TestGetFactValue_ConcurrentCacheWrite(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	callCount := 0
	var mutex sync.Mutex

	// Créer un fait dynamique sans cache pré-rempli
	dynamicFact := func() int {
		mutex.Lock()
		callCount++
		count := callCount
		mutex.Unlock()
		return count
	}

	// Ajouter sans cache pour forcer le calcul à chaque fois
	err := almanac.AddFact("dynamic_concurrent", dynamicFact, gorulesengine.WithoutCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Lancer plusieurs goroutines qui calculent et potentiellement cachent
	const numGoroutines = 50
	results := make(chan interface{}, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			val, err := almanac.GetFactValue("dynamic_concurrent", nil, "")
			if err != nil {
				errors <- err
				return
			}
			results <- val
		}()
	}

	// Collecter les résultats
	seenValues := make(map[int]bool)
	for i := 0; i < numGoroutines; i++ {
		select {
		case err := <-errors:
			t.Fatalf("Unexpected error from goroutine: %v", err)
		case val := <-results:
			intVal, ok := val.(int)
			if !ok {
				t.Errorf("Expected int value, got %T", val)
			}
			seenValues[intVal] = true
		}
	}

	// Vérifier que toutes les goroutines ont bien obtenu une valeur
	if len(seenValues) == 0 {
		t.Error("Expected to see some values")
	}

	// Le nombre d'appels devrait être égal au nombre de goroutines (pas de cache)
	mutex.Lock()
	finalCount := callCount
	mutex.Unlock()

	if finalCount != numGoroutines {
		t.Errorf("Expected %d calls without cache, got %d", numGoroutines, finalCount)
	}
}

func TestGetFactValue_DynamicCacheStorageAfterCalculation(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	callCount := 0
	// Créer un fait dynamique avec cache activé
	// Avec la modification de GetCacheKey, les faits dynamiques peuvent maintenant
	// avoir une cache key basée sur leur ID
	dynamicFact := func() int {
		callCount++
		return 42
	}

	// Ajouter le fait dynamique avec cache activé
	// Les faits dynamiques ne sont PAS pré-cachés dans AddFact (condition !fact.IsDynamic())
	err := almanac.AddFact("dynamic_cacheable", dynamicFact, gorulesengine.WithCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Première récupération - devrait calculer ET mettre en cache (lignes 148-150!)
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

	// Deuxième récupération - devrait venir du cache (pas de nouveau calcul)
	val2, err := almanac.GetFactValue("dynamic_cacheable", nil, "")
	if err != nil {
		t.Fatalf("Expected no error on second call, got %v", err)
	}

	if val2 != 42 {
		t.Errorf("Expected 42 from cache, got %v", val2)
	}

	// Le fait ne devrait PAS avoir été recalculé (toujours 1 appel)
	if callCount != 1 {
		t.Errorf("Expected 1 call total (second from cache), got %d", callCount)
	}

	// Troisième récupération - toujours du cache
	val3, err := almanac.GetFactValue("dynamic_cacheable", nil, "")
	if err != nil {
		t.Fatalf("Expected no error on third call, got %v", err)
	}

	if val3 != 42 {
		t.Errorf("Expected 42 from cache, got %v", val3)
	}

	if callCount != 1 {
		t.Errorf("Expected still 1 call total, got %d", callCount)
	}
}

func TestGetFactValue_PrimitiveWithPath_String(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter un fait avec une valeur string primitive
	err := almanac.AddFact("string_fact", "simple_string")
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Essayer d'appliquer un path à une valeur primitive string
	// Devrait retourner la valeur telle quelle (ligne 148-150)
	val, err := almanac.GetFactValue("string_fact", nil, "$.some.path")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val != "simple_string" {
		t.Errorf("Expected 'simple_string', got %v", val)
	}
}

func TestGetFactValue_PrimitiveWithPath_Int(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter un fait avec une valeur int primitive
	err := almanac.AddFact("int_fact", 42)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Essayer d'appliquer un path à une valeur primitive int
	// Devrait retourner la valeur telle quelle
	val, err := almanac.GetFactValue("int_fact", nil, "$.number.value")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val != 42 {
		t.Errorf("Expected 42, got %v", val)
	}
}

func TestGetFactValue_PrimitiveWithPath_Bool(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter un fait avec une valeur bool primitive
	err := almanac.AddFact("bool_fact", true)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Essayer d'appliquer un path à une valeur primitive bool
	// Devrait retourner la valeur telle quelle
	val, err := almanac.GetFactValue("bool_fact", nil, "$.boolean.flag")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val != true {
		t.Errorf("Expected true, got %v", val)
	}
}

func TestGetFactValue_PrimitiveWithPath_Float(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter un fait avec une valeur float primitive
	err := almanac.AddFact("float_fact", 3.14)
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Essayer d'appliquer un path à une valeur primitive float
	// Devrait retourner la valeur telle quelle
	val, err := almanac.GetFactValue("float_fact", nil, "$.decimal.value")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val != 3.14 {
		t.Errorf("Expected 3.14, got %v", val)
	}
}

func TestGetFactValue_NilValueWithPath(t *testing.T) {
	almanac := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})

	// Ajouter un fait dynamique qui retourne nil
	nilFunc := func() interface{} {
		return nil
	}

	err := almanac.AddFact("nil_fact", nilFunc, gorulesengine.WithoutCache())
	if err != nil {
		t.Fatalf("Unexpected error adding fact: %v", err)
	}

	// Essayer d'appliquer un path à une valeur nil retournée par la fonction
	// Le TypeOf(nil) retourne nil, donc valType sera nil
	// La condition valType != nil sera false, donc on retournera nil directement
	val, err := almanac.GetFactValue("nil_fact", nil, "$.some.path")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if val != nil {
		t.Errorf("Expected nil, got %v", val)
	}
}

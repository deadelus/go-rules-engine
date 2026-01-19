package gorulesengine

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestHTTPRuleProvider(t *testing.T) {
	rules := []*Rule{
		{
			Name: "test-rule",
			Conditions: ConditionSet{
				All: []ConditionNode{
					{
						Condition: &Condition{
							Fact:     "age",
							Operator: "equal",
							Value:    18,
						},
					},
				},
			},
		},
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(rules)
	}))
	defer server.Close()

	provider := NewHTTPRuleProvider(server.URL)
	fetchedRules, err := provider.FetchRules(context.Background())

	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}
	if len(fetchedRules) != 1 {
		t.Fatalf("Expected 1 rule, got %d", len(fetchedRules))
	}
	if fetchedRules[0].Name != "test-rule" {
		t.Fatalf("Expected rule name 'test-rule', got '%s'", fetchedRules[0].Name)
	}
}

func TestHotReloader(t *testing.T) {
	var callCount int32
	rules1 := []*Rule{{Name: "rule1"}}
	rules2 := []*Rule{{Name: "rule2"}}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)
		if count == 1 {
			json.NewEncoder(w).Encode(rules1)
		} else {
			json.NewEncoder(w).Encode(rules2)
		}
	}))
	defer server.Close()

	engine := NewEngine()
	provider := NewHTTPRuleProvider(server.URL)
	reloader := NewHotReloader(engine, provider, 100*time.Millisecond)

	var updateCount int32
	reloader.OnUpdate(func(rules []*Rule) {
		atomic.AddInt32(&updateCount, 1)
	})

	reloader.Start(context.Background())
	if !reloader.running {
		t.Error("Reloader should be running")
	}
	reloader.Start(context.Background()) // Should return immediately
	reloader.Stop()
	if reloader.running {
		t.Error("Reloader should NOT be running")
	}
	reloader.Stop() // Should return immediately
}

func TestHotReloader_ContextDone(t *testing.T) {
	provider := NewHTTPRuleProvider("http://localhost")
	reloader := NewHotReloader(NewEngine(), provider, 1*time.Hour)
	ctx, cancel := context.WithCancel(context.Background())

	// Start reloader
	reloader.Start(ctx)

	// Cancel context
	cancel()

	// Wait a bit to ensure the goroutine has time to see the cancellation
	// before we call Stop() which would close stopChan.
	time.Sleep(20 * time.Millisecond)

	reloader.Stop()
}

func TestHotReloader_NilCallbacks(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`[{"Name":"r1"}]`))
	}))
	defer server.Close()

	engine := NewEngine()
	provider := NewHTTPRuleProvider(server.URL)
	reloader := NewHotReloader(engine, provider, 10*time.Millisecond)
	// No OnUpdate, No OnError set

	reloader.Start(context.Background())
	time.Sleep(50 * time.Millisecond)
	reloader.Stop()
}

func TestHotReloader_Error_NilCallback(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	engine := NewEngine()
	provider := NewHTTPRuleProvider(server.URL)
	reloader := NewHotReloader(engine, provider, 10*time.Millisecond)
	// No OnError set

	reloader.Start(context.Background())
	time.Sleep(50 * time.Millisecond)
	reloader.Stop()
}

func TestHTTPRuleProvider_ErrorCases(t *testing.T) {
	t.Run("invalid url", func(t *testing.T) {
		provider := NewHTTPRuleProvider("http://invalid-url-that-does-not-exist-123.com")
		_, err := provider.FetchRules(context.Background())
		if err == nil {
			t.Error("Expected error for invalid URL")
		}
	})

	t.Run("malformed json", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Write([]byte("invalid json"))
		}))
		defer server.Close()

		provider := NewHTTPRuleProvider(server.URL)
		_, err := provider.FetchRules(context.Background())
		if err == nil {
			t.Error("Expected error for malformed JSON")
		}
	})

	t.Run("unexpected status code", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusNotFound)
		}))
		defer server.Close()

		provider := NewHTTPRuleProvider(server.URL)
		_, err := provider.FetchRules(context.Background())
		if err == nil {
			t.Error("Expected error for 404")
		}
	})

	t.Run("request creation failure", func(t *testing.T) {
		// A URL with a control character should fail NewRequestWithContext
		provider := NewHTTPRuleProvider("http://localhost\n")
		_, err := provider.FetchRules(context.Background())
		if err == nil {
			t.Error("Expected error for invalid URL character")
		} else {
			reErr, ok := err.(*RuleEngineError)
			if !ok || reErr.Msg == "" || reErr.Type != ErrLoader {
				t.Errorf("Expected RuleEngineError with ErrLoader, got %v", err)
			}
		}
	})

	t.Run("failed to read body", func(t *testing.T) {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "100")
			w.WriteHeader(http.StatusOK)
			// Connection will be closed when the handler returns,
			// but we didn't write 100 bytes.
		}))
		defer server.Close()

		provider := NewHTTPRuleProvider(server.URL)
		_, err := provider.FetchRules(context.Background())
		if err == nil {
			t.Error("Expected error for incomplete body")
		}
	})
}

func TestHTTPRuleProvider_ETag(t *testing.T) {
	etag := "v1"
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		if r.Header.Get("If-None-Match") == etag {
			w.WriteHeader(http.StatusNotModified)
			return
		}
		w.Header().Set("ETag", etag)
		w.Write([]byte(`[]`))
	}))
	defer server.Close()

	provider := NewHTTPRuleProvider(server.URL)

	// First call - should get rules and set ETag
	_, err := provider.FetchRules(context.Background())
	if err != nil {
		t.Fatalf("First call failed: %v", err)
	}
	if provider.ETag != etag {
		t.Errorf("Expected ETag %s, got %s", etag, provider.ETag)
	}

	// Second call - should send ETag and get 304 (returns nil, nil)
	rules, err := provider.FetchRules(context.Background())
	if err != nil {
		t.Fatalf("Second call failed: %v", err)
	}
	if rules != nil {
		t.Error("Expected nil rules for 304 response")
	}
	if callCount != 2 {
		t.Errorf("Expected 2 calls to server, got %d", callCount)
	}
}

func TestHotReloader_Error(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	engine := NewEngine()
	provider := NewHTTPRuleProvider(server.URL)
	reloader := NewHotReloader(engine, provider, 50*time.Millisecond)

	var errorCount int32
	reloader.OnError(func(err error) {
		atomic.AddInt32(&errorCount, 1)
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	reloader.Start(ctx)
	time.Sleep(150 * time.Millisecond)
	reloader.Stop()

	if atomic.LoadInt32(&errorCount) == 0 {
		t.Errorf("Expected errors during hot reload")
	}
}

func TestHotReloader_NoUpdate(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotModified)
	}))
	defer server.Close()

	engine := NewEngine()
	provider := NewHTTPRuleProvider(server.URL)
	reloader := NewHotReloader(engine, provider, 10*time.Millisecond)

	var updateCount int32
	reloader.OnUpdate(func(rules []*Rule) {
		atomic.AddInt32(&updateCount, 1)
	})

	reloader.Start(context.Background())
	time.Sleep(50 * time.Millisecond)
	reloader.Stop()

	if atomic.LoadInt32(&updateCount) > 0 {
		t.Errorf("Expected no updates when provider returns nil, got %d", updateCount)
	}
}

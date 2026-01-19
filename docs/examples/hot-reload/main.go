package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"time"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func main() {
	fmt.Println("🔥 Hot-reload Example")
	fmt.Println("====================")

	// 1. Setup a mock server to simulate a remote rule source
	var callCount int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		count := atomic.AddInt32(&callCount, 1)

		var rules []*gre.Rule
		if count == 1 {
			fmt.Println("🌐 [Server] Serving Version 1 of rules...")
			rules = []*gre.Rule{
				{
					Name: "rule-v1",
					Conditions: gre.ConditionSet{
						All: []gre.ConditionNode{
							{
								Condition: &gre.Condition{
									Fact: "x", Operator: "equal", Value: 1,
								},
							},
						},
					},
					OnSuccess: []gre.RuleEvent{{Name: "notify"}},
				},
			}
		} else {
			fmt.Println("🌐 [Server] Serving Version 2 of rules...")
			rules = []*gre.Rule{
				{
					Name: "rule-v2",
					Conditions: gre.ConditionSet{
						All: []gre.ConditionNode{
							{
								Condition: &gre.Condition{
									Fact: "x", Operator: "equal", Value: 1,
								},
							},
						},
					},
					OnSuccess: []gre.RuleEvent{{Name: "notify"}},
				},
			}
		}

		json.NewEncoder(w).Encode(rules)
	}))
	defer server.Close()

	// 2. Initialize Engine
	engine := gre.NewEngine()
	engine.RegisterEvent(gre.Event{
		Name: "notify",
		Action: func(ctx gre.EventContext) error {
			fmt.Printf("🔔 [Event] Rule '%s' triggered successfully!\n", ctx.RuleName)
			return nil
		},
	})

	// 3. Setup HotReloader (polling every 2 seconds)
	provider := gre.NewHTTPRuleProvider(server.URL)
	reloader := gre.NewHotReloader(engine, provider, 2*time.Second)

	reloader.OnUpdate(func(rules []*gre.Rule) {
		fmt.Printf("🔄 [Reloader] Successfully updated to %d rules\n", len(rules))
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	fmt.Println("🚀 Starting HotReloader...")
	reloader.Start(ctx)

	// 4. Run evaluation loop
	almanac := gre.NewAlmanac(nil)
	almanac.AddFact("x", 1)

	for i := 0; i < 3; i++ {
		fmt.Printf("\n--- Evaluation Step %d ---\n", i+1)
		_, err := engine.Run(almanac)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
		}

		fmt.Println("Waiting for next reload...")
		time.Sleep(2500 * time.Millisecond) // Slightly more than reload interval
	}

	reloader.Stop()
	fmt.Println("\n🏁 Done!")
}

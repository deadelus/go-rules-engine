package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

// Request defines the input for our API
type Request struct {
	Facts map[string]interface{} `json:"facts"`
}

// Response defines the output of our API
type Response struct {
	Success bool                       `json:"success"`
	Results map[string]*gre.RuleResult `json:"results"`
	Errors  []string                   `json:"errors"`
}

func main() {
	// 1. Initialize the engine with some optimizations
	engine := gre.NewEngine(
		gre.WithConditionCaching(),
		gre.WithAuditTrace(),
		gre.WithSmartSkip(),
	)

	// 2. Define some business rules (e.g., Eligibility & Pricing)
	engine.AddRule(&gre.Rule{
		Name: "vip-discount",
		Conditions: gre.ConditionSet{
			All: []gre.ConditionNode{
				{
					Condition: &gre.Condition{
						Fact:     "totalSpend",
						Operator: "greater_than_inclusive",
						Value:    1000,
					},
				},
				{
					Condition: &gre.Condition{
						Fact:     "accountAgeDays",
						Operator: "greater_than",
						Value:    365,
					},
				},
			},
		},
		OnSuccess: []gre.RuleEvent{{Name: "apply-vip-badge"}},
	})

	engine.AddRule(&gre.Rule{
		Name: "high-risk-transaction",
		Conditions: gre.ConditionSet{
			Any: []gre.ConditionNode{
				{
					Condition: &gre.Condition{
						Fact:     "amount",
						Operator: "greater_than",
						Value:    5000,
					},
				},
				{
					Condition: &gre.Condition{
						Fact:     "isFirstPurchase",
						Operator: "equal",
						Value:    true,
					},
				},
			},
		},
	})

	// 3. Define the HTTP Handler
	http.HandleFunc("/evaluate", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			sendJSON(w, http.StatusBadRequest, Response{
				Success: false,
				Errors:  []string{"Invalid JSON body: " + err.Error()},
			})
			return
		}

		// Create Almanac from request facts
		almanac := gre.NewAlmanac()
		for k, v := range req.Facts {
			almanac.AddFact(gre.FactID(k), v)
		}

		// Execute performance tracking
		start := time.Now()
		_, err := engine.Run(almanac)
		duration := time.Since(start)

		if err != nil {
			sendJSON(w, http.StatusInternalServerError, Response{
				Success: false,
				Errors:  []string{err.Error()},
			})
			return
		}

		// Return results with audit trace
		// We set the execution time in a header for visibility
		w.Header().Set("X-Execution-Time", duration.String())
		sendJSON(w, http.StatusOK, Response{
			Success: true,
			Results: engine.Results(),
		})
	})

	fmt.Println("🚀 Rules Engine API listening on :8080")
	fmt.Println("Example Payload:")
	fmt.Println("curl -X POST http://localhost:8080/evaluate -H \"Content-Type: application/json\" -d '{\"facts\": {\"totalSpend\": 1200, \"accountAgeDays\": 400, \"amount\": 6000, \"isFirstPurchase\": false}}'")

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func sendJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

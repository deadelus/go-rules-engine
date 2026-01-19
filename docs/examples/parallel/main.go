package main

import (
	"fmt"
	"time"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

func main() {
	fmt.Println("🚀 Example Parallel Execution - Advanced Performance")
	fmt.Println("=====================================================")

	// Creating an engine with Parallel Execution enabled (using 5 workers)
	// This is particularly useful when facts are dynamic and perform slow operations
	// like API calls or complex calculations.
	engine := gre.NewEngine(
		gre.WithParallelExecution(5),
		gre.WithAuditTrace(), // Always good to see what happened
	)

	// Registering an event to track execution
	engine.RegisterEvent(gre.Event{
		Name: "log-success",
		Action: func(ctx gre.EventContext) error {
			fmt.Printf("   ✅ Rule '%s' passed!\n", ctx.RuleName)
			return nil
		},
	})

	// Let's create multiple rules that depend on slow dynamic facts
	for i := 1; i <= 5; i++ {
		engine.AddRule(&gre.Rule{
			Name:     fmt.Sprintf("slow-rule-%d", i),
			Priority: i * 10,
			Conditions: gre.ConditionSet{
				All: []gre.ConditionNode{
					{
						Condition: &gre.Condition{
							Fact:     gre.FactID(fmt.Sprintf("slow-fact-%d", i)),
							Operator: "equal",
							Value:    true,
						},
					},
				},
			},
			OnSuccess: []gre.RuleEvent{{Name: "log-success"}},
		})
	}

	// Setting up the Almanac with slow dynamic facts
	almanac := gre.NewAlmanac()

	for i := 1; i <= 5; i++ {
		factID := gre.FactID(fmt.Sprintf("slow-fact-%d", i))
		// Each fact takes 200ms to compute
		almanac.AddFact(factID, func() interface{} {
			fmt.Printf("   ⏳ Computing %s...\n", factID)
			time.Sleep(200 * time.Millisecond)
			return true
		})
	}

	fmt.Println("\n🏃 Running evaluation...")
	start := time.Now()

	_, err := engine.Run(almanac)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	duration := time.Since(start)
	fmt.Printf("\n✨ Evaluation finished in: %v\n", duration)
	fmt.Println("   (In sequential mode, it would have taken at least 1000ms)")

	// The results are still consistent and accessible
	fmt.Printf("\n📊 Total rules evaluated: %d\n", len(engine.Results()))
}

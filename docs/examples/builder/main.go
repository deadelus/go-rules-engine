package main

import (
	"fmt"
	"log"
	"strings"

	gre "github.com/deadelus/go-rules-engine/v2/src"
)

// GlobalHandler handles all events
type GlobalHandler struct{}

// Handle processes events triggered by rules.
func (h *GlobalHandler) Handle(event gre.Event, ctx gre.EventContext) error {
	fmt.Printf("✅ Rule '%s' succeeded → Event '%s'\n", ctx.RuleName, event.Name)
	return nil
}

func main() {
	fmt.Println("🏗️  Go Rules Engine - Builder API Demo")
	fmt.Println(strings.Repeat("=", 52))
	fmt.Println()

	// Create engine with descending priority
	engine := gre.NewEngine()
	engine.SetEventHandler(&GlobalHandler{})

	// Register events with actions
	engine.RegisterEvent(gre.Event{
		Name: "send-welcome-email",
		Mode: gre.EventModeSync,
		Params: map[string]interface{}{
			"tier": "platinum",
		},
		Action: func(ctx gre.EventContext) error {
			email, _ := ctx.Almanac.GetFactValue("email", nil, "")
			tier := ctx.Params["tier"]
			fmt.Printf("📧 Sending welcome email to %v (tier: %s)\n", email, tier)
			return nil
		},
	})

	engine.RegisterEvent(gre.Event{
		Name: "apply-discount",
		Mode: gre.EventModeSync,
		Params: map[string]interface{}{
			"discount": 20,
		},
		Action: func(ctx gre.EventContext) error {
			discount := ctx.Params["discount"]
			fmt.Printf("💰 Applying %v%% discount\n", discount)
			return nil
		},
	})

	engine.RegisterEvent(gre.Event{
		Name: "access-denied",
		Mode: gre.EventModeSync,
		Params: map[string]interface{}{
			"reason": "User is underage",
		},
		Action: func(ctx gre.EventContext) error {
			reason := ctx.Params["reason"]
			age, _ := ctx.Almanac.GetFactValue("age", nil, "")
			fmt.Printf("🚫 Access denied: %v (age: %v)\n", reason, age)
			return nil
		},
	})

	// ========================================
	// Example 1: Simple Rule with Equal
	// ========================================
	fmt.Println("📝 Creating rules using Builder API...")

	// Register simple events
	engine.RegisterEvents(
		gre.Event{Name: "adult-verified", Mode: gre.EventModeSync,
			Params: map[string]interface{}{"message": "User is an adult"}},
		gre.Event{Name: "minor-detected", Mode: gre.EventModeSync,
			Params: map[string]interface{}{"message": "User is a minor"}},
		gre.Event{Name: "premium-benefits", Mode: gre.EventModeSync},
		gre.Event{Name: "vip-status-granted", Mode: gre.EventModeSync},
		gre.Event{Name: "valid-email", Mode: gre.EventModeSync},
		gre.Event{Name: "country-allowed", Mode: gre.EventModeSync},
		gre.Event{Name: "content-approved", Mode: gre.EventModeSync},
	)

	adultRule := gre.NewRuleBuilder().
		WithName("adult-verification").
		WithPriority(100).
		WithConditions(gre.ConditionNode{
			SubSet: &gre.ConditionSet{
				All: []gre.ConditionNode{
					{Condition: gre.GreaterThanInclusive("age", float64(18))},
					{Condition: gre.NotEqual("status", "banned")},
				},
			},
		}).
		WithOnSuccess("adult-verified").
		WithOnFailure("access-denied").
		Build()

	engine.AddRule(adultRule)

	// ========================================
	// Example 2: Premium User Rule with helpers
	// ========================================
	premiumRule := gre.NewRuleBuilder().
		WithName("premium-user-benefits").
		WithPriority(90).
		WithConditions(gre.ConditionNode{
			SubSet: &gre.ConditionSet{
				All: []gre.ConditionNode{
					{Condition: gre.Equal("membership", "premium")},
					{Condition: gre.GreaterThan("accountAge", float64(365))},
				},
			},
		}).
		WithOnSuccess("premium-benefits", "apply-discount").
		Build()

	engine.AddRule(premiumRule)

	// ========================================
	// Example 3: Complex nested conditions using helper functions
	// ========================================
	anySets := gre.AnySets(
		gre.All(
			gre.Equal("membership", "vip"),
			gre.GreaterThanInclusive("totalSpent", 10000),
		),
		gre.All(
			gre.Equal("membership", "premium"),
			gre.GreaterThan("totalSpent", 20000),
		),
	)

	vipRule := gre.NewRuleBuilder().
		WithName("vip-user-detection").
		WithPriority(95).
		WithConditions(gre.ConditionNode{
			SubSet: &anySets,
		}).
		WithOnSuccess("vip-status-granted", "send-welcome-email").
		Build()

	engine.AddRule(vipRule)

	// ========================================
	// Example 4: Email validation with Regex
	// ========================================
	emailRule := gre.NewRuleBuilder().
		WithName("email-validation").
		WithPriority(80).
		WithConditions(gre.ConditionNode{
			Condition: gre.Regex("email", "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"),
		}).
		WithOnSuccess("valid-email").
		Build()

	engine.AddRule(emailRule)

	// ========================================
	// Example 5: Country restrictions using In/NotIn
	// ========================================
	restrictedCountryRule := gre.NewRuleBuilder().
		WithName("country-restrictions").
		WithPriority(85).
		WithConditions(gre.ConditionNode{
			SubSet: &gre.ConditionSet{
				All: []gre.ConditionNode{
					{Condition: gre.NotIn("country", []string{"XX", "YY", "ZZ"})},
					{Condition: gre.In("region", []string{"EU", "NA", "APAC"})},
				},
			},
		}).
		WithOnSuccess("country-allowed").
		Build()

	engine.AddRule(restrictedCountryRule)

	// ========================================
	// Example 6: Content moderation with Contains
	// ========================================
	contentRule := gre.NewRuleBuilder().
		WithName("content-moderation").
		WithPriority(70).
		WithConditions(gre.ConditionNode{
			SubSet: &gre.ConditionSet{
				None: []gre.ConditionNode{
					{Condition: gre.Contains("tags", "spam")},
					{Condition: gre.Contains("tags", "offensive")},
				},
			},
		}).
		WithOnSuccess("content-approved").
		Build()

	engine.AddRule(contentRule)

	// ========================================
	// Test Scenario 1: Regular adult user
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🧪 Test Scenario 1: Regular Adult User")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	almanac1 := gre.NewAlmanac()
	almanac1.AddFact("age", 25)
	almanac1.AddFact("status", "active")
	almanac1.AddFact("membership", "free")
	almanac1.AddFact("email", "user@example.com")
	almanac1.AddFact("country", "FR")
	almanac1.AddFact("region", "EU")
	almanac1.AddFact("tags", []string{"verified", "active"})

	e, err := engine.Run(almanac1)

	results1 := e.ReduceResults()

	if err != nil {
		log.Fatal(err)
	}

	matched1 := 0
	for _, passed := range results1 {
		if passed {
			matched1++
		}
	}
	fmt.Printf("\n📊 Results: %d rules matched\n", matched1)

	// ========================================
	// Test Scenario 2: Premium user with benefits
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🧪 Test Scenario 2: Premium User")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	almanac2 := gre.NewAlmanac()
	almanac2.AddFact("age", 30)
	almanac2.AddFact("status", "active")
	almanac2.AddFact("membership", "premium")
	almanac2.AddFact("accountAge", 400)
	almanac2.AddFact("totalSpent", 5000.0)
	almanac2.AddFact("email", "premium@example.com")
	almanac2.AddFact("country", "US")
	almanac2.AddFact("region", "NA")
	almanac2.AddFact("tags", []string{"premium", "verified"})

	e, err = engine.Run(almanac2)

	results2 := e.ReduceResults()

	if err != nil {
		log.Fatal(err)
	}

	matched2 := 0
	for _, passed := range results2 {
		if passed {
			matched2++
		}
	}
	fmt.Printf("\n📊 Results: %d rules matched\n", matched2)

	// ========================================
	// Test Scenario 3: VIP user (high spender)
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🧪 Test Scenario 3: VIP User (High Spender)")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	almanac3 := gre.NewAlmanac()
	almanac3.AddFact("age", 35)
	almanac3.AddFact("status", "active")
	almanac3.AddFact("membership", "premium")
	almanac3.AddFact("accountAge", 800)
	almanac3.AddFact("totalSpent", 25000.0)
	almanac3.AddFact("email", "vip@example.com")
	almanac3.AddFact("country", "JP")
	almanac3.AddFact("region", "APAC")
	almanac3.AddFact("tags", []string{"vip", "verified"})

	e, err = engine.Run(almanac3)

	results3 := e.ReduceResults()

	if err != nil {
		log.Fatal(err)
	}

	matched3 := 0
	for _, passed := range results3 {
		if passed {
			matched3++
		}
	}
	fmt.Printf("\n📊 Results: %d rules matched\n", matched3)

	// ========================================
	// Test Scenario 4: Invalid email and restricted country
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🧪 Test Scenario 4: Invalid Email & Restricted Country")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	almanac4 := gre.NewAlmanac()
	almanac4.AddFact("age", 22)
	almanac4.AddFact("status", "active")
	almanac4.AddFact("membership", "free")
	almanac4.AddFact("email", "invalid-email")
	almanac4.AddFact("country", "XX")
	almanac4.AddFact("region", "OTHER")
	almanac4.AddFact("tags", []string{"spam"})

	e, err = engine.Run(almanac4)

	results4 := e.ReduceResults()

	if err != nil {
		log.Fatal(err)
	}

	matched4 := 0
	for _, passed := range results4 {
		if passed {
			matched4++
		}
	}
	fmt.Printf("\n📊 Results: %d rules matched\n", matched4)

	// ========================================
	// Test Scenario 5: Underage user (OnFailure example)
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🧪 Test Scenario 5: Underage User (OnFailure Event)")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	almanac5 := gre.NewAlmanac()
	almanac5.AddFact("age", 16)
	almanac5.AddFact("status", "active")
	almanac5.AddFact("membership", "free")
	almanac5.AddFact("email", "teen@example.com")
	almanac5.AddFact("country", "FR")
	almanac5.AddFact("region", "EU")
	almanac5.AddFact("tags", []string{"verified"})

	e, err = engine.Run(almanac5)

	results5 := e.ReduceResults()

	if err != nil {
		log.Fatal(err)
	}

	matched5 := 0
	for _, passed := range results5 {
		if passed {
			matched5++
		}
	}
	fmt.Printf("\n📊 Results: %d rules matched\n", matched5)

	// ========================================
	// Summary
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("📈 Summary")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\n✨ Builder API Benefits:")
	fmt.Println("  • Fluent interface for readable rule creation")
	fmt.Println("  • Helper functions: Equal, GreaterThan, Contains, Regex, etc.")
	fmt.Println("  • ConditionSet helpers: All, Any, None, AllSets, AnySets")
	fmt.Println("  • Type-safe and compile-time checked")
	fmt.Println("  • Ergonomic and intuitive API")
	fmt.Println("\n✅ Demo completed successfully!")
}

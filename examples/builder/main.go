package main

import (
	"fmt"
	"log"
	"strings"

	gorulesengine "github.com/deadelus/go-rules-engine/src"
)

func main() {
	fmt.Println("🏗️  Go Rules Engine - Builder API Demo")
	fmt.Println(strings.Repeat("=", 52))
	fmt.Println()

	// Create engine with descending priority
	engine := gorulesengine.NewEngine()

	// Register callbacks
	engine.RegisterCallback("send-welcome-email", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		userID := event.Params["userID"]
		tier := event.Params["tier"]
		fmt.Printf("📧 Sending welcome email to user %v (tier: %s)\n", userID, tier)
		return nil
	})

	engine.RegisterCallback("apply-discount", func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		discount := event.Params["discount"]
		fmt.Printf("💰 Applying %v%% discount\n", discount)
		return nil
	})

	// Global success handler
	engine.OnSuccess(func(event gorulesengine.Event, almanac *gorulesengine.Almanac, result gorulesengine.RuleResult) error {
		fmt.Printf("✅ Rule '%s' succeeded\n", result.Rule.Name)
		return nil
	})

	// ========================================
	// Example 1: Simple Rule with Equal
	// ========================================
	fmt.Println("📝 Creating rules using Builder API...")

	adultRule := gorulesengine.NewRuleBuilder().
		WithName("adult-verification").
		WithPriority(100).
		WithConditions(gorulesengine.ConditionNode{
			SubSet: &gorulesengine.ConditionSet{
				All: []gorulesengine.ConditionNode{
					{Condition: gorulesengine.GreaterThanInclusive("age", 18)},
					{Condition: gorulesengine.NotEqual("status", "banned")},
				},
			},
		}).
		WithEvent("adult-verified", map[string]interface{}{
			"message": "User is an adult",
		}).
		Build()

	engine.AddRule(adultRule)

	// ========================================
	// Example 2: Premium User Rule with helpers
	// ========================================
	premiumRule := gorulesengine.NewRuleBuilder().
		WithName("premium-user-benefits").
		WithPriority(90).
		WithConditions(gorulesengine.ConditionNode{
			SubSet: &gorulesengine.ConditionSet{
				All: []gorulesengine.ConditionNode{
					{Condition: gorulesengine.Equal("membership", "premium")},
					{Condition: gorulesengine.GreaterThan("accountAge", 365)},
				},
			},
		}).
		WithEvent("premium-benefits", map[string]interface{}{
			"discount": 20,
			"userID":   "{{userID}}",
			"tier":     "gold",
		}).
		WithOnSuccess("apply-discount").
		Build()

	engine.AddRule(premiumRule)

	// ========================================
	// Example 3: Complex nested conditions using helper functions
	// ========================================
	anySets := gorulesengine.AnySets(
		gorulesengine.All(
			gorulesengine.Equal("membership", "vip"),
			gorulesengine.GreaterThanInclusive("totalSpent", 10000),
		),
		gorulesengine.All(
			gorulesengine.Equal("membership", "premium"),
			gorulesengine.GreaterThan("totalSpent", 20000),
		),
	)

	vipRule := gorulesengine.NewRuleBuilder().
		WithName("vip-user-detection").
		WithPriority(95).
		WithConditions(gorulesengine.ConditionNode{
			SubSet: &anySets,
		}).
		WithEvent("vip-status-granted", map[string]interface{}{
			"tier":   "platinum",
			"userID": "{{userID}}",
		}).
		WithOnSuccess("send-welcome-email").
		Build()

	engine.AddRule(vipRule)

	// ========================================
	// Example 4: Email validation with Regex
	// ========================================
	emailRule := gorulesengine.NewRuleBuilder().
		WithName("email-validation").
		WithPriority(80).
		WithConditions(gorulesengine.ConditionNode{
			Condition: gorulesengine.Regex("email", "^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\\.[a-zA-Z]{2,}$"),
		}).
		WithEvent("valid-email", map[string]interface{}{
			"message": "Email format is valid",
		}).
		Build()

	engine.AddRule(emailRule)

	// ========================================
	// Example 5: Country restrictions using In/NotIn
	// ========================================
	restrictedCountryRule := gorulesengine.NewRuleBuilder().
		WithName("country-restrictions").
		WithPriority(85).
		WithConditions(gorulesengine.ConditionNode{
			SubSet: &gorulesengine.ConditionSet{
				All: []gorulesengine.ConditionNode{
					{Condition: gorulesengine.NotIn("country", []string{"XX", "YY", "ZZ"})},
					{Condition: gorulesengine.In("region", []string{"EU", "NA", "APAC"})},
				},
			},
		}).
		WithEvent("country-allowed", map[string]interface{}{
			"message": "User country is allowed",
		}).
		Build()

	engine.AddRule(restrictedCountryRule)

	// ========================================
	// Example 6: Content moderation with Contains
	// ========================================
	contentRule := gorulesengine.NewRuleBuilder().
		WithName("content-moderation").
		WithPriority(70).
		WithConditions(gorulesengine.ConditionNode{
			SubSet: &gorulesengine.ConditionSet{
				None: []gorulesengine.ConditionNode{
					{Condition: gorulesengine.Contains("tags", "spam")},
					{Condition: gorulesengine.Contains("tags", "offensive")},
				},
			},
		}).
		WithEvent("content-approved", map[string]interface{}{
			"message": "Content passed moderation",
		}).
		Build()

	engine.AddRule(contentRule)

	// ========================================
	// Test Scenario 1: Regular adult user
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🧪 Test Scenario 1: Regular Adult User")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	almanac1 := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac1.AddFact("age", 25)
	almanac1.AddFact("status", "active")
	almanac1.AddFact("membership", "free")
	almanac1.AddFact("email", "user@example.com")
	almanac1.AddFact("country", "FR")
	almanac1.AddFact("region", "EU")
	almanac1.AddFact("tags", []string{"verified", "active"})

	results1, err := engine.Run(almanac1)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n📊 Results: %d rules matched\n", len(results1))
	for _, result := range results1 {
		if result.Result {
			fmt.Printf("   ✓ %s → %s\n", result.Rule.Name, result.Event.Type)
		}
	}

	// ========================================
	// Test Scenario 2: Premium user with benefits
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🧪 Test Scenario 2: Premium User")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	almanac2 := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac2.AddFact("age", 30)
	almanac2.AddFact("status", "active")
	almanac2.AddFact("membership", "premium")
	almanac2.AddFact("accountAge", 400)
	almanac2.AddFact("totalSpent", 5000.0)
	almanac2.AddFact("email", "premium@example.com")
	almanac2.AddFact("country", "US")
	almanac2.AddFact("region", "NA")
	almanac2.AddFact("tags", []string{"premium", "verified"})

	results2, err := engine.Run(almanac2)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n📊 Results: %d rules matched\n", len(results2))
	for _, result := range results2 {
		if result.Result {
			fmt.Printf("   ✓ %s → %s\n", result.Rule.Name, result.Event.Type)
		}
	}

	// ========================================
	// Test Scenario 3: VIP user (high spender)
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🧪 Test Scenario 3: VIP User (High Spender)")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	almanac3 := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac3.AddFact("age", 35)
	almanac3.AddFact("status", "active")
	almanac3.AddFact("membership", "premium")
	almanac3.AddFact("accountAge", 800)
	almanac3.AddFact("totalSpent", 25000.0)
	almanac3.AddFact("email", "vip@example.com")
	almanac3.AddFact("country", "JP")
	almanac3.AddFact("region", "APAC")
	almanac3.AddFact("tags", []string{"vip", "verified"})

	results3, err := engine.Run(almanac3)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n📊 Results: %d rules matched\n", len(results3))
	for _, result := range results3 {
		if result.Result {
			fmt.Printf("   ✓ %s → %s\n", result.Rule.Name, result.Event.Type)
		}
	}

	// ========================================
	// Test Scenario 4: Invalid email and restricted country
	// ========================================
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("🧪 Test Scenario 4: Invalid Email & Restricted Country")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	almanac4 := gorulesengine.NewAlmanac([]*gorulesengine.Fact{})
	almanac4.AddFact("age", 22)
	almanac4.AddFact("status", "active")
	almanac4.AddFact("membership", "free")
	almanac4.AddFact("email", "invalid-email")
	almanac4.AddFact("country", "XX")
	almanac4.AddFact("region", "OTHER")
	almanac4.AddFact("tags", []string{"spam"})

	results4, err := engine.Run(almanac4)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\n📊 Results: %d rules matched\n", len(results4))
	for _, result := range results4 {
		if result.Result {
			fmt.Printf("   ✓ %s → %s\n", result.Rule.Name, result.Event.Type)
		}
	}

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

package gorulesengine

import (
	"fmt"
	"strings"
)

func PrintRules(rules []Rule) {
	for _, rule := range rules {
		fmt.Printf("Rule: %s (Priority: %d)\n", rule.Name, rule.Priority)
		fmt.Printf("Event: %s\n", rule.Event.Type)
		fmt.Println("Conditions:")
		PrintConditionSet(rule.Conditions, 1)
		fmt.Println()
	}
}

func PrintConditionSet(cs ConditionSet, indent int) {
	prefix := strings.Repeat("  ", indent)

	if len(cs.All) > 0 {
		fmt.Println(prefix + "ALL")
		for _, n := range cs.All {
			PrintNode(n, indent+1)
		}
	}

	if len(cs.Any) > 0 {
		fmt.Println(prefix + "ANY")
		for _, n := range cs.Any {
			PrintNode(n, indent+1)
		}
	}

	if len(cs.None) > 0 {
		fmt.Println(prefix + "NONE")
		for _, n := range cs.None {
			PrintNode(n, indent+1)
		}
	}
}

func PrintNode(n ConditionNode, indent int) {
	prefix := strings.Repeat("  ", indent)

	if n.Condition != nil {
		fmt.Printf(
			"%s- %s %s %v\n",
			prefix,
			n.Condition.Fact,
			n.Condition.Operator,
			n.Condition.Value,
		)
		return
	}

	if n.SubSet != nil {
		PrintConditionSet(*n.SubSet, indent)
	}
}

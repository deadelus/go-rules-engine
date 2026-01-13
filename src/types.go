package gorulesengine

type ConditionType string

const All ConditionType = "all"
const Any ConditionType = "any"
const None ConditionType = "none"

type OperatorType string

const Equal OperatorType = "equal"
const NotEqual OperatorType = "not_equal"
const LessThan OperatorType = "less_than"
const LessThanInclusive OperatorType = "less_than_inclusive"
const GreaterThan OperatorType = "greater_than"
const GreaterThanInclusive OperatorType = "greater_than_inclusive"
const In OperatorType = "in"
const NotIn OperatorType = "not_in"
const Contains OperatorType = "contains"
const NotContains OperatorType = "not_contains"

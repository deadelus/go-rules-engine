package gorulesengine

type Rule struct {
	Name       string       `json:"name,omitempty"`
	Priority   int          `json:"priority,omitempty"`
	Conditions ConditionSet `json:"conditions"`
	Event      Event        `json:"event"`
	OnSuccess  *string      `json:"on_success,omitempty"`
	OnFailure  *string      `json:"on_failure,omitempty"`
}

type Event struct {
	Type   string                 `json:"type"`
	Params map[string]interface{} `json:"params,omitempty"`
}

package validator

import "fmt"

type RuleConfig struct {
	Value        any
	MessageFunc  func(RuleConfig) string
	ValidateFunc func(RuleConfig) bool
}

func Min(min int) *RuleConfig {
	return &RuleConfig{
		MessageFunc: func(r RuleConfig) string {
			return fmt.Sprintf("'%s' does not meet the minimum length of %d", r.Value, min)
		},
		ValidateFunc: func(r RuleConfig) bool {
			val, ok := r.Value.(string)
			if !ok {
				return false
			}

			return len(val) >= min
		},
	}
}

func Max(max int) *RuleConfig {
	return &RuleConfig{
		MessageFunc: func(r RuleConfig) string {
			return fmt.Sprintf("'%s' exceeds the maximum length of %d", r.Value, max)
		},
		ValidateFunc: func(r RuleConfig) bool {
			val, ok := r.Value.(string)
			if !ok {
				return false
			}

			return len(val) <= max
		},
	}
}

func Must(predicate func() bool) *RuleConfig {
	return &RuleConfig{
		MessageFunc: func(r RuleConfig) string {
			return "Validation failed"
		},
		ValidateFunc: func(r RuleConfig) bool {
			return predicate()
		},
	}
}

func Required() *RuleConfig {
	return &RuleConfig{
		MessageFunc: func(r RuleConfig) string {
			return "Required"
		},
		ValidateFunc: func(r RuleConfig) bool {
			val, ok := r.Value.(string)
			if !ok {
				return false
			}

			return len(val) > 0
		},
	}
}

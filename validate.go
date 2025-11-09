package excel

import (
	"github.com/gookit/validate"
)

func Validate(policy *Policy, data map[string]any) *validate.Validation {
	validation := validate.New(data)
	validation.FilterRules(policy.Filters)
	validation.StringRules(policy.Rules)
	validation.AddMessages(policy.Messages)
	validation.AddTranslates(policy.Attributes)
	validation.Validate()
	return validation
}

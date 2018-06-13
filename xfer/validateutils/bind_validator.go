package validateutils

import "fmt"

type ValidatorContext struct {
	Errors map[string]string
}

func (c *ValidatorContext) AddErrorf(field, error string, params ...interface{}) {
	c.AddError(field, fmt.Sprintf(error, params...))
}

func (c *ValidatorContext) AddError(field, error string) {
	if c.Errors == nil {
		c.Errors = map[string]string{}
	}
	c.Errors[field] = error
}

func (c *ValidatorContext) HasErrors() bool {
	return c.Errors != nil && len(c.Errors) > 0
}

// Validable must be implemented if form input (used with bindAndValidate()) needs specific validations which can't
// be specified as part of "valid" tags.
type Validable interface {
	// Validate returns a map with errors. Keys are field names, values are messages.
	Validate(context *ValidatorContext)
}

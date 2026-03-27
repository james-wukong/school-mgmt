package tables

import (
	"fmt"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
)

// Global instances for performance (thread-safe)
var (
	decoder  = form.NewDecoder()
	validate = validator.New()
)

// MapAndValidate decodes GoAdmin values and runs struct-tag validation
func MapAndValidate[T any](values map[string][]string) (*T, error) {
	fmt.Printf("DEBUG: values map content: %#v\n", values)
	var result T

	// 1. Decode map[string][]string into the struct
	if err := decoder.Decode(&result, values); err != nil {
		return nil, err
	}

	// 2. Validate the struct based on 'validate' tags
	if err := validate.Struct(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

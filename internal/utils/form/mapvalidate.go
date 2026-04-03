package form

// MapAndValidate decodes GoAdmin form values into T and validates it.
// Use this in any SetInsertFn or SetUpdateFn handler.
func MapAndValidate[T any](values map[string][]string) (*T, error) {
	var result T

	if err := Decoder.Decode(&result, values); err != nil {
		return nil, err
	}

	if err := Validate.Struct(&result); err != nil {
		return nil, err
	}

	return &result, nil
}

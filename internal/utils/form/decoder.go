package form

import (
	"fmt"
	"reflect"
	"time"

	"github.com/go-playground/form/v4"
	"github.com/go-playground/validator/v10"
	"github.com/james-wukong/online-school-mgmt/internal/types"
)

var (
	Decoder  = form.NewDecoder()
	Validate = validator.New()
)

func init() {
	// Register custom type for CivilDate
	Decoder.RegisterCustomTypeFunc(func(vals []string) (interface{}, error) {
		if len(vals) == 0 || vals[0] == "" {
			return types.CivilDate{}, nil
		}
		t, err := time.Parse("2006-01-02", vals[0])
		if err != nil {
			return nil, fmt.Errorf("invalid date format, expected YYYY-MM-DD: %w", err)
		}
		return types.CivilDate(t), nil
	}, types.CivilDate{})

	// Use "form" tag for error messages instead of struct field names
	Validate.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := fld.Tag.Get("form")
		if name == "" || name == "-" {
			return fld.Name
		}
		return name
	})
}

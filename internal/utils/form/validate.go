package form

import (
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/james-wukong/online-school-mgmt/internal/types"
)

func init() {
	// Custom validator: ensures CivilDate is not zero
	if err := Validate.RegisterValidation("valid_date", func(fl validator.FieldLevel) bool {
		d, ok := fl.Field().Interface().(types.CivilDate)
		if !ok {
			return false
		}
		return !time.Time(d).IsZero()
	}); err != nil {
		panic("failed to register valid_date validator: " + err.Error())
	}
}

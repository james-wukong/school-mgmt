package validation

import (
	"github.com/go-playground/validator/v10"
	"github.com/james-wukong/online-school-mgmt/internal/dto"
)

var Validate *validator.Validate

func init() {
	var Validate = validator.New()

	Validate.RegisterStructValidation(
		RequirementValidation,
		dto.RequirementCreateRequest{},
	)
}

func RequirementValidation(sl validator.StructLevel) {
	req := sl.Current().Interface().(dto.RequirementCreateRequest)

	validateNestedIfIDMissing(sl, req.TeacherID, req.Teacher, "Teacher")
	validateNestedIfIDMissing(sl, req.SubjectID, req.Subject, "Subject")
	validateNestedIfIDMissing(sl, req.ClassID, req.Class, "Class")
}

func validateNestedIfIDMissing(sl validator.StructLevel, id int64, obj any, fieldName string) {
	if id != 0 {
		return
	}

	if err := Validate.Struct(obj); err != nil {
		for _, e := range err.(validator.ValidationErrors) {
			sl.ReportError(
				e.Value(),
				fieldName+"."+e.StructField(),
				e.StructField(),
				e.Tag(),
				"",
			)
		}
	}
}

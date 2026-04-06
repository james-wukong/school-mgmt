package dto

import (
	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/types"
)

type TeacherBase struct {
	// School context is usually required for every teacher
	SchoolID int64 `form:"school_id" csv:"school_id" json:"school_id" validate:"-"`

	FirstName string `form:"first_name" csv:"first_name" json:"first_name" validate:"required,min=1,max=100"`
	LastName  string `form:"last_name" csv:"last_name" json:"last_name" validate:"required,min=1,max=100"`

	// Email is optional in your model (*string), but often required in business logic
	Email string `form:"email" csv:"email" json:"email" validate:"omitempty,email"`
	Phone string `form:"phone" csv:"phone" json:"phone" validate:"omitempty,max=20"`

	// Dates from HTML forms come as strings (e.g., "2026-03-27")
	// The form decoder handles the conversion to time.Time if configured,
	// otherwise, use a string and parse it in the logic.
	HireDate types.CivilDate `form:"hire_date" csv:"hire_date" json:"hire_date" validate:"omitempty"`

	// Restricted set of values
	EmploymentType string `form:"employment_type" csv:"employment_type" json:"employment_type" validate:"omitempty,oneof=FullTime PartTime Contract Permanent"`

	// Constraints for scheduling logic
	MaxClassesPerDay int  `form:"max_classes_per_day" csv:"max_classes_per_day" json:"max_classes_per_day" validate:"omitempty,min=1,max=20"`
	IsActive         bool `form:"is_active" csv:"is_active" json:"is_active"` // Booleans don't need 'required' (false is a zero value)
}

type TeacherRelationBase struct {
	// Many-to-Many: GoAdmin sends multi-selects as a slice of strings or IDs
	SubjectIDs  []int64 `form:"subjects[]" csv:"-" json:"-" validate:"omitempty,unique"`
	TimeslotIDs []int64 `form:"timeslots[]" csv:"-" json:"-" validate:"omitempty,unique"`
}

type TeacherCreateRequest struct {
	TeacherBase
	TeacherRelationBase
	// Employee ID: Must be unique, so we ensure it's provided and positive
	EmployeeID int64 `form:"employee_id" csv:"employee_id" json:"employee_id" validate:"omitempty"`
}

type TeacherUpdateRequest struct {
	ID int64 `form:"id" json:"id" csv:"id" validate:"required"` // The ID is mandatory
	TeacherBase
	TeacherRelationBase
	// Employee ID: Must be unique, so we ensure it's provided and positive
	EmployeeID int64 `form:"employee_id" csv:"employee_id" json:"employee_id" validate:"omitempty"`
}

type TeacherStatusUpdateRequest struct {
	ID       int64 `form:"id" validate:"required"` // The ID is mandatory
	IsActive bool  `form:"is_active"`
}

type TeacherBatchCreateRequest struct {
	TeacherBase
	// Employee ID: Must be unique, so we ensure it's provided and positive
	EmployeeID int64 `form:"employee_id" csv:"employee_id" json:"employee_id" validate:"required"`
	// Many-to-Many: GoAdmin sends multi-selects as a slice of strings or IDs
	SubjectIDs types.Int64Slice `form:"subjects[]" csv:"subject_ids" json:"subject_ids" validate:"omitempty,unique"`
}

func (req *TeacherBase) toModel() (*models.Teachers, error) {
	return &models.Teachers{
		SchoolID:         req.SchoolID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            &req.Email, // Mapping string to *string
		Phone:            &req.Phone,
		HireDate:         req.HireDate,
		EmploymentType:   req.EmploymentType,
		MaxClassesPerDay: req.MaxClassesPerDay,
		IsActive:         req.IsActive,
	}, nil
}

func (req *TeacherCreateRequest) ToModel() (*models.Teachers, error) {
	m, err := req.toModel()
	if err != nil {
		return nil, err
	}
	m.EmployeeID = req.EmployeeID
	return m, nil
}

func (req *TeacherUpdateRequest) ToModel() (*models.Teachers, error) {
	m, err := req.toModel()
	if err != nil {
		return nil, err
	}
	m.EmployeeID = req.EmployeeID
	// Manually attach the ID since it's not in the Base
	m.ID = req.ID
	return m, nil
}

func (req *TeacherStatusUpdateRequest) ToModel() *models.Teachers {
	return &models.Teachers{
		ID:       req.ID,
		IsActive: req.IsActive,
	}
}

func (req *TeacherBatchCreateRequest) ToModel() (*models.Teachers, error) {
	m, err := req.toModel()
	if err != nil {
		return nil, err
	}
	if len(req.SubjectIDs) > 0 {
		for _, sub := range req.SubjectIDs {
			m.Subjects = append(m.Subjects, &models.Subjects{ID: sub})
		}
	}
	m.EmployeeID = req.EmployeeID

	return m, nil
}

package dto

import (
	"fmt"
	"time"

	"github.com/james-wukong/online-school-mgmt/internal/models"
)

type TeacherBase struct {
	// School context is usually required for every teacher
	SchoolID int64 `form:"school_id" validate:"required,gt=0"`

	// Employee ID: Must be unique, so we ensure it's provided and positive
	EmployeeID int64 `form:"employee_id" validate:"required,gt=0"`

	FirstName string `form:"first_name" validate:"omitempty,min=2,max=100"`
	LastName  string `form:"last_name" validate:"omitempty,min=2,max=100"`

	// Email is optional in your model (*string), but often required in business logic
	Email string `form:"email" validate:"omitempty,email"`
	Phone string `form:"phone" validate:"omitempty,max=20"`

	// Dates from HTML forms come as strings (e.g., "2026-03-27")
	// The form decoder handles the conversion to time.Time if configured,
	// otherwise, use a string and parse it in the logic.
	HireDate string `form:"hire_date" validate:"omitempty,datetime=2006-01-02"`

	// Restricted set of values
	EmploymentType string `form:"employment_type" validate:"omitempty,oneof=FullTime PartTime Contract Permanent"`

	// Constraints for scheduling logic
	MaxClassesPerDay int  `form:"max_classes_per_day" validate:"omitempty,min=1,max=20"`
	IsActive         bool `form:"is_active"` // Booleans don't need 'required' (false is a zero value)

	// Many-to-Many: GoAdmin sends multi-selects as a slice of strings or IDs
	SubjectIDs  []int64 `form:"subjects[]" validate:"omitempty,unique"`
	TimeslotIDs []int64 `form:"timeslots[]" validate:"omitempty,unique"`
}

type TeacherCreateRequest struct {
	TeacherBase
}

type TeacherUpdateRequest struct {
	ID int64 `form:"id" validate:"required"` // The ID is mandatory
	TeacherBase
}

type TeacherStatusUpdateRequest struct {
	ID       int64 `form:"id" validate:"required"` // The ID is mandatory
	IsActive bool  `form:"is_active"`
}

func (req *TeacherBase) toModel() (*models.Teachers, error) {
	// Handle complex conversions like strings to time.Time
	hireDate, err := time.Parse(models.TimeDateLayout, req.HireDate)
	if err != nil {
		return nil, fmt.Errorf("invalid hire date format: %v", err)
	}

	return &models.Teachers{
		SchoolID:         req.SchoolID,
		EmployeeID:       req.EmployeeID,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Email:            &req.Email, // Mapping string to *string
		Phone:            &req.Phone,
		HireDate:         hireDate,
		EmploymentType:   req.EmploymentType,
		MaxClassesPerDay: req.MaxClassesPerDay,
		IsActive:         req.IsActive,
	}, nil
}

func (req *TeacherCreateRequest) ToModel() (*models.Teachers, error) {
	return req.toModel()
}

func (req *TeacherUpdateRequest) ToModel() (*models.Teachers, error) {
	m, err := req.toModel()
	if err != nil {
		return nil, err
	}
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

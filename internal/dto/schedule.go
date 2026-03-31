package dto

import (
	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/shopspring/decimal"
)

type ScheduleBase struct {
	SchoolID      int64 `form:"school_id" validate:"required"`
	SemesterID    int64 `form:"semester_id" csv:"semester_id" validate:"required"`
	RequirementID int64 `form:"requirement_id" validate:"required"`
	RoomID        int64 `form:"room_id" validate:"required"`
	TimeslotID    int64 `form:"timeslot_id" validate:"required"`

	// Status uses a pointer to handle the "Default" value correctly if not sent
	Status models.ScheduleStatus `form:"status" validate:"oneof=Draft Published Active Archived"`

	// Version is kept as decimal to ensure 1.00 doesn't become 1
	Version decimal.Decimal `form:"version" validate:"required"`

	// relationships
	Requirement models.Requirements `form:"foreignKey:RequirementID" validate:"requirement,omitempty"`
	Semester    models.Semesters    `form:"foreignKey:SemeseterID" csv:"-" validate:"semester,omitempty"`
	Room        models.Rooms        `form:"foreignKey:RoomID" validate:"room,omitempty"`
	Timeslot    models.Timeslots    `form:"foreignKey:TimeslotID" validate:"timeslot,omitempty"`
}

type ScheduleCreateRequest struct {
	ScheduleBase
}

type ScheduleUpdateRequest struct {
	ID int64 `form:"id" validate:"required"` // The ID is mandatory
	ScheduleBase
}

func (req *ScheduleBase) toModel() (*models.Schedules, error) {
	return &models.Schedules{
		SchoolID:      req.SchoolID,
		RequirementID: req.RequirementID,
		RoomID:        req.RoomID,
		TimeslotID:    req.TimeslotID,
		Status:        req.Status,
		Version:       req.Version,
	}, nil
}

func (req *ScheduleCreateRequest) ToModel() (*models.Schedules, error) {
	return req.toModel()
}

func (req *ScheduleUpdateRequest) ToModel() (*models.Schedules, error) {
	m, err := req.toModel()
	if err != nil {
		return nil, err
	}
	// Manually attach the ID since it's not in the Base
	m.ID = req.ID
	return m, nil
}

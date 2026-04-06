// Package dto
package dto

import (
	"time"

	"github.com/james-wukong/online-school-mgmt/internal/models"
	"github.com/james-wukong/online-school-mgmt/internal/types"
)

type TimeslotsBase struct {
	SchoolID   int64            `form:"school_id" json:"school_id" csv:"school_id" validate:"-"`
	SemesterID int64            `form:"semester_id" json:"semester_id" csv:"semester_id" validate:"-"`
	Day        models.DayOfWeek `form:"day" json:"day" csv:"day" validate:"required"`
	StartTime  types.ClockTime  `form:"start_time" json:"start_time" csv:"start_time" validate:"required"`
	EndTime    types.ClockTime  `form:"end_time" json:"end_time" csv:"end_time" validate:"required"`
}

type TimeslotsCreateRequest struct {
	TimeslotsBase
}

type TimeslotsUpdateRequest struct {
	ID int64 `form:"id" json:"id" csv:"id" validate:"required"` // The ID is mandatory
	TimeslotsBase
}

func (s *TimeslotsBase) toModel() (*models.Timeslots, error) {
	return &models.Timeslots{
		SchoolID:   s.SchoolID,
		SemesterID: s.SemesterID,
		DayOfWeek:  s.Day,
		StartTime:  time.Time(s.StartTime),
		EndTime:    time.Time(s.EndTime),
	}, nil
}

func (s *TimeslotsCreateRequest) ToModel() (*models.Timeslots, error) {
	return s.toModel()
}

func (s *TimeslotsUpdateRequest) ToModel() (*models.Timeslots, error) {
	m, err := s.toModel()
	if err != nil {
		return nil, err
	}
	m.ID = s.ID
	return m, nil
}

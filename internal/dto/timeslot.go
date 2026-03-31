// Package dto
package dto

import (
	"fmt"
	"strings"
	"time"

	"github.com/james-wukong/online-school-mgmt/internal/models"
)

type HourMinute time.Time

type ClassTimes struct {
	StartTime HourMinute `json:"start_time"`
	EndTime   HourMinute `json:"end_time"`
}

// UnmarshalJSON handles the "09:00" -> HourMinute conversion
// When json.Unmarshal encounters a field of type HourMinute,
// it checks if that type has an UnmarshalJSON([]byte) error method
func (hm *HourMinute) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")
	if s == "" {
		return nil
	}
	t, err := time.Parse("15:04", s)
	if err != nil {
		return err
	}
	*hm = HourMinute(t)
	return nil
}

// Schedule maps the JSON where keys are days of the week
type Schedule map[string][]ClassTimes

// Validate checks for overlapping time slots within the same day.
func (s Schedule) Validate() error {
	for day, slots := range s {
		for i := 0; i < len(slots); i++ {
			// if end time is earlier than start time, return error
			if !time.Time(slots[i].EndTime).After(time.Time(slots[i].StartTime)) {
				return fmt.Errorf(
					"conflict detected on %s at period: %s - %s",
					day, time.Time(slots[i].StartTime), time.Time(slots[i].EndTime),
				)
			}
			for j := i + 1; j < len(slots); j++ {
				// Simple logic: if Start of one is between Start/End of another
				// This assumes you've converted strings to integers for comparison
				if time.Time(slots[i].StartTime).Equal(time.Time(slots[j].StartTime)) {
					return fmt.Errorf(
						"conflict detected on %s at period %s",
						day, time.Time(slots[i].StartTime),
					)
				}
			}
		}
	}
	return nil
}

func (s Schedule) MapToTimeslots(semesterID, schoolID int64) []*models.Timeslots {
	var t []*models.Timeslots
	var day models.DayOfWeek
	for key, value := range s {
		if len(value) == 0 {
			continue
		}

		day = ParseDay(key)
		for _, slot := range value {
			t = append(t, &models.Timeslots{
				SchoolID:   schoolID,
				SemesterID: semesterID,
				DayOfWeek:  day,
				StartTime:  time.Time(slot.StartTime),
				EndTime:    time.Time(slot.EndTime),
			})
		}
	}
	return t
}

func ParseDay(key string) models.DayOfWeek {
	var dayMap = map[string]models.DayOfWeek{
		"monday":    models.Monday,
		"tuesday":   models.Tuesday,
		"wednesday": models.Wednesday,
		"thursday":  models.Thursday,
		"friday":    models.Friday,
		"saturday":  models.Saturday,
		"sunday":    models.Sunday,
	}
	normalized := strings.ToLower(key)

	if day, ok := dayMap[normalized]; ok {
		return day
	}

	return models.Monday
}

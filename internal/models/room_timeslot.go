package models

type RoomTimeslots struct {
	// Foreign Key to Rooms
	RoomID int64 `gorm:"not null;primaryKey" json:"room_id"`
	// Foreign Key to Timeslots
	TimeslotID int64 `gorm:"not null;primaryKey" json:"timeslot_id"`
}

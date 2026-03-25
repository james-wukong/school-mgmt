// Package models defines the subjects entity and related value objects.
// It represents how data looks in the database or business rules.
package models

// Subjects represents the subjects table in PostgreSQL.
// It uses GORM tags to handle identity columns and automatic timestamps.
type Subjects struct {
	// id is GENERATED ALWAYS. We use <-:false to prevent GORM from
	// including it in INSERT or UPDATE statements.
	ID       int64 `gorm:"column:id;primaryKey;<-:false" json:"id"`
	SchoolID int64 `gorm:"column:school_id;not null" json:"school_id"`

	Name        string `gorm:"column:name;not null;unique" json:"name"`
	Code        string `gorm:"column:code;not null;unique" json:"code"`
	Description string `gorm:"column:description" json:"description"`

	RequiresLab bool `gorm:"column:requires_lab;default:false" json:"requires_lab"`
	IsHeavy     bool `gorm:"column:is_heavy;default:false" json:"is_heavy"`
	// Relationships (Optional but recommended for Eager Loading)
	School *Schools `gorm:"foreignKey:SchoolID;constraint:OnDelete:CASCADE;" json:"-"`
}

func NewSubjects(schoolID int64,
	name, code string,
	requiresLab, isHeavy bool,
) *Subjects {
	return &Subjects{
		SchoolID:    schoolID,
		Name:        name,
		Code:        code,
		RequiresLab: requiresLab,
		IsHeavy:     isHeavy,
	}
}

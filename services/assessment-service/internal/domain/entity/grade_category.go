package entity

import (
	"time"

	"github.com/google/uuid"
)

type GradeCategory struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Weight      int       `json:"weight"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (GradeCategory) TableName() string {
	return "grade_categories"
}

func (g *GradeCategory) Validate() map[string]string {
	errors := make(map[string]string)
	if g.Name == "" {
		errors["name"] = "Name is required"
	}
	if g.Weight < 0 {
		errors["weight"] = "Weight cannot be negative"
	}
	return errors
}

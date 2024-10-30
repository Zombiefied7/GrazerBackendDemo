// models.go
package main

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

// Custom type for handling []string as JSON
type StringArray []string

// Value marshals the []string to JSON for storing in the database
func (s StringArray) Value() (driver.Value, error) {
	return json.Marshal(s)
}

// Scan unmarshals JSON from the database to a []string
func (s *StringArray) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &s)
}

// Preferences represents user matching preferences
type Preferences struct {
	ID        uint        `gorm:"primaryKey"`
	UserID    uint        `gorm:"uniqueIndex"` // Each user has only one preferences entry
	Genders   StringArray `gorm:"type:json"`   // Acceptable genders stored as JSON
	MinAge    int
	MaxAge    int
	DietTypes StringArray `gorm:"type:json"` // Multiple diet preferences stored as JSON
	Distance  int
	CreatedAt time.Time
	UpdatedAt time.Time
}

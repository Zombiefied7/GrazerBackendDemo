package main

import (
	"time"
)

// User represents a user profile with attributes needed for matching
type User struct {
	ID        uint   `gorm:"primaryKey"`
	Name      string `gorm:"size:100"`
	Gender    string `gorm:"size:10"` // "male", "female", or "other"
	Location  string `gorm:"size:50"` // Coordinates as "latitude,longitude"
	DietType  string `gorm:"size:20"` // "vegan", "vegetarian"
	Age       int    `gorm:"size:3"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

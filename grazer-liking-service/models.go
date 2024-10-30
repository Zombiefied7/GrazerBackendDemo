// models.go
package main

import (
	"time"
)

// Like represents a like relationship between two users
type Like struct {
	ID          uint `gorm:"primaryKey"`
	UserID      uint // The user who liked
	LikedUserID uint // The user being liked
	CreatedAt   time.Time
}

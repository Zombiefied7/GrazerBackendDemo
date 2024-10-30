// models.go
package main

import (
	"time"
)

// Like represents a like relationship between two users (for incoming events)
type Like struct {
	UserID      uint      `json:"userID"`
	LikedUserID uint      `json:"likedUserID"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Match represents a mutual match between two users
type Match struct {
	ID            uint `gorm:"primaryKey"`
	UserID        uint // One of the matched users
	MatchedUserID uint // The other matched user
	CreatedAt     time.Time
}

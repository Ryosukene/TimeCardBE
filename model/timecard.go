package model

import (
	"time"
)

type AttendanceRecord struct {
	ID           uint      `json:"id" gorm:"primaryKey"`
	UserID       uint      `json:"user_id" gorm:"not null"`
	ClockInTime  time.Time `json:"clock_in_time"`
	ClockOutTime time.Time `json:"clock_out_time"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	User         User      `json:"user" gorm:"foreignKey:UserID; constraint:OnDelete:CASCADE"`
}

type AttendanceRecordResponse struct {
	ID           uint         `json:"id"`
	UserID       uint         `json:"user_id"`
	ClockInTime  time.Time    `json:"clock_in_time"`
	ClockOutTime time.Time    `json:"clock_out_time"`
	CreatedAt    time.Time    `json:"created_at"`
	UpdatedAt    time.Time    `json:"updated_at"`
	User         UserResponse `json:"user"`
}

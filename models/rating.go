package models

import (
	"time"
)

type Rating struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	RideID    uint      `gorm:"not null" json:"ride_id"`
	Score     int       `gorm:"type:int" json:"score"` // 1-5
	Comment   string    `gorm:"type:text" json:"comment"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`

	Ride Ride `gorm:"foreignKey:RideID" json:"-"`
}

func (Rating) TableName() string {
	return "ratings"
}

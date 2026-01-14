package models

import (
	"time"

	"gorm.io/gorm"
)

type Client struct {
	ID          uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID      uint           `gorm:"not null;uniqueIndex" json:"user_id"`
	FullName    string         `gorm:"type:varchar(120);not null" json:"full_name"`
	Email       string         `gorm:"type:varchar(100);not null" json:"email"`
	Phone       string         `gorm:"type:varchar(20);not null" json:"phone"`
	Preferences string         `gorm:"type:json" json:"preferences"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`

	// Esta línea es la que fallaba porque no encontraba "User"
	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (Client) TableName() string {
	return "clients"
}

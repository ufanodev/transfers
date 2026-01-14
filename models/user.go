package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	Username     string         `gorm:"type:varchar(50);not null;uniqueIndex" json:"username"`
	Email        string         `gorm:"type:varchar(100);not null;uniqueIndex" json:"email"`
	PasswordHash string         `gorm:"column:password_hash;type:varchar(255);not null" json:"-"`
	Role         string         `gorm:"type:enum('client', 'company', 'admin', 'driver');not null;default:'client';index" json:"role"`
	Phone        string         `gorm:"type:varchar(20)" json:"phone"`
	IsActive     bool           `gorm:"column:is_active;default:true" json:"is_active"`
	ResetToken   string         `gorm:"column:reset_token;type:varchar(255)" json:"-"`
	ResetExpires *time.Time     `gorm:"column:reset_expires" json:"-"`
	CreatedAt    time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (User) TableName() string {
	return "users"
}

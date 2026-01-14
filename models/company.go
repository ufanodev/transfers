package models

import (
	"time"

	"gorm.io/gorm"
)

type Company struct {
	ID         uint           `gorm:"primaryKey;autoIncrement" json:"id"`
	UserID     uint           `gorm:"not null;uniqueIndex" json:"user_id"`
	Name       string         `gorm:"type:varchar(100);not null" json:"name"`
	TaxID      string         `gorm:"type:varchar(20);not null;uniqueIndex" json:"tax_id"`
	Address    string         `gorm:"type:text" json:"address"`
	PostalCode string         `gorm:"type:varchar(10)" json:"postal_code"`
	Website    string         `gorm:"type:varchar(100)" json:"website"`
	CreatedAt  time.Time      `json:"created_at"`
	UpdatedAt  time.Time      `json:"updated_at"`
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`

	User User `gorm:"foreignKey:UserID" json:"-"`
}

func (Company) TableName() string {
	return "companies"
}

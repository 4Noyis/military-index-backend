package models

import "time"

// Country represents a country entity in the database
type Country struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name      string    `gorm:"type:varchar(100);not null;unique" json:"name" binding:"required"`
	Code      string    `gorm:"type:varchar(3);not null;unique" json:"code" binding:"required"`
	FlagURL   string    `gorm:"type:varchar(255)" json:"flag_url"`
	CreatedAt time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

// TableName specifies the table name for the Country model
func (Country) TableName() string {
	return "countries"
}

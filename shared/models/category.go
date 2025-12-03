package models

import "time"

// Category represents a technology category entity in the database
type Category struct {
	ID          uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name        string    `gorm:"type:varchar(100);not null;unique" json:"name" binding:"required"`
	Description string    `gorm:"type:text" json:"description"`
	IconURL     string    `gorm:"type:varchar(255)" json:"icon_url"`
	CreatedAt   time.Time `gorm:"type:autoCreateTime" json:"created_at"`
}

// TableName specifies the table name for the Category Model
func (Category) TableName() string {
	return "tech_categories"
}

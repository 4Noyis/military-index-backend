package models

import "time"

// Technology represents a military technology entity in the database
type Technology struct {
	ID            uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	CountryID     uint      `gorm:"not null;index" json:"country_id" binding:"required"`
	CategoryID    uint      `gorm:"not null;index" json:"category_id" binding:"required"`
	Name          string    `gorm:"type:varchar(200);not null" json:"name" binding:"required"`
	Description   string    `gorm:"type:text" json:"description"`
	Designer      string    `gorm:"type:varchar(200)" json:"designer"`
	YearDeveloped *int      `gorm:"type:integer;index" json:"year_developed"`
	YearDeployed  *int      `gorm:"type:integer;index" json:"year_deployed"`
	Manufacturer  string    `gorm:"type:varchar(250)" json:"manufacturer"`
	UnitCost      *int64    `gorm:"type:bigint" json:"unit_cost"`
	Mass          *int      `gorm:"type:integer" json:"mass"`
	Length        *int      `gorm:"type:integer" json:"length"`
	Width         *int      `gorm:"type:integer" json:"width"`
	Height        *int      `gorm:"type:integer" json:"height"`
	Status        string    `gorm:"type:varchar(50)" json:"status"`
	ImageURL      string    `gorm:"type:varchar(255)" json:"image_url"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`

	// Relationships
	Country  Country  `gorm:"foreignKey:CountryID;constraint:OnDelete:CASCADE" json:"country,omitempty"`
	Category Category `gorm:"foreignKey:CategoryID;constraint:OnDelete:RESTRICT" json:"category,omitempty"`
}

// TableName specifies the table name for the Technology Model
func (Technology) TableName() string {
	return "technologies"
}

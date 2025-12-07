package models

type AircraftSpecs struct {
	ID          uint   `gorm:"primaryKey" json:"id"`
	EquipmentID uint   `gorm:"not null;uniqueIndex" json:"equipment_id"`

	Type        string `gorm:"size:100" json:"type"`           // fighter, bomber, cargo, helicopter...
	Crew        int    `json:"crew"`
	Wingspan    string `gorm:"size:50" json:"wingspan"`
	MaxSpeed    string `gorm:"size:100" json:"max_speed"`
	RangeKm     string `gorm:"size:100" json:"range_km"`
	ServiceCeiling string `gorm:"size:100" json:"service_ceiling"`
	EngineType  string `gorm:"size:150" json:"engine_type"`
	EngineCount int    `json:"engine_count"`
	Armament    string `gorm:"type:text" json:"armament"`

	Technology   Technology `gorm:"foreignKey:ID" json:"-"`
}
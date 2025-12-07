package model

import "time"


type GroundVehicle struct {
    TechnologyID uint `gorm:"primaryKey" json:"technology_id"`

    EnginePower *int   `json:"engine_power"` // hp
    ArmorLevel  string `gorm:"type:varchar(150)" json:"armor_level"`
    MaxSpeed    *int   `json:"max_speed"`    // km/h
    Crew        *int   `json:"crew"`
    Range       *int   `json:"range"`        // km

    Technology Technology `gorm:"foreignKey:TechnologyID;references:ID;constraint:OnDelete:CASCADE"`
}

func (GroundVehicle) TableName() string {
    return "ground_vehicle_specs"
}

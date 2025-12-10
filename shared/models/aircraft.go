package models

type Aircraft struct {
    TechnologyID uint    `gorm:"primaryKey" json:"technology_id"`

    MaxSpeed     *int    `json:"max_speed"`        // km/h
    Ceiling      *int    `json:"ceiling"`          // max irtifa (m)
    Range        *int    `json:"range"`            // menzil (km)
    Crew         *int    `json:"crew"`             // mürettebat
    EngineType   string  `gorm:"type:varchar(150)" json:"engine_type"`
    Wingspan     *float64 `json:"wingspan"`        // kanat açıklığı (m)    

    Technology Technology `gorm:"foreignKey:TechnologyID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Aircraft) TableName() string {
    return "aircraft_specs"
}
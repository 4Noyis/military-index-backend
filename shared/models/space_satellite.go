package model

type Space_satellite struct {
    TechnologyID uint `gorm:"primaryKey"`

    OrbitType   string
    MissionType string
    Payload     string

    Technology Technology `gorm:"foreignKey:TechnologyID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Satellite) TableName() string {
    return "satellite_specs"
}

package model

type Drone struct {
    TechnologyID uint `gorm:"primaryKey"`

    Endurance   *int
    Payload     *int
    ControlType string
    Range       *int

    Technology Technology `gorm:"foreignKey:TechnologyID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Drone) TableName() string {
    return "drone_specs"
}

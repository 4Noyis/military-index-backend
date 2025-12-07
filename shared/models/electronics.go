package model

type Electronics struct {
    TechnologyID uint `gorm:"primaryKey"`

    SystemType string
    PowerUsage *int
    Frequency  string

    Technology Technology `gorm:"foreignKey:TechnologyID;references:ID;constraint:OnDelete:CASCADE"`
}

func (Electronics) TableName() string {
    return "electronics_specs"
}

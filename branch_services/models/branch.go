package models

type Branch struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string `gorm:"type:varchar(100);not null" json:"name"`
	IsActive bool   `gorm:"not null;default:false" json:"is_active"`
}

func (Branch) TableName() string {
	return "branches"
}
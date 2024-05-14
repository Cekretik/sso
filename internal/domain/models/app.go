package models

type App struct {
	ID     int    `gorm:"column:id" json:"id"`
	Name   string `gorm:"column:name" json:"name"`
	Secret string `gorm:"column:secret" json:"secret"`
}

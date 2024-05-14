package models

type User struct {
	ID       int    `gorm:"column:id" json:"id"`
	Email    string `gorm:"column:email" json:"email"`
	PassHash []byte `gorm:"column:pass_hash" json:"pass_hash"`
	IsAdmin  bool   `gorm:"column:is_admin" json:"is_admin"`
}

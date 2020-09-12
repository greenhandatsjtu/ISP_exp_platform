package models

type Role struct {
	ID          uint    `json:"id"`
	Description string  `gorm:"type:varchar(100);unique;not_null" json:"description"`
	Users       []*User `gorm:"many2many:user_roles" json:"users"`
}

package model

import "gorm.io/gorm"

type Group struct {
	gorm.Model

	Name    string  `json:"group_name" gorm:"column:name;size:255;uniqueIndex"`
	Members []*User `json:"members" gorm:"many2many:group_members;"`
}

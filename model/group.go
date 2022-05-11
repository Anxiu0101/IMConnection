package model

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model

	Name    string  `json:"group_name" gorm:"column:name;size:255;uniqueIndex"`
	Members []*User `json:"members" gorm:"many2many:group_members;"`
}

// GetUserList 获取组用户列表
func GetUserList(name string, pageNum, pageSize int) {
	var data []User
	err := DB.Select("Members").Model(Group{}).
		Where("name = ?", name).
		Offset(pageNum).Limit(pageSize).
		Find(&data).Error
	if err != nil {

	}

}

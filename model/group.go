package model

import (
	"gorm.io/gorm"
)

type Group struct {
	gorm.Model

	Name    string  `json:"group_name" gorm:"column:name;size:255;uniqueIndex"`
	Members []*User `json:"members" gorm:"many2many:group_members;"`
}

type GroupMembers struct {
	GroupId uint `json:"group_id" gorm:"comment:GroupID"`
	UserId  uint `json:"user_id" gorm:"comment:UserID"`
}

// GetUserList 获取组用户列表
//func GetUserList(name string, pageSize int) map[string]interface{} {
//	var data []User
//
//	err := DB.Select("Members").Model(Group{}).
//		Where("name = ?", name).
//		Limit(pageSize).Find(&data).Error
//	if err != nil {
//		log.Printf("Error: %s", err)
//	}
//
//	var group Group
//	gid := DB.Model(Group{}).Select("id").Where("name = ?", name)
//
//	var result map[string]interface{}
//	DB.Table("imc_group_members").Select("user_id").Where("group_id = ?")
//
//	return data
//}

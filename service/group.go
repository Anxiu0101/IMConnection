package service

import (
	"IMConnection/model"
	"IMConnection/pkg/e"
	"fmt"
)

type GroupService struct {
	Name string `json:"name" form:"name"`
}

// Create 创建组
// 1. 查询是否重名
// 2. 在数据库创建组
// 3. 返回成功结果
func (service *GroupService) Create() model.Response {
	code := e.Success
	var group model.Group

	if err := model.DB.Select("name").Where("name = ?", service.Name).First(&group).Error; err != nil {
		code = e.InvalidParams
		return model.Response{
			Code: code,
			Msg:  e.GetMsg(code),
			Data: "组名冲突，请使用其他组名",
		}
	}

	group.Name = service.Name
	if err := model.DB.Create(&service); err != nil {
		code = e.Error
		return model.Response{
			Code: code,
			Msg:  e.GetMsg(code),
			Data: "数据库异常",
		}
	}

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: fmt.Sprintf("successful create %s", service.Name),
	}
}

// Invite 邀请群成员
// 1. 查询用户是否存在
// 2. 查询
// TODO group member invite service haven't finish
func (service *GroupService) Invite(uid uint) model.Response {
	code := e.Success
	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: fmt.Sprintf("successful create %s", service.Name),
	}
}

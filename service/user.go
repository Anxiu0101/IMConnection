package service

import (
	"IMConnection/model"
	"IMConnection/pkg/e"
)

// AccountService 账户服务，关于用户的登陆和注册
type AccountService struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
}

func (service *AccountService) Login() model.Response {
	code := e.Success

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: "Successful login",
	}
}

func (service *AccountService) Register() model.Response {
	code := e.Success

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: "Successful login",
	}
}

func (service *AccountService) SetPassword() model.Response {
	code := e.Success

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: "Successful login",
	}
}

// UserService 用户服务，关于用户的信息以及好友关系等。
type UserService struct {
	Email  string `json:"email" form:"email"`
	Gender int    `json:"gender" form:"gender"`
	Age    int    `json:"age" form:"age"`
	Tel    int    `json:"tel" form:"tel"`
}

func (service *UserService) GetInfo() model.Response {
	code := e.Success

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: "Successful login",
	}
}

func (service *UserService) UpdateInfo() model.Response {
	code := e.Success

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: "Successful login",
	}
}

func (service *UserService) GetRelationList() model.Response {
	code := e.Success

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: "Successful login",
	}
}

func (service *UserService) GetRelationByGroup(group string) model.Response {
	code := e.Success

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: "Successful login",
	}
}

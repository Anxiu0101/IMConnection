package service

import (
	"IMConnection/model"
	"IMConnection/pkg/e"
	"IMConnection/pkg/logging"
	"IMConnection/pkg/util"
)

// AccountService 账户服务，关于用户的登陆和注册
type AccountService struct {
	Username string `json:"username" form:"username"`
	Password string `json:"password" form:"password"`
	Avatars  string `json:"avatars"  form:"avatars"`
}

// Register
// 1. 先查询用户名是否已存在
// 2. 为新用户设置密码
// 3. 在数据库中创建新用户
func (service *AccountService) Register() model.Response {
	code := e.Success
	var user model.User
	var count int64

	// 查询用户是否已存在
	model.DB.Model(&model.User{}).Where("username = ?", service.Username).Find(&user).Count(&count)
	if count > 0 {
		code = e.InvalidParams
		return model.Response{
			Code: code,
			Msg:  e.GetMsg(code),
			Data: "用户已存在",
		}
	}

	// 为新用户设置密码
	user.UserName = service.Username
	if err := user.SetPassword(service.Password); err != nil {
		logging.Info(err)
		code = e.Error
		return model.Response{
			Code: code,
			Msg:  e.GetMsg(code),
			Data: "密码加密失败",
		}
	}

	// 创建新用户
	if err := model.DB.Create(&user).Error; err != nil {
		logging.Info(err)
		code = e.Error
		return model.Response{
			Code: code,
			Msg:  e.GetMsg(code),
			Data: "数据库错误",
		}
	}

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: "Successful Register",
	}
}

// Login
// 1. 查询用户
// 2. 验证用户密码
// 3. 生成 token
// 4. 返回结果
func (service *AccountService) Login() model.Response {
	code := e.Success

	// 查询用户是否存在
	// 错误情况：用户被封禁，用户不存在，数据库错误
	var user model.User
	if err := model.DB.Where("username = ? AND state = 0", service.Username).Find(&user).Error; err != nil {
		code = e.InvalidParams
		logging.Info(err)
		return model.Response{
			Code: code,
			Msg:  e.GetMsg(code),
			Data: "用户查询错误",
		}
	}

	// 验证用户密码是否正确
	if !user.CheckPassword(service.Password) {
		code = e.InvalidParams
		return model.Response{
			Code: code,
			Msg:  e.GetMsg(code),
			Data: "用户密码错误",
		}
	}

	// 修改为使用 refresh token 和 access token 组成的 token 对
	// TODO
	token, err := util.GenerateToken(user.ID, service.Username, 0)
	if err != nil {
		logging.Info(err)
		code = e.Error
		return model.Response{
			Code: code,
			Msg:  e.GetMsg(code),
		}
	}

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: model.TokenData{User: model.BuildUser(user), AccessToken: token},
	}
}

// ResetPassword
// 1. 查询用户
// 2. 设置用户密码
// 3. 返回结果
func (service *AccountService) ResetPassword(id uint, newPassword string) model.Response {
	code := e.Success

	// 查询用户是否存在
	var user model.User
	if err := model.DB.Where("id = ?", id).Find(&user).Error; err != nil {
		code = e.InvalidParams
		logging.Info(err)
		return model.Response{
			Code: code,
			Msg:  e.GetMsg(code),
			Data: "用户不存在",
		}
	}

	// 重设密码
	if err := user.SetPassword(newPassword); err != nil {
		logging.Info(err)
		code = e.Error
		return model.Response{
			Code: code,
			Msg:  e.GetMsg(code),
			Data: "密码设置失败",
		}
	}

	// 返回结果
	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: "密码已成功修改",
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

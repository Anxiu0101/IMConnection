package service

import (
	"IMConnection/cache"
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
// 4. 将用户 ID 加载入内存
// 5. 返回结果
func (service *AccountService) Login() model.Response {
	code := e.Success

	// 查询用户是否存在
	// 错误情况：用户被封禁，用户不存在，数据库错误
	var user model.User
	if err := model.DB.Where("username = ? AND state = TRUE", service.Username).Find(&user).Error; err != nil {
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

	// TODO 登陆时将用户加载入内存
	if u1, _ := cache.RedisClient.Get(cache.Ctx, string(user.ID)).Result(); u1 == "0" {
		// 不设置过期
		cache.RedisClient.Set(cache.Ctx, string(user.ID), 1, 0)
	}

	// 将用户好友关系载入缓存，若加载失败则退出并返回 加载用户关系失败
	// TODO 修改其使得用户关系加载失败也可以正常登录
	if err := model.LoadRelation(); err != nil {
		logging.Info(err)
		code = e.Error
		return model.Response{
			Code: code,
			Msg:  e.GetMsg(code),
			Data: "加载用户关系失败",
		}
	}

	// 修改为使用 refresh token 和 access token 组成的 token 对
	accessToken, refreshToken, err := util.GenerateTokenPair(user.ID, service.Username, 0)
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
		Data: model.TokenData{
			User:         model.BuildAccountInfo(user),
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		},
	}
}

// ResetPassword
// 1. 查询用户
// 2. 设置用户密码
// 3. 删除 refresh token
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

// GetInfo 获取用户信息
// 1. 检查用户
// 2. 获取用户信息
// 3. 返回结果
func (service *UserService) GetInfo(uid uint) model.Response {
	code := e.Success

	var user model.User
	if err := model.DB.Model(model.User{}).First(&user).Error; err != nil {
		code = e.Error
		logging.Info(err)
		return model.Response{
			Code: code,
			Msg:  e.GetMsg(code),
			Data: "用户不存在",
		}
	}

	var userInfo model.UserInfo
	if err := model.DB.Model(user).Find(&userInfo).Error; err != nil {
		code = e.Error
		logging.Info(err)
		return model.Response{
			Code: code,
			Msg:  e.GetMsg(code),
			Data: "查询用户信息失败",
		}
	}

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: userInfo,
	}
}

// UpdateInfo 更新用户信息
// 1. 检查用户
// 2. 更新用户信息
func (service *UserService) UpdateInfo() model.Response {
	code := e.Success

	// TODO UpdateInfo Service

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: "Successful login",
	}
}

// GetRelationList 获取用户关系列表
// 1. 检查用户
// 2. 从 redis 中读取用户关系列表
// 	3. 若 redis 中为空，则向 SQL 数据库中发送请求获取 relation
// 4. 返回结果
func (service *UserService) GetRelationList() model.Response {
	code := e.Success

	// TODO GetRelationList Service

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: "Successful login",
	}
}

func (service *UserService) GetRelationByGroup(group string) model.Response {
	code := e.Success

	// TODO GetRelationByGroup

	return model.Response{
		Code: code,
		Msg:  e.GetMsg(code),
		Data: "Successful login",
	}
}

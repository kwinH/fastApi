package service

import (
	"fastApi/app/http/middleware"
	"fastApi/app/http/request"
	"fastApi/app/http/response"
	"fastApi/app/model"
	"fastApi/core/global"
	"github.com/gin-gonic/gin"
)

// UserService 管理用户登录的服务
type UserService struct {
}

// Login 用户登录函数
func (service *UserService) Login(c *gin.Context, loginRequest request.LoginRequest) response.Response {
	var userModel model.User

	if err := global.DB.Where("user_name = ?", loginRequest.UserName).First(&userModel).Error; err != nil {
		return response.ParamErr("账号或密码错误", nil, nil)
	}

	if userModel.CheckPassword(loginRequest.Password) == false {
		return response.ParamErr("账号或密码错误", nil, nil)
	}

	token, err := middleware.GenerateToken(userModel.ID)
	if err != nil {
		return response.ParamErr("token生成失败", nil, nil)
	}

	return response.BuildLoginResponse(userModel, token)
}

// valid 验证表单
func (service *UserService) valid(registerRequest request.RegisterRequest) *response.Response {
	if registerRequest.PasswordConfirm != registerRequest.Password {
		return &response.Response{
			Code: 40001,
			Msg:  "两次输入的密码不相同",
		}
	}

	count := int64(0)
	global.DB.Model(&model.User{}).Where("nickname = ?", registerRequest.Nickname).Count(&count)
	if count > 0 {
		return &response.Response{
			Code: 40001,
			Msg:  "昵称被占用",
		}
	}

	count = 0
	global.DB.Model(&model.User{}).Where("user_name = ?", registerRequest.UserName).Count(&count)
	if count > 0 {
		return &response.Response{
			Code: 40001,
			Msg:  "用户名已经注册",
		}
	}

	return nil
}

// Register 用户注册
func (service *UserService) Register(registerRequest request.RegisterRequest) response.Response {
	user := model.User{
		Nickname: registerRequest.Nickname,
		UserName: registerRequest.UserName,
		Status:   model.Active,
	}

	// 表单验证
	if err := service.valid(registerRequest); err != nil {
		return *err
	}

	// 加密密码
	if err := user.SetPassword(registerRequest.Password); err != nil {
		return response.Err(
			response.CodeEncryptError,
			"密码加密失败",
			nil,
			err,
		)
	}

	// 创建用户
	if err := global.DB.Create(&user).Error; err != nil {
		return response.ParamErr("注册失败", nil, err)
	}

	token, err := middleware.GenerateToken(user.ID)
	if err != nil {
		return response.ParamErr("token生成失败", nil, nil)
	}

	return response.BuildLoginResponse(user, token)
}

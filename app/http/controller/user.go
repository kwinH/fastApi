package controller

import (
	"fastApi/app/http/request"
	"fastApi/app/http/response"
	"fastApi/app/http/service"
	"github.com/gin-gonic/gin"
)

type UserController struct {
	Api
}

// UserRegister 用户注册接口
// @Summary     用户注册接口
// @Description 用户注册接口
// @Tags        用户相关接口
// @Accept      application/json
// @Produce     application/json
// @Param       Authorization header string                  true "用户令牌"
// @Param       object        query  request.RegisterRequest false "查询参数"
// @Success     200 {object} response.Response
// @Router      /user/register [post] [get]
func (a UserController) UserRegister(c *gin.Context) (res interface{}, err error) {
	var param request.RegisterRequest
	var service service.UserService
	if err = c.ShouldBind(&param); err != nil {
		return
	}

	res = service.Register(param)
	return
}

// UserLogin 用户登录接口
// @Summary     用户登录接口1
// @Description 用户登录接口2
// @Tags        用户相关接口
// @Accept      application/json
// @Produce     application/json
// @Param       object query request.LoginRequest false "查询参数"
// @Success     200 {object} response.UserResponse
// @Router      /user/login [post]
func (a UserController) UserLogin(c *gin.Context) (res interface{}, err error) {
	var param = request.LoginRequest{}
	var service service.UserService
	if err = c.ShouldBind(&param); err != nil {
		return
	}

	res = service.Login(c, param)
	return
}

// UserMe 用户详情
func (a UserController) UserMe(c *gin.Context) (res interface{}, err error) {
	user := a.currentUser(c)
	res = response.BuildUserResponse(*user)
	return
}

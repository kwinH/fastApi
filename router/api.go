package router

import (
	"fastApi/app/http/controller"
	"fastApi/app/http/middleware"
	"github.com/gin-gonic/gin"
)

func apiRoute(r *gin.Engine) {
	// 路由
	v1 := r.Group("/api/v1")
	{
		UserController := controller.UserController{}

		// 用户登录
		v1.POST("user/register", wrap(UserController.UserRegister))

		// 用户登录
		v1.POST("user/login", wrap(UserController.UserLogin))

		// 需要登录保护的
		auth := v1.Group("")
		auth.Use(middleware.JwtAuth())
		{
			// User Routing
			auth.GET("user/me", wrap(UserController.UserMe))
		}
	}
}

package controller

import (
	"fastApi/app/model"
	"github.com/gin-gonic/gin"
)

type Api struct {

}

// CurrentUser 获取当前用户
func (a Api)currentUser(c *gin.Context) *model.User {
	if userId, _ := c.Get("userId"); userId != nil {
		user, err := model.GetUser(userId)
		if err == nil {
			return &user
		}
	}
	return nil
}
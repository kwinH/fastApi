package response

import (
	"fastApi/app/model"
)

// User 用户序列化器
type User struct {
	ID        uint   `json:"id"`         //ID
	UserName  string `json:"user_name"`  //用户名
	Nickname  string `json:"nickname"`   //昵称
	Status    string `json:"status"`     //状态
	Avatar    string `json:"avatar"`     //头像
	CreatedAt int64  `json:"created_at"` //注册时间
}

type Data struct {
	User  User   `json:"userInfo"`
	Token string `json:"token"`
}

type UserResponse struct {
	Response
	Data User
}

// BuildUserResponse 序列化用户响应
func BuildUserResponse(user model.User) UserResponse {
	return UserResponse{
		Data: User{
			ID:        user.ID,
			UserName:  user.UserName,
			Nickname:  user.Nickname,
			Status:    user.Status,
			Avatar:    user.Avatar,
			CreatedAt: user.CreatedAt.Unix(),
		},
	}
}

func BuildLoginResponse(user model.User, token string) Response {
	return Response{
		Data: Data{
			User: User{
				ID:        user.ID,
				UserName:  user.UserName,
				Nickname:  user.Nickname,
				Status:    user.Status,
				Avatar:    user.Avatar,
				CreatedAt: user.CreatedAt.Unix(),
			},
			Token: token,
		},
	}
}

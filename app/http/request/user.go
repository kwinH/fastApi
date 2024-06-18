package request

type LoginRequest struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=5,max=30" trans:"用户名"` //用户名
	Password string `form:"password" json:"password" binding:"required,min=8,max=40" trans:"密码"`    //密码
}

type RegisterRequest struct {
	Nickname        string `form:"nickname" json:"nickname" binding:"required,min=2,max=30" trans:"昵称"`                    //昵称
	UserName        string `form:"user_name" json:"user_name" binding:"required,min=5,max=30" trans:"用户名"`                 //用户名
	Password        string `form:"password" json:"password" binding:"required,min=8,max=40" trans:"密码"`                    //密码
	PasswordConfirm string `form:"password_confirm" json:"password_confirm" binding:"required,min=8,max=40"  trans:"确认密码"` //确认密码
}

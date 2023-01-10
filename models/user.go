package models

type User struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Password string `json:"password"`
	AddTime  int64  `json:"add_time"`
	Status   int    `json:"status"`
	Mobile   string `json:"mobile"`
	Avatar   string `json:"avatar"`
}

type RegisterForm struct {
	UserName   string `json:"user_name" binding:"required"`                                 // 用户名，必填
	Mobile     string `json:"mobile" binding:"required,mobile"`                             //手机号，必填
	PassWord   string `json:"password" binding:"required,min=6,max=20"`                     // 密码，长度范围6到20位
	RePassWord string `json:"re_password" binding:"required,min=6,max=20,eqfield=PassWord"` // 确认密码，必须与密码一致
}

type PasswordLoginForm struct {
	Mobile   string `json:"mobile" binding:"required,mobile"`
	PassWord string `json:"password" binding:"required,min=6,max=20"`
}

type LoginReply struct {
	Id   int64 `json:"id"`
	Role int   `json:"role"`
}

type TokenReply struct {
	Id        int64  `json:"id"`         // 用户id
	Token     string `json:"token"`      // 颁发的token
	ExpiredAt int64  `json:"expired_at"` // token失效时间
}

type UserListForm struct {
	PageNum  int64 `json:"page_num"`
	PageSize int64 `json:"page_size"`
}

type UserListReply struct {
	Id      int64  `json:"id"`       // 用户id
	Name    string `json:"name"`     // 用户名
	AddTime int64  `json:"add_time"` // 注册时间
	Status  int    `json:"status"`   // 用户状态
	Mobile  string `json:"mobile"`   // 用户手机号
}

package model

type User struct {
	Model
	Email    string `json:"email" form:"email" gorm:"unique_index:idx_email"`
	UserName string `json:"user_name" form:"user_name" gorm:"unique_index:idx_user_name"`
	Password string `json:"password" form:"password"`
}

type LoginResp struct {
	User  User   `json:"user" form:"user"`
	Token string `json:"token" form:"token"`
}

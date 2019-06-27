package auth_service

import "go-gin-example/models"

type Auth struct {
	Username string
	Password string
}

//查询用户
func (a *Auth) Check() (bool, error) {
	return models.CheckAuth(a.Username, a.Password)
}

//添加用户
func (a *Auth) AddAuth() (bool, error) {
	return models.AddAuth(a.Username, a.Password)
}

//修改密码
func (a *Auth) ResetPassword(newPassword string) (bool, error) {
	return models.ResetPassword(a.Username, a.Password, newPassword)
}

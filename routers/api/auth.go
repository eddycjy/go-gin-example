package api

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"

	"go-gin-example/pkg/app"
	"go-gin-example/pkg/e"
	"go-gin-example/pkg/util"
	"go-gin-example/service/auth_service"
)

type auth struct {
	Username string `valid:"Required; MaxSize(50)"`
	Password string `valid:"Required; MaxSize(50)"`
}

func queryAuth(username, password string) (bool, error) {
	valid := validation.Validation{}
	a := auth{Username: username, Password: password}
	ok, err := valid.Valid(&a)
	if !ok {
		app.MarkErrors(valid.Errors)
		return ok, err
	}
	authService := auth_service.Auth{Username: username, Password: password}
	return authService.Check()
}

//重置密码
func ResetPassword(c *gin.Context) {
	appG := app.Gin{C: c}
	username := c.PostForm("username")
	oldPassword := c.PostForm("old_password")
	newPassword := c.PostForm("new_password")
	isExist, err := queryAuth(username, oldPassword)
	if isExist && (err == nil) && newPassword != "" {
		authService := auth_service.Auth{Username: username, Password: oldPassword}
		_, err := authService.ResetPassword(newPassword)
		if err != nil {
			appG.Response(http.StatusForbidden, e.ERROR, []rune{})
			return
		}
		appG.Response(http.StatusOK, e.SUCCESS, []rune{})
		return
	}
	appG.Response(http.StatusInternalServerError, e.InvalidParams, []rune{})
	return
}

//注册
func Register(c *gin.Context) {
	appG := app.Gin{C: c}

	username := c.PostForm("username")
	password := c.PostForm("password")
	isExist, err := queryAuth(username, password)

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorSQL, nil)
		return
	}

	if !isExist {
		authService := auth_service.Auth{Username: username, Password: password}
		_, _ = authService.AddAuth()
		appG.Response(http.StatusOK, e.SUCCESS, nil)
		return
	}
	appG.Response(http.StatusForbidden, e.ErrorAuthExist, nil)
	return

}

// @Summary Get Auth
// @Produce  json
// @Param username query string true "userName"
// @Param password query string true "password"
// @Success 200 {object} app.Response
// @Failure 500 {object} app.Response
// @Router /auth [get]
func GetAuth(c *gin.Context) {
	appG := app.Gin{C: c}

	username := c.Query("username")
	password := c.Query("password")
	isExist, err := queryAuth(username, password)

	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorAuthCheckTokenFail, nil)
		return
	}

	if !isExist {
		appG.Response(http.StatusUnauthorized, e.ErrorAuth, nil)
		return
	}

	token, err := util.GenerateToken(username, password)
	if err != nil {
		appG.Response(http.StatusInternalServerError, e.ErrorAuthToken, nil)
		return
	}

	appG.Response(http.StatusOK, e.SUCCESS, map[string]string{
		"token": token,
	})
}

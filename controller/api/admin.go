package api

import (
	"github.com/gin-gonic/gin"
)

type admin struct {
}

// Login .
// @Summary 管理员登录
// @Description 后台管理员登录
// @Accept  json
// @Produce  json
// @Param	username		formData	string		true	"用户名"
// @Param	password		formData	string		true	"密码"
// @Success 200 {string} string	"ok"
// @Failure 400 {object} api.ResultError "error"
// @Router /admin/login [post]
func (*admin) Login(c *gin.Context) {
	err := ResultError{}
	c.JSON(400, err)
}

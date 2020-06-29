package api

import (
	"github.com/gin-gonic/gin"
	"github.com/nilorg/naas/internal/service"
)

// 备注：
// 早期思路，后台管理系统不走OAuth2认证

type admin struct {
}

// Login .
// @Tags 		管理员（已弃用）
// @Summary 	管理员登录
// @Description 后台管理员登录
// @Accept  json
// @Produce  json
// @Param	username		formData	string		true	"用户名"
// @Param	password		formData	string		true	"密码"
// @Success 200 {string} string	"ok"
// @Failure 400 {object} api.Result "error"
// @Router /admin/login [post]
func (*admin) Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	su, err := service.Admin.Login(username, password)
	if err != nil {
		c.JSON(400, err)
	}
	c.JSON(200, su)
}

package controller

import (
	"net/http"
	"web_app/logic"
	"web_app/models"
	"web_app/utils"

	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// RegisterUserController 处理用户注册业务逻辑
// @Summary 用户注册接口
// @Description 用户填写用户名，手机号，密码并确认密码无误即可注册
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param object query models.RegisterForm true "查询参数"
// @Success 200
// @Router /users/register [post]
func RegisterUserController(c *gin.Context) {
	// 获取与校验参数
	f := new(models.RegisterForm)
	if err := c.ShouldBindJSON(f); err != nil {
		zap.L().Error("signup with invalid param", zap.Error(err))
		utils.HandleValidatorError(c, err)
		return
	}
	// 业务逻辑处理
	if e := logic.SignUp(c, f); e != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": e.Error(),
		})
		return
	}
	// 返回响应
	c.JSON(http.StatusOK, "register success")
}

type LoginReply *models.TokenReply

// LoginController  用户登录接口
// @Summary 用户登录接口
// @Description 用户通过手机号和密码登录
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param object query models.PasswordLoginForm true "查询参数"
// @Success 200 {object} LoginReply
// @Router /user/login [post]
func LoginController(c *gin.Context) {
	// 表单验证
	f := new(models.PasswordLoginForm)
	if err := c.ShouldBind(f); err != nil {
		utils.HandleValidatorError(c, err)
		return
	}
	// 业务逻辑处理
	err, tp := logic.SignIn(c, f)
	if err != nil {
		return
	}
	// 返回响应数据
	c.JSON(http.StatusOK, tp)
}

type ResponseUserList []*models.UserListReply

// GetUserListController  获取用户信息列表接口
// @Summary 用户信息列表接口
// @Description 可分页查询用户信息
// @Tags 用户相关接口
// @Accept application/json
// @Produce application/json
// @Param Authorization header string true "用户令牌"
// @Param object query models.UserListForm false "查询参数"
// @Security ApiKeyAuth
// @Success 200 {object} ResponseUserList
// @Router /user/list [get]
func GetUserListController(c *gin.Context) {
	// 参数验证
	var pageNum, pageSize int64
	pageNum = c.GetInt64("page_num")
	pageSize = c.GetInt64("page_size")
	if pageNum == 0 {
		pageNum = 1
	}
	if pageSize == 0 {
		pageSize = 5
	}
	// 业务逻辑处理
	err, userList := logic.GetUserList(c, pageNum, pageSize)
	if err != nil {
		return
	}
	// 返回响应
	c.JSON(http.StatusOK, userList)
}

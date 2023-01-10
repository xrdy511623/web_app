package logic

import (
	"fmt"
	"net/http"
	"time"
	"web_app/dao/mysql"
	"web_app/middlewares"
	"web_app/models"
	"web_app/pkg"
	"web_app/utils"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"

	"go.uber.org/zap"
)

func SignUp(ctx *gin.Context, f *models.RegisterForm) (err error) {
	// 判断用户是否存在
	err = mysql.CheckUserExist(f.UserName)
	if err != nil {
		zap.L().Error("user already exist", zap.Error(err))
		return
	}
	// 生成uid
	uid := pkg.GenId()
	// 密码加密
	encrypted := utils.EncryptPassword(f.PassWord)
	u := &models.User{
		Id:       uid,
		Name:     f.UserName,
		Password: encrypted,
		AddTime:  time.Now().Unix(),
		Status:   1,
		Mobile:   f.Mobile,
	}
	// 保存进数据库
	if err = mysql.SaveUser(u); err != nil {
		fmt.Printf("save user failed, err:%v", err)
		return
	}
	fmt.Printf("save user successful, user_name:%v", u.Name)
	return nil
}

func SignIn(ctx *gin.Context, f *models.PasswordLoginForm) (error, *models.TokenReply) {
	// 密码加密
	encrypted := utils.EncryptPassword(f.PassWord)
	u := &models.User{
		Password: encrypted,
		Mobile:   f.Mobile,
	}
	// 校验用户身份
	err, r := mysql.CheckUser(u)
	if err != nil {
		fmt.Printf("identify user failed, err:%v", err)
		return err, nil
	}
	// 颁发token
	j := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(r.Id),
		AuthorityId: uint(r.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: time.Now().Unix() + 60*60*24,
			Issuer:    "web_app",
		},
	}
	token, err := j.CreateToken(claims)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return err, nil
	}
	tp := &models.TokenReply{
		Id:        r.Id,
		Token:     token,
		ExpiredAt: (time.Now().Unix() + 60*60*24) * 1000,
	}
	return nil, tp
}

func GetUserList(ctx *gin.Context, pageNum, pageSize int64) (error, []*models.UserListReply) {
	// 获取用户数据信息
	err, userList := mysql.GetUserList(pageNum, pageSize)
	if err != nil {
		zap.L().Error("get user list failed", zap.Error(err))
		return err, nil
	}
	return nil, userList
}

package api

import (
	"CollabDoc-go/global"
	"CollabDoc-go/model/database"
	"CollabDoc-go/model/request"
	"CollabDoc-go/model/response"
	"fmt"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"time"
)

type UserApi struct {
}

func (userApi *UserApi) Register(c *gin.Context) {
	var req request.Register
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	session := sessions.Default(c)
	// 两次邮箱一致性判断
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("This email doesn't match the email to be verified", c)
		return
	}

	// 获取会话中存储的邮箱验证码
	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("Invalid verification code", c)
		return
	}

	// 判断邮箱验证码是否过期
	savedTime := session.Get("expire_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("The verification code has expired, please resend it", c)
		return
	}

	u := database.User{Username: req.Username, Password: req.Password, Email: req.Email}

	user, err := userService.Register(u)
	if err != nil {
		global.Log.Error("Failed to register user:", zap.Error(err))
		response.FailWithMessage("Failed to register user", c)
		return
	}

	// 注册成功后，生成 token 并返回
	//userApi.TokenNext(c, user)
	response.OkWithData(user, c)
}
func (userApi *UserApi) UserLogin(c *gin.Context) {
	var req database.User
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	if req, err = userService.GetUser(req); err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	response.OkWithData(req, c)
}

// EmailLogin 邮箱登录
func (userApi *UserApi) EmailLogin(c *gin.Context) {
	var req request.RegisterRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	fmt.Println(req)

	// 校验验证码
	//if store.Verify(req.CaptchaID, req.Captcha, true) {
	//u := database.User{Email: req.Email, Password: req.Password}
	//user, err := userService.EmailLogin(u)
	//if err != nil {
	//	global.Log.Error("Failed to login:", zap.Error(err))
	//	response.FailWithMessage("Failed to login", c)
	//	return
	//}

	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"accessToken": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MCwicGFzc3dvcmQiOiIxMjM0NTYiLCJyZWFsTmFtZSI6IlZiZW4iLCJyb2xlcyI6WyJzdXBlciJdLCJ1c2VybmFtZSI6InZiZW4iLCJpYXQiOjE3NTA2NTI4NzksImV4cCI6MTc1MTI1NzY3OX0.k3NlQ5Tv9De2y3HZTEYRp-iu5TAtvBD1rdjuFQjbdnE",
			"id":          0,
			"password":    "123456",
			"realName":    "Vben",
			"roles":       []string{"super"},
			"username":    "vben",
		},
		"error":   nil,
		"message": "ok",
	})
	// 登录成功后生成 token
	//userApi.TokenNext(c, user)
	//response.OkWithData(user, c)
	//return
	//}
	//response.FailWithMessage("Incorrect verification code", c)
}
func (userApi *UserApi) Login(c *gin.Context) {
	switch c.Query("flag") {
	case "email":
		userApi.EmailLogin(c)
	case "qq":
		userApi.UserLogin(c)
	default:
		userApi.UserLogin(c)
	}
}

// UserInfo 获取个人信息
func (userApi *UserApi) UserInfo(c *gin.Context) {
	//userID := utils.GetUserID(c)
	//var user database.User
	//err := c.ShouldBindQuery(&user)
	//if err != nil {
	//	response.FailWithMessage(err.Error(), c)
	//	return
	//}
	//user, err = userService.UserInfo(user.ID)
	//if err != nil {
	//	global.Log.Error("Failed to get user information:", zap.Error(err))
	//	response.FailWithMessage("Failed to get user information", c)
	//	return
	//}
	c.JSON(http.StatusOK, gin.H{
		"code": 0,
		"data": gin.H{
			"id":       123,              //user.ID,
			"realName": "wlc",            //user.Username,    // 映射前端需要的 realName
			"roles":    []string{"user"}, // 写死一个角色，或者查库
			"username": "wlc2",           //user.Username,
		},
		"error":   nil,
		"message": "ok",
	})
	//response.OkWithData(user, c)
}
func (userApi *UserApi) Codes(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"code":    0,
		"data":    []string{"AC_100100", "AC_100110", "AC_100120", "AC_100010"},
		"error":   nil,
		"message": "ok",
	})
}

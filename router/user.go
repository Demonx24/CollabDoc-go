package router

import (
	"CollabDoc-go/api"

	"github.com/gin-gonic/gin"
)

type UserRouter struct {
}

func (u *UserRouter) InitUserRouter(Router *gin.RouterGroup) {
	userRouter := Router.Group("user")
	//userPublicRouter := PublicRouter.Group("user")
	//userLoginRouter := PublicRouter.Group("user").Use(middleware.LoginRecord())
	//userAdminRouter := AdminRouter.Group("user")
	userApi := api.ApiGroupApp.UserApi
	//{
	//	userRouter.POST("logout", userApi.Logout)
	//	userRouter.PUT("resetPassword", userApi.UserResetPassword)
	userRouter.GET("info", userApi.UserInfo)
	//	userRouter.PUT("changeInfo", userApi.UserChangeInfo)
	//	userRouter.GET("weather", userApi.UserWeather)
	//	userRouter.GET("chart", userApi.UserChart)
	//}
	//{
	//	userPublicRouter.POST("forgotPassword", userApi.ForgotPassword)
	//	userPublicRouter.GET("card", userApi.UserCard)
	//}
	{
		userRouter.POST("emaillogin", userApi.EmailLogin)
		userRouter.POST("register", userApi.Register)
		userRouter.POST("login", userApi.Login)
		userRouter.GET("codes", userApi.Codes)
		userRouter.POST("logout", userApi.Logout)
	}
	//{
	//	userAdminRouter.GET("list", userApi.UserList)
	//	userAdminRouter.PUT("freeze", userApi.UserFreeze)
	//	userAdminRouter.PUT("unfreeze", userApi.UserUnfreeze)
	//	userAdminRouter.GET("loginList", userApi.UserLoginList)
	//}
}

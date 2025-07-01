package router

import (
	"CollabDoc-go/api"
	"github.com/gin-gonic/gin"
)

type RouterRouter struct{}

func (r *RouterRouter) InitRouterRouter(Router *gin.RouterGroup) {
	routerRouter := Router.Group("routes")
	//userPublicRouter := PublicRouter.Group("user")
	//userLoginRouter := PublicRouter.Group("user").Use(middleware.LoginRecord())
	//userAdminRouter := AdminRouter.Group("user")
	routerApi := api.ApiGroupApp.RouterApi
	routerRouter.GET("", routerApi.Router)
}

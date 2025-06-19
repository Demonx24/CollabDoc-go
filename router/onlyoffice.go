package router

import (
	"CollabDoc-go/api"
	"github.com/gin-gonic/gin"
)

type OnlyofficeRouter struct{}

func (r *OnlyofficeRouter) InitOnlyOfficeRouter(Router *gin.RouterGroup) {
	commentRouter := Router.Group("onlyoffice")
	commentApi := api.ApiGroupApp.OnlyofficeApi
	{

		commentRouter.POST("/callback", commentApi.Callback)
	}
}

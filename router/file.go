package router

import (
	"CollabDoc-go/api"

	"github.com/gin-gonic/gin"
)

type FileRouter struct {
}

func (b *BaseRouter) InitFileRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("file")

	FileApi := api.ApiGroupApp.FileApi
	{
		baseRouter.GET("", FileApi.FilePath)

	}
}

package router

import (
	"CollabDoc-go/api"

	"github.com/gin-gonic/gin"
)

type DocumentRouter struct {
}

func (b *BaseRouter) InitDocumentRouter(Router *gin.RouterGroup) {
	baseRouter := Router.Group("document")

	documentApi := api.ApiGroupApp.DocumentApi
	{
		baseRouter.POST("create", documentApi.CreateDocumentByUserUUid)
		baseRouter.GET("get", documentApi.GetDocumentByUserUUId)
		baseRouter.PUT("update", documentApi.UpdateDocument)
		baseRouter.GET("public", documentApi.GetPublicDocuments)
		baseRouter.GET("getversions", documentApi.GetVersions)
		baseRouter.GET("diff", documentApi.GetDiff)
	}
}

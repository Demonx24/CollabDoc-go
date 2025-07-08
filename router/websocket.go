package router

import (
	"CollabDoc-go/api"
	"github.com/gin-gonic/gin"
)

type WebsocketRouter struct {
}

func (u *WebsocketRouter) InitWebsocketRouter(Router *gin.RouterGroup) {
	websocketRouter := Router.Group("ws")
	//websocketApi := api.ApiGroupApp.WebsocketApi
	websocketRouter.GET("", api.DocWebSocketHandler)
}

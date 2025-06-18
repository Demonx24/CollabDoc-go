package service

import (
	"CollabDoc-go/api"
	"CollabDoc-go/store"
	"CollabDoc-go/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Run() {
	store.InitRedis()
	store.InitMongo()

	r := gin.Default()
	r.Use(cors.Default())

	r.GET("/ws/:docId", ws.HandleWebSocket)
	r.POST("/doc", api.CreateDoc)
	r.GET("/doc/:id", api.GetDoc)

	r.Run(":8080")
}

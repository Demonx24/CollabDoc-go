package api

import (
	"CollabDoc-go/model/database"
	"CollabDoc-go/model/request"
	"CollabDoc-go/model/response"
	"CollabDoc-go/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func DocWebSocketHandler(c *gin.Context) {
	var req request.Websocketlogin
	if err := c.ShouldBindQuery(&req); err != nil {
		response.Fail(c)
		return
	}
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	room := service.GetRoom(req.DocId)

	client := &database.Client{
		Useruuid: req.Useruuid, // 这里可替换为真实用户ID
		Send:     make(chan database.Operation, 256),
	}

	room.AddClient(client)

	initOp := database.Operation{
		Type:    "sync",
		Text:    room.Text,
		Version: room.Version,
	}

	if err := conn.WriteJSON(initOp); err != nil {
		log.Println("发送初始化文档状态失败:", err)
	}

	go func() {
		for op := range client.Send {
			if err := conn.WriteJSON(op); err != nil {
				break
			}
		}
		conn.Close()
	}()

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println("读取消息失败:", err)
			break
		}

		var op database.Operation
		if err := json.Unmarshal(msg, &op); err != nil {
			log.Println("消息反序列化失败:", err)
			continue
		}

		newOp, err := room.ProcessOperation(op)
		if err != nil {
			if newOp.Type == "sync" {
				_ = conn.WriteJSON(newOp)
			}
			continue
		}

		log.Printf("收到客户端操作: %+v", op)

		room.Broadcast(newOp, client)
	}

	room.RemoveClient(client)
}

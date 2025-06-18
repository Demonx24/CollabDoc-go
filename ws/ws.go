package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	User    string `json:"user"`
	Content string `json:"content"`
	OpType  string `json:"op_type"`
}

func HandleWebSocket(c *gin.Context) {
	docId := c.Param("docId")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	defer conn.Close()

	for {
		var msg Message
		err := conn.ReadJSON(&msg)
		if err != nil {
			log.Println("read error:", err)
			break
		}
		// 示例：广播到所有客户端，后续需构建房间管理
		log.Printf("[%s] %s: %s", docId, msg.User, msg.Content)
	}
}

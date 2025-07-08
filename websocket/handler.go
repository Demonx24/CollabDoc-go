package websocket

import (
	"CollabDoc-go/ot"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"sync"
)

type Client struct {
	Conn   *websocket.Conn
	UserID string
}

type DocumentRoom struct {
	ID      string
	Text    string
	Version int
	Clients map[*Client]bool
	Mu      sync.Mutex
	History []ot.Operation
}

var Rooms = make(map[string]*DocumentRoom)
var RoomsMu sync.Mutex

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func getOrCreateRoom(docId string) *DocumentRoom {
	RoomsMu.Lock()
	defer RoomsMu.Unlock()
	if r, ok := Rooms[docId]; ok {
		return r
	}
	r := &DocumentRoom{
		ID:      docId,
		Clients: make(map[*Client]bool),
		Text:    "",
		Version: 0,
	}
	Rooms[docId] = r
	return r
}

func HandleWebSocket(c *gin.Context) {
	docId := c.Query("docId")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}
	client := &Client{Conn: conn}
	r := getOrCreateRoom(docId)

	r.Mu.Lock()
	r.Clients[client] = true
	r.Mu.Unlock()

	// 发送初始文本与版本
	conn.WriteJSON(ot.Operation{
		Type:    "sync",
		Text:    r.Text,
		Version: r.Version,
	})

	go readPump(client, r)
}

func readPump(client *Client, room *DocumentRoom) {
	defer func() {
		room.Mu.Lock()
		delete(room.Clients, client)
		room.Mu.Unlock()
		client.Conn.Close()
	}()

	for {
		_, msg, err := client.Conn.ReadMessage()
		if err != nil {
			log.Println("Read error:", err)
			break
		}
		var op ot.Operation
		if err := json.Unmarshal(msg, &op); err != nil {
			log.Println("Unmarshal error:", err)
			continue
		}
		applyAndBroadcast(room, op, client)
	}
}

func applyAndBroadcast(room *DocumentRoom, op ot.Operation, sender *Client) {
	room.Mu.Lock()
	defer room.Mu.Unlock()

	if op.Version != room.Version+1 && len(room.History) > 0 {
		op = ot.Transform(op, room.History[len(room.History)-1])
	}

	newText, err := op.Apply(room.Text)
	if err != nil {
		log.Println("Apply failed:", err)
		return
	}

	room.Text = newText
	room.Version++
	op.Version = room.Version
	room.History = append(room.History, op)

	for c := range room.Clients {
		if c != sender {
			c.Conn.WriteJSON(op)
		}
	}
}

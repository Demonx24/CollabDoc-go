package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"CollabDoc-go/global"
	"CollabDoc-go/model/database"
	"CollabDoc-go/ot"
)

type DocumentRoom struct {
	*database.DocumentRoom
}

var rooms = map[string]*DocumentRoom{}
var roomsMu sync.Mutex

func GetRoom(id string) *DocumentRoom {
	roomsMu.Lock()
	defer roomsMu.Unlock()

	if r, ok := rooms[id]; ok {
		return r
	}
	// 先尝试从 Redis 读文本和版本
	text := getRedisText(id)
	version := getRedisVersion(id)
	if text == "" || version == 0 {
		// Redis缓存不存在，去Mongo查
		docMd, err := GetDocumentMd(id)
		if err == nil && docMd != nil {
			text = docMd.Text
			version = docMd.Version
			// 缓存回 Redis
			global.Redis.Set(context.Background(), "doc:"+id+":text", text, 0)
			global.Redis.Set(context.Background(), "doc:"+id+":version", version, 0)
		}
	}
	// 初始化数据库层的 DocumentRoom 实例
	dbRoom := &database.DocumentRoom{
		DocID:   id,
		Text:    text, // 复用前面取到的缓存或者mongo数据
		Version: version,
		History: []database.Operation{},
		Clients: make(map[*database.Client]bool),
	}

	// 用 service 层的 DocumentRoom 包装 database.DocumentRoom
	room := &DocumentRoom{dbRoom}

	rooms[id] = room

	return room
}

func getRedisVersion(docID string) int {
	val, err := global.Redis.Get(context.Background(), "doc:"+docID+":version").Result()
	if err != nil {
		return 0
	}
	v, _ := strconv.Atoi(val)
	return v
}

func getRedisText(docID string) string {
	val, err := global.Redis.Get(context.Background(), "doc:"+docID+":text").Result()
	if err != nil {
		return ""
	}
	return val
}

// 应用操作（版本校验、OT转换、保存Redis）
func (r *DocumentRoom) ProcessOperation(op database.Operation) (database.Operation, error) {
	r.Mu.Lock()
	defer r.Mu.Unlock()

	log.Printf("ProcessOperation 开始：op.version=%d，当前版本=%d", op.Version, r.Version)

	if op.Version != r.Version+1 {
		if op.Version <= r.Version {
			transformed, err, _ := ot.ApplyAndTransform(op, r.Text, r.History[op.Version:])
			fmt.Println("前端版本低于服务端"+transformed, err)
		}
		log.Printf("版本不一致，当前=%d，收到=%d", r.Version, op.Version)
		return database.Operation{
			Type:    "sync",
			Text:    r.Text,
			Version: r.Version,
		}, errors.New("version mismatch, need sync")
	}

	baseIdx := op.Version - len(r.History) - 1
	if baseIdx < 0 {
		baseIdx = 0
	}
	if baseIdx > len(r.History) {
		baseIdx = len(r.History)
	}
	historySlice := r.History[baseIdx:]

	newText, transformedOp, err := ot.ApplyAndTransform(op, r.Text, historySlice)
	if err != nil {
		log.Printf("操作应用失败：%v", err)
		return database.Operation{}, err
	}

	r.Text = newText
	r.Version = transformedOp.Version
	r.History = append(r.History, transformedOp)

	// 写 Redis创建一个带超时的 context，2秒后自动取消
	ctx1, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := global.Redis.Set(ctx1, "doc:"+r.DocID+":text", r.Text, 0).Err(); err != nil {
		log.Printf("写 Redis 文本失败: %v", err)
	}

	ctx2, cancel2 := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel2()

	if err := global.Redis.Set(ctx2, "doc:"+r.DocID+":version", r.Version, 0).Err(); err != nil {
		log.Printf("写 Redis 版本失败: %v", err)
	}
	if err = SaveDocumentMd(r.DocID, r.Text, r.Version, transformedOp); err != nil {
		log.Printf("MD保存monggo错误 %v", err)
		return database.Operation{}, err
	}

	log.Printf("操作完成：type=%s, pos=%d, text=%q, version=%d",
		transformedOp.Type, transformedOp.Position, transformedOp.Text, transformedOp.Version)
	return transformedOp, nil
}

// 添加客户端
func (r *DocumentRoom) AddClient(c *database.Client) {
	r.Mu.Lock()
	r.Clients[c] = true
	r.Mu.Unlock()
	log.Printf("文档 %s 新增客户端，当前连接数：%d", r.DocID, len(r.Clients))
}

// 移除客户端
func (r *DocumentRoom) RemoveClient(c *database.Client) {
	r.Mu.Lock()
	delete(r.Clients, c)
	r.Mu.Unlock()
	log.Printf("文档 %s 客户端断开，剩余连接数：%d", r.DocID, len(r.Clients))
}

// 广播操作到所有客户端（排除exclude）
func (r *DocumentRoom) Broadcast(op database.Operation, exclude *database.Client) {
	r.Mu.Lock()
	defer r.Mu.Unlock()
	for c := range r.Clients {
		if c == exclude {
			continue
		}
		select {
		case c.Send <- op:
			log.Printf("已广播版本 %d", op.Version)
		default:
			log.Printf("客户端发送队列已满，广播失败，version=%d", op.Version)
		}
	}
}

package database

import (
	"sync"
	"time"
)

// 操作结构，记录一次插入/删除/同步行为
type Operation struct {
	UserID   string    `json:"userId" bson:"userId"`     // 操作用户 UUID
	Type     string    `json:"type" bson:"type"`         // insert/delete/sync
	Position int       `json:"position" bson:"position"` // 光标位置（rune索引）
	Text     string    `json:"text" bson:"text"`         // 操作文本
	Version  int       `json:"version" bson:"version"`   // 操作对应版本号
	Time     time.Time `json:"time" bson:"time"`         // 操作时间戳
}

// 客户端连接
type Client struct {
	Useruuid string         // 用户ID，便于追踪是谁连入
	Send     chan Operation // 发给该客户端的消息
}

// 文档房间，代表一个文档的协作空间
type DocumentRoom struct {
	DocID     string           `json:"docId" bson:"docId"`         // 文档 UUID
	OwnerID   string           `json:"ownerId" bson:"ownerId"`     // 文档所有者
	Title     string           `json:"title" bson:"title"`         // 文档标题
	Text      string           `json:"text" bson:"text"`           // 当前最新内容
	Version   int              `json:"version" bson:"version"`     // 当前版本号
	History   []Operation      `json:"history" bson:"history"`     // 保存在内存的部分历史操作
	Clients   map[*Client]bool `json:"-" bson:"-"`                 // 当前活跃连接（不保存到 mongo）
	Mu        sync.Mutex       `json:"-" bson:"-"`                 // 并发锁
	UpdatedAt time.Time        `json:"updatedAt" bson:"updatedAt"` // 最后更新时间
	CreatedAt time.Time        `json:"createdAt" bson:"createdAt"` // 创建时间
}
type DocumentMd struct {
	DocID     string      `bson:"doc_id"`
	Text      string      `bson:"text"`
	Version   int         `bson:"version"`
	History   []Operation `bson:"history"`
	UpdatedAt time.Time   `bson:"updated_at"`
}

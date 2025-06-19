package database

import "CollabDoc-go/global"

// Document 文档表
type User_Documents struct {
	global.MODEL
	Title            string `gorm:"size:255;not null" json:"title"`
	OwnerID          uint   `gorm:"not null" json:"owner_id"`
	TeamID           *uint  `json:"team_id"`            // 允许为空
	CurrentVersionID *uint  `json:"current_version_id"` // 允许为空
}

// DocumentVersion 文档版本表
type DocumentVersion struct {
	global.MODEL
	DocumentID  uint   `gorm:"not null" json:"document_id"`
	VersionName string `gorm:"size:64" json:"version_name"`
	FilePath    string `gorm:"size:512;not null" json:"file_path"`
	CreatedBy   uint   `gorm:"not null" json:"created_by"`
	// CreatedAt 可用 global.MODEL 的 CreatedAt
}

// DocumentPermission 文档权限表
type DocumentPermission struct {
	global.MODEL
	DocumentID uint   `gorm:"not null" json:"document_id"`
	UserID     uint   `gorm:"not null" json:"user_id"`
	Permission string `gorm:"type:enum('owner','edit','view');not null" json:"permission"`
}

// EditLog 编辑日志表
type EditLog struct {
	global.MODEL
	DocumentID uint   `gorm:"not null" json:"document_id"`
	UserID     uint   `gorm:"not null" json:"user_id"`
	Action     string `gorm:"size:64;not null" json:"action"`
	WordCount  int    `json:"word_count"`
	// CreatedAt 用 global.MODEL 的 CreatedAt
}

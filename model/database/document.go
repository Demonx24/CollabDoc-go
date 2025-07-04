package database

import (
	"CollabDoc-go/global"
	"time"
)

// Document 文档表
type User_Documents struct {
	global.MODEL
	Title            string  `gorm:"size:255;not null" json:"title"`
	OwnerID          string  `gorm:"not null" json:"owner_id"`
	DocUUID          string  `gorm:"type:char(36);uniqueIndex;not null" json:"doc_uuid"`                      // 文档唯一标识
	TeamID           *uint   `json:"team_id"`                                                                 // 允许为空
	CurrentVersionID *uint   `json:"current_version_id"`                                                      // 允许为空
	DocType          string  `gorm:"size:50;not null;default:'general'" json:"doc_type"`                      // 文档类型
	Status           string  `gorm:"type:enum('active','archived','deleted');default:'active'" json:"status"` // 文档状态
	Description      *string `gorm:"type:text" json:"description,omitempty"`                                  // 文档描述，可空
	IsPublic         *bool   `gorm:"default:false" json:"is_public"`
}

// DocumentVersion 文档版本表
type DocumentVersion struct {
	global.MODEL
	DocumentID    uint   `gorm:"not null" json:"document_id"`
	VersionName   string `gorm:"size:64" json:"version_name"`
	FilePath      string `gorm:"size:512;not null" json:"file_path"`
	CreatedBy     string `gorm:"not null" json:"created_by"`
	VersionNumber uint   `gorm:"not null" json:"version_number"`    // 版本号
	Checksum      string `gorm:"size:64" json:"checksum,omitempty"` // 文件校验码（可选）
	// CreatedAt 可用 global.MODEL 的 CreatedAt
}

// DocumentPermission 文档权限表
type DocumentPermission struct {
	global.MODEL
	DocumentID uint   `gorm:"not null" json:"document_id"`
	UserID     uint   `gorm:"not null" json:"user_id"`
	Permission string `gorm:"type:enum('owner','edit','view');not null" json:"permission"`
}

type DocDiff struct {
	DocUUID       string                 `bson:"doc_uuid"`
	FromVersion   int                    `bson:"from_version"`
	ToVersion     int                    `bson:"to_version"`
	ChangedFields map[string]interface{} `bson:"changed_fields"` // 差异内容
	CreatedAt     time.Time              `bson:"created_at"`
}
type DiffMessage struct {
	DocUUID     string `json:"doc_uuid"`
	FromVersion int    `json:"from_version"`
	ToVersion   int    `json:"to_version"`
}
type DiffItem struct {
	Operation string `json:"operation"` // "INSERT", "DELETE", "EQUAL"
	Text      string `json:"text"`
}

// 表名：documents
func (User_Documents) TableName() string {
	return "documents"
}

// 表名：document_versions
func (DocumentVersion) TableName() string {
	return "document_versions"
}

// 表名：document_permissions
func (DocumentPermission) TableName() string {
	return "document_permissions"
}

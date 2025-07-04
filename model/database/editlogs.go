package database

import "time"

type DocumentEditLog struct {
	ID            int64     `gorm:"primaryKey;autoIncrement" json:"id"`
	DocumentUUID  string    `gorm:"type:char(36);not null" json:"document_uuid"` // 文档 UUID
	UserUUID      string    `gorm:"type:char(36);not null" json:"user_uuid"`     // 用户 UUID
	Action        string    `gorm:"type:varchar(64);not null" json:"action"`     // 编辑操作类型
	WordCount     *int      `gorm:"type:int" json:"word_count"`                  // 影响的字数（可为空）
	VersionNumber int       `gorm:"not null" json:"version_number"`              // 文档版本号
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`            // 创建时间
}

func (DocumentEditLog) TableName() string {
	return "document_edit_logs"
}

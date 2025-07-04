package request

import "time"

type CreateDocRequest struct {
	UserUUID string `json:"owner_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	DocType  string `json:"doc_type" binding:"required,oneof=docx pptx xlsx"`
}

type DocumentByUserRequest struct {
	UserUUID string `json:"uuid" binding:"required"`
}
type UpdateDocument struct {
	DocUUID     string    `gorm:"type:char(36);uniqueIndex;not null" json:"doc_uuid"`
	Title       string    `gorm:"size:255;not null" json:"title"`
	DocType     string    `gorm:"size:50;not null;default:'general'" json:"doc_type"`                      // 文档类型
	Status      string    `gorm:"type:enum('active','archived','deleted');default:'active'" json:"status"` // 文档状态
	Description *string   `gorm:"type:text" json:"description,omitempty"`                                  // 文档描述，可空
	IsPublic    bool      `gorm:"default:false" json:"is_public"`
	UpdatedAt   time.Time `json:"updated_at"` // 更新时间
}
type GetVersions struct {
	DocumentID uint `gorm:"not null" json:"document_id"form:"document_id"`
}
type GetDiff struct {
	DocUUID string `gorm:"type:char(36);uniqueIndex;not null" json:"doc_uuid" form:"doc_uuid"`
	FromVer int    `gorm:"size:255;not null" json:"from_ver" form:"from"`
	ToVer   int    `gorm:"size:255;not null" json:"to_ver"form:"to"`
}

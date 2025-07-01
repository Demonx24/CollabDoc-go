package request

import "github.com/gofrs/uuid"

type CommentInfoByArticleID struct {
	ArticleID string `json:"article_id" uri:"article_id" binding:"required"`
}

type CommentCreate struct {
	UserUUID  uuid.UUID `json:"-"`
	ArticleID string    `json:"article_id" binding:"required"`
	PID       *uint     `json:"p_id"`
	Content   string    `json:"content" binding:"required,max=320"`
}

type CommentDelete struct {
	IDs []uint `json:"ids"`
}

type CommentList struct {
	ArticleID *string `json:"article_id" form:"article_id"`
	UserUUID  *string `json:"user_uuid" form:"user_uuid"`
	Content   *string `json:"content" form:"content"`
	PageInfo
}
type FilePath struct {
	OwnerID string `gorm:"not null" json:"owner_id" form:"owner_id"`
	DocUUID string `gorm:"type:char(36);uniqueIndex;not null" json:"doc_uuid" form:"doc_uuid"`
	ID      uint   `form:"id" json:"id" gorm:"primarykey"` // 主键 ID
	Ext     string `form:"ext" json:"ext" binding:"required,max=32"`
}

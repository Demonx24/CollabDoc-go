package response

import (
	"CollabDoc-go/model/database"
)

type CreateDocResponse struct {
	Document database.User_Documents `json:"document"`
}

type GetDocResponse struct {
	ID       uint    `form:"id" json:"id" gorm:"primarykey"` // 主键 ID
	DocUUID  string  `json:"doc_uuid"`                       // 用于跳转和分享
	Title    string  `json:"title"`                          // 标题
	DocType  string  `json:"doc_type"`                       // 类型 docx/pptx/xlsx
	Status   string  `json:"status"`                         // 状态 active/archived/deleted
	IsPublic bool    `json:"is_public"`                      // 是否公开
	Summary  *string `json:"summary,omitempty"`              // 文档简介（可以截断 description）
	Updated  string  `json:"updated_at"`
}

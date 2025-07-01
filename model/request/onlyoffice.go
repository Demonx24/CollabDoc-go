package request

type OnlyofficeUser struct {
	UserUUID string `json:"id" form:"user_uuid"`
	Name     string `json:"name" form:"name"`
	DocUUID  string `gorm:"type:char(36);uniqueIndex;not null" json:"doc_uuid" form:"doc_uuid"`
}

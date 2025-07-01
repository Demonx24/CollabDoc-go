package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/gofrs/uuid"
	"time"
)

type User struct {
	ID           uint           `gorm:"primaryKey;autoIncrement" json:"id" form:"id"`
	UUID         uuid.UUID      `gorm:"type:char(36);uniqueIndex;not null" json:"uuid" form:"uuid"`
	Username     string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"username" form:"username"`
	Password     string         `gorm:"type:varchar(255);not null" json:"-" form:"password"` // 不序列化到JSON，允许绑定表单
	Nickname     string         `gorm:"type:varchar(50)" json:"nickname" form:"nickname"`
	Email        string         `gorm:"type:varchar(100);uniqueIndex" json:"email" form:"email"`
	Avatar       string         `gorm:"type:varchar(255)" json:"avatar" form:"avatar"`
	Address      string         `gorm:"type:varchar(255)" json:"address" form:"address"`
	OpenID       string         `gorm:"column:openid;type:varchar(100)" json:"openid" form:"openid"`
	AccessToken  string         `gorm:"type:text" json:"accessToken" form:"accessToken"`
	RefreshToken string         `gorm:"type:text" json:"refreshToken" form:"refreshToken"`
	Permissions  JSONStringList `gorm:"type:json" json:"permissions" form:"permissions"`
	Roles        JSONStringList `gorm:"type:json" json:"roles" form:"roles"`
	Expires      *time.Time     `json:"expires" form:"expires"`
	IsFrozen     bool           `gorm:"default:false" json:"isFrozen" form:"isFrozen"`
	CreatedAt    time.Time      `json:"createdAt" form:"createdAt"`
	UpdatedAt    time.Time      `json:"updatedAt" form:"updatedAt"`
	DeletedAt    *time.Time     `gorm:"index" json:"deletedAt,omitempty" form:"deletedAt"`
}

func (User) TableName() string {
	return "users"
}

type JSONStringList []string

func (j *JSONStringList) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(bytes, j)
}

func (j JSONStringList) Value() (driver.Value, error) {
	return json.Marshal(j)
}

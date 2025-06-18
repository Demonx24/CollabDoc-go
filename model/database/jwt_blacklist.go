package database

import "CollabDoc-go/global"

// JwtBlacklist JWT 黑名单表
type JwtBlacklist struct {
	global.MODEL
	Jwt string `json:"jwt" gorm:"type:text"` // Jwt
}

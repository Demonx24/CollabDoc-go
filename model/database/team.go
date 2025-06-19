package database

import "CollabDoc-go/global"

type Team struct {
	global.MODEL
	Name     string `gorm:"size:128;not null" json:"name"`
	LeaderID uint   `json:"leader_id"`
}
type TeamMember struct {
	global.MODEL
	TeamID uint   `gorm:"not null" json:"team_id"`
	UserID uint   `gorm:"not null" json:"user_id"`
	Role   string `gorm:"size:32;not null" json:"role"`
	// JoinedAt 可用 CreatedAt 代替
}

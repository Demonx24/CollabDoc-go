package database

import (
	"encoding/json"
	"time"
)

// RouteMenu 对应数据库中的 route_menu 表
// RouteMenu 对应数据库中的 route_menu 表
type RouteMenu struct {
	ID        int             `gorm:"column:id;primaryKey;autoIncrement"`
	ParentID  int             `gorm:"column:parent_id;not null"`
	Path      string          `gorm:"column:path;size:128;not null"`
	Name      *string         `gorm:"column:name;size:64"`  // 改为指针以区分空值
	Title     *string         `gorm:"column:title;size:64"` // 同上
	Icon      *string         `gorm:"column:icon;size:64"`
	Roles     json.RawMessage `gorm:"column:roles;type:json"` // e.g. ["user","admin"]
	KeepAlive bool            `gorm:"column:keep_alive;type:boolean;default:false"`
	Redirect  *string         `gorm:"column:redirect;size:128"`
	SortOrder int             `gorm:"column:sort_order;type:int;default:0"`
	Component *string         `gorm:"column:component;size:128"`
	FrameSrc  *string         `gorm:"column:frame_src;size:255"`
	CreatedAt time.Time       `gorm:"column:created_at"`
	UpdatedAt time.Time       `gorm:"column:updated_at"`
}

// TableName 指定表名
func (RouteMenu) TableName() string {
	return "route_menu"
}

// RouteMeta 是前端约定的 meta 信息
type RouteMeta struct {
	Title      string   `json:"title"`
	Icon       string   `json:"icon,omitempty"`
	Roles      []string `json:"roles,omitempty"`
	KeepAlive  bool     `json:"keepAlive,omitempty"`
	FrameSrc   *string  `json:"frameSrc,omitempty"`
	ShowLink   bool     `json:"showLink,omitempty"` // 新增字段
	ActivePath string   `json:"activePath,omitempty"`
}

// Route 是返回给前端的路由结构
type Route struct {
	Path      string    `json:"path"`
	Name      string    `json:"name"`
	Component *string   `json:"component,omitempty"`
	Redirect  string    `json:"redirect,omitempty"`
	Meta      RouteMeta `json:"meta"`
	Children  []Route   `json:"children,omitempty"`
}

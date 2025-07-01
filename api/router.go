package api

import (
	"CollabDoc-go/global"
	"CollabDoc-go/model/database"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

type RouterApi struct{}

var ctx = context.Background()

func (routerApi *RouterApi) Router(c *gin.Context) {
	var all []database.RouteMenu
	const cacheKey = "route:tree"

	// 先尝试读取缓存
	if cached, err := global.Redis.Get(ctx, cacheKey).Result(); err == nil && cached != "" {
		var cachedRoutes []database.Route
		if err := json.Unmarshal([]byte(cached), &cachedRoutes); err == nil {
			fmt.Println("从缓存读取路由成功")
			c.JSON(http.StatusOK, cachedRoutes)
			return
		}
	}

	// 从数据库查询所有路由
	db := global.DB
	if err := db.Order("parent_id, sort_order").Find(&all).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 打印所有菜单ID及parent_id，确认数据关系
	fmt.Println("数据库中路由数据:")
	for _, m := range all {
		fmt.Printf("ID=%d, ParentID=%d, Path=%s, Name=%s\n", m.ID, m.ParentID, m.Path, m.Name)
	}

	// 按 parent_id 分组，方便递归构建树
	groups := make(map[int][]database.RouteMenu)
	for _, m := range all {
		groups[m.ParentID] = append(groups[m.ParentID], m)
	}

	fmt.Printf("分组后菜单数量: %d\n", len(groups))

	// 递归构建树形结构函数
	var build func(pid int) []database.Route
	build = func(pid int) []database.Route {
		var list []database.Route
		childrenMenus := groups[pid]
		fmt.Printf("构建parent_id=%d的子菜单数量=%d\n", pid, len(childrenMenus))

		for _, m := range childrenMenus {
			// 解析 roles JSON
			var roles []string
			if len(m.Roles) > 0 {
				err := json.Unmarshal(m.Roles, &roles)
				if err != nil {
					fmt.Printf("解析Roles失败，ID=%d, err=%v\n", m.ID, err)
				}
			}

			node := database.Route{
				Path: m.Path,
				Name: "",
				Meta: database.RouteMeta{
					Title:     " ptr(m.Title)",
					Roles:     roles,
					KeepAlive: m.KeepAlive,
				},
				Children: build(m.ID), // 递归构建子菜单
			}

			// 赋值可选字段 Icon, Redirect, Component, FrameSrc
			if m.Icon != nil {
				node.Meta.Icon = *m.Icon
			}
			if m.Redirect != nil {
				node.Redirect = *m.Redirect
			}
			if m.Component != nil && *m.Component != "" {
				node.Component = m.Component
			}
			if m.FrameSrc != nil && *m.FrameSrc != "" {
				// 假设 RouteMeta 有 FrameSrc 字段是 *string
				node.Meta.FrameSrc = m.FrameSrc
			}

			fmt.Printf("添加路由: Path=%s, Name=%s, 子菜单数量=%d\n", node.Path, node.Name, len(node.Children))

			list = append(list, node)
		}
		return list
	}

	routes := build(0)

	//	缓存序列化后的路由树
	if bytes, err := json.Marshal(routes); err == nil {
		global.Redis.Set(ctx, cacheKey, bytes, 10*time.Minute)
	}

	fmt.Println(routes)
	c.JSON(http.StatusOK, routes)
}

func ptr(s string) *string {
	return &s
}

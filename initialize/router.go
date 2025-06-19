package initialize

import (
	"CollabDoc-go/global"
	"CollabDoc-go/router"
	"fmt"

	//"CollabDoc-go/router"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"time"
)

func InitRouter() *gin.Engine {
	// 设置gin模式
	gin.SetMode(global.Config.System.Env)
	Router := gin.Default()
	Router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:8080"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization", "Token", "Lang", "X-Requested-With"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	//使用gin会话路由
	var store = cookie.NewStore([]byte(global.Config.System.SessionsSecret))
	Router.Use(sessions.Sessions("session", store))
	//静态网页
	Router.Static("/static", "./static")
	// 将指定目录下的文件提供给客户端
	// "uploads" 是URL路径前缀，http.Dir("uploads")是实际文件系统中存储文件的目录
	Router.Static(global.Config.Upload.Path, "."+global.Config.Upload.Path)
	//Router.StaticFS(global.Config.Upload.Path, http.Dir(global.Config.Upload.Path))
	// 创建路由组

	routerGroup := router.RouterGroupApp
	//
	publicGroup := Router.Group(global.Config.System.RouterPrefix)
	//
	routerGroup.InitOnlyOfficeRouter(publicGroup)
	routerGroup.InitBaseRouter(publicGroup)
	routerGroup.InitUserRouter(publicGroup)
	for _, route := range Router.Routes() {
		fmt.Println(route.Method, route.Path)
	}
	return Router
}

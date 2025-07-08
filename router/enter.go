package router

type RouterGroup struct {
	OnlyofficeRouter
	BaseRouter
	UserRouter
	RouterRouter
	FileRouter
	WebsocketRouter
}

var RouterGroupApp = new(RouterGroup)

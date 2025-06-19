package router

type RouterGroup struct {
	OnlyofficeRouter
	BaseRouter
	UserRouter
}

var RouterGroupApp = new(RouterGroup)

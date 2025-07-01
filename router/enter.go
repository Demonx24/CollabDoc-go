package router

type RouterGroup struct {
	OnlyofficeRouter
	BaseRouter
	UserRouter
	RouterRouter
	FileRouter
}

var RouterGroupApp = new(RouterGroup)

package api

import "CollabDoc-go/service"

type ApiGroup struct {
	OnlyofficeApi
	UserApi
	BaseApi
}

var ApiGroupApp = new(ApiGroup)
var userService = service.ServiceGroupApp.UserService
var onlyofficeService = service.ServiceGroupApp.OnlyofficeService
var baseService = service.ServiceGroupApp.BaseService

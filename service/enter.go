package service

type ServiceGroup struct {
	UserService
	OnlyofficeService
	BaseService
}

var ServiceGroupApp = new(ServiceGroup)

package service

type ServiceGroup struct {
	UserService
	OnlyofficeService
	BaseService
	DocumentService
	Document_vService
	MinioService
}

var ServiceGroupApp = new(ServiceGroup)

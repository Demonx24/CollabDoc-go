package service

type ServiceGroup struct {
	UserService
	OnlyofficeService
	BaseService
	DocumentService
	Document_vService
	MinioService
	KafkaService
	MongoService
	EditLogService
	DocumentRoom
	JwtService
}

var ServiceGroupApp = new(ServiceGroup)

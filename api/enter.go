package api

import "CollabDoc-go/service"

type ApiGroup struct {
	OnlyofficeApi
	UserApi
	BaseApi
	RouterApi
	DocumentApi
	FileApi
}

var ApiGroupApp = new(ApiGroup)
var userService = service.ServiceGroupApp.UserService
var onlyofficeService = service.ServiceGroupApp.OnlyofficeService
var baseService = service.ServiceGroupApp.BaseService
var documentService = service.ServiceGroupApp.DocumentService
var documnet_vService = service.ServiceGroupApp.Document_vService
var minioService = service.ServiceGroupApp.MinioService
var kafkaService = service.ServiceGroupApp.KafkaService
var mongoService = service.ServiceGroupApp.MongoService
var editlogService = service.ServiceGroupApp.EditLogService

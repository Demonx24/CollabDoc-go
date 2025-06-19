package main

import (
	"CollabDoc-go/core"
	"CollabDoc-go/global"
	"CollabDoc-go/initialize"
)

func main() {
	global.Config = core.InitConf()
	global.Log = core.InitLogger()
	global.DB = initialize.InitGorm()
	//global.Redis = initialize.ConnectRedis()

	defer global.Redis.Close()

	core.RunServer()

}

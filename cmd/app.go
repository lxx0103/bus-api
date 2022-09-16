package cmd

import (
	"bus-api/api/v1/auth"
	"bus-api/api/v1/setting"
	"bus-api/core/cache"
	"bus-api/core/config"
	"bus-api/core/database"
	"bus-api/core/log"
	"bus-api/core/router"
)

func Run(args []string) {
	config.LoadConfig(args[1])
	log.ConfigLogger()
	cache.ConfigCache()
	database.ConfigMysql()
	r := router.InitRouter()
	router.InitPublicRouter(r, auth.Routers)
	router.InitAuthRouter(r, auth.AuthRouter, setting.AuthRouter)
	router.RunServer(r)
}

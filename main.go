package main

import (
	_ "app-service/bussiness-service/routers"

	daoApi "api/dao_service"
	"github.com/astaxie/beego"
)

func main() {
	var cfg = beego.AppConfig
	daoApi.ActionDaoApi.Init(cfg.String("ActionDaoService"))
	daoApi.AlgorithmDaoApi.Init(cfg.String("AlgorithmDaoService"))
	daoApi.BussinessDaoApi.Init(cfg.String("BussinessDaoService"))
	daoApi.ResourceDaoApi.Init(cfg.String("ResourceDaoService"))
	daoApi.UserDaoApi.Init(cfg.String("UserDaoService"))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

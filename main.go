package main

import (
	"casbin-demo/filters"
	"casbin-demo/models"
	_ "casbin-demo/routers"

	beego "github.com/beego/beego/v2/server/web"
	_ "github.com/go-sql-driver/mysql"
)

func main() {

	beego.InsertFilter("/v1/user/*", beego.BeforeRouter, filters.NewAuthorizer(models.Enforcer))

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

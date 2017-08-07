package main

import (
	_ "dev.cloud.360baige.com/routers"
	_ "dev.cloud.360baige.com/filters"
	"github.com/astaxie/beego"

)

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

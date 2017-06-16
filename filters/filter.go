package filters

import (
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego"
)

func init() {
	log := logs.NewLogger()
	log.SetLogger(logs.AdapterConsole)
	log.Debug("filter run start")

	beego.InsertFilter("/v1/company/*",beego.BeforeRouter, UrlOriginManager)
	beego.InsertFilter("/v1/user/*",beego.BeforeRouter, UrlOriginManager)

	log.Debug("filter run end")
}

var UrlManager = func(ctx *context.Context) {
	//ctx.Redirect(302, "/v1/user")
}

var UrlOriginManager = func(ctx *context.Context) {
	r := ctx.Request
	w := ctx.ResponseWriter
	Origin := r.Header.Get("Origin")
	if Origin != "" {
		w.Header().Add("Access-Control-Allow-Origin", Origin)
		w.Header().Add("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
		w.Header().Add("Access-Control-Allow-Headers", "x-requested-with,content-type")
		w.Header().Add("Access-Control-Allow-Credentials", "true")
	}
}
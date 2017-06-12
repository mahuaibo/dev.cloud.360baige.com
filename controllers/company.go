package controllers

import (
	"github.com/astaxie/beego"
)

type CompanyController struct {
	beego.Controller
}

// @Title 企业信息
// @Description 企业信息
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   userId     query   string true       "用户ID"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /detail [get]
func (c *CompanyController) Detail() {

}

// @Title 企业信息修改
// @Description 企业信息修改
// @Success 200 {"code":200,"messgae":"ok", "data":{ ... ... }}
// @Param   userId     query   string true       "用户ID"
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"..."}
// @router /modify [post]
func (c *CompanyController) Modify() {

}
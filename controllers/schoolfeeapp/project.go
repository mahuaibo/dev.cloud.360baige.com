package schoolfeeapp

import (
	"github.com/astaxie/beego"
)

// Project API
type ProjectController struct {
	beego.Controller
}

// @Title 校园收费列表接口
// @Description Project List 校园收费列表接口
// @Success 200 {"code":200,"messgae":"获取缴费项目成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取缴费项目失败"}
// @router /list [get]
func (c *ProjectController) ListOfProject() {

}

package schoolfeewin

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/mobile"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/account"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/utils"
	"time"
)

// Record API
type RecordController struct {
	beego.Controller
}

// @Title 校园收费列表接口
// @Description Project List 校园收费列表接口
// @Success 200 {"code":200,"messgae":"获取缴费项目成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取缴费项目失败"}
// @router /list [get]
func (c *RecordController) ListOfRecord() {

}

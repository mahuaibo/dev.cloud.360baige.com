package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/logger"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/http/mobile/center"
	"time"
	"dev.model.360baige.com/action"
)

type LoggerController struct {
	beego.Controller
}

// @Title 新增add
// @Description post user by uid
// @Param	content		path 	string	true		"内容"
// @Param	remark		path 	string	true		"描述"
// @Param	ip		path 	string	true		"IP地址	"
// @Param	type		path 	int	true		"类别（增、删、改、查）"
// @Param	ownerId		path 	int64	true		"操作者ID"
// @Param	companyId		path 	int64	true		"公司ID"
// @Success 200 {object} models.logger
// @Failure 403 :uid is empty
// @router /add [post]
func (c *LoggerController) Add()  {
	res := LoggerAddResponse{}
	access_token := c.GetString("access_token")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
	}
	//检测 accessToken
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type:  "And",
		Key: "accessToken",
		Val:  access_token,
	})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Messgae = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		com_id := replyAccessToken.CompanyId
		user_id := replyAccessToken.UserId
		user_position_id := replyAccessToken.Id
		user_position_type := replyAccessToken.Type
		if com_id == 0 || user_id == 0 || user_position_id == 0 {
			res.Code = ResponseSystemErr
			res.Messgae = "获取应用信息失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			res := LoggerAddResponse{}
			Type, _ := c.GetInt8("type")
			var reply Logger
			args := Logger{
				CreateTime:       time.Now().Unix(),
				Content:          c.GetString("content"),
				Remark:           c.GetString("remark"),
				Type:             Type,
				UserId:           user_id,
				CompanyId:        com_id,
				UserPositionId:   user_position_id,
				UserPositionType: user_position_type,
			}
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Logger", "Add", args, &reply)
			if err == nil {
				res.Code = ResponseNormal
				res.Messgae = "新增成功"
				res.Data.Id = reply.Id
			} else {
				res.Code = ResponseSystemErr
				res.Messgae = "新增失败"
			}
			c.Data["json"] = res
			c.ServeJSON()
		}
	}

}

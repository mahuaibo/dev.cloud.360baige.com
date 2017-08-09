package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/mobile"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/message"
	"time"
	"dev.model.360baige.com/action"
	"encoding/json"
)

// MessageTemp API
type MessageTempController struct {
	beego.Controller
}

// @Title 消息列表接口
// @Description 消息列表接口
// @Success 200 {"code":200,"messgae":"获取消息列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   user_position_id     query   string true       "user_position_id"
// @Param   current     query   string true       "当前页"
// @Param   page_size     query   string true       "每页数量"
// @Failure 400 {"code":400,"message":"获取消息失败"}
// @router /list [get]
func (c *MessageTempController) List() {
	res := MessageListResponse{}
	access_token := c.GetString("access_token")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Messgae = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
	}
	// 检测 accessToken
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type: "And",
		Key:  "accessToken",
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
		// company_id、user_id、user_position_id、user_position_type
		com_id := replyAccessToken.CompanyId
		user_id := replyAccessToken.UserId
		user_position_id := replyAccessToken.Id
		if com_id == 0 || user_id == 0 || user_position_id == 0 {
			res.Code = ResponseSystemErr
			res.Messgae = "获取信息失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			var reply action.PageByCond
			var cond1 []action.CondValue
			cond1 = append(cond1, action.CondValue{
				Type: "And",
				Key:  "receive_user_id",
				Val:  user_id,
			})
			getuser_position_id, _ := c.GetInt64("user_position_id")

			if getuser_position_id > 0 {
				cond1 = append(cond1, action.CondValue{
					Type: "And",
					Key:  "receive_user_position_id",
					Val:  user_position_id,
				})
			}
			currentPage, _ := c.GetInt64("current")
			pageSize, _ := c.GetInt64("page_size")
			err = client.Call(beego.AppConfig.String("EtcdURL"), "MessageTemp", "PageByCond", &action.PageByCond{
				CondList: cond1,
				Cols:     []string{"id", "create_time", "content", },
				OrderBy:  []string{"id"},
				PageSize: pageSize,
				Current:  currentPage,
			}, &reply)
			if err != nil {
				res.Code = ResponseSystemErr
				res.Messgae = "获取消息信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			} else {

				reply2List := []MessageTemp{}
				err = json.Unmarshal([]byte(reply.Json), &reply2List)
				// List 循环赋值
				for _, value := range reply2List {
					re := time.Unix(value.CreateTime/1000, 0).Format("2006-01-02")
					res.Data.List = append(res.Data.List, MessageListValue{
						Id:         value.Id,
						CreateTime: re,
						Content:    value.Content,
					})
				}
				res.Data.Total = reply.Total
				res.Data.Current = currentPage
				res.Data.CurrentSize = reply.CurrentSize
				res.Data.OrderBy = reply.OrderBy
				res.Data.PageSize = pageSize
				res.Code = ResponseNormal
				res.Messgae = "获取消息信息成功"
				c.Data["json"] = res
				c.ServeJSON()
			}
		}
	}
}

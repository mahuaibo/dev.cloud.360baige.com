package center

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/mobile/center"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/order"
	"time"
	"dev.model.360baige.com/action"
	"encoding/json"
)

type OrderController struct {
	beego.Controller
}

// @Title 订单列表接口
// @Description 订单列表接口
// @Success 200 {"code":200,"message":"获取订单列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   current     query   string true       "当前页"
// @Param   page_size     query   string true       "每页数量"
// @Param   status     query   string true       "订单状态：-2 全部 0:撤回 1：待审核 2：已通过 3：未通过 4：发货中 5：完成"
// @Failure 400 {"code":400,"message":"获取订单列表信息失败"}
// @router /list [get]
func (c *OrderController) List() {
	res := OrderListResponse{}
	access_token := c.GetString("access_token")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Message = "访问令牌无效"
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
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		//company_id、user_id、user_position_id、user_position_type
		com_id := replyAccessToken.CompanyId
		user_id := replyAccessToken.UserId
		user_position_id := replyAccessToken.Id
		user_position_type := replyAccessToken.Type
		if com_id == 0 || user_id == 0 || user_position_id == 0 {
			res.Code = ResponseSystemErr
			res.Message = "获取信息失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			var reply action.PageByCond
			var cond1 []action.CondValue
			cond1 = append(cond1, action.CondValue{
				Type: "And",
				Key:  "company_id",
				Val:  com_id,
			})
			cond1 = append(cond1, action.CondValue{
				Type: "And",
				Key:  "user_id",
				Val:  user_id,
			})
			cond1 = append(cond1, action.CondValue{
				Type: "And",
				Key:  "user_position_id",
				Val:  user_position_id,
			})
			cond1 = append(cond1, action.CondValue{
				Type: "And",
				Key:  "user_position_type",
				Val:  user_position_type,
			})
			status, _ := c.GetInt8("status")
			if status != -2 {
				cond1 = append(cond1,  action.CondValue{
					Type:  "And",
					Key: "status",
					Val:  status,
				})
			} else {
				cond1 = append(cond1, action.CondValue{
					Type:  "And",
					Key: "status__gt",
					Val:  -1,
				})
			}
			currentPage, _ := c.GetInt64("current")
			pageSize, _ := c.GetInt64("page_size")
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "PageByCond", &action.PageByCond{
				CondList:     cond1,
				Cols:     []string{"id", "create_time", "code", "price", "type", "pay_type", "brief", "status"},
				OrderBy:  []string{"id"},
				PageSize: pageSize,
				Current:  currentPage,
			}, &reply)
			if err != nil {
				res.Code = ResponseSystemErr
				res.Message = "获取订单信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			} else {
				replyList := []Order{}
				err = json.Unmarshal([]byte(reply.Json), &replyList)
				//List 循环赋值
				for _, value := range replyList {
					re := time.Unix(value.CreateTime/1000, 0).Format("2006-01-02")
					var rPayType, rStatus string
					if value.PayType == 1 {
						rPayType = "在线支付"
					} else {
						rPayType = "线下支付"
					}
					rStatus = GetStatus(value.Status)
					res.Data.List = append(res.Data.List, OrderValue{
						Id:         value.Id,
						CreateTime: re,
						Code:       value.Code,
						Price:      value.Price,
						Type:       value.Type,
						PayType:    rPayType,
						Brief:      value.Brief,
						Status:     rStatus,
					})
				}
				res.Data.Total = reply.Total
				res.Data.Current = currentPage
				res.Data.CurrentSize = reply.CurrentSize
				res.Data.OrderBy = reply.OrderBy
				res.Data.PageSize = pageSize
				res.Data.Status = status
				res.Code = ResponseNormal
				res.Message = "获取账务统计信息成功"
				c.Data["json"] = res
				c.ServeJSON()

			}
		}
	}
}
func GetStatus(status int8) string {
	var rStatus string
	switch  status {
	case 0:
		rStatus = "撤回"
	case 1:
		rStatus = "待审核"
	case 2:
		rStatus = "已通过"
	case 3:
		rStatus = "未通过"
	case 4:
		rStatus = "发货中"
	default:
		rStatus = "完成"
	}
	return rStatus
}

// @Title 详情接口
// @Description 账务详情接口
// @Success 200 {"code":200,"message":"获取订单详情成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   string true       "id"
// @Failure 400 {"code":400,"message":"获取取订单信息失败"}
// @router /detail [get]
func (c *OrderController) Detail() {
	res := OrderDetailResponse{}
	access_token := c.GetString("access_token")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Message = "访问令牌无效"
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
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		ai_id, _ := c.GetInt64("id")
		if ai_id == 0 {
			res.Code = ResponseSystemErr
			res.Message = "获取信息失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {

			var reply Order
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "FindById", &Order{
				Id: ai_id,
			}, &reply)
			if err != nil {
				res.Code = ResponseSystemErr
				res.Message = "获取信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			} else {
				re := time.Unix(reply.CreateTime/1000, 0).Format("2006-01-02")
				var rPayType, rStatus string
				rStatus = GetStatus(reply.Status)
				if reply.PayType == 1 {
					rPayType = "在线支付"
				} else {
					rPayType = "线下支付"
				}
				res.Code = ResponseSystemErr
				res.Message = "获取账户信息成功"
				res.Data.CreateTime = re
				res.Data.Code = reply.Code
				res.Data.Price = reply.Price
				res.Data.Type = reply.Type
				res.Data.PayType = rPayType
				res.Data.Brief = reply.Brief
				res.Data.Status = rStatus
				c.Data["json"] = res
				c.ServeJSON()
			}
		}
	}
}

// @Title 账务详情接口
// @Description 账务详情接口
// @Success 200 {"code":200,"message":"获取账务详情成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   code     query   string true       "code"
// @Failure 400 {"code":400,"message":"获取账务统计信息失败"}
// @router /detailbycode [get]
func (c *OrderController) DetailByCode() {
	res := OrderDetailResponse{}
	access_token := c.GetString("access_token")
	if access_token == "" {
		res.Code = ResponseSystemErr
		res.Message = "访问令牌无效"
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
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		code := c.GetString("code")
		if code == "" {
			res.Code = ResponseSystemErr
			res.Message = "获取信息失败"
			c.Data["json"] = res
			c.ServeJSON()
		} else {
			var reply Order
			var args2 action.FindByCond
			args2.CondList = append(args2.CondList, action.CondValue{
				Type:  "And",
				Key: "code",
				Val:  code,
			})
			args2.Fileds = []string{"id","create_time", "code", "price", "type","pay_type","brief","status"}
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Order", "FindByCond", args2, &reply)
			if err != nil {
				res.Code = ResponseSystemErr
				res.Message = "获取信息失败"
				c.Data["json"] = res
				c.ServeJSON()
			} else {
				re := time.Unix(reply.CreateTime/1000, 0).Format("2006-01-02")
				var rPayType, rStatus string
				rStatus = GetStatus(reply.Status)
				if reply.PayType == 1 {
					rPayType = "在线支付"
				} else {
					rPayType = "线下支付"
				}
				res.Code = ResponseSystemErr
				res.Message = "获取账户信息成功"
				res.Data.CreateTime = re
				res.Data.Code = reply.Code
				res.Data.Price = reply.Price
				res.Data.Type = reply.Type
				res.Data.PayType = rPayType
				res.Data.Brief = reply.Brief
				res.Data.Status = rStatus
				c.Data["json"] = res
				c.ServeJSON()
			}
		}
	}
}

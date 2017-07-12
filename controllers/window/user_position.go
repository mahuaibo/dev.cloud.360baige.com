package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/user"
	. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/http/window"
	//. "dev.model.360baige.com/http/company"
	"time"
	"fmt"
)

// USER API
type UserPositionController struct {
	beego.Controller
}

// @Title 获取用户身份接口
// @Description 获取用户身份接口
// @Success 200 {"code":200,"messgae":"获取用户身份成功","data":{"accessToken":"ok","expire_in":0}}
// @Param access_ticket     query   string true       "访问票据"
// @Failure 400 {"code":400,"message":"获取用户身份失败"}
// @router /positionlist [get]
func (c *UserPositionController) PositionList() {
	res := UserPositionResponse{}
	access_ticket := c.GetString("access_ticket")

	var replyUser User
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByAccessTicket", &User{
		AccessTicket: access_ticket,
	}, &replyUser)
	fmt.Println(replyUser)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "访问票据无效"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		//判断时效是否超时 TODO
		timestamp := time.Now().UnixNano() / 1e6
		if (replyUser.ExpireIn == 0 || timestamp < replyUser.ExpireIn) {

		} else {
			res.Code = ResponseSystemErr
			res.Messgae = "访问票据超时"
			c.Data["json"] = res
			c.ServeJSON()
		}
	}

	var replyUserPosition *UserPositionPaginator
	var cond1 []CondValue
	cond1 = append(cond1, CondValue{
		Type:  "And",
		Exprs: "user_id",
		Args:  replyUser.Id,
	})
	cond1 = append(cond1, CondValue{
		Type:  "And",
		Exprs: "status__gt",
		Args:  -1,
	})
	fmt.Println(cond1)
	//cond := orm.NewCondition()
	//cond1 := cond.And("user_id__exact", replyUser.Id).And("status__gt", -1)

	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "ListAll", &UserPositionPaginator{
		Cond:     cond1,
		Cols:     []string{"id", "user_id", "company_id", "type", "person_id"},
		OrderBy:  []string{"id"},
		PageSize: 0,
	}, &replyUserPosition)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取用户身份失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		//var idarg []int64
		//for _, value := range replyUserPosition.List {
		//	idarg = append(idarg, value.CompanyId)
		//}
		//fmt.Println(idarg)
		//var replyUserCompany *CompanyPaginator
		//cond2 := orm.NewCondition()
		//cond3 := cond2.And("id__in",idarg)
		//err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "ListAll", &CompanyPaginator{
		//	Cond: cond3,
		//	Cols: []string{"id", "name", "short_name"},
		//	OrderBy:  []string{"id"},
		//	PageSize: 0,
		//}, &replyUserCompany)
		//fmt.Println("dfddfsc-------------")
		//fmt.Println(replyUserCompany)
		//if err != nil {
		res.Code = ResponseNormal
		res.Messgae = "获取用户身份成功"
		//	//循环赋值
		//
		//}else{
		//
		//}
		//
		//fmt.Println(replyUserPosition)
		//res.Data = &replyUserPosition.List
		c.Data["json"] = res
		c.ServeJSON()
	}
}

// @Title 获取一个登录用户身份权限接口
// @Description 获取一个登录用户身份权限接口
// @Success 200 {"code":200,"messgae":"获取一个登录用户身份权限成功","data":{"accessToken":"ok","expire_in":0}}
// @Param user_position_id     query   string true       "用户身份Id"
// @Failure 400 {"code":400,"message":"获取一个登录用户身份权限失败"}
// @router /detail [get]
func (c *UserPositionController) Detail() {
	res := Response{}
	user_position_id, _ := c.GetInt64("user_position_id", 0)

	var replyUserPosition UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindById", &UserPosition{
		Id: user_position_id,
	}, &replyUserPosition)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取用户身份失败"
		c.Data["json"] = res
		c.ServeJSON()
	} else {
		res.Code = ResponseNormal
		res.Messgae = "获取用户身份成功"
		res.Data = replyUserPosition
		c.Data["json"] = res
		c.ServeJSON()
	}
}

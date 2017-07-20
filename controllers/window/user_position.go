package window

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/models/user"
	//. "dev.model.360baige.com/models/response"
	. "dev.model.360baige.com/http/window"
	"time"
	"fmt"
	. "dev.model.360baige.com/models/company"
	"strconv"
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
		if replyUser.ExpireIn == 0 || timestamp < replyUser.ExpireIn {

		} else {
			res.Code = ResponseSystemErr
			res.Messgae = "访问票据超时"
			c.Data["json"] = res
			c.ServeJSON()
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
			//获取公司名称
			var idarg []int64
			idmap := make(map[int64]int64)
			for _, value := range replyUserPosition.List {
				idmap[value.CompanyId] = value.CompanyId
			}
			for _, value := range idmap {
				idarg = append(idarg, value)
			}
			fmt.Println(idarg)
			var cond2 []CondValue
			cond2 = append(cond2, CondValue{
				Type:  "And",
				Exprs: "id__in",
				Args:  idarg,
			})
			var replyUserCompany *CompanyPaginator
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "ListAll", &CompanyPaginator{
				Cond:     cond2,
				Cols:     []string{"id", "name", "short_name", "status"},
				OrderBy:  []string{"id"},
				PageSize: 0,
			}, &replyUserCompany)
			var resData []UserPositionListItem
			//循环赋值
			if err != nil {
				res.Code = ResponseSystemErr
				res.Messgae = "获取用户身份失败"
				c.Data["json"] = res
				c.ServeJSON()
			}
			companyByIds := make(map[int64]Company)
			for _, value := range replyUserCompany.List {
				companyByIds[value.Id] = value
			}
			for _, value := range replyUserPosition.List {
				if (companyByIds[value.CompanyId].Status != -1) {
					resData = append(resData, UserPositionListItem{
						Id:               value.Id,
						CompanyName:      companyByIds[value.CompanyId].Name,
						CompanyShortName: companyByIds[value.CompanyId].ShortName,
						CompanyStatus:    companyByIds[value.CompanyId].Status,
						CompanyId:        value.CompanyId,
						Type:             value.Type,
						PersonId:         value.PersonId,
					})
				}

			}
			res.Code = ResponseNormal
			res.Messgae = "获取用户身份成功"
			res.Data = resData
			fmt.Println(res)
			c.Data["json"] = res
			c.ServeJSON()
		}

	}

}

// @Title 获取一个登录用户身份权限接口
// @Description 获取一个登录用户身份权限接口
// @Success 200 {"code":200,"messgae":"获取一个登录用户身份权限成功","data":{"accessToken":"ok","expire_in":0}}
// @Param user_position_id     query   string true       "用户身份Id"
// @Failure 400 {"code":400,"message":"获取一个登录用户身份权限失败"}
// @router /positiontoken [get]
func (c *UserPositionController) PositionToken() {
	res := UserPositionTokenResponse{}
	user_position_id, _ := c.GetInt64("user_position_id", 0)
	var replyUserPosition UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindById", &UserPosition{
		Id: user_position_id,
	}, &replyUserPosition)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取身份失败"
		c.Data["json"] = res
		c.ServeJSON()
	}
	timestamp := time.Now().UnixNano() / 1e6
	newAccessTicket := strconv.FormatInt(replyUserPosition.Id, 10) + strconv.FormatInt(timestamp, 10)
	replyUserPosition.UpdateTime = timestamp
	replyUserPosition.AccessToken = newAccessTicket //更新token 应该判断时效，再做更新
	err2 := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "UpdateById", replyUserPosition, nil)
	if err2 != nil {
		res.Code = ResponseSystemErr
		res.Messgae = "获取身份失败"
		c.Data["json"] = res
		c.ServeJSON()
	}
	res.Code = ResponseNormal
	res.Messgae = "获取身份成功"
	res.Data.AccessToken = replyUserPosition.AccessToken
	res.Data.ExpireIn = replyUserPosition.ExpireIn
	c.Data["json"] = res
	c.ServeJSON()
}

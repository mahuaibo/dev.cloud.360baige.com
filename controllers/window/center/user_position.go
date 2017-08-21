package center

import (
	. "dev.model.360baige.com/http/window/center"
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/company"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/log"
	"strconv"
	"time"
)

//  UserPosition API
type UserPositionController struct {
	beego.Controller
}

// @Title 获取用户身份接口
// @Description 获取用户身份接口
// @Success 200 {"code":200,"message":"获取用户身份成功"}
// @Param access_ticket     query   string true       "访问票据"
// @Failure 400 {"code":400,"message":"获取用户身份失败"}
// @router /list [post]
func (c *UserPositionController) PositionList() {
	res := UserPositionResponse{}
	access_ticket := c.GetString("access_ticket")

	var replyUser user.User
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "access_ticket", Val: access_ticket },
		},
		Fileds: []string{"id", "expire_in"},
	}, &replyUser)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "系统异常[验证访问票据失败]"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	if replyUser.Id == 0 {
		res.Code = ResponseLogicErr
		res.Message = "访问票据无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var replyUserPosition []user.UserPosition
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "ListByCond", &action.ListByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "user_id", Val: replyUser.Id },
			action.CondValue{Type: "And", Key: "status__gt", Val: -1 },
		},
		Cols:     []string{"id", "user_id", "company_id", "type", "person_id"},
		OrderBy:  []string{"id"},
		PageSize: -1,
	}, &replyUserPosition)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "获取用户身份失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	//获取公司名称
	var company_ids []int64
	for _, value := range replyUserPosition {
		company_ids = append(company_ids, value.CompanyId)
	}
	var replyUserCompany []company.Company
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "ListByCond", &action.ListByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "id__in", Val: company_ids},
		},
		Cols:     []string{"id", "name", "logo"},
		OrderBy:  []string{"id"},
		PageSize: -1,
	}, &replyUserCompany)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "获取用户身份失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var resData []UserPositionListItem
	companys := make(map[int64]company.Company)
	for _, value := range replyUserCompany {
		companys[value.Id] = value
	}
	for _, value := range replyUserPosition {
		if companys[value.CompanyId].Status != -1 {
			resData = append(resData, UserPositionListItem{
				Id:          value.Id,
				Type:        value.Type,
				CompanyLogo: companys[value.CompanyId].Logo,
				CompanyName: companys[value.CompanyId].Name,
			})
		}
	}
	res.Code = ResponseNormal
	res.Message = "获取用户身份成功"
	res.Data = resData
	c.Data["json"] = res
	c.ServeJSON()
}

// @Title 获取访问令牌
// @Description 获取访问令牌
// @Success 200 {"code":200,"message":"获取访问令牌成功"}
// @Param user_position_id     query   string true       "用户身份Id"
// @Param access_ticket     query   string true       "访问票据"
// @Failure 400 {"code":400,"message":"获取访问令牌失败"}
// @router /getAccessToken [post]
func (c *UserPositionController) GetAccessToken() {
	res := UserPositionTokenResponse{}
	user_position_id, _ := c.GetInt64("user_position_id", 0)
	access_ticket := c.GetString("access_ticket", "======")
	log.Println("access_ticket:", access_ticket)

	var replyUser user.User
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "access_ticket", Val: access_ticket },
		},
		Fileds: []string{"id"},
	}, &replyUser)
	log.Println("replyUser:", replyUser)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "系统异常[验证访问票据失败]"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	if replyUser.Id == 0 {
		res.Code = ResponseLogicErr
		res.Message = "访问票据无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	log.Println("通过")
	var replyUserPosition user.UserPosition
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "id", Val: user_position_id },
			action.CondValue{Type: "And", Key: "user_id", Val: replyUser.Id },
			action.CondValue{Type: "And", Key: "status", Val: 0 },
		},
		Fileds: []string{"id", "expire_in"},
	}, &replyUserPosition)

	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "获取访问令牌失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	if replyUserPosition.Id == 0 {
		res.Code = ResponseLogicErr
		res.Message = "访问票据权限不足"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	update_time := time.Now().UnixNano() / 1e6
	newAccessTicket := strconv.FormatInt(replyUserPosition.Id, 10) + strconv.FormatInt(update_time, 10)
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "UpdateById", &action.UpdateByIdCond{
		Id: []int64{user_position_id},
		UpdateList: [] action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: update_time },
			action.UpdateValue{Key: "access_token", Val: newAccessTicket},
		},
	}, nil)
	if err != nil {
		res.Code = ResponseSystemErr
		res.Message = "获取访问令牌失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	res.Code = ResponseNormal
	res.Message = "获取访问令牌成功"
	res.Data.AccessToken = newAccessTicket
	res.Data.ExpireIn = replyUserPosition.ExpireIn
	c.Data["json"] = res
	c.ServeJSON()
	return
}

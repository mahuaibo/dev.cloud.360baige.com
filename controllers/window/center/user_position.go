package center

import (
	. "dev.model.360baige.com/http/window/center"
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/company"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/log"
	"dev.cloud.360baige.com/utils"
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
// @Param accessTicket     query   string true       "访问票据"
// @Failure 400 {"code":400,"message":"获取用户身份失败"}
// @router /list [post]
func (c *UserPositionController) PositionList() {
	type data UserPositionResponse
	accessTicket := c.GetString("accessTicket")

	var replyUser user.User
	err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "access_ticket", Val: accessTicket },
		},
		Fileds: []string{"id", "expire_in"},
	}, &replyUser)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常[验证访问票据失败]"}
		c.ServeJSON()
		return
	}
	if replyUser.Id == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问票据无效"}
		c.ServeJSON()
		return
	}

	var replyUserPosition []user.UserPosition
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "ListByCond", &action.ListByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "user_id", Val: replyUser.Id },
			action.CondValue{Type: "And", Key: "status__gt", Val: -1 },
		},
		OrderBy:  []string{"id"},
		PageSize: -1,
		Cols:     []string{"id", "user_id", "company_id", "type", "person_id"},
	}, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取用户身份失败"}
		c.ServeJSON()
		return
	}

	var company_ids []int64
	for _, value := range replyUserPosition {
		company_ids = append(company_ids, value.CompanyId)
	}
	var replyCompany []company.Company
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "ListByCond", &action.ListByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "id__in", Val: company_ids},
		},
		Cols:     []string{"id", "name", "logo", "type"},
		OrderBy:  []string{"id"},
		PageSize: -1,
	}, &replyCompany)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取用户身份失败"}
		c.ServeJSON()
		return
	}
	var resData []UserPositionListItem
	listOfCompany := make(map[int64]company.Company)
	for _, value := range replyCompany {
		listOfCompany[value.Id] = value
	}
	for _, value := range replyUserPosition {
		if listOfCompany[value.CompanyId].Status != -1 {
			logoUrl := utils.SignURLSample(listOfCompany[value.CompanyId].Logo)
			resData = append(resData, UserPositionListItem{
				UserPositionId:   value.Id,
				UserPositionName: UserPositionName(value.Type),
				CompanyLogo:      logoUrl,
				CompanyName:      listOfCompany[value.CompanyId].Name,
			})
		}
	}
	c.Data["json"] = data{Code: Normal, Message: "获取用户身份成功", Data: resData}
	c.ServeJSON()
	return
}

// @Title 获取访问令牌
// @Description 获取访问令牌
// @Success 200 {"code":200,"message":"获取访问令牌成功"}
// @Param userPositionId     query   string true       "用户身份Id"
// @Param accessType     query   string true       "访问值类型: 0 accessTicket 1 accessToken 默认为 0"
// @Param accessValue     query   string true       "访问值"
// @Failure 400 {"code":400,"message":"获取访问令牌失败"}
// @router /getAccessToken [post]
func (c *UserPositionController) GetAccessToken() {
	type data UserPositionTokenResponse
	userPositionId, _ := c.GetInt64("userPositionId", 0)
	accessType := c.GetString("accessType", "0")
	accessValue := c.GetString("accessValue", "(=@*&%^!)")
	currentTime := time.Now().UnixNano() / 1e6

	var err error
	var replyUserId int64
	if accessType == "1" {
		var replyUserPosition user.UserPosition
		err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
			CondList: []action.CondValue{
				action.CondValue{Type: "And", Key: "access_token", Val: accessValue },
				action.CondValue{Type: "And", Key: "expire_in__gt", Val: currentTime },
			},
			Fileds: []string{"user_id"},
		}, &replyUserPosition)

		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常[验证访问令牌失败]"}
			c.ServeJSON()
			return
		}
		if replyUserPosition.UserId == 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: "访问令牌失败"}
			c.ServeJSON()
			return
		}
		replyUserId = replyUserPosition.UserId

	} else {
		var replyUser user.User
		err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", action.FindByCond{
			CondList: []action.CondValue{
				action.CondValue{Type: "And", Key: "access_ticket", Val: accessValue },
			},
			Fileds: []string{"id"},
		}, &replyUser)
		log.Println("replyUser:", replyUser)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常[验证访问票据失败]"}
			c.ServeJSON()
			return
		}

		if replyUser.Id == 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: "访问票据无效"}
			c.ServeJSON()
			return
		}
		replyUserId = replyUser.Id

	}

	var replyUserPosition user.UserPosition
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "id", Val: userPositionId },
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserId },
			action.CondValue{Type: "And", Key: "status", Val: 0 },
		},
		Fileds: []string{"id", "access_token", "expire_in"},
	}, &replyUserPosition)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "获取访问令牌失败"}
		c.ServeJSON()
		return
	}

	if replyUserPosition.Id == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "访问票据权限不足"}
		c.ServeJSON()
		return
	}

	if currentTime > replyUserPosition.ExpireIn {
		createAccessToken := utils.CreateAccessValue(strconv.FormatInt(replyUserPosition.Id, 10) + "#" + strconv.FormatInt(replyUserPosition.UserId, 10) + "#" + strconv.FormatInt(currentTime, 10))
		expireIn := currentTime + 3600 * 1000 * 24 * 30
		err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "UpdateById", &action.UpdateByIdCond{
			Id: []int64{userPositionId},
			UpdateList: [] action.UpdateValue{
				action.UpdateValue{Key: "expire_in", Val: expireIn },
				action.UpdateValue{Key: "access_token", Val: createAccessToken},
			},
		}, nil)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "获取访问令牌失败"}
			c.ServeJSON()
			return
		} else {
			replyUserPosition.AccessToken = createAccessToken
			replyUserPosition.ExpireIn = expireIn
		}
	}

	c.Data["json"] = data{Code: Normal, Message: "获取访问令牌成功", Data: UserPositionToken{
		AccessToken: replyUserPosition.AccessToken,
		ExpireIn:    replyUserPosition.ExpireIn,
	}}
	c.ServeJSON()
	return
}

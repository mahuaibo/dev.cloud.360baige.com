package authority

import (
	. "dev.model.360baige.com/http/window/authority"
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/company"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/log"
	"dev.cloud.360baige.com/utils"
	"strconv"
	"dev.model.360baige.com/models/application"
	"fmt"
)

var AuthorizeAccessToken string
var AuthorizeExpireIn int64

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
	currentTimestamp := utils.CurrentTimestamp()
	accessType := c.GetString("accessType", "0")
	accessValue := c.GetString("accessValue", "(=@*&%^!)")

	err := utils.Unable(map[string]string{"accessValue": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	var replyUserId int64
	if accessType == "1" {
		replyUserPosition, err := utils.UserPosition(accessValue, currentTimestamp)
		if err != nil {
			c.Data["json"] = data{Code: ErrorPower, Message: err.Error()}
			c.ServeJSON()
			return
		}
		replyUserId = replyUserPosition.UserId
	} else {
		var replyUser user.User
		err := client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", action.FindByCond{
			CondList: []action.CondValue{
				action.CondValue{Type: "And", Key: "access_ticket", Val: accessValue },
			},
			Fileds: []string{"id", "expire_in"},
		}, &replyUser)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}
		if replyUser.Id == 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
			c.ServeJSON()
			return
		}
		replyUserId = replyUser.Id
	}

	var replyUserPosition []user.UserPosition
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "ListByCond", &action.ListByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "user_id", Val: replyUserId },
			action.CondValue{Type: "And", Key: "status__gt", Val: -1 },
		},
		OrderBy:  []string{"id"},
		PageSize: -1,
		Cols:     []string{"id", "user_id", "company_id", "type", "person_id"},
	}, &replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
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
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
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
			logoUrl := utils.SignURLSample(listOfCompany[value.CompanyId].Logo, 60)
			resData = append(resData, UserPositionListItem{
				UserPositionId:   value.Id,
				UserPositionName: user.UserPositionName(value.Type),
				CompanyLogo:      logoUrl,
				CompanyName:      listOfCompany[value.CompanyId].Name,
			})
		}
	}
	fmt.Println("resData", resData)
	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: resData}
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
	currentTimestamp := utils.CurrentTimestamp()
	userPositionId, _ := c.GetInt64("userPositionId", 0)
	accessType := c.GetString("accessType", "0")
	accessValue := c.GetString("accessValue", "(=@*&%^!)")

	err := utils.Unable(map[string]string{"accessValue": "string:true", "userPositionId": "int:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	var replyUserId int64
	if accessType == "1" {
		replyUserPosition, err := utils.UserPosition(accessValue, currentTimestamp)
		if err != nil {
			c.Data["json"] = data{Code: ErrorPower, Message: err.Error()}
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
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}

		if replyUser.Id == 0 {
			c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
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
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	if replyUserPosition.Id == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "FAIL"}
		c.ServeJSON()
		return
	}

	if currentTimestamp > replyUserPosition.ExpireIn {
		createAccessToken := utils.CreateAccessValue(strconv.FormatInt(replyUserPosition.Id, 10) + "#" + strconv.FormatInt(replyUserPosition.UserId, 10) + "#" + strconv.FormatInt(currentTimestamp, 10))
		expireIn := currentTimestamp + user.UserPositionExpireIn
		err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "UpdateById", &action.UpdateByIdCond{
			Id: []int64{userPositionId},
			UpdateList: [] action.UpdateValue{
				action.UpdateValue{Key: "expire_in", Val: expireIn },
				action.UpdateValue{Key: "access_token", Val: createAccessToken},
			},
		}, nil)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}
		replyUserPosition.AccessToken = createAccessToken
		replyUserPosition.ExpireIn = expireIn
	}

	AuthorizeAccessToken = replyUserPosition.AccessToken
	AuthorizeExpireIn = replyUserPosition.ExpireIn
	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: UserPositionToken{
		AccessToken: replyUserPosition.AccessToken,
		ExpireIn:    replyUserPosition.ExpireIn,
	}}
	c.ServeJSON()
	return
}

// @Title 授权
// @Description 授权接口
// @Success 200 {"code":200,"message":"登录成功"}
// @Param   username     query   string true       "用户名：百鸽账号、邮箱、手机号码 三种登录方式"
// @Param   password query   string true       "密码"
// @Failure 400 {"code":400,"message":"登录失败"}
// @router /authorize [post]
func (c *UserPositionController) Authorize() {
	type data AuthorityResponse
	redirectUri := c.GetString("redirectUri")
	err := utils.Unable(map[string]string{"redirectUri": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	var replyApplicationTpl application.ApplicationTpl
	err = client.Call(beego.AppConfig.String("EtcdURL"), "ApplicationTpl", "FindByCond", &action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "site", Val: redirectUri},
		},
	}, &replyApplicationTpl)
	fmt.Println("replyApplicationTpl", replyApplicationTpl)
	if err != nil || replyApplicationTpl.Id == 0 {
		c.Data["json"] = data{Code: ErrorPower, Message: "应用未授权"}
		c.ServeJSON()
		return
	}

	if AuthorizeAccessToken == "" && AuthorizeExpireIn == 0 {
		c.Data["json"] = data{Code: WaitLogin, Message: "等待登陆"}
		c.ServeJSON()
		return
	}

	var replyUserPosition user.UserPosition
	err = client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", &action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "access_token", Val: AuthorizeAccessToken },
			action.CondValue{Type: "And", Key: "status__gt", Val: -1 },
		},
		Fileds:     []string{"id", "user_id", "company_id", "type", "person_id"},
	}, &replyUserPosition)
	fmt.Println("replyUserPosition", replyUserPosition)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyUserPosition.Id == 0 {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	var replyUser user.User
	err = client.Call(beego.AppConfig.String("EtcdURL"), "User", "FindByCond", action.FindByCond{
		CondList: []action.CondValue{
			action.CondValue{Type: "And", Key: "id", Val: replyUserPosition.UserId},
		},
	}, &replyUser)
	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyUser.Id == 0 {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}

	headUrl := utils.SignURLSample(replyUser.Head, 3600)
	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: AuthorityUserPositionToken{
		Head:headUrl,
		Username:replyUser.Username,
		AccessToken: AuthorizeAccessToken,
		ExpireIn:    AuthorizeExpireIn,
	}}
	c.ServeJSON()
	AuthorizeAccessToken = ""
	AuthorizeExpireIn = 0
	return
}
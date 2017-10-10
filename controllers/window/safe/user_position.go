package safe

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
		} else {
			replyUserPosition.AccessToken = createAccessToken
			replyUserPosition.ExpireIn = expireIn
		}
	}

	c.Data["json"] = data{Code: Normal, Message: "SUCCESS", Data: UserPositionToken{
		AccessToken: replyUserPosition.AccessToken,
		ExpireIn:    replyUserPosition.ExpireIn,
	}}
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
// @router /changeAccessToken [post]
func (c *UserPositionController) ChangeAccessToken() {
	type data ChangeAccessTokenResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken", "(=@*&%^!)")
	err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
	if err != nil {
		c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
		c.ServeJSON()
		return
	}

	replyUserPosition, err := utils.UserPosition(accessToken, currentTimestamp)
	if err != nil {
		c.Data["json"] = data{Code: ErrorPower, Message: err.Error()}
		c.ServeJSON()
		return
	}

	if accessToken == replyUserPosition.AccessToken && currentTimestamp >= replyUserPosition.ExpireIn - user.UserPositionTransitExpireIn {
		// 变更
		createAccessToken := utils.CreateAccessValue(strconv.FormatInt(replyUserPosition.Id, 10) + "#" + strconv.FormatInt(replyUserPosition.UserId, 10) + "#" + strconv.FormatInt(currentTimestamp, 10))
		var replyNum action.Num
		client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "UpdateById", &action.UpdateByIdCond{
			Id: []int64{replyUserPosition.Id},
			UpdateList: []action.UpdateValue{
				action.UpdateValue{"update_time", currentTimestamp},
				action.UpdateValue{"transit_token", replyUserPosition.AccessToken},
				action.UpdateValue{"transit_expire_in", replyUserPosition.ExpireIn},
				action.UpdateValue{"access_token", createAccessToken},
				action.UpdateValue{"expire_in", replyUserPosition.ExpireIn + user.UserPositionExpireIn},
			},
		}, &replyNum)
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}
		replyUserPosition.AccessToken = createAccessToken
		replyUserPosition.ExpireIn = replyUserPosition.ExpireIn + user.UserPositionExpireIn
	}

	c.Data["json"] = data{Code: Normal, Message: "SECCESS", Data: UserPositionToken{
		AccessToken: replyUserPosition.AccessToken,
		ExpireIn:    replyUserPosition.ExpireIn,
	}}
	c.ServeJSON()
	return
}

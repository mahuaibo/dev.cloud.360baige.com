package website

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	"dev.model.360baige.com/models/company"
	. "dev.model.360baige.com/http/window/website"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/utils"
)

// Company API
type CompanyController struct {
	beego.Controller
}

// @Title 企业信息接口
// @Description 企业信息接口
// @Success 200 {"code":200,"message":"获取企业信息成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取企业信息失败"}
// @router /detail [post]
func (c *CompanyController) Detail() {
	type data CompanyHeadDetailResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")

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

	var replyCompany company.Company
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "FindById", &company.Company{
		Id: replyUserPosition.CompanyId,
	}, &replyCompany)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常[公司信息失败]"}
		c.ServeJSON()
		return
	}
	logoUrl := utils.SignURLSample(replyCompany.Logo, 3600)
	c.Data["json"] = data{Code: Normal, Message: "获取公司信息成功",
		Data: CompanyDetail{
			Id:         replyCompany.Id,
			Logo:       logoUrl,
			Name:       replyCompany.Name,
			WebsiteImg: replyCompany.WebsiteImg,
		},
	}
	c.ServeJSON()
	return
}

// @Title 企业信息修改接口
// @Description 企业信息修改接口
// @Success 200 {"code":200,"message":"企业信息修改成功"}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"企业信息修改失败"}
// @router /modify [post]
func (c *CompanyController) Modify() {
	type data CompanyModifyResponse
	currentTimestamp := utils.CurrentTimestamp()
	accessToken := c.GetString("accessToken")
	name := c.GetString("name")
	logo := c.GetString("logo")
	websiteImg := c.GetString("websiteImg")

	err := utils.Unable(map[string]string{"accessToken": "string:true", "name": "string:true", "logo": "string:true", "websiteImg": "string:true"}, c.Ctx.Input)
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

	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Company", "UpdateById", action.UpdateByIdCond{
		UpdateList: []action.UpdateValue{
			action.UpdateValue{Key: "update_time", Val: currentTimestamp },
			action.UpdateValue{Key: "name", Val: name},
			action.UpdateValue{Key: "logo", Val: logo},
			action.UpdateValue{Key: "website_img", Val: websiteImg},
		},
		Id: []int64{replyUserPosition.CompanyId},
	}, &replyNum)

	if err != nil {
		c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
		c.ServeJSON()
		return
	}
	if replyNum.Value == 0 {
		c.Data["json"] = data{Code: ErrorLogic, Message: "企业信息修改失败"}
		c.ServeJSON()
		return
	}

	c.Data["json"] = data{Code: Normal, Message: "企业信息修改成功"}
	c.ServeJSON()
	return
}

// @Title 用户注册
// @Description 用户注册
// @Success 200 {"code":200,"message":"企业logo上传失败"}
// @Param   accessToken     query   string true       "方文令牌"
// @Param   id     query   string true       "企业ID"
// @Param   uploadFile query   string true       "file"
// @Failure 400 {"code":400,"message":"企业logo上传成功"}
// @router /uploadImage [options,post]
func (c *CompanyController) UploadImage() {
	currentTimestamp := utils.CurrentTimestamp()
	requestType := c.Ctx.Request.Method
	if requestType == "POST" {
		type data UploadImageResponse
		accessToken := c.GetString("accessToken")
		_, handle, _ := c.Ctx.Request.FormFile("uploadFile")

		err := utils.Unable(map[string]string{"accessToken": "string:true"}, c.Ctx.Input)
		if err != nil {
			c.Data["json"] = data{Code: ErrorLogic, Message: err.Error()}
			c.ServeJSON()
			return
		}

		_, err = utils.UserPosition(accessToken, currentTimestamp)
		if err != nil {
			c.Data["json"] = data{Code: ErrorPower, Message: err.Error()}
			c.ServeJSON()
			return
		}

		objectKey, err := utils.UploadImage(handle, "Company/logoImages/")
		if err != nil {
			c.Data["json"] = data{Code: ErrorSystem, Message: "系统异常，请稍后重试"}
			c.ServeJSON()
			return
		}

		c.Data["json"] = data{Code: Normal, Data: ImageData{
			ObjectKey:objectKey,
			Image: utils.SignURLSample(objectKey, 3600),
		}, Message: "企业logo上传成功"}
		c.ServeJSON()
		return
	}
}

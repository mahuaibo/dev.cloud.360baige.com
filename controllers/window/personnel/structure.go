package personnel

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/personnel"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/personnel"
	"dev.model.360baige.com/action"
	"strings"
	"time"
	"dev.cloud.360baige.com/utils"
)

// Person API
type StructureController struct {
	beego.Controller
}

// @Title 校园收费列表接口
// @Description Structure List 校园收费列表接口
// @Success 200 {"code":200,"Message":"获取成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"获取失败"}
// @router /list [post]
func (c *StructureController) ListOfStructure() {
	res := ListOfStructureResponse{}
	accessToken := c.GetString("accessToken")
	if accessToken == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "access_token", Val: accessToken})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 2.
	var args2 action.PageByCond
	args2.CondList = append(args2.CondList,
		action.CondValue{Type: "And", Key: "company_id", Val: replyAccessToken.CompanyId},
		action.CondValue{Type: "And", Key: "parent_id", Val: 0},
		action.CondValue{Type: "And", Key: "status__gt", Val: -1},
	)
	args2.Cols = []string{"id", "company_id", "parent_id", "name"}
	classList := GetStructureList(args2).List

	if len(classList) <= 0 {
		res.Code = ResponseLogicErr
		res.Message = "获取组织结构失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	classList = append(classList, StructureData{Id: 0, Label: "无组织人员"})
	res.Code = ResponseNormal
	res.Message = "获取组织结构成功"
	res.Data.List = classList
	c.Data["json"] = res
	c.ServeJSON()
	return
}

func GetStructureList(args action.PageByCond) ListOfStructure {
	var list ListOfStructure
	var replyStructure []personnel.Structure
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "ListByCond", args, &replyStructure)
	if err != nil {
		return list
	}

	for _, val := range replyStructure {
		var data StructureData
		data.Id = val.Id
		data.Label = val.Name

		var args2 action.PageByCond
		args2.CondList = append(args2.CondList,
			action.CondValue{Type: "And", Key: "company_id", Val: val.CompanyId},
			action.CondValue{Type: "And", Key: "parent_id", Val: val.Id},
			action.CondValue{Type: "And", Key: "status__gt", Val: -1},
		)
		args2.Cols = []string{"id", "company_id", "parent_id", "name"}
		data.Children = GetStructureList(args2).List
		list.List = append(list.List, data)
	}
	return list
}

// @Title 添加校园收费项目接口
// @Description Structure Add 添加校园收费项目接口
// @Success 200 {"code":200,"Message":"添加成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   name     query   string true       "项目名称"
// @Param   parent_id     query   int64 true       "上级ID"
// @Param   Type     query   int true       "类型 1.班级 2.部门"
// @Failure 400 {"code":400,"message":"添加成功"}
// @router /add [post]
func (c *StructureController) AddStructure() {
	res := AddStructureResponse{}
	accessToken := c.GetString("accessToken")
	nameLists := strings.Split(strings.Replace(c.GetString("name"), "；", ";", -1), ";")

	if accessToken == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "access_token", Val: accessToken})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	//2.
	var count int64
	for _, value := range nameLists {
		parentId, _ := c.GetInt64("parent_id")
		for _, val := range strings.Split(value, ">") {
			// 判断当前节点该结构是否存在
			var args2 action.FindByCond
			args2.CondList = append(args2.CondList,
				action.CondValue{Type: "And", Key: "company_id", Val: replyAccessToken.CompanyId},
				action.CondValue{Type: "And", Key: "name", Val: val},
				action.CondValue{Type: "And", Key: "parent_id", Val: parentId},
				action.CondValue{Type: "And", Key: "status", Val: 1})
			args2.Fileds = []string{"id"}
			var replyStructure personnel.Structure
			err := client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "FindByCond", args2, &replyStructure)
			if err == nil && replyStructure.Id > 0 {
				parentId = replyStructure.Id
				continue
			}
			// 添加组织结构
			operationTime := time.Now().UnixNano() / 1e6
			var args3 personnel.Structure
			args3.CreateTime = operationTime
			args3.UpdateTime = operationTime
			args3.CompanyId = replyAccessToken.CompanyId
			args3.Name = val
			args3.ParentId = parentId
			args3.Status = 1
			var addReplyStructure personnel.Structure
			err = client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "Add", args3, &addReplyStructure)
			if err != nil {
				res.Code = ResponseLogicErr
				res.Message = "添加组织结构失败"
				c.Data["json"] = res
				c.ServeJSON()
				return
			}
			parentId = addReplyStructure.Id
			count++
		}
	}
	res.Code = ResponseNormal
	res.Message = "添加组织结构成功"
	res.Data.Count = count
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 修改校园收费项目接口
// @Description Structure Add 修改校园收费项目接口
// @Success 200 {"code":200,"Message":"修改缴费项目成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   id     query   int64 true       ""
// @Param   name     query   string true       "项目名称"
// @Failure 400 {"code":400,"message":"修改缴费项目失败"}
// @router /modify [post]
func (c *StructureController) ModifyStructure() {
	res := ModifyStructureResponse{}
	accessToken := c.GetString("accessToken")
	id, _ := c.GetInt64("id", 0)
	name := c.GetString("name")
	parentId, _ := c.GetInt64("parentId")

	if accessToken == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "access_token", Val: accessToken})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	// 2.
	args2 := &personnel.Structure{
		Id: id,
	}
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "FindById", args2, args2)
	if args2.CompanyId != replyAccessToken.CompanyId {
		res.Code = ResponseLogicErr
		res.Message = "非法操作"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	// 3.
	args3 := action.UpdateByIdCond{
		Id: []int64{args2.Id},
	}
	operationTime := time.Now().UnixNano() / 1e6
	args3.UpdateList = append(args3.UpdateList,
		action.UpdateValue{Key: "UpdateTime", Val: operationTime},
		action.UpdateValue{Key: "Name", Val: name},
	)

	if parentId > 0 {
		args3.UpdateList = append(args3.UpdateList, action.UpdateValue{Key: "parent_id", Val: parentId})
	}
	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "UpdateById", args3, &replyNum)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "修改失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseNormal
	res.Message = "修改成功"
	res.Data.Count = replyNum.Value
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 删除校园收费记录接口
// @Description Delete Structure 删除校园收费记录接口
// @Success 200 {"code":200,"Message":"删除缴费项目记录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   ids     query   int true       "项目记录ids"
// @Failure 400 {"code":400,"message":"删除缴费项目记录失败"}
// @router /delete [post]
func (c *StructureController) DeleteStructure() {
	res := DeleteStrectureResponse{}
	accessToken := c.GetString("accessToken")
	ids := utils.StrArrToInt64Arr(strings.Split(c.GetString("ids"), ","))
	if accessToken == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "access_token", Val: accessToken})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 2.
	args2 := action.DeleteByIdCond{Value: ids}
	var replyStructure action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "DeleteById", args2, &replyStructure)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "删除组织结构失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	res.Code = ResponseNormal
	res.Message = "删除组织结构成功"
	res.Data.Count = replyStructure.Value
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 获取人员组织结构接口
// @Description Delete Structure 获取人员组织结构接口
// @Success 200 {"code":200,"Message":"获取人员组织结构成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   accessToken     query   string true       "访问令牌"
// @Param   ids     query   int true       "项目记录ids"
// @Failure 400 {"code":400,"message":"获取人员组织结构失败"}
// @router /getStructureIds [post]
func (c *StructureController) GetStructureIds() {
	res := GetStrectureIdsResponse{}
	accessToken := c.GetString("accessToken")
	person_id, _ := c.GetInt64("personId")
	if accessToken == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "access_token", Val: accessToken})
	args.Fileds = []string{"id", "user_id", "company_id", "type"}
	var replyAccessToken user.UserPosition
	err := client.Call(beego.AppConfig.String("EtcdURL"), "UserPosition", "FindByCond", args, &replyAccessToken)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌失效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 2.
	var args2 action.ListByCond
	args2.CondList = append(args2.CondList,
		action.CondValue{Type: "And", Key: "company_id", Val: replyAccessToken.CompanyId},
		action.CondValue{Type: "And", Key: "person_id", Val: person_id},
		action.CondValue{Type: "And", Key: "status", Val: 1})
	args2.Cols = []string{"structure_id"}
	var replyPersonStructure []personnel.PersonStructure
	err = client.Call(beego.AppConfig.String("EtcdURL"), "PersonStructure", "ListByCond", args2, &replyPersonStructure)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "获取组织结构失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var ids []int64 = make([]int64, len(replyPersonStructure))
	for key, val := range replyPersonStructure {
		ids[key] = val.StructureId
	}

	res.Code = ResponseNormal
	res.Message = "获取组织结构成功"
	res.Data.Ids = ids
	c.Data["json"] = res
	c.ServeJSON()
	return
}

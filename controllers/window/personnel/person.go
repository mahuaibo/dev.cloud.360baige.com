package personnel

import (
	"github.com/astaxie/beego"
	"dev.cloud.360baige.com/rpc/client"
	. "dev.model.360baige.com/http/window/personnel"
	"dev.model.360baige.com/models/user"
	"dev.model.360baige.com/models/personnel"
	"dev.model.360baige.com/action"
	"dev.cloud.360baige.com/utils"
	"fmt"
	"strings"
	"time"
	"os"
	"io"
	"github.com/xuri/excelize"
	"strconv"
	"github.com/tealeg/xlsx"
	"net/http"
	"encoding/json"
)

// Person API
type PersonController struct {
	beego.Controller
}

// @Title 校园收费记录列表接口
// @Description Project List 校园收费记录列表接口
// @Success 200 {"code":200,"message":"获取人员列表成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   structure_id     query   int63 true       "组织结构ID"
// @Param   page_size     query   int64 true       "每页显示条数"
// @Param   current     query   int64 true       "页码"
// @Failure 400 {"code":400,"message":"获取人员列表失败"}
// @router /list [post]
func (c *PersonController) ListOfPerson() {
	res := ListOfPersonResponse{}
	access_token := c.GetString("access_token")
	structureId, _ := c.GetInt64("structure_id")
	pageSize, _ := c.GetInt64("page_size")
	currentPage, _ := c.GetInt64("current")

	if access_token == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "access_token", Val: access_token})
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
	// 查询组织结构下 人员ids
	var args2 action.PageByCond
	if structureId != 0 {
		args2.CondList = append(args2.CondList, action.CondValue{Type: "And", Key: "structure_id", Val: structureId})
	}
	args2.CondList = append(args2.CondList, action.CondValue{Type: "And", Key: "company_id", Val: replyAccessToken.CompanyId})
	args2.Cols = []string{"person_id"}
	var replyPersonStructure []personnel.PersonStructure
	err = client.Call(beego.AppConfig.String("EtcdURL"), "PersonStructure", "ListByCond", args2, &replyPersonStructure)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "获取人员列表失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	var personIds []int64 = make([]int64, len(replyPersonStructure))
	for key, val := range replyPersonStructure {
		personIds[key] = val.PersonId
	}
	if len(personIds) <= 0 {
		res.Code = ResponseNormal
		res.Message = "当前结构暂无人员！"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 查询该公司全部人员ID
	if (structureId == 0) {
		var args3 action.PageByCond
		args3.CondList = append(args3.CondList,
			action.CondValue{Type: "And", Key:  "company_id", Val:  replyAccessToken.CompanyId},
			action.CondValue{Type: "And", Key:  "status__gt", Val:  -1},
		)
		args3.Cols = []string{"id"}
		var replyPersonIdList []personnel.Person
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Person", "ListByCond", args3, &replyPersonIdList)
		var ids []int64 = make([]int64, len(replyPersonIdList))
		for key, val := range replyPersonIdList {
			ids[key] = val.Id
		}
		personIds = utils.Minus(ids, personIds)
	}
	// 查询人员
	var args4 action.PageByCond
	args4.CondList = append(args4.CondList,
		action.CondValue{Type: "And", Key: "id__in", Val: personIds},
		action.CondValue{Type: "And", Key: "company_id", Val: replyAccessToken.CompanyId},
		action.CondValue{Type: "And", Key: "status__gt", Val: -1})
	args4.OrderBy = []string{"id"}
	args4.Cols = []string{"id", "create_time", "update_time", "company_id", "code", "name", "sex", "birthday", "type", "phone", "contact", "status" }
	args4.PageSize = pageSize
	args4.Current = currentPage
	var replyPerson []personnel.Person
	fmt.Println("args4", args4)
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Person", "ListByCond", args4, &replyPerson)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "获取人员列表失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 数据拼接
	listOfPerson := make([]Person, len(replyPerson), len(replyPerson))
	for index, rec := range replyPerson {
		listOfPerson[index] = Person{
			Id:         rec.Id,
			CreateTime: time.Unix(rec.CreateTime/1000, 0).Format("2006-01-02"),
			CompanyId:  rec.CompanyId,
			Code:       rec.Code,
			Name:       rec.Name,
			Sex:        rec.Sex,
			Birthday:   time.Unix(rec.CreateTime/1000, 0).Format("2006-01-02"),
			Type:       rec.Type,
			Phone:      rec.Phone,
			Contact:    rec.Contact,
			Status:     rec.Status,
		}
	}
	res.Code = ResponseNormal
	res.Message = "获取人员列表成功"
	res.Data.List = listOfPerson
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 添加人员接口
// @Description Person Add 添加收费名单接口
// @Success 200 {"code":200,"message":"添加收费名单成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   name     query   string true       "姓名"
// @Param   sex     query   string true       "性别"
// @Param   type     query   int8 true       "职位"
// @Param   phone     query   string true       "手机号码"
// @Param   birthday     query   int64 true       "生日"
// @Param   contact     query   string true       "联系人"
// @Param   structureId     query   int64 true       "组织结构ID"
// @Failure 400 {"code":400,"message":"添加收费名单失败"}
// @router /add [post]
func (c *PersonController) AddPerson() {
	res := AddPersonResponse{}
	access_token := c.GetString("access_token")
	code := c.GetString("code")
	name := c.GetString("name")
	sex := c.GetString("sex")
	Type, _ := c.GetInt8("type")
	phone := c.GetString("phone")
	birthday, _ := c.GetInt64("birthday")
	contact := c.GetString("contact")
	structureIds := utils.StrArrToInt64Arr(strings.Split(c.GetString("structure_ids"), ","))

	fmt.Println("birthday", birthday)
	if access_token == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "access_token", Val: access_token})
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
	// 添加人员
	operationTime := time.Now().UnixNano() / 1e6
	args2 := &personnel.Person{
		CreateTime: operationTime,
		UpdateTime: operationTime,
		CompanyId:  replyAccessToken.CompanyId,
		Code:       code,
		Name:       name,
		Sex:        sex,
		Birthday:   birthday,
		Type:       Type,
		Phone:      phone,
		Contact:    contact,
		Status:     1,
	}
	var replyPerson personnel.Person
	fmt.Println("args2", args2)
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Person", "Add", args2, &replyPerson)
	fmt.Println("replyPerson", replyPerson)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "添加人员失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	fmt.Println("structureIds", structureIds)
	if len(structureIds) > 0 {
		var args3 []personnel.PersonStructure = make([]personnel.PersonStructure, len(structureIds))
		for key, val := range structureIds {
			args3[key] = personnel.PersonStructure{
				CreateTime:  operationTime,
				UpdateTime:  operationTime,
				CompanyId:   replyAccessToken.CompanyId,
				PersonId:    replyPerson.Id,
				StructureId: val,
				Type:        1,
				Status:      1,
			}
		}
		// 添加组织结构
		var replyPersonStructure personnel.PersonStructure
		fmt.Println("args3", args3)
		err = client.Call(beego.AppConfig.String("EtcdURL"), "PersonStructure", "AddMultiple", args3, &replyPersonStructure)
		fmt.Println("replyPersonStructure", replyPersonStructure)
		if err != nil {
			res.Code = ResponseLogicErr
			res.Message = "组织结构添加失败"
			c.Data["json"] = res
			c.ServeJSON()
			return
		}
	}

	res.Code = ResponseNormal
	res.Message = "添加人员成功"
	res.Data.Id = replyPerson.Id
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 修改缴费记录接口
// @Description Project Add 修改缴费记录接口
// @Success 200 {"code":200,"message":"修改缴费记录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   id     query   int64 true       "员工id"
// @Param   name     query   string true       "姓名"
// @Param   sex     query   string true       "性别"
// @Param   type     query   int8 true       "职位"
// @Param   phone     query   string true       "手机号码"
// @Param   birthday     query   int64 true       "生日"
// @Param   contact     query   string true       "联系人"
// @Param   structureId     query   int64 true       "组织结构ID"
// @Failure 400 {"code":400,"message":"修改缴费记录失败"}
// @router /modify [post]
func (c *PersonController) ModifyPerson() {
	res := ModifyPersonResponse{}
	access_token := c.GetString("access_token")
	id, _ := c.GetInt64("id")
	code := c.GetString("code")
	name := c.GetString("name")
	sex := c.GetString("sex")
	Type, _ := c.GetInt8("type")
	phone := c.GetString("phone")
	birthday := c.GetString("birthday")
	contact := c.GetString("contact")
	structureIds := utils.StrArrToInt64Arr(strings.Split(c.GetString("structure_ids"), ","))

	if access_token == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{
		Type: "And",
		Key:  "access_token",
		Val:  access_token,
	})
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
	// 查询人员信息判断有效性
	args2 := &personnel.Person{
		Id: id,
	}
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Person", "FindById", args2, args2)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "修改人员信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	if args2.CompanyId != replyAccessToken.CompanyId {
		res.Code = ResponseLogicErr
		res.Message = "非法操作"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	// 修改人员信息
	args3 := action.UpdateByIdCond{}
	args3.Id = []int64{args2.Id}
	args3.UpdateList = append(args3.UpdateList,
		action.UpdateValue{Key: "update_time", Val: time.Now().UnixNano() / 1e6},
		action.UpdateValue{Key: "code", Val: code},
		action.UpdateValue{Key: "name", Val: name},
		action.UpdateValue{Key: "sex", Val: sex},
		action.UpdateValue{Key: "birthday", Val: birthday},
		action.UpdateValue{Key: "type", Val: Type},
		action.UpdateValue{Key: "phone", Val: phone},
		action.UpdateValue{Key: "contact", Val: contact},
		action.UpdateValue{Key: "type", Val: Type},
	)
	var replyNum action.Num
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Person", "UpdateById", args3, &replyNum)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "修改人员信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}

	// 查询原有人员组织
	var args4 action.ListByCond
	args4.CondList = append(args4.CondList,
		action.CondValue{Type: "And", Key: "company_id", Val: replyAccessToken.CompanyId},
		action.CondValue{Type: "And", Key: "person_id", Val: args2.Id},
		action.CondValue{Type: "And", Key: "status", Val: 1})
	args4.Cols = []string{"structure_id"}
	var replyPersonStructure []personnel.PersonStructure
	err = client.Call(beego.AppConfig.String("EtcdURL"), "PersonStructure", "ListByCond", args4, &replyPersonStructure)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "人员组织添加失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	var ids []int64
	for _, value := range replyPersonStructure {
		ids = append(ids, value.StructureId)
	}

	// 修改人员组织信息
	modifyStructureIds := utils.Minus(ids, structureIds)
	if len(modifyStructureIds) > 0 {
		args5 := action.DeleteByCond{}
		args5.CondList = append(args5.CondList,
			action.CondValue{Type: "And", Key: "structure_id__in", Val: modifyStructureIds},
			action.CondValue{Type: "And", Key: "person_id", Val: id})
		var replyPerson action.Num
		err = client.Call(beego.AppConfig.String("EtcdURL"), "PersonStructure", "DeleteByCond", args5, &replyPerson)
		if err != nil {
			res.Code = ResponseLogicErr
			res.Message = "修改人员组织失败"
			c.Data["json"] = res
			c.ServeJSON()
			return
		}
	}
	// 新增人员组织信息
	addStructureIds := utils.Minus(structureIds, ids)
	if len(addStructureIds) > 0 {
		operationTime := time.Now().UnixNano() / 1e6
		var args6 []personnel.PersonStructure = make([]personnel.PersonStructure, len(addStructureIds))
		for key, val := range addStructureIds {
			args6[key] = personnel.PersonStructure{
				CreateTime:  operationTime,
				UpdateTime:  operationTime,
				CompanyId:   args2.CompanyId,
				PersonId:    args2.Id,
				StructureId: val,
				Type:        1,
				Status:      1,
			}
		}
		var replyPersonStructure1 personnel.PersonStructure
		err = client.Call(beego.AppConfig.String("EtcdURL"), "PersonStructure", "AddMultiple", args6, &replyPersonStructure1)
		if err != nil {
			res.Code = ResponseLogicErr
			res.Message = "修改人员组织失败"
			c.Data["json"] = res
			c.ServeJSON()
			return
		}
	}

	res.Code = ResponseNormal
	res.Message = "修改人员信息成功"
	res.Data.Count = replyNum.Value
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 删除校园收费记录接口
// @Description Delete Person 删除校园收费记录接口
// @Success 200 {"code":200,"message":"删除缴费项目记录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   ids     query   string true       "员工ids"
// @Failure 400 {"code":400,"message":"删除缴费项目记录失败"}
// @router /delete [post]
func (c *PersonController) DeletePerson() {
	res := DeletePersonResponse{}
	access_token := c.GetString("access_token")
	personIds := utils.StrArrToInt64Arr(strings.Split(c.GetString("person_ids"), ","))
	fmt.Println("personIds", personIds)
	fmt.Println("111111", c.GetString("person_ids"))
	if access_token == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "access_token", Val: access_token, })
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
	args2 := action.DeleteByIdCond{Value: personIds}
	var replyPerson action.Num
	fmt.Println("args2", args2)
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Person", "DeleteById", args2, &replyPerson)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "删除人员失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	fmt.Println("2:", replyPerson)
	res.Code = ResponseNormal
	res.Message = "删除人员成功"
	res.Data.Count = replyPerson.Value
	c.Data["json"] = res
	c.ServeJSON()
	return
}

// @Title 导入人员接口
// @Description Delete Person 导入人员接口
// @Success 200 {"code":200,"message":"导出成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Failure 400 {"code":400,"message":"导入失败"}
// @router /upload [options,post]
func (c *PersonController) UploadPerson() {
	res := UploadPersonResponse{}
	access_token := c.GetString("access_token")
	if access_token == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "access_token", Val: access_token })
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

	var count int64 = 0
	requestType := c.Ctx.Request.Method
	if requestType == "POST" {
		// 获取file文件
		formFile, header, err := c.Ctx.Request.FormFile("uploadFile")
		if err != nil {
			fmt.Println("Get form file failed: %s\n", err)
			return
		}
		defer formFile.Close()
		objectKey := "./" + strconv.FormatInt(time.Now().UnixNano()/1e6, 10) + header.Filename
		// 创建保存文件
		destFile, err := os.Create(objectKey)
		if err != nil {
			fmt.Println("Create failed: %s\n", err)
			return
		}
		defer destFile.Close()
		// 读取表单文件，写入保存文件
		_, err = io.Copy(destFile, formFile)
		if err != nil {
			fmt.Println("Write file failed: %s\n", err)
			return
		}
		// 解析文件
		xlsx, err := excelize.OpenFile(objectKey)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// 读取文件内容
		rows := xlsx.GetRows("sheet" + strconv.Itoa(xlsx.GetSheetIndex("Sheet1")))
		for key, row := range rows {
			if key == 0 {
				continue
			}
			// 添加人员
			personId, err := addPerson(row, replyAccessToken.CompanyId)
			if err != nil {
				res.Code = ResponseLogicErr
				res.Message = "上传人员失败"
				c.Data["json"] = res
				c.ServeJSON()
				return
			}
			// 添加组织结构
			nameLists := strings.Split(strings.Replace(row[6], "；", ";", -1), ";")
			if len(nameLists) <= 0 {
				continue
			}
			for _, value := range nameLists {
				var parentId int64 = 0
				for _, val := range strings.Split(value, ">") {
					// 判断当前节点该结构是否存在
					var args3 action.FindByCond
					args3.CondList = append(args3.CondList,
						action.CondValue{Type: "And", Key:  "company_id", Val:  replyAccessToken.CompanyId},
						action.CondValue{Type: "And", Key:  "name", Val:  val},
						action.CondValue{Type: "And", Key:  "parent_id", Val:  parentId},
						action.CondValue{Type: "And", Key:  "status", Val:  1})
					args3.Fileds = []string{"id"}
					var replyStructure personnel.Structure
					err := client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "FindByCond", args3, &replyStructure)
					if err == nil && replyStructure.Id > 0 {
						parentId = replyStructure.Id
						continue
					}
					// 添加组织结构
					structureId, err := addStructure(replyAccessToken.CompanyId, val, parentId)
					if err != nil {
						res.Code = ResponseLogicErr
						res.Message = "添加组织结构失败"
						c.Data["json"] = res
						c.ServeJSON()
						return
					}
					parentId = structureId
				}
				// 添加人员组织
				err = addPersonStructure(replyAccessToken.CompanyId, personId, parentId)
				if err != nil {
					res.Code = ResponseLogicErr
					res.Message = "修改人员组织失败"
					c.Data["json"] = res
					c.ServeJSON()
					return
				}
			}
			count ++
		}

		err = os.Remove(objectKey) // 删除文件
		if err != nil {
			fmt.Println("file remove Error!", err)
		} else {
			fmt.Print("file remove OK!")
		}
	}

	// 3.
	res.Code = ResponseNormal
	res.Message = "上传人员成功"
	res.Data.Count = count
	c.Data["json"] = res
	c.ServeJSON()
	return
}
// 添加人员
func addPerson(data []string, companyId int64) (int64, error) {
	timestamp := time.Now().UnixNano() / 1e6
	var Type int8 = 1
	if data[4] == "学生" {
		Type = 2
	}
	var birthday int64 = 0
	if data[3] != "" {
		time, _ := time.Parse("2006-01-02", data[3])
		birthday = time.Unix() * 1000
	}
	// 添加人员
	var args2 personnel.Person
	args2.CreateTime = timestamp
	args2.UpdateTime = timestamp
	args2.CompanyId = companyId
	args2.Code = data[0]
	args2.Name = data[1]
	args2.Sex = data[2]
	args2.Birthday = birthday
	args2.Type = Type
	args2.Phone = data[5]
	args2.Status = 1
	var replyPerson personnel.Person
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Person", "Add", args2, &replyPerson)
	return replyPerson.Id, err
}
// 添加组织
func addStructure(companyId int64, name string, parentId int64) (int64, error) {
	timestamp := time.Now().UnixNano() / 1e6
	var args4 personnel.Structure
	args4.CreateTime = timestamp
	args4.UpdateTime = timestamp
	args4.CompanyId = companyId
	args4.Name = name
	args4.ParentId = parentId
	args4.Status = 1
	var replyStructure personnel.Structure
	err := client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "Add", args4, &replyStructure)
	return replyStructure.Id, err
}
// 添加人员组织
func addPersonStructure(companyId int64, personId int64, structureId int64) error {
	timestamp := time.Now().UnixNano() / 1e6
	var args personnel.PersonStructure
	args.CreateTime = timestamp
	args.UpdateTime = timestamp
	args.CompanyId = companyId
	args.PersonId = personId
	args.StructureId = structureId
	args.Type = 1
	args.Status = 1
	var replyPersonStructure personnel.PersonStructure
	return client.Call(beego.AppConfig.String("EtcdURL"), "PersonStructure", "Add", args, &replyPersonStructure)
}

// @Title 下载缴费记录接口
// @Description Delete Person 下载缴费记录接口
// @Success 200 {"code":200,"message":"下载缴费记录成功","data":{"access_ticket":"xxxx","expire_in":0}}
// @Param   access_token     query   string true       "访问令牌"
// @Param   Person_ids     query   int true       "项目记录IDs"
// @Failure 400 {"code":400,"message":"下载缴费记录失败"}
// @router /download [get,post]
func (c *PersonController) DownloadPerson() {
	res := DownloadPersonResponse{}
	access_token := c.GetString("access_token")
	structureId, _ := c.GetInt64("structureId")
	if access_token == "" {
		res.Code = ResponseLogicErr
		res.Message = "访问令牌无效"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 1.
	var args action.FindByCond
	args.CondList = append(args.CondList, action.CondValue{Type: "And", Key: "access_token", Val: access_token })
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

	// 查询组织结构下 人员ids
	var args2 action.PageByCond
	if structureId != 0 {
		args2.CondList = append(args2.CondList, action.CondValue{Type: "And", Key: "structure_id", Val: structureId})
	}
	args2.CondList = append(args2.CondList, action.CondValue{Type: "And", Key: "company_id", Val: replyAccessToken.CompanyId})
	args2.Cols = []string{"person_id", "structure_id"}
	var replyPersonStructure []personnel.PersonStructure
	err = client.Call(beego.AppConfig.String("EtcdURL"), "PersonStructure", "ListByCond", args2, &replyPersonStructure)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "下载人员信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// TODO 去重
	var personIds []string = make([]string, len(replyPersonStructure))
	for key, val := range replyPersonStructure {
		personIds[key] = strconv.FormatInt(val.PersonId, 10)
	}
	// 2.
	args3 := action.ListByCond{CondList: []action.CondValue{
		action.CondValue{Type: "And", Key: "company_id", Val: replyAccessToken.CompanyId},
		action.CondValue{Type: "And", Key: "id__in", Val: RemoveDuplicatesAndEmpty(personIds)},
	}}
	var replyPerson []personnel.Person
	err = client.Call(beego.AppConfig.String("EtcdURL"), "Person", "ListByCond", args3, &replyPerson)
	if err != nil {
		res.Code = ResponseLogicErr
		res.Message = "下载人员信息失败"
		c.Data["json"] = res
		c.ServeJSON()
		return
	}
	// 3.
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	row := sheet.AddRow()
	row.SetHeightCM(1) //设置每行的高度
	row.AddCell().Value = "编号"
	row.AddCell().Value = "姓名"
	row.AddCell().Value = "性别"
	row.AddCell().Value = "职位"
	row.AddCell().Value = "联系电话"
	row.AddCell().Value = "组织结构"
	//编号,姓名、性别、职位、联系电话、组织结构
	for _, rec := range replyPerson {
		row := sheet.AddRow()
		row.SetHeightCM(1) //设置每行的高度
		row.AddCell().Value = rec.Code
		row.AddCell().Value = rec.Name
		row.AddCell().Value = rec.Sex
		var Type string
		if rec.Type == 1 {
			Type = "教师"
		} else {
			Type = "学生"
		}
		row.AddCell().Value = Type
		row.AddCell().Value = rec.Phone
		// 获取人员所在组织架构名称
		var args4 action.ListByCond
		args4.CondList = append(args4.CondList, action.CondValue{Type: "And", Key: "person_id", Val: rec.Id})
		args4.Cols = []string{"structure_id"}
		var replyPersonStructure []personnel.PersonStructure
		err = client.Call(beego.AppConfig.String("EtcdURL"), "PersonStructure", "ListByCond", args4, &replyPersonStructure)
		fmt.Println("replyPersonStructure", replyPersonStructure)
		if err != nil {
			res.Code = ResponseLogicErr
			res.Message = "下载人员信息失败"
			c.Data["json"] = res
			c.ServeJSON()
			return
		}
		var StructureIds []int64 = make([]int64, len(replyPersonStructure))
		for key, val := range replyPersonStructure {
			StructureIds[key] = val.StructureId
		}
		var args5 action.ListByCond
		args5.CondList = append(args5.CondList, action.CondValue{Type: "And", Key: "id__in", Val: StructureIds})
		args5.Cols = []string{"name"}
		var replyStructure []personnel.Structure
		err = client.Call(beego.AppConfig.String("EtcdURL"), "Structure", "ListByCond", args5, &replyStructure)
		if err != nil {
			res.Code = ResponseLogicErr
			res.Message = "下载人员信息失败"
			c.Data["json"] = res
			c.ServeJSON()
			return
		}
		var structureNames []string = make([]string, len(replyStructure))
		for key, val := range replyStructure {
			structureNames[key] = val.Name
		}
		names, _ := json.Marshal(structureNames)
		row.AddCell().Value = string(names)
	}
	objectKey := strconv.FormatInt(time.Now().UnixNano()/1e6, 10) + "file.xlsx"
	err = file.Save(objectKey)
	if err != nil {
		panic(err)
	}

	c.Ctx.Output.Header("Accept-Ranges", "bytes")
	c.Ctx.Output.Header("Content-Disposition", "attachment; filename="+fmt.Sprintf("%s", objectKey)) //文件名
	c.Ctx.Output.Header("Cache-Control", "must-revalidate, post-check=0, pre-check=0")
	c.Ctx.Output.Header("Pragma", "no-cache")
	c.Ctx.Output.Header("Expires", "0")
	//最主要的一句
	http.ServeFile(c.Ctx.ResponseWriter, c.Ctx.Request, objectKey)
	res.Code = ResponseNormal
	res.Message = "下载缴费记录成功"
	c.Data["json"] = res
	c.ServeJSON()
	return
}

func RemoveDuplicatesAndEmpty(a []string) (ret []string) {
	a_len := len(a)
	for i := 0; i < a_len; i++ {
		if (i > 0 && a[i-1] == a[i]) || len(a[i]) == 0 {
			continue;
		}
		ret = append(ret, a[i])
	}
	return
}

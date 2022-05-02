/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
* @HomePage http://www.sysadm.cn
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at:
* http://www.apache.org/licenses/LICENSE-2.0
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and  limitations under the License.
* @License GNU Lesser General Public License  https://www.sysadm.cn/lgpl.html
 */

package server

import (
	//	"encoding/json"
	//	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/wangyysde/sysadm/db"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadm/utils"
	"github.com/wangyysde/sysadmServer"
)


var  projectActions = []string{"list","add","getcount","del","getinfo"}


 func (p Project)ModuleName()string{
	return "project"
}

func (p Project) ActionHanderCaller(action string, c *sysadmServer.Context){
	switch action{
		case "list":
			p.listHandler(c)
		case "add":
			p.addHandler(c)
		case "getcount":
			p.getCountHandler(c)
		case "del":
			p.delHandler(c)
		case "getinfo":
			p.getInfoHandler(c)
	}
	
}

/* 
	listHandler list project information according to start and num provided by rquest's URL
	response the client with Status: false, Erorrcode: int, and Message: string if operation is failed
	otherwise response the client with Status: true, Erorrcode: 0, and Message: "list of project" if operation is successful
*/
func (p Project) listHandler(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080001,"debug","now handling project list handler through api."))
	conditionKey, _ := c.GetQuery("conditionKey")
	conditionValue, _ := c.GetQuery("conditionValue")
	where := make(map[string]string)
	if strings.TrimSpace(conditionKey) != "" && strings.TrimSpace(conditionValue) != ""{
		conditionKey = strings.ToLower(strings.TrimSpace(conditionKey))
		conditionValue = strings.ToLower(strings.TrimSpace(conditionValue))
		if strings.EqualFold(conditionKey,"name") || strings.EqualFold(conditionKey,"comment"){
			where[conditionKey] = " like '%" + conditionValue + "%'"
		}else{
			where[conditionKey] = "=" + conditionValue
		}
	}
	
	deleted, _ := c.GetQuery("deleted")
	deleted = strings.ToLower(strings.TrimSpace(deleted))
	if deleted == "n" || deleted == "0" || deleted == "" {
		where[deleted] = "deleted=0 " 
	}

	start, _ := c.GetQuery("start")
	num, _ := c.GetQuery("num")
	start = strings.ToLower(strings.TrimSpace(start))
	num = strings.ToLower(strings.TrimSpace(num))
	var limit []int
	if start != "" {
		tmpStart, _ := strconv.Atoi(start)
		limit = append(limit, tmpStart)
		if num != "" {
			tmpNum, _ := strconv.Atoi(num)
			limit = append(limit, tmpNum)
		}
	}
	
	orderField, _ := c.GetQuery("orderfield")
	order, _ := c.GetQuery("order")
	orderField = strings.ToLower(strings.TrimSpace(orderField))
	order = strings.ToLower(strings.TrimSpace(order))
	orderInt := 1
	if order == "" || order != "desc" {
		orderInt = 0
	} 
	orderData := make([]db.OrderData,0)
	if orderField != ""  {
		if(p.isField(orderField)) {
			tmpData := db.OrderData{Key: orderField, Order: orderInt}
			orderData = append(orderData, tmpData)
		}
	}

	// Qeurying data from DB
	selectData := db.SelectData{
		Tb: []string{"project"},
		OutFeilds: []string{"projectid","ownerid","name","comment","deleted","creation_time","update_time",},
		Where: where,
		Order: orderData,
		Limit: limit,
	}
	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	retData,err := dbEntity.QueryData(&selectData)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error"){
		errs = append(errs,err...)
		logErrors(errs)
		ret := buildResponse(1080002,false,"database query error")
		c.JSON(http.StatusOK, ret)
		return 
	} 

	// no project has be queried. 
	if len(retData) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080101,"debug","no project has be queried."))
		logErrors(errs)
		ret := buildResponse(1080101,false,"no project has be queried.")
		c.JSON(http.StatusOK, ret)
		return 
	}
	
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080003,"debug","send response data to the client."))
	logErrors(errs)
	ret := ApiResponseStatus {
		Status: true,
		Errorcode: 1080006,
		Message: retData,
	}
	c.JSON(http.StatusOK, ret)
	
}

func (p Project) isField(name string) bool {
	obj := reflect.TypeOf(p)
	if obj.Kind() == reflect.Ptr {
		obj = obj.Elem()
	}

	found := false
	for i :=0 ; i < obj.NumField() ; i++ {
		fieldName := obj.Field(i).Name
		if strings.EqualFold(fieldName,name) {
			found = true
			return found 
		}
	}

	return found
}

/* 
	addHandler get parametes from the request and check their validity
	insert the data of the project into the DB if they are validly. 
*/
func (p Project) addHandler(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080003,"debug","now try to add project infromation into the DB."))
	keys := []string{"name","comment","ownerid"}
	datas,err := utils.GetRequestData(c,keys)
	errs = append(errs,err...)
	logErrors(errs)
	errs = errs[0:0]

	if datas == nil || datas["name"] == "" || datas["ownerid"] == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080004,"error","data is not valid."))
		logErrors(errs)
		ret := buildResponse(1080004,false,"数据错误！")
		c.JSON(http.StatusOK, ret)
		return
	}
	
	where := make(map[string]string)
	where["name"] = "=" + "\"" + datas["name"] + "\""
	outField := "count(\"name\") as num"
	selectData := db.SelectData{
		Tb: []string{"project"},
		OutFeilds: []string{outField},
		Where: where,
	}
	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	retData,err := dbEntity.QueryData(&selectData)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error"){
		errs = append(errs,err...)
		logErrors(errs)
		ret := buildResponse(108001008,false,"数据查询错误，稍后请重试！")
		c.JSON(http.StatusOK, ret)
		return 
	}
	line := retData[0]
	numStr := utils.Interface2String(line["num"])
	numInt,_ := strconv.Atoi(numStr)
	if numInt > 0  {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080014,"error","系统中不能有同名的项目信息"))
		logErrors(errs)
		ret := buildResponse(1080014,false,"系统中不能有同名的项目信息！")
		c.JSON(http.StatusOK, ret)
		return
	} 
	creation_time := time.Now().Unix()
	update_time := creation_time
	insertData := make(map[string]interface{})
	insertData["name"] = datas["name"]
	insertData["ownerid"] = datas["ownerid"]
	insertData["comment"] = datas["comment"]
	insertData["deleted"] = 0
	insertData["creation_time"] = creation_time
	insertData["update_time"] = update_time
	_,err = dbEntity.InsertData("project",insertData)
	errs = append(errs, err...)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error"){
		logErrors(errs)
		ret := buildResponse(1080005,false,"数据查询错误，稍后请重试！")
		c.JSON(http.StatusOK, ret)
		return
	}
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080006,"debug","project has be added."))
	logErrors(errs)
	ret := buildResponse(0,true,"项目添加成功！")
	c.JSON(http.StatusOK, ret)
}

/* 
	getCountHandler get total number rows of  project information according to contions by rquest's URL
	response the client with Status: false, Erorrcode: int, and Message: string if operation is failed
	otherwise response the client with Status: true, Erorrcode: 0, and Message: number if operation is successful
*/
func (p Project) getCountHandler(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080007,"debug","now getting count of  projects through api."))
	conditionKey, _ := c.GetQuery("conditionKey")
	conditionValue, _ := c.GetQuery("conditionValue")
	

	where := make(map[string]string)
	if strings.TrimSpace(conditionKey) != "" && strings.TrimSpace(conditionValue) != ""{
		conditionKey = strings.ToLower(strings.TrimSpace(conditionKey))
		conditionValue = strings.ToLower(strings.TrimSpace(conditionValue))
		if strings.EqualFold(conditionKey,"name") || strings.EqualFold(conditionKey,"comment"){
			where[conditionKey] = " like '%" + conditionValue + "%'"
		}else{
			where[conditionKey] = "=" + conditionValue
		}
	}
	
	deleted, _ := c.GetQuery("deleted")
	deleted = strings.ToLower(strings.TrimSpace(deleted))
	if deleted == "n" || deleted == "0" || deleted == "" {
		where[deleted] = "deleted=0 " 
	}

	field, okField := c.GetQuery("field")
	if !okField {
		field = "*"
	} else {
		field = strings.TrimSpace(field)
	}
	outField := "count(" + field + ") as num"
	// Qeurying data from DB
	selectData := db.SelectData{
		Tb: []string{"project"},
		OutFeilds: []string{outField},
		Where: where,
	}
	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	retData,err := dbEntity.QueryData(&selectData)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error"){
		errs = append(errs,err...)
		logErrors(errs)
		ret := ApiResponseStatus {
			Status: false,
			Errorcode: 1080008,
			Message: "database query error",
		}
		c.JSON(http.StatusOK, ret)
		return 
	} 
	
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080009,"debug","send response data to the client."))
	logErrors(errs)
	ret := ApiResponseStatus {
		Status: true,
		Errorcode: 1080006,
		Message: retData,
	}
	c.JSON(http.StatusOK, ret)
	
}

/* 
	delHandler delete project information from DB according to condition by rquest's URL
	response the client with Status: false, Erorrcode: int, and Message: string if operation is failed
	otherwise response the client with Status: true, Erorrcode: 0, and Message: number if operation is successful
*/
func (p Project) delHandler(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080010,"debug","now deleting projects."))
	keys := []string{"projectid[]"}
	datas,err := utils.GetRequestDataArray(c,keys)
	errs = append(errs,err...)
	data,okdata := datas["projectid[]"]
	if !okdata || len(data) <1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080011,"debug","parameters are error."))
		logErrors(errs)
		ret := buildResponse(1080011,false,"parameters are error")
		c.JSON(http.StatusOK, ret)
		return 
	}
	logErrors(errs)
	errs = errs[0:0]
	var ids = ""
	if len(data) < 2 {
		ids = ids + " =" + data[0]
	}else {
		ids = " in ("
		first := true
		for _,id := range data {
			if first {
				ids += id
				first = false
			} else {
				ids = ids + "," +id
			}
		}
		ids += ")"
	}
	
	where := make(map[string]string)
	where["projectid"] = ids	
	deletetData := db.SelectData{
		Tb: []string{"project"},
		Where: where,
	}
	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	retData,err := dbEntity.DeleteData(&deletetData)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error"){
		errs = append(errs,err...)
		logErrors(errs)
		ret := buildResponse(1080012,false,"database query error")
		c.JSON(http.StatusOK, ret)
		return 
	} 
	
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080013,"debug","send response data to the client."))
	logErrors(errs)
	retNum := strconv.Itoa(int(retData))
	ret := buildResponse(0,true, retNum + "个项目信息已被删除")
	c.JSON(http.StatusOK, ret)
}

/* 
	getInfoHandler get a project information according projectid or project name
*/
func (p Project) getInfoHandler(c *sysadmServer.Context){
	var errs []sysadmerror.Sysadmerror
		   
	projectid,okProjectid := c.GetQuery("id")
	projectname,okProjectname := c.GetQuery("name")

	if (projectid == "" || !okProjectid ) && (projectname == "" || !okProjectname) {
		ret := ApiResponseStatus {
			Status: false,
			Errorcode: 1080014,
			Message: "projectid and project name are both empty or invalid",
		}
		c.JSON(http.StatusOK, ret)
		return 
	}

	// Qeurying data from DB
	whereMap :=  make(map[string]string,0)
	if projectid != "" {
		var ids = ""
		projectids := strings.Split(projectid, ",")
		if len(projectids) >1 {
			ids = " in ("
			first := true
			for _,id := range projectids {
				if first {
					ids += id
					first = false
				} else {
					ids = ids + "," +id
				}
			}
			ids += ")"
		} else {
			ids = ids + " =" + projectid
		}
		whereMap["projectid"] = ids
	}

	if projectname != "" {
		var names = ""
		projectnames := strings.Split(projectname, ",")
		if len(projectnames)>1 {
			names = " in ("
			first := true
			for _,p := range projectnames {
				if first {
					names = names + "'" + p + "'"
					first = false
				} else {
					names = names + ",'" + p + "'"
				}
			}
			names += ")"
		} else {
			names = names + "='"+projectname+"'"
		}
		whereMap["name"] = names
	} 

	selectData := db.SelectData{
		Tb: []string{"project"},
		OutFeilds: []string{"projectid","ownerid","name", "comment","deleted","creation_time","update_time"},
		Where: whereMap,
	}

	dbEntity := RuntimeData.RuningParas.DBConfig.Entity
	retData,err := dbEntity.QueryData(&selectData)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error"){
		errs = append(errs,err...)
		logErrors(errs)
		ret := ApiResponseStatus {
			Status: false,
			Errorcode: 1080015,
			Message: "database query error",
		}
		c.JSON(http.StatusOK, ret)
		return 
	} 

	// if the project is not exist in DB 
	if len(retData) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080016,"debug","no data"))
		logErrors(errs)
		ret := ApiResponseStatus {
			Status: false,
			Errorcode: 1080016,
			Message: "no project",
		}
		c.JSON(http.StatusOK, ret)
		return 
	}
	
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080017,"debug","send response data to the client."))
	logErrors(errs)

	ret := ApiResponseStatus {
		Status: true,
		Errorcode: 1080017,
		Message: retData,
	}

	c.JSON(http.StatusOK, ret)
}
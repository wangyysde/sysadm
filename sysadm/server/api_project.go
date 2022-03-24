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

	"github.com/wangyysde/sysadm/db"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)


var  projectActions = []string{"list"}


 func (p Project)ModuleName()string{
	return "project"
}

func (p Project) ActionHanderCaller(action string, c *sysadmServer.Context){
	switch action{
		case "list":
			p.listHandler(c)
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
		ret := ApiResponseStatus {
			Status: false,
			Errorcode: 1080002,
			Message: "database query error",
		}
		c.JSON(http.StatusOK, ret)
		return 
	} 

	// no project has be queried. 
	if len(retData) < 1 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080002,"debug","no project has be queried."))
		logErrors(errs)
		ret := ApiResponseStatus {
			Status: true,
			Errorcode: 1080002,
			Message: "[]",
		}
		c.JSON(http.StatusOK, ret)
		return 
	}
	
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080003,"debug","send response data to the client."))
	logErrors(errs)
/*
	retJson, e := json.Marshal(retData)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(1080004,"error","marsha result data error: %s",e))
		ret := ApiResponseStatus {
			Status: false,
			Errorcode: 1080004,
			Message: fmt.Sprintf("marsha result data error: %s",e),
		}
		c.JSON(http.StatusOK, ret)
		return 
	}
*/
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
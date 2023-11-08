/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2023 Bzhy Network. All rights reserved.
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

package app

import (
	"fmt"
	"github.com/wangyysde/sysadmServer"
	"net/http"
	"strconv"
	"strings"
	datacenter "sysadm/datacenter/app"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	"sysadm/sysadmLog"
	"sysadm/user"
	"sysadm/utils"
)

// order fields data of cluster list page
var allOrderFields = map[string]string{"TD1": "cnName", "TD2": "enName", "TD3": "datacenterid", "TD5": "status"}

// which field will be order if user has not selected
var defaultOrderField = "TD1"

// 1 for DESC 0 for ASC
var defaultOrderDirection = "1"

// define all list items(cols) name
var allListItems = map[string]string{"TD1": "名称", "TD2": "Name", "TD3": "数据中心", "TD4": "值班电话", "TD5": "状态"}

// all popmenu items defined Format:
// item name, action name, action method
var allPopMenuItems = []string{"查看详情,detail,GET,page", "编辑可用区,edit,GET,page", "删除可用区,del,POST,tip", "启用可用区,enable,POST,tip", "禁用可用区,disable,POST,tip"}

func listHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(700090001, "debug", "now handling available Zone list"))
	messageTemplateFile := "showmessage.html"
	listTemplateFile := "objectlistNew.html"
	messageTplData := make(map[string]interface{}, 0)

	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700090002, "error", "user should login %s", e))
		runData.logEntity.LogErrors(errs)

		tplData := map[string]interface{}{
			"errormessage": "user should login",
		}
		templateFile := "showmessage.html"
		c.HTML(http.StatusOK, templateFile, tplData)
		return
	}

	// get request data
	requestData, e := getRequestData(c)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700090003, "error", "get request data error %s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "系统内部出错，请稍后再试或联系系统管理员"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}

	// 初始化模板数据
	tplData, e := objectsUI.InitTemplateData("/"+DefaultObjectName+"/", "基础设施", "可用区列表", "添加可用区", "yes",
		allPopMenuItems, []string{}, requestData)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700090004, "error", "initate template data error %s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "系统内部出错，请稍后再试或联系系统管理员"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}

	// preparing datacenter data
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = datacenter.New()
	conditions := make(map[string]string, 0)
	conditions["isDeleted"] = "=0"
	order := make(map[string]string, 0)
	var emptyString []string
	dcList, e := dcEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700090005, "error", "get datacenter data error %s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "系统内部出错，请稍后再试或联系系统管理员"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}

	// preparing select group data
	var firstOptions []objectsUI.SelectOption
	firstSelectedID := "0"
	if strings.TrimSpace(requestData["groupSelectID"]) != "" {
		firstSelectedID = strings.TrimSpace(requestData["groupSelectID"])
	}
	firstSelect := objectsUI.SelectData{
		Title:      "数据中心",
		SelectedId: firstSelectedID,
	}
	for _, line := range dcList {
		dcData, ok := line.(datacenter.DatacenterSchema)
		if !ok {
			outPutErrorMsg(c, 700090006, errs, fmt.Errorf("data is not datacenter"))
			return
		}
		option := objectsUI.SelectOption{
			Id:       strconv.Itoa(int(dcData.Id)),
			Text:     dcData.CnName,
			ParentID: "0",
		}
		firstOptions = append(firstOptions, option)
	}
	firstSelect.Options = firstOptions
	secondSelect := objectsUI.SelectData{}
	thirdSelect := objectsUI.SelectData{}
	if e := objectsUI.BuildMultiSelectData(firstSelect, secondSelect, thirdSelect, tplData); e != nil {
		outPutErrorMsg(c, 700090007, errs, e)
		return
	}

	// build table header for list objects
	objectsUI.BuildThData(requestData, allOrderFields, allListItems, tplData, defaultOrderField, defaultOrderDirection)

	searchContent := objectsUI.GetSearchContentFromRequest(requestData)
	ids := objectsUI.GetObjectIdsFromRequest(requestData)
	searchKeys := []string{"id", "cnName", "enName"}
	startPos := objectsUI.GetStartPosFromRequest(requestData)
	azConditions := objectsUI.BuildCondition(requestData, "=0", "datacenterid")

	// get total number of list objects
	var azEntity sysadmObjects.ObjectEntity
	azEntity = New()
	azCount, e := azEntity.GetObjectCount(searchContent, ids, searchKeys, azConditions)
	if e != nil || azCount < 1 {
		outPutErrorMsg(c, 700090008, errs, e)
		return
	}

	// get list data
	requestOrder := objectsUI.BuildOrderDataForQuery(requestData, allOrderFields, defaultOrderField, defaultOrderDirection)
	azList, e := azEntity.GetObjectList(searchContent, ids, searchKeys, azConditions, startPos, runData.pageInfo.NumPerPage, requestOrder)
	if e != nil {
		outPutErrorMsg(c, 700090009, errs, e)
		return
	}

	// prepare cluster list data
	objListData, e := prepareObjectData(dcList, azList)
	if e != nil {
		outPutErrorMsg(c, 700090010, errs, e)
		return
	}
	tplData["objListData"] = objListData

	// prepare page number information
	objectsUI.BuildPageNumInfo(tplData, requestData, azCount, startPos, runData.pageInfo.NumPerPage, defaultOrderField, defaultOrderDirection)

	runData.logEntity.LogErrors(errs)
	c.HTML(http.StatusOK, listTemplateFile, tplData)
}

func getRequestData(c *sysadmServer.Context) (map[string]string, error) {
	requestData, e := utils.NewGetRequestData(c, []string{"groupSelectID", "searchContent", "start", "orderfield", "direction"})
	if e != nil {
		return requestData, e
	}

	objectIds := ""
	objectIDMap, _ := utils.GetRequestDataArray(c, []string{"objectid[]"})
	if objectIDMap != nil {
		objectIDSlice, ok := objectIDMap["objectid[]"]
		if ok {
			objectIds = strings.Join(objectIDSlice, ",")
		}
	}
	requestData["objectIds"] = objectIds
	if strings.TrimSpace(requestData["start"]) == "" {
		requestData["start"] = "0"
	}

	return requestData, nil
}

func outPutErrorMsg(c *sysadmServer.Context, errcode int, errs []sysadmLog.Sysadmerror, e error) {
	messageTemplateFile := "showmessage.html"
	messageTplData := make(map[string]interface{}, 0)
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(errcode, "error", "%s", e))
	runData.logEntity.LogErrors(errs)
	messageTplData["errormessage"] = "系统内部出错，请稍后再试或联系系统管理员"
	c.HTML(http.StatusOK, messageTemplateFile, messageTplData)

	return
}

func prepareObjectData(dcList, azList []interface{}) ([]map[string]interface{}, error) {
	var dataList []map[string]interface{}
	dcData := make(map[uint]datacenter.DatacenterSchema)
	for _, line := range dcList {
		lineData, ok := line.(datacenter.DatacenterSchema)
		if !ok {
			return dataList, fmt.Errorf("handle datacenter data error")
		}
		dcData[lineData.Id] = lineData
	}

	for _, line := range azList {
		lineMap := make(map[string]interface{}, 0)
		lineData, ok := line.(AvailablezoneSchema)
		if !ok {
			return dataList, fmt.Errorf("AZ data is not valid")
		}
		lineMap["objectID"] = lineData.Id
		lineMap["TD1"] = lineData.CnName
		lineMap["TD2"] = lineData.EnName
		lineMap["TD3"] = getDCName(lineData.Datacenterid, dcData)
		lineMap["TD4"] = lineData.DutyTel
		statusStr := "未知"
		popmenuitems := ""
		switch lineData.Status {
		case StatusUnused:
			statusStr = "未启用"
			popmenuitems = "0,1,2,3"
		case StatusEnabled:
			statusStr = "启用"
			popmenuitems = "0,1,4"
		case StatusDisabled:
			statusStr = "已禁用"
			popmenuitems = "0,1,2,3"
		}
		lineMap["TD5"] = statusStr
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return dataList, nil
}

func getDCName(dcID uint, dcData map[uint]datacenter.DatacenterSchema) string {
	dcLine, ok := dcData[dcID]
	if !ok {
		return "未知"
	}

	return dcLine.CnName
}

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
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	"sysadm/sysadmLog"
	"sysadm/user"
	"sysadm/utils"
	sysadmVersion "sysadm/version/app"
)

// order fields data of cluster list page
var allOrderFields = map[string]string{"TD1": "name"}

// which field will be order if user has not selected
var defaultOrderField = "TD1"

// 1 for DESC 0 for ASC
var defaultOrderDirection = "1"

// all popmenu items defined Format:
// item name, action name, action method
var allPopMenuItems = []string{"删除,del,POST,tip"}

// define all list items(cols) name
var allListItems = map[string]string{"TD1": "名称", "TD2": "体系结构", "TD3": "位数", "TD4": "版本", "TD5": "描述"}

func listHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(7000180001, "debug", "now handling OS list"))
	listTemplateFile := "objectlistNew.html"

	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000180002, errs, e)
		return
	}

	// get request data
	requestData, e := getRequestData(c)
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000180003, errs, e)
		return
	}

	// 初始化模板数据
	tplData, e := objectsUI.InitTemplateData("/"+defaultObjectName+"/", "基础设施", "操作系统列表", "添加操作系统信息", "yes",
		allPopMenuItems, []string{}, requestData)
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000180004, errs, e)
		return
	}

	firstGroupSelect := objectsUI.SelectData{}
	secondGroupSelect := objectsUI.SelectData{}
	thirdGroupSelect := objectsUI.SelectData{}
	if e := objectsUI.BuildMultiSelectData(firstGroupSelect, secondGroupSelect, thirdGroupSelect, tplData); e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000180009, errs, e)
		return
	}

	// build table header for list objects
	objectsUI.BuildThData(requestData, allOrderFields, allListItems, tplData, defaultOrderField, defaultOrderDirection)

	searchContent := objectsUI.GetSearchContentFromRequest(requestData)
	ids := objectsUI.GetObjectIdsFromRequest(requestData)
	searchKeys := []string{"name"}
	startPos := objectsUI.GetStartPosFromRequest(requestData)

	// preparing os data
	var osEntity sysadmObjects.ObjectEntity
	osEntity = New()
	conditions := make(map[string]string, 0)
	order := make(map[string]string, 0)
	var emptyString []string

	osCount, e := osEntity.GetObjectCount(searchContent, ids, searchKeys, conditions)
	if e != nil || osCount < 1 {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000180005, errs, e)
		return
	}

	// get list data
	requestOrder := objectsUI.BuildOrderDataForQuery(requestData, allOrderFields, defaultOrderField, defaultOrderDirection)
	osList, e := osEntity.GetObjectList(searchContent, ids, searchKeys, conditions, startPos, runData.pageInfo.NumPerPage, requestOrder)
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000180006, errs, e)
		return
	}

	var versionEntity sysadmObjects.ObjectEntity
	versionEntity = sysadmVersion.New()
	conditions["typeID"] = "='" + strconv.Itoa(int(sysadmVersion.VersionTypeOS)) + "'"
	versionList, e := versionEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000180007, errs, e)
		return
	}

	// prepare cluster list data
	objListData, e := prepareObjectData(osList, versionList)
	if e != nil {
		objectsUI.OutPutErrorMsg(c, "", runData.logEntity, 7000180008, errs, e)
		return
	}

	tplData["objListData"] = objListData

	// prepare page number information
	objectsUI.BuildPageNumInfo(tplData, requestData, osCount, startPos, runData.pageInfo.NumPerPage, defaultOrderField, defaultOrderDirection)

	runData.logEntity.LogErrors(errs)
	c.HTML(http.StatusOK, listTemplateFile, tplData)

}

func getRequestData(c *sysadmServer.Context) (map[string]string, error) {
	requestData, e := utils.NewGetRequestData(c, []string{"searchContent", "start", "orderfield", "direction"})
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

func prepareObjectData(osList, versionList []interface{}) ([]map[string]interface{}, error) {
	var dataList []map[string]interface{}

	for _, line := range osList {
		osData, ok := line.(OSSchema)
		if !ok {
			return dataList, fmt.Errorf("data is not OS defined schema")
		}

		lineMap := make(map[string]interface{}, 0)
		lineMap["TD1"] = osData.Name
		lineMap["TD2"] = osData.Architecture
		lineMap["TD3"] = strconv.Itoa(osData.Bit)
		lineMap["TD4"] = getVerName(osData.OSID, versionList)
		lineMap["TD5"] = osData.Description
		popmenuitems := "0"
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}
	return dataList, nil
}

func getVerName(osID int, versionList []interface{}) string {
	for _, line := range versionList {
		verData, ok := line.(sysadmVersion.VersionSchema)
		if !ok {
			return "未知"
		}

		if verData.OSID == osID {
			return verData.Name
		}
	}

	return "未知"
}

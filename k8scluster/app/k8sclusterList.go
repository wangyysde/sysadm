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
	az "sysadm/availablezone/app"
	datacenter "sysadm/datacenter/app"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	"sysadm/sysadmLog"
	"sysadm/user"
	"sysadm/utils"
)

// order fields data of cluster list page
var allOrderFields = map[string]string{"TD1": "id", "TD3": "cnName", "TD5": "status"}

// which field will be order if user has not selected
var defaultOrderField = "TD1"

// 1 for DESC 0 for ASC
var defaultOrderDirection = "1"

// all popmenu items defined Format:
// item name, action name, action method
var allPopMenuItems = []string{"查看详情,detail,GET", "节点列表,list,GET", "删除集群,del,POST", "查看日志,getlog,GET", "kubeconf下载,getkubeconf,GET"}

// define all list items(cols) name
var allListItems = map[string]string{"TD1": "集群ID", "TD2": "数据中心/可用区", "TD3": "集群名", "TD4": "版本", "TD5": "状态"}

func listHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(700050001, "debug", "now handling k8s cluster list"))
	messageTemplateFile := "showmessage.html"
	listTemplateFile := "objectllist.html"
	messageTplData := make(map[string]interface{}, 0)
	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700050002, "error", "user should login %s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "user should login"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}

	// get request data
	requestData, e := getRequestData(c)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700050003, "error", "get request data error %s", e))
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
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700050004, "error", "get datacenter data error %s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "系统内部出错，请稍后再试或联系系统管理员"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}

	// preparing AZ data
	var azEntity sysadmObjects.ObjectEntity
	azEntity = az.New()
	azList, e := azEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700050005, "error", "get available zone data error %s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "系统内部出错，请稍后再试或联系系统管理员"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}

	searchContent := objectsUI.GetSearchContentFromRequest(requestData)
	ids := objectsUI.GetObjectIdsFromRequest(requestData)
	searchKeys := []string{"id", "cnName", "version"}
	startPos := objectsUI.GetStartPosFromRequest(requestData)
	clusterConditions := objectsUI.BuildCondition(requestData, "=0", "dcid")

	// get total number of list objects
	var clusterEntity sysadmObjects.ObjectEntity
	clusterEntity = New()
	clusterCount, e := clusterEntity.GetObjectCount(searchContent, ids, searchKeys, clusterConditions)
	if e != nil || clusterCount < 1 {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700050006, "error", "get cluster data error %s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "系统中没有查到集群信息数据"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}

	// get list data
	requestOrder := objectsUI.BuildOrderDataForQuery(requestData, allOrderFields, defaultOrderField, defaultOrderDirection)
	clusterList, e := clusterEntity.GetObjectList(searchContent, ids, searchKeys, clusterConditions, startPos, runData.pageInfo.NumPerPage, requestOrder)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700050007, "error", "get cluster data error %s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "系统内部出错，请稍后再试或联系系统管理员"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}

	// set template data for select and search block
	tplData, e := objectsUI.InitTemplateData("/"+DefaultObjectName+"/", "集群管理", "集群信息列表", "添加集群", "yes",
		allPopMenuItems, []string{}, requestData)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700050008, "error", "initate template data error %s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "系统内部出错，请稍后再试或联系系统管理员"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}

	// set template data for all select items
	tplData["groupSelect"] = buildDatacenterSelectData(dcList)
	objectsUI.BuildThData(requestData, allOrderFields, allListItems, tplData, defaultOrderField, defaultOrderDirection)

	// prepare cluster list data
	objListData, e := prepareClusterData(clusterList, dcList, azList)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700050009, "error", "prepare cluster list data error %s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "系统内部出错，请稍后再试或联系系统管理员"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}
	tplData["objListData"] = objListData

	// prepare page number information
	objectsUI.BuildPageNumInfo(tplData, requestData, clusterCount, startPos, runData.pageInfo.NumPerPage, defaultOrderField, defaultOrderDirection)

	runData.logEntity.LogErrors(errs)
	c.HTML(http.StatusOK, listTemplateFile, tplData)
}

func buildDatacenterSelectData(dzList []interface{}) []map[string]interface{} {
	var ret []map[string]interface{}
	// default item
	row := make(map[string]interface{}, 0)
	row["id"] = "0"
	row["text"] = "所有数据中心"
	ret = append(ret, row)
	for _, v := range dzList {
		row := make(map[string]interface{}, 0)
		dzData := v.(datacenter.DatacenterSchema)
		row["id"] = strconv.Itoa(int(dzData.Id))
		row["text"] = dzData.CnName
		ret = append(ret, row)
	}

	return ret
}

func prepareClusterData(clustData, dcData, azData []interface{}) ([]map[string]interface{}, error) {
	var dataList []map[string]interface{}
	var dcSchemaData []datacenter.DatacenterSchema
	var azSchemaData []az.AvailablezoneSchema
	for _, line := range dcData {
		lineData, ok := line.(datacenter.DatacenterSchema)
		if !ok {
			return dataList, fmt.Errorf("handle datacenter data error")
		}
		dcSchemaData = append(dcSchemaData, lineData)
	}

	for _, line := range azData {
		lineData, ok := line.(az.AvailablezoneSchema)
		if !ok {
			return dataList, fmt.Errorf("handle availablezone data error")
		}
		azSchemaData = append(azSchemaData, lineData)
	}

	for _, line := range clustData {
		lineMap := make(map[string]interface{}, 0)
		lineData, ok := line.(K8sclusterSchema)
		if !ok {
			return dataList, fmt.Errorf("handle cluster  data error")
		}
		lineMap["objectID"] = lineData.Id
		lineMap["TD1"] = lineData.Id
		lineMap["TD2"] = getDcAzName(lineData.Azid, dcSchemaData, azSchemaData)
		lineMap["TD3"] = lineData.CnName
		lineMap["TD4"] = lineData.Version
		lineMap["TD5"] = GetStatusText(lineData.Status)
		lineMap["popmenuitems"] = getPopMenuItemsId(lineData.Status, lineData.IsDeleted)
		dataList = append(dataList, lineMap)
	}

	return dataList, nil
}

func getDcAzName(azid uint, dcData []datacenter.DatacenterSchema, azData []az.AvailablezoneSchema) string {
	dcName := "未知"
	azName := "未知"
	dcID := uint(0)
	for _, azLine := range azData {
		if azLine.Id == azid {
			azName = azLine.CnName
			dcID = azLine.Datacenterid
			break
		}
	}

	for _, dcLine := range dcData {
		if dcLine.Id == dcID {
			dcName = dcLine.CnName
			break
		}
	}

	return dcName + "/" + azName
}

func getPopMenuItemsId(status, isDeleted int) string {
	popmenuitemidstr := ""
	switch status {
	case 0:
		popmenuitemidstr = "0"
	case 1:
		popmenuitemidstr = "0,1,2,3,4"
	case 2:
		popmenuitemidstr = "0,1,2,3,4"
	default:
		popmenuitemidstr = "0"
	}

	if isDeleted != 1 {
		popmenuitemidstr = popmenuitemidstr + ",2"
	}

	return popmenuitemidstr
}

func getRequestData(c *sysadmServer.Context) (map[string]string, error) {
	requestData, e := utils.NewGetRequestData(c, []string{"groupSelectID", "searchContent", "objectid", "start", "orderfield", "direction"})
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

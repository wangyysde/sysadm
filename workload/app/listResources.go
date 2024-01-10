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
	"github.com/wangyysde/sysadmServer"
	"net/http"
	"strings"
	datacenter "sysadm/datacenter/app"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	"sysadm/sysadmLog"
	"sysadm/user"
)

func listResourceHandler(c *sysadmServer.Context, module, action string) {
	var errs []sysadmLog.Sysadmerror
	var additionalJs = []string{"js/sysadmfunctions.js", "/js/workloadList.js"}
	var additionalCss = []string{}

	var listTemplateFile = "workloadlist.html"

	errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000800002, "debug", "now handling %s list", module))

	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		objectsUI.OutPutMsg(c, "", "您未登录或超时", runData.logEntity, 8000800003, errs, e)
		return
	}

	// get request data
	requestKeys := []string{"dcID", "clusterID", "namespace", "start", "orderfield", "direction", "searchContent", "objID"}
	requestData, e := getRequestData(c, requestKeys)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "参数错误,请确认是从正确地方连接过来的", runData.logEntity, 8000800004, errs, e)
		return
	}

	// 为前端下拉菜单的数据中心部分准备数据
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = datacenter.New()
	conditions := make(map[string]string, 0)
	conditions["isDeleted"] = "=0"
	order := make(map[string]string, 0)
	var emptyString []string
	dcList, e := dcEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000600005, errs, e)
		return
	}

	// preparing object list data
	selectedCluster := strings.TrimSpace(requestData["clusterID"])
	if selectedCluster == "" {
		selectedCluster = "0"
	}
	selectedNamespace := strings.TrimSpace(requestData["namespace"])
	if action == "QuotaList" {
		selectedNamespace = requestData["objID"]
		requestData["namespace"] = requestData["objID"]
	}
	if selectedNamespace == "" {
		selectedNamespace = "0"
	}

	objEntity, e := newObjEntity(module)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000600006, errs, e)
		return
	}
	objEntity.setObjectInfo()

	// 初始化模板数据
	subCategory := objEntity.getModuleName() + "列表"
	addButtonTitle := objEntity.getAddButtonTitle()
	if selectedCluster == "0" {
		addButtonTitle = ""
	}

	isSearchForm := objEntity.getIsSearchForm()
	allPopMenuItems := objEntity.getAllPopMenuItems()
	if action == "QuotaList" {
		subCategory = "命名空间 >> 资源配额列表"
		addButtonTitle = ""
		isSearchForm = ""
		allPopMenuItems = quotaListPagePopmenu
	}

	tplData, e := objectsUI.InitTemplateDataForWorkload("/"+defaultObjectName+"/", objEntity.getMainModuleName(), subCategory, addButtonTitle, isSearchForm,
		allPopMenuItems, additionalJs, additionalCss, requestData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000700006, errs, e)
		return
	}
	tplData["objName"] = module

	if objEntity.getNamespaced() {
		e = buildSelectDataWithNs(tplData, dcList, requestData)
	} else {
		if action == "QuotaList" {
			e = buildSelectDataWithNs(tplData, dcList, requestData)
		} else {
			e = buildSelectData(tplData, dcList, requestData)
		}

	}
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000600007, errs, e)
		return
	}
	startPos := objectsUI.GetStartPosFromRequest(requestData)

	var count int = 0
	var objListData []map[string]interface{}
	if action == "QuotaList" {
		count, objListData, e = listQuotaData(selectedCluster, selectedNamespace, startPos, requestData)
	} else {
		count, objListData, e = objEntity.listObjectData(selectedCluster, selectedNamespace, startPos, requestData)
	}
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000600008, errs, e)
		return
	}

	if count == 0 {
		tplData["noData"] = "当前系统中没有数据"
	} else {
		tplData["noData"] = ""
		tplData["objListData"] = objListData

		// build table header for list objects
		allListItems := objEntity.getAllListItems()
		defaultOrderField := objEntity.getDefaultOrderField()
		objDefaultOrderDirection := objEntity.getDefaultOrderDirection()
		allorderFields := objEntity.getAllorderFields()
		if action == "QuotaList" {
			allListItems = quotaListAllListItems
			defaultOrderField = quotaListDefaultOrderField
			objDefaultOrderDirection = quotaListDefaultOrderDirection
			allorderFields = quotaListAllOrderFields
		}
		objectsUI.BuildThDataWithOrderFunc(requestData, allListItems, tplData, defaultOrderField, objDefaultOrderDirection, allorderFields)

		// prepare page number information
		objectsUI.BuildPageNumInfoForWorkloadList(tplData, requestData, count, startPos, runData.pageInfo.NumPerPage, defaultOrderField, objDefaultOrderDirection)
	}
	runData.logEntity.LogErrors(errs)
	if objEntity.getTemplateFile(action) != "" {
		listTemplateFile = objEntity.getTemplateFile(action)
	}
	c.HTML(http.StatusOK, listTemplateFile, tplData)

}

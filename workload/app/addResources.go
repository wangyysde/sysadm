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
	datacenter "sysadm/datacenter/app"
	sysadmCluster "sysadm/k8scluster/app"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	"sysadm/sysadmLog"
	"sysadm/sysadmapi/apiutils"
	"sysadm/user"
)

func addFormResourceHandler(c *sysadmServer.Context, module, action string) {
	var errs []sysadmLog.Sysadmerror
	var defaultemplateFile = "addObjTabs.html"
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000900001, "debug", "now display %s add form", module))

	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		objectsUI.OutPutMsg(c, "", "您未登录或超时", runData.logEntity, 8000900002, errs, e)
		return
	}

	// get request data
	requestKeys := []string{"dcID", "clusterID", "namespace"}
	requestData, e := getRequestData(c, requestKeys)
	if requestData["clusterID"] == "" || requestData["clusterID"] == "0" {
		objectsUI.OutPutMsg(c, "", "请选择需要添加资源的集群", runData.logEntity, 8000900003, errs, e)
		return
	}
	objEntity, e := newObjEntity(module)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900004, errs, e)
		return
	}
	objEntity.setObjectInfo()

	// preparing datacenter data
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = datacenter.New()
	dcInfo, e := dcEntity.GetObjectInfoByID(requestData["dcID"])
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900006, errs, e)
		return
	}
	dcInfoData, ok := dcInfo.(datacenter.DatacenterSchema)
	if !ok {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900007, errs, fmt.Errorf("data is not datacenter schema"))
		return
	}
	dcName := dcInfoData.CnName

	// get cluster name
	var clusterEntity sysadmObjects.ObjectEntity
	clusterEntity = sysadmCluster.New()
	clusterInfo, e := clusterEntity.GetObjectInfoByID(requestData["clusterID"])
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900008, errs, e)
		return
	}
	cluserInfoData, ok := clusterInfo.(sysadmCluster.K8sclusterSchema)
	if !ok {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900009, errs, fmt.Errorf("data is not k8scluster schema"))
		return
	}
	clusterName := cluserInfoData.CnName

	if requestData["namespace"] == "0" {
		requestData["namespace"] = ""
	}

	// 初始化模板数据
	tplData, e := objectsUI.InitTemplateDataForWorkload("/"+defaultObjectName+"/", objEntity.getMainModuleName(), objEntity.getModuleName(), "", "no",
		[]string{}, objEntity.getAdditionalJs(), objEntity.getAdditionalCss(), make(map[string]string, 0))
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900010, errs, e)
		return
	}

	tplData["enctype"] = true
	tplData["dcID"] = requestData["dcID"]
	tplData["clusterID"] = requestData["clusterID"]
	tplData["namespace"] = requestData["namespace"]
	tplData["objID"] = objEntity.getModuleID()
	tplData["dcName"] = dcName
	tplData["clusterName"] = clusterName
	tplData["objName"] = objEntity.getModuleName()
	tplData["apiVersion"] = apiVersion
	e = objEntity.buildAddFormData(tplData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "系统内容错误，请稍后再试，如果仍有问题，请联系系统管理员", runData.logEntity, 8000900011, errs, e)
		return
	}
	runData.logEntity.LogErrors(errs)

	pageTemplateFile := defaultemplateFile
	if objEntity.getTemplateFile(action) != "" {
		pageTemplateFile = objEntity.getTemplateFile(action)
	}
	c.HTML(http.StatusOK, pageTemplateFile, tplData)
}

func addNewQuotaFormHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	var listTemplateFile = "addObjTabs.html"
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000900012, "debug", "now adding quota to namespace"))

	requestKeys := []string{"dcID", "clusterID", "namespace", "objID"}
	requestData, e := checkUserLoginAndGetRequestData(c, requestKeys)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "您未登录或系统出错，请稍后再试，如仍有问题请联系系统管理员", runData.logEntity, 8000900013, errs, e)
		return
	}

	if requestData["objID"] == "" {
		objectsUI.OutPutMsg(c, "", "添加配额需要指定命名空间", runData.logEntity, 8000900014, errs, e)
		return
	}

	nsEntity := &namespace{}
	nsEntity.setObjectInfo()
	subModuleName := nsEntity.getModuleName() + ">>新增资源配额"
	tplData, e := objectsUI.InitTemplateDataForWorkload("/"+defaultObjectName+"/", nsEntity.getMainModuleName(), subModuleName, "", "no",
		[]string{}, nsEntity.additionalJs, nsEntity.additionalCss, make(map[string]string, 0))
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900015, errs, e)
		return
	}

	// preparing datacenter data
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = datacenter.New()
	dcInfo, e := dcEntity.GetObjectInfoByID(requestData["dcID"])
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900016, errs, e)
		return
	}
	dcInfoData, ok := dcInfo.(datacenter.DatacenterSchema)
	if !ok {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900017, errs, fmt.Errorf("data is not datacenter schema"))
		return
	}
	dcName := dcInfoData.CnName

	// get cluster name
	var clusterEntity sysadmObjects.ObjectEntity
	clusterEntity = sysadmCluster.New()
	clusterInfo, e := clusterEntity.GetObjectInfoByID(requestData["clusterID"])
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900018, errs, e)
		return
	}
	cluserInfoData, ok := clusterInfo.(sysadmCluster.K8sclusterSchema)
	if !ok {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900019, errs, fmt.Errorf("data is not k8scluster schema"))
		return
	}
	clusterName := cluserInfoData.CnName

	tplData["enctype"] = true
	tplData["dcID"] = requestData["dcID"]
	tplData["clusterID"] = requestData["clusterID"]
	tplData["namespace"] = requestData["objID"]
	tplData["objID"] = nsEntity.getModuleID()
	tplData["dcName"] = dcName
	tplData["clusterName"] = clusterName
	tplData["objName"] = "资源配额"

	e = nsEntity.buildAddQuotaFormData(tplData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "系统内容错误，请稍后再试，如果仍有问题，请联系系统管理员", runData.logEntity, 8000900020, errs, e)
		return
	}
	runData.logEntity.LogErrors(errs)
	c.HTML(http.StatusOK, listTemplateFile, tplData)
}

func addNewQuotaHandler(c *sysadmServer.Context, action string) {
	var errs []sysadmLog.Sysadmerror
	var response apiutils.ApiResponseData

	errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000900020, "debug", "try to add new object to namespace"))
	islogin, _, _ := user.IsLogin(c, runData.sessionName)
	if !islogin {
		response = apiutils.BuildResponseDataForError(8000900021, "您没有登录或者没有权限执行本操作")
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000900021, "info", "user has not login or not permission when adding new quota"))
		runData.logEntity.LogErrors(errs)
		c.JSON(http.StatusOK, response)
		return
	}

	nsEntity := namespace{}
	nsEntity.setObjectInfo()
	var e error
	switch action {
	case "addNewQuota":
		e = nsEntity.addNewQuota(c)
	case "addNewLimitRange":
		e = nsEntity.addNewLimitRange(c)
	}
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000900022, "info", "%s", e))
		runData.logEntity.LogErrors(errs)
		response = apiutils.BuildResponseDataForError(8000900022, "系统出错，请稍后再试或者联系系统管理员")
		c.JSON(http.StatusOK, response)
		return
	}

	response = apiutils.BuildResponseDataForSuccess("新的配额信息已经添加成功")
	c.JSON(http.StatusOK, response)
	return
}

func addNewLimitRangehandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	var listTemplateFile = "addObjTabs.html"
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(8000900023, "debug", "now adding LimitRange to namespace"))

	requestKeys := []string{"dcID", "clusterID", "namespace", "objID"}
	requestData, e := checkUserLoginAndGetRequestData(c, requestKeys)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "您未登录或系统出错，请稍后再试，如仍有问题请联系系统管理员", runData.logEntity, 8000900024, errs, e)
		return
	}

	if requestData["objID"] == "" {
		objectsUI.OutPutMsg(c, "", "添加默认资源配额需要指定命名空间", runData.logEntity, 8000900025, errs, e)
		return
	}

	nsEntity := &namespace{}
	nsEntity.setObjectInfo()
	subModuleName := nsEntity.getModuleName() + "  >>  新增默认资源配额"
	tplData, e := objectsUI.InitTemplateDataForWorkload("/"+defaultObjectName+"/", nsEntity.getMainModuleName(), subModuleName, "", "no",
		[]string{}, nsEntity.additionalJs, nsEntity.additionalCss, make(map[string]string, 0))
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900026, errs, e)
		return
	}

	// preparing datacenter data
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = datacenter.New()
	dcInfo, e := dcEntity.GetObjectInfoByID(requestData["dcID"])
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900027, errs, e)
		return
	}
	dcInfoData, ok := dcInfo.(datacenter.DatacenterSchema)
	if !ok {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900028, errs, fmt.Errorf("data is not datacenter schema"))
		return
	}
	dcName := dcInfoData.CnName

	// get cluster name
	var clusterEntity sysadmObjects.ObjectEntity
	clusterEntity = sysadmCluster.New()
	clusterInfo, e := clusterEntity.GetObjectInfoByID(requestData["clusterID"])
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900029, errs, e)
		return
	}
	cluserInfoData, ok := clusterInfo.(sysadmCluster.K8sclusterSchema)
	if !ok {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 8000900030, errs, fmt.Errorf("data is not k8scluster schema"))
		return
	}
	clusterName := cluserInfoData.CnName

	tplData["enctype"] = true
	tplData["dcID"] = requestData["dcID"]
	tplData["clusterID"] = requestData["clusterID"]
	tplData["namespace"] = requestData["objID"]
	tplData["objID"] = nsEntity.getModuleID()
	tplData["dcName"] = dcName
	tplData["clusterName"] = clusterName
	tplData["objName"] = "默认资源配额"

	e = nsEntity.buildAddLimitRangeFormData(tplData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", fmt.Sprintf("%s", e), runData.logEntity, 8000900031, errs, e)
		return
	}
	runData.logEntity.LogErrors(errs)
	c.HTML(http.StatusOK, listTemplateFile, tplData)
}

/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2024 Bzhy Network. All rights reserved.
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
	"context"
	"fmt"
	"github.com/wangyysde/sysadmServer"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"net/http"
	"sort"
	"strconv"
	"strings"
	sysadmAZ "sysadm/availablezone/app"
	datacenter "sysadm/datacenter/app"
	sysadmK8sClient "sysadm/k8sclient"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	sysadmOS "sysadm/os/app"
	"sysadm/sysadmLog"
	"sysadm/sysadmapi/apiutils"
	"sysadm/user"
	"sysadm/utils"
	sysadmVersion "sysadm/version/app"
)

func listHostHandler(c *sysadmServer.Context) {
	var listTemplateFile = "hostlist.html"
	var errs []sysadmLog.Sysadmerror

	errs = append(errs, sysadmLog.NewErrorWithStringLevel(7001400007, "debug", "now list host data"))
	requestKeys := []string{"dcID", "azID", "clusterID", "userid", "start", "orderfield", "direction", "searchContent", "objID"}
	requestData, e := objectsUI.GetRequestData(c, requestKeys)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "系统出错，请稍后再试或联系系统管理员", runData.logEntity, 7001400008, errs, e)
		return
	}

	tplData, e := objectsUI.InitTemplateDataForWorkload("/"+hostObjectName+"/", "集群管理",
		"节点管理", "添加节点", "",
		hostAllPopMenuItems, []string{"/js/hostlist.js", "/js/sysadmfunctions.js"}, []string{}, requestData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "系统出错，请稍后再试或联系系统管理员", runData.logEntity, 7001400009, errs, e)
		return
	}
	tplData["objName"] = "节点列表"

	// 为前端下拉菜单的数据中心部分准备数据
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = datacenter.New()
	conditions := make(map[string]string, 0)
	conditions["isDeleted"] = "=0"
	order := make(map[string]string, 0)
	var emptyString []string
	dcList, e := dcEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, order)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "系统出错，请稍后再试或联系系统管理员", runData.logEntity, 7001400009, errs, e)
		return
	}
	e = buildSelectData(tplData, dcList, requestData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "系统出错，请稍后再试或联系系统管理员", runData.logEntity, 7001400010, errs, e)
		return
	}

	startPos := objectsUI.GetStartPosFromRequest(requestData)

	dcSelect := tplData["dcSelect"].(objectsUI.SelectData)
	azSelect := tplData["azSelect"].(objectsUI.SelectData)
	clusterSelect := tplData["clusterSelect"].(objectsUI.SelectData)
	count, objListData, e := listHostData(dcSelect.SelectedId, azSelect.SelectedId, clusterSelect.SelectedId, startPos, requestData)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "", runData.logEntity, 7001400011, errs, e)
		return
	}
	if count == 0 {
		tplData["noData"] = "当前系统中没有数据"
	} else {
		tplData["noData"] = ""
		tplData["objListData"] = objListData
		objectsUI.BuildThDataWithOrderFunc(requestData, hostAllListItems, tplData, "TD1", "1", hostAllOrderFields)
		// prepare page number information
		objectsUI.BuildPageNumInfoForWorkloadList(tplData, requestData, count, startPos, runData.pageInfo.NumPerPage, "TD1", "1")
	}

	c.HTML(http.StatusOK, listTemplateFile, tplData)
}

func buildSelectData(tplData map[string]interface{}, dcList []interface{}, requestData map[string]string) error {

	selectedDC := strings.TrimSpace(requestData["dcID"])
	if selectedDC == "" {
		selectedDC = "0"
	}
	selectedAZ := strings.TrimSpace(requestData["azID"])
	if selectedAZ == "" {
		selectedAZ = "0"
	}
	selectedCluster := strings.TrimSpace(requestData["clusterID"])
	if selectedCluster == "" {
		selectedCluster = "0"
	}

	// preparing datacenter options
	var dcOptions []objectsUI.SelectOption
	dcOption := objectsUI.SelectOption{
		Id:       "0",
		Text:     "===选择数据中心===",
		ParentID: "0",
	}
	dcOptions = append(dcOptions, dcOption)
	for _, line := range dcList {
		dcData, ok := line.(datacenter.DatacenterSchema)
		if !ok {
			return fmt.Errorf("the data is not datacenter schema")
		}
		dcOption := objectsUI.SelectOption{
			Id:       strconv.Itoa(int(dcData.Id)),
			Text:     dcData.CnName,
			ParentID: "0",
		}
		dcOptions = append(dcOptions, dcOption)
	}
	dcSelect := objectsUI.SelectData{Title: "数据中心", SelectedId: selectedDC, Options: dcOptions}
	tplData["dcSelect"] = dcSelect

	// preparing AZ options
	var azOptions []objectsUI.SelectOption
	if selectedDC != "0" {
		azOption := objectsUI.SelectOption{
			Id:       "0",
			Text:     "===所有可用区===",
			ParentID: "0",
		}
		azOptions = append(azOptions, azOption)
		var azEntity sysadmObjects.ObjectEntity
		azEntity = sysadmAZ.New()
		conditions := make(map[string]string, 0)
		var emptyString []string
		conditions["isDeleted"] = "='0'"
		conditions["datacenterid"] = "='" + selectedDC + "'"
		azList, e := azEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, make(map[string]string))
		if e != nil {
			return e
		}
		for _, line := range azList {
			azData, ok := line.(sysadmAZ.AvailablezoneSchema)
			if !ok {
				return fmt.Errorf("the data is not available zone schema")
			}
			azOption := objectsUI.SelectOption{
				Id:       strconv.Itoa(int(azData.Id)),
				Text:     azData.CnName,
				ParentID: strconv.Itoa(int(azData.Datacenterid)),
			}
			azOptions = append(azOptions, azOption)
		}
	} else {
		azOption := objectsUI.SelectOption{
			Id:       "0",
			Text:     "===选择可用区===",
			ParentID: "0",
		}
		azOptions = append(azOptions, azOption)
	}
	azSelect := objectsUI.SelectData{Title: "可用区", SelectedId: selectedAZ, Options: azOptions}
	tplData["azSelect"] = azSelect

	// preparing cluster options
	var clusterOptions []objectsUI.SelectOption
	if selectedAZ != "0" {
		clusterOption := objectsUI.SelectOption{
			Id:       "0",
			Text:     "===全部集群===",
			ParentID: "0",
		}
		clusterOptions = append(clusterOptions, clusterOption)
		var k8sclusterEntity sysadmObjects.ObjectEntity
		k8sclusterEntity = New()
		conditions := make(map[string]string, 0)
		var emptyString []string
		conditions["isDeleted"] = "='0'"
		conditions["azid"] = "='" + selectedAZ + "'"
		clusterList, e := k8sclusterEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, make(map[string]string))
		if e != nil {
			return e
		}
		for _, line := range clusterList {
			clusterData, ok := line.(K8sclusterSchema)
			if !ok {
				return fmt.Errorf("the data is not cluster schema")
			}
			clusterOption := objectsUI.SelectOption{
				Id:       clusterData.Id,
				Text:     clusterData.CnName,
				ParentID: strconv.Itoa(int(clusterData.Dcid)),
			}
			clusterOptions = append(clusterOptions, clusterOption)
		}
	} else {
		clusterOption := objectsUI.SelectOption{
			Id:       "0",
			Text:     "===选择集群===",
			ParentID: "0",
		}
		clusterOptions = append(clusterOptions, clusterOption)
	}
	clusterSelect := objectsUI.SelectData{Title: "集群", SelectedId: selectedCluster, Options: clusterOptions}
	tplData["clusterSelect"] = clusterSelect

	return nil
}

func listHostData(selectedDC, selectedAZ, selectedCluster string, startPos int,
	requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	var hostEntity sysadmObjects.ObjectEntity
	hostEntity, e := HostNew(runData.dbConf, runData.workingRoot)
	if e != nil {
		return 0, dataList, e
	}
	conditions := make(map[string]string, 0)
	if selectedCluster != "0" {
		conditions["k8sclusterid"] = "='" + selectedCluster + "'"
	} else {
		if selectedAZ != "0" {
			conditions["azid"] = "='" + selectedAZ + "'"
		} else {
			if selectedDC != "0" {
				conditions["dcid"] = "='" + selectedDC + "'"
			}
		}
	}
	hostList, e := hostEntity.GetObjectList("", []string{}, []string{}, conditions, 0, 0, make(map[string]string))
	if e != nil {
		return 0, dataList, e
	}
	totalNum := len(hostList)
	if totalNum < 1 {
		return 0, dataList, nil
	}

	orderfield := requestData["orderfield"]
	direction := requestData["direction"]
	if orderfield == "" {
		orderfield = "TD1"
	}
	if direction == "" || (direction != "0" && direction != "1") {
		direction = "1"
	}

	sortHostListData(hostList, direction, orderfield, hostAllOrderFields)
	endCount := totalNum
	if endCount > startPos+runData.pageInfo.NumPerPage {
		endCount = startPos + runData.pageInfo.NumPerPage
	}
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = datacenter.New()
	var azEntity sysadmObjects.ObjectEntity
	azEntity = sysadmAZ.New()
	var osEntity sysadmObjects.ObjectEntity
	osEntity = sysadmOS.New()
	var versionEntity sysadmObjects.ObjectEntity
	versionEntity = sysadmVersion.New()
	var k8sEntity sysadmObjects.ObjectEntity
	k8sEntity = New()
	for i := startPos; i < endCount; i++ {
		interfaceData := hostList[i]
		hostData, ok := interfaceData.(HostSchema)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not host schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = hostData.HostId
		lineMap["TD1"] = hostData.Hostname

		lineMap["TD2"] = "未知数据中心"
		dcID := strings.TrimSpace(strconv.Itoa(int(hostData.Dcid)))
		if dcID != "" && dcID != "0" {
			dcInfo, e := dcEntity.GetObjectInfoByID(dcID)
			if e == nil {
				dcData, ok := dcInfo.(datacenter.DatacenterSchema)
				if ok {
					lineMap["TD2"] = dcData.CnName
				}
			}
		}

		lineMap["TD3"] = "未知可用区"
		azID := strings.TrimSpace(strconv.Itoa(int(hostData.Azid)))
		if azID != "" && azID != "0" {
			azInfo, e := azEntity.GetObjectInfoByID(azID)
			if e == nil {
				azData, ok := azInfo.(sysadmAZ.AvailablezoneSchema)
				if ok {
					lineMap["TD3"] = azData.CnName
				}
			}
		}

		osStr := "未知操作系统"
		osID := strings.TrimSpace(strconv.Itoa(hostData.OSID))
		if osID != "" && osID != "0" {
			osInfo, e := osEntity.GetObjectInfoByID(osID)
			if e == nil {
				osData, ok := osInfo.(sysadmOS.OSSchema)
				if ok {
					osStr = osData.Name
				}
			}
		}

		verStr := "未知版本"
		verID := strings.TrimSpace(strconv.Itoa(hostData.OSVersionID))
		if verID != "" && verID != "0" {
			versionInfo, e := versionEntity.GetObjectInfoByID(verID)
			if e == nil {
				versionData, ok := versionInfo.(sysadmVersion.VersionSchema)
				if ok {
					verStr = versionData.Name
				}
			}
		}
		lineMap["TD4"] = osStr + "/" + verStr

		k8sClusterID := strings.TrimSpace(hostData.K8sClusterID)
		lineMap["TD5"] = "未知集群"
		if k8sClusterID != "" && k8sClusterID != "0" {
			k8sInfo, e := k8sEntity.GetObjectInfoByID(hostData.K8sClusterID)
			if e == nil {
				k8sData, ok := k8sInfo.(K8sclusterSchema)
				if ok {
					lineMap["TD5"] = k8sData.CnName
				}
			}
		}

		popmenuitems := "0"
		hostStatus := ""
		switch hostData.Status {
		case string(HostStatusRunning):
			hostStatus = "运行中"
			popmenuitems = popmenuitems + ",1,6"
		case string(HostStatusMaintenance):
			hostStatus = "维护中"
			popmenuitems = popmenuitems + ",1,6"
		case string(HostStatusOffline):
			hostStatus = "已下线"
			popmenuitems = popmenuitems + ",6"
		case string(HostStatusDeleted):
			hostStatus = "已删除"
		case string(HostStatusReady):
			hostStatus = "正常"
			popmenuitems = popmenuitems + ",1,6"
		}
		if hostData.Status == string(HostStatusMemoryPressure) {
			if hostStatus == "" {
				hostStatus = "内存压力"
			} else {
				hostStatus = hostStatus + "," + "内存压力"
			}
		}
		if hostData.Status == string((HostStatusDiskPressure)) {
			if hostStatus == "" {
				hostStatus = "磁盘压力"
			} else {
				hostStatus = hostStatus + "," + "磁盘压力"
			}
		}
		if hostData.Status == string((HostPIDPressure)) {
			if hostStatus == "" {
				hostStatus = "PID压力"
			} else {
				hostStatus = hostStatus + "," + "PID压力"
			}
		}
		if hostData.Status == string((HostNetworkUnavailable)) {
			if hostStatus == "" {
				hostStatus = "网络不可用"
			} else {
				hostStatus = hostStatus + "," + "网络不可用"
			}
		}
		lineMap["TD6"] = hostStatus
		popmenuitems = buildPopItemList(popmenuitems, hostData)
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortHostListData(data []interface{}, direction, orderField string, allOrderFields map[string]objectsUI.SortBy) {
	for k, fn := range allOrderFields {
		if k == orderField {
			if direction == "1" {
				sort.Sort(objectsUI.SortData{Data: data, By: fn})
			} else {
				sort.Sort(sort.Reverse(objectsUI.SortData{Data: data, By: fn}))
			}
			break
		}
	}

	return
}

func buildPopItemList(popItems string, hostData HostSchema) string {
	clusterID := strings.TrimSpace(hostData.K8sClusterID)
	if clusterID == "" {
		popItems = popItems + ",7"
	}
	var k8sEntity sysadmObjects.ObjectEntity
	k8sEntity = New()
	k8sInfo, e := k8sEntity.GetObjectInfoByID(hostData.K8sClusterID)
	if e != nil {
		return popItems
	}
	k8sData, ok := k8sInfo.(K8sclusterSchema)
	if !ok {
		return popItems
	}

	restConf, e := sysadmK8sClient.BuildConfigFromParasWithConnectType(k8sData.ConnectType, k8sData.Apiserver, k8sData.K8sClusterID, k8sData.ClusterUser, "",
		k8sData.Ca, k8sData.Cert, k8sData.Key, k8sData.Token, k8sData.KubeConfig)
	if e != nil {
		return popItems
	}
	k8sClient, e := kubernetes.NewForConfig(restConf)
	if e != nil {
		return popItems
	}
	hostName := strings.TrimSpace(hostData.Hostname)
	nodeData, e := k8sClient.CoreV1().Nodes().Get(context.Background(), hostName, metav1.GetOptions{})
	if e != nil {
		return popItems
	}

	unschedule := nodeData.Spec.Unschedulable
	if unschedule {
		popItems = popItems + ",4,5"
	} else {
		popItems = popItems + ",2,3,5"
	}

	return popItems
}

func getAZOptionsByDCHandler(c *sysadmServer.Context) {
	// order fields data of cluster list page
	var errs []sysadmLog.Sysadmerror
	requestData, e := utils.NewGetRequestData(c, []string{"objID"})
	if e != nil || requestData["objID"] == "0" {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(7001400012, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		response := apiutils.BuildResponseDataForError(7001400012, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	var azEntity sysadmObjects.ObjectEntity
	azEntity = sysadmAZ.New()
	conditions := make(map[string]string, 0)
	conditions["isDeleted"] = "='0'"
	conditions["datacenterid"] = "='" + requestData["objID"] + "'"
	azList, e := azEntity.GetObjectList("", []string{}, []string{}, conditions, 0, 0, make(map[string]string))
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(7001400013, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		response := apiutils.BuildResponseDataForError(7001400013, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	msg := "0:===选择可用区==="
	for _, line := range azList {
		lineAZData := line.(sysadmAZ.AvailablezoneSchema)
		lineStr := strconv.Itoa(int(lineAZData.Id)) + ":" + lineAZData.CnName
		msg = msg + "," + lineStr
	}

	response := apiutils.BuildResponseDataForSuccess(msg)
	c.JSON(http.StatusOK, response)

	return
}

func getClusterOptionsByAZ(c *sysadmServer.Context) {
	// order fields data of cluster list page
	var errs []sysadmLog.Sysadmerror
	requestData, e := utils.NewGetRequestData(c, []string{"objID"})
	if e != nil || requestData["objID"] == "0" {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(7001400014, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		response := apiutils.BuildResponseDataForError(7001400014, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	var k8sclusterEntity sysadmObjects.ObjectEntity
	k8sclusterEntity = New()
	conditions := make(map[string]string, 0)
	var emptyString []string
	conditions["isDeleted"] = "='0'"
	conditions["azid"] = "='" + requestData["objID"] + "'"
	clusterList, e := k8sclusterEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, make(map[string]string))
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(7001400015, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		response := apiutils.BuildResponseDataForError(7001400015, "数据处理错误")
		c.JSON(http.StatusOK, response)
		return
	}

	msg := "0:===选择集群==="
	for _, line := range clusterList {
		lineClusterData := line.(K8sclusterSchema)
		lineStr := lineClusterData.Id + ":" + lineClusterData.CnName
		msg = msg + "," + lineStr
	}

	response := apiutils.BuildResponseDataForSuccess(msg)
	c.JSON(http.StatusOK, response)

	return
}

func addHostFormHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(7001400016, "debug", "displaying add host form page"))

	_, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "您没有登陆或没有权限执行本操作", runData.logEntity, 7001400017, errs, e)
		return
	}

	requestKeys := []string{"dcID", "azID", "clusterID"}
	requestData, e := objectsUI.GetRequestData(c, requestKeys)
	if e != nil {
		objectsUI.OutPutMsg(c, "", "系统出错，请稍后再试或联系系统管理员", runData.logEntity, 7001400018, errs, e)
		return
	}

	if requestData["dcID"] == "" || requestData["dcID"] == "0" || requestData["azID"] == "" || requestData["azID"] == "" {
		objectsUI.OutPutMsg(c, "", "操作错误，请稍后再试或联系系统管理员", runData.logEntity, 7001400019, errs, e)
		return
	}

	// 余下的仅显示一个帮助信息
	return
}

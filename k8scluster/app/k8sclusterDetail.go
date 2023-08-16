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
	"strconv"
	"strings"
	az "sysadm/availablezone/app"
	datacenter "sysadm/datacenter/app"
	"sysadm/k8sclient"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	"sysadm/sysadmLog"
	"sysadm/user"
	userApp "sysadm/user/app"
	"sysadm/utils"
)

func clusterDetailHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	errs = append(errs, sysadmLog.NewErrorWithStringLevel(700070001, "debug", "now display cluster details"))
	messageTemplateFile := "poppageerror.html"
	detailTemplateFile := "objectDetailPage.html"
	baseUri := "/" + DefaultModuleName + "/"
	redirectUri := "list"

	messageTplData := make(map[string]interface{}, 0)
	// get userid
	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700070002, "error", "user should login %s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "你没有登录或没有权限查看集群详情"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}

	requestData, e := utils.NewGetRequestData(c, []string{"objID"})
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700070003, "error", "get request data error %s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "参数错误，请确认操作正确"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}

	clusterID := strings.TrimSpace(requestData["objID"])
	if clusterID == "" {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700070004, "error", "cluster ID is empty"))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "你没有选择要查看祥情的集群"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}

	// get cluster infromation by clusterID
	var clusterEntity sysadmObjects.ObjectEntity
	clusterEntity = New()

	data, e := clusterEntity.GetObjectInfoByID(clusterID)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700070005, "error", "there is an error occurred when get cluster information:%s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "查询集群祥情出错，请稍后再试或联系系统管理员"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}
	clusterData := data.(K8sclusterSchema)

	var tplDataLines []objectsUI.LineDataForDetail
	separator := objectsUI.ItemForDetail{Label: "基础信息", IsSeparator: true}
	lineData := objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{separator}}
	tplDataLines = append(tplDataLines, lineData)

	lines, e := buildTplDataForBasicInfo(clusterData)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700070006, "error", "there is an error occurred when build cluster information:%s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "查询集群祥情出错，请稍后再试或联系系统管理员"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}
	tplDataLines = append(tplDataLines, lines...)

	separator = objectsUI.ItemForDetail{Label: "集群连接信息", IsSeparator: true}
	lineData = objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{separator}}
	tplDataLines = append(tplDataLines, lineData)

	lines = buildTplDataForConnectInfo(clusterData)
	tplDataLines = append(tplDataLines, lines...)

	separator = objectsUI.ItemForDetail{Label: "业务统计数据", IsSeparator: true}
	lineData = objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{separator}}
	tplDataLines = append(tplDataLines, lineData)

	lines, e = buildTplDataForCountInfo(clusterData)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700070010, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "查询集群祥情出错，请稍后再试或联系系统管理员"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}
	tplDataLines = append(tplDataLines, lines...)

	tplData, e := objectsUI.InitTemplateForShowObjectDetails("集群管理", "集群详情", redirectUri, baseUri)
	if e != nil {
		errs = append(errs, sysadmLog.NewErrorWithStringLevel(700070019, "error", "%s", e))
		runData.logEntity.LogErrors(errs)
		messageTplData["errormessage"] = "查询集群祥情出错，请稍后再试或联系系统管理员"
		c.HTML(http.StatusOK, messageTemplateFile, messageTplData)
		return
	}
	tplData["data"] = tplDataLines

	c.HTML(http.StatusOK, detailTemplateFile, tplData)
}

func buildTplDataForBasicInfo(data K8sclusterSchema) ([]objectsUI.LineDataForDetail, error) {
	var tplDataLines []objectsUI.LineDataForDetail

	clusterID := objectsUI.ItemForDetail{Label: "集群ID: ", Value: data.Id}
	cnName := objectsUI.ItemForDetail{Label: "中文名: ", Value: data.CnName}
	enName := objectsUI.ItemForDetail{Label: "Englisth Name: ", Value: data.EnName}

	lineData := objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{clusterID, cnName, enName}}
	tplDataLines = append(tplDataLines, lineData)

	dcID := utils.Interface2String(data.Dcid)

	// try to get datacenter name
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = datacenter.New()
	dataDC, e := dcEntity.GetObjectInfoByID(dcID)
	if e != nil {
		return tplDataLines, e
	}
	dcData := dataDC.(datacenter.DatacenterSchema)
	dcName := strings.TrimSpace(dcData.CnName)

	// try to get availablezone name
	azID := utils.Interface2String(data.Azid)
	var azEntity sysadmObjects.ObjectEntity
	azEntity = az.New()
	dataAZ, e := azEntity.GetObjectInfoByID(azID)
	if e != nil {
		return tplDataLines, e
	}
	azData := dataAZ.(az.AvailablezoneSchema)
	azName := strings.TrimSpace(azData.CnName)

	clusterStatus := GetStatusText(data.Status)
	isDeletedStr := "正常"
	if data.IsDeleted != 0 {
		isDeletedStr = "已删除"
	}

	dc := objectsUI.ItemForDetail{Label: "所属数据中心: ", Value: dcName}
	az := objectsUI.ItemForDetail{Label: "所属可用区: ", Value: azName}
	status := objectsUI.ItemForDetail{Label: "集群状态: ", Value: clusterStatus}
	isDeleted := objectsUI.ItemForDetail{Label: "删除状态: ", Value: isDeletedStr}
	lineData = objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{dc, az, status, isDeleted}}
	tplDataLines = append(tplDataLines, lineData)

	createBy := utils.Interface2String(data.CreateBy)
	updateBy := utils.Interface2String(data.UpdateBy)
	var userEntity sysadmObjects.ObjectEntity
	userEntity = userApp.New()
	dataCreateUser, e := userEntity.GetObjectInfoByID(createBy)
	if e != nil {
		return tplDataLines, e
	}
	userDataCreate := dataCreateUser.(userApp.UserSchema)
	createByUser := userDataCreate.Username
	updateByUser := ""
	updateTime := ""
	if updateBy != "0" {
		dataUpdateUser, e := userEntity.GetObjectInfoByID(updateBy)
		if e != nil {
			return tplDataLines, e
		}
		userDataUpdate := dataUpdateUser.(userApp.UserSchema)
		updateByUser = userDataUpdate.Username
		updateTime = data.UpdateTime
	}
	createTime := objectsUI.ItemForDetail{Label: "创建时间: ", Value: data.CreateTime}
	createByData := objectsUI.ItemForDetail{Label: "创建者: ", Value: createByUser}
	updateTimeData := objectsUI.ItemForDetail{Label: "更新时间: ", Value: updateTime}
	updateByData := objectsUI.ItemForDetail{Label: "更新者: ", Value: updateByUser}
	lineData = objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{createTime, createByData, updateTimeData, updateByData}}
	tplDataLines = append(tplDataLines, lineData)

	version := objectsUI.ItemForDetail{Label: "K8S版本: ", Value: data.Version}
	cri := objectsUI.ItemForDetail{Label: "容器运行时: ", Value: data.Cri}
	podcidr := objectsUI.ItemForDetail{Label: "Pod CIDR: ", Value: data.Podcidr}
	svcCIDR := objectsUI.ItemForDetail{Label: "Service CIDR: ", Value: data.Servicecidr}
	lineData = objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{version, cri, podcidr, svcCIDR}}
	tplDataLines = append(tplDataLines, lineData)

	dutyTel := objectsUI.ItemForDetail{Label: "负责人联系电话: ", Value: data.DutyTel}
	lineData = objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{dutyTel}}
	tplDataLines = append(tplDataLines, lineData)

	remarkTile := objectsUI.ItemForDetail{Label: "描述: ", Value: ""}
	lineData = objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{remarkTile}}
	tplDataLines = append(tplDataLines, lineData)

	remarkValue := objectsUI.ItemForDetail{Label: "", Value: data.Remark}
	lineData = objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{remarkValue}}
	tplDataLines = append(tplDataLines, lineData)

	return tplDataLines, nil
}

func buildTplDataForConnectInfo(data K8sclusterSchema) []objectsUI.LineDataForDetail {
	var lines []objectsUI.LineDataForDetail

	user := objectsUI.ItemForDetail{Label: "连接用户: ", Value: data.ClusterUser}
	lineData := objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{user}}
	lines = append(lines, lineData)

	ca := objectsUI.ItemForDetail{Label: "根证书: ", Value: "下载根证书", ActionUrl: "download?type=ca&objID=" + data.Id}
	cert := objectsUI.ItemForDetail{Label: "证书: ", Value: "下载证书", ActionUrl: "download?type=cert&objID=" + data.Id}
	key := objectsUI.ItemForDetail{Label: "密钥: ", Value: "下载密钥", ActionUrl: "download?type=cert&objID=" + data.Id}
	lineData = objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{ca, cert, key}}
	lines = append(lines, lineData)

	return lines
}

func buildTplDataForCountInfo(data K8sclusterSchema) ([]objectsUI.LineDataForDetail, error) {
	var lines []objectsUI.LineDataForDetail

	restConf, e := k8sclient.BuildConfigFromParametes([]byte(data.Ca), []byte(data.Cert), []byte(data.Key), data.Apiserver, data.Id, data.ClusterUser, "default")
	if e != nil {
		return lines, e
	}

	nodeCountForCP, e := k8sclient.GetNodeCount(restConf, k8sclient.NodeRoleCP)
	if e != nil {
		return lines, e
	}
	nodeCountForWK, e := k8sclient.GetNodeCount(restConf, k8sclient.NodeRoleWK)
	if e != nil {
		return lines, e
	}
	cpStr := strconv.Itoa(int(nodeCountForCP.Total)) + "/" + strconv.Itoa(int(nodeCountForCP.Ready)) + "/" + strconv.Itoa(int(nodeCountForCP.Unready))
	wkStr := strconv.Itoa(int(nodeCountForWK.Total)) + "/" + strconv.Itoa(int(nodeCountForWK.Ready)) + "/" + strconv.Itoa(int(nodeCountForWK.Unready))
	cpNode := objectsUI.ItemForDetail{Label: "CP节点统计(总数/健康/异常): ", Value: cpStr}
	wkNode := objectsUI.ItemForDetail{Label: "Work节点统计(总数/健康/异常): ", Value: wkStr}
	lineData := objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{cpNode, wkNode}}
	lines = append(lines, lineData)

	deployCount, e := k8sclient.GetDeploymentCount(restConf, "")
	if e != nil {
		return lines, e
	}
	deployStr := strconv.Itoa(int(deployCount.Total)) + "/" + strconv.Itoa(int(deployCount.Ready)) + "/" + strconv.Itoa(int(deployCount.Unready))
	deployData := objectsUI.ItemForDetail{Label: "Deployment统计(总数/健康/异常): ", Value: deployStr}
	statefulSetCount, e := k8sclient.GetStatefulSetCount(restConf, "")
	if e != nil {
		return lines, e
	}
	stsStr := strconv.Itoa(int(statefulSetCount.Total)) + "/" + strconv.Itoa(int(statefulSetCount.Ready)) + "/" + strconv.Itoa(int(statefulSetCount.Unready))
	stsData := objectsUI.ItemForDetail{Label: "StatefulSet统计(总数/健康/异常): ", Value: stsStr}
	lineData = objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{deployData, stsData}}
	lines = append(lines, lineData)

	dsCount, e := k8sclient.GetDaemonSetCount(restConf, "")
	if e != nil {
		return lines, e
	}
	jobCount, e := k8sclient.GetJobCount(restConf, "")
	if e != nil {
		return lines, e
	}
	dsStr := strconv.Itoa(int(dsCount.Total)) + "/" + strconv.Itoa(int(dsCount.Ready)) + "/" + strconv.Itoa(int(dsCount.Unready))
	jobStr := strconv.Itoa(int(jobCount.Total)) + "/" + strconv.Itoa(int(jobCount.Ready))
	dsData := objectsUI.ItemForDetail{Label: "daemonSet统计(总数/健康/异常): ", Value: dsStr}
	jobData := objectsUI.ItemForDetail{Label: "job统计(总数/已完成): ", Value: jobStr}
	lineData = objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{dsData, jobData}}
	lines = append(lines, lineData)

	podCount, e := k8sclient.GetPodCount(restConf, "")
	if e != nil {
		return lines, e
	}
	podStr := strconv.Itoa(int(podCount.Total)) + "/" + strconv.Itoa(int(podCount.Ready)) + "/" + strconv.Itoa(int(podCount.Unready))
	podData := objectsUI.ItemForDetail{Label: "pod统计(总数/健康/异常): ", Value: podStr}
	lineData = objectsUI.LineDataForDetail{Items: []objectsUI.ItemForDetail{podData}}
	lines = append(lines, lineData)

	return lines, nil
}

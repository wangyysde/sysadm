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
	"k8s.io/client-go/rest"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"sysadm/k8sclient"
	sysadmObjects "sysadm/objects/app"
	"sysadm/sysadmLog"
	"sysadm/sysadmapi/apiutils"
	"sysadm/user"
	"sysadm/utils"
)

func validCNNameHandler(c *sysadmServer.Context) {
	var response apiutils.ApiResponseData

	islogin, _, _ := user.IsLogin(c, runData.sessionName)
	if !islogin {
		response = apiutils.BuildResponseDataForError(700070007, "您没有登录或者没有权限执行本操作")
		c.JSON(http.StatusOK, response)
		return
	}

	requestData, e := utils.NewGetRequestData(c, []string{"objvalue"})

	if e != nil || !validCNName(requestData["objvalue"]) {
		response = apiutils.BuildResponseDataForError(700070001, "集群的中文名称为必埴项，且其长度不得大于255个字符")
	} else {
		response = apiutils.BuildResponseDataForSuccess("ok")
	}

	c.JSON(http.StatusOK, response)
}

func validENNameHandler(c *sysadmServer.Context) {
	var response apiutils.ApiResponseData

	islogin, _, _ := user.IsLogin(c, runData.sessionName)
	if !islogin {
		response = apiutils.BuildResponseDataForError(700070008, "您没有登录或者没有权限执行本操作")
		c.JSON(http.StatusOK, response)
		return
	}

	requestData, e := utils.NewGetRequestData(c, []string{"objvalue"})
	if e != nil || !validENName(requestData["objvalue"]) {
		response = apiutils.BuildResponseDataForError(700070002, "集群的英文名称为选埴项，且其长度不得大于255个字符")
	} else {
		response = apiutils.BuildResponseDataForSuccess("ok")
	}

	c.JSON(http.StatusOK, response)
}

func validApiserverAddressHandler(c *sysadmServer.Context) {
	var response apiutils.ApiResponseData

	islogin, _, _ := user.IsLogin(c, runData.sessionName)
	if !islogin {
		response = apiutils.BuildResponseDataForError(700070009, "您没有登录或者没有权限执行本操作")
		c.JSON(http.StatusOK, response)
		return
	}

	requestData, e := utils.NewGetRequestData(c, []string{"objvalue"})
	if e != nil || !validApiserverAddress(requestData["objvalue"]) {
		response = apiutils.BuildResponseDataForError(700070003, "集群的kube-apiserver地址为必埴项，且其长度不得大于255个字符.形式为为x.x.x.x:6443")
	} else {
		response = apiutils.BuildResponseDataForSuccess("ok")
	}

	c.JSON(http.StatusOK, response)
}

func validClusterUserHandler(c *sysadmServer.Context) {
	var response apiutils.ApiResponseData

	islogin, _, _ := user.IsLogin(c, runData.sessionName)
	if !islogin {
		response = apiutils.BuildResponseDataForError(7000700010, "您没有登录或者没有权限执行本操作")
		c.JSON(http.StatusOK, response)
		return
	}

	requestData, e := utils.NewGetRequestData(c, []string{"objvalue"})
	if e != nil || !validClusterUser(requestData["objvalue"]) {
		response = apiutils.BuildResponseDataForError(700070005, "集群用户名必埴项，且其长度不得大于255个字符")
	} else {
		response = apiutils.BuildResponseDataForSuccess("ok")
	}

	c.JSON(http.StatusOK, response)
}

func validDutyTelHandler(c *sysadmServer.Context) {
	var response apiutils.ApiResponseData

	islogin, _, _ := user.IsLogin(c, runData.sessionName)
	if !islogin {
		response = apiutils.BuildResponseDataForError(700070011, "您没有登录或者没有权限执行本操作")
		c.JSON(http.StatusOK, response)
		return
	}

	requestData, e := utils.NewGetRequestData(c, []string{"objvalue"})
	if e != nil || !validDutyTel(requestData["objvalue"]) {
		response = apiutils.BuildResponseDataForError(700070006, "值班电话为选填项，且其长度不得大于20个字符")
	} else {
		response = apiutils.BuildResponseDataForSuccess("ok")
	}

	c.JSON(http.StatusOK, response)
}

func addPostHandler(c *sysadmServer.Context) {
	var errs []sysadmLog.Sysadmerror
	var response apiutils.ApiResponseData

	errs = append(errs, sysadmLog.NewErrorWithStringLevel(700070012, "debug", "try to add cluster information to DB"))

	islogin, userid, e := user.IsLogin(c, runData.sessionName)
	if !islogin {
		response = apiutils.BuildResponseDataForError(700070013, "您没有登录或者没有权限执行本操作")
		err := sysadmLog.NewErrorWithStringLevel(700070013, "error", "%s", e)
		errs = append(errs, err)
		c.JSON(http.StatusOK, response)
		runData.logEntity.LogErrors(errs)
		return
	}

	formData, e := utils.GetMultipartData(c, []string{"dcid", "azid", "cnName", "enName", "connectType", "apiserver", "ca", "cert", "token", "kubeConfig", "key", "dutyTel", "remark"})
	if e != nil {
		response = apiutils.BuildResponseDataForError(700070014, "数据处理错误，请稍后再试或联系平台管理员")
		err := sysadmLog.NewErrorWithStringLevel(700070014, "error", "%s", e)
		errs = append(errs, err)
		c.JSON(http.StatusOK, response)
		runData.logEntity.LogErrors(errs)
		return
	}

	dcid := strings.TrimSpace(formData["dcid"].([]string)[0])
	azid := strings.TrimSpace(formData["azid"].([]string)[0])
	if dcid == "" || dcid == "0" || azid == "" || azid == "0" {
		response = apiutils.BuildResponseDataForError(700070015, "请选择添加的集群所属于的数据中心和可用区")
		err := sysadmLog.NewErrorWithStringLevel(700070015, "info", "no datacenter or AZ has be selected when add cluster")
		errs = append(errs, err)
		c.JSON(http.StatusOK, response)
		runData.logEntity.LogErrors(errs)
		return
	}
	dcidInt, e := strconv.Atoi(dcid)
	azidInt, e1 := strconv.Atoi(azid)
	if e != nil || e1 != nil {
		response = apiutils.BuildResponseDataForError(700070015, "数据中心或可用区数据不合法，请确认操作是否正确")
		err := sysadmLog.NewErrorWithStringLevel(700070015, "info", "dcid or azid is not valid")
		errs = append(errs, err)
		c.JSON(http.StatusOK, response)
		runData.logEntity.LogErrors(errs)
		return
	}

	cnName := strings.TrimSpace(formData["cnName"].([]string)[0])
	if !validCNName(cnName) {
		response = apiutils.BuildResponseDataForError(700070016, "集群的中文名称为必埴项，且其长度不得大于255个字符")
		err := sysadmLog.NewErrorWithStringLevel(700070016, "info", "cnName is not valid when add cluster")
		errs = append(errs, err)
		c.JSON(http.StatusOK, response)
		runData.logEntity.LogErrors(errs)
		return
	}

	enName := strings.TrimSpace(formData["enName"].([]string)[0])
	if !validENName(enName) {
		response = apiutils.BuildResponseDataForError(700070017, "集群的英文名称为选埴项，且其长度不得大于255个字符")
		err := sysadmLog.NewErrorWithStringLevel(700070017, "info", "enName is not valid when add cluster")
		errs = append(errs, err)
		c.JSON(http.StatusOK, response)
		runData.logEntity.LogErrors(errs)
		return
	}

	connectType := strings.TrimSpace(formData["connectType"].([]string)[0])
	var connectTypeDefined []string
	for k, _ := range k8sClusterConnectType {
		connectTypeDefined = append(connectTypeDefined, k)
	}
	if !utils.FoundStrInSlice(connectTypeDefined, connectType, true) {
		response = apiutils.BuildResponseDataForError(700070018, "集群连接类型不正确，请确认你操作正确或联系系统管理员")
		err := sysadmLog.NewErrorWithStringLevel(700070018, "error", "connect type %s is not valid", connectType)
		errs = append(errs, err)
		c.JSON(http.StatusOK, response)
		runData.logEntity.LogErrors(errs)
		return
	}

	token := strings.TrimSpace(formData["token"].([]string)[0])
	if connectType == "1" && token == "" {
		response = apiutils.BuildResponseDataForError(700070019, "集群连接类型已选择为通过token连接，但是token值为空")
		err := sysadmLog.NewErrorWithStringLevel(700070019, "error", "content of token is empty")
		errs = append(errs, err)
		c.JSON(http.StatusOK, response)
		runData.logEntity.LogErrors(errs)
		return
	}

	apiserver := ""
	if connectType == "1" || connectType == "2" {
		apiserver = strings.TrimSpace(formData["apiserver"].([]string)[0])
		if apiserver != "" && !validApiserverAddress(apiserver) {
			response = apiutils.BuildResponseDataForError(700070020, "集群的kube-apiserver地址为必埴项，且其长度不得大于255个字符.形式为为x.x.x.x:6443")
			err := sysadmLog.NewErrorWithStringLevel(700070020, "info", "apiServer address is not valid when add cluster")
			errs = append(errs, err)
			c.JSON(http.StatusOK, response)
			runData.logEntity.LogErrors(errs)
			return
		}
	}

	dutyTel := strings.TrimSpace(formData["dutyTel"].([]string)[0])
	if !validDutyTel(dutyTel) {
		response = apiutils.BuildResponseDataForError(700070021, "值班电话为选填项，且其长度不得大于20个字符")
		err := sysadmLog.NewErrorWithStringLevel(700070021, "info", "duty Tel is not valid when add cluster")
		errs = append(errs, err)
		c.JSON(http.StatusOK, response)
		runData.logEntity.LogErrors(errs)
		return
	}

	ca := strings.TrimSpace(formData["ca"].([]string)[0])
	cert := strings.TrimSpace(formData["cert"].([]string)[0])
	key := strings.TrimSpace(formData["key"].([]string)[0])
	if connectType == "0" && (ca == "" || cert == "" || key == "") {
		response = apiutils.BuildResponseDataForError(700070022, "连接集群所需要的证书和密钥文件必须填写")
		err := sysadmLog.NewErrorWithStringLevel(700070022, "info", "ca, cert or key is empty")
		errs = append(errs, err)
		c.JSON(http.StatusOK, response)
		runData.logEntity.LogErrors(errs)
		return
	}

	if connectType == "1" && ca == "" {
		response = apiutils.BuildResponseDataForError(700070023, "使用Token连接时,根证书用必须填写")
		err := sysadmLog.NewErrorWithStringLevel(700070023, "info", "ca is empty")
		errs = append(errs, err)
		c.JSON(http.StatusOK, response)
		runData.logEntity.LogErrors(errs)
		return
	}

	var kubeConfigContent []byte
	if connectType == "2" {
		if formData["kubeConfig"] == nil {
			response = apiutils.BuildResponseDataForError(700070026, "连接方式已选择为通过kubeConfig方式连接，但是没有上传文件")
			err := sysadmLog.NewErrorWithStringLevel(700070026, "info", "kubeconfig file is empty")
			errs = append(errs, err)
			c.JSON(http.StatusOK, response)
			runData.logEntity.LogErrors(errs)
			return
		}
		var e error = nil
		kubeConfigContent, e = utils.ReadUploadedFile(formData["kubeConfig"].(*multipart.FileHeader))
		if e != nil {
			response = apiutils.BuildResponseDataForError(700070027, "上传kubeConfig文件失败,请确认操作是否正确")
			err := sysadmLog.NewErrorWithStringLevel(700070027, "error", "upload kubeconfig file  error %s", e)
			errs = append(errs, err)
			c.JSON(http.StatusOK, response)
			runData.logEntity.LogErrors(errs)
			return
		}
	}

	k8sClusterID, clusterID, kubeVersion, cri, podcidr, svccidr, restConf, e := getClusterInfo(dcidInt, azidInt, token, apiserver, ca, cert, key,
		string(kubeConfigContent), connectType)
	if e != nil {
		response = apiutils.BuildResponseDataForError(700070028, "连接集群出错，请确认集群连接数据是否正确")
		err := sysadmLog.NewErrorWithStringLevel(700070028, "error", "%s", e)
		errs = append(errs, err)
		c.JSON(http.StatusOK, response)
		runData.logEntity.LogErrors(errs)
		return
	}

	conditions := make(map[string]string, 0)
	conditions["k8sClusterID"] = "='" + k8sClusterID + "'"

	// try to add cluster data into DB
	var clusterEntity sysadmObjects.ObjectEntity
	clusterEntity = New()
	var emptyString []string
	existCount, e := clusterEntity.GetObjectCount("", emptyString, emptyString, conditions)
	if existCount > 0 || e != nil {
		response = apiutils.BuildResponseDataForError(700070026, "需要添加的集群已存在或查询出现错误")
		err := sysadmLog.NewErrorWithStringLevel(700070026, "error", "k8s cluster exist or got an error %s", e)
		errs = append(errs, err)
		c.JSON(http.StatusOK, response)
		runData.logEntity.LogErrors(errs)
		return
	}
	clusterUser := k8sclient.GetClusterDefaultUser()
	addData := K8sclusterSchema{
		Id:           clusterID,
		K8sClusterID: k8sClusterID,
		Dcid:         uint(dcidInt),
		Azid:         uint(azidInt),
		CnName:       cnName,
		EnName:       enName,
		Apiserver:    apiserver,
		ClusterUser:  clusterUser,
		Ca:           ca,
		Cert:         cert,
		Key:          key,
		ConnectType:  connectType,
		Token:        token,
		KubeConfig:   string(kubeConfigContent),
		Version:      kubeVersion,
		Cri:          cri,
		Podcidr:      podcidr,
		Servicecidr:  svccidr,
		DutyTel:      dutyTel,
		Status:       1,
		IsDeleted:    0,
		CreateBy:     uint(userid),
		Remark:       formData["remark"].([]string)[0],
	}

	e = clusterEntity.AddObject(addData)
	if e != nil {
		response = apiutils.BuildResponseDataForError(700070027, fmt.Sprintf("添加集群信息出错，错误信息为:%s", e))
		err := sysadmLog.NewErrorWithStringLevel(700070027, "error", "add cluster information error %s", e)
		errs = append(errs, err)
		c.JSON(http.StatusOK, response)
		runData.logEntity.LogErrors(errs)
		return
	}

	response = apiutils.BuildResponseDataForSuccess("集群已经添加成功")
	err := sysadmLog.NewErrorWithStringLevel(700070029, "info", "cluster infromation has be added")
	errs = append(errs, err)
	c.JSON(http.StatusOK, response)

	e = tryReconizeHostsInCluster(addData, restConf, userid)
	if e != nil {
		err := sysadmLog.NewErrorWithStringLevel(700070030, "error", "add host to db error %s", e)
		errs = append(errs, err)
	}

	runData.logEntity.LogErrors(errs)
}

// 为了避免重复添加，需要获取K8S集群的ID进行判断。当前k8s集群不支持集群层的ID,但是建议使用kube-system命名空间的uid代替集群的ID
// 参见：https://github.com/open-telemetry/semantic-conventions/blob/156f9424fe5d83d8543119224c3af6ae9af518cf/specification/resource/semantic_conventions/k8s.md?plain=1#L28-L51
func getClusterInfo(dcid, azid int, token, apiserver, ca, cert, key, kubeConfig, connectType string) (string, string, string, string, string, string, *rest.Config, error) {
	idData, e := utils.NewWorker(uint64(dcid), uint64(azid))
	if e != nil {
		return "", "", "", "", "", "", nil, e
	}
	clusterID, e := idData.GetID()
	if e != nil {
		return "", "", "", "", "", "", nil, e
	}

	apiserver = strings.TrimSpace(apiserver)
	ca = strings.TrimSpace(ca)
	cert = strings.TrimSpace(cert)
	key = strings.TrimSpace(key)
	restConf, e := k8sclient.BuildConfigFromParasWithConnectType(connectType, apiserver, clusterID, "", "", ca, cert, key, token, kubeConfig)
	if e != nil {
		return "", "", "", "", "", "", nil, e
	}

	kubernetesClusterID, e := k8sclient.GetKubernetesClusterID(restConf)
	if e != nil {
		return "", "", "", "", "", "", nil, e
	}

	kubeVersion, e := k8sclient.GetKubernetesVersion(restConf)
	if e != nil {
		return "", "", "", "", "", "", nil, e
	}

	cri, e := k8sclient.GetCRIInfo(restConf)
	if e != nil {
		return "", "", "", "", "", "", nil, e
	}

	podcidr, e := k8sclient.GetPodCIDR(restConf)
	if e != nil {
		return "", "", "", "", "", "", nil, e
	}

	svccidr, e := k8sclient.GetSvcCIDR(restConf)
	if e != nil {
		return "", "", "", "", "", "", nil, e
	}

	return kubernetesClusterID, clusterID, kubeVersion, cri, podcidr, svccidr, restConf, nil
}

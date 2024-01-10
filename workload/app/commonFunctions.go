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
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/wangyysde/sysadmServer"
	apiCoreV1 "k8s.io/api/core/v1"
	resourceapi "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	"k8s.io/client-go/kubernetes"
	"mime/multipart"
	"sort"
	"strconv"
	"strings"
	datacenter "sysadm/datacenter/app"
	sysadmK8sClient "sysadm/k8sclient"
	sysadmK8sCluster "sysadm/k8scluster/app"
	sysadmObjects "sysadm/objects/app"
	"sysadm/objectsUI"
	"sysadm/user"
	"sysadm/utils"
)

func buildSelectDataWithNs(tplData map[string]interface{}, dcList []interface{}, requestData map[string]string) error {

	selectedDC := strings.TrimSpace(requestData["dcID"])
	if selectedDC == "" {
		selectedDC = "0"
	}
	selectedCluster := strings.TrimSpace(requestData["clusterID"])
	if selectedCluster == "" {
		selectedCluster = "0"
	}

	selectedNamespace := strings.TrimSpace(requestData["namespace"])
	if selectedNamespace == "" {
		selectedNamespace = "0"
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

	// preparing cluster options
	var clusterOptions []objectsUI.SelectOption
	clusterOption := objectsUI.SelectOption{
		Id:       "0",
		Text:     "===选择集群===",
		ParentID: "0",
	}
	clusterOptions = append(clusterOptions, clusterOption)
	if selectedDC != "0" {
		var k8sclusterEntity sysadmObjects.ObjectEntity
		k8sclusterEntity = sysadmK8sCluster.New()
		conditions := make(map[string]string, 0)
		var emptyString []string
		conditions["isDeleted"] = "='0'"
		conditions["dcid"] = "='" + selectedDC + "'"
		clusterList, e := k8sclusterEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, make(map[string]string))
		if e != nil {
			return e
		}

		for _, line := range clusterList {
			clusterData, ok := line.(sysadmK8sCluster.K8sclusterSchema)
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
	}
	clusterSelect := objectsUI.SelectData{Title: "集群", SelectedId: selectedCluster, Options: clusterOptions}
	tplData["clusterSelect"] = clusterSelect

	// preparing namespace options
	var namespaceOptions []objectsUI.SelectOption
	if selectedCluster != "0" {
		namespaceOption := objectsUI.SelectOption{
			Id:       "0",
			Text:     "所有命名空间",
			ParentID: "0",
		}
		namespaceOptions = append(namespaceOptions, namespaceOption)
		var k8sclusterEntity sysadmObjects.ObjectEntity
		k8sclusterEntity = sysadmK8sCluster.New()
		clusterInfo, e := k8sclusterEntity.GetObjectInfoByID(selectedCluster)
		if e != nil {
			return e
		}
		clusterData, ok := clusterInfo.(sysadmK8sCluster.K8sclusterSchema)
		if !ok {
			return fmt.Errorf("the data is not K8S data schema")
		}
		ca := []byte(clusterData.Ca)
		cert := []byte(clusterData.Cert)
		key := []byte(clusterData.Key)
		restConf, e := sysadmK8sClient.BuildConfigFromParametes(ca, cert, key, clusterData.Apiserver, clusterData.Id, clusterData.ClusterUser, "default")
		if e != nil {
			return e
		}

		clientSet, e := sysadmK8sClient.BuildClientset(restConf)
		if e != nil {
			return e
		}
		nsList, e := clientSet.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
		if e != nil {
			return e
		}
		for _, line := range nsList.Items {
			name := line.Name
			nsOption := objectsUI.SelectOption{
				Id:       name,
				Text:     name,
				ParentID: selectedCluster,
			}
			namespaceOptions = append(namespaceOptions, nsOption)
		}
	} else {
		namespaceOption := objectsUI.SelectOption{
			Id:       "0",
			Text:     "===选择命名空间===",
			ParentID: "0",
		}
		namespaceOptions = append(namespaceOptions, namespaceOption)
	}
	nsSelect := objectsUI.SelectData{Title: "命名空间", SelectedId: selectedNamespace, Options: namespaceOptions}
	tplData["nsSelect"] = nsSelect

	return nil
}

func buildSelectData(tplData map[string]interface{}, dcList []interface{}, requestData map[string]string) error {

	selectedDC := strings.TrimSpace(requestData["dcID"])
	if selectedDC == "" {
		selectedDC = "0"
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

	// preparing cluster options
	var clusterOptions []objectsUI.SelectOption
	clusterOption := objectsUI.SelectOption{
		Id:       "0",
		Text:     "===选择集群===",
		ParentID: "0",
	}
	clusterOptions = append(clusterOptions, clusterOption)
	if selectedDC != "0" {
		var k8sclusterEntity sysadmObjects.ObjectEntity
		k8sclusterEntity = sysadmK8sCluster.New()
		conditions := make(map[string]string, 0)
		var emptyString []string
		conditions["isDeleted"] = "='0'"
		conditions["dcid"] = "='" + selectedDC + "'"
		clusterList, e := k8sclusterEntity.GetObjectList("", emptyString, emptyString, conditions, 0, 0, make(map[string]string))
		if e != nil {
			return e
		}

		for _, line := range clusterList {
			clusterData, ok := line.(sysadmK8sCluster.K8sclusterSchema)
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
	}
	clusterSelect := objectsUI.SelectData{Title: "集群", SelectedId: selectedCluster, Options: clusterOptions}
	tplData["clusterSelect"] = clusterSelect

	nsSelect := objectsUI.SelectData{}
	tplData["nsSelect"] = nsSelect

	return nil
}

func getRequestData(c *sysadmServer.Context, keys []string) (map[string]string, error) {
	requestData, e := utils.NewGetRequestData(c, keys)
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

func newObjEntity(module string) (objectEntity, error) {
	for n, v := range modulesDefined {
		if n == module {
			return v, nil
		}
	}

	return nil, fmt.Errorf("module %s  has not found", module)
}

func buildClientSetByClusterID(clusterID string) (*kubernetes.Clientset, error) {
	clusterID = strings.TrimSpace(clusterID)
	if clusterID == "" {
		return nil, fmt.Errorf("cluster ID is empty")
	}

	var k8sclusterEntity sysadmObjects.ObjectEntity
	k8sclusterEntity = sysadmK8sCluster.New()
	clusterInfo, e := k8sclusterEntity.GetObjectInfoByID(clusterID)
	if e != nil {
		return nil, e
	}
	clusterData, ok := clusterInfo.(sysadmK8sCluster.K8sclusterSchema)
	if !ok {
		return nil, fmt.Errorf("the data is not K8S data schema")
	}
	ca := []byte(clusterData.Ca)
	cert := []byte(clusterData.Cert)
	key := []byte(clusterData.Key)
	restConf, e := sysadmK8sClient.BuildConfigFromParametes(ca, cert, key, clusterData.Apiserver, clusterData.Id, clusterData.ClusterUser, "default")
	if e != nil {
		return nil, e
	}

	clientSet, e := sysadmK8sClient.BuildClientset(restConf)
	if e != nil {
		return nil, e
	}

	return clientSet, nil
}

func checkUserLoginAndGetRequestData(c *sysadmServer.Context, keys []string) (map[string]string, error) {
	requestData, e := utils.NewGetRequestData(c, keys)
	if e != nil {
		return requestData, e
	}

	userid, e := user.GetSessionValue(c, "userid", runData.sessionName)
	if e != nil || userid == nil {
		return requestData, fmt.Errorf("user has not login or not permission")
	}

	requestData["userid"] = userid.(string)

	return requestData, nil
}

func getFormDataFromMultipartForm(c *sysadmServer.Context, keys []string) (map[string]string, error) {
	data := make(map[string]string, 0)

	keys = append(keys, "dcid")
	keys = append(keys, "clusterID")
	keys = append(keys, "namespace")
	keys = append(keys, "addType")
	keys = append(keys, "objContent")

	formData, e := utils.GetMultipartData(c, keys)
	if e != nil {
		return data, e
	}

	for _, k := range keys {
		keyData := formData[k].([]string)
		value := ""
		if len(keyData) > 1 {
			for _, v := range keyData {
				if value == "" {
					value = v
				} else {
					value = value + "," + v
				}
			}
		} else {
			if len(keyData) > 0 {
				value = keyData[0]
			} else {
				value = ""
			}
		}
		data[k] = value
	}

	objContent := formData["objContent"].([]string)
	yamlContent := objContent[0]
	if data["addType"] == "1" {
		yamlFile, e := utils.GetMultipartData(c, []string{"objFile"})
		if e != nil {
			return data, e
		}
		yamlByte, e := utils.ReadUploadedFile(yamlFile["objFile"].(*multipart.FileHeader))
		if e != nil {
			return data, e
		}
		yamlContent = utils.Interface2String(yamlByte)
	}

	data["yamlContent"] = yamlContent

	return data, nil
}

func builtContainerFormData(tplData map[string]interface{}) error {
	//var formItemNames = []string{"containerName"}
	var containerData []interface{}
	lineData := objectsUI.InitLineData("containerBasicInfoLine", false, false, false)
	_ = objectsUI.AddLabelData("basicInfoLabelID", "mid", "Left", "基本信息", false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData("containerNameLine", false, false, false)
	_ = objectsUI.AddTextData("containerNameID", "containerName", "", "容器名称", "", "", "长度小于63个字母数字或-且以字母数据开开头和结尾的字符串", 30, false, false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData("containerTypeLine", false, false, false)
	var checkboxOptions []objectsUI.Option
	checkboxOptions, _ = objectsUI.AddCheckBoxOption("特权容器", "1", false, false, checkboxOptions)
	checkboxOptions, _ = objectsUI.AddCheckBoxOption("初始化容器", "2", false, false, checkboxOptions)
	_ = objectsUI.AddCheckBoxData("containerTypeID", "containerType", "容器类型", "", false, checkboxOptions, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData("containerImageLine", false, false, false)
	_ = objectsUI.AddTextData("containerImageID", "containerImage", "", "容器镜像", "", "", "类似于hb.sysadm.cn/application/nginx:1.23.1地址", 30, false, false, lineData)
	containerData = append(containerData, lineData)

	imagePullPolicyOptions := make(map[string]string, 0)
	imagePullPolicyOptions["IfNotPresent"] = "本地不存在时拉取"
	imagePullPolicyOptions["Always"] = "总是拉取"
	imagePullPolicyOptions["Never"] = "从不拉取"
	lineData = objectsUI.InitLineData("imagePullPolicySelectLine", false, false, false)
	_ = objectsUI.AddSelectData("imagePullPolicyID", "imagePullPolicy", "IfNotPresent", "", "", "镜像摘取策略", "", 1, false, false, imagePullPolicyOptions, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData("commandLine", false, false, false)
	_ = objectsUI.AddTextData("startCommandID", "startCommand", "", "命令参数", "", "", "如果填写，则替代ENTRYPOINT和CMD值.由命令和参数通过;连接起来的字符串.如/bin/mkdir;-p;/var/data/new/subdir", 30, false, false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData("environmentLine", false, false, false)
	_ = objectsUI.AddTextData("environmentID", "environment", "", "环境变量", "", "", "由KEY1,value1;KEY2,value2组成的字符串", 30, false, false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData("containerQuotaInfoLine", false, false, false)
	_ = objectsUI.AddLabelData("containerQuotaInfoLineID", "mid", "Left", "容器资源配额", false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData("CPU资源配额", false, false, false)
	_ = objectsUI.AddTextData("cpuRequestID", "cpuRequest", "", "CPU 请求", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("cpuRequestUnit", "cpuRequestUnit", "Core", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("cpuLimitID", "cpuLimit", "", "CPU 上限", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("cpuLimitUnit", "cpuLimitUnit", "Core", "", "", false, false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData("内存资源配额", false, false, false)
	_ = objectsUI.AddTextData("memRequestID", "memRequest", "", "内存 请求", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("memRequestUnit", "memRequestUnit", "Mi", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("memLimitID", "memLimit", "", "内存 上限", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("memLimitUnit", "memLimitUnit", "Mi", "", "", false, false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData("containerHealthyCheckLine", false, false, false)
	_ = objectsUI.AddLabelData("containerHealthyCheckLabelID", "mid", "Left", "容器健康探测", false, lineData)
	containerData = append(containerData, lineData)

	containerData, e := buildProbeData("startup", "启动探测", containerData)
	if e != nil {
		return e
	}

	containerData, e = buildProbeData("readiness", "就绪探测", containerData)
	if e != nil {
		return e
	}

	containerData, e = buildProbeData("liveness", "存活探测", containerData)
	if e != nil {
		return e
	}

	tplData["ContainerData"] = containerData

	return nil
}

func buildProbeData(kind, kindTitle string, containerData []interface{}) ([]interface{}, error) {
	lineData := objectsUI.InitLineData((kind + "ProbeBlock"), false, false, false)
	_ = objectsUI.AddLabelData((kind + "ProbeBlockID"), "small", "", kindTitle, false, lineData)
	containerData = append(containerData, lineData)

	var probeTypeOptions []objectsUI.Option
	probeOptionType := objectsUI.Option{Text: "不启用", Value: "0", Checked: true, Disabled: false}
	probeTypeOptions = append(probeTypeOptions, probeOptionType)
	probeOptionType = objectsUI.Option{Text: "HTTP探测", Value: "1", Checked: false, Disabled: false}
	probeTypeOptions = append(probeTypeOptions, probeOptionType)
	probeOptionType = objectsUI.Option{Text: "TCP探测", Value: "2", Checked: false, Disabled: false}
	probeTypeOptions = append(probeTypeOptions, probeOptionType)
	probeOptionType = objectsUI.Option{Text: "命令探测", Value: "3", Checked: false, Disabled: false}
	probeTypeOptions = append(probeTypeOptions, probeOptionType)
	lineData = objectsUI.InitLineData((kind + "ProbeTypeLineID"), false, false, false)
	_ = objectsUI.AddRadioData((kind + "ProbeTypeID"), (kind + "ProbeType"), "探测方法", (kind + "ProbeTypeChanged"), false, probeTypeOptions, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData((kind + "ProbeHTTPGetLineID"), true, false, true)
	startupHTTPProtocolOptions := make(map[string]string, 0)
	startupHTTPProtocolOptions["HTTP"] = "HTTP"
	startupHTTPProtocolOptions["HTTPS"] = "HTTPS"
	e := objectsUI.AddSelectData((kind + "HttpProtocolID"), (kind + "HttpProtocol"), "HTTP", "", "", "协议", "", 1, false, false, startupHTTPProtocolOptions, lineData)
	if e != nil {
		return containerData, e
	}
	_ = objectsUI.AddTextData((kind + "HttpPortID"), (kind + "HttpPort"), "", "端口", "", "", "对于HTTP协议,不填表示80端口,对于HTTPS协议不填表示443端口", 10, false, false, lineData)
	_ = objectsUI.AddTextData((kind + "HttpPathID"), (kind + "HttpPath"), "", "路径", "", "", "以/开头的HTTP(S)请求URI", 30, false, false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData((kind + "HttpHeaderLineID"), true, true, true)
	_ = objectsUI.AddTextData((kind + "HttpHeaderKey"), (kind + "HttpHeaderKey"), "", "HttpHeader Key", "", "", "", 30, false, false, lineData)
	_ = objectsUI.AddTextData((kind + "HttpHeaderValue"), (kind + "HttpHeaderValue"), "", "值", "", "", "", 30, false, false, lineData)
	_ = objectsUI.AddWordsInputData((kind + "HttpHeaderDelID"), (kind + "HttpHeaderDel[]"), "fa-trash", "#", "workloadDelHttpHeader", false, true, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData((kind + "HttpHeaderAnchorLine"), false, false, true)
	_ = objectsUI.AddWordsInputData((kind + "HttpHeaderAnchorID"), (kind + "HttpHeaderAnchor"), "添加HTTP Header", "#", (kind + "HttpHeaderAdd"), false, false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData((kind + "TcpPortLineID"), false, false, false)
	_ = objectsUI.AddTextData((kind + "TcpPortID"), (kind + "TcpPort"), "", "TCP端口", "", "", "", 30, false, false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData((kind + "CommandLineID"), false, false, false)
	_ = objectsUI.AddTextData((kind + "CommandID"), (kind + "CommandValue"), "", "Command", "", "", "用;号连接起来的命令及其所需要的参数所组成的字符串.如/bin/sh;-c;sleep;3600", 30, false, false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData((kind + "InitialDelaySecondsLine"), false, false, false)
	_ = objectsUI.AddTextData((kind + "InitialDelaySecondsID"), (kind + "InitialDelaySeconds"), "", "初始延迟(秒)", "", "", "在容器创建后，等待多少秒再进行存活状态探测", 10, false, false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData((kind + "PeriodSecondsLine"), false, false, false)
	_ = objectsUI.AddTextData((kind + "PeriodSecondsID"), (kind + "PeriodSeconds"), "", "探测周期(秒)", "", "", "探测周期，默认是10秒,最小值为1秒", 10, false, false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData((kind + "TimeoutSecondstLine"), false, false, false)
	_ = objectsUI.AddTextData((kind + "TimeoutSecondsID"), (kind + "TimeoutSeconds"), "", "超时时间(秒)", "", "", "探测超时时间", 10, false, false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData((kind + "FailureThresholdtLine"), false, false, false)
	_ = objectsUI.AddTextData((kind + "FailureThresholdID"), (kind + "FailureThreshold"), "", "不健康阀值", "", "", "表示之前是健康状态后连续检测到多少次不健康即为视为容器不健康，最小值1，默认值为3", 10, false, false, lineData)
	containerData = append(containerData, lineData)

	lineData = objectsUI.InitLineData((kind + "SuccessThresholdtLine"), false, true, false)
	_ = objectsUI.AddTextData((kind + "SuccessThresholdID"), (kind + "SuccessThreshold[]"), "", "健康阀值", "", "", "表示容器不健康后,连续探测到多少次健康就表示容器已下午健康状态", 10, false, false, lineData)
	containerData = append(containerData, lineData)

	return containerData, nil
}

func buildStorageFormData(tplData map[string]interface{}) error {
	var storageData []interface{}

	lineData := objectsUI.InitLineData("volumeDefineLineID", false, false, false)
	_ = objectsUI.AddLabelData("volumeDefineID", "mid", "Left", "数据卷定义", false, lineData)
	storageData = append(storageData, lineData)

	lineData = objectsUI.InitLineData("volumeLineID", false, false, false)
	_ = objectsUI.AddTextData("volumeNameID", "volumeName", "", "数据卷名称", "", "", "长度小于63个字母数字或-且以字母数据开开头和结尾的字符串", 30, false, false, lineData)
	storageData = append(storageData, lineData)

	var volumeTypeOptions []objectsUI.Option
	volumeOptionType := objectsUI.Option{Text: "不挂载数据卷", Value: "0", Checked: true, Disabled: false}
	volumeTypeOptions = append(volumeTypeOptions, volumeOptionType)
	volumeOptionType = objectsUI.Option{Text: "持久数据卷", Value: "1", Checked: false, Disabled: false}
	volumeTypeOptions = append(volumeTypeOptions, volumeOptionType)
	volumeOptionType = objectsUI.Option{Text: "临时目录(EmptyDir)", Value: "2", Checked: false, Disabled: false}
	volumeTypeOptions = append(volumeTypeOptions, volumeOptionType)
	volumeOptionType = objectsUI.Option{Text: "HostPath", Value: "3", Checked: false, Disabled: false}
	volumeTypeOptions = append(volumeTypeOptions, volumeOptionType)
	volumeOptionType = objectsUI.Option{Text: "配置字典(ConfigMap)", Value: "4", Checked: false, Disabled: false}
	volumeTypeOptions = append(volumeTypeOptions, volumeOptionType)
	volumeOptionType = objectsUI.Option{Text: "密文(Secret)", Value: "5", Checked: false, Disabled: false}
	volumeTypeOptions = append(volumeTypeOptions, volumeOptionType)
	lineData = objectsUI.InitLineData("volumeTypeLineID", false, false, false)
	_ = objectsUI.AddRadioData("volumeTypeID", "volumeType", "数据卷类型", "volumeTypeChanged", false, volumeTypeOptions, lineData)
	storageData = append(storageData, lineData)

	pvcOptions := make(map[string]string, 0)
	lineData = objectsUI.InitLineData("pvcSelectLineID", false, false, true)
	_ = objectsUI.AddSelectData("pvcSelectID", "pvcSelect", "", "", "", "选择持久化数据卷声明(PVC)", "", 1, false, false, pvcOptions, lineData)
	storageData = append(storageData, lineData)

	lineData = objectsUI.InitLineData("emptyDirtLineID", false, false, true)
	_ = objectsUI.AddLabelData("emptyDirtID", "small", "", "临时目录没有参数需要配置", false, lineData)
	storageData = append(storageData, lineData)

	lineData = objectsUI.InitLineData("hostPathLineID", true, false, true)
	_ = objectsUI.AddTextData("hostPathID", "hostPath", "", "主机上路径", "", "", "HostPath类型卷将应用可能调度的节点上本字段所指定的目录或文件挂载到容器内", 30, false, false, lineData)
	storageData = append(storageData, lineData)

	hostPathTypeOptions := make(map[string]string, 0)
	hostPathTypeOptions["DirectoryOrCreate"] = "DirectoryOrCreate"
	hostPathTypeOptions["Directory"] = "Directory"
	hostPathTypeOptions["FileOrCreate"] = "FileOrCreate"
	hostPathTypeOptions["File"] = "File"
	hostPathTypeOptions["Socket"] = "Socket"
	hostPathTypeOptions["CharDevice"] = "CharDevice"
	hostPathTypeOptions["BlockDevice"] = "BlockDevice"
	lineData = objectsUI.InitLineData("hostPathTypeSelectLineID", false, true, false)
	_ = objectsUI.AddSelectData("hostPathTypeID", "hostPathType", "DirectoryOrCreate", "", "", "类型", "", 1, false, false, hostPathTypeOptions, lineData)
	storageData = append(storageData, lineData)

	cmOptions := make(map[string]string, 0)
	lineData = objectsUI.InitLineData("cmSelectLineID", false, false, true)
	_ = objectsUI.AddSelectData("cmSelectID", "cmSelect", "", "", "", "选择ConfigMap", "", 1, false, false, cmOptions, lineData)
	storageData = append(storageData, lineData)

	secretOptions := make(map[string]string, 0)
	lineData = objectsUI.InitLineData("secretSelectLineID", false, false, true)
	_ = objectsUI.AddSelectData("secretSelectID", "secretSelect", "", "", "", "选择Secret", "", 1, false, false, secretOptions, lineData)
	storageData = append(storageData, lineData)

	lineData = objectsUI.InitLineData("volumeMountLineID", false, false, true)
	_ = objectsUI.AddLabelData("volumeMountID", "Mid", "Left", "卷挂载配置", false, lineData)
	storageData = append(storageData, lineData)

	containerOptions := make(map[string]string, 0)
	lineData = objectsUI.InitLineData("containerSelectLineID", true, false, true)
	_ = objectsUI.AddSelectData("containerSelectID", "containerSelect", "", "", "", "选择容器", "", 1, false, false, containerOptions, lineData)
	storageData = append(storageData, lineData)

	lineData = objectsUI.InitLineData("volumeMountPathLineID", false, false, false)
	_ = objectsUI.AddTextData("volumeMountPathID", "volumeMountPath", "", "容器内挂载路径", "", "", "卷在容器内的挂载路径，为必填项", 30, false, false, lineData)
	storageData = append(storageData, lineData)

	lineData = objectsUI.InitLineData("volumeMountSubPathLineID", false, true, false)
	_ = objectsUI.AddTextData("volumeMountSubPathID", "volumeMountSubPath", "", "挂载路径SubPath", "", "", "挂载路径SubPath，为选填项", 30, false, false, lineData)
	storageData = append(storageData, lineData)

	tplData["StorageData"] = storageData
	return nil
}

func buildPodTemplateSpecApplyConfig(formData map[string]interface{}, matchLabels, annotations map[string]string) (*corev1.PodTemplateSpecApplyConfiguration, error) {
	// 准备Pod Spec里的Volumes数据
	volumes, volumeMounts, e := buildPodVolumes(formData)
	if e != nil {
		return nil, e
	}

	initContainers, containers, e := buildContainers(formData, volumeMounts)
	if e != nil {
		return nil, e
	}

	podSpecApplyConfiguration := corev1.PodSpecApplyConfiguration{Volumes: volumes, InitContainers: initContainers, Containers: containers}

	podTemplageSpecApplyConfiguration := &corev1.PodTemplateSpecApplyConfiguration{Spec: &podSpecApplyConfiguration}
	podTemplageSpecApplyConfiguration = podTemplageSpecApplyConfiguration.WithLabels(matchLabels)
	podTemplageSpecApplyConfiguration = podTemplageSpecApplyConfiguration.WithAnnotations(annotations)

	return podTemplageSpecApplyConfiguration, nil
}

func buildPodVolumes(formData map[string]interface{}) ([]corev1.VolumeApplyConfiguration, []VolumeMount, error) {
	var volumes []corev1.VolumeApplyConfiguration
	var volumeMounts []VolumeMount
	storageMountData, ok := formData["storageMountData[]"].([]string)
	if !ok {
		return volumes, volumeMounts, fmt.Errorf("there is an error occurred when getting storage data")
	}

	volumeNameSlice := []string{}
	for _, data := range storageMountData {
		dataJson, e := base64.StdEncoding.DecodeString(data)
		if e != nil {
			return volumes, volumeMounts, e
		}
		volumeMountData := VolumeMount{}
		e = json.Unmarshal(dataJson, &volumeMountData)
		if e != nil {
			return volumes, volumeMounts, e
		}
		volumeMounts = append(volumeMounts, volumeMountData)
		var hostPath *corev1.HostPathVolumeSourceApplyConfiguration = nil
		var emptyDir *corev1.EmptyDirVolumeSourceApplyConfiguration = nil
		var secret *corev1.SecretVolumeSourceApplyConfiguration = nil
		var pvc *corev1.PersistentVolumeClaimVolumeSourceApplyConfiguration = nil
		var configMap *corev1.ConfigMapVolumeSourceApplyConfiguration = nil
		switch volumeMountData.BasicInfo.VolumeType {
		case "1":
			pvc = &corev1.PersistentVolumeClaimVolumeSourceApplyConfiguration{}
			pvcName := volumeMountData.PvcData.Name
			readOnly := false
			pvc.ClaimName = &pvcName
			pvc.ReadOnly = &readOnly
		case "2":
			emptyDir = &corev1.EmptyDirVolumeSourceApplyConfiguration{}
		case "3":
			hostPath = &corev1.HostPathVolumeSourceApplyConfiguration{}
			path := volumeMountData.HostPathData.HostPath
			hostPathType := (apiCoreV1.HostPathType)(volumeMountData.HostPathData.HostPathType)
			hostPath.Path = &path
			hostPath.Type = &hostPathType
		case "4":
			configMap = &corev1.ConfigMapVolumeSourceApplyConfiguration{}
			name := volumeMountData.CmName
			configMap.Name = &name
		case "5":
			secret = &corev1.SecretVolumeSourceApplyConfiguration{}
			name := volumeMountData.SecretName
			secret.SecretName = &name
		default:
			return volumes, volumeMounts, fmt.Errorf("volume type %s is not valid", volumeMountData.BasicInfo.VolumeType)
		}

		volumeName := volumeMountData.BasicInfo.VolumeName
		volumeApplyConfiguration := corev1.VolumeApplyConfiguration{Name: &volumeName}
		volumeApplyConfiguration.HostPath = hostPath
		volumeApplyConfiguration.EmptyDir = emptyDir
		volumeApplyConfiguration.Secret = secret
		volumeApplyConfiguration.PersistentVolumeClaim = pvc
		volumeApplyConfiguration.ConfigMap = configMap

		volumeNameSlice = append(volumeNameSlice, volumeName)
		volumes = append(volumes, volumeApplyConfiguration)
	}

	duplicateStr := utils.IsDuplicateElementInSlice(volumeNameSlice, true)
	if duplicateStr != "" {
		return volumes, volumeMounts, fmt.Errorf("the volume name %s is duplicated", duplicateStr)
	}

	return volumes, volumeMounts, nil
}

func buildContainers(formData map[string]interface{}, volumeMounts []VolumeMount) ([]corev1.ContainerApplyConfiguration, []corev1.ContainerApplyConfiguration, error) {
	var initContainers []corev1.ContainerApplyConfiguration
	var containers []corev1.ContainerApplyConfiguration

	containerFormData, ok := formData["containerData[]"].([]string)
	if !ok {
		return initContainers, containers, fmt.Errorf("get data from client error")
	}

	var containerNameSlice []string
	for _, data := range containerFormData {
		dataJson, e := base64.StdEncoding.DecodeString(data)
		if e != nil {
			return initContainers, containers, e
		}
		containerData := ContainerData{}
		e = json.Unmarshal(dataJson, &containerData)
		if e != nil {
			return initContainers, containers, e
		}
		if utils.FoundStrInSlice(containerNameSlice, containerData.BasicInfo.ContainerName, true) {
			return initContainers, containers, fmt.Errorf("name %s of container is duplicated", containerData.BasicInfo.ContainerName)
		}
		containerNameSlice = append(containerNameSlice, containerData.BasicInfo.ContainerName)
		containerApplyConfiguration, e := buildContainerApplyConfiguraion(containerData, volumeMounts)
		if e != nil {
			return initContainers, containers, e
		}

		containerType := containerData.BasicInfo.ContainerType
		containerTypeSlice := strings.Split(containerType, ",")
		initContainer := false
		for _, v := range containerTypeSlice {
			if v == "2" {
				initContainer = true
				break
			}
		}
		if initContainer {
			initContainers = append(initContainers, containerApplyConfiguration)
		} else {
			containers = append(containers, containerApplyConfiguration)
		}
	}

	return initContainers, containers, nil
}

func buildContainerApplyConfiguraion(data ContainerData, volumeMounts []VolumeMount) (corev1.ContainerApplyConfiguration, error) {
	containerData := corev1.ContainerApplyConfiguration{}
	name := data.BasicInfo.ContainerName
	containerData.Name = &name
	image := data.BasicInfo.ContainerImage
	containerData.Image = &image
	formCommands := data.BasicInfo.StartCommand
	if formCommands != "" {
		formCommandSlice := strings.Split(formCommands, ";")
		var command []string
		var args []string
		for i, v := range formCommandSlice {
			if i == 0 {
				command = append(command, v)
			} else {
				args = append(args, v)
			}
		}
		containerData.Command = command
		containerData.Args = args
	}

	formEnvs := data.BasicInfo.Environment
	if formEnvs != "" {
		var envs []corev1.EnvVarApplyConfiguration
		formEnvsSlice := strings.Split(formEnvs, ";")
		for _, v := range formEnvsSlice {
			vSlice := strings.Split(v, ",")
			if len(vSlice) < 2 {
				continue
			}
			name := vSlice[0]
			value := vSlice[1]
			env := corev1.EnvVarApplyConfiguration{Name: &name, Value: &value}
			envs = append(envs, env)
		}
		containerData.Env = envs
	}

	resourceQuota, e := buildContainerQuota(data.Quota)
	if e != nil {
		return containerData, e
	}
	containerData.Resources = resourceQuota

	containerVolumeMounts, e := buildContainerVolumeMounts(name, volumeMounts)
	if e != nil {
		return containerData, e
	}
	containerData.VolumeMounts = containerVolumeMounts

	livenessProbe, e := buildContainerProbe(data.Liveness)
	if e != nil {
		return containerData, e
	}
	containerData.LivenessProbe = livenessProbe

	readinessProbe, e := buildContainerProbe(data.Readiness)
	if e != nil {
		return containerData, e
	}
	containerData.ReadinessProbe = readinessProbe

	startupProbe, e := buildContainerProbe(data.Startup)
	if e != nil {
		return containerData, e
	}
	containerData.StartupProbe = startupProbe

	imagePullPolicy := apiCoreV1.PullIfNotPresent
	if strings.ToUpper(strings.TrimSpace(data.BasicInfo.ImagePullPolicy)) == strings.ToUpper(string(apiCoreV1.PullAlways)) {
		imagePullPolicy = apiCoreV1.PullAlways
	}
	if strings.ToUpper(strings.TrimSpace(data.BasicInfo.ImagePullPolicy)) == strings.ToUpper(string(apiCoreV1.PullNever)) {
		imagePullPolicy = apiCoreV1.PullNever
	}

	containerData.ImagePullPolicy = &imagePullPolicy
	containerType := data.BasicInfo.ContainerType
	privileged := false
	containerTypeSlice := strings.Split(containerType, ",")
	for _, t := range containerTypeSlice {
		if t == "1" {
			privileged = true
			break
		}
	}

	if privileged {
		securityContext := corev1.SecurityContextApplyConfiguration{Privileged: &privileged}
		containerData.SecurityContext = &securityContext
	}

	return containerData, nil
}

func buildContainerQuota(data ContainerQuota) (*corev1.ResourceRequirementsApplyConfiguration, error) {
	var resourceRequirementsApplyConfiguration *corev1.ResourceRequirementsApplyConfiguration = nil

	cpuRequest := data.CpuRequest
	memRequest := data.MemRequest
	requestList := make(apiCoreV1.ResourceList, 0)
	if cpuRequest != "" {
		requestQuantity, e := resourceapi.ParseQuantity(cpuRequest)
		if e != nil {
			return resourceRequirementsApplyConfiguration, e
		}
		requestList[apiCoreV1.ResourceCPU] = requestQuantity
	}

	if memRequest != "" {
		if !strings.HasSuffix(memRequest, "Mi") {
			memRequest = memRequest + "Mi"
		}
		requestQuantity, e := resourceapi.ParseQuantity(memRequest)
		if e != nil {
			return resourceRequirementsApplyConfiguration, e
		}
		requestList[apiCoreV1.ResourceMemory] = requestQuantity
	}
	if len(requestList) > 0 {
		if resourceRequirementsApplyConfiguration == nil {
			resourceRequirementsApplyConfiguration = &corev1.ResourceRequirementsApplyConfiguration{}
		}
		resourceRequirementsApplyConfiguration.Requests = &requestList
	}

	limitList := make(apiCoreV1.ResourceList, 0)
	cpuLimit := data.CpuLimit
	memLimit := data.MemLimit
	if cpuLimit != "" {
		limitQuantity, e := resourceapi.ParseQuantity(cpuLimit)
		if e != nil {
			return resourceRequirementsApplyConfiguration, e
		}
		limitList[apiCoreV1.ResourceCPU] = limitQuantity
	}

	if memLimit != "" {
		if !strings.HasSuffix(memLimit, "Mi") {
			memLimit = memLimit + "Mi"
		}
		limitQuantity, e := resourceapi.ParseQuantity(memLimit)
		if e != nil {
			return resourceRequirementsApplyConfiguration, e
		}
		limitList[apiCoreV1.ResourceMemory] = limitQuantity
	}

	if len(limitList) > 0 {
		if resourceRequirementsApplyConfiguration == nil {
			resourceRequirementsApplyConfiguration = &corev1.ResourceRequirementsApplyConfiguration{}
		}
		resourceRequirementsApplyConfiguration.Limits = &limitList
	}

	return resourceRequirementsApplyConfiguration, nil
}

func buildContainerVolumeMounts(containerName string, volumeMounts []VolumeMount) ([]corev1.VolumeMountApplyConfiguration, error) {
	var containerVolumeMounts []corev1.VolumeMountApplyConfiguration
	var haveMountVolumes []string

	for _, v := range volumeMounts {
		container := v.ContainerData.Name
		volumeName := v.BasicInfo.VolumeName
		mountPath := v.ContainerData.MountPath
		mountSubPath := v.ContainerData.SubPath
		if container == containerName {
			if utils.FoundStrInSlice(haveMountVolumes, volumeName, true) {
				return containerVolumeMounts, fmt.Errorf("volume %s has mounted in container %s", volumeName, containerName)
			}
			haveMountVolumes = append(haveMountVolumes, volumeName)
			volumeMount := corev1.VolumeMountApplyConfiguration{Name: &volumeName, MountPath: &mountPath, SubPath: &mountSubPath}
			containerVolumeMounts = append(containerVolumeMounts, volumeMount)
		}
	}

	return containerVolumeMounts, nil
}

func buildContainerProbe(data ContainerProbe) (*corev1.ProbeApplyConfiguration, error) {
	var ret *corev1.ProbeApplyConfiguration = nil
	probeHandler := corev1.ProbeHandlerApplyConfiguration{}
	switch data.ProbeType {
	case "0":
		return ret, nil
	case "1":
		path := data.HttpPath

		host := "127.0.0.1"
		schema := apiCoreV1.URISchemeHTTP
		if strings.ToUpper(strings.TrimSpace(data.HttpProtocol)) == "HTTPS" {
			schema = apiCoreV1.URISchemeHTTPS
		}
		if path == "" {
			path = "/"
		}
		port := data.HttpPort
		if port == "" {
			if schema == apiCoreV1.URISchemeHTTP {
				port = "80"
			} else {
				port = "443"
			}
		}
		portInt, e := strconv.Atoi(port)
		if e != nil {
			return ret, fmt.Errorf("HTTP Port %s is not valid", path)
		}
		portInt32 := int32(portInt)
		portIntOrString := intstr.IntOrString{Type: intstr.Int, IntVal: portInt32}

		headerKeys := data.HttpHeaderKey
		headerValues := data.HttpHeaderValue
		if len(headerKeys) != len(headerValues) {
			return ret, fmt.Errorf("the count of HTTP Header keys is not equal to the count of HTTP Header values")
		}

		var httpHeaders []corev1.HTTPHeaderApplyConfiguration = nil
		for i, k := range headerKeys {
			v := headerValues[i]
			httpHeader := corev1.HTTPHeaderApplyConfiguration{Name: &k, Value: &v}
			httpHeaders = append(httpHeaders, httpHeader)
		}
		httpGet := &corev1.HTTPGetActionApplyConfiguration{Path: &path, Port: &portIntOrString, Host: &host, Scheme: &schema}
		if httpHeaders != nil {
			for _, h := range httpHeaders {
				httpGet = httpGet.WithHTTPHeaders(&h)
			}
		}
		probeHandler.HTTPGet = httpGet
	case "2":
		tcpPortInt, e := strconv.Atoi(data.TcpPort)
		if e != nil {
			return ret, fmt.Errorf("the tcp port %s is not valid", data.TcpPort)
		}
		tcpPortInt32 := int32(tcpPortInt)
		portIntOrString := intstr.IntOrString{Type: intstr.Int, IntVal: tcpPortInt32}
		host := "127.0.0.1"
		tcpSocketActionApplyConfiguration := corev1.TCPSocketActionApplyConfiguration{Port: &portIntOrString, Host: &host}
		probeHandler.TCPSocket = &tcpSocketActionApplyConfiguration
	case "3":
		commandStr := data.Command
		command := corev1.ExecActionApplyConfiguration{}
		var commands []string
		commandSlice := strings.Split(commandStr, ";")
		for _, v := range commandSlice {
			commands = append(commands, v)
		}
		command.Command = commands
		probeHandler.Exec = &command
	default:
		return ret, fmt.Errorf("the probe type %s is not valid", data.ProbeType)
	}

	ret = &corev1.ProbeApplyConfiguration{}
	ret.Exec = probeHandler.Exec
	ret.HTTPGet = probeHandler.HTTPGet
	ret.TCPSocket = probeHandler.TCPSocket

	if data.InitialDelaySeconds != "" {
		initialDelaySecondsInt, e := strconv.Atoi(data.InitialDelaySeconds)
		if e != nil {
			return ret, fmt.Errorf("InitialDelaySeconds %s is not valid", data.InitialDelaySeconds)
		}
		initialDelaySecondsInt32 := int32(initialDelaySecondsInt)
		ret.InitialDelaySeconds = &initialDelaySecondsInt32
	}

	if data.TimeoutSeconds != "" {
		timeoutSecondsInt, e := strconv.Atoi(data.TimeoutSeconds)
		if e != nil {
			return ret, fmt.Errorf("TimeoutSeconds %s is not valid", data.TimeoutSeconds)
		}
		timeoutSecondsInt32 := int32(timeoutSecondsInt)
		ret.TimeoutSeconds = &timeoutSecondsInt32
	}

	if data.PeriodSeconds != "" {
		periodSecondsInt, e := strconv.Atoi(data.PeriodSeconds)
		if e != nil {
			return ret, fmt.Errorf("PeriodSeconds %s is not valid", data.PeriodSeconds)
		}
		periodSecondsInt32 := int32(periodSecondsInt)
		ret.PeriodSeconds = &periodSecondsInt32
	}

	if data.SuccessThreshold != "" {
		successThresholdInt, e := strconv.Atoi(data.SuccessThreshold)
		if e != nil {
			return ret, fmt.Errorf("SuccessThreshold %s is not valid", data.SuccessThreshold)
		}
		successThresholdInt32 := int32(successThresholdInt)
		ret.SuccessThreshold = &successThresholdInt32
	}

	if data.FailureThreshold != "" {
		failureThresholdInt, e := strconv.Atoi(data.FailureThreshold)
		if e != nil {
			return ret, fmt.Errorf("FailureThreshold %s is not valid", data.FailureThreshold)
		}
		failureThresholdInt32 := int32(failureThresholdInt)
		ret.FailureThreshold = &failureThresholdInt32
	}

	return ret, nil
}

func sortWorkloadData(data []interface{}, direction, orderField string, allOrderFields map[string]objectsUI.SortBy) {
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

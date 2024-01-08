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
	applyconfigAppv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	corev1 "k8s.io/client-go/applyconfigurations/core/v1"
	applyconfigMetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"strconv"
	"strings"
	"sysadm/k8sclient"
	"sysadm/objectsUI"
	"sysadm/utils"
)

func (d *deployment) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": nil, "TD2": nil}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "创建时间"}
	additionalJs := []string{}
	additionalCss := []string{}
	templateFile := "addWorkload.html"

	d.mainModuleName = "工作负载"
	d.moduleName = "Deployment"
	d.allPopMenuItems = allPopMenuItems
	d.allListItems = allListItems
	d.addButtonTile = ""
	d.isSearchForm = "no"
	d.allOrderFields = allOrderFields
	d.defaultOrderField = "TD1"
	d.defaultOrderDirection = "1"
	d.namespaced = true
	d.moduleID = "deployment"
	d.additionalJs = additionalJs
	d.additionalCss = additionalCss
	d.templateFile = templateFile

}

func (d *deployment) getMainModuleName() string {
	return d.mainModuleName
}

func (d *deployment) getModuleName() string {
	return d.moduleName
}

func (d *deployment) getAddButtonTitle() string {
	return d.addButtonTile
}

func (d *deployment) getIsSearchForm() string {
	return d.isSearchForm
}

func (d *deployment) getAllPopMenuItems() []string {
	return d.allPopMenuItems
}

func (d *deployment) getAllListItems() map[string]string {
	return d.allListItems
}

func (d *deployment) getDefaultOrderField() string {
	return d.defaultOrderField
}

func (d *deployment) getDefaultOrderDirection() string {
	return d.defaultOrderDirection
}

func (d *deployment) getAllorderFields() map[string]objectsUI.SortBy {
	return d.allOrderFields
}

func (d *deployment) getNamespaced() bool {
	return d.namespaced
}

// for ingressclass
func (d *deployment) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	// TODO

	return 0, dataList, nil
}

func (d *deployment) getModuleID() string {
	return d.moduleID
}

func (d *deployment) buildAddFormData(tplData map[string]interface{}) error {
	tplData["thirdCategory"] = "创建Deployment"
	formData, e := objectsUI.InitFormData("addDeployment", "addDeployment", "POST", "_self", "yes", "addWorkload", "")
	if e != nil {
		return e
	}
	tplData["formData"] = formData

	e = buildBasiceFormData(tplData)
	if e != nil {
		return e
	}

	e = builtContainerFormData(tplData)
	if e != nil {
		return e
	}

	e = buildStorageFormData(tplData)
	if e != nil {
		return e
	}

	e = buildMoreDataFormData(tplData)
	if e != nil {
		return e
	}
	e = buildServiceDataFormData(tplData)
	if e != nil {
		return e
	}

	return nil
}

func (d *deployment) getAdditionalJs() []string {
	return d.additionalJs
}
func (d *deployment) getAdditionalCss() []string {
	return d.additionalCss
}

func (d *deployment) addNewResource(c *sysadmServer.Context, module string) error {
	requestKeys := []string{"dcid", "clusterID", "namespace", "addType", "nsSelected", "name", "replics", "labelKey[]", "labelValue[]", "annotationKey[]", "annotationValue[]"}
	requestKeys = append(requestKeys, "selectorKey[]")
	requestKeys = append(requestKeys, "selectorValue[]")
	requestKeys = append(requestKeys, "containerData[]")
	requestKeys = append(requestKeys, "storageMountData[]")
	formData, e := utils.GetMultipartData(c, requestKeys)
	if e != nil {
		return e
	}

	ns := formData["nsSelected"].([]string)
	name := formData["name"].([]string)
	deployApplyConfig := applyconfigAppv1.Deployment(name[0], ns[0])

	// 配置labels
	labelKeys := formData["labelKey[]"].([]string)
	labelValue := formData["labelValue[]"].([]string)
	if len(labelKeys) != len(labelValue) {
		return fmt.Errorf("label's key is not equal to label's value")
	}

	labels := make(map[string]string, 0)
	if len(labelKeys) > 0 {
		for i, k := range labelKeys {
			value := labelValue[i]
			labels[k] = value
		}
	} else {
		labels[defaultLabelKey] = name[0]
	}
	for k, v := range extraLabels {
		labels[k] = v
	}
	deployApplyConfig = deployApplyConfig.WithLabels(labels)

	// 配置注解
	annotationKey := formData["annotationKey[]"].([]string)
	annotationValue := formData["annotationValue[]"].([]string)
	if len(annotationKey) != len(annotationValue) {
		return fmt.Errorf("annotation's key is not equal to annotation's value")
	}
	annotations := make(map[string]string, 0)
	for i, k := range annotationKey {
		value := annotationValue[i]
		annotations[k] = value
	}
	deployApplyConfig = deployApplyConfig.WithAnnotations(annotations)

	// 准备副本数
	deploySpecApplyConfig := applyconfigAppv1.DeploymentSpecApplyConfiguration{}
	replicsSlice := formData["replics"].([]string)
	replicsStr := replicsSlice[0]
	replicsInt, e := strconv.Atoi(replicsStr)
	if e != nil {
		return e
	}
	replicsInt32 := int32(replicsInt)
	deploySpecApplyConfig.Replicas = &replicsInt32

	// 配置selector
	selectorKeys := formData["selectorKey[]"].([]string)
	selectorValues := formData["selectorValue[]"].([]string)
	if len(selectorKeys) != len(selectorValues) {
		return fmt.Errorf("selector's key is not equal to selector's value")
	}
	matchLabels := make(map[string]string, 0)
	if len(selectorKeys) > 0 {
		for i, k := range selectorKeys {
			matchLabels[k] = selectorValues[i]
		}
	} else {
		matchLabels = labels
	}
	labelSelector := applyconfigMetav1.LabelSelectorApplyConfiguration{MatchLabels: matchLabels}
	deploySpecApplyConfig.Selector = &labelSelector

	podTemplateSpecApplyConfiguration, e := buildPodTemplateSpecApplyConfig(formData, matchLabels, annotations)
	if e != nil {
		return e
	}
	deploySpecApplyConfig.Template = podTemplateSpecApplyConfiguration
	deployApplyConfig = deployApplyConfig.WithSpec(&deploySpecApplyConfig)

	clusterIDSlice := formData["clusterID"].([]string)
	clusterID := clusterIDSlice[0]
	clientSet, e := buildClientSetByClusterID(clusterID)
	if e != nil {
		return e
	}
	applyOption := metav1.ApplyOptions{
		Force:        true,
		FieldManager: k8sclient.FieldManager,
	}
	_, e = clientSet.AppsV1().Deployments(ns[0]).Apply(context.Background(), deployApplyConfig, applyOption)

	return e

}

func (d *deployment) delResource(s *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (d *deployment) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

func (d *deployment) getTemplateFile(action string) string {

	return d.templateFile
}

func buildBasiceFormData(tplData map[string]interface{}) error {
	clusterID := tplData["clusterID"].(string)
	if clusterID == "" || clusterID == "0" {
		return fmt.Errorf("cluster must be specified when add a new deployment")
	}
	clientSet, e := buildClientSetByClusterID(clusterID)
	if e != nil {
		return e
	}
	nsObjectList, e := clientSet.CoreV1().Namespaces().List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return e
	}
	nsExistList := []string{}
	for _, item := range nsObjectList.Items {
		found := false
		for _, n := range denyDeployWokloadNSList {
			if strings.TrimSpace(strings.ToLower(item.Name)) == n {
				found = true
			}
		}

		if found {
			continue
		}
		nsExistList = append(nsExistList, item.Name)

	}

	// 准备命名空间下拉菜单option数据
	nsOptions := make(map[string]string, 0)
	defaultSelectedNs := ""
	for _, n := range nsExistList {
		nsOptions[n] = n
		if defaultSelectedNs == "" {
			defaultSelectedNs = n
		}
	}

	// 准备基本信息表单内容
	var basicData []interface{}
	lineData := objectsUI.InitLineData("namespaceSelectLine", false, false, false)
	e = objectsUI.AddSelectData("nsSelectedID", "nsSelected", defaultSelectedNs, "", "", "选择命名空间", "", 1, false, false, nsOptions, lineData)
	if e != nil {
		return e
	}
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("nameLine", false, false, false)
	_ = objectsUI.AddTextData("name", "name", "", "应用名称", "validateNewName", "addWorkloadValidateNewName", "长度小于63个字母数字或-且以字母数据开开头和结尾的字符串", 30, false, false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("replicsLine", false, false, false)
	_ = objectsUI.AddTextData("replics", "replics", "", "副本数", "", "", "大于等于0的整数", 10, false, false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("newNsLabel", true, true, false)
	_ = objectsUI.AddTextData("labelKey", "labelKey[]", "", "标签", "", "", "", 30, false, false, lineData)
	_ = objectsUI.AddWordsInputData("equal", "equal", "=", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("labelValue", "labelValue[]", "", "值", "", "", "", 30, false, false, lineData)
	_ = objectsUI.AddWordsInputData("delLabel", "delLabel", "fa-trash", "#", "workloadDelLabel", false, true, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("addlabelanchor", false, false, false)
	_ = objectsUI.AddWordsInputData("addLabel", "addLabel", "增加标签", "#", "workloadAddLabel", false, false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("newAnnotationLabel", true, true, false)
	_ = objectsUI.AddTextData("annotationKey", "annotationKey[]", "", "注解", "", "", "", 30, false, false, lineData)
	_ = objectsUI.AddWordsInputData("equal", "equal", "=", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("annotationValue", "annotationValue[]", "", "值", "", "", "", 30, false, false, lineData)
	_ = objectsUI.AddWordsInputData("annotationLabel", "annotationLabel", "fa-trash", "#", "workloadDelAnnotaion", false, true, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("annotationanchor", false, false, false)
	_ = objectsUI.AddWordsInputData("annotationLabel", "annotationLabel", "增加注解", "#", "workloadAddAnnotation", false, false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("selectorLabel", true, true, false)
	_ = objectsUI.AddTextData("selectornKey", "selectorKey[]", "", "标签选择器", "", "", "", 30, false, false, lineData)
	_ = objectsUI.AddWordsInputData("equal", "equal", "=", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("selectorValue", "selectorValue[]", "", "值", "", "", "", 30, false, false, lineData)
	_ = objectsUI.AddWordsInputData("selectorLabel", "selectorLabel", "fa-trash", "#", "workloadDelSelector", false, true, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("selectoranchor", false, false, false)
	_ = objectsUI.AddWordsInputData("selectorLabel", "selectorLabel", "添加匹配条件", "#", "workloadAddSelector", false, false, lineData)
	basicData = append(basicData, lineData)
	tplData["BasicData"] = basicData

	return nil
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

func buildMoreDataFormData(tplData map[string]interface{}) error {
	var moreData []interface{}

	// TODO

	tplData["MoreData"] = moreData

	return nil
}

func buildServiceDataFormData(tplData map[string]interface{}) error {
	var serviceData []interface{}

	// TODO

	tplData["ServiceData"] = serviceData

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

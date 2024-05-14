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
	"fmt"
	"github.com/wangyysde/sysadmServer"
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applyconfigAppv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	applyconfigMetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"strconv"
	"strings"
	"sysadm/k8sclient"
	"sysadm/objectsUI"
	"sysadm/utils"
)

func (d *deployment) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": deploymentSortByName, "TD6": deploymentSortByCreatetime}
	allPopMenuItems := []string{"Scale,scale,POST,tip", "编辑,edit,GET,page", "重启,restart,POST,tip", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "状态", "TD4": "标签", "TD5": "Pods", "TD6": "创建时间"}
	additionalJs := []string{"js/sysadmfunctions.js", "/js/workloadList.js"}
	additionalCss := []string{}
	templateFile := "addWorkload.html"

	d.mainModuleName = "工作负载"
	d.moduleName = "无状态服务"
	d.allPopMenuItems = allPopMenuItems
	d.allListItems = allListItems
	d.addButtonTile = "创建无状态服务"
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

// for deployment
func (d *deployment) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	if selectedCluster == "" || selectedCluster == "0" {
		return 0, dataList, nil
	}
	nsStr := ""
	if selectedNS != "0" {
		nsStr = selectedNS
	}

	clientSet, e := buildClientSetByClusterID(selectedCluster)
	if e != nil {
		return 0, dataList, e
	}

	deployList, e := clientSet.AppsV1().Deployments(nsStr).List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return 0, dataList, e
	}

	totalNum := len(deployList.Items)
	if totalNum < 1 {
		return 0, dataList, nil
	}

	orderfield := requestData["orderfield"]
	direction := requestData["direction"]
	if orderfield == "" {
		orderfield = d.getDefaultOrderField()
	}
	if direction == "" || (direction != "0" && direction != "1") {
		direction = d.getDefaultOrderDirection()
	}

	var deployItems []interface{}
	for _, item := range deployList.Items {
		deployItems = append(deployItems, item)
	}
	sortWorkloadData(deployItems, direction, orderfield, d.getAllorderFields())

	endCount := totalNum
	if endCount > startPos+runData.pageInfo.NumPerPage {
		endCount = startPos + runData.pageInfo.NumPerPage
	}

	for i := startPos; i < endCount; i++ {
		interfaceData := deployItems[i]
		deployData, ok := interfaceData.(appsv1.Deployment)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not Deployment schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = deployData.Name
		lineMap["TD1"] = deployData.Name
		lineMap["TD2"] = deployData.Namespace
		statusStr := "运行中"
		if deployData.Status.ReadyReplicas == 0 {
			statusStr = "未运行"
		}
		if deployData.Status.ReadyReplicas < deployData.Status.Replicas {
			statusStr = "部分运行"
		}
		lineMap["TD3"] = statusStr
		lineMap["TD4"] = objectsUI.ConvertMap2HTML(deployData.Labels)
		lineMap["TD5"] = strconv.Itoa(int(deployData.Status.ReadyReplicas)) + "/" + strconv.Itoa(int(deployData.Status.Replicas))
		lineMap["TD6"] = deployData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		popmenuitems := ""
		if int(deployData.Status.Replicas) > 0 {
			popmenuitems = "0,1,3"
		} else {
			popmenuitems = "0,1,2,3"
		}
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
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

	e = buildDeployBasiceFormData(tplData)
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
	switch action {
	case "list":
		return deploymentTemplateFiles["list"]
	case "addform":
		return deploymentTemplateFiles["addform"]
	default:
		return ""
	}

	return d.templateFile
}

func buildDeployBasiceFormData(tplData map[string]interface{}) error {
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

func deploymentSortByName(p, q interface{}) bool {
	pData, ok := p.(appsv1.Deployment)
	if !ok {
		return false
	}
	qData, ok := q.(appsv1.Deployment)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func deploymentSortByCreatetime(p, q interface{}) bool {
	pData, ok := p.(appsv1.Deployment)
	if !ok {
		return false
	}
	qData, ok := q.(appsv1.Deployment)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}
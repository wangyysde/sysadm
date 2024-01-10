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
	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applyconfigAppv1 "k8s.io/client-go/applyconfigurations/apps/v1"
	applyconfigMetav1 "k8s.io/client-go/applyconfigurations/meta/v1"
	"sort"
	"strconv"
	"strings"
	"sysadm/k8sclient"
	"sysadm/objectsUI"
	"sysadm/utils"
)

func (s *statefulSet) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": nil, "TD6": nil}
	allPopMenuItems := []string{"Scale,scale,POST,tip", "编辑,edit,GET,page", "重启,restart,POST,tip", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "状态", "TD4": "标签", "TD5": "Pods", "TD6": "创建时间"}
	additionalJs := []string{"js/sysadmfunctions.js", "/js/workloadList.js"}
	additionalCss := []string{}
	templateFile := "workloadlist.html"

	s.mainModuleName = "工作负载"
	s.moduleName = "有状态服务"
	s.allPopMenuItems = allPopMenuItems
	s.allListItems = allListItems
	s.addButtonTile = "添加有状态服务"
	s.isSearchForm = "no"
	s.allOrderFields = allOrderFields
	s.defaultOrderField = "TD1"
	s.defaultOrderDirection = "1"
	s.namespaced = true
	s.moduleID = "statefulSet"
	s.additionalJs = additionalJs
	s.additionalCss = additionalCss
	s.templateFile = templateFile

}

func (s *statefulSet) getMainModuleName() string {
	return s.mainModuleName
}

func (s *statefulSet) getModuleName() string {
	return s.moduleName
}

func (s *statefulSet) getAddButtonTitle() string {
	return s.addButtonTile
}

func (s *statefulSet) getIsSearchForm() string {
	return s.isSearchForm
}

func (s *statefulSet) getAllPopMenuItems() []string {
	return s.allPopMenuItems
}

func (s *statefulSet) getAllListItems() map[string]string {
	return s.allListItems
}

func (s *statefulSet) getDefaultOrderField() string {
	return s.defaultOrderField
}

func (s *statefulSet) getDefaultOrderDirection() string {
	return s.defaultOrderDirection
}

func (s *statefulSet) getAllorderFields() map[string]objectsUI.SortBy {
	return s.allOrderFields
}

func (s *statefulSet) getNamespaced() bool {
	return s.namespaced
}

// for ingressclass
func (s *statefulSet) listObjectData(selectedCluster, selectedNS string,
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

	stsList, e := clientSet.AppsV1().StatefulSets(nsStr).List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return 0, dataList, e
	}

	totalNum := len(stsList.Items)
	if totalNum < 1 {
		return 0, dataList, nil
	}

	orderfield := requestData["orderfield"]
	direction := requestData["direction"]
	if orderfield == "" {
		orderfield = s.defaultOrderField
	}
	if direction == "" || (direction != "0" && direction != "1") {
		direction = s.defaultOrderDirection
	}

	var stsItems []interface{}
	for _, item := range stsList.Items {
		stsItems = append(stsItems, item)
	}

	if direction == "1" {
		if orderfield == "TD1" {
			sort.Sort(objectsUI.SortData{Data: stsItems, By: sortStsByName})
		} else {
			sort.Sort(objectsUI.SortData{Data: stsItems, By: sortStsByCreatetime})
		}
	} else {
		if orderfield == "TD1" {
			sort.Sort(sort.Reverse(objectsUI.SortData{Data: stsItems, By: sortStsByName}))
		} else {
			sort.Sort(sort.Reverse(objectsUI.SortData{Data: stsItems, By: sortStsByCreatetime}))
		}
	}

	endCount := totalNum
	if endCount > startPos+runData.pageInfo.NumPerPage {
		endCount = startPos + runData.pageInfo.NumPerPage
	}

	for i := startPos; i < endCount; i++ {
		interfaceData := stsItems[i]
		stsData, ok := interfaceData.(appsv1.StatefulSet)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not StatefulSet schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = stsData.Name
		lineMap["TD1"] = stsData.Name
		lineMap["TD2"] = stsData.Namespace
		statusStr := "运行中"
		if stsData.Status.ReadyReplicas == 0 {
			statusStr = "未运行"
		}
		if stsData.Status.ReadyReplicas < stsData.Status.Replicas {
			statusStr = "部分运行"
		}
		lineMap["TD3"] = statusStr
		lineMap["TD4"] = objectsUI.ConvertMap2HTML(stsData.Labels)
		lineMap["TD5"] = strconv.Itoa(int(stsData.Status.ReadyReplicas)) + "/" + strconv.Itoa(int(stsData.Status.Replicas))
		lineMap["TD6"] = stsData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		popmenuitems := ""
		if int(stsData.Status.Replicas) > 0 {
			popmenuitems = "0,1,3"
		} else {
			popmenuitems = "0,1,2,3"
		}
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil

}

func (s *statefulSet) getModuleID() string {
	return s.moduleID
}

func (s *statefulSet) buildAddFormData(tplData map[string]interface{}) error {
	tplData["thirdCategory"] = "创建有状态服务"
	formData, e := objectsUI.InitFormData("addStatefulSet", "addStatefulSet", "POST", "_self", "yes", "addWorkload", "")
	if e != nil {
		return e
	}
	tplData["formData"] = formData

	e = buildStatefulSetBasiceFormData(tplData)
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

func (s *statefulSet) getAdditionalJs() []string {
	return s.additionalJs
}
func (s *statefulSet) getAdditionalCss() []string {
	return s.additionalCss
}

func (s *statefulSet) addNewResource(c *sysadmServer.Context, module string) error {
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
	statefulSetApplyConfig := applyconfigAppv1.StatefulSet(name[0], ns[0])

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
		labels[defaultStatefulSetLabelKey] = name[0]
	}
	for k, v := range extraLabels {
		labels[k] = v
	}
	statefulSetApplyConfig = statefulSetApplyConfig.WithLabels(labels)

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
	statefulSetApplyConfig = statefulSetApplyConfig.WithAnnotations(annotations)

	// 准备副本数
	statefulSetSpecApplyConfig := applyconfigAppv1.StatefulSetSpecApplyConfiguration{}
	replicsSlice := formData["replics"].([]string)
	replicsStr := replicsSlice[0]
	replicsInt, e := strconv.Atoi(replicsStr)
	if e != nil {
		return e
	}
	replicsInt32 := int32(replicsInt)
	statefulSetSpecApplyConfig.Replicas = &replicsInt32

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
	statefulSetSpecApplyConfig.Selector = &labelSelector

	podTemplateSpecApplyConfiguration, e := buildPodTemplateSpecApplyConfig(formData, matchLabels, annotations)
	if e != nil {
		return e
	}
	statefulSetSpecApplyConfig.Template = podTemplateSpecApplyConfiguration
	statefulSetApplyConfig = statefulSetApplyConfig.WithSpec(&statefulSetSpecApplyConfig)

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
	_, e = clientSet.AppsV1().StatefulSets(ns[0]).Apply(context.Background(), statefulSetApplyConfig, applyOption)

	return e

}

func (s *statefulSet) delResource(ss *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (s *statefulSet) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

func (s *statefulSet) getTemplateFile(action string) string {
	switch action {
	case "list":
		return statefulSetTemplateFiles["list"]
	case "addform":
		return statefulSetTemplateFiles["addform"]
	}

	return ""
}

func buildStatefulSetBasiceFormData(tplData map[string]interface{}) error {
	clusterID := tplData["clusterID"].(string)
	if clusterID == "" || clusterID == "0" {
		return fmt.Errorf("cluster must be specified when add a new statefulSet")
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
		for _, n := range denyStatefulSetWokloadNSList {
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

func sortStsByName(p, q interface{}) bool {
	pData, ok := p.(appsv1.StatefulSet)
	if !ok {
		return false
	}
	qData, ok := q.(appsv1.StatefulSet)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortStsByCreatetime(p, q interface{}) bool {
	pData, ok := p.(appsv1.StatefulSet)
	if !ok {
		return false
	}
	qData, ok := q.(appsv1.StatefulSet)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

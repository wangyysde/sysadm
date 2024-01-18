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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applyconfigCorev1 "k8s.io/client-go/applyconfigurations/core/v1"
	"sort"
	"strconv"
	"strings"
	"sysadm/k8sclient"
	"sysadm/objectsUI"
	"sysadm/utils"
)

func (c *configmap) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortCmByName, "TD6": sortCmByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "标签", "TD4": "数据项数", "TD5": "是否可修改", "TD6": "创建时间"}
	additionalJs := []string{}
	additionalCss := []string{}
	templateFile := ""

	c.mainModuleName = "配置和存储"
	c.moduleName = "配置字典"
	c.allPopMenuItems = allPopMenuItems
	c.allListItems = allListItems
	c.addButtonTile = "创建配置字典"
	c.isSearchForm = "no"
	c.allOrderFields = allOrderFields
	c.defaultOrderField = "TD1"
	c.defaultOrderDirection = "1"
	c.namespaced = true
	c.moduleID = "configmap"
	c.additionalJs = additionalJs
	c.additionalCss = additionalCss
	c.templateFile = templateFile
}

func (c *configmap) getMainModuleName() string {
	return c.mainModuleName
}

func (c *configmap) getModuleName() string {
	return c.moduleName
}

func (c *configmap) getAddButtonTitle() string {
	return c.addButtonTile
}

func (c *configmap) getIsSearchForm() string {
	return c.isSearchForm
}

func (c *configmap) getAllPopMenuItems() []string {
	return c.allPopMenuItems
}

func (c *configmap) getAllListItems() map[string]string {
	return c.allListItems
}

func (c *configmap) getDefaultOrderField() string {
	return c.defaultOrderField
}

func (c *configmap) getDefaultOrderDirection() string {
	return c.defaultOrderDirection
}

func (c *configmap) getAllorderFields() map[string]objectsUI.SortBy {
	return c.allOrderFields
}

func (c *configmap) getNamespaced() bool {
	return c.namespaced
}

// for configmap
func (c *configmap) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	clientSet, e := createClientSet(selectedCluster)
	if clientSet == nil {
		return 0, dataList, e
	}

	nsStr, orderField, direction := checkRequestData(selectedNS, c.defaultOrderField, c.defaultOrderDirection, requestData)
	cmList, e := clientSet.CoreV1().ConfigMaps(nsStr).List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return 0, dataList, e
	}
	totalNum := len(cmList.Items)
	if totalNum < 1 {
		return 0, dataList, nil
	}
	endCount := totalNum
	if endCount > startPos+runData.pageInfo.NumPerPage {
		endCount = startPos + runData.pageInfo.NumPerPage
	}

	var cmItems []interface{}
	for _, item := range cmList.Items {
		cmItems = append(cmItems, item)
	}

	moduleAllOrderFields := c.allOrderFields
	for field, fn := range moduleAllOrderFields {
		if field == orderField {
			if direction == "1" {
				sort.Sort(objectsUI.SortData{Data: cmItems, By: fn})
			} else {
				sort.Sort(sort.Reverse(objectsUI.SortData{Data: cmItems, By: fn}))
			}

		}
	}

	for i := startPos; i < endCount; i++ {
		interfaceData := cmItems[i]
		cmIData, ok := interfaceData.(corev1.ConfigMap)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not Ingress schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = cmIData.Name
		lineMap["TD1"] = cmIData.Name
		lineMap["TD2"] = cmIData.Namespace
		lineMap["TD3"] = objectsUI.ConvertMap2HTML(cmIData.Labels)
		dataCount := len(cmIData.Data)
		binaryDataCount := len(cmIData.BinaryData)
		totalDataCount := dataCount + binaryDataCount
		lineMap["TD4"] = strconv.Itoa(totalDataCount)
		editable := "是"
		popmenuitems := "0,1"
		if cmIData.Immutable != nil && *cmIData.Immutable {
			editable = "否"
			popmenuitems = "0,1"
		}
		lineMap["TD5"] = editable
		lineMap["TD6"] = cmIData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)

		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortCmByName(p, q interface{}) bool {
	pData, ok := p.(corev1.ConfigMap)
	if !ok {
		return false
	}
	qData, ok := q.(corev1.ConfigMap)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortCmByCreatetime(p, q interface{}) bool {
	pData, ok := p.(corev1.ConfigMap)
	if !ok {
		return false
	}
	qData, ok := q.(corev1.ConfigMap)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

func (c *configmap) getModuleID() string {
	return c.moduleID
}

func (c *configmap) buildAddFormData(tplData map[string]interface{}) error {
	tplData["thirdCategory"] = "创建数据字典"
	formData, e := objectsUI.InitFormData("addConfigMap", "addConfigMap", "POST", "_self", "yes", "addWorkload", "")
	if e != nil {
		return e
	}
	tplData["formData"] = formData

	e = buildConfigMapAddFormData(tplData)
	if e != nil {
		return e
	}

	return nil
}

func (c *configmap) getAdditionalJs() []string {
	return c.additionalJs
}
func (c *configmap) getAdditionalCss() []string {
	return c.additionalCss
}

func (c *configmap) addNewResource(s *sysadmServer.Context, module string) error {
	requestKeys := []string{"dcid", "clusterID", "namespace", "addType", "nsSelected", "name", "labelKey[]", "labelValue[]", "annotationKey[]", "annotationValue[]"}
	requestKeys = append(requestKeys, "configmapkey[]", "isBinaryData[]", "configmapData[]")
	formData, e := utils.GetMultipartData(s, requestKeys)
	if e != nil {
		return e
	}
	ns := formData["nsSelected"].([]string)
	name := formData["name"].([]string)
	configmapApplyConfig := applyconfigCorev1.ConfigMap(name[0], ns[0])
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
		labels[defaultConfigMapLabelKey] = name[0]
	}
	for k, v := range extraLabels {
		labels[k] = v
	}
	configmapApplyConfig = configmapApplyConfig.WithLabels(labels)

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
	configmapApplyConfig = configmapApplyConfig.WithAnnotations(annotations)

	configmapkeys := formData["configmapkey[]"].([]string)
	isBinaryDatas := formData["isBinaryData[]"].([]string)
	configmapDatas := formData["configmapData[]"].([]string)
	if len(configmapkeys) < 1 {
		return fmt.Errorf("there is not data in the new configmap")
	}
	asciiData := make(map[string]string)
	binaryData := make(map[string][]byte)
	var kSlice []string
	for i, k := range configmapkeys {
		if utils.FoundStrInSlice(kSlice, k, true) {
			return fmt.Errorf("the k %s is duplicate", k)
		}
		kSlice = append(kSlice, k)
		isBinary := isBinaryDatas[i]
		if isBinary == "1" {
			cmDataByte := []byte(strings.TrimSpace(configmapDatas[i]))
			k = strings.TrimSpace(k)
			binaryData[k] = cmDataByte
		} else {
			strData := configmapDatas[i]
			asciiData[k] = strData
		}
	}

	configmapApplyConfig = configmapApplyConfig.WithBinaryData(binaryData)
	configmapApplyConfig = configmapApplyConfig.WithData(asciiData)

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

	_, e = clientSet.CoreV1().ConfigMaps(ns[0]).Apply(context.Background(), configmapApplyConfig, applyOption)

	return e
}

func (c *configmap) delResource(s *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (c *configmap) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

func (c *configmap) getTemplateFile(action string) string {

	switch action {
	case "list":
		return configMapTemplateFiles["list"]
	case "addform":
		return configMapTemplateFiles["addform"]
	}

	return ""
}

func buildConfigMapAddFormData(tplData map[string]interface{}) error {
	clusterID := tplData["clusterID"].(string)
	if clusterID == "" || clusterID == "0" {
		return fmt.Errorf("cluster must be specified when add a new configMap")
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
	lineData := objectsUI.InitLineData("basicInfoID", false, false, false)
	_ = objectsUI.AddLabelData("basicInfoID", "mid", "Left", "基本信息", false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("namespaceSelectLine", false, false, false)
	e = objectsUI.AddSelectData("nsSelectedID", "nsSelected", defaultSelectedNs, "", "", "选择命名空间", "", 1, false, false, nsOptions, lineData)
	if e != nil {
		return e
	}
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("nameLine", false, false, false)
	_ = objectsUI.AddTextData("name", "name", "", "字典名称", "validateNewName", "addWorkloadValidateNewName", "长度小于63个字母数字或-且以字母数据开开头和结尾的字符串", 30, false, false, lineData)
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

	lineData = objectsUI.InitLineData("configMapDataZoneID", false, false, false)
	_ = objectsUI.AddLabelData("configMapDataZoneID", "mid", "Left", "字典数据", false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("configmapkeyID", true, true, false)
	_ = objectsUI.AddTextData("configmapkeyID[]", "configmapkey[]", "", "Key", "", "", "", 20, false, false, lineData)

	dataTypeOptions := make(map[string]string, 0)
	dataTypeOptions["0"] = "文本数据"
	dataTypeOptions["1"] = "二进制数据"
	_ = objectsUI.AddSelectData("isBinaryData", "isBinaryData[]", "0", "", "", "数据类型： ", "", 1, false, false, dataTypeOptions, lineData)
	_ = objectsUI.AddTextareaData("configmapDataID[]", "configmapData[]", "", "  数据", "", "", "", 60, 5, false, false, lineData)
	_ = objectsUI.AddWordsInputData("selectorLabel[]", "selectorLabel", "fa-trash", "#", "workloadDelSelector", false, true, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("configmapkeyID", false, false, true)
	_ = objectsUI.AddWordsInputData("configmapkeyID", "configmapkey", "添加数据项", "#", "workloadAddSelectorBlock", false, false, lineData)
	basicData = append(basicData, lineData)
	tplData["BasicData"] = basicData

	return nil
}

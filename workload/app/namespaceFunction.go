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
	"github.com/wangyysde/yaml"
	corev1 "k8s.io/api/core/v1"
	resourceapi "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	appconfigv1 "k8s.io/client-go/applyconfigurations/core/v1"
	"mime/multipart"
	"sort"
	"strings"
	"sysadm/k8sclient"
	"sysadm/objectsUI"
	"sysadm/utils"
)

var mainModuleName = "集群管理"
var moduleName = "命名空间"

func (n *namespace) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortNamespaceByName, "TD2": sortNamespaceByCreatetime}
	allPopMenuItems := []string{"删除,del,POST,tip", "新增配额,addQuota,GET,page", "配额列表,listQuota,Get,page", "新增默认配额,addListRange,GET,page", "默认配额列表,addListRange,GET,page"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "创建时间"}
	additionalJs := []string{"/js/namespace.js"}
	additionalCss := []string{}
	templateFile := ""

	n.mainModuleName = mainModuleName
	n.moduleName = moduleName
	n.allPopMenuItems = allPopMenuItems
	n.allListItems = allListItems
	n.addButtonTile = ""
	n.isSearchForm = "no"
	n.allOrderFields = allOrderFields
	n.defaultOrderField = "TD1"
	n.defaultOrderDirection = "1"
	n.namespaced = false
	n.moduleID = "namespace"
	n.additionalJs = additionalJs
	n.additionalCss = additionalCss
	n.templateFile = templateFile
}

func (n *namespace) getMainModuleName() string {
	return n.mainModuleName
}

func (n *namespace) getModuleName() string {
	return n.moduleName
}

func (n *namespace) getAddButtonTitle() string {
	return n.addButtonTile
}

func (n *namespace) getIsSearchForm() string {
	return n.isSearchForm
}

func (n *namespace) getAllPopMenuItems() []string {
	return n.allPopMenuItems
}

func (n *namespace) getAllListItems() map[string]string {
	return n.allListItems
}

func (n *namespace) getDefaultOrderField() string {
	return n.defaultOrderField
}

func (n *namespace) getDefaultOrderDirection() string {
	return n.defaultOrderDirection
}

func (n *namespace) getAllorderFields() map[string]objectsUI.SortBy {
	return n.allOrderFields
}

func (n *namespace) getNamespaced() bool {
	return n.namespaced
}

// for ingressclass
func (n *namespace) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	//TODO

	return 0, dataList, nil
}

func sortNamespaceByName(p, q interface{}) bool {
	// TODO

	return true
}

func sortNamespaceByCreatetime(p, q interface{}) bool {
	// TODO

	return true
}

func (n *namespace) getModuleID() string {
	return n.moduleID
}

func (n *namespace) buildAddFormData(tplData map[string]interface{}) error {
	formData, e := objectsUI.InitFormData("addNamespace", "addNamespace", "POST", "_self", "yes", "addNewNamespace", "")
	if e != nil {
		return e
	}

	lineData := objectsUI.InitLineData("newns", false, false, false)
	e = objectsUI.AddTextData("newns", "newns", "", "命名空间名称", "validateNewName", "", "长度小于63个字母数字或-且以字母数据开开头和结尾的字符串", 30, false, false, lineData)
	if e != nil {
		return e
	}
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("newNsLabel", true, true, false)
	e = objectsUI.AddTextData("labelKey", "labelKey[]", "", "标签", "", "", "", 30, false, false, lineData)
	if e != nil {
		return e
	}

	e = objectsUI.AddWordsInputData("equal", "equal", "=", "", "", false, false, lineData)
	if e != nil {
		return e
	}

	e = objectsUI.AddTextData("labelValue", "labelValue[]", "", "值", "", "", "", 30, false, false, lineData)
	if e != nil {
		return e
	}

	_ = objectsUI.AddWordsInputData("delLabel", "delLabel", "fa-trash", "#", "nsDelLabel", false, true, lineData)
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("addlabelanchor", false, false, false)
	_ = objectsUI.AddWordsInputData("addLabel", "addLabel", "增加标签", "#", "nsAddLabel", false, false, lineData)
	formData.Data = append(formData.Data, lineData)

	tplData["formData"] = formData
	return nil
}

func (n *namespace) getAdditionalJs() []string {
	return n.additionalJs
}
func (n *namespace) getAdditionalCss() []string {
	return n.additionalCss
}

func (n *namespace) addNewResource(c *sysadmServer.Context, module string) error {
	formData, e := utils.GetMultipartData(c, []string{"dcid", "clusterID", "addType", "objContent", "objFile", "newns", "labelKey[]", "labelValue[]"})
	if e != nil {
		return e
	}

	addTypeSlice := formData["addType"].([]string)
	addType := strings.TrimSpace(addTypeSlice[0])
	if addType != "0" && addType != "1" && addType != "2" {
		return fmt.Errorf("add type(%s) is error", addType)
	}

	yamlContent := ""
	if addType == "0" {
		objContentSlice := formData["objContent"].([]string)
		yamlContent = objContentSlice[0]
	}
	if addType == "1" {
		yamlByte, e := utils.ReadUploadedFile(formData["objFile"].(*multipart.FileHeader))
		if e != nil {
			return e
		}
		yamlContent = utils.Interface2String(yamlByte)
	}

	clusterIDSlice := formData["clusterID"].([]string)
	clusterID := clusterIDSlice[0]
	clientSet, e := buildClientSetByClusterID(clusterID)
	if e != nil {
		return e
	}
	if addType == "0" || addType == "1" {
		return k8sclient.ApplyFromYamlByClientSet(yamlContent, clientSet)
	}

	newNsSlice := formData["newns"].([]string)
	newNs := newNsSlice[0]
	nsApplyConfig := appconfigv1.Namespace(newNs)
	lablesMap := make(map[string]string, 0)
	labelKeys := formData["labelKey[]"].([]string)
	labelValues := formData["labelValue[]"].([]string)
	for i, k := range labelKeys {
		value := labelValues[i]
		lablesMap[k] = value
	}
	nsApplyConfig = nsApplyConfig.WithLabels(lablesMap)
	applyOption := metav1.ApplyOptions{
		Force:        true,
		FieldManager: k8sclient.FieldManager,
	}
	_, e = clientSet.CoreV1().Namespaces().Apply(context.Background(), nsApplyConfig, applyOption)

	return e
}

func (n *namespace) delResource(c *sysadmServer.Context, module string, requestData map[string]string) error {
	clusterID := requestData["clusterID"]
	ns := requestData["objID"]
	systemNS := []string{"default", "kube-node-lease", "kube-public", "kube-system"}
	for _, v := range systemNS {
		if ns == v {
			return fmt.Errorf("namespace %s is a system namespace. it was not be deleted", ns)
		}
	}
	clientSet, e := buildClientSetByClusterID(clusterID)
	if e != nil {
		return e
	}

	return clientSet.CoreV1().Namespaces().Delete(context.Background(), ns, metav1.DeleteOptions{})

}

func (n *namespace) buildAddQuotaFormData(tplData map[string]interface{}) error {
	formData, e := objectsUI.InitFormData("addQuota", "addQuota", "POST", "_self", "yes", "addNewQuotaFunc", "")
	if e != nil {
		return e
	}

	lineData := objectsUI.InitLineData("quotanameLine", false, false, false)
	e = objectsUI.AddTextData("quotaname", "quotaname", "", "配额名称", "validateNewName", "", "长度小于63个字母数字或-且以字母数据开开头和结尾的字符串", 30, false, false, lineData)
	if e != nil {
		return e
	}
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("computingLabel", false, false, false)
	e = objectsUI.AddLabelData("computingResource", "mid", "Left", "计算资源配额", false, lineData)
	if e != nil {
		return e
	}
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("cpuLine", false, false, false)
	_ = objectsUI.AddTextData("cpuRequest", "cpuRequest", "", "CPU 请求", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("cpuRequestUnit", "cpuRequestUnit", "Core", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("cpuLimit", "cpuLimit", "", "CPU 上限", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("cpuLimitUnit", "cpuLimitUnit", "Core", "", "", false, false, lineData)
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("memoryLine", false, false, false)
	_ = objectsUI.AddTextData("memRequest", "memRequest", "", "内存 请求", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("memRequestUnit", "memRequestUnit", "Mi", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("memLimit", "memLimit", "", "内存 上限", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("memLimitUnit", "memLimitUnit", "Mi", "", "", false, false, lineData)
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("StorageLabel", false, false, false)
	_ = objectsUI.AddLabelData("StorageResource", "mid", "Left", "存储资源配额", false, lineData)
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("storageLine", false, false, false)
	_ = objectsUI.AddTextData("storageRequest", "storageRequest", "", "存储请求总量", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("storageRequestUnit", "storageRequestUnit", "Gi", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("pvcNum", "pvcNum", "", "存储卷声明数量", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("pvNumUnit", "pvNumUnit", "个", "", "", false, false, lineData)
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("NumberLimitLabel", false, false, false)
	_ = objectsUI.AddLabelData("numberLimit", "mid", "Left", "对象数量配额", false, lineData)
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("numberLimitLine", false, false, false)
	_ = objectsUI.AddTextData("podNum", "podNum", "", "容器组Pod", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("podNumUnit", "podNumUnit", "个", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("serviceNum", "serviceNum", "", "服务Service", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("serviceNumUnit", "serviceNumUnit", "个", "", "", false, false, lineData)
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("numberLimitLine2", false, false, false)
	_ = objectsUI.AddTextData("secretNum", "secretNum", "", "保密字典Secret", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("secretNumUnit", "secretNumUnit", "个", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("configMapNum", "configMapNum", "", "配置项ConfigMap", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("configMapNumUnit", "configMapNumUnit", "个", "", "", false, false, lineData)
	formData.Data = append(formData.Data, lineData)

	tplData["formData"] = formData
	return nil
}

func (n *namespace) addNewQuota(c *sysadmServer.Context) error {
	keys := []string{"quotaname"}
	for k, _ := range allQuotaFormItems {
		keys = append(keys, k)
	}

	formData, e := getFormDataFromMultipartForm(c, keys)
	if e != nil {
		return e
	}

	clientSet, e := buildClientSetByClusterID(formData["clusterID"])
	if e != nil {
		return e
	}
	if formData["addType"] == "0" || formData["addType"] == "1" {
		return k8sclient.ApplyFromYamlByClientSet(formData["yamlContent"], clientSet)
	}

	quotaApplyConfig := appconfigv1.ResourceQuota(formData["quotaname"], formData["namespace"])
	resourceQuotaSpec := appconfigv1.ResourceQuotaSpec()
	resourceList := make(corev1.ResourceList, 0)
	e = parseFormData(formData, resourceList)
	if e != nil {
		return e
	}

	resourceQuotaSpec.Hard = &resourceList
	quotaApplyConfig.WithSpec(resourceQuotaSpec)

	applyOption := metav1.ApplyOptions{
		Force:        true,
		FieldManager: k8sclient.FieldManager,
	}
	_, e = clientSet.CoreV1().ResourceQuotas(formData["namespace"]).Apply(context.Background(), quotaApplyConfig, applyOption)

	return e
}

func parseFormData(formData map[string]string, resourceList corev1.ResourceList) error {
	for k, v := range allQuotaFormItems {
		value, ok := formData[k]
		if !ok {
			return fmt.Errorf("form key %s has not set", k)
		}

		if strings.TrimSpace(value) == "" {
			continue
		}

		value = value + v.UintStr
		requestQuantity, e := resourceapi.ParseQuantity(value)
		if e != nil {
			return e
		}
		resourceList[v.ResourceName] = requestQuantity
	}

	return nil
}

func (n *namespace) buildAddLimitRangeFormData(tplData map[string]interface{}) error {
	clientSet, e := buildClientSetByClusterID(tplData["clusterID"].(string))
	if e != nil {
		return e
	}

	limitRangeList, e := clientSet.CoreV1().LimitRanges(tplData["namespace"].(string)).List(context.Background(), metav1.ListOptions{})
	if len(limitRangeList.Items) > 0 {
		return fmt.Errorf("there is a limitrange in %s namespace", tplData["namespace"].(string))
	}

	formData, e := objectsUI.InitFormData("addLimitRange", "addLimitRange", "POST", "_self", "yes", "addNewLimitRangeFunc", "")
	if e != nil {
		return e
	}

	lineData := objectsUI.InitLineData("limitRangenameLine", false, false, false)
	_ = objectsUI.AddTextData("limitRangeName", "limitRangeName", "", "配额名称", "validateNewName", "", "长度小于63个字母数字或-且以字母数据开开头和结尾的字符串", 30, false, false, lineData)
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("computingLabel", false, false, false)
	_ = objectsUI.AddLabelData("computingResource", "mid", "Left", "计算资源配额", false, lineData)
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("cpuLine", false, false, false)
	_ = objectsUI.AddTextData("cpuMin", "cpuMin", "", "CPU最小值", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("cpuRequestUnit", "cpuRequestUnit", "Core", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("cpuMax", "cpuMax", "", "CPU最大值", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("cpuLimitUnit", "cpuLimitUnit", "Core", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("cpuDefault", "cpuDefault", "", "CPU默认值", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("cpuLimitUnit", "cpuLimitUnit", "Core", "", "", false, false, lineData)
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("memoryLine", false, false, false)
	_ = objectsUI.AddTextData("memMin", "memMin", "", "内存最小值", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("memRequestUnit", "memRequestUnit", "Mi", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("memMax", "memMax", "", "内存最大值", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("memLimitUnit", "memLimitUnit", "Mi", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("memDefault", "memDefault", "", "内存默认值", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("memLimitUnit", "memLimitUnit", "Mi", "", "", false, false, lineData)
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("StorageLabel", false, false, false)
	_ = objectsUI.AddLabelData("storageResource", "mid", "Left", "存储资源配额", false, lineData)
	formData.Data = append(formData.Data, lineData)

	lineData = objectsUI.InitLineData("storageLine", false, false, false)
	_ = objectsUI.AddTextData("storageMin", "storageMin", "", "存储总量最小值", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("storageRequestUnit", "storageRequestUnit", "Gi", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("storageMax", "storageMax", "", "存储总量最大值", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("storageRequestUnit", "storageRequestUnit", "Gi", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("storageDefault", "storageDefault", "", "存储总量默认大值", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("storageRequestUnit", "storageRequestUnit", "Gi", "", "", false, false, lineData)
	formData.Data = append(formData.Data, lineData)

	tplData["formData"] = formData
	return nil
}

func (n *namespace) addNewLimitRange(c *sysadmServer.Context) error {
	keys := []string{"limitRangeName"}
	for k, _ := range allLimitRangeFormItems {
		keys = append(keys, k)
	}

	formData, e := getFormDataFromMultipartForm(c, keys)
	if e != nil {
		return e
	}

	clientSet, e := buildClientSetByClusterID(formData["clusterID"])
	if e != nil {
		return e
	}
	if formData["addType"] == "0" || formData["addType"] == "1" {
		return k8sclient.ApplyFromYamlByClientSet(formData["yamlContent"], clientSet)
	}

	limitType := corev1.LimitTypeContainer
	limitItems := []appconfigv1.LimitRangeItemApplyConfiguration{}
	var minResourceList corev1.ResourceList = nil
	var maxResourceList corev1.ResourceList = nil
	var defaultResourceList corev1.ResourceList = nil
	if formData["cpuMin"] != "" {
		if minResourceList == nil {
			minResourceList = make(corev1.ResourceList, 0)
		}
		requestQuantity, e := resourceapi.ParseQuantity(formData["cpuMin"])
		if e != nil {
			return e
		}
		minResourceList[corev1.ResourceCPU] = requestQuantity
	}

	if formData["cpuMax"] != "" {
		if maxResourceList == nil {
			maxResourceList = make(corev1.ResourceList, 0)
		}
		requestQuantity, e := resourceapi.ParseQuantity(formData["cpuMax"])
		if e != nil {
			return e
		}
		maxResourceList[corev1.ResourceCPU] = requestQuantity
	}
	if formData["cpuDefault"] != "" {
		if defaultResourceList == nil {
			defaultResourceList = make(corev1.ResourceList, 0)
		}
		requestQuantity, e := resourceapi.ParseQuantity(formData["cpuDefault"])
		if e != nil {
			return e
		}
		defaultResourceList[corev1.ResourceCPU] = requestQuantity
	}
	if formData["memMin"] != "" {
		if minResourceList == nil {
			minResourceList = make(corev1.ResourceList, 0)
		}
		requestQuantity, e := resourceapi.ParseQuantity((formData["memMin"] + "Mi"))
		if e != nil {
			return e
		}
		minResourceList[corev1.ResourceMemory] = requestQuantity
	}

	if formData["memMax"] != "" {
		if maxResourceList == nil {
			maxResourceList = make(corev1.ResourceList, 0)
		}
		requestQuantity, e := resourceapi.ParseQuantity((formData["memMax"] + "Mi"))
		if e != nil {
			return e
		}
		maxResourceList[corev1.ResourceMemory] = requestQuantity
	}
	if formData["memDefault"] != "" {
		if defaultResourceList == nil {
			defaultResourceList = make(corev1.ResourceList, 0)
		}
		requestQuantity, e := resourceapi.ParseQuantity((formData["memDefault"] + "Mi"))
		if e != nil {
			return e
		}
		defaultResourceList[corev1.ResourceMemory] = requestQuantity
	}
	limitRangeItemMem := appconfigv1.LimitRangeItemApplyConfiguration{
		Type:    &limitType,
		Max:     &maxResourceList,
		Min:     &minResourceList,
		Default: &defaultResourceList,
	}
	limitItems = append(limitItems, limitRangeItemMem)

	var storageMinResourceList corev1.ResourceList = nil
	var storageMaxResourceList corev1.ResourceList = nil
	var storageDefaultResourceList corev1.ResourceList = nil

	if formData["storageMin"] != "" {
		if storageMinResourceList == nil {
			storageMinResourceList = make(corev1.ResourceList, 0)
		}
		requestQuantity, e := resourceapi.ParseQuantity((formData["storageMin"] + "Gi"))
		if e != nil {
			return e
		}
		storageMinResourceList[corev1.ResourceStorage] = requestQuantity
	}
	if formData["storageMax"] != "" {
		if storageMaxResourceList == nil {
			storageMaxResourceList = make(corev1.ResourceList, 0)
		}
		requestQuantity, e := resourceapi.ParseQuantity((formData["storageMax"] + "Gi"))
		if e != nil {
			return e
		}
		storageMaxResourceList[corev1.ResourceStorage] = requestQuantity
	}
	if formData["storageDefault"] != "" {
		if storageDefaultResourceList == nil {
			storageDefaultResourceList = make(corev1.ResourceList, 0)
		}

		requestQuantity, e := resourceapi.ParseQuantity((formData["storageDefault"] + "Gi"))
		if e != nil {
			return e
		}
		storageDefaultResourceList[corev1.ResourceStorage] = requestQuantity
	}
	storageType := corev1.LimitTypePersistentVolumeClaim
	limitRangeItemStorage := appconfigv1.LimitRangeItemApplyConfiguration{
		Type:    &storageType,
		Max:     &storageMaxResourceList,
		Min:     &storageMinResourceList,
		Default: &storageDefaultResourceList,
	}
	limitItems = append(limitItems, limitRangeItemStorage)

	limitRangeApplyConfig := appconfigv1.LimitRange(formData["limitRangeName"], formData["namespace"])
	limitRangeSpec := appconfigv1.LimitRangeSpec()
	limitRangeSpec.Limits = limitItems
	limitRangeApplyConfig.Spec = limitRangeSpec

	applyOption := metav1.ApplyOptions{
		Force:        true,
		FieldManager: k8sclient.FieldManager,
	}
	_, e = clientSet.CoreV1().LimitRanges(formData["namespace"]).Apply(context.Background(), limitRangeApplyConfig, applyOption)

	return e
}

// for list Quota data
func listQuotaData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {

	var dataList []map[string]interface{}

	clientSet, e := createClientSet(selectedCluster)
	if clientSet == nil {
		return 0, dataList, e
	}

	nsStr, orderField, direction := checkRequestData(selectedNS, quotaListDefaultOrderField, quotaListDefaultOrderDirection, requestData)
	objList, e := clientSet.CoreV1().ResourceQuotas(nsStr).List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return 0, dataList, e
	}
	totalNum := len(objList.Items)
	if totalNum < 1 {
		return 0, dataList, nil
	}
	endCount := totalNum
	if endCount > startPos+runData.pageInfo.NumPerPage {
		endCount = startPos + runData.pageInfo.NumPerPage
	}

	var objItems []interface{}
	for _, item := range objList.Items {
		objItems = append(objItems, item)
	}

	for field, fn := range quotaListAllOrderFields {
		if field == orderField {
			if direction == "1" {
				sort.Sort(objectsUI.SortData{Data: objItems, By: fn})
			} else {
				sort.Sort(sort.Reverse(objectsUI.SortData{Data: objItems, By: fn}))
			}

		}
	}

	for i := startPos; i < endCount; i++ {
		interfaceData := objItems[i]
		objData, ok := interfaceData.(corev1.ResourceQuota)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not ResourceQuota schema")
		}
		lineMap := make(map[string]interface{}, 0)
		hard := objData.Spec.Hard
		lineMap["objectID"] = objData.Name
		lineMap["TD1"] = objData.Name
		lineMap["TD2"] = objData.Namespace
		requestCPU := hard[corev1.ResourceRequestsCPU]
		lineMap["TD3"] = requestCPU.String()
		limitCPU := hard[corev1.ResourceLimitsCPU]
		lineMap["TD4"] = limitCPU.String()
		requestMemory := hard[corev1.ResourceRequestsMemory]
		lineMap["TD5"] = requestMemory.String()
		limitMemory := hard[corev1.ResourceLimitsMemory]
		lineMap["TD6"] = limitMemory.String()
		lineMap["TD7"] = objData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		popmenuitems := "0,1,2"
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortQuotaByName(p, q interface{}) bool {
	pData, ok := p.(corev1.ResourceQuota)
	if !ok {
		return false
	}
	qData, ok := q.(corev1.ResourceQuota)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortQuotaByCreatetime(p, q interface{}) bool {
	pData, ok := p.(corev1.ResourceQuota)
	if !ok {
		return false
	}
	qData, ok := q.(corev1.ResourceQuota)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

func limitRangeDel(c *sysadmServer.Context, module string, requestData map[string]string) error {
	ns := requestData["objID"]
	clusterID := requestData["clusterID"]

	clientSet, e := buildClientSetByClusterID(clusterID)
	if e != nil {
		return e
	}

	objList, e := clientSet.CoreV1().LimitRanges(ns).List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return e
	}
	if len(objList.Items) < 1 {
		return fmt.Errorf("there is not LimitRange in %s namespace to be delete", ns)
	}

	if len(objList.Items) > 1 {
		return fmt.Errorf("there are more than on  LimitRange in %s namespace to be delete", ns)
	}

	item := objList.Items[0]
	limitRangeName := item.Name

	return clientSet.CoreV1().LimitRanges(ns).Delete(context.Background(), limitRangeName, metav1.DeleteOptions{})

}

func (n *namespace) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	tplData["mainCategory"] = n.mainModuleName
	tplData["subCategory"] = n.getModuleName()

	switch action {
	case "quotaDetail":
		tplData["thirdCategory"] = "资源配额详情"
	case "limitRangeDetail":
		tplData["thirdCategory"] = "默认资源配额详情"
	case "detail":
		tplData["thirdCategory"] = "命名空间详情"
	default:

	}

	objID := strings.TrimSpace(requestData["objID"])
	if objID == "" {
		return fmt.Errorf("can not get resource information without resource name")
	}

	clusterID := requestData["clusterID"]
	clientSet, e := createClientSet(clusterID)
	if e != nil {
		return e
	}
	var objEntity interface{} = nil
	switch action {
	case "quotaDetail", "limitRangeDetail":

		if action == "quotaDetail" {
			ns := strings.TrimSpace(requestData["namespace"])
			quotaInfo, e := clientSet.CoreV1().ResourceQuotas(ns).Get(context.Background(), objID, metav1.GetOptions{})
			if e != nil {
				return e
			}
			objEntity = quotaInfo
		}
		if action == "limitRangeDetail" {
			ns := strings.TrimSpace(requestData["objID"])
			if ns == "" {
				return fmt.Errorf("can not get resource information without name of namespace")
			}
			limitRangeList, e := clientSet.CoreV1().LimitRanges(ns).List(context.Background(), metav1.ListOptions{})
			if e != nil {
				return e
			}
			if len(limitRangeList.Items) != 1 {
				return fmt.Errorf("there is not LimitRange or there are more than one LimitRange in %s namespace", ns)
			}
			items := limitRangeList.Items
			name := items[0].Name
			limitRangeInfo, e := clientSet.CoreV1().LimitRanges(ns).Get(context.Background(), name, metav1.GetOptions{})
			if e != nil {
				return e
			}
			limitRangeInfo.ManagedFields = nil
			objEntity = limitRangeInfo
		}
	case "detail":
		nsInfo, e := clientSet.CoreV1().Namespaces().Get(context.Background(), objID, metav1.GetOptions{})
		if e != nil {
			return e
		}
		objEntity = nsInfo
	}

	switch action {
	case "quotaDetail":
		return buildQuotaDetailData(objEntity, tplData)
	case "limitRangeDetail":
		return buildLimitRangeData(objEntity, tplData)
	case "detail":
		return buildNamespaceData(objEntity, tplData)
	}

	return fmt.Errorf("no action %s was defined in module namespace")
}

func buildQuotaDetailData(objEntity interface{}, tplData map[string]interface{}) error {
	quotaData, ok := objEntity.(*corev1.ResourceQuota)
	if !ok {
		return fmt.Errorf("the data is not ResourceQuota schema")
	}

	quotaData.ManagedFields = nil
	scheme := runtime.NewScheme()
	corev1.AddToScheme(scheme)
	codec := serializer.NewCodecFactory(scheme).LegacyCodec(corev1.SchemeGroupVersion)
	output, _ := runtime.Encode(codec, quotaData)
	resourceYamlContent, e := yaml.JSONToYAML(output)
	if e != nil {
		return e
	}
	tplData["resourceYamlContent"] = strings.TrimSpace(string(resourceYamlContent))
	data := []*objectsUI.LineData{}

	lineData := objectsUI.InitLineData("generalInfoLabel", false, false, false)
	_ = objectsUI.AddLabelData("generalInfo", "mid", "Left", "通用信息", false, lineData)
	data = append(data, lineData)

	lineData = objectsUI.InitLineData("generalInfoLine", false, false, false)
	_ = objectsUI.AddTitleValueData("quotaName", "配额名称", quotaData.Name, "", "", lineData)
	_ = objectsUI.AddTitleValueData("quotaName", "命名空间", quotaData.Namespace, "", "", lineData)
	_ = objectsUI.AddTitleValueData("createTime", "创建时间", quotaData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat), "", "", lineData)
	data = append(data, lineData)

	labelsStr := ""
	for k, v := range quotaData.Labels {
		if labelsStr == "" {
			labelsStr = k + ":  " + v
		} else {
			labelsStr = labelsStr + ", " + k + ":  " + v
		}
	}
	lineData = objectsUI.InitLineData("labelsLine", false, false, false)
	_ = objectsUI.AddTitleValueData("quotaLabelsTitle", "标签", labelsStr, "", "", lineData)
	data = append(data, lineData)

	lineData = objectsUI.InitLineData("computingResourceLabelLine", false, false, false)
	_ = objectsUI.AddLabelData("computingResourceLabel", "mid", "Left", "计算资源配额", false, lineData)
	data = append(data, lineData)

	hard := quotaData.Spec.Hard
	lineData = objectsUI.InitLineData("CpuQuotaLine", false, false, false)
	cpuRequest := hard[corev1.ResourceRequestsCPU]
	limitCPU := hard[corev1.ResourceLimitsCPU]
	_ = objectsUI.AddTitleValueData("cpuRequest", "CPU请求", cpuRequest.String(), "", "", lineData)
	_ = objectsUI.AddTitleValueData("cpuLimit", "CPU上限", limitCPU.String(), "", "", lineData)
	data = append(data, lineData)

	lineData = objectsUI.InitLineData("memQuotaLine", false, false, false)
	requestMemory := hard[corev1.ResourceRequestsMemory]
	limitMemory := hard[corev1.ResourceLimitsMemory]
	_ = objectsUI.AddTitleValueData("memRequest", "内存请求", requestMemory.String(), "", "", lineData)
	_ = objectsUI.AddTitleValueData("memLimit", "内存上限", limitMemory.String(), "", "", lineData)
	data = append(data, lineData)

	lineData = objectsUI.InitLineData("storageResourceLabelLine", false, false, false)
	_ = objectsUI.AddLabelData("storageResourceLabel", "mid", "Left", "存储资源配额", false, lineData)
	data = append(data, lineData)

	lineData = objectsUI.InitLineData("storageQuotaLine", false, false, false)
	storageRequest := hard[corev1.ResourceRequestsStorage]
	pvcNum := hard[corev1.ResourcePersistentVolumeClaims]
	_ = objectsUI.AddTitleValueData("storageRequest", "存储请求总量", storageRequest.String(), "", "", lineData)
	_ = objectsUI.AddTitleValueData("pvcNum", "存储卷声明数量", pvcNum.String(), "", "", lineData)
	data = append(data, lineData)

	lineData = objectsUI.InitLineData("objectNumLine", false, false, false)
	_ = objectsUI.AddLabelData("objectLabel", "mid", "Left", "对象数量配额", false, lineData)
	data = append(data, lineData)

	podNum := hard[corev1.ResourcePods]
	serviceNum := hard[corev1.ResourceServices]
	secretsNum := hard[corev1.ResourceSecrets]
	cmNum := hard[corev1.ResourceConfigMaps]

	lineData = objectsUI.InitLineData("podQuotaLine", false, false, false)
	_ = objectsUI.AddTitleValueData("podNumber", "容器组Pod数量", podNum.String(), "", "", lineData)
	_ = objectsUI.AddTitleValueData("serviceNumber", "服务Service数量", serviceNum.String(), "", "", lineData)
	data = append(data, lineData)

	lineData = objectsUI.InitLineData("cmQuotaLine", false, false, false)
	_ = objectsUI.AddTitleValueData("secretNumber", "保密字典Secret数量", secretsNum.String(), "", "", lineData)
	_ = objectsUI.AddTitleValueData("cmNumber", "配置项configMap数量", cmNum.String(), "", "", lineData)
	data = append(data, lineData)

	tplData["data"] = data

	return nil
}

func buildLimitRangeData(objEntity interface{}, tplData map[string]interface{}) error {
	limitRangeData, ok := objEntity.(*corev1.LimitRange)
	if !ok {
		return fmt.Errorf("the data is not LimitRange schema")
	}

	limitRangeData.ManagedFields = nil
	scheme := runtime.NewScheme()
	corev1.AddToScheme(scheme)
	codec := serializer.NewCodecFactory(scheme).LegacyCodec(corev1.SchemeGroupVersion)
	output, _ := runtime.Encode(codec, limitRangeData)
	resourceYamlContent, e := yaml.JSONToYAML(output)
	if e != nil {
		return e
	}
	tplData["resourceYamlContent"] = strings.TrimSpace(string(resourceYamlContent))

	data := []*objectsUI.LineData{}
	lineData := objectsUI.InitLineData("generalInfoLabel", false, false, false)
	_ = objectsUI.AddLabelData("generalInfo", "mid", "Left", "通用信息", false, lineData)
	data = append(data, lineData)

	lineData = objectsUI.InitLineData("generalInfoLine", false, false, false)
	_ = objectsUI.AddTitleValueData("quotaName", "配额名称", limitRangeData.Name, "", "", lineData)
	_ = objectsUI.AddTitleValueData("quotaName", "命名空间", limitRangeData.Namespace, "", "", lineData)
	_ = objectsUI.AddTitleValueData("createTime", "创建时间", limitRangeData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat), "", "", lineData)
	data = append(data, lineData)

	labelsStr := ""
	for k, v := range limitRangeData.Labels {
		if labelsStr == "" {
			labelsStr = k + ":  " + v
		} else {
			labelsStr = labelsStr + ", " + k + ":  " + v
		}
	}
	lineData = objectsUI.InitLineData("labelsLine", false, false, false)
	_ = objectsUI.AddTitleValueData("quotaLabelsTitle", "标签", labelsStr, "", "", lineData)
	data = append(data, lineData)

	var computingLimitRangeItem corev1.LimitRangeItem
	var storageLimitRangeItem corev1.LimitRangeItem
	limitItems := limitRangeData.Spec.Limits
	for _, item := range limitItems {
		if item.Type == corev1.LimitTypeContainer || item.Type == corev1.LimitTypePod {
			computingLimitRangeItem = item
		} else {
			storageLimitRangeItem = item
		}

	}
	computingMax := computingLimitRangeItem.Max
	computingMin := computingLimitRangeItem.Min
	computingDefault := computingLimitRangeItem.Default
	storageMax := storageLimitRangeItem.Max
	storageMin := storageLimitRangeItem.Min
	storageDefault := storageLimitRangeItem.Default

	cpuMaxStr := "-"
	if cpuMax, ok := computingMax[corev1.ResourceCPU]; ok {
		cpuMaxStr = cpuMax.String()
	}
	cpuMinStr := "-"
	if cpuMin, ok := computingMin[corev1.ResourceCPU]; ok {
		cpuMinStr = cpuMin.String()
	}
	cpuDefaultStr := "-"
	if cpuDefault, ok := computingDefault[corev1.ResourceCPU]; ok {
		cpuDefaultStr = cpuDefault.String()
	}

	memMaxStr := "-"
	if memMax, ok := computingMax[corev1.ResourceMemory]; ok {
		memMaxStr = memMax.String()
	}
	memMinStr := "-"
	if memMin, ok := computingMin[corev1.ResourceMemory]; ok {
		memMinStr = memMin.String()
	}
	memDefaultStr := "-"
	if memDefault, ok := computingDefault[corev1.ResourceMemory]; ok {
		memDefaultStr = memDefault.String()
	}

	storageMaxStr := "-"
	if sMax, ok := storageMax[corev1.ResourceStorage]; ok {
		storageMaxStr = sMax.String()
	}
	storageMinStr := "-"
	if sMin, ok := storageMin[corev1.ResourceStorage]; ok {
		storageMinStr = sMin.String()
	}
	storageDefaultStr := "-"
	if sDefault, ok := storageDefault[corev1.ResourceStorage]; ok {
		storageDefaultStr = sDefault.String()
	}

	lineData = objectsUI.InitLineData("computingResourceLabelLine", false, false, false)
	_ = objectsUI.AddLabelData("computingResourceLabel", "mid", "Left", "计算资源配额", false, lineData)
	data = append(data, lineData)

	lineData = objectsUI.InitLineData("CpuQuotaLine", false, false, false)
	_ = objectsUI.AddTitleValueData("cpuMin", "CPU最小值", cpuMinStr, "", "", lineData)
	_ = objectsUI.AddTitleValueData("cpuMax", "CPU最大值", cpuMaxStr, "", "", lineData)
	_ = objectsUI.AddTitleValueData("cpuDefault", "CPU默认值", cpuDefaultStr, "", "", lineData)
	data = append(data, lineData)

	lineData = objectsUI.InitLineData("memQuotaLine", false, false, false)
	_ = objectsUI.AddTitleValueData("memMin", "内存最小值", memMinStr, "", "", lineData)
	_ = objectsUI.AddTitleValueData("memMax", "内存最大值", memMaxStr, "", "", lineData)
	_ = objectsUI.AddTitleValueData("memDefault", "内存默认值", memDefaultStr, "", "", lineData)
	data = append(data, lineData)

	lineData = objectsUI.InitLineData("storageResourceLabelLine", false, false, false)
	_ = objectsUI.AddLabelData("storageResourceLabel", "mid", "Left", "存储资源配额", false, lineData)
	data = append(data, lineData)

	lineData = objectsUI.InitLineData("storageQuotaLine", false, false, false)
	_ = objectsUI.AddTitleValueData("storageMin", "存储总量最小值", storageMinStr, "", "", lineData)
	_ = objectsUI.AddTitleValueData("storageMax", "存储总量最大值", storageMaxStr, "", "", lineData)
	_ = objectsUI.AddTitleValueData("storageDefault", "存储总量默认值", storageDefaultStr, "", "", lineData)
	data = append(data, lineData)

	tplData["data"] = data

	return nil
}

func buildNamespaceData(objEntity interface{}, tplData map[string]interface{}) error {
	nsData, ok := objEntity.(*corev1.Namespace)
	if !ok {
		return fmt.Errorf("the data is not Namespace schema")
	}

	nsData.ManagedFields = nil
	scheme := runtime.NewScheme()
	corev1.AddToScheme(scheme)
	codec := serializer.NewCodecFactory(scheme).LegacyCodec(corev1.SchemeGroupVersion)
	output, _ := runtime.Encode(codec, nsData)
	resourceYamlContent, e := yaml.JSONToYAML(output)
	if e != nil {
		return e
	}
	tplData["resourceYamlContent"] = strings.TrimSpace(string(resourceYamlContent))

	data := []*objectsUI.LineData{}

	lineData := objectsUI.InitLineData("generalInfoLine", false, false, false)
	_ = objectsUI.AddTitleValueData("nsName", "命名空间名称", nsData.Name, "", "", lineData)
	data = append(data, lineData)

	labelsStr := ""
	for k, v := range nsData.Labels {
		if labelsStr == "" {
			labelsStr = k + ":  " + v
		} else {
			labelsStr = labelsStr + ", " + k + ":  " + v
		}
	}
	lineData = objectsUI.InitLineData("labelsLine", false, false, false)
	_ = objectsUI.AddTitleValueData("quotaLabelsTitle", "标签", labelsStr, "", "", lineData)
	data = append(data, lineData)

	tplData["data"] = data

	return nil
}

func quotaDel(c *sysadmServer.Context, module string, requestData map[string]string) error {
	quotaName := requestData["objID"]
	clusterID := requestData["clusterID"]
	ns := requestData["namespace"]

	clientSet, e := buildClientSetByClusterID(clusterID)
	if e != nil {
		return e
	}

	return clientSet.CoreV1().ResourceQuotas(ns).Delete(context.Background(), quotaName, metav1.DeleteOptions{})

}

func (n *namespace) getTemplateFile(action string) string {
	// TODO

	return n.templateFile
}
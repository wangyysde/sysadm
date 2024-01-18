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
	"html/template"
	apiNetworkingV1 "k8s.io/api/networking/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	applyConfigNetworkingV1 "k8s.io/client-go/applyconfigurations/networking/v1"
	"sort"
	"strconv"
	"strings"
	"sysadm/k8sclient"
	"sysadm/objectsUI"
	"sysadm/utils"
)

func (i *ingress) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortIngressByName, "TD6": sortIngressByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "Hosts", "TD4": "标签", "TD5": "IngressClassName", "TD6": "创建时间"}

	i.mainModuleName = "服务管理"
	i.moduleName = "ingress"
	i.allPopMenuItems = allPopMenuItems
	i.allListItems = allListItems
	i.addButtonTile = "添加Ingress"
	i.isSearchForm = "no"
	i.allOrderFields = allOrderFields
	i.defaultOrderField = "TD1"
	i.defaultOrderDirection = "1"
	i.namespaced = true
	i.moduleID = "ingress"
	i.additionalJs = []string{}
	i.additionalCss = []string{}
	i.templateFile = ""

}

func (i *ingress) getMainModuleName() string {
	return i.mainModuleName
}

func (i *ingress) getModuleName() string {
	return i.moduleName
}

func (i *ingress) getAddButtonTitle() string {
	return i.addButtonTile
}

func (i *ingress) getIsSearchForm() string {
	return i.isSearchForm
}

func (i *ingress) getAllPopMenuItems() []string {
	return i.allPopMenuItems
}

func (i *ingress) getAllListItems() map[string]string {
	return i.allListItems
}

func (i *ingress) getDefaultOrderField() string {
	return i.defaultOrderField
}

func (i *ingress) getDefaultOrderDirection() string {
	return i.defaultOrderDirection
}

func (i *ingress) getAllorderFields() map[string]objectsUI.SortBy {
	return i.allOrderFields
}

func (i *ingress) getNamespaced() bool {
	return i.namespaced
}

// for Ingress
func (i *ingress) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	clientSet, e := createClientSet(selectedCluster)
	if clientSet == nil {
		return 0, dataList, e
	}

	nsStr, orderField, direction := checkRequestData(selectedNS, i.defaultOrderField, i.defaultOrderDirection, requestData)
	ingressList, e := clientSet.NetworkingV1().Ingresses(nsStr).List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return 0, dataList, e
	}
	totalNum := len(ingressList.Items)
	if totalNum < 1 {
		return 0, dataList, nil
	}
	endCount := totalNum
	if endCount > startPos+runData.pageInfo.NumPerPage {
		endCount = startPos + runData.pageInfo.NumPerPage
	}

	var ingressItems []interface{}
	for _, item := range ingressList.Items {
		ingressItems = append(ingressItems, item)
	}

	moduleAllOrderFields := i.allOrderFields
	for field, fn := range moduleAllOrderFields {
		if field == orderField {
			if direction == "1" {
				sort.Sort(objectsUI.SortData{Data: ingressItems, By: fn})
			} else {
				sort.Sort(sort.Reverse(objectsUI.SortData{Data: ingressItems, By: fn}))
			}

		}
	}

	for i := startPos; i < endCount; i++ {
		interfaceData := ingressItems[i]
		ingressData, ok := interfaceData.(networkingv1.Ingress)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not Ingress schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = ingressData.Name
		lineMap["TD1"] = ingressData.Name
		lineMap["TD2"] = ingressData.Namespace
		hostStr := ""
		for _, rules := range ingressData.Spec.Rules {
			hostStr = hostStr + "<div>" + rules.Host + "</div>"
		}
		lineMap["TD3"] = template.HTML(hostStr)
		lineMap["TD4"] = objectsUI.ConvertMap2HTML(ingressData.Labels)
		lineMap["TD5"] = *ingressData.Spec.IngressClassName
		lineMap["TD6"] = ingressData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		popmenuitems := "0,1"
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func (i *ingress) getModuleID() string {
	return i.moduleID
}

func (i *ingress) buildAddFormData(tplData map[string]interface{}) error {
	tplData["thirdCategory"] = "创建Ingress"
	formData, e := objectsUI.InitFormData("addIngress", "addIngress", "POST", "_self", "yes", "addWorkload", "")
	if e != nil {
		return e
	}
	tplData["formData"] = formData

	e = buildIngressBasiceFormData(tplData)
	if e != nil {
		return e
	}

	return nil

	return nil
}

func (i *ingress) getAdditionalJs() []string {
	return i.additionalJs
}
func (i *ingress) getAdditionalCss() []string {
	return i.additionalCss
}

func (i *ingress) addNewResource(c *sysadmServer.Context, module string) error {
	requestKeys := []string{"dcid", "clusterID", "namespace", "addType", "nsSelected", "name", "labelKey[]", "labelValue[]", "annotationKey[]", "annotationValue[]"}
	requestKeys = append(requestKeys, "nsSelected", "ingressClasses", "enableDefaultBackend", "defaultService", "defaultServicePort", "domain", "matchMethod[]")
	requestKeys = append(requestKeys, "enabledHttps", "secret", "matchUriPath[]", "matchService[]", "matchPort[]")
	formData, e := utils.GetMultipartData(c, requestKeys)
	if e != nil {
		return e
	}

	ns := formData["nsSelected"].([]string)
	name := formData["name"].([]string)
	ingressApplyConfiguration := applyConfigNetworkingV1.Ingress(name[0], ns[0])

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
		labels[defaultIngressLabelKey] = name[0]
	}
	for k, v := range extraLabels {
		labels[k] = v
	}
	ingressApplyConfiguration = ingressApplyConfiguration.WithLabels(labels)

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
	ingressApplyConfiguration = ingressApplyConfiguration.WithAnnotations(annotations)

	ingressClasses := strings.TrimSpace(formData["ingressClasses"].([]string)[0])
	if ingressClasses == "" || ingressClasses == "0" {
		return fmt.Errorf("you should select an ingressClass for new ingress")
	}
	ingressSpecApplyConfgiguration := applyConfigNetworkingV1.IngressSpecApplyConfiguration{}
	ingressSpecApplyConfgiguration.IngressClassName = &ingressClasses

	enableDefaultBackend := strings.TrimSpace(formData["enableDefaultBackend"].([]string)[0])
	if enableDefaultBackend == "1" {
		defaultService := strings.TrimSpace(formData["defaultService"].([]string)[0])
		defaultServicePort := strings.TrimSpace(formData["defaultServicePort"].([]string)[0])
		defaultServicePortInt, e := strconv.Atoi(defaultServicePort)
		if e != nil {
			return e
		}
		serviceBackendPortApplyConfiguration := applyConfigNetworkingV1.ServiceBackendPort()
		serviceBackendPortApplyConfiguration = serviceBackendPortApplyConfiguration.WithNumber(int32(defaultServicePortInt))
		if defaultService == "" || defaultService == "0" {
			return fmt.Errorf("default service name must not be empty")
		}
		defaultServiceBackendApplyConfiguration := applyConfigNetworkingV1.IngressServiceBackendApplyConfiguration{Name: &defaultService, Port: serviceBackendPortApplyConfiguration}
		defaultIngressBackendApplyConfiguration := applyConfigNetworkingV1.IngressBackendApplyConfiguration{Service: &defaultServiceBackendApplyConfiguration}
		ingressSpecApplyConfgiguration.DefaultBackend = &defaultIngressBackendApplyConfiguration
	}

	domain := strings.TrimSpace(formData["domain"].([]string)[0])
	enabledHttps := strings.TrimSpace(formData["enabledHttps"].([]string)[0])
	if domain == "" {
		return fmt.Errorf("domain must be specified")
	}
	if enabledHttps != "" && enabledHttps != "0" {
		secret := strings.TrimSpace(formData["secret"].([]string)[0])
		if secret == "" || secret == "0" {
			return fmt.Errorf("secret name must be specified when enabled HTTPS")
		}
		var tlsConfigs []applyConfigNetworkingV1.IngressTLSApplyConfiguration
		tlsConfig := applyConfigNetworkingV1.IngressTLSApplyConfiguration{Hosts: []string{domain}, SecretName: &secret}
		tlsConfigs = append(tlsConfigs, tlsConfig)
		ingressSpecApplyConfgiguration.TLS = tlsConfigs
	}

	matchMethods := formData["matchMethod[]"].([]string)
	matchUriPaths := formData["matchUriPath[]"].([]string)
	matchServices := formData["matchService[]"].([]string)
	matchPorts := formData["matchPort[]"].([]string)
	if len(matchUriPaths) < 1 {
		return fmt.Errorf("at least one rule for an ingress")
	}
	var httpIngressPaths []applyConfigNetworkingV1.HTTPIngressPathApplyConfiguration
	for i, p := range matchUriPaths {
		method := matchMethods[i]
		s := strings.TrimSpace(matchServices[i])
		port := strings.TrimSpace(matchPorts[i])
		portInt, e := strconv.Atoi(port)
		if e != nil {
			return e
		}
		p = strings.TrimSpace(p)
		methodPath := apiNetworkingV1.PathType(method)
		serviceBackendPortApplyConfiguration := applyConfigNetworkingV1.ServiceBackendPort()
		serviceBackendPortApplyConfiguration = serviceBackendPortApplyConfiguration.WithNumber(int32(portInt))
		serviceBackendApplyConfiguration := applyConfigNetworkingV1.IngressServiceBackendApplyConfiguration{Name: &s, Port: serviceBackendPortApplyConfiguration}
		serviceIngressBackendApplyConfiguration := applyConfigNetworkingV1.IngressBackendApplyConfiguration{Service: &serviceBackendApplyConfiguration}
		httpIngressPath := applyConfigNetworkingV1.HTTPIngressPathApplyConfiguration{}
		httpIngressPath.Path = &p
		httpIngressPath.PathType = &methodPath
		httpIngressPath.Backend = &serviceIngressBackendApplyConfiguration
		httpIngressPaths = append(httpIngressPaths, httpIngressPath)
	}
	httpIngressRuleValue := applyConfigNetworkingV1.HTTPIngressRuleValueApplyConfiguration{Paths: httpIngressPaths}
	ingressRuleValueApplyConfiguration := applyConfigNetworkingV1.IngressRuleValueApplyConfiguration{HTTP: &httpIngressRuleValue}
	ingressRuleApplyConfiguration := applyConfigNetworkingV1.IngressRuleApplyConfiguration{Host: &domain, IngressRuleValueApplyConfiguration: ingressRuleValueApplyConfiguration}
	var ingressRuleApplyConfigurations []applyConfigNetworkingV1.IngressRuleApplyConfiguration
	ingressRuleApplyConfigurations = append(ingressRuleApplyConfigurations, ingressRuleApplyConfiguration)
	ingressSpecApplyConfgiguration.Rules = ingressRuleApplyConfigurations

	ingressApplyConfiguration = ingressApplyConfiguration.WithSpec(&ingressSpecApplyConfgiguration)
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
	_, e = clientSet.NetworkingV1().Ingresses(ns[0]).Apply(context.Background(), ingressApplyConfiguration, applyOption)

	return e
}

func (i *ingress) delResource(c *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (i *ingress) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

func (i *ingress) getTemplateFile(action string) string {
	switch action {
	case "list":
		return ingressTemplateFiles["list"]
	case "addform":
		return ingressTemplateFiles["addform"]
	default:
		return ""

	}
	return ""

}

func buildIngressBasiceFormData(tplData map[string]interface{}) error {
	clusterID := tplData["clusterID"].(string)
	if clusterID == "" || clusterID == "0" {
		return fmt.Errorf("cluster must be specified when add a new ingress")
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
		for _, n := range denyIngressNSList {
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
	e = objectsUI.AddSelectData("nsSelectedID", "nsSelected", defaultSelectedNs, "getNameList", "nsChangedForIngressAdd", "选择命名空间", "", 1, false, false, nsOptions, lineData)
	if e != nil {
		return e
	}
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("nameLine", false, false, false)
	_ = objectsUI.AddTextData("name", "name", "", "Ingress名称", "validateNewName", "addWorkloadValidateNewName", "长度小于63个字母数字或-且以字母数据开开头和结尾的字符串", 30, false, false, lineData)
	basicData = append(basicData, lineData)

	ingressClassesList, e := clientSet.NetworkingV1().IngressClasses().List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return e
	}

	ingressClassesOptions := make(map[string]string, 0)
	for _, item := range ingressClassesList.Items {
		ingressClassesOptions[item.Name] = item.Name
	}
	if len(ingressClassesOptions) < 1 {
		ingressClassesOptions["0"] = "===请首先创建IngressClasses对象==="
	}
	lineData = objectsUI.InitLineData("ingressClassesLine", false, false, false)
	_ = objectsUI.AddSelectData("ingressClassesID", "ingressClasses", defaultSelectedNs, "", "", "选择IngressClass", "", 1, false, false, ingressClassesOptions, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("routeRuleLabelID", false, false, false)
	_ = objectsUI.AddLabelData("routeRuleID", "mid", "Left", "路由规则", false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("enableDefaultBackendID", false, false, false)
	var checkboxOptions []objectsUI.Option
	checkboxOptions, _ = objectsUI.AddCheckBoxOption("启用默认后端", "1", false, false, checkboxOptions)
	_ = objectsUI.AddCheckBoxData("enableDefaultBackendID", "enableDefaultBackend", "", "switchDisplayStatusForCheckBoxBlock", false, checkboxOptions, lineData)
	serviceOptions := make(map[string]string, 0)
	serviceOptions["0"] = "===选择后端服务==="
	_ = objectsUI.AddSelectData("defaultServiceID", "defaultService", "0", "", "", "服务", "", 1, false, true, serviceOptions, lineData)
	_ = objectsUI.AddTextData("defaultServicePortID", "defaultServicePort", "", "端口", "", "", "", 10, false, true, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("domainLineID", false, false, false)
	_ = objectsUI.AddTextData("domain", "domain", "", "域名", "", "", "", 30, false, false, lineData)
	var httpsOptions []objectsUI.Option
	httpsOptions, _ = objectsUI.AddCheckBoxOption("启用HTTPS", "1", false, false, httpsOptions)
	_ = objectsUI.AddCheckBoxData("enabledHttpsID", "enabledHttps", "", "switchDisplaySecretForIngressAdd", false, httpsOptions, lineData)
	secretOptions := make(map[string]string, 0)
	secretOptions["0"] = "===选择存储TLS证书的Secret==="
	_ = objectsUI.AddSelectData("secretID", "secret", "0", "", "", "密文", "", 1, false, true, secretOptions, lineData)
	basicData = append(basicData, lineData)

	matchMethodOptions := make(map[string]string, 0)
	matchMethodOptions["Prefix"] = "前缀匹配"
	matchMethodOptions["Exact"] = "精确匹配"
	matchMethodOptions["ImplementationSpecific"] = "由控制器决定匹配方式"

	matchServiceOptions := make(map[string]string, 0)
	matchServiceOptions["0"] = "===选择后端服务==="

	lineData = objectsUI.InitLineData("PathMatchRule", true, true, false)
	_ = objectsUI.AddSelectData("matchMethodID", "matchMethod[]", "Prefix", "", "", "匹配方法", "", 1, false, false, matchMethodOptions, lineData)
	_ = objectsUI.AddTextData("matchUriPathID", "matchUriPath[]", "", "映射路径", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddSelectData("matchServiceID", "matchService[]", "0", "", "", "后端服务", "", 1, false, false, matchServiceOptions, lineData)
	_ = objectsUI.AddTextData("matchPortID", "matchPort[]", "", "端口", "", "", "", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("delMatchRuleID", "delMatchRule", "fa-trash", "#", "workloadDelSelector", false, true, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("PathMatchRule", false, false, false)
	_ = objectsUI.AddWordsInputData("PathMatchRuleAnchor", "PathMatchRuleAnchor", "添加匹配条件", "#", "workloadAddSelectorBlock", false, false, lineData)
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

	tplData["BasicData"] = basicData

	return nil
}

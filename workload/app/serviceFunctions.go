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
	"html/template"
	coreV1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	applyconfigCoreV1 "k8s.io/client-go/applyconfigurations/core/v1"
	"strconv"
	"strings"
	"sysadm/k8sclient"
	"sysadm/objectsUI"
	"sysadm/utils"
)

func (s *service) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortServiceByName, "TD8": sortServiceByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "标签", "TD4": "类别", "TD5": "集群IP", "TD6": "外部IP", "TD7": "端口", "TD8": "创建时间"}
	additionalJs := []string{"js/sysadmfunctions.js", "/js/workloadList.js"}
	additionalCss := []string{}
	templateFile := "addWorkload.html"

	s.mainModuleName = "服务管理"
	s.moduleName = "服务"
	s.allPopMenuItems = allPopMenuItems
	s.allListItems = allListItems
	s.addButtonTile = "创建服务"
	s.isSearchForm = "no"
	s.allOrderFields = allOrderFields
	s.defaultOrderField = "TD1"
	s.defaultOrderDirection = "1"
	s.namespaced = true
	s.moduleID = "service"
	s.additionalJs = additionalJs
	s.additionalCss = additionalCss
	s.templateFile = templateFile

}

func (s *service) getMainModuleName() string {
	return s.mainModuleName
}

func (s *service) getModuleName() string {
	return s.moduleName
}

func (s *service) getAddButtonTitle() string {
	return s.addButtonTile
}

func (s *service) getIsSearchForm() string {
	return s.isSearchForm
}

func (s *service) getAllPopMenuItems() []string {
	return s.allPopMenuItems
}

func (s *service) getAllListItems() map[string]string {
	return s.allListItems
}

func (s *service) getDefaultOrderField() string {
	return s.defaultOrderField
}

func (s *service) getDefaultOrderDirection() string {
	return s.defaultOrderDirection
}

func (s *service) getAllorderFields() map[string]objectsUI.SortBy {
	return s.allOrderFields
}

func (s *service) getNamespaced() bool {
	return s.namespaced
}

// for daemonSet
func (s *service) listObjectData(selectedCluster, selectedNS string,
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

	serviceList, e := clientSet.CoreV1().Services(nsStr).List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return 0, dataList, e
	}

	totalNum := len(serviceList.Items)
	if totalNum < 1 {
		return 0, dataList, nil
	}

	orderfield := requestData["orderfield"]
	direction := requestData["direction"]
	if orderfield == "" {
		orderfield = s.getDefaultOrderField()
	}
	if direction == "" || (direction != "0" && direction != "1") {
		direction = s.getDefaultOrderDirection()
	}

	var serviceItems []interface{}
	for _, item := range serviceList.Items {
		serviceItems = append(serviceItems, item)
	}

	sortWorkloadData(serviceItems, direction, orderfield, s.getAllorderFields())

	endCount := totalNum
	if endCount > startPos+runData.pageInfo.NumPerPage {
		endCount = startPos + runData.pageInfo.NumPerPage
	}

	for i := startPos; i < endCount; i++ {
		interfaceData := serviceItems[i]
		serviceData, ok := interfaceData.(coreV1.Service)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not Service schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = serviceData.Name
		lineMap["TD1"] = serviceData.Name
		lineMap["TD2"] = serviceData.Namespace
		lineMap["TD3"] = objectsUI.ConvertMap2HTML(serviceData.Labels)
		lineMap["TD4"] = serviceData.Spec.Type
		lineMap["TD5"] = serviceData.Spec.ClusterIP
		externalIPs := "-"
		for _, ip := range serviceData.Spec.ExternalIPs {
			externalIPs = externalIPs + "<div>" + ip + "</div>"
		}
		lineMap["TD6"] = template.HTML(externalIPs)
		portStr := ""
		for _, v := range serviceData.Spec.Ports {
			nodePortStr := ""
			nodePort := int(v.NodePort)
			if nodePort > 0 && nodePort < 65535 {
				nodePortStr = strconv.Itoa(nodePort)
			} else {
				nodePortStr = "-"
			}
			portStr = portStr + "<div>" + string(v.Protocol) + ": " + strconv.Itoa(int(v.Port)) + ":" + nodePortStr + "</div>"
		}
		lineMap["TD7"] = template.HTML(portStr)
		lineMap["TD8"] = serviceData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		popmenuitems := "0,1"
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func (s *service) getModuleID() string {
	return s.moduleID
}

func (s *service) buildAddFormData(tplData map[string]interface{}) error {
	tplData["thirdCategory"] = "创建服务"
	formData, e := objectsUI.InitFormData("addService", "addService", "POST", "_self", "yes", "addWorkload", "")
	if e != nil {
		return e
	}
	tplData["formData"] = formData

	e = buildServiceBasiceFormData(tplData)
	if e != nil {
		return e
	}

	return nil
}

func (s *service) getAdditionalJs() []string {
	return s.additionalJs
}
func (s *service) getAdditionalCss() []string {
	return s.additionalCss
}

func (s *service) addNewResource(c *sysadmServer.Context, module string) error {
	requestKeys := []string{"dcid", "clusterID", "namespace", "addType", "nsSelected", "name", "serviceType", "portName[]", "protocol[]", "port[]", "targetPort[]", "nodePort[]"}
	requestKeys = append(requestKeys, "externalIPs", "selectorKey[]")
	requestKeys = append(requestKeys, "selectorValue[]")
	requestKeys = append(requestKeys, "labelKey[]", "labelValue[]", "annotationKey[]", "annotationValue[]")
	formData, e := utils.GetMultipartData(c, requestKeys)
	if e != nil {
		return e
	}
	ns := formData["nsSelected"].([]string)
	name := formData["name"].([]string)
	serviceApplyConfig := applyconfigCoreV1.Service(name[0], ns[0])

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
		labels[defaultServiceLabelKey] = name[0]
	}
	for k, v := range extraLabels {
		labels[k] = v
	}
	serviceApplyConfig = serviceApplyConfig.WithLabels(labels)

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
	serviceApplyConfig = serviceApplyConfig.WithAnnotations(annotations)

	serviceSpecApplyConfiguration := applyconfigCoreV1.ServiceSpecApplyConfiguration{}

	serviceTypeStr := strings.TrimSpace(formData["serviceType"].([]string)[0])
	serviceType := coreV1.ServiceTypeClusterIP
	switch coreV1.ServiceType(serviceTypeStr) {
	case coreV1.ServiceTypeClusterIP:
		serviceType = coreV1.ServiceTypeClusterIP
	case coreV1.ServiceTypeNodePort:
		serviceType = coreV1.ServiceTypeNodePort
	case coreV1.ServiceTypeLoadBalancer:
		serviceType = coreV1.ServiceTypeLoadBalancer
	case coreV1.ServiceTypeExternalName:
		serviceType = coreV1.ServiceTypeExternalName
	}
	serviceSpecApplyConfiguration.Type = &serviceType

	externalIP := strings.TrimSpace(formData["externalIPs"].([]string)[0])
	if externalIP != "" {
		externalIPs := strings.Split(externalIP, ";")
		serviceSpecApplyConfiguration.ExternalIPs = externalIPs
	}

	var applyPorts []applyconfigCoreV1.ServicePortApplyConfiguration
	portNames := formData["portName[]"].([]string)
	protocols := formData["protocol[]"].([]string)
	ports := formData["port[]"].([]string)
	targetPorts := formData["targetPort[]"].([]string)
	nodePorts := formData["nodePort[]"].([]string)
	for i, n := range portNames {
		servicePortApplyConfiguration := applyconfigCoreV1.ServicePortApplyConfiguration{}
		protocolStr := strings.TrimSpace(protocols[i])
		n = strings.TrimSpace(n)
		portStr := ports[i]
		if n != "" {
			portName := n
			servicePortApplyConfiguration.Name = &portName
		} else {
			portName := strings.ToLower(protocolStr) + "-" + portStr
			servicePortApplyConfiguration.Name = &portName
		}

		protocol := coreV1.ProtocolTCP
		switch protocolStr {
		case "TCP":
			protocol = coreV1.ProtocolTCP
		case "UDP":
			protocol = coreV1.ProtocolUDP
		case "SCTP":
			protocol = coreV1.ProtocolSCTP
		default:
			protocol = coreV1.ProtocolTCP
		}
		servicePortApplyConfiguration.Protocol = &protocol

		if strings.TrimSpace(portStr) != "" {
			portInt, e := strconv.Atoi(portStr)
			if e != nil {
				return e
			}
			portInt32 := int32(portInt)
			servicePortApplyConfiguration.Port = &portInt32
		}

		targetPortStr := targetPorts[i]
		if strings.TrimSpace(targetPortStr) != "" {
			targetPortInt, e := strconv.Atoi(targetPortStr)
			if e != nil {
				return e
			}
			targetPortInt32 := int32(targetPortInt)
			intOrStr := intstr.IntOrString{Type: intstr.Int, IntVal: targetPortInt32}
			servicePortApplyConfiguration.TargetPort = &intOrStr
		}

		if serviceType == coreV1.ServiceTypeNodePort {
			nodePortStr := nodePorts[i]
			if nodePortStr != "" {
				nodePortInt, e := strconv.Atoi(nodePortStr)
				if e != nil {
					return e
				}
				nodePortInt32 := int32(nodePortInt)
				servicePortApplyConfiguration.NodePort = &nodePortInt32
			}
		}

		applyPorts = append(applyPorts, servicePortApplyConfiguration)
	}
	serviceSpecApplyConfiguration.Ports = applyPorts

	selectorKeys := formData["selectorKey[]"].([]string)
	selectorValue := formData["selectorValue[]"].([]string)
	if len(selectorKeys) != len(selectorValue) {
		return fmt.Errorf("selector's key is not equal to selector's value")
	}
	selectors := make(map[string]string, 0)
	for i, k := range selectorKeys {
		value := selectorValue[i]
		selectors[k] = value
	}
	serviceSpecApplyConfiguration.Selector = selectors
	serviceApplyConfig = serviceApplyConfig.WithSpec(&serviceSpecApplyConfiguration)

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

	_, e = clientSet.CoreV1().Services(ns[0]).Apply(context.Background(), serviceApplyConfig, applyOption)

	return e

}

func (s *service) delResource(c *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (s *service) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

func (s *service) getTemplateFile(action string) string {
	switch action {
	case "list":
		return serviceTemplateFiles["list"]
	case "addform":
		return serviceTemplateFiles["addform"]
	default:
		return ""

	}
	return ""
}

func buildServiceBasiceFormData(tplData map[string]interface{}) error {
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
	_ = objectsUI.AddTextData("name", "name", "", "服务名称", "validateNewName", "addWorkloadValidateNewName", "长度小于63个字母数字或-且以字母数据开开头和结尾的字符串", 30, false, false, lineData)
	basicData = append(basicData, lineData)

	serviceTypeOptions := make(map[string]string, 0)
	serviceTypeOptions[string(coreV1.ServiceTypeClusterIP)] = "ClusterIP"
	serviceTypeOptions[string(coreV1.ServiceTypeNodePort)] = "NodePort"
	serviceTypeOptions[string(coreV1.ServiceTypeLoadBalancer)] = "LoadBalancer"
	serviceTypeOptions[string(coreV1.ServiceTypeExternalName)] = "ExternalName"
	lineData = objectsUI.InitLineData("serviceTypeLine", false, false, false)
	_ = objectsUI.AddSelectData("serviceType", "serviceType", string(coreV1.ServiceTypeClusterIP), "", "", "服务类型", "", 1, false, false, serviceTypeOptions, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("externalIPsLine", false, false, false)
	_ = objectsUI.AddTextData("externalIPs", "externalIPs", "", "外部IP地址", "", "", "可不填，多个时用;分隔，例IP1;IP2", 30, false, false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("portData", true, true, false)
	_ = objectsUI.AddTextData("portName", "portName[]", "", "端口名称", "", "", "可以不填", 10, false, false, lineData)
	portOptions := make(map[string]string, 0)
	portOptions["TCP"] = "TCP"
	portOptions["UDP"] = "UDP"
	portOptions["SCTP"] = "SCTP"
	_ = objectsUI.AddSelectData("protocol", "protocol[]", "TCP", "", "", "协议", "", 1, false, false, portOptions, lineData)
	_ = objectsUI.AddTextData("port", "port[]", "", "外部端口", "", "", "集群内不同应用互访问和LB,ClusterIP外部端口号", 10, false, false, lineData)
	_ = objectsUI.AddTextData("targetPort", "targetPort[]", "", "目标端口", "", "", "目标端口，即容器服务侦听的端口号", 10, false, false, lineData)
	_ = objectsUI.AddTextData("nodePort", "nodePort[]", "", "NodePort端口", "", "", "NodePort.当服务类型为NodePort时必填", 10, false, false, lineData)
	_ = objectsUI.AddWordsInputData("selectorLabel", "selectorLabel", "fa-trash", "#", "workloadDelSelector", false, true, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("portData", false, false, false)
	_ = objectsUI.AddWordsInputData("portData", "portData", "添加端口数据", "#", "workloadAddSelectorBlock", false, false, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("selectorLabel", true, true, false)
	_ = objectsUI.AddTextData("selectornKey", "selectorKey[]", "", "标签选择器", "", "", "", 30, false, false, lineData)
	_ = objectsUI.AddWordsInputData("equal", "equal", "=", "", "", false, false, lineData)
	_ = objectsUI.AddTextData("selectorValue", "selectorValue[]", "", "值", "", "", "", 30, false, false, lineData)
	_ = objectsUI.AddWordsInputData("selectorLabel", "selectorLabel", "fa-trash", "#", "workloadDelSelector", false, true, lineData)
	basicData = append(basicData, lineData)

	lineData = objectsUI.InitLineData("selectoranchor", false, false, false)
	_ = objectsUI.AddWordsInputData("selectorLabel", "selectorLabel", "添加选择器", "#", "workloadAddSelector", false, false, lineData)
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

func sortServiceByName(p, q interface{}) bool {
	pData, ok := p.(coreV1.Service)
	if !ok {
		return false
	}
	qData, ok := q.(coreV1.Service)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortServiceByCreatetime(p, q interface{}) bool {
	pData, ok := p.(coreV1.Service)
	if !ok {
		return false
	}
	qData, ok := q.(coreV1.Service)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

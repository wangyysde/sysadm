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
	"strconv"
	"sysadm/objectsUI"
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

	//TODO

	return nil
}

func (s *service) getAdditionalJs() []string {
	return s.additionalJs
}
func (s *service) getAdditionalCss() []string {
	return s.additionalCss
}

func (s *service) addNewResource(c *sysadmServer.Context, module string) error {
	// TODO

	return nil

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
	// TODO

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

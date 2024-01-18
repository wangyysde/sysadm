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
	"sort"
	"strconv"
	"sysadm/objectsUI"
)

func (s *sa) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortSaByName, "TD4": sortSaByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "关联的Secrets数", "TD4": "创建时间"}
	additionalJs := []string{}
	additionalCss := []string{}
	templateFile := ""

	s.mainModuleName = "帐号与角色"
	s.moduleName = "服务帐号"
	s.allPopMenuItems = allPopMenuItems
	s.allListItems = allListItems
	s.addButtonTile = "添加服务帐号"
	s.isSearchForm = "no"
	s.allOrderFields = allOrderFields
	s.defaultOrderField = "TD1"
	s.defaultOrderDirection = "1"
	s.namespaced = true
	s.moduleID = "serviceaccount"
	s.additionalJs = additionalJs
	s.additionalCss = additionalCss
	s.templateFile = templateFile
}

func (s *sa) getMainModuleName() string {
	return s.mainModuleName
}

func (s *sa) getModuleName() string {
	return s.moduleName
}

func (s *sa) getAddButtonTitle() string {
	return s.addButtonTile
}

func (s *sa) getIsSearchForm() string {
	return s.isSearchForm
}

func (s *sa) getAllPopMenuItems() []string {
	return s.allPopMenuItems
}

func (s *sa) getAllListItems() map[string]string {
	return s.allListItems
}

func (s *sa) getDefaultOrderField() string {
	return s.defaultOrderField
}

func (s *sa) getDefaultOrderDirection() string {
	return s.defaultOrderDirection
}

func (s *sa) getAllorderFields() map[string]objectsUI.SortBy {
	return s.allOrderFields
}

func (s *sa) getNamespaced() bool {
	return s.namespaced
}

// for Ingress
func (s *sa) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	clientSet, e := createClientSet(selectedCluster)
	if clientSet == nil {
		return 0, dataList, e
	}

	nsStr, orderField, direction := checkRequestData(selectedNS, s.defaultOrderField, s.defaultOrderDirection, requestData)
	objList, e := clientSet.CoreV1().ServiceAccounts(nsStr).List(context.Background(), metav1.ListOptions{})
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

	moduleAllOrderFields := s.allOrderFields
	for field, fn := range moduleAllOrderFields {
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
		objData, ok := interfaceData.(corev1.ServiceAccount)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not SA schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = objData.Name
		lineMap["TD1"] = objData.Name
		lineMap["TD2"] = objData.Namespace
		secretsNum := len(objData.Secrets)
		imagePullSecretNum := len(objData.ImagePullSecrets)
		totalSecretNum := secretsNum + imagePullSecretNum
		lineMap["TD3"] = strconv.Itoa(totalSecretNum)
		lineMap["TD4"] = objData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		lineMap["popmenuitems"] = "0,1"
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortSaByName(p, q interface{}) bool {
	pData, ok := p.(corev1.ServiceAccount)
	if !ok {
		return false
	}
	qData, ok := q.(corev1.ServiceAccount)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortSaByCreatetime(p, q interface{}) bool {
	pData, ok := p.(corev1.ServiceAccount)
	if !ok {
		return false
	}
	qData, ok := q.(corev1.ServiceAccount)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

func (s *sa) getModuleID() string {
	return s.moduleID
}

func (s *sa) buildAddFormData(tplData map[string]interface{}) error {
	tplData["thirdCategory"] = "创建服务帐号"
	formData, e := objectsUI.InitFormData("addSa", "addSa", "POST", "_self", "yes", "addWorkload", "")
	if e != nil {
		return e
	}
	tplData["formData"] = formData

	// TODO
	return nil
}

func (s *sa) getAdditionalJs() []string {
	return s.additionalJs
}
func (s *sa) getAdditionalCss() []string {
	return s.additionalCss
}

func (s *sa) addNewResource(c *sysadmServer.Context, module string) error {
	// TODO

	return nil
}

func (s *sa) delResource(c *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (s *sa) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

func (s *sa) getTemplateFile(action string) string {
	switch action {
	case "list":
		return saTemplateFiles["list"]
	case "addform":
		return saTemplateFiles["addform"]
	default:
		return ""
	}
	return ""
}

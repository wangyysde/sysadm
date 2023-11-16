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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
	"sysadm/objectsUI"
)

func (s *secret) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortSecretByName, "TD6": sortSecretByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "类型", "TD4": "标签", "TD5": "数据项数", "TD6": "是否可编辑", "TD7": "创建时间"}

	s.mainModuleName = "配置和存储"
	s.moduleName = "Secrets"
	s.allPopMenuItems = allPopMenuItems
	s.allListItems = allListItems
	s.addButtonTile = ""
	s.isSearchForm = "no"
	s.allOrderFields = allOrderFields
	s.defaultOrderField = "TD1"
	s.defaultOrderDirection = "1"
}

func (s *secret) getMainModuleName() string {
	return s.mainModuleName
}

func (s *secret) getModuleName() string {
	return s.moduleName
}

func (s *secret) getAddButtonTitle() string {
	return s.addButtonTile
}

func (s *secret) getIsSearchForm() string {
	return s.isSearchForm
}

func (s *secret) getAllPopMenuItems() []string {
	return s.allPopMenuItems
}

func (s *secret) getAllListItems() map[string]string {
	return s.allListItems
}

func (s *secret) getDefaultOrderField() string {
	return s.defaultOrderField
}

func (s *secret) getDefaultOrderDirection() string {
	return s.defaultOrderDirection
}

func (s *secret) getAllorderFields() map[string]objectsUI.SortBy {
	return s.allOrderFields
}

func (s *secret) getNamespaced() bool {
	return s.namespaced
}

// for Ingress
func (s *secret) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	clientSet, e := createClientSet(selectedCluster)
	if clientSet == nil {
		return 0, dataList, e
	}

	nsStr, orderField, direction := checkRequestData(selectedNS, s.defaultOrderField, s.defaultOrderDirection, requestData)
	secretList, e := clientSet.CoreV1().Secrets(nsStr).List(context.Background(), metav1.ListOptions{})
	if e != nil {
		return 0, dataList, e
	}
	totalNum := len(secretList.Items)
	if totalNum < 1 {
		return 0, dataList, nil
	}
	endCount := totalNum
	if endCount > startPos+runData.pageInfo.NumPerPage {
		endCount = startPos + runData.pageInfo.NumPerPage
	}

	var secretItems []interface{}
	for _, item := range secretList.Items {
		secretItems = append(secretItems, item)
	}

	moduleAllOrderFields := s.allOrderFields
	for field, fn := range moduleAllOrderFields {
		if field == orderField {
			if direction == "1" {
				sort.Sort(objectsUI.SortData{Data: secretItems, By: fn})
			} else {
				sort.Sort(sort.Reverse(objectsUI.SortData{Data: secretItems, By: fn}))
			}

		}
	}

	for i := startPos; i < endCount; i++ {
		interfaceData := secretItems[i]
		secretData, ok := interfaceData.(corev1.Secret)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not Secret schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = secretData.Name
		lineMap["TD1"] = secretData.Name
		lineMap["TD2"] = secretData.Namespace
		lineMap["TD3"] = secretData.Type
		lineMap["TD4"] = objectsUI.ConvertMap2HTML(secretData.Labels)
		dataCount := len(secretData.Data)
		stringDataCount := len(secretData.StringData)
		totalNum := dataCount + stringDataCount
		lineMap["TD5"] = totalNum
		editable := "是"
		popmenuitems := "0,1"
		if secretData.Immutable != nil && *secretData.Immutable {
			editable = "否"
			popmenuitems = "0,1"
		}
		lineMap["TD6"] = editable
		lineMap["TD7"] = secretData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortSecretByName(p, q interface{}) bool {
	pData, ok := p.(corev1.Secret)
	if !ok {
		return false
	}
	qData, ok := q.(corev1.Secret)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortSecretByCreatetime(p, q interface{}) bool {
	pData, ok := p.(corev1.Secret)
	if !ok {
		return false
	}
	qData, ok := q.(corev1.Secret)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

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
	"strconv"
	"sysadm/objectsUI"
)

func (c *configmap) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortCmByName, "TD6": sortCmByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "标签", "TD4": "数据项数", "TD5": "是否可修改", "TD6": "创建时间"}

	c.mainModuleName = "配置和存储"
	c.moduleName = "configmap"
	c.allPopMenuItems = allPopMenuItems
	c.allListItems = allListItems
	c.addButtonTile = ""
	c.isSearchForm = "no"
	c.allOrderFields = allOrderFields
	c.defaultOrderField = "TD1"
	c.defaultOrderDirection = "1"
	c.namespaced = true
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

// for Ingress
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

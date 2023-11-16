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
	"html/template"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
	"sysadm/objectsUI"
)

func (i *ingress) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortIngressByName, "TD6": sortIngressByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "Hosts", "TD4": "标签", "TD5": "IngressClassName", "TD6": "创建时间"}

	i.mainModuleName = "服务管理"
	i.moduleName = "ingress"
	i.allPopMenuItems = allPopMenuItems
	i.allListItems = allListItems
	i.addButtonTile = ""
	i.isSearchForm = "no"
	i.allOrderFields = allOrderFields
	i.defaultOrderField = "TD1"
	i.defaultOrderDirection = "1"
	i.namespaced = true
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

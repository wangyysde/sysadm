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
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
	"sysadm/objectsUI"
)

func (i *ingressclass) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortIngressClassByName, "TD3": sortIngressClassByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD3": "创建时间"}
	additionalJs := []string{}
	additionalCss := []string{}

	i.mainModuleName = "服务管理"
	i.moduleName = "Ingress Classes"
	i.allPopMenuItems = allPopMenuItems
	i.allListItems = allListItems
	i.addButtonTile = ""
	i.isSearchForm = "no"
	i.allOrderFields = allOrderFields
	i.defaultOrderField = "TD1"
	i.defaultOrderDirection = "1"
	i.namespaced = false
	i.moduleID = "ingressclass"
	i.additionalJs = additionalJs
	i.additionalCss = additionalCss
}

func (i *ingressclass) getMainModuleName() string {
	return i.mainModuleName
}

func (i *ingressclass) getModuleName() string {
	return i.moduleName
}

func (i *ingressclass) getAddButtonTitle() string {
	return i.addButtonTile
}

func (i *ingressclass) getIsSearchForm() string {
	return i.isSearchForm
}

func (i *ingressclass) getAllPopMenuItems() []string {
	return i.allPopMenuItems
}

func (i *ingressclass) getAllListItems() map[string]string {
	return i.allListItems
}

func (i *ingressclass) getDefaultOrderField() string {
	return i.defaultOrderField
}

func (i *ingressclass) getDefaultOrderDirection() string {
	return i.defaultOrderDirection
}

func (i *ingressclass) getAllorderFields() map[string]objectsUI.SortBy {
	return i.allOrderFields
}

func (i *ingressclass) getNamespaced() bool {
	return i.namespaced
}

// for ingressclass
func (i *ingressclass) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	clientSet, e := createClientSet(selectedCluster)
	if clientSet == nil {
		return 0, dataList, e
	}

	_, orderField, direction := checkRequestData(selectedNS, i.defaultOrderField, i.defaultOrderDirection, requestData)
	objList, e := clientSet.NetworkingV1().IngressClasses().List(context.Background(), metav1.ListOptions{})
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

	moduleAllOrderFields := i.allOrderFields
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
		objData, ok := interfaceData.(networkingv1.IngressClass)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not IngressClass schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = objData.Name
		lineMap["TD1"] = objData.Name
		lineMap["TD2"] = objData.Spec.Controller
		lineMap["TD3"] = objData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		popmenuitems := "0,1"
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortIngressClassByName(p, q interface{}) bool {
	pData, ok := p.(networkingv1.IngressClass)
	if !ok {
		return false
	}
	qData, ok := q.(networkingv1.IngressClass)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortIngressClassByCreatetime(p, q interface{}) bool {
	pData, ok := p.(networkingv1.IngressClass)
	if !ok {
		return false
	}
	qData, ok := q.(networkingv1.IngressClass)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

func (i *ingressclass) getModuleID() string {
	return i.moduleID
}

func (i *ingressclass) buildAddFormData(tplData map[string]interface{}) error {
	// TODO
	return nil
}

func (i *ingressclass) getAdditionalJs() []string {
	return i.additionalJs
}
func (i *ingressclass) getAdditionalCss() []string {
	return i.additionalCss
}

func (i *ingressclass) addNewResource(c *sysadmServer.Context, module string) error {
	// TODO

	return nil
}

func (i *ingressclass) delResource(c *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (i *ingressclass) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

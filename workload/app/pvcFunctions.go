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
	"sysadm/objectsUI"
)

func (p *pvc) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortPvcByName, "TD6": sortPvcByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "状态", "TD4": "VolumeName", "TD5": "容量", "TD6": "访问模式", "TD7": "StorageClassName", "TD8": "创建时间"}
	additionalJs := []string{}
	additionalCss := []string{}

	p.mainModuleName = "配置和存储"
	p.moduleName = "pvc"
	p.allPopMenuItems = allPopMenuItems
	p.allListItems = allListItems
	p.addButtonTile = ""
	p.isSearchForm = "no"
	p.allOrderFields = allOrderFields
	p.defaultOrderField = "TD1"
	p.defaultOrderDirection = "1"
	p.namespaced = true
	p.moduleID = "pvc"
	p.additionalJs = additionalJs
	p.additionalCss = additionalCss
}

func (p *pvc) getMainModuleName() string {
	return p.mainModuleName
}

func (p *pvc) getModuleName() string {
	return p.moduleName
}

func (p *pvc) getAddButtonTitle() string {
	return p.addButtonTile
}

func (p *pvc) getIsSearchForm() string {
	return p.isSearchForm
}

func (p *pvc) getAllPopMenuItems() []string {
	return p.allPopMenuItems
}

func (p *pvc) getAllListItems() map[string]string {
	return p.allListItems
}

func (p *pvc) getDefaultOrderField() string {
	return p.defaultOrderField
}

func (p *pvc) getDefaultOrderDirection() string {
	return p.defaultOrderDirection
}

func (p *pvc) getAllorderFields() map[string]objectsUI.SortBy {
	return p.allOrderFields
}

func (p *pvc) getNamespaced() bool {
	return p.namespaced
}

// for Ingress
func (p *pvc) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	clientSet, e := createClientSet(selectedCluster)
	if clientSet == nil {
		return 0, dataList, e
	}

	nsStr, orderField, direction := checkRequestData(selectedNS, p.defaultOrderField, p.defaultOrderDirection, requestData)
	objList, e := clientSet.CoreV1().PersistentVolumeClaims(nsStr).List(context.Background(), metav1.ListOptions{})
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

	moduleAllOrderFields := p.allOrderFields
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
		objData, ok := interfaceData.(corev1.PersistentVolumeClaim)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not Secret schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = objData.Name
		lineMap["TD1"] = objData.Name
		lineMap["TD2"] = objData.Namespace
		lineMap["TD3"] = objData.Status.Phase
		lineMap["TD4"] = objData.Spec.VolumeName
		lineMap["TD5"] = objData.Spec.Resources.Requests.Storage().String()
		accessMode := ""
		for _, mode := range objData.Spec.AccessModes {
			switch mode {
			case corev1.ReadWriteOnce:
				if accessMode == "" {
					accessMode = "RWO"
				} else {
					accessMode = accessMode + "," + "RWO"
				}
			case corev1.ReadOnlyMany:
				if accessMode == "" {
					accessMode = "ROM"
				} else {
					accessMode = accessMode + "," + "ROM"
				}
			case corev1.ReadWriteMany:
				if accessMode == "" {
					accessMode = "RWM"
				} else {
					accessMode = accessMode + "," + "RWM"
				}
			case corev1.ReadWriteOncePod:
				if accessMode == "" {
					accessMode = "RWOP"
				} else {
					accessMode = accessMode + "," + "RWOP"
				}
			}
		}
		lineMap["TD6"] = accessMode
		lineMap["TD7"] = objData.Spec.StorageClassName
		lineMap["TD8"] = objData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		popmenuitems := ""
		if objData.Status.Phase != corev1.ClaimBound {
			popmenuitems = "0,1"
		}
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortPvcByName(p, q interface{}) bool {
	pData, ok := p.(corev1.PersistentVolumeClaim)
	if !ok {
		return false
	}
	qData, ok := q.(corev1.PersistentVolumeClaim)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortPvcByCreatetime(p, q interface{}) bool {
	pData, ok := p.(corev1.PersistentVolumeClaim)
	if !ok {
		return false
	}
	qData, ok := q.(corev1.PersistentVolumeClaim)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

func (p *pvc) getModuleID() string {
	return p.moduleID
}

func (p *pvc) buildAddFormData(tplData map[string]interface{}) error {
	// TODO
	return nil
}

func (p *pvc) getAdditionalJs() []string {
	return p.additionalJs
}
func (p *pvc) getAdditionalCss() []string {
	return p.additionalCss
}

func (p *pvc) addNewResource(c *sysadmServer.Context, module string) error {
	// TODO

	return nil
}

func (p *pvc) delResource(c *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (p *pvc) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

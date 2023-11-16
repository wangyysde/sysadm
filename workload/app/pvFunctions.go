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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
	"sysadm/objectsUI"
)

func (p *pv) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortPvByName, "TD9": sortPvByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "容量", "TD3": "访问模式", "TD4": "回收策略", "TD5": "状态", "TD6": "绑定对象", "TD7": "Storage Class", "TD8": "卷模式", "TD9": "创建时间"}

	p.mainModuleName = "配置和存储"
	p.moduleName = "Persistent Volumes"
	p.allPopMenuItems = allPopMenuItems
	p.allListItems = allListItems
	p.addButtonTile = ""
	p.isSearchForm = "no"
	p.allOrderFields = allOrderFields
	p.defaultOrderField = "TD1"
	p.defaultOrderDirection = "1"
	p.namespaced = false
}

func (p *pv) getMainModuleName() string {
	return p.mainModuleName
}

func (p *pv) getModuleName() string {
	return p.moduleName
}

func (p *pv) getAddButtonTitle() string {
	return p.addButtonTile
}

func (p *pv) getIsSearchForm() string {
	return p.isSearchForm
}

func (p *pv) getAllPopMenuItems() []string {
	return p.allPopMenuItems
}

func (p *pv) getAllListItems() map[string]string {
	return p.allListItems
}

func (p *pv) getDefaultOrderField() string {
	return p.defaultOrderField
}

func (p *pv) getDefaultOrderDirection() string {
	return p.defaultOrderDirection
}

func (p *pv) getAllorderFields() map[string]objectsUI.SortBy {
	return p.allOrderFields
}

func (p *pv) getNamespaced() bool {
	return p.namespaced
}

// for ingressclass
func (p *pv) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	clientSet, e := createClientSet(selectedCluster)
	if clientSet == nil {
		return 0, dataList, e
	}

	_, orderField, direction := checkRequestData(selectedNS, p.defaultOrderField, p.defaultOrderDirection, requestData)
	objList, e := clientSet.CoreV1().PersistentVolumes().List(context.Background(), metav1.ListOptions{})
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
		objData, ok := interfaceData.(corev1.PersistentVolume)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not PersistentVolumes schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = objData.Name
		lineMap["TD1"] = objData.Name
		capacity := ""
		for k, v := range objData.Spec.Capacity {
			capacity = capacity + "<div>" + string(k) + ": " + v.String() + "</div>"
		}
		lineMap["TD2"] = template.HTML(capacity)
		accessMode := ""
		for _, v := range objData.Spec.AccessModes {
			if accessMode != "" {
				accessMode = accessMode + ","
			}
			switch v {
			case corev1.ReadWriteOnce:
				accessMode = accessMode + "RWO"
			case corev1.ReadOnlyMany:
				accessMode = accessMode + "ROM"
			case corev1.ReadWriteMany:
				accessMode = accessMode + "RWM"
			case corev1.ReadWriteOncePod:
				accessMode = accessMode + "RWOP"
			}
		}
		lineMap["TD3"] = accessMode
		lineMap["TD4"] = objData.Spec.PersistentVolumeReclaimPolicy
		status := objData.Status.Phase
		lineMap["TD5"] = status
		boundPvc := objData.Spec.ClaimRef.Namespace + "/" + objData.Spec.ClaimRef.Name
		if boundPvc == "/" {
			boundPvc = "-"
		}
		lineMap["TD6"] = boundPvc
		lineMap["TD7"] = objData.Spec.StorageClassName
		lineMap["TD8"] = objData.Spec.VolumeMode
		lineMap["TD9"] = objData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		popmenuitems := ""
		if status != corev1.VolumeBound {
			popmenuitems = "0,1"
		}

		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortPvByName(p, q interface{}) bool {
	pData, ok := p.(corev1.PersistentVolume)
	if !ok {
		return false
	}
	qData, ok := q.(corev1.PersistentVolume)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortPvByCreatetime(p, q interface{}) bool {
	pData, ok := p.(corev1.PersistentVolume)
	if !ok {
		return false
	}
	qData, ok := q.(corev1.PersistentVolume)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

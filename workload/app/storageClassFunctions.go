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
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
	"sysadm/objectsUI"
)

func (s *storageclass) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortStorageClassByName, "TD7": sortStorageClassByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "提供者", "TD3": "参数", "TD4": "回收策略", "TD5": "绑定模式", "TD6": "是否允许扩容", "TD7": "创建时间"}
	additionalJs := []string{}
	additionalCss := []string{}

	s.mainModuleName = "配置和存储"
	s.moduleName = "Storage Classes"
	s.allPopMenuItems = allPopMenuItems
	s.allListItems = allListItems
	s.addButtonTile = ""
	s.isSearchForm = "no"
	s.allOrderFields = allOrderFields
	s.defaultOrderField = "TD1"
	s.defaultOrderDirection = "1"
	s.namespaced = false
	s.moduleID = "storageclass"
	s.additionalJs = additionalJs
	s.additionalCss = additionalCss
}

func (s *storageclass) getMainModuleName() string {
	return s.mainModuleName
}

func (s *storageclass) getModuleName() string {
	return s.moduleName
}

func (s *storageclass) getAddButtonTitle() string {
	return s.addButtonTile
}

func (s *storageclass) getIsSearchForm() string {
	return s.isSearchForm
}

func (s *storageclass) getAllPopMenuItems() []string {
	return s.allPopMenuItems
}

func (s *storageclass) getAllListItems() map[string]string {
	return s.allListItems
}

func (s *storageclass) getDefaultOrderField() string {
	return s.defaultOrderField
}

func (s *storageclass) getDefaultOrderDirection() string {
	return s.defaultOrderDirection
}

func (s *storageclass) getAllorderFields() map[string]objectsUI.SortBy {
	return s.allOrderFields
}

func (s *storageclass) getNamespaced() bool {
	return s.namespaced
}

// for ingressclass
func (s *storageclass) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	clientSet, e := createClientSet(selectedCluster)
	if clientSet == nil {
		return 0, dataList, e
	}

	_, orderField, direction := checkRequestData(selectedNS, s.defaultOrderField, s.defaultOrderDirection, requestData)
	objList, e := clientSet.StorageV1().StorageClasses().List(context.Background(), metav1.ListOptions{})
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
		objData, ok := interfaceData.(storagev1.StorageClass)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not StorageClass schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = objData.Name
		lineMap["TD1"] = objData.Name
		lineMap["TD2"] = objData.Provisioner
		params := ""
		for k, v := range objData.Parameters {
			params = params + "<div> Key: " + k + " value: " + v + "</div>"
		}
		if params == "" {
			params = "-"
		}
		lineMap["TD3"] = template.HTML(params)
		lineMap["TD4"] = objData.ReclaimPolicy
		lineMap["TD5"] = objData.VolumeBindingMode
		allowVolumeExpansion := "否"
		if objData.AllowVolumeExpansion != nil && *objData.AllowVolumeExpansion {
			allowVolumeExpansion = "是"
		}
		lineMap["TD6"] = allowVolumeExpansion
		lineMap["TD7"] = objData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		popmenuitems := "0,1"
		lineMap["popmenuitems"] = popmenuitems
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortStorageClassByName(p, q interface{}) bool {
	pData, ok := p.(storagev1.StorageClass)
	if !ok {
		return false
	}
	qData, ok := q.(storagev1.StorageClass)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortStorageClassByCreatetime(p, q interface{}) bool {
	pData, ok := p.(storagev1.StorageClass)
	if !ok {
		return false
	}
	qData, ok := q.(storagev1.StorageClass)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

func (s *storageclass) getModuleID() string {
	return s.moduleID
}

func (s *storageclass) buildAddFormData(tplData map[string]interface{}) error {
	// TODO
	return nil
}

func (s *storageclass) getAdditionalJs() []string {
	return s.additionalJs
}
func (s *storageclass) getAdditionalCss() []string {
	return s.additionalCss
}

func (s *storageclass) addNewResource(c *sysadmServer.Context, module string) error {
	// TODO

	return nil
}

func (s *storageclass) delResource(c *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (s *storageclass) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

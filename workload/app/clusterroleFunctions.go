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
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
	"sysadm/objectsUI"
)

func (c *clusterrole) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortClusterroleByName, "TD2": sortClusterroleByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "创建时间"}
	additionalJs := []string{}
	additionalCss := []string{}

	c.mainModuleName = "帐号与角色"
	c.moduleName = "集群角色"
	c.allPopMenuItems = allPopMenuItems
	c.allListItems = allListItems
	c.addButtonTile = ""
	c.isSearchForm = "no"
	c.allOrderFields = allOrderFields
	c.defaultOrderField = "TD1"
	c.defaultOrderDirection = "1"
	c.namespaced = false
	c.moduleID = "clusterrole"
	c.additionalJs = additionalJs
	c.additionalCss = additionalCss

}

func (c *clusterrole) getMainModuleName() string {
	return c.mainModuleName
}

func (c *clusterrole) getModuleName() string {
	return c.moduleName
}

func (c *clusterrole) getAddButtonTitle() string {
	return c.addButtonTile
}

func (c *clusterrole) getIsSearchForm() string {
	return c.isSearchForm
}

func (c *clusterrole) getAllPopMenuItems() []string {
	return c.allPopMenuItems
}

func (c *clusterrole) getAllListItems() map[string]string {
	return c.allListItems
}

func (c *clusterrole) getDefaultOrderField() string {
	return c.defaultOrderField
}

func (c *clusterrole) getDefaultOrderDirection() string {
	return c.defaultOrderDirection
}

func (c *clusterrole) getAllorderFields() map[string]objectsUI.SortBy {
	return c.allOrderFields
}

func (c *clusterrole) getNamespaced() bool {
	return c.namespaced
}

// for ingressclass
func (c *clusterrole) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	clientSet, e := createClientSet(selectedCluster)
	if clientSet == nil {
		return 0, dataList, e
	}

	_, orderField, direction := checkRequestData(selectedNS, c.defaultOrderField, c.defaultOrderDirection, requestData)
	objList, e := clientSet.RbacV1().ClusterRoles().List(context.Background(), metav1.ListOptions{})
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

	moduleAllOrderFields := c.allOrderFields
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
		objData, ok := interfaceData.(rbacv1.ClusterRole)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not ClusterRole schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = objData.Name
		lineMap["TD1"] = objData.Name
		lineMap["TD2"] = objData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		lineMap["popmenuitems"] = "0,1"
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortClusterroleByName(p, q interface{}) bool {
	pData, ok := p.(rbacv1.ClusterRole)
	if !ok {
		return false
	}
	qData, ok := q.(rbacv1.ClusterRole)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortClusterroleByCreatetime(p, q interface{}) bool {
	pData, ok := p.(rbacv1.ClusterRole)
	if !ok {
		return false
	}
	qData, ok := q.(rbacv1.ClusterRole)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

func (c *clusterrole) getModuleID() string {
	return c.moduleID
}

func (c *clusterrole) buildAddFormData(tplData map[string]interface{}) error {
	// TODO
	return nil
}

func (c *clusterrole) getAdditionalJs() []string {
	return c.additionalJs
}
func (c *clusterrole) getAdditionalCss() []string {
	return c.additionalCss
}

func (c *clusterrole) addNewResource(s *sysadmServer.Context, module string) error {
	// TODO

	return nil
}

func (c *clusterrole) delResource(s *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (c *clusterrole) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

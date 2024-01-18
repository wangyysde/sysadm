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

func (r *role) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortRoleByName, "TD3": sortRoleByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "命名空间", "TD3": "创建时间"}
	additionalJs := []string{}
	additionalCss := []string{}
	templateFile := ""

	r.mainModuleName = "帐号与角色"
	r.moduleName = "角色"
	r.allPopMenuItems = allPopMenuItems
	r.allListItems = allListItems
	r.addButtonTile = "创建角色"
	r.isSearchForm = "no"
	r.allOrderFields = allOrderFields
	r.defaultOrderField = "TD1"
	r.defaultOrderDirection = "1"
	r.namespaced = true
	r.moduleID = "role"
	r.additionalJs = additionalJs
	r.additionalCss = additionalCss
	r.templateFile = templateFile
}

func (r *role) getMainModuleName() string {
	return r.mainModuleName
}

func (r *role) getModuleName() string {
	return r.moduleName
}

func (r *role) getAddButtonTitle() string {
	return r.addButtonTile
}

func (r *role) getIsSearchForm() string {
	return r.isSearchForm
}

func (r *role) getAllPopMenuItems() []string {
	return r.allPopMenuItems
}

func (r *role) getAllListItems() map[string]string {
	return r.allListItems
}

func (r *role) getDefaultOrderField() string {
	return r.defaultOrderField
}

func (r *role) getDefaultOrderDirection() string {
	return r.defaultOrderDirection
}

func (r *role) getAllorderFields() map[string]objectsUI.SortBy {
	return r.allOrderFields
}

func (r *role) getNamespaced() bool {
	return r.namespaced
}

// for Ingress
func (r *role) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	clientSet, e := createClientSet(selectedCluster)
	if clientSet == nil {
		return 0, dataList, e
	}

	nsStr, orderField, direction := checkRequestData(selectedNS, r.defaultOrderField, r.defaultOrderDirection, requestData)
	objList, e := clientSet.RbacV1().Roles(nsStr).List(context.Background(), metav1.ListOptions{})
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

	moduleAllOrderFields := r.allOrderFields
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
		objData, ok := interfaceData.(rbacv1.Role)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not RBAC schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = objData.Name
		lineMap["TD1"] = objData.Name
		lineMap["TD2"] = objData.Namespace
		lineMap["TD3"] = objData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		lineMap["popmenuitems"] = "0,1"
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortRoleByName(p, q interface{}) bool {
	pData, ok := p.(rbacv1.Role)
	if !ok {
		return false
	}
	qData, ok := q.(rbacv1.Role)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortRoleByCreatetime(p, q interface{}) bool {
	pData, ok := p.(rbacv1.Role)
	if !ok {
		return false
	}
	qData, ok := q.(rbacv1.Role)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

func (r *role) getModuleID() string {
	return r.moduleID
}

func (r *role) buildAddFormData(tplData map[string]interface{}) error {
	tplData["thirdCategory"] = "创建角色"
	formData, e := objectsUI.InitFormData("addRole", "addRole", "POST", "_self", "yes", "addWorkload", "")
	if e != nil {
		return e
	}
	tplData["formData"] = formData

	// TODO
	return nil
}

func (r *role) getAdditionalJs() []string {
	return r.additionalJs
}
func (r *role) getAdditionalCss() []string {
	return r.additionalCss
}

func (r *role) addNewResource(c *sysadmServer.Context, module string) error {
	// TODO

	return nil
}

func (r *role) delResource(c *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (r *role) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

func (r *role) getTemplateFile(action string) string {
	switch action {
	case "list":
		return roleTemplateFiles["list"]
	case "addform":
		return roleTemplateFiles["addform"]
	default:
		return ""
	}
	return ""

}

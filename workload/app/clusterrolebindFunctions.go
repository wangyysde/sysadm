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
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sort"
	"sysadm/objectsUI"
)

func (c *clusterrolebind) setObjectInfo() {
	allOrderFields := map[string]objectsUI.SortBy{"TD1": sortClusterrolebindByName, "TD4": sortClusterrolebindByCreatetime}
	allPopMenuItems := []string{"编辑,edit,GET,page", "删除,del,POST,tip"}
	allListItems := map[string]string{"TD1": "名称", "TD2": "关联的角色", "TD3": "关联的对象", "TD4": "创建时间"}
	additionalJs := []string{}
	additionCss := []string{}
	templateFile := ""

	c.mainModuleName = "帐号与角色"
	c.moduleName = "集群角色绑定"
	c.allPopMenuItems = allPopMenuItems
	c.allListItems = allListItems
	c.addButtonTile = "添加集群角色绑定"
	c.isSearchForm = "no"
	c.allOrderFields = allOrderFields
	c.defaultOrderField = "TD1"
	c.defaultOrderDirection = "1"
	c.namespaced = false
	c.moduleID = "clusterrolebind"
	c.additionalJs = additionalJs
	c.additionalCss = additionCss
	c.templateFile = templateFile
}

func (c *clusterrolebind) getMainModuleName() string {
	return c.mainModuleName
}

func (c *clusterrolebind) getModuleName() string {
	return c.moduleName
}

func (c *clusterrolebind) getAddButtonTitle() string {
	return c.addButtonTile
}

func (c *clusterrolebind) getIsSearchForm() string {
	return c.isSearchForm
}

func (c *clusterrolebind) getAllPopMenuItems() []string {
	return c.allPopMenuItems
}

func (c *clusterrolebind) getAllListItems() map[string]string {
	return c.allListItems
}

func (c *clusterrolebind) getDefaultOrderField() string {
	return c.defaultOrderField
}

func (c *clusterrolebind) getDefaultOrderDirection() string {
	return c.defaultOrderDirection
}

func (c *clusterrolebind) getAllorderFields() map[string]objectsUI.SortBy {
	return c.allOrderFields
}

func (c *clusterrolebind) getNamespaced() bool {
	return c.namespaced
}

// for Ingress
func (c *clusterrolebind) listObjectData(selectedCluster, selectedNS string,
	startPos int, requestData map[string]string) (int, []map[string]interface{}, error) {
	var dataList []map[string]interface{}

	clientSet, e := createClientSet(selectedCluster)
	if clientSet == nil {
		return 0, dataList, e
	}

	_, orderField, direction := checkRequestData(selectedNS, c.defaultOrderField, c.defaultOrderDirection, requestData)
	objList, e := clientSet.RbacV1().ClusterRoleBindings().List(context.Background(), metav1.ListOptions{})
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
		objData, ok := interfaceData.(rbacv1.ClusterRoleBinding)
		if !ok {
			return 0, dataList, fmt.Errorf("the data is not ClusterRoleBinding schema")
		}
		lineMap := make(map[string]interface{}, 0)
		lineMap["objectID"] = objData.Name
		lineMap["TD1"] = objData.Name
		roleRef := objData.RoleRef.Kind + "/" + objData.RoleRef.Name
		lineMap["TD2"] = roleRef
		objRef := ""
		for _, obj := range objData.Subjects {
			objRef = objRef + "<div>" + obj.Kind + "/" + obj.Name + "</div>"
		}
		lineMap["TD3"] = template.HTML(objRef)
		lineMap["TD4"] = objData.CreationTimestamp.Format(objectsUI.DefaultTimeStampFormat)
		lineMap["popmenuitems"] = "0,1"
		dataList = append(dataList, lineMap)
	}

	return totalNum, dataList, nil
}

func sortClusterrolebindByName(p, q interface{}) bool {
	pData, ok := p.(rbacv1.ClusterRoleBinding)
	if !ok {
		return false
	}
	qData, ok := q.(rbacv1.ClusterRoleBinding)
	if !ok {
		return false
	}

	return pData.Name < qData.Name
}

func sortClusterrolebindByCreatetime(p, q interface{}) bool {
	pData, ok := p.(rbacv1.ClusterRoleBinding)
	if !ok {
		return false
	}
	qData, ok := q.(rbacv1.ClusterRoleBinding)
	if !ok {
		return false
	}

	return pData.CreationTimestamp.String() < qData.CreationTimestamp.String()
}

func (c *clusterrolebind) getModuleID() string {
	return c.moduleID
}

func (c *clusterrolebind) buildAddFormData(tplData map[string]interface{}) error {
	tplData["thirdCategory"] = "创建集群角色绑定"
	formData, e := objectsUI.InitFormData("addClusterRoleBinding", "addClusterRoleBinding", "POST", "_self", "yes", "addWorkload", "")
	if e != nil {
		return e
	}
	tplData["formData"] = formData

	//TODO
	return nil
}

func (c *clusterrolebind) getAdditionalJs() []string {
	return c.additionalJs
}
func (c *clusterrolebind) getAdditionalCss() []string {
	return c.additionalCss
}

func (c *clusterrolebind) addNewResource(s *sysadmServer.Context, module string) error {
	// TODO

	return nil
}

func (c *clusterrolebind) delResource(s *sysadmServer.Context, module string, requestData map[string]string) error {
	// TODO

	return nil
}

func (c *clusterrolebind) showResourceDetail(action string, tplData map[string]interface{}, requestData map[string]string) error {
	// TODO

	return nil
}

func (c *clusterrolebind) getTemplateFile(action string) string {
	switch action {
	case "list":
		return roleBindingsTemplateFiles["list"]
	case "addform":
		return roleBindingsTemplateFiles["addform"]
	default:
		return ""
	}
}

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
	"sysadm/objectsUI"
)

type moduleInfo struct {
	mainModuleName  string
	moduleName      string
	allPopMenuItems []string
	allListItems    map[string]string
	addButtonTile   string
	isSearchForm    string
	namespaced      bool
}

type ingress struct {
	moduleInfo
	orderInfo
}

type configmap struct {
	moduleInfo
	orderInfo
}

type pvc struct {
	moduleInfo
	orderInfo
}

type rolebindings struct {
	moduleInfo
	orderInfo
}

type role struct {
	moduleInfo
	orderInfo
}

type sa struct {
	moduleInfo
	orderInfo
}

type ingressclass struct {
	moduleInfo
	orderInfo
}

type storageclass struct {
	moduleInfo
	orderInfo
}

type pv struct {
	moduleInfo
	orderInfo
}

type secret struct {
	moduleInfo
	orderInfo
}

type clusterrole struct {
	moduleInfo
	orderInfo
}

type clusterrolebind struct {
	moduleInfo
	orderInfo
}

type orderInfo struct {
	allOrderFields        map[string]objectsUI.SortBy
	defaultOrderField     string
	defaultOrderDirection string
}

type objectEntity interface {
	setObjectInfo()
	listObjectData(string, string, int, map[string]string) (int, []map[string]interface{}, error)
	getMainModuleName() string
	getModuleName() string
	getAddButtonTitle() string
	getIsSearchForm() string
	getAllPopMenuItems() []string
	getAllListItems() map[string]string
	getDefaultOrderField() string
	getDefaultOrderDirection() string
	getAllorderFields() map[string]objectsUI.SortBy
	getNamespaced() bool
}

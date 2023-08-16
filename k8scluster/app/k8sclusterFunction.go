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
	"fmt"
	"strings"
	sysadmObjects "sysadm/objects/app"
	"sysadm/utils"
)

func New() K8scluster {
	ret := K8scluster{}
	ret.Name = DefaultObjectName
	ret.TableName = DefaultTableName
	ret.PkName = DefaultPkName
	return ret
}

func (k K8scluster) GetObjectInfoByID(id string) (interface{}, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("can not get cluster information with empty ID")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(k.TableName, k.PkName, id)
	if e != nil {
		return nil, e
	}

	k8sclusterData := K8sclusterSchema{}
	e = sysadmObjects.Unmarshal(dbData, &k8sclusterData)

	return k8sclusterData, e
}

func (k K8scluster) GetObjectCount(searchContent string, ids, searchKeys []string, conditions map[string]string) (int, error) {
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &K8sclusterSchema{}); !ok {
		return -1, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &K8sclusterSchema{}); !ok {
		return -1, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	return sysadmObjects.GetObjectCount(k.TableName, k.PkName, searchContent, ids, searchKeys, conditions)
}

func (k K8scluster) GetObjectList(searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]interface{}, error) {

	var ret []interface{}
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &K8sclusterSchema{}); !ok {
		return ret, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &K8sclusterSchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	var orderKeys []string
	for k, _ := range orders {
		orderKeys = append(orderKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(orderKeys, &K8sclusterSchema{}); !ok {
		return ret, fmt.Errorf("the keys of orders must be the object fields name.error %s", e)
	}

	dbData, e := sysadmObjects.GetObjectList(k.TableName, k.PkName, searchContent, ids, searchKeys, conditions, startPos, step, orders)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		cluster := &K8sclusterSchema{}
		if e := sysadmObjects.Unmarshal(v, cluster); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *cluster)
	}

	return tmpRes, nil
}

func GetStatusText(status int) string {
	statusText := "未知"
	if str, ok := allStatus[status]; ok {
		statusText = str
	}

	return statusText
}

func (k K8scluster) AddObject(data interface{}) error {
	k8sSchemaData, ok := data.(K8sclusterSchema)
	if !ok {
		return fmt.Errorf("the data try to add is not valid")
	}

	idData, e := utils.NewWorker(uint64(k8sSchemaData.Dcid), uint64(k8sSchemaData.Azid))
	if e != nil {
		return e
	}
	clusterID, e := idData.GetID()
	if e != nil {
		return e
	}

	k8sSchemaData.Id = clusterID
	insertData, e := sysadmObjects.Marshal(k8sSchemaData)
	if e != nil {
		return e
	}

	return sysadmObjects.AddObject(k.TableName, "", insertData)
}

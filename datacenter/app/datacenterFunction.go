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
)

func New() Datacenter {
	ret := Datacenter{}
	ret.Name = DefaultObjectName
	ret.TableName = DefaultTableName
	ret.PkName = DefaultPkName
	return ret
}

func (d Datacenter) GetObjectInfoByID(id string) (interface{}, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("can not get datacenter information with empty ID")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(d.TableName, d.PkName, id)
	if e != nil {
		return nil, e
	}

	availablezoneData := Datacenter{}
	e = sysadmObjects.Unmarshal(dbData, &availablezoneData)

	return availablezoneData, e
}

func (d Datacenter) GetObjectCount(searchContent string, ids, searchKeys []string, conditions map[string]string) (int, error) {
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &DatacenterSchema{}); !ok {
		return -1, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &DatacenterSchema{}); !ok {
		return -1, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	return sysadmObjects.GetObjectCount(d.TableName, d.PkName, searchContent, ids, searchKeys, conditions)
}

func (d Datacenter) GetObjectList(searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]interface{}, error) {

	var ret []interface{}
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &DatacenterSchema{}); !ok {
		return ret, fmt.Errorf("search key are not valid. error: %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &DatacenterSchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name. %s", e)
	}

	var orderKeys []string
	for k, _ := range orders {
		orderKeys = append(orderKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(orderKeys, &DatacenterSchema{}); !ok {
		return ret, fmt.Errorf("the keys of orders must be the object fields name.error %s", e)
	}

	dbData, e := sysadmObjects.GetObjectList(d.TableName, d.PkName, searchContent, ids, searchKeys, conditions, startPos, step, orders)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		datacenter := &DatacenterSchema{}
		if e := sysadmObjects.Unmarshal(v, datacenter); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *datacenter)
	}

	return tmpRes, nil
}

func (d Datacenter) AddObject(data interface{}) error {
	dcSchemaData, ok := data.(DatacenterSchema)
	if !ok {
		return fmt.Errorf("the data try to add is not valid")
	}

	insertData, e := sysadmObjects.Marshal(dcSchemaData)
	if e != nil {
		return e
	}

	return sysadmObjects.AddObject(d.TableName, d.PkName, insertData)
}

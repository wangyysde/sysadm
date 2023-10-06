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

func New() Availablezone {
	ret := Availablezone{}
	ret.Name = DefaultObjectName
	ret.TableName = DefaultTableName
	ret.PkName = DefaultPkName
	return ret
}

func (a Availablezone) GetObjectInfoByID(id string) (interface{}, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("can not get available zone information with empty ID")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(a.TableName, a.PkName, id)
	if e != nil {
		return nil, e
	}

	availablezoneData := AvailablezoneSchema{}
	e = sysadmObjects.Unmarshal(dbData, &availablezoneData)

	return availablezoneData, e
}

func (a Availablezone) GetObjectCount(searchContent string, ids, searchKeys []string, conditions map[string]string) (int, error) {
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &AvailablezoneSchema{}); !ok {
		return -1, fmt.Errorf("search key are not valid.error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &AvailablezoneSchema{}); !ok {
		return -1, fmt.Errorf("the keys of conditions must be the object fields name.error %s", e)
	}

	return sysadmObjects.GetObjectCount(a.TableName, a.PkName, searchContent, ids, searchKeys, conditions)
}

func (a Availablezone) GetObjectList(searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]interface{}, error) {

	var ret []interface{}
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &AvailablezoneSchema{}); !ok {
		return ret, fmt.Errorf("search key are not valid.error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &AvailablezoneSchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	var orderKeys []string
	for k, _ := range orders {
		orderKeys = append(orderKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(orderKeys, &AvailablezoneSchema{}); !ok {
		return ret, fmt.Errorf("the keys of orders must be the object fields name. error %s", e)
	}

	dbData, e := sysadmObjects.GetObjectList(a.TableName, a.PkName, searchContent, ids, searchKeys, conditions, startPos, step, orders)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		az := &AvailablezoneSchema{}
		if e := sysadmObjects.Unmarshal(v, az); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *az)
	}

	return tmpRes, nil
}

func (a Availablezone) GetObjectListByDCID(dcID string, conditions map[string]string) ([]interface{}, error) {
	var ret []interface{}
	var emptyString []string

	dcID = strings.TrimSpace(dcID)
	if dcID == "" {
		return ret, fmt.Errorf("get availablezone list should specified datacenter ID")
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &AvailablezoneSchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	newConditoins := make(map[string]string, 0)
	for k, v := range conditions {
		newConditoins[k] = v
	}
	newConditoins["datacenterid"] = "=" + dcID

	dbData, e := sysadmObjects.GetObjectList(a.TableName, a.PkName, "", emptyString, emptyString, newConditoins, 0, 0, nil)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		az := &AvailablezoneSchema{}
		if e := sysadmObjects.Unmarshal(v, az); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *az)
	}

	return tmpRes, nil
}

func (a Availablezone) AddObject(data interface{}) error {
	azSchemaData, ok := data.(AvailablezoneSchema)
	if !ok {
		return fmt.Errorf("the data try to add is not valid")
	}

	insertData, e := sysadmObjects.Marshal(azSchemaData)
	if e != nil {
		return e
	}

	return sysadmObjects.AddObject(a.TableName, a.PkName, insertData)
}

func (a Availablezone) AddObjectByTx(data interface{}) (map[string]interface{}, string, error) {
	addData := make(map[string]interface{}, 0)

	schemaData, ok := data.(AvailablezoneSchema)
	if !ok {
		return addData, "", fmt.Errorf("there is an error occurred when coverting data to Availablezone Schema schema")
	}

	addData, e := sysadmObjects.Marshal(schemaData)
	if e != nil {
		return addData, "", e
	}

	return addData, a.TableName, nil
}

func (a Availablezone) GetObjectIDFieldName() (string, string, error) {
	return a.TableName, a.PkName, nil
}

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

func New() Version {
	ret := Version{}
	ret.Name = defaultObjectName
	ret.TableName = defaultTableName
	ret.PkName = defaultPkName
	return ret
}

func (v Version) GetObjectInfoByID(id string) (interface{}, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("can not get Version information with empty ID")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(v.TableName, v.PkName, id)
	if e != nil {
		return nil, e
	}

	versionData := VersionSchema{}
	e = sysadmObjects.Unmarshal(dbData, &versionData)

	return versionData, e
}

func (v Version) GetObjectInfoByName(name string) (interface{}, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return nil, fmt.Errorf("can not get Version information with empty name")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(v.TableName, "name", name)
	if e != nil {
		return nil, e
	}

	versionData := VersionSchema{}
	e = sysadmObjects.Unmarshal(dbData, &versionData)

	return versionData, e
}

func (v Version) GetObjectCount(searchContent string, ids, searchKeys []string, conditions map[string]string) (int, error) {
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &VersionSchema{}); !ok {
		return -1, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &VersionSchema{}); !ok {
		return -1, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	return sysadmObjects.GetObjectCount(v.TableName, v.PkName, searchContent, ids, searchKeys, conditions)
}

func (v Version) GetObjectList(searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]interface{}, error) {

	var ret []interface{}
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &VersionSchema{}); !ok {
		return ret, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &VersionSchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	var orderKeys []string
	for k, _ := range orders {
		orderKeys = append(orderKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(orderKeys, &VersionSchema{}); !ok {
		return ret, fmt.Errorf("the keys of orders must be the object fields name.error %s", e)
	}

	dbData, e := sysadmObjects.GetObjectList(v.TableName, v.PkName, searchContent, ids, searchKeys, conditions, startPos, step, orders)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		version := &VersionSchema{}
		if e := sysadmObjects.Unmarshal(v, version); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *version)
	}

	return tmpRes, nil
}

func (v Version) AddObject(data interface{}) error {
	versionSchemaData, ok := data.(VersionSchema)
	if !ok {
		return fmt.Errorf("the data try to add is not valid")
	}

	insertData, e := sysadmObjects.Marshal(versionSchemaData)
	if e != nil {
		return e
	}

	return sysadmObjects.AddObject(v.TableName, "", insertData)
}

func (v Version) AddObjectByTx(data interface{}) (map[string]interface{}, string, error) {
	addData := make(map[string]interface{}, 0)

	schemaData, ok := data.(VersionSchema)
	if !ok {
		return addData, "", fmt.Errorf("there is an error occurred when coverting data to Version Schema schema")
	}

	addData, e := sysadmObjects.Marshal(schemaData)
	if e != nil {
		return addData, "", e
	}

	return addData, v.TableName, nil
}

func (v Version) GetObjectIDFieldName() (string, string, error) {
	return v.TableName, v.PkName, nil
}

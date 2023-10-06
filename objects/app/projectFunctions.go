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
)

func New() Project {
	ret := Project{}
	ret.Name = DefaultObjectName
	ret.TableName = DefaultTableName
	ret.PkName = DefaultPkName
	return ret
}

func (p Project) GetObjectInfoByID(id string) (interface{}, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("can not get project information with empty ID")
	}

	dbData, e := GetObjectInfoByID(p.TableName, p.PkName, id)
	if e != nil {
		return nil, e
	}

	projectData := ProjectSchema{}
	e = Unmarshal(dbData, &projectData)

	return projectData, e
}

func (p Project) GetObjectCount(searchContent string, ids, searchKeys []string, conditions map[string]string) (int, error) {
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := ValidKeysInSchema(searchKeys, &ProjectSchema{}); !ok {
		return -1, fmt.Errorf("search key are not valid, error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := ValidKeysInSchema(conditionKeys, &ProjectSchema{}); !ok {
		return -1, fmt.Errorf("the keys of conditions must be the object fields name.error %s", e)
	}

	return GetObjectCount(p.TableName, p.PkName, searchContent, ids, searchKeys, conditions)
}

func (p Project) GetObjectList(searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]interface{}, error) {

	var ret []interface{}
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := ValidKeysInSchema(searchKeys, &ProjectSchema{}); !ok {
		return ret, fmt.Errorf("search key are not valid.error %s ", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := ValidKeysInSchema(conditionKeys, &ProjectSchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name.error %s", e)
	}

	var orderKeys []string
	for k, _ := range orders {
		orderKeys = append(orderKeys, k)
	}
	if ok, e := ValidKeysInSchema(orderKeys, &ProjectSchema{}); !ok {
		return ret, fmt.Errorf("the keys of orders must be the object fields name.error %s", e)
	}

	dbData, e := GetObjectList(p.TableName, p.PkName, searchContent, ids, searchKeys, conditions, startPos, step, orders)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		project := &ProjectSchema{}
		if e := Unmarshal(v, project); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *project)
	}

	return tmpRes, nil
}

func (p Project) AddObject(data interface{}) error {
	projectSchemaData, ok := data.(ProjectSchema)
	if !ok {
		return fmt.Errorf("the data try to add is not valid")
	}

	insertData, e := Marshal(projectSchemaData)
	if e != nil {
		return e
	}

	return AddObject(p.TableName, p.PkName, insertData)
}

func (p Project) AddObjectByTx(data interface{}) (map[string]interface{}, string, error) {
	addData := make(map[string]interface{}, 0)

	schemaData, ok := data.(ProjectSchema)
	if !ok {
		return addData, "", fmt.Errorf("there is an error occurred when coverting data to Project Schema schema")
	}

	addData, e := Marshal(schemaData)
	if e != nil {
		return addData, "", e
	}

	return addData, p.TableName, nil
}

func (p Project) GetObjectIDFieldName() (string, string, error) {
	return p.TableName, p.PkName, nil
}

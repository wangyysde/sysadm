/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2024 Bzhy Network. All rights reserved.
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

func New() Event {
	ret := Event{}
	ret.Name = defaultObjectName
	ret.TableName = defaultTableName
	ret.PkName = defaultPkName
	return ret
}

func (ev Event) GetObjectInfoByID(id string) (interface{}, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("can not get event information with empty ID")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(ev.TableName, ev.PkName, id)
	if e != nil {
		return nil, e
	}

	eventData := EventSchema{}
	e = sysadmObjects.Unmarshal(dbData, &eventData)

	return eventData, e
}

func (ev Event) GetObjectCount(searchContent string, ids, searchKeys []string, conditions map[string]string) (int, error) {
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &EventSchema{}); !ok {
		return -1, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &EventSchema{}); !ok {
		return -1, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	return sysadmObjects.GetObjectCount(ev.TableName, ev.PkName, searchContent, ids, searchKeys, conditions)
}

func (ev Event) GetObjectList(searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]interface{}, error) {

	var ret []interface{}
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &EventSchema{}); !ok {
		return ret, fmt.Errorf("search key are not valid. error: %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &EventSchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name. %s", e)
	}

	var orderKeys []string
	for k, _ := range orders {
		orderKeys = append(orderKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(orderKeys, &EventSchema{}); !ok {
		return ret, fmt.Errorf("the keys of orders must be the object fields name.error %s", e)
	}

	dbData, e := sysadmObjects.GetObjectList(ev.TableName, ev.PkName, searchContent, ids, searchKeys, conditions, startPos, step, orders)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		datacenter := &EventSchema{}
		if e := sysadmObjects.Unmarshal(v, datacenter); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *datacenter)
	}

	return tmpRes, nil
}

func (ev Event) AddObject(data interface{}) error {
	eventSchemaData, ok := data.(EventSchema)
	if !ok {
		return fmt.Errorf("the data try to add is not valid")
	}

	insertData, e := sysadmObjects.Marshal(eventSchemaData)
	if e != nil {
		return e
	}

	return sysadmObjects.AddObject(ev.TableName, ev.PkName, insertData)
}

func (ev Event) AddObjectByTx(data interface{}) (map[string]interface{}, string, error) {
	addData := make(map[string]interface{}, 0)

	schemaData, ok := data.(EventSchema)
	if !ok {
		return addData, "", fmt.Errorf("there is an error occurred when coverting data to Datacenter Schema schema")
	}

	addData, e := sysadmObjects.Marshal(schemaData)
	if e != nil {
		return addData, "", e
	}

	return addData, ev.TableName, nil
}

func (ev Event) GetObjectIDFieldName() (string, string, error) {
	return ev.TableName, ev.PkName, nil
}

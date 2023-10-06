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
	sysadmDB "sysadm/db"
	sysadmObjects "sysadm/objects/app"
)

func New(dbConfig *sysadmDB.DbConfig, workingRoot string) (Yum, error) {
	ret := Yum{}
	if dbConfig == nil && runData.dbConf == nil {
		return ret, fmt.Errorf("DB configuration has not be set")
	}

	if runData.dbConf == nil {
		runData.dbConf = dbConfig
	}

	if workingRoot == "" && runData.workingRoot == "" {
		return ret, fmt.Errorf("working root path has not be set")
	}

	if runData.workingRoot == "" {
		runData.workingRoot = workingRoot
	}

	if sysadmObjects.GetRunDataForDBConf() == nil {
		if e := sysadmObjects.SetRunDataForDBConf(runData.dbConf); e != nil {
			return ret, e
		}
	}

	if sysadmObjects.GetWorkingRoot() == "" {
		if e := sysadmObjects.SetWorkingRoot(runData.workingRoot); e != nil {
			return ret, e
		}
	}
	ret.Name = defaultObjectName
	ret.TableName = defaultTableName
	ret.PkName = defaultPkName
	return ret, nil
}

func (y Yum) GetObjectInfoByID(id string) (interface{}, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("the value of id is empty")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(y.TableName, y.PkName, id)
	if e != nil {
		return nil, e
	}

	versionData := YumSchema{}
	e = sysadmObjects.Unmarshal(dbData, &versionData)

	return versionData, e
}

func (y Yum) GetObjectInfoByName(name string) (interface{}, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(y.TableName, "name", name)
	if e != nil {
		return nil, e
	}

	versionData := YumSchema{}
	e = sysadmObjects.Unmarshal(dbData, &versionData)

	return versionData, e
}

func (y Yum) GetObjectCount(searchContent string, ids, searchKeys []string, conditions map[string]string) (int, error) {
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &YumSchema{}); !ok {
		return -1, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &YumSchema{}); !ok {
		return -1, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	return sysadmObjects.GetObjectCount(y.TableName, y.PkName, searchContent, ids, searchKeys, conditions)
}

func (y Yum) GetObjectList(searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]interface{}, error) {

	var ret []interface{}
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &YumSchema{}); !ok {
		return ret, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &YumSchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	var orderKeys []string
	for k, _ := range orders {
		orderKeys = append(orderKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(orderKeys, &YumSchema{}); !ok {
		return ret, fmt.Errorf("the keys of orders must be the object fields name.error %s", e)
	}

	dbData, e := sysadmObjects.GetObjectList(y.TableName, y.PkName, searchContent, ids, searchKeys, conditions, startPos, step, orders)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		yumdata := &YumSchema{}
		if e := sysadmObjects.Unmarshal(v, yumdata); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *yumdata)
	}

	return tmpRes, nil
}

func (y Yum) AddObject(data interface{}) error {
	versionSchemaData, ok := data.(YumSchema)
	if !ok {
		return fmt.Errorf("the data try to add is not valid")
	}

	insertData, e := sysadmObjects.Marshal(versionSchemaData)
	if e != nil {
		return e
	}

	return sysadmObjects.AddObject(y.TableName, "", insertData)
}

func (y Yum) AddObjectByTx(data interface{}) (map[string]interface{}, string, error) {
	addData := make(map[string]interface{}, 0)

	schemaData, ok := data.(YumSchema)
	if !ok {
		return addData, "", fmt.Errorf("there is an error occurred when coverting data to Yum Schema schema")
	}

	addData, e := sysadmObjects.Marshal(schemaData)
	if e != nil {
		return addData, "", e
	}

	return addData, y.TableName, nil
}

func (y Yum) GetObjectIDFieldName() (string, string, error) {
	return y.TableName, y.PkName, nil
}

func GetYumListByOSIDAndVersionID(osid, versionid int) ([]YumSchema, error) {
	var ret []YumSchema

	if osid < 1 || versionid < 1 {
		return ret, fmt.Errorf("osid or versionid is not valid")
	}

	conditions := make(map[string]string, 0)
	conditions["osid"] = string(osid)
	conditions["versionid"] = string(versionid)
	var emptyString []string
	dbData, e := sysadmObjects.GetObjectList(defaultTableName, defaultPkName, "", emptyString, emptyString, conditions, 0, 0, map[string]string{})
	if e != nil {
		return ret, e
	}

	for _, v := range dbData {
		yumdata := &YumSchema{}
		if e := sysadmObjects.Unmarshal(v, yumdata); e != nil {
			return ret, e
		}
		ret = append(ret, *yumdata)
	}

	return ret, nil
}

func GetYumIDsByObjectList(data []YumSchema) string {
	ids := ""
	for _, item := range data {
		if ids == "" {
			ids = string(item.YumID)
		} else {
			ids = ids + "," + string(item.YumID)
		}
	}

	return ids
}

func GetObjectInfos() (string, string, string, string, string) {
	return defaultObjectName, defaultTableName, defaultPkName, DefaultModuleName, DefaultApiVersion
}

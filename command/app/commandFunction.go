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
	"math/rand"
	"strings"

	sysadmDB "sysadm/db"
	sysadmObjects "sysadm/objects/app"
	sysadmUtils "sysadm/utils"
	"time"
)

func New(dbConfig *sysadmDB.DbConfig, workingRoot string) (Command, error) {
	ret := Command{}

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

func (c Command) GetObjectInfoByID(id string) (interface{}, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("can not get Command defined information with empty ID")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(c.TableName, c.PkName, id)
	if e != nil {
		return nil, e
	}

	commandData := CommandDefinedSchema{}
	e = sysadmObjects.Unmarshal(dbData, &commandData)

	return commandData, e
}

func (c Command) GetObjectInfoByName(name string) (interface{}, error) {
	name = strings.TrimSpace(strings.ToLower(name))
	if name == "" {
		return nil, fmt.Errorf("can not get command information with empty name")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(c.TableName, "name", name)
	if e != nil {
		return nil, e
	}

	commandData := CommandDefinedSchema{}
	e = sysadmObjects.Unmarshal(dbData, &commandData)

	return commandData, e
}

func (c Command) GetObjectCount(searchContent string, ids, searchKeys []string, conditions map[string]string) (int, error) {
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &CommandDefinedSchema{}); !ok {
		return -1, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &CommandDefinedSchema{}); !ok {
		return -1, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	return sysadmObjects.GetObjectCount(c.TableName, c.PkName, searchContent, ids, searchKeys, conditions)
}

func (c Command) GetObjectList(searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]interface{}, error) {

	var ret []interface{}
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &CommandDefinedSchema{}); !ok {
		return ret, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &CommandDefinedSchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	var orderKeys []string
	for k, _ := range orders {
		orderKeys = append(orderKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(orderKeys, &CommandDefinedSchema{}); !ok {
		return ret, fmt.Errorf("the keys of orders must be the object fields name.error %s", e)
	}

	dbData, e := sysadmObjects.GetObjectList(c.TableName, c.PkName, searchContent, ids, searchKeys, conditions, startPos, step, orders)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		command := &CommandDefinedSchema{}
		if e := sysadmObjects.Unmarshal(v, command); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *command)
	}

	return tmpRes, nil
}

func (c Command) AddObject(data interface{}) error {
	versionSchemaData, ok := data.(CommandDefinedSchema)
	if !ok {
		return fmt.Errorf("the data try to add is not valid")
	}

	insertData, e := sysadmObjects.Marshal(versionSchemaData)
	if e != nil {
		return e
	}

	return sysadmObjects.AddObject(c.TableName, "", insertData)
}

func (c Command) AddObjectByTx(data interface{}) (map[string]interface{}, string, error) {
	addData := make(map[string]interface{}, 0)

	schemaData, ok := data.(CommandDefinedSchema)
	if !ok {
		return addData, "", fmt.Errorf("there is an error occurred when coverting data to CommandDefined Schema")
	}

	addData, e := sysadmObjects.Marshal(schemaData)
	if e != nil {
		return addData, "", e
	}

	return addData, c.TableName, nil
}

func (c Command) GetObjectIDFieldName() (string, string, error) {
	return c.TableName, c.PkName, nil
}

func GenerateCommandID() (string, error) {
	// 当前设计要求command是不区分数据中心和可用区的，同时为了使用现有id生成组件，
	// 这里产生伪数据中心和可用区ID，并用其生成command ID
	randInst := rand.NewSource(time.Now().UnixMicro())
	randInst.Seed(time.Now().UnixMicro())
	fakeDCID := rand.Intn(15)
	fakeAZID := rand.Intn(15)

	idData, e := sysadmUtils.NewWorker(uint64(fakeDCID), uint64(fakeAZID))
	if e != nil {
		return "", e
	}

	id, e := idData.GetID()
	if e != nil {
		return "", e
	}

	return string(id), nil
}

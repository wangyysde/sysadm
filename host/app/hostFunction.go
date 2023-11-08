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
	"strconv"
	"strings"
	sysadmCommand "sysadm/command/app"
	sysadmDB "sysadm/db"
	sysadmObjects "sysadm/objects/app"
	sysadmUtils "sysadm/utils"
	sysadmYum "sysadm/yum/app"
)

func New(dbConfig *sysadmDB.DbConfig, workingRoot string) (Host, error) {
	ret := Host{}
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

	ret.Name = DefaultObjectName
	ret.TableName = DefaultTableName
	ret.PkName = DefaultPkName
	return ret, nil
}

func (h Host) GetObjectInfoByID(id string) (interface{}, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return nil, fmt.Errorf("can not get host information with empty ID")
	}

	dbData, e := sysadmObjects.GetObjectInfoByID(h.TableName, h.PkName, id)
	if e != nil {
		return nil, e
	}

	hostData := HostSchema{}
	e = sysadmObjects.Unmarshal(dbData, &hostData)

	return hostData, e
}

func (h Host) GetObjectCount(searchContent string, ids, searchKeys []string, conditions map[string]string) (int, error) {
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &HostSchema{}); !ok {
		return -1, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &HostSchema{}); !ok {
		return -1, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	return sysadmObjects.GetObjectCount(h.TableName, h.PkName, searchContent, ids, searchKeys, conditions)
}

func (h Host) GetObjectList(searchContent string, ids, searchKeys []string, conditions map[string]string,
	startPos, step int, orders map[string]string) ([]interface{}, error) {

	var ret []interface{}
	searchContent = strings.TrimSpace(searchContent)
	if ok, e := sysadmObjects.ValidKeysInSchema(searchKeys, &HostSchema{}); !ok {
		return ret, fmt.Errorf("search key are not valid. error %s", e)
	}

	var conditionKeys []string
	for k, _ := range conditions {
		conditionKeys = append(conditionKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(conditionKeys, &HostSchema{}); !ok {
		return ret, fmt.Errorf("the keys of conditions must be the object fields name. error %s", e)
	}

	var orderKeys []string
	for k, _ := range orders {
		orderKeys = append(orderKeys, k)
	}
	if ok, e := sysadmObjects.ValidKeysInSchema(orderKeys, &HostSchema{}); !ok {
		return ret, fmt.Errorf("the keys of orders must be the object fields name.error %s", e)
	}

	dbData, e := sysadmObjects.GetObjectList(h.TableName, h.PkName, searchContent, ids, searchKeys, conditions, startPos, step, orders)
	if e != nil {
		return ret, e
	}

	var tmpRes []interface{}
	for _, v := range dbData {
		cluster := &HostSchema{}
		if e := sysadmObjects.Unmarshal(v, cluster); e != nil {
			return ret, e
		}
		tmpRes = append(tmpRes, *cluster)
	}

	return tmpRes, nil
}

func (h Host) AddObject(data interface{}) error {
	hostSchemaData, ok := data.(HostSchema)
	if !ok {
		return fmt.Errorf("the data try to add is not valid")
	}

	insertData, e := sysadmObjects.Marshal(hostSchemaData)
	if e != nil {
		return e
	}

	return sysadmObjects.AddObject(h.TableName, "", insertData)
}

func (h Host) GetObjectMaxID() (uint, error) {
	return sysadmObjects.GetObjectMaxID(runData.dbConf.Entity, h.TableName, h.PkName)
}

func (h Host) AddObjectByTx(data interface{}) (map[string]interface{}, string, error) {
	addData := make(map[string]interface{}, 0)

	hostSchemaData, ok := data.(HostSchema)
	if !ok {
		return addData, "", fmt.Errorf("there is an error occurred when coverting data to host schema")
	}

	addData, e := sysadmObjects.Marshal(hostSchemaData)
	if e != nil {
		return addData, "", e
	}

	return addData, h.TableName, nil
}

func (h Host) GetObjectIDFieldName() (string, string, error) {
	return h.TableName, h.PkName, nil
}

func AddHostFromCluster(userid, osID, osversionid, dcid, azid int, hostname, status, k8sclusterid, machineID, systemID, architecture, kernelVersion string, ips []string) error {
	hostIP := ""
	hostIPType := HostTypeIPTypeV4
	if len(ips) > 0 {
		hostIP = ips[0]
		_, hostIPType = sysadmUtils.JudgeIpv4OrIpv6(hostIP)
	}

	hostSchemaData := HostSchema{
		UserId:        strconv.Itoa(userid),
		ProjectID:     0,
		Hostname:      hostname,
		OSID:          osID,
		OSVersionID:   osversionid,
		Status:        status,
		AgentIP:       hostIP,
		IpType:        hostIPType,
		K8sClusterID:  k8sclusterid,
		Dcid:          uint(dcid),
		Azid:          uint(azid),
		MachineID:     machineID,
		SystemID:      systemID,
		Architecture:  architecture,
		KernelVersion: kernelVersion,
	}

	objHost, e := New(runData.dbConf, runData.workingRoot)
	if e != nil {
		return e
	}

	tx, e := sysadmObjects.BeginTx(runData.dbConf.Entity, objHost)
	if e != nil {
		return e
	}

	e = tx.AddObject(hostSchemaData)
	if e != nil {
		tx.Rollback()
		return e
	}

	e = tx.UpdateObjectNextID()
	if e != nil {
		tx.Rollback()
		return e
	}

	hostid, e := sysadmObjects.GetObjectMaxID(tx.Tx.Entity, DefaultTableName, DefaultPkName)
	if e != nil {
		tx.Rollback()
		return e
	}

	// 为新创建的主机创建需要执行的命令及其参数信息
	commandInst, e := sysadmCommand.New(runData.dbConf, runData.workingRoot)
	if e != nil {
		tx.Rollback()
		return e
	}
	yumObjName, _, _, _, _ := sysadmYum.GetObjectInfos()
	commandDefs, e := commandInst.GetDefinitionListForCreateObj(DefaultObjectName, yumObjName, osID, osversionid)
	if e != nil {
		tx.Rollback()
		return e
	}

	yumIDs := ""
	for _, line := range commandDefs {
		commandDefined, ok := line.(sysadmCommand.CommandDefinedSchema)
		if !ok {
			tx.Rollback()
			return fmt.Errorf("the data is not a command definition")
		}
		_, tmpYumIDs, e := commandInst.AddCommandForHostByTx(tx, hostid, commandDefined)
		if e != nil {
			tx.Rollback()
			return e
		}

		tmpYumIDs = strings.TrimSpace(tmpYumIDs)
		if tmpYumIDs != "" {
			if yumIDs == "" {
				yumIDs = tmpYumIDs
			} else {
				yumIDs = yumIDs + "," + tmpYumIDs
			}
		}
	}

	// 增加主机与Yum的关联信息
	if yumIDs != "" {
		tmpYumIDsSlice := strings.Split(yumIDs, ",")
		yumIDsSlice := sysadmUtils.UniqueStringSlice(tmpYumIDsSlice)
		data := make(map[string]interface{}, 0)
		for _, id := range yumIDsSlice {
			data["hostid"] = hostid
			data["yumid"] = id
			e := tx.AddObjectWithMap(defaultHostYumTableName, data)
			if e != nil {
				tx.Rollback()
				return e
			}
		}
	}

	return tx.Commit()
}

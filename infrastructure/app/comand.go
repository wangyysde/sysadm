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
	"time"

	"github.com/wangyysde/sysadm/db"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadm/utils"
)

/*
addCommandToDB add the information of a command into DB
return 1 and []sysadmerror.Sysadmerror
otherwise return 0  and []sysadmerror.Sysadmerror
*/
func addCommandToDB(tx *db.Tx, command string, passiveMode int, hostid int) (int, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(30304001, "debug", "try to add the information of command %s to DB", command))
	// get the last commandid which will be added
	whereStatement := make(map[string]string, 0)
	whereStatement["tableName"] = "command"
	whereStatement["fieldName"] = "commandID"
	selectData := db.SelectData{
		Tb:        []string{"ids"},
		OutFeilds: []string{"nextValue"},
		Where:     whereStatement,
	}

	dbEntity := WorkingData.dbConf.Entity
	retData, err := dbEntity.QueryData(&selectData)
	errs = append(errs, err...)
	if retData == nil {
		return 0, errs
	}

	commandID := 0
	line := retData[0]
	if v, ok := line["nextValue"]; ok {
		id, e := utils.Interface2Int(v)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(303040002, "error", "got commandID error %s", e))
		} else {
			commandID = id
		}
	} else {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(303040003, "error", "got commandID error"))
	}

	insertData := make(db.FieldData, 0)
	insertData["commandID"] = commandID
	insertData["command"] = command
	insertData["hostID"] = hostid
	if passiveMode == 1 {
		insertData["agentPassive"] = 1
	} else {
		insertData["agentPassive"] = 0
	}

	now := time.Now()
	nowInt64 := now.Unix()
	insertData["createTime"] = nowInt64
	insertData["tryTimes"] = 0
	insertData["status"] = 0

	_, err = tx.InsertData("command", insertData)
	errs = append(errs, err...)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error") {
		return 0, errs
	}

	// update commandID value in ids table using transaction
	updateData := make(db.FieldData, 0)
	updateData["nextValue"] = commandID + 1
	whereStatement["tableName"] = "command"
	whereStatement["fieldName"] = "commandID"

	_, err = tx.UpdateData("ids", updateData, whereStatement)
	errs = append(errs, err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
		return 0, errs
	}

	return commandID, errs

}

/*
addCommandParametersToDB add the information of command parameters into DB
return 1 and []sysadmerror.Sysadmerror
otherwise return 0  and []sysadmerror.Sysadmerror
*/
func addCommandParametersToDB(tx *db.Tx, commandID int, paras map[string]string) (int, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(30304002, "debug", "try to add the information of parameters for command to DB"))
	for key, value := range paras {
		insertData := make(db.FieldData, 0)
		insertData["name"] = key
		insertData["value"] = value
		insertData["commandID"] = commandID

		_, err := tx.InsertData("commandParameters", insertData)
		errs = append(errs, err...)
		if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error") {
			return 0, errs
		}
	}

	return 1, errs
}

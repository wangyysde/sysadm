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
*
 */

package app

import (
	"fmt"
	"github.com/wangyysde/sysadmServer"
	"strings"
	"sysadm/db"
	"sysadm/sysadmapi/apiutils"
	"sysadm/sysadmerror"
	"sysadm/utils"
	"time"
)

/*
handler for handling list of the infrastructure
Query parameters of request are below:
conditionKey: key name for DB query ,such as hostid, userid,hostname....
conditionValue: the value of the conditionKey.for hostid, userid,hostname using =, for name, comment using like.
deleted: 0 :normarl 1: deleted
start: start number of the result will be returned.
num: lines of the result will be returned.
*/
func delHost(c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror

	hostMap, err := utils.GetRequestDataArray(c, []string{"hostid[]"})
	hostids := hostMap["hostid[]"]
	errs = append(errs, err...)

	if hostids == nil {
		err = apiutils.SendResponseForErrorMessage(c, 3040001, "no host requested to delete")
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	dbConf := WorkingData.dbConf
	dbEntity := dbConf.Entity
	tx, err := db.Begin(dbEntity)
	errs = append(errs, err...)
	if tx == nil {
		err := apiutils.SendResponseForErrorMessage(c, 3040002, "start transaction on DB query error. check logs for details.")
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	for _, hostid := range hostids {
		if e := tryDelHostData(hostid, tx); e != nil {
			tx.Rollback()
			msg := fmt.Sprintf("%s", e)
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(3040003, "error", "%s", e))
			err := apiutils.SendResponseForErrorMessage(c, 3040003, msg)
			errs = append(errs, err...)
			logErrors(errs)
			return
		}

		whereMap := make(map[string]string, 0)
		whereMap["hostid"] = hostid
		updataData := make(db.FieldData, 0)
		updataData["status"] = "\"deleted\""
		deleteTimeStr := time.Now().Format("2006-01-02 15:04:05")
		updataData["deletetime"] = "\"" + deleteTimeStr + "\""
		_, err := tx.UpdateData("host", updataData, whereMap)
		errs = append(errs, err...)
		if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("error") {
			tx.Rollback()
			errs = append(errs, err...)
			err := apiutils.SendResponseForErrorMessage(c, 3040004, "delete host error")
			errs = append(errs, err...)
			logErrors(errs)
			return
		}
	}

	if e := tx.Commit(); e != nil {
		tx.Rollback()
		msg := fmt.Sprintf("%s", e)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(3040005, "error", "%s", e))
		err := apiutils.SendResponseForErrorMessage(c, 3040005, msg)
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	msg := "host has be deleted successful"
	err = apiutils.SendResponseForSuccessMessage(c, msg)
	errs = append(errs, err...)
	logErrors(errs)
}

func tryDelHostData(hostid string, tx *db.Tx) error {
	hostid = strings.TrimSpace(hostid)
	if hostid == "" {
		return fmt.Errorf("No host requested to delete")
	}

	// delete ip address from DB
	e := tryDelIPFromDBByHostID(hostid, tx)
	if e != nil {
		return e
	}

	// delete command data from DB for host
	e = tryDelCommandDataFromDBByHostID(hostid, tx)
	if e != nil {
		return e
	}

	// delete yum data from DB for host
	e = tryDelYumDataFromDBByHostID(hostid, tx)
	if e != nil {
		return e
	}

	// delete yum data from DB for host
	e = tryDelHostMacFromDBByHostID(hostid, tx)
	if e != nil {
		return e
	}

	return nil
}

func tryDelIPFromDBByHostID(hostid string, tx *db.Tx) error {
	hostid = strings.TrimSpace(hostid)
	if hostid == "" {
		return fmt.Errorf("No host requested to delete")
	}

	whereMap := make(map[string]string, 0)
	whereMap["hostid"] = "=\"" + hostid + "\""

	selectData := db.SelectData{
		Tb:    []string{"hostIP"},
		Where: whereMap,
	}

	affectedRow, _ := tx.DeleteData(&selectData)
	if affectedRow == -1 {
		return fmt.Errorf("delete ip from db error")
	}

	return nil

}

func tryDelCommandDataFromDBByHostID(hostid string, tx *db.Tx) error {
	hostid = strings.TrimSpace(hostid)
	if hostid == "" {
		return fmt.Errorf("No host requested to delete")
	}

	// Qeurying data from DB
	whereMap := make(map[string]string, 0)
	whereMap["hostID"] = "=\"" + hostid + "\""
	selectData := db.SelectData{
		Tb:        []string{"command"},
		OutFeilds: []string{"commandID"},
		Where:     whereMap,
	}
	dbEntity := WorkingData.dbConf.Entity
	dbData, _ := dbEntity.QueryData(&selectData)
	if dbData == nil {
		return fmt.Errorf("delete command data from db error")
	}

	for _, row := range dbData {
		commandID := utils.Interface2String(row["commandID"])
		if commandID != "" {
			whereMap := make(map[string]string, 0)
			whereMap["commandID"] = "=\"" + commandID + "\""

			// delete command history
			deleteData := db.SelectData{
				Tb:    []string{"commandHistory"},
				Where: whereMap,
			}
			affectedRow, _ := tx.DeleteData(&deleteData)
			if affectedRow == -1 {
				return fmt.Errorf("delete command history error")
			}

			// delete command parameters
			deleteData = db.SelectData{
				Tb:    []string{"commandParameters"},
				Where: whereMap,
			}
			affectedRow, _ = tx.DeleteData(&deleteData)
			if affectedRow == -1 {
				return fmt.Errorf("delete command parameters error")
			}

			// delete command status data
			deleteData = db.SelectData{
				Tb:    []string{"commandStatusHistory"},
				Where: whereMap,
			}
			affectedRow, _ = tx.DeleteData(&deleteData)
			if affectedRow == -1 {
				return fmt.Errorf("delete command status error")
			}

			// delete command logs
			deleteData = db.SelectData{
				Tb:    []string{"commandLogs"},
				Where: whereMap,
			}
			affectedRow, _ = tx.DeleteData(&deleteData)
			if affectedRow == -1 {
				return fmt.Errorf("delete command logs error")
			}
		}
	}

	return nil
}

func tryDelYumDataFromDBByHostID(hostid string, tx *db.Tx) error {
	hostid = strings.TrimSpace(hostid)
	if hostid == "" {
		return fmt.Errorf("No host requested to delete")
	}

	whereMap := make(map[string]string, 0)
	whereMap["hostid"] = "=\"" + hostid + "\""

	deleteData := db.SelectData{
		Tb:    []string{"hostYum"},
		Where: whereMap,
	}
	affectedRow, _ := tx.DeleteData(&deleteData)
	if affectedRow == -1 {
		return fmt.Errorf("delete yum data error")
	}

	return nil
}

func tryDelHostMacFromDBByHostID(hostid string, tx *db.Tx) error {
	hostid = strings.TrimSpace(hostid)
	if hostid == "" {
		return fmt.Errorf("No host requested to delete")
	}

	whereMap := make(map[string]string, 0)
	whereMap["hostid"] = "=\"" + hostid + "\""

	deleteData := db.SelectData{
		Tb:    []string{"hostMAC"},
		Where: whereMap,
	}
	affectedRow, _ := tx.DeleteData(&deleteData)
	if affectedRow == -1 {
		return fmt.Errorf("delete host mac error")
	}

	return nil
}

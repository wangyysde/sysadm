/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2023 Bzhy Network. All rights reserved.
* @HomePage http://www.sysadm.cn
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at:
* https://www.sysadm.cn/licenses/apache-2.0.txt
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and  limitations under the License.
*
* Note: 这个文件用于定义处理命令和数据库相关的函数
 */

package app

import (
	"strconv"
	"strings"
	"sysadm/db"
	"sysadm/sysadmerror"
	"sysadm/utils"
	"time"
)

/* 将已经成功下发，但是在规定的时间内未收到命令状态的命令数据移到history表中 */
func moveCommandForOvertime(tx *db.Tx) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	// getting the data of over time command from command table
	nowSecond := int(time.Now().Unix()) + defaultMaxExecuteTime
	whereStatement := make(map[string]string, 0)
	whereStatement["sendTime"] = " > '" + string(nowSecond) + "'"
	whereStatement["status"] = " = '" + string(CommandStatusSent) + "' or status = '" + string(CommandStatusRunning) + "'"
	selectData := db.SelectData{
		Tb:        []string{"command"},
		OutFeilds: []string{"*"},
		Where:     whereStatement,
	}

	dbEntity := runData.dbEntity
	retData, err := dbEntity.QueryData(&selectData)
	errs = append(errs, err...)
	if retData == nil {
		return true, errs
	}

	for _, line := range retData {
		// moving command data from command table to commandHistory
		insertData := make(db.FieldData, 0)
		insertData["commandID"] = line["commandID"]
		insertData["command"] = line["command"]
		insertData["hostID"] = line["hostID"]
		insertData["synchronized"] = line["synchronized"]
		insertData["createTime"] = line["createTime"]
		insertData["sendTime"] = line["sendTime"]
		insertData["completeTime"] = line["completeTime"]
		insertData["tryTimes"] = line["tryTimes"]
		insertData["status"] = CommandStatusTimeout
		insertData["statusMsg"] = "命令执行超时"
		_, err = tx.InsertData("commandHistory", insertData)
		errs = append(errs, err...)
		if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
			return false, errs
		}

		whereStatement := make(map[string]string, 0)
		whereStatement["commandID"] = " = '" + utils.Interface2String(insertData["commandID"]) + "'"
		selectData := db.SelectData{
			Tb:    []string{"command"},
			Where: whereStatement,
		}
		_, err = tx.DeleteData(&selectData)
		errs = append(errs, err...)
		if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
			return false, errs
		}

		selectData = db.SelectData{
			Tb:    []string{"commandParameters"},
			Where: whereStatement,
		}
		_, err = tx.DeleteData(&selectData)
		errs = append(errs, err...)
		if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
			return false, errs
		}
	}

	return true, errs
}

// moveCommandByCommandID 根据commandID将对应的command数据从command表中移到commandHistory表中
// 如果移动成功，则相应的删除commandParameters表中对应command的参数数据
func moveCommandByCommandID(tx *db.Tx, commandID int, msg string) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	if commandID == 0 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20050006, "error", "can not get command data with zero ID"))
		return false, errs
	}

	ok, commandData, err := getCommandFromDBByID(commandID)
	errs = append(errs, err...)
	if !ok {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20050007, "error", "command with ID %d was not found", commandID))
		return true, errs
	}

	commandData["statusMsg"] = msg
	_, err = tx.InsertData("commandHistory", commandData)
	errs = append(errs, err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
		return false, errs
	}

	whereStatement := make(map[string]string, 0)
	whereStatement["commandID"] = " = '" + utils.Interface2String(commandID) + "'"
	selectData := db.SelectData{
		Tb:    []string{"command"},
		Where: whereStatement,
	}
	_, err = tx.DeleteData(&selectData)
	errs = append(errs, err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
		return false, errs
	}
	selectData = db.SelectData{
		Tb:    []string{"commandParameters"},
		Where: whereStatement,
	}
	_, err = tx.DeleteData(&selectData)
	errs = append(errs, err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
		return false, errs
	}

	return true, errs

}

// get command data from DB what will be sent to agent by apiserver which is runnint in active mode
// return true, []commandDataBeSent and []sysadmerror.Sysadmerror if successful
// otherwise retrun false, []commandDataBeSent and []sysadmerror.Sysadmerror
func getCommandDataByStatus(status []CommandStatusCode) (bool, []commandDataBeSent, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	var ret []commandDataBeSent

	if len(status) == 0 {
		for _, v := range AllCommandStatusCode {
			status = append(status, v)
		}
	}

	whereStatement := make(map[string]string, 0)
	statusStr := ""
	for _, v := range status {
		if statusStr == "" {
			statusStr = " = '" + string(v) + "'"
		} else {
			statusStr = statusStr + " or '" + string(v) + "'"
		}
	}
	whereStatement["status"] = statusStr
	whereStatement["tryTimes"] = "< '" + string(defaultCommandExecuteMaxTryTimes) + "'"
	selectData := db.SelectData{
		Tb:        []string{"command"},
		OutFeilds: []string{"*"},
		Where:     whereStatement,
	}

	dbEntity := runData.dbEntity
	retData, err := dbEntity.QueryData(&selectData)
	errs = append(errs, err...)
	if retData == nil {
		return false, ret, errs
	}

	for _, line := range retData {
		lineCommand := commandDataBeSent{
			CommandData: CommandData{},
		}
		commandID := utils.Interface2String(line["commandID"])
		lineCommand.Command.CommandSeq = BuildCommandSeqByID(commandID)
		lineCommand.Command.Command = utils.Interface2String(line["command"])
		synchronized, e := utils.Interface2Int(line["synchronized"])
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20050002, "error", "get command data error %s", e))
			return false, ret, errs
		}
		if synchronized == 0 {
			lineCommand.Command.Synchronized = false
		} else {
			lineCommand.Command.Synchronized = true
		}

		whereStatement := make(map[string]string, 0)
		whereStatement["commandID"] = " = '" + commandID + "'"
		selectData := db.SelectData{
			Tb:        []string{"commandParameters"},
			OutFeilds: []string{"name", "value", "commandID"},
			Where:     whereStatement,
		}
		retData, err := dbEntity.QueryData(&selectData)
		paras := make(map[string]string, 0)
		errs = append(errs, err...)
		if retData != nil {
			for _, line := range retData {
				name := utils.Interface2String(line["name"])
				value := utils.Interface2String(line["value"])
				paras[name] = value
			}
		}
		lineCommand.Command.Parameters = paras

		hostID := utils.Interface2String(line["hostID"])
		hostWhere := make(map[string]string, 0)
		hostWhere["hostid"] = "='" + hostID + "'"
		hostSelect := db.SelectData{
			Tb:        []string{"host"},
			OutFeilds: []string{"agentAddress", "commandUri", "agentIsTls", "agentCa", "agentCert", "agentKey", "insecureSkipVerify", "agentPort"},
			Where:     hostWhere,
		}
		hostData, err := dbEntity.QueryData(&hostSelect)
		errs = append(errs, err...)
		if hostData == nil {
			return false, ret, errs
		}
		hostidInt, _ := utils.Interface2Int(hostID)
		lineCommand.hostID = int32(hostidInt)
		hostLineData := hostData[0]
		lineCommand.agentAddress = utils.Interface2String(hostLineData["agentAddress"])
		commandUri := strings.TrimSpace(utils.Interface2String(hostLineData["commandUri"]))
		if commandUri == "" {
			commandUri = activeSendCommandUri
		}
		lineCommand.commandUri = commandUri
		agentIsTls, _ := utils.Interface2Int(hostLineData["agentIsTls"])
		if agentIsTls == 0 {
			lineCommand.agentIsTls = false
		} else {
			lineCommand.agentIsTls = true
		}
		lineCommand.agentCa = utils.Interface2String(hostLineData["agentCa"])
		lineCommand.agentCert = utils.Interface2String(hostLineData["agentCert"])
		lineCommand.agentKey = utils.Interface2String(hostLineData["agentKey"])
		insecureSkipVerify, _ := utils.Interface2Int(hostLineData["insecureSkipVerify"])
		if insecureSkipVerify == 0 {
			lineCommand.insecureSkipVerify = false
		} else {
			lineCommand.insecureSkipVerify = true
		}
		agentPort, _ := utils.Interface2Int(hostLineData["agentPort"])
		lineCommand.agentPort = agentPort

		ret = append(ret, lineCommand)

	}

	return true, ret, errs
}

// updateCommandStatus update command status in command table with status
func updateCommandStatusToDB(tx *db.Tx, commandID, status, trytimes int) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	// update commandID value in ids table using transaction
	updateData := make(db.FieldData, 0)
	updateData["status"] = status
	updateData["tryTimes"] = trytimes
	whereStatement := make(map[string]string, 0)
	whereStatement["commandID"] = strconv.Itoa(commandID)

	_, err := tx.UpdateData("command", updateData, whereStatement)
	errs = append(errs, err...)
	if sysadmerror.GetMaxLevel(errs) > sysadmerror.GetLevelNum("error") {
		return false, errs
	}
	return true, errs
}

// getCommandFromDBByID get command data from command table in DB server by commandID
func getCommandFromDBByID(commandID int) (bool, map[string]interface{}, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	var ret = make(map[string]interface{}, 0)

	if commandID == 0 {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20050005, "error", "can not get command data from DB with zero ID"))
		return false, ret, errs
	}

	whereStatement := make(map[string]string, 0)
	whereStatement["commandID"] = " = '" + string(commandID) + "'"
	selectData := db.SelectData{
		Tb:        []string{"command"},
		OutFeilds: []string{"*"},
		Where:     whereStatement,
	}

	dbEntity := runData.dbEntity
	retData, err := dbEntity.QueryData(&selectData)
	errs = append(errs, err...)
	if len(retData) < 1 {
		return false, ret, errs
	}

	ret = retData[0]

	return true, ret, errs
}

// addCommandStatusHistory insert command status data into commandStatusHistory table in DB
func addCommandStatusHistory(tx *db.Tx, commandID, hostID, status int, command, statusMsg string) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	nowSecond := int(time.Now().Unix())

	insertData := make(db.FieldData, 0)
	insertData["commandID"] = commandID
	insertData["command"] = command
	insertData["hostID"] = hostID
	insertData["receivedTime"] = nowSecond
	insertData["status"] = status
	insertData["statusMsg"] = statusMsg
	_, err := tx.InsertData("commandStatusHistory", insertData)
	errs = append(errs, err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
		return false, errs
	}

	return true, errs
}

// getCommandDataForGetStatus get command data from DB. apiserver try to get command status or command logs for
// these command
func getCommandDataForGetStatus(status []CommandStatusCode) (bool, []clientRequestData, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	var ret []clientRequestData

	// get all command if status parameters is empty
	if len(status) == 0 {
		for _, v := range AllCommandStatusCode {
			status = append(status, v)
		}
	}

	// join where statement with status parameters
	whereStatement := make(map[string]string, 0)
	statusStr := ""
	for _, v := range status {
		if statusStr == "" {
			statusStr = " = '" + string(v) + "'"
		} else {
			statusStr = statusStr + " or '" + string(v) + "'"
		}
	}
	whereStatement["status"] = statusStr
	selectData := db.SelectData{
		Tb:        []string{"command"},
		OutFeilds: []string{"*"},
		Where:     whereStatement,
	}

	dbEntity := runData.dbEntity
	retData, err := dbEntity.QueryData(&selectData)
	errs = append(errs, err...)
	if retData == nil {
		return false, ret, errs
	}

	for _, line := range retData {
		lineCommand := clientRequestData{
			Command: Command{},
		}
		commandID := utils.Interface2String(line["commandID"])
		lineCommand.Command.CommandSeq = BuildCommandSeqByID(commandID)
		lineCommand.Command.Command = utils.Interface2String(line["command"])
		synchronized, e := utils.Interface2Int(line["synchronized"])
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20050006, "error", "get command data error %s", e))
			return false, ret, errs
		}
		if synchronized == 0 {
			lineCommand.Command.Synchronized = false
		} else {
			lineCommand.Command.Synchronized = true
		}

		hostID := utils.Interface2String(line["hostID"])
		hostWhere := make(map[string]string, 0)
		hostWhere["hostid"] = "='" + hostID + "'"
		hostSelect := db.SelectData{
			Tb:        []string{"host"},
			OutFeilds: []string{"agentAddress", "commandUri", "agentIsTls", "agentCa", "agentCert", "agentKey", "insecureSkipVerify", "agentPort"},
			Where:     hostWhere,
		}
		hostData, err := dbEntity.QueryData(&hostSelect)
		errs = append(errs, err...)
		if hostData == nil {
			return false, ret, errs
		}
		hostidInt, _ := utils.Interface2Int(hostID)
		lineCommand.hostID = int32(hostidInt)
		hostLineData := hostData[0]
		lineCommand.agentAddress = utils.Interface2String(hostLineData["agentAddress"])
		commandUri := strings.TrimSpace(utils.Interface2String(hostLineData["commandUri"]))
		commandStatusUri := strings.TrimSpace(utils.Interface2String(hostLineData["commandStatusUri"]))
		commandLogsUri := strings.TrimSpace(utils.Interface2String(hostLineData["commandLogsUri"]))
		lineCommand.commandUri = commandUri
		lineCommand.commandStatusUri = commandStatusUri
		lineCommand.commandLogsUri = commandLogsUri
		agentIsTls, _ := utils.Interface2Int(hostLineData["agentIsTls"])
		if agentIsTls == 0 {
			lineCommand.agentIsTls = false
		} else {
			lineCommand.agentIsTls = true
		}
		lineCommand.agentCa = utils.Interface2String(hostLineData["agentCa"])
		lineCommand.agentCert = utils.Interface2String(hostLineData["agentCert"])
		lineCommand.agentKey = utils.Interface2String(hostLineData["agentKey"])
		insecureSkipVerify, _ := utils.Interface2Int(hostLineData["insecureSkipVerify"])
		if insecureSkipVerify == 0 {
			lineCommand.insecureSkipVerify = false
		} else {
			lineCommand.insecureSkipVerify = true
		}
		agentPort, _ := utils.Interface2Int(hostLineData["agentPort"])
		lineCommand.agentPort = agentPort

		ret = append(ret, lineCommand)
	}

	return true, ret, errs
}

// addCommandLogToDB add  command log  into commandLogs table in DB
func addCommandLogToDB(tx *db.Tx, log Log, commandID, operation int) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	endTransaction := false
	if tx == nil {
		dbEntity := runData.dbEntity
		tx, err := db.Begin(dbEntity)
		if tx == nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20050007, "error", "start  "+
				"a transaction error %s", err))
			return false, errs
		}
		endTransaction = true
	}

	if commandID == 0 || !IsLogSeqValid(log.LogSeq) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20050008, "error", "log data is not valid "))
		return false, errs
	}

	createTime := int(time.Now().Unix())
	insertData := make(db.FieldData, 0)
	insertData["logSeq"] = log.LogSeq
	insertData["commandID"] = commandID
	insertData["createTime"] = createTime
	insertData["level"] = log.Level
	insertData["operation"] = operation
	insertData["logMessage"] = log.Message
	_, err := tx.InsertData("commandLogs", insertData)
	errs = append(errs, err...)
	if sysadmerror.GetMaxLevel(err) >= sysadmerror.GetLevelNum("error") {
		if endTransaction {
			_ = tx.Rollback()
		}
		return false, errs
	}

	if endTransaction {
		_ = tx.Commit()
	}

	return true, errs
}

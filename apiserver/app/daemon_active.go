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
* Note: this file holds functions of daemon for active mode
 */

package app

import (
	"encoding/json"
	"fmt"
	"math"
	"strconv"
	"strings"
	"sysadm/db"
	"sysadm/httpclient"
	"sysadm/redis"
	"sysadm/sysadmerror"
	sysadmutils "sysadm/utils"
	"time"
)

func RunDaemonActive() {

	go sendCommandsLoop()

	go getCommandStatusLoop()

	getCommandLogsLoop()
}

func sendCommandsLoop() {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040001, "debug", "now try to start main loop for sending command"))

	dbEntity := runData.dbEntity
	for {
		tx, err := db.Begin(dbEntity)
		errs = append(errs, err...)
		if tx != nil {
			ok, err := moveCommandForOvertime(tx)
			errs = append(errs, err...)
			if !ok {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040002, "warning", "cleaning history data of command error"))
				tx.Rollback()
			}
			tx.Commit()

			ok, commandData, err := getCommandDataByStatus([]CommandStatusCode{CommandStatusCreated, ComandStatusSendError, CommandStatusError})
			errs = append(errs, err...)
			logErrors(errs)
			errs = errs[0:0]
			if !ok {
				continue
			}
			for _, lineData := range commandData {
				go sendCommand(&lineData)
				if runData.command.concurrencySendCommand == 0 {
					runData.command.concurrencySendCommand = defaultConcurrencySendCommand
				}
				delay := math.Ceil(float64(1000 / runData.command.concurrencySendCommand))
				time.Sleep(time.Duration(delay) * time.Millisecond)
			}
		}

		if runData.runConf.ConfGlobal.CheckCommandInterval == 0 {
			runData.runConf.ConfGlobal.CheckCommandInterval = defaultCheckCommandInterval
		}
		time.Sleep(time.Duration(runData.runConf.ConfGlobal.CheckCommandInterval) * time.Second)
	}

}

// sendCommand try to send command data to client when apiserver is running in active mode
// additional: the following things will be to do:
// 1. move command data from command table into commandHistory if any error has occurred and trytimes is bigger defined.
// 2. update the status for the command in command table if the status is not the status of task
// 3. insert the status for the command and its subtask into commandStatusHistory table
func sendCommand(data *commandDataBeSent) {
	var errs []sysadmerror.Sysadmerror
	if data == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040003, "warning", "no command  "+
			"data will be send to client"))
		logErrors(errs)
		return
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040004, "debug", "try to send "+
		"command data of command %s to client %s", data.Command.CommandSeq, data.agentAddress))
	dbEntity := runData.dbEntity
	tx, err := db.Begin(dbEntity)
	if tx == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040005, "error", "start  "+
			"a transaction error %s", err))
		logErrors(errs)
		return
	}
	commandID, e := GetCommandIDFromSeq(data.CommandSeq)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040006, "error", "%s", e))
		logErrors(errs)
		return
	}
	httpClient, e := createHttpClient(0, 0, 0, 0, "",
		[]byte(data.agentCa), []byte(data.agentCert), []byte(data.agentKey), data.insecureSkipVerify, data.agentIsTls)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040007, "error", "can not create"+
			"http request client. error %s", e))

		ok, err := updateCommandStatus(tx, commandID, int(data.hostID), int(ComandStatusSendError), data.Command.Command, "can not connect to agent "+data.agentAddress)
		errs = append(errs, err...)
		if ok {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
		logErrors(errs)
		return
	}

	ok, requestParas, err := buildSendCommandRequestParas(data)
	errs = append(errs, err...)
	if !ok {
		ok, err := updateCommandStatus(tx, commandID, int(data.hostID), int(ComandStatusSendError), data.Command.Command, "can not connect to agent "+data.agentAddress)
		errs = append(errs, err...)
		if ok {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
		logErrors(errs)
		return
	}

	commandDataJson, e := json.Marshal(data.CommandData)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040008, "error", "marshal command"+
			" data to json format error %s", e))
		ok, err := updateCommandStatus(tx, commandID, int(data.hostID), int(ComandStatusSendError), data.Command.Command, "can not connect to agent "+data.agentAddress)
		errs = append(errs, err...)
		if ok {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
		logErrors(errs)
		return
	}

	body, e := httpclient.NewSendRequest(requestParas, httpClient, strings.NewReader(sysadmutils.Bytes2str(commandDataJson)))
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040009, "error", "send command data"+
			" to client error %s", e))
		ok, err := updateCommandStatus(tx, commandID, int(data.hostID), int(ComandStatusSendError), data.Command.Command, "can not connect to agent "+data.agentAddress)
		errs = append(errs, err...)
		if ok {
			_ = tx.Commit()
		} else {
			_ = tx.Rollback()
		}
		logErrors(errs)
		return
	}

	handleResponseForSendCommand(data, body)
}

// updateCommandStatus update command status in command table and insert the command status into commandStatusHistory table
// updateCommandStatus will roll back or commit the transaction of DB operation if tx is nil. otherwise,the caller
// should roll back or commit the transaction.
// updateCommandStatus will move the command data from command table into commandHistory if the try times is bigger the maxvalue
func updateCommandStatus(tx *db.Tx, commandID, hostID, status int, command, statusMsg string) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror
	var endTransaction bool = false

	if tx == nil {
		dbEntity := runData.dbEntity
		tx, err := db.Begin(dbEntity)
		if tx == nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040009, "error", "start  "+
				"a transaction error %s", err))
			return false, errs
		}
		endTransaction = true
	}

	// update status in command table
	if status != int(CommandStatusTaskOk) && status != int(CommandStatusTaskError) {
		newTrytimes := 0
		commandMoved := false
		if status == int(ComandStatusSendError) || status == int(CommandStatusError) {
			ok, ret, err := getCommandFromDBByID(commandID)
			errs = append(errs, err...)
			if !ok {
				return false, errs
			}

			tryTimes, e := sysadmutils.Interface2Int(ret["tryTimes"])
			if e != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040010, "error", "a error %s"+
					"has occurred", e))
				return false, errs
			}
			maxTrytimes := runData.command.maxTryTimes
			if maxTrytimes == 0 {
				runData.command.maxTryTimes = defaultMaxCommandTrytimes
				maxTrytimes = defaultMaxCommandTrytimes
			}
			tryTimes = tryTimes + 1
			if tryTimes >= maxTrytimes {
				ok, err := moveCommandByCommandID(tx, commandID, "命令已经尝试最大重试次数")
				errs = append(errs, err...)
				if !ok {
					if endTransaction {
						tx.Rollback()
					}
					return false, errs
				} else {
					commandMoved = true
				}
			}
			newTrytimes = tryTimes
		}

		if !commandMoved {
			ok, err := updateCommandStatusToDB(tx, commandID, status, newTrytimes)
			errs = append(errs, err...)
			if !ok {
				if endTransaction {
					tx.Rollback()
				}
				return false, errs
			}
		}
	}

	ok, err := addCommandStatusHistory(tx, commandID, hostID, status, command, statusMsg)
	errs = append(errs, err...)
	if !ok {
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

// handleResponseForSendCommand handle the response what were responded  from the client
func handleResponseForSendCommand(data *commandDataBeSent, body []byte) {
	var errs []sysadmerror.Sysadmerror

	commandID, e := GetCommandIDFromSeq(data.CommandSeq)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040010, "error", "%s", e))
		logErrors(errs)
		return
	}

	res := &RepStatus{}
	e = json.Unmarshal(body, res)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040011, "error", "send command data"+
			" to client error %s", e))
		_, err := updateCommandStatus(nil, commandID, int(data.hostID), int(ComandStatusSendError), data.Command.Command, "客户端%s 返回错误 "+data.agentAddress)
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	status := res.StatusCode
	statusMsg := ""
	if status == ComandStatusSendError {
		statusMsg = "命令下发到客户" + data.agentAddress + "出错"
	} else {
		statusMsg = "命令已经成功下发到客户端" + data.agentAddress
	}

	_, err := updateCommandStatus(nil, commandID, int(data.hostID), int(status), data.Command.Command, statusMsg)
	errs = append(errs, err...)
	logErrors(errs)

}

// main loop function for get command status when apiserver is running in active mode
func getCommandStatusLoop() {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040012, "debug", "now try to start main loop for getting command status"))

	for {
		ok, clientRequestData, err := getCommandDataForGetStatus([]CommandStatusCode{ComandStatusReceived, CommandStatusSent, CommandStatusRunning})
		errs = append(errs, err...)
		logErrors(errs)
		errs = errs[0:0]
		if !ok {
			continue
		}
		for _, lineData := range clientRequestData {
			go getCommandStatus(&lineData)
			if runData.command.concurrencyGetCommandStatus == 0 {
				runData.command.concurrencyGetCommandStatus = defaultConcurrencyGetCommandStatus
			}
			delay := math.Ceil(float64(1000 / runData.command.concurrencyGetCommandStatus))
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		if runData.runConf.ConfGlobal.GetStatusInterval == 0 {
			runData.runConf.ConfGlobal.GetStatusInterval = defaultGetStatusInterval
		}
		time.Sleep(time.Duration(runData.runConf.ConfGlobal.GetStatusInterval) * time.Second)
	}
}

// getCommandStatus try to get command status from a client when apiserver is running in active mode
func getCommandStatus(data *clientRequestData) {
	var errs []sysadmerror.Sysadmerror
	if data == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040013, "warning", "no command status will be got status "))
		logErrors(errs)
		return
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040014, "debug", "try to get "+
		"status of command %s from client %s", data.Command.CommandSeq, data.agentAddress))

	httpClient, e := createHttpClient(0, 0, 0, 0, "",
		[]byte(data.agentCa), []byte(data.agentCert), []byte(data.agentKey), data.insecureSkipVerify, data.agentIsTls)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040017, "error", "can not create"+
			"http request client. error %s", e))

		logErrors(errs)
		return
	}

	ok, requestParas, err := buildClientRequestParas(data.agentAddress, data.commandStatusUri, data.agentPort, data.agentIsTls)
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		return
	}

	commandSeq := data.Command.CommandSeq
	commandStatusReq := CommandStatusReq{
		CommandSeq:       commandSeq,
		NodeIdentiferStr: "",
	}
	commandStatusJson, e := json.Marshal(commandStatusReq)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040018, "error", "marshal command"+
			" data to json format error %s", e))
		logErrors(errs)
		return
	}

	body, e := httpclient.NewSendRequest(requestParas, httpClient, strings.NewReader(sysadmutils.Bytes2str(commandStatusJson)))
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040019, "error", "send command data"+
			" to client error %s", e))
		logErrors(errs)
		return
	}

	handleResponseForGetStatus(data, body)
}

// handleResponseForGetStatus handle the response what were responded  from the client for get command status request
func handleResponseForGetStatus(data *clientRequestData, body []byte) {
	var errs []sysadmerror.Sysadmerror

	commandStatus := &CommandStatus{}
	e := json.Unmarshal(body, commandStatus)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040020, "error", "get command status error %s", e))
		logErrors(errs)
		return
	}

	commandSeq := data.Command.CommandSeq
	commandID, e := GetCommandIDFromSeq(commandSeq)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040021, "error", "get command status error %s", e))
		logErrors(errs)
		return
	}

	// if the client return command is not found, then we should set command status to CommandStatusError
	if commandStatus.NotCommand {
		statusMsg := fmt.Sprintf("get command status error %s", e)
		_, err := updateCommandStatus(nil, commandID, int(data.hostID), int(CommandStatusError), data.Command.Command, statusMsg)
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	StatusCode := commandStatus.StatusCode
	_, err := updateCommandStatus(nil, commandID, int(data.hostID), int(StatusCode), data.Command.Command, commandStatus.StatusMessage)
	errs = append(errs, err...)
	logErrors(errs)

}

// main loop function for get command logs when apiserver is running in active mode
func getCommandLogsLoop() {
	var errs []sysadmerror.Sysadmerror
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040022, "debug", "now try to start main loop for getting command logs"))

	for {
		ok, clientRequestData, err := getCommandDataForGetStatus([]CommandStatusCode{ComandStatusReceived, CommandStatusSent, CommandStatusRunning})
		errs = append(errs, err...)
		logErrors(errs)
		errs = errs[0:0]
		if !ok {
			continue
		}
		for _, lineData := range clientRequestData {
			go getCommandLog(&lineData)
			if runData.command.concurrencyGetCommandLog == 0 {
				runData.command.concurrencyGetCommandLog = defaultConcurrencyGetCommandLog
			}
			delay := math.Ceil(float64(1000 / runData.command.concurrencyGetCommandLog))
			time.Sleep(time.Duration(delay) * time.Millisecond)
		}

		if runData.runConf.ConfGlobal.GetLogInterval == 0 {
			runData.runConf.ConfGlobal.GetLogInterval = defaultGetLogInterval
		}
		time.Sleep(time.Duration(runData.runConf.ConfGlobal.GetLogInterval) * time.Second)
	}
}

// getCommandLog try to get command log from a client when apiserver is running in active mode
func getCommandLog(data *clientRequestData) {
	var errs []sysadmerror.Sysadmerror
	if data == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040023, "warning", "no command logs will be got status "))
		logErrors(errs)
		return
	}

	httpClient, e := createHttpClient(0, 0, 0, 0, "",
		[]byte(data.agentCa), []byte(data.agentCert), []byte(data.agentKey), data.insecureSkipVerify, data.agentIsTls)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040024, "error", "can not create"+
			"http request client. error %s", e))

		logErrors(errs)
		return
	}

	if strings.TrimSpace(data.commandLogsUri) == "" {
		data.commandLogsUri = defaultCommandLogsUri
	}
	ok, requestParas, err := buildClientRequestParas(data.agentAddress, data.commandLogsUri, data.agentPort, data.agentIsTls)
	errs = append(errs, err...)
	if !ok {
		logErrors(errs)
		return
	}

	commandSeq := data.CommandSeq
	if runData.command.logRootPathInRedis == "" {
		runData.command.logRootPathInRedis = defaultLogRootPathInRedis
	}
	redisKey := runData.command.logRootPathInRedis + "/" + commandSeq

	if runData.command.maxGetLogNumPerTime == 0 {
		runData.command.maxGetLogNumPerTime = defaultMaxGetLogNumPerTime
	}

	ok, e = redis.Exists(runData.redisEntity, runData.redisCtx, redisKey)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040025, "error", "an error has occurred"+
			"%s", e))
		logErrors(errs)
		return
	}
	dateStr := time.Now().Local().Format("20060102")
	startPos := dateStr + "000000"
	if ok {
		lastPos, e := redis.Get(runData.redisEntity, runData.redisCtx, redisKey)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040026, "error", "get last position error "+
				"%s", e))
			logErrors(errs)
			return
		}

		dateStr, id, e := GetLogIDFromSeq(lastPos)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040027, "error", "get last position error "+
				"%s", e))
			logErrors(errs)
			return
		}

		id = id + 1
		startPos = dateStr + strconv.Itoa(id)
	}

	logReq := LogReq{
		CommandSeq:       commandSeq,
		NodeIdentiferStr: "",
		StartSeq:         startPos,
		Num:              runData.command.maxGetLogNumPerTime,
	}
	logJson, e := json.Marshal(logReq)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040028, "error", "marshal log"+
			" request  to json format error %s", e))
		logErrors(errs)
		return
	}

	body, e := httpclient.NewSendRequest(requestParas, httpClient, strings.NewReader(sysadmutils.Bytes2str(logJson)))
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040029, "error", "get log from client"+
			" error %s", e))
		logErrors(errs)
		return
	}

	handleResponseForGetLogs(commandSeq, body)
}

// handleResponseForGetLogs handle the response what were responded  from the client
func handleResponseForGetLogs(commandSeq string, body []byte) {
	var errs []sysadmerror.Sysadmerror

	logData := &LogData{}
	e := json.Unmarshal(body, logData)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040030, "error", "get command log error %s", e))
		logErrors(errs)
		return
	}

	// if the client return command is not found,
	if logData.NotCommand {
		logErrors(errs)
		return
	}

	if !IsCommandSeqValid(commandSeq) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040031, "error", "command sequence is not valid", commandSeq))
		logErrors(errs)
		return
	}
	commandID, e := GetCommandIDFromSeq(commandSeq)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040032, "error", "command sequence is not valid", commandSeq))
		logErrors(errs)
		return
	}

	dbEntity := runData.dbEntity
	tx, err := db.Begin(dbEntity)
	if tx == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(20040033, "error", "start  "+
			"a transaction error %s", err))
		logErrors(errs)
		return
	}

	logs := logData.Logs
	for _, l := range logs {
		ok, err := addCommandLogToDB(tx, l, commandID, 3)
		errs = append(errs, err...)
		if !ok {
			_ = tx.Rollback()
			logErrors(errs)
			return
		}
	}

	_ = tx.Commit()
	logErrors(errs)
}

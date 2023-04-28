/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
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
 */

package app

import (
	"fmt"
	"net/http"
	"strings"
	"encoding/json"

	apiserver "sysadm/apiserver/app"
	"sysadm/redis"
	"sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

func addHandlers(r *sysadmServer.Engine) (errs []sysadmerror.Sysadmerror) {

	if e := addRootHandler(r); e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100803001, "fatal", "add root handler error %s", e))
		return errs
	}

	// we should build nodeIdentifier after start listen
	nodeIdentifier, e := apiserver.BuildNodeIdentifer(RunConf.Global.NodeIdentifer)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(100803001, "fatal", "build node identifier error %s", e))
		return errs
	}
	runData.nodeIdentifer = &nodeIdentifier

	if !RunConf.Agent.Passive {
		// add handler for the path of uri specifing if agent running in active mode
		if e := addReceiveCommandHandler(r); e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(100803002, "fatal", "add receive command  handler error %s", e))
			return errs
		}

		if e := addGetCommandStatusHandler(r); e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(100803003, "fatal", "add get command status  handler error %s", e))
			return errs
		}

		if e := addGetLogsHandler(r); e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(100803004, "fatal", "add get logs  handler error %s", e))
			return errs
		}
	}

	return errs
}

// addRootHandler adding handler for root path
func addRootHandler(r *sysadmServer.Engine) error {
	if r == nil {
		return fmt.Errorf("router is nil")
	}

	r.Any("/", func(c *sysadmServer.Context) {
		c.JSON(200, sysadmServer.H{
			"status": "ok",
		})
	})

	return nil
}

// listener handler for receiving command data what the apiserver send to the agent
func addReceiveCommandHandler(r *sysadmServer.Engine) error {
	if strings.TrimSpace(RunConf.Global.Uri) == "" {
		RunConf.Global.Uri = defaultReceiveCommandUri
	}

	listenUri := RunConf.Global.Uri
	if listenUri[0:1] != "/" {
		listenUri = "/" + listenUri
	}

	r.POST(listenUri, receivedCommand)
	r.GET(listenUri, receivedCommand)

	return nil
}

// listener handler for receiving CommandStatusReq data  what the apiserver send to the agent
func addGetCommandStatusHandler(r *sysadmServer.Engine) error {
	if strings.TrimSpace(RunConf.Global.CommandStatusUri) == "" {
		RunConf.Global.CommandStatusUri = defaultGetCommandStatus
	}

	listenUri := RunConf.Global.CommandStatusUri
	if listenUri[0:1] != "/" {
		listenUri = "/" + listenUri
	}

	r.POST(listenUri, getCommandStatus)
	r.GET(listenUri, getCommandStatus)

	return nil
}

// // listener handler for receiving LogReq data  what the apiserver send to the agent
func addGetLogsHandler(r *sysadmServer.Engine) error {
	if strings.TrimSpace(RunConf.Global.CommandLogsUri) == "" {
		RunConf.Global.CommandLogsUri = defaultGetCommandLogs
	}

	listenUri := RunConf.Global.CommandStatusUri
	if listenUri[0:1] != "/" {
		listenUri = "/" + listenUri
	}

	r.POST(listenUri, getCommandLogs)
	r.GET(listenUri, getCommandLogs)

	return nil
}

func receivedCommand(c *sysadmServer.Context) {
	var cmd apiserver.CommandData = apiserver.CommandData{}
	var errs []sysadmerror.Sysadmerror

	err := c.BindJSON(&cmd)
	if err != nil {
		data := make(map[string]interface{}, 0)
		msg := fmt.Sprintf("receive command error %s", err)
		commandStatus, e := apiserver.BuildCommandStatus("", RunConf.Global.NodeIdentifer, msg, *runData.nodeIdentifer, apiserver.ComandStatusSendError, data, true)
		c.JSON(http.StatusOK, commandStatus)

		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082001, "error", msg))
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082001, "error", fmt.Sprintf("build command status error: %s", e)))
		logErrors(errs)
		return
	}

	doRouteCommand(&cmd, c)
}

func getCommandStatus(c *sysadmServer.Context) {
	var cmdStatusReq apiserver.CommandStatusReq = apiserver.CommandStatusReq{}
	var errs []sysadmerror.Sysadmerror

	err := c.BindJSON(&cmdStatusReq)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082001, "error", "the request for getting command status is not valid %s ", err))
		statusData := make(map[string]interface{}, 0)
		commandData := apiserver.CommandData{
			NodeIdentiferStr: "",
			Command: apiserver.Command{
				CommandSeq: "0000000000000000000",
			},
		}

		_, err := handleCommandStatus(c, &commandData, fmt.Sprintf("the request for getting command status is not valid %s ", err), statusData, apiserver.ComandStatusSendError, true, true)
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	e := responseForGetCommandStatus(c, cmdStatusReq)
	errs = append(errs, e...)
	logErrors(errs)

}

// response command status data to apiserver when apiserver get the data
func responseForGetCommandStatus(c *sysadmServer.Context, req apiserver.CommandStatusReq) (errs []sysadmerror.Sysadmerror) {
	commandSeq := req.CommandSeq

	commandData := apiserver.CommandData{
		NodeIdentiferStr: "",
		Command: apiserver.Command{
			CommandSeq: "0000000000000000000",
		},
	}
	statusData := make(map[string]interface{}, 0)

	if !apiserver.IsCommandSeqValid(commandSeq) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082002, "error", "command sequence %s is not valid ", commandSeq))
		_, err := responseCommandStatusToServer(c, &commandData, fmt.Sprintf("command sequence %s is not valid ", commandSeq), statusData, apiserver.ComandStatusSendError, true, false)
		errs = append(errs, err...)
		return errs
	}
	commandData.CommandSeq = commandSeq

	// we should delete data of command status in redis server
	key := defaultRootPathCommandStatus + commandSeq
	exist, e := redis.Exists(runData.redisEntity, runData.redisctx, key)
	if e != nil || !exist {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082003, "error", "command status data for command %s is not exist ", commandSeq))
		_, err := responseCommandStatusToServer(c, &commandData, fmt.Sprintf("command status data for command %s is not exist ", commandSeq), statusData, apiserver.ComandStatusSendError, true, false)
		errs = append(errs, err...)
		return errs
	}

	isChanged, err := isNodeIdentiferChanged(req.NodeIdentiferStr)
	errs = append(errs, err...)
	if isChanged {
		newNodeIdentifer, err := apiserver.BuildNodeIdentifer(strings.ToUpper(strings.TrimSpace(req.NodeIdentiferStr)))
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082004, "error", "change node identifer error %s", err))
			_, err := responseCommandStatusToServer(c, &commandData, fmt.Sprintf("change node identifer error %s", err), statusData, apiserver.ComandStatusSendError, false, false)
			errs = append(errs, err...)
			return errs
		} else {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082004, "debug", "node identifier has be changed "))
			runData.nodeIdentifer = &newNodeIdentifer
		}
	}

	gotStatusData, e := redis.HGetAll(runData.redisEntity, runData.redisctx, key)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082005, "error", "an error has occurred %s ", e))
		_, err := responseCommandStatusToServer(c, &commandData, fmt.Sprintf("an error has occurred %s ", e), statusData, apiserver.CommandStatusError, false, false)
		errs = append(errs, err...)
		return errs
	}

	for v, k := range gotStatusData {
		statusData[k] = interface{}(v)
	}

	_, err = responseCommandStatusToServer(c, &commandData, "", statusData, apiserver.CommandStatusOK, false, true)
	errs = append(errs, err...)
	return errs

}

// this function gets logs of a command and set them to the apiserver
func getCommandLogs(c *sysadmServer.Context) {
	var cmdLogReq apiserver.LogReq = apiserver.LogReq{}
	var errs []sysadmerror.Sysadmerror

	logs := make([]apiserver.Log,0)
	logData := apiserver.LogData{
		CommandSeq: "0000000000000000000",
		NodeIdentifier: *runData.nodeIdentifer,
		Logs: logs,
		Total: 0,
		EndFlag: false,
		NotCommand: true,
	}

	e := c.BindJSON(&cmdLogReq)
	if e != nil { 
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082006, "error", "the request for getting command logs is not valid %s ", e))
		_, err := responseCommandLogToServer(c, logData, false) 
		errs = append(errs, err...)
		logErrors(errs) 
		return
	}

	commandSeq := strings.TrimSpace(cmdLogReq.CommandSeq)
	if strings.TrimSpace(commandSeq) == "" || !apiserver.IsCommandSeqValid(commandSeq) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082007, "error", "command sequence %s is not valid ", commandSeq))
		_, err := responseCommandLogToServer(c, logData, false) 
		errs = append(errs, err...)
		logErrors(errs) 
		return
	}
	logData.CommandSeq = commandSeq
	logData.NotCommand = false

	nodeIdentiferStr := strings.TrimSpace(cmdLogReq.NodeIdentiferStr)
	isChanged, err := isNodeIdentiferChanged(nodeIdentiferStr)
	errs = append(errs, err...)
	if isChanged {
		newNodeIdentifer, err := apiserver.BuildNodeIdentifer(strings.ToUpper(strings.TrimSpace(nodeIdentiferStr)))
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082008, "error", "change node identifer error %s", err))
			_, err := responseCommandLogToServer(c, logData, false) 
			errs = append(errs, err...)
			logErrors(errs)
			return
		} else {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082008, "debug", "node identifier has be changed "))
			runData.nodeIdentifer = &newNodeIdentifer
			logData.NodeIdentifier = newNodeIdentifer
		}
	}

	key := defaultRootPathCommandLog + commandSeq
	exist, e := redis.Exists(runData.redisEntity, runData.redisctx, key)
	if !exist {
		logData.EndFlag = true
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082009, "debug", "no command los for command %s ", commandSeq))
		_, err := responseCommandLogToServer(c, logData, false) 
		errs = append(errs, err...)
		logErrors(errs)
		return
	}
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082010, "error", "an error has occurred when check whether command logs of command %s is exist. error %s ", commandSeq,e))
		_, err := responseCommandLogToServer(c, logData, false) 
		errs = append(errs, err...)
		logErrors(errs)
	}

	maxNum := 0
	if cmdLogReq.Num < 1 || cmdLogReq.Num > maxLogNumPerRequest {
		maxNum = maxLogNumPerRequest
	}
	listLen, _ := redis.LLen(runData.redisEntity, runData.redisctx, key)
	endFlag := false
	if listLen < maxNum {
		maxNum = listLen
		endFlag = true
	}
	logData.EndFlag = endFlag

	total := 0
	for i :=0; i<maxNum; i++ {
		logJson, e := redis.LPop(runData.redisEntity, runData.redisctx, key)
		if e != nil {
			continue
		}

		log := apiserver.Log{}
		e = json.Unmarshal([]byte(logJson), &log)
		if e != nil {
			continue
		}
		
		logs = append(logs,log)
		total = total + 1
	}

	logData.Total = total
	logData.Logs = logs
	_, err = responseCommandLogToServer(c, logData, true) 
	errs = append(errs,err...)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(10082011, "info", "logs of command %s has be sent to apiserver ", commandSeq))
	logErrors(errs)
}




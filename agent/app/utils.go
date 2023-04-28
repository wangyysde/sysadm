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
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	apiserver "sysadm/apiserver/app"
	"sysadm/config"
	"sysadm/httpclient"
	redis "sysadm/redis"
	"sysadm/sysadmerror"
	"sysadm/utils"
	"github.com/wangyysde/sysadmLog"
	"github.com/wangyysde/sysadmServer"
)

func SetVersion(version *config.Version) {
	if version == nil {
		return
	}

	version.Version = ver
	version.Author = author

	CliOps.Version = *version
	RunConf.Version = *version
}

func GetVersion() *config.Version {
	if CliOps.Version.Version != "" {
		return &CliOps.Version
	}

	return nil
}

/*
log log messages to logfile or stdout
*/
func logErrors(errs []sysadmerror.Sysadmerror) {

	for _, e := range errs {
		l := sysadmerror.GetErrorLevelString(e)
		no := e.ErrorNo
		sysadmServer.Logf(l, "erroCode: %d Msg: %s", no, e.ErrorMsg)
	}

}

// set parameters to accessLogger and errorLooger
func setLogger() []sysadmerror.Sysadmerror {
	var errs []sysadmerror.Sysadmerror

	sysadmServer.SetLoggerKind(RunConf.Global.Log.Kind)
	sysadmServer.SetLogLevel(RunConf.Global.Log.Level)
	sysadmServer.SetTimestampFormat(RunConf.Global.Log.TimeStampFormat)
	_, fp, err := sysadmServer.SetAccessLogFile(RunConf.Global.Log.AccessLog)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081002, "error", "can not set access log file(%s) error: %s", RunConf.Global.Log.AccessLog, err))
	} else {
		RunConf.Global.Log.AccessLogFp = fp
	}

	if RunConf.Global.Log.SplitAccessAndError && RunConf.Global.Log.ErrorLog != "" {
		_, fp, err := sysadmServer.SetErrorLogFile(RunConf.Global.Log.ErrorLog)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081003, "error", "can not set error log file(%s) error: %s", RunConf.Global.Log.ErrorLog, err))
		} else {
			RunConf.Global.Log.ErrorLogFp = fp
		}
	}

	sysadmServer.SetIsSplitLog(RunConf.Global.Log.SplitAccessAndError)
	if RunConf.Global.DebugMode {
		sysadmServer.SetMode(sysadmServer.DebugMode)
	}

	return errs
}

// close access log file descriptor and error log file descriptor
// set AccessLogger  and ErrorLogger to nil
func closeLogger() {
	if RunConf.Global.Log.AccessLogFp != nil {
		fp := RunConf.Global.Log.AccessLogFp
		fp.Close()
		RunConf.Global.Log.AccessLogFp = nil
	}

	if RunConf.Global.Log.ErrorLogFp != nil {
		fp := RunConf.Global.Log.ErrorLogFp
		fp.Close()
		RunConf.Global.Log.ErrorLogFp = nil
	}

}

/*
checkLogLevel check level if is a log level string.
return true if it is a log level string otherwise return false
*/
func checkLogLevel(level string) bool {
	if len(level) < 1 {
		return false
	}

	for _, l := range sysadmServer.Levels {
		if strings.EqualFold(strings.ToLower(level), strings.ToLower(l)) {
			return true
		}
	}

	return false
}

/*
checkLogTimeFormat check valid of log format.
return true if format is a log time format string otherwise return false
*/
func checkLogTimeFormat(format string) bool {
	if len(format) < 1 {
		return false
	}

	for _, v := range sysadmServer.TimestampFormat {
		if strings.EqualFold(format, v) {
			return true
		}
	}

	return false
}

/*
buildGetCommandUrl build complete url address which will be send to a server
*/
func buildUrl(uri string) string {
	var url string = ""

	if strings.TrimSpace(uri) == "" {
		uri = "/"
	}

	svr := RunConf.Global.Server.Address
	port := RunConf.Global.Server.Port
	tls := RunConf.Global.Tls

	if tls.IsTls {
		if port == 443 {
			if uri[0:1] == "/" {
				url = "https://" + svr + uri
			} else {
				url = "https://" + svr + "/" + uri
			}
		} else {
			portStr := strconv.Itoa(port)
			if uri[0:1] == "/" {
				url = "https://" + svr + ":" + portStr + uri
			} else {
				url = "https://" + svr + ":" + portStr + "/" + uri
			}
		}

		return url
	}

	if port == 80 {
		if uri[0:1] == "/" {
			url = "http://" + svr + uri
		} else {
			url = "http://" + svr + "/" + uri
		}
	} else {
		portStr := strconv.Itoa(port)
		if uri[0:1] == "/" {
			url = "http://" + svr + ":" + portStr + uri
		} else {
			url = "http://" + svr + ":" + portStr + "/" + uri
		}
	}

	return url
}

/*
route the command to a handler function which will execute the command according to gotCommand.
c is nil when a command is got in passive mode or  by CLI. otherwise, c should not be nil.
*/
func doRouteCommand(gotCommand *apiserver.CommandData, c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	// if server ask agent to change node identifer
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090123, "debug", "try to run command  %+v", gotCommand))
	isChanged, err := isNodeIdentiferChanged(gotCommand.NodeIdentiferStr)
	errs = append(errs, err...)
	if isChanged {
		newNodeIdentifer, err := apiserver.BuildNodeIdentifer(strings.ToUpper(strings.TrimSpace(gotCommand.NodeIdentiferStr)))
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090215, "error", "change node identifer error %s", err))
			data := make(map[string]interface{}, 0)
			tmpC, err := handleCommandStatus(c, gotCommand, fmt.Sprintf("change node identifer error %s", err), data, apiserver.CommandStatusUnrecognized, false, true)
			if tmpC == nil {
				c = nil
			}
			errs = append(errs, err...)
			logErrors(errs)
			return
		} else {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090216, "debug", "node identifier has be changed "))
			runData.nodeIdentifer = &newNodeIdentifer
		}
	}

	if !apiserver.IsCommandSeqValid(gotCommand.CommandSeq) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090217, "error", "command sequence %s is not valid ", gotCommand.CommandSeq))
		data := make(map[string]interface{}, 0)
		tmpC, err := handleCommandStatus(c, gotCommand, "change node identifer error", data, apiserver.CommandStatusUnrecognized, true, true)
		if tmpC == nil {
			c = nil
		}
		errs = append(errs, err...)
		logErrors(errs)
		return
	}

	switch {
	case strings.Compare(strings.ToLower(gotCommand.Command.Command), "gethostip") == 0:
		runGetHostIP(gotCommand, c)
	case strings.Compare(strings.ToLower(gotCommand.Command.Command), "addyum") == 0:
		runAddyum(gotCommand, c)
	default:
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090266, "error", "command %s is not in the list of agent supported", gotCommand.Command.Command))
		data := make(map[string]interface{}, 0)
		tmpC, err := handleCommandStatus(c, gotCommand, fmt.Sprintf("command %s is not in the list of agent supported", gotCommand.Command.Command), data, 
		apiserver.CommandStatusUnrecognized, false, true)
		if tmpC == nil {
			c = nil
		}
		errs = append(errs, err...)
		logErrors(errs)
	}

}

// sendCommandStatusToServer send the of command status to server actively.
func sendCommandStatusToServer(data []byte) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	if runData.httpClient == nil {
		if err := buildHttpClient(); err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090223, "error", "build http client error %s.", err))
			return false, errs
		}

	}

	url := buildUrl(RunConf.Global.CommandStatusUri)
	requestParams := &httpclient.RequestParams{}
	requestParams.Method = http.MethodPost
	requestParams.Url = url

	body, err := httpclient.NewSendRequest(requestParams, runData.httpClient, strings.NewReader(utils.Bytes2str(data)))
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090226, "error", "can not send the status of command to server %s", err))
		return false, errs
	}

	repStatusData := apiserver.RepStatus{}
	if err := json.Unmarshal(body, &repStatusData); err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090227, "error", "can not unmarshal response message what come from server %s", err))
		return false, errs
	}

	if repStatusData.StatusCode == apiserver.ComandStatusReceived {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090228, "info", "the status of command %s has be send to server", repStatusData.CommandSeq))
		key := defaultRootPathCommandStatus + repStatusData.CommandSeq
		exist, e := redis.Exists(runData.redisEntity, runData.redisctx, key)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090229, "error", "check whether a key is exist in redis server error %s ", e))
		}
		if exist && e == nil {
			e := redis.Del(runData.redisEntity, runData.redisctx, key)
			if e != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090230, "error", "delete data of command status for command(%s) in redis server error %s ", repStatusData.CommandSeq, e))
			} else {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090231, "debug", "delete data of command status for command(%s) in redis server successful", repStatusData.CommandSeq))
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090232, "info", "command status data for command %s has be send to server", repStatusData.CommandSeq))
			}
		} else {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090233, "debug", "data of command status for command(%s) in redis server is not exist or a error occurred %s", repStatusData.CommandSeq, e))
		}

		return true, errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090234, "error", "send command status data for command %s to server error", repStatusData.CommandSeq))
	return false, errs

}

// trySendCommandStatusToServer try to send command status data to server when agent is running in passive mode
// the max try time is defined by sendCommandStatusMaxTryTimes
func trySendCommandStatusToServer(gotCommand *apiserver.CommandData, msg string, data map[string]interface{}, statusCode apiserver.CommandStatusCode, notCommand bool) {
	var errs []sysadmerror.Sysadmerror

	if !RunConf.Agent.Passive {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090222, "debug", "agent can not send command status data to apiserver actively when is running in active mode."))
		logErrors(errs)
		return
	}

	commandStausData, e := apiserver.BuildCommandStatus(gotCommand.CommandSeq, "", msg, *runData.nodeIdentifer, statusCode, data, notCommand)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090224, "error", "build command %s status error %s ", gotCommand.CommandSeq, e))
		logErrors(errs)
		return
	}

	datajson, err := json.Marshal(commandStausData)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090225, "error", "encoding response data to json string error %s", err))
		logErrors(errs)
		return
	}

	for i := 0; i < sendCommandStatusMaxTryTimes; i++ {
		ok, err := sendCommandStatusToServer(datajson)
		if ok {
			errs = append(errs, err...)
			logErrors(errs)
			return
		}

		intervalTimes := math.Pow(2, float64(i))
		errs = append(errs, err...)
		time.Sleep(time.Duration(sendCommandStatusTryInterval*int(intervalTimes)) * time.Second)
	}
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090226, "error", "can not set command status data for command %s to server error", gotCommand.CommandSeq))

	key := defaultRootPathCommandStatus + gotCommand.CommandSeq
	exist, e := redis.Exists(runData.redisEntity, runData.redisctx, key)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090250, "error", "check whether a key is exist in redis server error %s ", e))
		logErrors(errs)
		return
	}

	if exist && e == nil {
		e := redis.Del(runData.redisEntity, runData.redisctx, key)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090251, "error", "delete data of command status for command(%s) in redis server error %s ", gotCommand.CommandSeq, e))
		} else {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090252, "debug", "delete data of command status for command(%s) in redis server successful", gotCommand.CommandSeq))
		}
	} else {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090253, "debug", "data of command status for command(%s) in redis server is not exist or a error occurred %s", gotCommand.CommandSeq, e))
	}

	logErrors(errs)
}

func isNodeIdentiferChanged(nodeIdentifierStr string) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	nodeIdentifierStr = strings.ToUpper(strings.TrimSpace(nodeIdentifierStr))
	if nodeIdentifierStr == "" {
		return false, errs
	}

	oldIdentifier := RunConf.Global.NodeIdentifer
	var oldIP, oldMAC, oldHostName, oldCustomize, newIP, newMAC, newHostName, newCustomize string
	oldSlice := strings.Split(oldIdentifier, "")
	for _, v := range oldSlice {
		switch {
		case strings.Compare(strings.TrimSpace(strings.ToUpper(v)), "IP") == 0:
			oldIP = "IP"
		case strings.Compare(strings.TrimSpace(strings.ToUpper(v)), "MAC") == 0:
			oldMAC = "MAC"
		case strings.Compare(strings.TrimSpace(strings.ToUpper(v)), "HOSTNAME") == 0:
			oldHostName = "HOSTNAME"
		default:
			oldCustomize = v
		}
	}
	oldStr := oldIP + oldMAC + oldHostName + oldCustomize

	newSlice := strings.Split(nodeIdentifierStr, "")
	for _, v := range newSlice {
		switch {
		case strings.Compare(strings.TrimSpace(strings.ToUpper(v)), "IP") == 0:
			newIP = "IP"
		case strings.Compare(strings.TrimSpace(strings.ToUpper(v)), "MAC") == 0:
			newMAC = "MAC"
		case strings.Compare(strings.TrimSpace(strings.ToUpper(v)), "HOSTNAME") == 0:
			newHostName = "HOSTNAME"
		default:
			newCustomize = v
		}
	}
	newStr := newIP + newMAC + newHostName + newCustomize

	if oldStr == newStr {
		return false, errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090214, "debug", "node identifer will be changed to  %s", newStr))
	return true, errs
}

// setCommandStatusIntoRedis build CommandStatus and try to save it into redis server
func setCommandStatusIntoRedis(gotCommand *apiserver.CommandData, msg string, data map[string]interface{}, statusCode apiserver.CommandStatusCode, notCommand bool) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	key := defaultRootPathCommandStatus + gotCommand.CommandSeq
	commandStatus, e := apiserver.BuildCommandStatus(gotCommand.CommandSeq, "", msg, *runData.nodeIdentifer, statusCode, data, notCommand)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090215, "error", "build command status error %s", e))
		return false, errs
	}

	commandStatusMap, e := apiserver.ConvCommandStatus2Map(&commandStatus)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090216, "error", "convert CommandStatus to map error %s", e))
		return false, errs
	}

	e = redis.HSet(runData.redisEntity, runData.redisctx, key, commandStatusMap)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090217, "error", "save command status data to redis server error %s", e))
		return false, errs
	}

	return true, errs
}

// handleCommandStatus handle the following things:
// 1. it try to send command status data to server if agent is running in passive
// 2. it try to save command status data to redis server
// 3. it try to response to server with command status data if command is synchronized when agent is running in active
func handleCommandStatus(c *sysadmServer.Context, gotCommand *apiserver.CommandData, msg string, data map[string]interface{}, statusCode apiserver.CommandStatusCode, notCommand, lastone bool) (*sysadmServer.Context, []sysadmerror.Sysadmerror) {

	var errs []sysadmerror.Sysadmerror

	// data of command statuses are not save to redis server when agent is running in passive mode
	if RunConf.Agent.Passive {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090235, "debug", " try to send status of command  %s to server ", gotCommand.CommandSeq))
		go trySendCommandStatusToServer(gotCommand, msg, data, statusCode, notCommand)
		return nil, errs
	}

	_, err := setCommandStatusIntoRedis(gotCommand, msg, data, statusCode, notCommand)
	errs = append(errs, err...)

	if gotCommand.Synchronized {
		if c != nil && lastone {
			_, err := responseCommandStatusToServer(c, gotCommand, msg, data, statusCode, notCommand, true)
			errs = append(errs, err...)
			return nil, errs
		}
	} else {
		if c != nil {
			_, err := responseCommandStatusToServer(c, gotCommand, msg, data, statusCode, notCommand, false)
			errs = append(errs, err...)
			return nil, errs
		}
	}

	return c, errs
}

// responseCommandStatusToServer response to server with command status data and delete the key named commandSeq from redis server
func responseCommandStatusToServer(c *sysadmServer.Context, gotCommand *apiserver.CommandData, msg string, data map[string]interface{}, statusCode apiserver.CommandStatusCode, notCommand, deleteKey bool) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	if c == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090240, "error", "can not response command status on a nil connection "))
		return false, errs
	}

	commandStausData, e := apiserver.BuildCommandStatus(gotCommand.CommandSeq, "", msg, *runData.nodeIdentifer, statusCode, data, notCommand)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090241, "error", "build command %s status error %s ", gotCommand.CommandSeq, e))
		return false, errs
	}

	c.JSON(http.StatusOK, commandStausData)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090242, "debug", "command status of %s has be send to server ", gotCommand.CommandSeq))

	if !deleteKey {
		return true, errs
	}

	// we should delete data of command status in redis server
	key := defaultRootPathCommandStatus + gotCommand.CommandSeq
	exist, e := redis.Exists(runData.redisEntity, runData.redisctx, key)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090243, "error", "check whether a key is exist in redis server error %s ", e))
	}

	if exist && e == nil {
		e := redis.Del(runData.redisEntity, runData.redisctx, key)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090244, "error", "delete data of command status for command %s error %s ", gotCommand.CommandSeq, e))
	}

	return true, errs

}

// command log message will be save into the redis server as list. the key is commandSeq
// setCommandLogIntoRedis get length of log list from redis server if the key is exist.
// logSeq is the length + 1 if the key is exist. otherwise logSeq is zero
// then setCommandLogIntoRedis set log message into redis server
func setCommandLogIntoRedis(commandSeq, message string, level sysadmLog.Level) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	if !apiserver.IsCommandSeqValid(commandSeq) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090245, "error", "can not set a log message with invalid command sequence(%s) into redis server", commandSeq))
		return false, errs
	}

	key := defaultRootPathCommandLog + commandSeq
	exist, e := redis.Exists(runData.redisEntity, runData.redisctx, key)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090246, "error", "check whether a key is exist in redis server error %s ", e))
		return false, errs
	}

	logSeq := 1
	if exist {
		logLen, e := redis.LLen(runData.redisEntity, runData.redisctx, key)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090247, "error", "get the length of logs list for command %s in redis server error %s", commandSeq, e))
			return false, errs
		}

		logSeq = logLen
	}

	logSeq = logSeq + 1
	log := apiserver.BuildLog(logSeq, message, level)
	logJson, e := json.Marshal(log)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090248, "error", "can not convert log struct to json error: %s", e))
		return false, errs
	}
	
	e = redis.LPush(runData.redisEntity, runData.redisctx, key, logJson)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090249, "error", "can not save log message into redis server error:%s", e))
		return false, errs
	}

	return true, errs
}

// handleCommandLogs get log data from redis server then 
// response the log data to the server if agent is running acitve mode
// otherwise send the log data to the server if agent is running passive mode
func handleCommandLogs(commandSeq string, maxNum, tryTimes int, c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror

	if tryTimes >= sendCommandLogMaxTryTimes {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090268, "error", "times for trying send logs to server is more %d", sendCommandLogMaxTryTimes))
		logErrors(errs)
		return
	}

	if maxNum < 1 || maxNum > maxLogNumPerRequest {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090246, "warning", "max number(%d) of logs is larger than defined ", maxNum))
		maxNum = maxLogNumPerRequest
	}

	if !RunConf.Agent.Passive && c == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090247, "error", "can not reponse logs data to the server on nil"))
		logErrors(errs)
		return
	}

	logs := make([]apiserver.Log,0)
	var logData = apiserver.LogData{
		CommandSeq:     "",
		NodeIdentifier: *runData.nodeIdentifer,
		Logs:           logs,
	}

	if strings.TrimSpace(commandSeq) == "" || !apiserver.IsCommandSeqValid(commandSeq) {
		if RunConf.Agent.Passive {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090248, "error", "try to send command logs to server with invalid command sequence"))
			logErrors(errs)
			return
		}

		logData.CommandSeq = "0000000000000000000"
		logData.EndFlag = true
		logData.NotCommand = true

		_, err := responseCommandLogToServer(c, logData, false)
		errs = append(errs, err...)
		logErrors(errs)

		return
	}

	key := defaultRootPathCommandLog + commandSeq
	exist, e := redis.Exists(runData.redisEntity, runData.redisctx, key)
	if !exist {
		if RunConf.Agent.Passive {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090266, "error", "no command los for command %s ", commandSeq))
			logErrors(errs)
			return
		} else {
			logData.CommandSeq = commandSeq
			logData.EndFlag = true
			logData.NotCommand = true
			_, err := responseCommandLogToServer(c, logData, false)
			errs = append(errs, err...)
			logErrors(errs)
			return
		}
	}

	if e != nil {
		if !RunConf.Agent.Passive {
			logData.CommandSeq = commandSeq
			logData.EndFlag = false
			logData.NotCommand = false
			_, err := responseCommandLogToServer(c, logData, false)
			errs = append(errs, err...)
			logErrors(errs)
			return
		} else {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090267, "error", "an error has occurred %s ", e))
			logErrors(errs)
			intervalTimes := math.Pow(2, float64(tryTimes))
			time.Sleep(time.Duration(sendCommandLogTryInterval*int(intervalTimes)) * time.Second)
			handleCommandLogs(commandSeq, maxNum, (tryTimes + 1), c)
			return
		}
	}

	listLen, _ := redis.LLen(runData.redisEntity, runData.redisctx, key)
	endFlag := false
	if listLen < maxNum {
		maxNum = listLen
		endFlag = true
	}

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

	logData.CommandSeq = commandSeq
	logData.EndFlag = endFlag
	logData.NotCommand = false
	logData.Logs = logs
	logData.Total = total

	if RunConf.Agent.Passive {
		go trySendCommandLogToServer(logData, true)
		return 
	}
	
	_, err := responseCommandLogToServer(c, logData, true)
	errs = append(errs, err...)
	logErrors(errs)

}

// responseCommandStatusToServer response to server with command log data and delete the key named commandSeq from redis server
func responseCommandLogToServer(c *sysadmServer.Context, logData apiserver.LogData, deleteKey bool) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	if c == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090249, "error", "can not response command logs on a nil connection "))
		return false, errs
	}

	c.JSON(http.StatusOK, logData)
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090242, "50090250", "command logs of %s has be send to server ", logData.CommandSeq))

	if !deleteKey {
		return true, errs
	}

	// we should delete data of command status in redis server
	key := defaultRootPathCommandLog + logData.CommandSeq
	exist, e := redis.Exists(runData.redisEntity, runData.redisctx, key)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090251, "error", "check whether a key is exist in redis server error %s ", e))
	}

	if exist && e == nil {
		e := redis.Del(runData.redisEntity, runData.redisctx, key)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090252, "error", "delete data of command log for command %s error %s ", logData.CommandSeq, e))
		}
	}

	return true, errs
}

// sendCommandLogToServer send command logs to server when agent is running in actively mode.
func sendCommandLogToServer(data []byte, deleteKey bool) (bool, []sysadmerror.Sysadmerror) {
	var errs []sysadmerror.Sysadmerror

	if runData.httpClient == nil {
		if err := buildHttpClient(); err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090253, "error", "build http client error %s.", err))
			return false, errs
		}

	}

	url := buildUrl(RunConf.Global.CommandStatusUri)
	requestParams := &httpclient.RequestParams{}
	requestParams.Method = http.MethodPost
	requestParams.Url = url

	body, err := httpclient.NewSendRequest(requestParams, runData.httpClient, strings.NewReader(utils.Bytes2str(data)))
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090254, "error", "can not send the logs of command to server %s", err))
		return false, errs
	}

	repStatusData := apiserver.RepStatus{}
	if err := json.Unmarshal(body, &repStatusData); err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090255, "error", "can not unmarshal response message what come from server %s", err))
		return false, errs
	}

	if repStatusData.StatusCode == apiserver.ComandStatusReceived {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090256, "info", "the logs of command %s has be send to server", repStatusData.CommandSeq))
		if !deleteKey {
			return true, errs
		}

		key := defaultRootPathCommandLog + repStatusData.CommandSeq
		exist, e := redis.Exists(runData.redisEntity, runData.redisctx, key)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090257, "error", "check whether a key is exist in redis server error %s ", e))
		}
		if exist && e == nil {
			e := redis.Del(runData.redisEntity, runData.redisctx, key)
			if e != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090258, "error", "delete data of command status for command(%s) in redis server error %s ", repStatusData.CommandSeq, e))
			} else {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090259, "debug", "delete data of command log for command(%s) in redis server successful", repStatusData.CommandSeq))
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090260, "info", "command logs data for command %s has be send to server", repStatusData.CommandSeq))
			}
		} else {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090261, "debug", "data of command logs for command(%s) in redis server is not exist or a error occurred ", repStatusData.CommandSeq))
		}

		return true, errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090262, "error", "send command logs data for command %s to server error", repStatusData.CommandSeq))

	return false, errs
}

// trySendCommandLogToServer try to send command logs to server when agent is running in passive mode
// the max try time is defined by sendCommandLogMaxTryTimes
func trySendCommandLogToServer(logData apiserver.LogData, deleteKey bool) {
	var errs []sysadmerror.Sysadmerror

	if !RunConf.Agent.Passive {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090263, "debug", "agent can not send command logs to apiserver when it is running in active mode."))
		logErrors(errs)
		return 
	}

	datajson, err := json.Marshal(logData)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090264, "error", "encoding response data to json string error %s", err))
		logErrors(errs)
		return 
	}

	for i := 0; i < sendCommandLogMaxTryTimes; i++ {
		ok, err := sendCommandLogToServer(datajson, deleteKey)
		if ok {
			errs = append(errs, err...)
			logErrors(errs)
			return 
		}

		intervalTimes := math.Pow(2, float64(i))
		errs = append(errs, err...)
		time.Sleep(time.Duration(sendCommandLogTryInterval*int(intervalTimes)) * time.Second)
	}
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090265, "error", "can not send command logs for command to server error", logData.CommandSeq))
	logErrors(errs)
}

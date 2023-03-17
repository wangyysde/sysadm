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
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/wangyysde/sysadm/config"
	"github.com/wangyysde/sysadm/httpclient"
	"github.com/wangyysde/sysadm/sysadmapi/apiutils"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadm/utils"
	"github.com/wangyysde/sysadmServer"
	apiserver "github.com/wangyysde/sysadm/apiserver/app"
	redis "github.com/wangyysde/sysadm/redis"
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
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090215, "error", "change node identifer error %s",err))
			data := make(map[string]interface{},0)
			outPutCommandStatus(c,&gotCommand.Command,fmt.Sprintf("change node identifer error %s",err),data, apiserver.CommandStatusUnrecognized,false)
			logErrors(errs)
			return 
		} else {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090216, "debug", "node identifier has be changed "))
			runData.nodeIdentifer = &newNodeIdentifer
		}
	}
	
	if !apiserver.IsCommandSeqValid(gotCommand.CommandSeq) {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090217, "error", "command sequence %s is not valid ",gotCommand.CommandSeq))
		data := make(map[string]interface{},0)
		outPutCommandStatus(c,&gotCommand.Command,fmt.Sprintf("command sequence %s is not valid ",gotCommand.CommandSeq),data, apiserver.ComandStatusSendError,true)
		logErrors(errs)
		return 
	}

	switch {
	case strings.Compare(strings.ToLower(gotCommand.Command.Command), "gethostip") == 0:
		runGetHostIP(gotCommand.Command, c)
	case strings.Compare(strings.ToLower(gotCommand.Command.Command), "addyum") == 0:
		runAddyum(gotCommand.Command, c)
	default:
		data := make(map[string]interface{},0)
		outPutCommandStatus(c,&gotCommand.Command,fmt.Sprintf("command %s is not in the list of agent supported",strings.ToLower(gotCommand.Command.Command)),data, apiserver.CommandStatusUnrecognized,true )
	}
	
}


func outPutCommandStatus(c *sysadmServer.Context, gotCommand *apiserver.Command, msg string, data map[string]interface{},statusCode apiserver.CommandStatusCode, notCommand bool) {
	var errs []sysadmerror.Sysadmerror

	// that meaning is a response action if c is not nil. so agent must response the server with apiserver.CommandStatus whatever 
	// RunConf.Global.Output is server and apiserver.Command.Synchronized is true
	if c != nil {
		commandStausData, e := apiserver.BuildCommandStatus(gotCommand.CommandSeq,"",msg,*runData.nodeIdentifer,statusCode,data,notCommand)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090218, "error", "build command %s status error %s ",gotCommand.CommandSeq,e))
			logErrors(errs)
			return 
		}
		c.JSON(http.StatusOK,commandStausData)
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090219, "debug", "command status of %s has be send to server ",gotCommand.CommandSeq))
		// we should delete data of command status in redis server if agent runing in passive or command is asynchronous
		if RunConf.Agent.Passive || !gotCommand.Synchronized {
			key := defaultRootPathCommandStatus + gotCommand.CommandSeq
			e := redis.Del(runData.redisEntity,runData.redisctx,key)
			if e != nil {
				errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090220, "error", "delete data of command status for command %s error %s ",gotCommand.CommandSeq,e))
			}
		}
		logErrors(errs)
		return 
	}

	sendCommandStatusToServer(gotCommand, msg, data, statusCode, notCommand)

}


//  sendCommandStatusToServer send the of command status to server actively. 
func sendCommandStatusToServer(gotCommand *apiserver.Command, msg string, data map[string]interface{},statusCode apiserver.CommandStatusCode, notCommand bool) {
	var errs []sysadmerror.Sysadmerror

	if !RunConf.Agent.Passive {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090220, "error", "agent can not send command status data to apiserver actively when is running in active mode."))
		logErrors(errs)
		return 
	} 

	if runData.httpClient == nil {
		if err := buildHttpClient(); err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090221, "error", "build http client error %s.", err))
			logErrors(errs)
			return 
		}

	}

	url := buildUrl(RunConf.Global.CommandStatusUri)
	requestParams := &httpclient.RequestParams{}
	requestParams.Method = http.MethodPost
	requestParams.Url = url

	commandStausData, e := apiserver.BuildCommandStatus(gotCommand.CommandSeq,"",msg,*runData.nodeIdentifer,statusCode,data,notCommand)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090222, "error", "build command %s status error %s ",gotCommand.CommandSeq,e))
		logErrors(errs)
		return 
	}

	datajson, err := json.Marshal(commandStausData)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090223, "error", "encoding response data to json string error %s", err))
		logErrors(errs)
		return
	}

	body, err := httpclient.NewSendRequest(requestParams, runData.httpClient, strings.NewReader(utils.Bytes2str(datajson)))
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090224, "error", "can not send the status of command to server %s", err))
		logErrors(errs)
		return
	}

	repStatusData := apiserver.RepStatus{}
	if err := json.Unmarshal(body, &repStatusData); err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090205, "error", "can not unmarshal response message what come from server %s", err))
		logErrors(errs)
		return
	}

	if repStatusData.StatusCode == apiserver.ComandStatusReceived {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090206, "info", "the status of command %s has be send to server", repStatusData.CommandSeq))
		key := defaultRootPathCommandStatus + repStatusData.CommandSeq
		e := redis.Del(runData.redisEntity,runData.redisctx,key)
		if e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090207, "error", "delete command status data for command %s from redis server error %s", repStatusData.CommandSeq,e))
			logErrors(errs)
			return
		}
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090208, "debug", "command status data for command %s has be send to server", repStatusData.CommandSeq))
		logErrors(errs)
		return
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090209, "error", "send command status data for command %s to server error", repStatusData.CommandSeq))
	logErrors(errs)
	return

}

func isNodeIdentiferChanged(nodeIdentifierStr string) (bool, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	nodeIdentifierStr = strings.ToUpper(strings.TrimSpace(nodeIdentifierStr))
	if nodeIdentifierStr == ""{
		return false, errs
	}

	oldIdentifier := RunConf.Global.NodeIdentifer
	var oldIP,oldMAC,oldHostName,oldCustomize,newIP, newMAC, newHostName, newCustomize string 
	oldSlice := strings.Split(oldIdentifier,"")
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

	newSlice := strings.Split(nodeIdentifierStr,"")
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

	if oldStr ==  newStr {
		return false, errs
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090214, "debug", "node identifer will be changed to  %s", newStr))
	return true, errs
}

// setCommandStatusIntoRedis build CommandStatus and try to save it into redis server
func setCommandStatusIntoRedis(gotCommand *apiserver.Command,msg string, data map[string]interface{},statusCode apiserver.CommandStatusCode, notCommand bool)(bool, []sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	key := defaultRootPathCommandStatus + gotCommand.CommandSeq
	commandStatus, e := apiserver.BuildCommandStatus(gotCommand.CommandSeq,"",msg,*runData.nodeIdentifer,statusCode,data,notCommand)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090215, "error", "build command status error %s", e))
		return false, errs
	}
	
	commandStatusMap, e := apiserver.ConvCommandStatus2Map(&commandStatus)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090216, "error", "convert CommandStatus to map error %s", e))
		return false, errs
	}
	
	e = redis.HSet(runData.redisEntity,runData.redisctx,key,commandStatusMap)
	if e != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090217, "error", "save command status data to redis server error %s", e))
		return false, errs
	}

	return true, errs
}
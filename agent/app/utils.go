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

func getNodeIdentifer(confNodeIdentifer string) (*NodeIdentifer, error) {
	if strings.TrimSpace(confNodeIdentifer) == "" {
		confNodeIdentifer = "IP,HOSTNAME,MAC"
	}

	ret := NodeIdentifer{}

	identiferSlice := strings.Split(confNodeIdentifer, ",")
	isCustomize := true

	for _, value := range identiferSlice {
		switch {
		case strings.ToUpper(strings.TrimSpace(value)) == "IP":
			ips, err := utils.GetLocalIPs()
			if err != nil {
				return nil, fmt.Errorf("get local host ip address error %s", err)
			}
			ret.Ips = ips
			isCustomize = false
		case strings.ToUpper(strings.TrimSpace(value)) == "MAC":
			macs, err := utils.GetLocalMacs()
			if err != nil {
				return nil, fmt.Errorf("get local host mac information error %s", err)
			}
			ret.Macs = macs
			isCustomize = false
		case strings.ToUpper(strings.TrimSpace(value)) == "HOSTNAME":
			hostname, err := os.Hostname()
			if err != nil {
				return nil, fmt.Errorf("can not get hostname %s", err)
			}
			ret.Hostname = hostname
			isCustomize = false
		default:
			if strings.TrimSpace(value) != "" {
				ret.Customize = strings.TrimSpace(value)
				isCustomize = true
			} else {
				return nil, fmt.Errorf("node identifer %s is not valid", value)
			}
		}
	}

	if strings.TrimSpace(confNodeIdentifer) != "" && isCustomize {
		ret.Customize = strings.TrimSpace(confNodeIdentifer)
	}

	return &ret, nil
}

/*
buildGetCommandUrl build complete url address which will be send to a server
*/
func buildGetCommandUrl(newUri string) string {
	var url string = ""

	if strings.TrimSpace(RunConf.Global.Uri) == "" {
		RunConf.Global.Uri = "/"
	}

	var uri string = ""
	if strings.TrimSpace(newUri) == "" {
		uri = strings.TrimSpace(RunConf.Global.Uri)
	} else {
		uri = strings.TrimSpace(newUri)
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
func doRouteCommand(gotCommand *Command, c *sysadmServer.Context) {
	var errs []sysadmerror.Sysadmerror
	// if server ask agent to change node identifer
	errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090123, "debug", "get command from server %+v", gotCommand))
	logErrors(errs)
	errs = errs[0:0]
	if strings.TrimSpace(gotCommand.NodeIdentifer) != "" && !strings.EqualFold(strings.ToLower(gotCommand.NodeIdentifer), strings.ToLower(RunConf.Global.NodeIdentifer)) {
		RunConf.Global.NodeIdentifer = strings.ToLower(gotCommand.NodeIdentifer)
		nodeIdentifer, err := getNodeIdentifer(gotCommand.NodeIdentifer)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090109, "error", "get node identifier error %s", err))
			logErrors(errs)
			return
		}
		runData.nodeIdentifer = nodeIdentifer
	}

	switch {
	case strings.Compare(strings.ToLower(gotCommand.Command), "gethostip") == 0:
		runGetHostIP(gotCommand, c)
	}
}

func outPutResult(c *sysadmServer.Context, gotCommand *Command, data apiutils.ApiResponseData) {
	// it is that agent running in active mode if c is not nil,then we should response the client with data
	if c != nil {

		nodeIdentifer := runData.nodeIdentifer
		repData := map[string]interface{}{
			"nodeIdentifer": *nodeIdentifer,
			"data":          data,
		}
		c.JSON(http.StatusOK, repData)
		return
	}

	outPut := RunConf.Global.Output
	switch {
	case strings.Compare(strings.TrimSpace(strings.ToLower(outPut)), "server") == 0:
		outPutResultToServer(gotCommand, data)
	case strings.Compare(strings.TrimSpace(strings.ToLower(outPut)), "file") == 0:
		// if outputFile is not set,the we will output the result to server
		if RunConf.Global.OutputFile == "" {
			outPutResultToServer(gotCommand, data)
		} else {
			outPutResultToFile(gotCommand, data)
		}
	case strings.Compare(strings.TrimSpace(strings.ToLower(outPut)), "stdout") == 0:
		outPutResultToStdout(gotCommand, data)
	default:
		outPutResultToServer(gotCommand, data)
	}

}

func outPutMessage(c *sysadmServer.Context, gotCommand *Command, msg string, data ...interface{}) {
	var msgData []map[string]interface{}

	m := fmt.Sprintf(msg, data...)
	ret := map[string]interface{}{
		"msg": m,
	}
	msgData = append(msgData, ret)
	retData := apiutils.NewBuildResponseDataForMap(false, 50090201, msgData)
	outPutResult(c, gotCommand, retData)

}

/*
outPutResultToServer output the result of a command execution to the server when agent running in passive mode or CLI.
agent response the result to the server in outPutResult function if it is running active mode.
*/
func outPutResultToServer(gotCommand *Command, data apiutils.ApiResponseData) {
	var errs []sysadmerror.Sysadmerror

	if runData.httpClient == nil {
		if err := buildHttpClient(); err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090202, "fatal", "build http client error %s.", err))
			logErrors(errs)
			os.Exit(22)
		}

	}

	var url string = ""
	if strings.TrimSpace(gotCommand.ResponseUri) != "" {
		url = buildGetCommandUrl(strings.TrimSpace(gotCommand.ResponseUri))
	} else {
		url = buildGetCommandUrl(strings.TrimSpace(RunConf.Global.Uri))
	}

	requestParams := &httpclient.RequestParams{}
	requestParams.Method = http.MethodPost
	requestParams.Url = url

	nodeIdentifer := runData.nodeIdentifer
	repData := map[string]interface{}{
		"nodeIdentifer": *nodeIdentifer,
		"data":          data,
	}

	datajson, err := json.Marshal(repData)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090203, "error", "encoding response data to json string error %s", err))
		logErrors(errs)
		return
	}

	body, err := httpclient.NewSendRequest(requestParams, runData.httpClient, strings.NewReader(utils.Bytes2str(datajson)))
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090204, "error", "can not send the result of command execution to server %s", err))
		logErrors(errs)
		return
	}

	responseData := apiutils.ApiResponseData{}
	if err := json.Unmarshal(body, &responseData); err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090205, "error", "can not unmarshal response message what come from server %s", err))
		logErrors(errs)
		return
	}

	if responseData.Status {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090206, "info", "the result of the command execution has be sent to the server"))
		logErrors(errs)
		return
	}

	if len(responseData.Message) > 0 {
		msgSlice := responseData.Message
		msg := msgSlice[1]
		errMsg, okMsg := msg["msg"]
		if okMsg {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090207, "error", "send the result of the command execution error %s", errMsg))
			logErrors(errs)
			return
		}
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090208, "error", "unknow error has occurred when send the result of command execution to server"))
	logErrors(errs)
}

func outPutResultToFile(gotCommand *Command, data apiutils.ApiResponseData) {
	var errs []sysadmerror.Sysadmerror

	if strings.TrimSpace(RunConf.Global.OutputFile) == "" {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090209, "warning", "output has be set to file, but outputFile is not exist. we will output the result to stdout."))
		logErrors(errs)
		outPutResultToStdout(gotCommand, data)
		return
	}

	nodeIdentifer := runData.nodeIdentifer
	repData := map[string]interface{}{
		"nodeIdentifer": *nodeIdentifer,
		"data":          data,
	}

	datajson, err := json.Marshal(repData)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090210, "error", "encoding response data to json string error %s", err))
		logErrors(errs)
		return
	}

	fp, err := os.OpenFile(RunConf.Global.OutputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090214, "error", "can not open output file %s error %s", RunConf.Global.OutputFile, err))
		logErrors(errs)
		return
	}
	n, err := fp.Write(datajson)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090211, "error", "write the result of command to file %s error %s", RunConf.Global.OutputFile, err))
		logErrors(errs)
		return
	}

	errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090212, "info", "write the result of command to file %s successful. total bytes %d has be write to file", RunConf.Global.OutputFile, n))
	logErrors(errs)
}

func outPutResultToStdout(gotCommand *Command, data apiutils.ApiResponseData) {
	var errs []sysadmerror.Sysadmerror

	nodeIdentifer := runData.nodeIdentifer
	repData := map[string]interface{}{
		"nodeIdentifer": *nodeIdentifer,
		"data":          data,
	}

	datajson, err := json.Marshal(repData)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(50090213, "error", "encoding response data to json string error %s", err))
		logErrors(errs)
		return
	}

	fmt.Printf("%s", datajson)
}

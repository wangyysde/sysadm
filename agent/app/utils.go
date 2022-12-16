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
	"os"
	"strconv"
	"strings"

	"github.com/wangyysde/sysadm/config"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadm/utils"
	"github.com/wangyysde/sysadmServer"
)

func SetVersion(version *config.Version){
	if version == nil {
		return
	}

	version.Version = ver
	version.Author = author

	CliOps.Version = *version 
	RunConf.Version =  *version
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
func logErrors(errs []sysadmerror.Sysadmerror){

	for _,e := range errs {
		l := sysadmerror.GetErrorLevelString(e)
		no := e.ErrorNo
		sysadmServer.Logf(l,"erroCode: %d Msg: %s",no,e.ErrorMsg)
	}
	
}


// set parameters to accessLogger and errorLooger
func setLogger()([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	sysadmServer.SetLoggerKind(RunConf.Global.Log.Kind)
	sysadmServer.SetLogLevel(RunConf.Global.Log.Level)
	sysadmServer.SetTimestampFormat(RunConf.Global.Log.TimeStampFormat)
	_,fp,err := sysadmServer.SetAccessLogFile(RunConf.Global.Log.AccessLog)
	if err != nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081002,"error","can not set access log file(%s) error: %s",RunConf.Global.Log.AccessLog,err))
	}else{
		RunConf.Global.Log.AccessLogFp = fp
	}
	
	if RunConf.Global.Log.SplitAccessAndError && RunConf.Global.Log.ErrorLog != "" {
		_,fp,err := sysadmServer.SetErrorLogFile(RunConf.Global.Log.ErrorLog)
		if err != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(10081003,"error","can not set error log file(%s) error: %s",RunConf.Global.Log.ErrorLog ,err))
		}else{
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
func closeLogger(){
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

	for _,l := range sysadmServer.Levels {
		if strings.EqualFold(strings.ToLower(level),strings.ToLower(l)) {
			return true
		}
	}

	return false
}

/*
   checkLogTimeFormat check valid of log format.
   return true if format is a log time format string otherwise return false
*/
func checkLogTimeFormat(format string) bool{
	if len(format) < 1 {
		return false
	}

	for _,v := range sysadmServer.TimestampFormat {
		if strings.EqualFold(format,v) {
			return true
		}
	}

	return false
}

func getNodeIdentifer(confNodeIdentifer string) (*NodeIdentifer, error){
	if strings.TrimSpace(confNodeIdentifer) == "" {
		confNodeIdentifer = "IP,HOSTNAME,MAC"
	}

	ret :=  NodeIdentifer{}

	identiferSlice :=  strings.Split(confNodeIdentifer,",")
	isCustomize := true

	for _, value := range identiferSlice {
		switch {
		case strings.ToUpper(strings.TrimSpace(value)) == "IP":
			ips,err :=  utils.GetLocalIPs()
			if err != nil  {
				return nil, fmt.Errorf("get local host ip address error %s", err)
			}
			ret.Ips = ips
			isCustomize = false
		case strings.ToUpper(strings.TrimSpace(value)) == "MAC":
			macs,err :=  utils.GetLocalMacs()
			if err != nil  {
				return nil, fmt.Errorf("get local host mac information error %s", err)
			}
			ret.Macs = macs
			isCustomize = false
		case strings.ToUpper(strings.TrimSpace(value)) == "HOSTNAME":
			hostname, err :=  os.Hostname()
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
func buildGetCommandUrl() string {
	var url string = ""

	if strings.TrimSpace(RunConf.Global.Uri) == "" {
		RunConf.Global.Uri = "/"
	}

	uri := strings.TrimSpace(RunConf.Global.Uri)
	prefixStr := uri[0:5]
	if strings.Compare(strings.ToUpper(prefixStr),"HTTP:") == 0 || strings.Compare(strings.ToUpper(prefixStr),"HTTPS") == 0 {
		if strings.Compare(strings.ToUpper(prefixStr),"HTTP:") == 0 {
			tls := RunConf.Global.Tls
			if tls.IsTls {
				tmpUri := uri[5:(len(uri)-1)]
				url = "https" + tmpUri
			} else {
				url = uri
			}
		}

		if  strings.Compare(strings.ToUpper(prefixStr),"HTTPS") == 0 {
			tls := RunConf.Global.Tls
			if ! tls.IsTls {
				tmpUri := uri[6:(len(uri)-1)]
				url = "http" + tmpUri
			} else {
				url = uri
			}
		}
	} else {
		svr := RunConf.Global.Server.Address
		port := RunConf.Global.Server.Port
		tls := RunConf.Global.Tls
		if tls.IsTls {
			if port == 443 {
				if uri[0:1] == "/" {
					url = "https://" +  svr + uri 
				} else {
					url = "https://" +  svr + "/" + uri 
				}
			} else {
				portStr := strconv.Itoa(port)
				if uri[0:1] == "/" {
					url = "https://" + svr + ":" + portStr + uri
				} else {
					url = "https://" + svr + ":" + portStr + "/" + uri
				}
			}
		} else {
			if port == 80 {
				if uri[0:1] == "/" {
					url = "http://" +  svr + uri 
				} else {
					url = "http://" +  svr + "/" + uri 
				}
			} else {
				portStr := strconv.Itoa(port)
				if uri[0:1] == "/" {
					url = "http://" + svr + ":" + portStr + uri
				} else {
					url = "http://" + svr + ":" + portStr + "/" + uri
				}
			}
		}
	}
	
	return url
}

func handleHTTPBody(body []byte) {
	
}
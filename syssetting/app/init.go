/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
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
	"strings"

	"github.com/wangyysde/sysadm/config"
	sysadmDB "github.com/wangyysde/sysadm/db"
	"github.com/wangyysde/sysadm/sysadmerror"
	log "github.com/wangyysde/sysadmLog"
	"github.com/wangyysde/sysadmServer"
)

func (s SysSetting)GetModuleName() string{
	if strings.TrimSpace(s.ModuleName) == "" {
		s.ModuleName =  moduleName
	}

	return s.ModuleName
}

// do nothing for match old api interface request
func (s SysSetting)GetActionList()[]string{
	return []string{}
}

/*
	set dbConfig(*sysadmDB.DbConfig) and working root path to the global variable WorkingData
	the value of variable are not be instead of by the new values if them are not empty or nil. 
*/
func (s SysSetting)SetWorkingData(dbConf *sysadmDB.DbConfig, logConf *config.Log, workingRoot string)(errs []sysadmerror.Sysadmerror){
	if WorkingData.dbConf ==  nil {
		if dbConf == nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(4010001,"fatal","Can not set DB Conf to working data with nil" ))
			return errs
		}
		WorkingData.dbConf = dbConf
	}

	if WorkingData.logConf == nil {
		if logConf == nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(4010002,"fatal","Can not set Log Conf to working data with nil" ))
			return errs
		}
		WorkingData.logConf = logConf

		e := setLogger()
		errs = append(errs,e...)
	}
	
	if strings.TrimSpace(WorkingData.workingRoot) == "" {
		if strings.TrimSpace(workingRoot) == "" {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(4010003,"warn","working root path is empty" ))
		}

		WorkingData.workingRoot = workingRoot
	}

	return errs
}

// New an instance of Infrastructure
func NewSyssetting() *SysSetting{
	return &SysSetting{
		ModuleName: moduleName,
		ApiVersion: apiVersion,
	}
}


// set the value of global varible in SysadmServer with logger
func setLogger() (errs []sysadmerror.Sysadmerror) {
	logConf :=  WorkingData.logConf
	if strings.TrimSpace(logConf.Kind) != "" {
		sysadmServer.SetLoggerKind(logConf.Kind)
	}

	if strings.TrimSpace(logConf.Level) != "" {
		sysadmServer.SetLogLevel(logConf.Level)
	}
	
	if strings.TrimSpace(logConf.TimeStampFormat) != "" {
		sysadmServer.SetTimestampFormat(logConf.TimeStampFormat)
	}

	if logConf.AccessLogFp != nil {
		if e := sysadmServer.SetAccessLoggerWithFp(logConf.AccessLogFp); e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(4010004,"error","set access logger error %s", e ))
		}
	}

	if logConf.ErrorLogFp != nil {
		if e := sysadmServer.SetErrorLoggerWithFp(logConf.ErrorLogFp); e != nil {
			errs = append(errs, sysadmerror.NewErrorWithStringLevel(4010005,"error","set error logger error %s", e ))
		}
	}

	sysadmServer.SetIsSplitLog(logConf.SplitAccessAndError)

	level, e := log.ParseLevel(logConf.Level)
	if e != nil {
		sysadmServer.SetMode(sysadmServer.DebugMode)
	}else {
		if level >= log.DebugLevel {
			sysadmServer.SetMode(sysadmServer.DebugMode)
		}else{
			sysadmServer.SetMode(sysadmServer.ReleaseMode)
		}
	}

	return errs
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
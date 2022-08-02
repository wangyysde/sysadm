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

package server

import (
	"github.com/wangyysde/sysadm/config"
	"github.com/wangyysde/sysadm/infrastructure/app"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

func AddInfrastructureHandlers(r *sysadmServer.Engine)([]sysadmerror.Sysadmerror){
	var errs []sysadmerror.Sysadmerror

	if r == nil {
		errs = append(errs, sysadmerror.NewErrorWithStringLevel(700020001,"fatal","can not add handlers to nil" ))
		return errs
	}

	dbConf := RuntimeData.RuningParas.DBConfig 
	logConf :=  config.Log{}
	definedConfig := RuntimeData.RuningParas.DefinedConfig
	logConf.AccessLog = definedConfig.Log.AccessLog
	logConf.AccessLogFp = RuntimeData.RuningParas.AccessLogFp
	logConf.ErrorLog = definedConfig.Log.ErrorLog
	logConf.ErrorLogFp = RuntimeData.RuningParas.ErrorLogFp
	logConf.Kind = definedConfig.Log.Kind
	logConf.Level = definedConfig.Log.Level
	logConf.SplitAccessAndError = definedConfig.Log.SplitAccessAndError
	logConf.TimeStampFormat = definedConfig.Log.TimeStampFormat
	workingRoot := RuntimeData.StartParas.SysadmRootPath 

	infrastructure := app.NewInfrastructure()
	err := infrastructure.SetWorkingData(dbConf,&logConf,workingRoot)
	errs = append(errs,err...)
	if sysadmerror.GetMaxLevel(errs) >= sysadmerror.GetLevelNum("fatal"){
		return errs
	}

	err = infrastructure.AddHandlers(r)
	errs = append(errs,err...)
	
	
	return errs
}


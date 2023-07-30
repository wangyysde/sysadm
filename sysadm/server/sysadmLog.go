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
	"sysadm/sysadmLog"
)

// set parameters to accessLogger and errorLooger
func setSysadmLogger() {
	logger := sysadmLog.NewSysadmLogger()
	logger.SetLoggerKind(RuntimeData.RuningParas.DefinedConfig.Log.Kind)
	logger.SetLoggerLevel(RuntimeData.RuningParas.DefinedConfig.Log.Level)
	logger.SetTimestampFormat(RuntimeData.RuningParas.DefinedConfig.Log.TimeStampFormat)
	if RuntimeData.RuningParas.DefinedConfig.Log.AccessLog != "" {
		err := logger.SetAccessLogFile(RuntimeData.RuningParas.DefinedConfig.Log.AccessLog)
		if err != nil {
			logger.Logf("error", "%s", err)
		}
	}

	if RuntimeData.RuningParas.DefinedConfig.Log.SplitAccessAndError && RuntimeData.RuningParas.DefinedConfig.Log.ErrorLog != "" {
		err := logger.SetErrorLogFile(RuntimeData.RuningParas.DefinedConfig.Log.ErrorLog)
		if err != nil {
			logger.Logf("error", "%s", err)
		}
	}
	logger.SetIsSplitLog(RuntimeData.RuningParas.DefinedConfig.Log.SplitAccessAndError)

	RuntimeData.sysadmLogEntity = logger
}
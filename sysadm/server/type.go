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
	"os"

	"github.com/wangyysde/sysadmServer"
	sysadmDB "sysadm/db"
	"sysadm/sysadm/config"
	"sysadm/sysadmLog"
)

// Start parameters of the program
type StartParas struct {
	// Point to configuration file path of server
	ConfigPath string
	// root path of sysadm executeable package
	SysadmRootPath string
}

type RunningParas struct {
	// descriptor of access log file which will be used to close logger when system exit
	AccessLogFp *os.File
	// descriptor of error log file which will be used to close logger when system exit
	ErrorLogFp *os.File
	// the configurations what have been parsed. these configurations come from environment, configure file or default value
	DefinedConfig *config.Config
	// the DB configurations what have been parsed. these configurations come from environment, configure file or default value
	DBConfig *sysadmDB.DbConfig
}

type RuningData struct {
	// the start parameters of program such as configuration file , root path of program
	StartParas *StartParas
	// the runing parameters of program what come from environment, configure file or default value
	RuningParas *RunningParas
	// the old start parameters of program after the program reload or restart
	OldStartParas *StartParas
	// the old runing parameters of program after the program reload or restart
	OldRunningParas *RunningParas
	// log entity
	sysadmLogEntity *sysadmLog.LoggerConfig
}

type actionHandler struct {
	name         string
	templateFile string
	handler      sysadmServer.HandlerFunc
	method       []string
}

var RuntimeData RuningData = RuningData{StartParas: &StartParas{}, RuningParas: &RunningParas{}, OldStartParas: nil, OldRunningParas: nil, sysadmLogEntity: nil}
var CliData StartParas = StartParas{ConfigPath: ""}

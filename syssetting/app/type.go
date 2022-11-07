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
	"github.com/wangyysde/sysadm/config"
	sysadmDB "github.com/wangyysde/sysadm/db"
	"github.com/wangyysde/sysadm/sysadmerror"
	"github.com/wangyysde/sysadmServer"
)

// Saving running data of an instance
type RuningData struct {
	dbConf *sysadmDB.DbConfig
	logConf *config.Log
	workingRoot string
}

// Initating working data for an instance
var WorkingData RuningData = RuningData{
	dbConf: nil,
	logConf: nil,
	workingRoot: "",
}

type handlerAdder func (*sysadmServer.Engine, string, SysSetting)([]sysadmerror.Sysadmerror)

// SysSetting
type SysSetting struct {
	ModuleName string 
	ApiVersion string 
}

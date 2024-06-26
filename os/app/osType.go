/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2023 Bzhy Network. All rights reserved.
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
	sessions "github.com/wangyysde/sysadmSessions"
	sysadmDB "sysadm/db"
	sysadmObjects "sysadm/objects/app"
	"sysadm/sysadmLog"
	sysadmSetting "sysadm/syssetting/app"
)

type OS struct {
	// object name. the value should be "os"
	Name string
	// table name which hold host data in DB
	TableName string
	// field name of primary key in the table
	PkName string
}

// 操作系统信息表结构
type OSSchema struct {
	// 操作系统发行版ID
	OSID int `form:"osID" json:"osID" yaml:"osID" xml:"id" db:"osID"`
	// 操作系统发行版名称，如centos,readhat, ubantu等，不区分大小写
	Name string `form:"name" json:"name" yaml:"name" xml:"name" db:"name"`
	// 体系统架构
	Architecture string `form:"architecture" json:"architecture" yaml:"architecture" xml:"architecture" db:"architecture"`
	// 位数
	Bit int `form:"bit" json:"bit" yaml:"bit" xml:"bit" db:"bit"`
	// 描述
	Description string `form:"description" json:"description" yaml:"description" xml:"description" db:"description"`
}

// 存储运行期数据
type runingData struct {
	dbConf        *sysadmDB.DbConfig
	logEntity     *sysadmLog.LoggerConfig
	workingRoot   string
	sessionName   string
	sessionOption sessions.Options
	pageInfo      sysadmSetting.PageInfo
	objectEntiy   sysadmObjects.ObjectEntity
}

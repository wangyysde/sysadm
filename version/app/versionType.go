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

type Version struct {
	// object name. the value should be "version"
	Name string
	// table name which hold host data in DB
	TableName string
	// field name of primary key in the table
	PkName string
}

// 版本信息表结构
type VersionSchema struct {
	// 版本ID
	VersionID int `form:"versionID" json:"osID" yaml:"versionID" xml:"id" db:"versionID"`
	// 版本名称
	Name string `form:"name" json:"name" yaml:"name" xml:"name" db:"name"`
	// 操作系统ID
	OSID int `form:"osid" json:"osid" yaml:"osid" xml:"osid" db:"osid"`
	// 版本类型ID ，对应见常量定义
	TypeID int `form:"typeID" json:"typeID" yaml:"typeID" xml:"typeID" db:"typeID"`
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

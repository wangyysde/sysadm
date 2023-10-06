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

type Yum struct {
	// object name. the value should be "version"
	Name string
	// table name which hold host data in DB
	TableName string
	// field name of primary key in the table
	PkName string
}

// Yum信息表结构
type YumSchema struct {
	// Yum ID
	YumID uint `form:"yumid" json:"yumid" yaml:"yumid" xml:"id" db:"yumid"`
	// Yum名称，作为yum的标识符
	Name string `form:"name" json:"name" yaml:"name" xml:"name" db:"name"`
	// 操作系统ID
	OSID int `form:"osid" json:"osid" yaml:"osid" xml:"osid" db:"osid"`
	// 版本版本ID ，
	VersionID int `form:"versionid" json:"versionid" yaml:"versionid" xml:"versionid" db:"versionid"`
	// what type of the yum is it,such as os, docker, kubernetes,对应类型定义
	TypeID int `form:"typeid" json:"typeid" yaml:"typeid" xml:"typeid" db:"typeid"`
	// which catalog  of the yum is it,such as base, update,plus,......
	Catalog string `form:"catalog" json:"catalog" yaml:"catalog" xml:"catalog" db:"catalog"`
	// kind of the yum,such as local,remote.
	Kind string `form:"kind" json:"kind" yaml:"kind" xml:"kind" db:"kind"`
	// the url of yum if its kind is remote
	BaseUrl string `form:"base_url" json:"base_url" yaml:"base_url" xml:"base_url" db:"base_url"`
	// whether enabled this yum. 1 for enabled, otherwise for disabled
	Enabled int `form:"enabled" json:"enabled" yaml:"enabled" xml:"enabled" db:"enabled"`
	// whether gpg check. 1 for check, otherwise for not check
	Gpgcheck int `form:"gpgcheck" json:"gpgcheck" yaml:"gpgcheck" xml:"gpgcheck" db:"gpgcheck"`
	// the path of the gpgkey file for local, the url of the gpgkey file for remote
	Gpgkey string `form:"gpgkey" json:"gpgkey" yaml:"gpgkey" xml:"gpgkey" db:"gpgkey"`
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

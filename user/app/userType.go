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

type User struct {
	// object name. the value should be "user"
	Name string
	// table name which hold project data in DB
	TableName string
	// field name of primary key in the table
	PkName string
}

// 用户表结构
type UserSchema struct {
	// 用户ID，自增长
	Id string `form:"userid" json:"userid" yaml:"userid" xml:"userid" db:"userid"`
	// 用户名，最长不超过255个字符
	Username string `form:"username" json:"username" yaml:"username" xml:"username" db:"username"`
	// 用户的电子邮件地址，最长不超过255个字符
	Email string `form:"email" json:"email" yaml:"email" xml:"email" db:"email"`
	// 用户帐号密码，最长不超长40个字符
	Password string `form:"password" json:"password" yaml:"password" xml:"password" db:"password"`
	// 用户真实名字,最长不超过255个字符
	Realname string `form:"realname" json:"realname" yaml:"realname" xml:"realname" db:"realname"`
	// 描述，最找不超过255个字符
	Comment string `form:"comment" json:"comment" yaml:"comment" xml:"comment" db:"comment"`
	// 删除标识，0表示正常， 1表示已删除
	Deleted int `form:"deleted" json:"deleted" yaml:"deleted" xml:"deleted" db:"deleted"`
	// 上次用户被重置密码的操作者对应的ID
	ResetUuid string `form:"reset_uuid" json:"reset_uuid" yaml:"reset_uuid" xml:"reset_uuid" db:"reset_uuid"`
	// 加密密码的Salt
	Salt string `form:"salt" json:"salt" yaml:"salt" xml:"salt" db:"salt"`
	// 是否是管理员的标识 0表示普通用户 1表示是管理员
	SysadminFlag int `form:"sysadmin_flag" json:"sysadmin_flag" yaml:"sysadmin_flag" xml:"sysadmin_flag" db:"sysadmin_flag"`
	// 用户创建的时间截
	CreationTime int `form:"creation_time" json:"creation_time" yaml:"creation_time" xml:"creation_time" db:"creation_time"`
	// 用户更新的时间截
	UpdateTime int `form:"update_time" json:"update_time" yaml:"update_time" xml:"update_time" db:"update_time"`
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

/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2024 Bzhy Network. All rights reserved.
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

type Event struct {
	// object name. the value should be "event"
	Name string
	// table name which hold project data in DB
	TableName string
	// field name of primary key in the table
	PkName string
}

// 事件信息表结构
type EventSchema struct {
	// 事件ID，自增
	Id int `form:"id" json:"id" yaml:"id" xml:"id" db:"id"`
	// 事件类别
	Class int `form:"class" json:"class" yaml:"class" xml:"class" db:"class"`
	// 产生事件对象的范围,定义见eventDefined.go文件
	Scope int `form:"scope" json:"scope" yaml:"scope" xml:"scope" db:"scope"`
	// 事件发生时的时间截
	StartTime int `form:"startTime" json:"startTime" yaml:"startTime" xml:"startTime" db:"startTime"`
	// 发生事件的原因
	ReasonMessage string `form:"reasonMessage" json:"reasonMessage" yaml:"reasonMessage" xml:"reasonMessage" db:"reasonMessage"`
	// 产生事件的对象并及其ID或标识
	Object string `form:"object" json:"object" yaml:"object" xml:"object" db:"object"`
	// 产生事件如果有子对象，则本字段记录子对象及它的ID或标识
	SubObject string `form:"subObject" json:"subObject" yaml:"subObject" xml:"subObject" db:"subObject"`
	// 事件发生时的动作
	Action string `form:"action" json:"action" yaml:"action" xml:"action" db:"action"`
	// 记录事件后续处理所必须的数据
	Data string `form:"data" json:"data" yaml:"data" xml:"data" db:"data"`
	// 事件产生时的操作者,如果是系统动作产生的事件,此字段为空
	UserID int `form:"userID" json:"userID" yaml:"userID" xml:"userID" db:"userID"`
	// 是否删除 0正常 1已删除
	IsDeleted int `form:"isDeleted" json:"isDeleted" yaml:"isDeleted" xml:"isDeleted" db:"isDeleted"`
	//事件的删除时间截
	DeletedTime int `form:"deletedTime" json:"deletedTime" yaml:"deletedTime" xml:"deletedTime" db:"deletedTime"`
}

// 事件信息表结构
type EventUserSchema struct {
	// 关系ID值
	Id int `form:"id" json:"id" yaml:"id" xml:"id" db:"id"`
	// 事件ID值
	EventID int `form:"eventID" json:"eventID" yaml:"eventID" xml:"eventID" db:"eventID"`
	//用户ID值
	UserID int `form:"userID" json:"userID" yaml:"userID" xml:"userID" db:"userID"`
	//用户是否已经阅读了事件,0表示未读，1表示已读
	IsRead int `form:"isRead" json:"isRead" yaml:"isRead" xml:"isRead" db:"isRead"`
	//如果用户已阅读了事件，记录用户阅读事件的时间截
	ReadTime int `form:"readTime" json:"readTime" yaml:"readTime" xml:"readTime" db:"readTime"`
	// 是否删除 0正常 1已删除
	IsDeleted int `form:"isDeleted" json:"isDeleted" yaml:"isDeleted" xml:"isDeleted" db:"isDeleted"`
	//事件的删除时间截
	DeletedTime int `form:"deletedTime" json:"deletedTime" yaml:"deletedTime" xml:"deletedTime" db:"deletedTime"`
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

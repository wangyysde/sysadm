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
)

type Syssetting struct {
	// object name. the value should be "syssetting"
	Name string
	// table name which hold project data in DB
	TableName string
	// field name of primary key in the table
	PkName string
}

type SysSettingSchema struct {
	// 配置项ID,数据库中自增整数值
	Id uint `form:"id" json:"id" yaml:"id" xml:"id" db:"id"`
	// 配置项的应用范围，0 全局 1 k8s集群 2 节点级别 3 项目级别 4 用户组级别 5 用户级别
	// 对应的值参见syssettingDefined.go文件中的定义
	Scope int `form:"scope" json:"scope" yaml:"scope" xml:"scope" db:"scope"`
	// 配置项所适用的对象ID.如scope为0时,本字段值为0,如果scope字段是1时，则本字段为设置项所适用的k8s集群的集群ID值
	ObjectID string `form:"objectID" json:"objectID" yaml:"objectID" xml:"objectID" db:"objectID"`
	// 配置的key值, 不能为空。同一个级别的不能重复，大小写不敏感
	Key string `form:"key" json:"key" yaml:"key" xml:"key" db:"key"`
	// 配置项的默认值
	DefaultValue string `form:"defaultValue" json:"defaultValue" yaml:"defaultValue" xml:"defaultValue" db:"defaultValue"`
	// 配置项值
	Value string `form:"value" json:"value" yaml:"value" xml:"value" db:"value"`
	// 上次修改配置项值的用户,对应user表的userid.0表示系统自动设置的
	LastModifiedBy int `form:"lastModifiedBy" json:"lastModifiedBy" yaml:"lastModifiedBy" xml:"lastModifiedBy" db:"lastModifiedBy"`
	// 上次修改配置项值的时间截
	LastModifiedTime int `form:"lastModifiedTime" json:"lastModifiedTime" yaml:"lastModifiedTime" xml:"lastModifiedTime" db:"lastModifiedTime"`
	// 上次修改配置项值的原因
	LastModifiedReason string `form:"lastModifiedReason" json:"lastModifiedReason" yaml:"lastModifiedReason" xml:"lastModifiedReason" db:"lastModifiedReason"`
	// 上次修改前的值
	LastValue string `form:"lastValue" json:"lastValue" yaml:"lastValue" xml:"lastValue" db:"lastValue"`
}

// 存储运行期数据
type runingData struct {
	dbConf        *sysadmDB.DbConfig
	logEntity     *sysadmLog.LoggerConfig
	workingRoot   string
	sessionName   string
	sessionOption sessions.Options
	pageInfo      PageInfo
	objectEntiy   sysadmObjects.ObjectEntity
}
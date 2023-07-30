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

type Availablezone struct {
	// object name. the value should be "availablezone"
	Name string
	// table name which hold project data in DB
	TableName string
	// field name of primary key in the table
	PkName string
}

// 可用区信息表结构
type AvailablezoneSchema struct {
	// 集群ID，非自增，由雪花算法生成
	Id uint `form:"id" json:"id" yaml:"id" xml:"id" db:"id"`
	// 可用区中文名称
	CnName string `form:"cnName" json:"cnName" yaml:"cnName" xml:"cnName" db:"cnName"`
	// 可用区英文名称
	EnName string `form:"enName" json:"enName" yaml:"enName" xml:"enName" db:"enName"`
	// 可用区所属的数据中心ID
	Datacenterid uint `form:"datacenterid" json:"datacenterid" yaml:"datacenterid" xml:"datacenterid" db:"datacenterid"`
	// 值班电话
	DutyTel string `form:"dutyTel" json:"dutyTel" yaml:"dutyTel" xml:"dutyTel" db:"dutyTel"`
	// 状态0未启用 1已启用 2 已停用
	Status int `form:"status" json:"status" yaml:"status" xml:"status" db:"status"`
	// 是否删除 0正常 1已删除
	IsDeleted int `form:"isDeleted" json:"isDeleted" yaml:"isDeleted" xml:"isDeleted" db:"isDeleted"`

	// 创建时间
	CreateTime string `form:"createTime" json:"createTime" yaml:"createTime" xml:"createTime" db:"createTime"`
	// 更新时间
	UpdateTime string `form:"updateTime" json:"updateTime" yaml:"updateTime" xml:"updateTime" db:"updateTime"`
	// 创建人,对应user表的userid
	CreateBy uint `form:"createBy" json:"createBy" yaml:"createBy" xml:"createBy" db:"createBy"`
	// 更新人,对应user表的userid
	UpdateBy uint `form:"updateBy" json:"updateBy" yaml:"updateBy" xml:"updateBy" db:"updateBy"`
	// 备注
	Remark string `form:"remark" json:"remark" yaml:"remark" xml:"remark" db:"remark"`
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

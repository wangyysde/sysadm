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

type ObjectInfoSchema struct {
	// 对象ID
	ID int `form:"id" json:"id" yaml:"id" xml:"id" db:"id"`

	// 对象中文名称，平台整个环境下不可重复
	CnName string `form:"cnName" json:"cnName" yaml:"cnName" xml:"cnName" db:"cnName"`

	//对象英文名称，平台整个环境下不可重复
	EnName string `form:"enName" json:"enName" yaml:"enName" xml:"enName" db:"enName"`

	//对象所对应的在数据库中表名
	TableName string `form:"tableName" json:"tableName" yaml:"tableName" xml:"tableName" db:"tableName"`

	// 对象数据库表中标识对象的主键的字段名
	PkName string `form:"pkName" json:"pkName" yaml:"pkName" xml:"pkName" db:"pkName"`

	// 对象上是否可以运行命令.0表示否,1表示是
	CanRunCommand int `form:"canRunCommand" json:"canRunCommand" yaml:"canRunCommand" xml:"canRunCommand" db:"canRunCommand"`

	// 对象是否可以与命令关联, 即对象的状态是否可以影响命令的执行.0 表示否, 1表示是
	IsCommandRelated int `form:"isCommandRelated" json:"isCommandRelated" yaml:"isCommandRelated" xml:"isCommandRelated" db:"isCommandRelated"`

	// 是否废弃，0表示否，1表示是
	Deprecated int `form:"deprecated" json:"deprecated" yaml:"deprecated" xml:"deprecated" db:"deprecated"`
}

type ObjectTableSchema struct {
	// 对象ID
	ID int `form:"id" json:"id" yaml:"id" xml:"id" db:"id"`

	// 主对象的ID，即对应objectinfo表中id字段的值
	ObjectID int `form:"objectID" json:"objectID" yaml:"objectID" xml:"objectID" db:"objectID"`

	// 对象中文名称，平台整个环境下不可重复
	CnName string `form:"cnName" json:"cnName" yaml:"cnName" xml:"cnName" db:"cnName"`

	//对象英文名称，平台整个环境下不可重复
	EnName string `form:"enName" json:"enName" yaml:"enName" xml:"enName" db:"enName"`

	//对象所对应的在数据库中表名
	TableName string `form:"tableName" json:"tableName" yaml:"tableName" xml:"tableName" db:"tableName"`

	// 对象数据库表中标识对象的主键的字段名
	PkName string `form:"pkName" json:"pkName" yaml:"pkName" xml:"pkName" db:"pkName"`

	// 对象上是否可以运行命令.0表示否,1表示是
	CanRunCommand int `form:"canRunCommand" json:"canRunCommand" yaml:"canRunCommand" xml:"canRunCommand" db:"canRunCommand"`

	// 对象是否可以与命令关联, 即对象的状态是否可以影响命令的执行.0 表示否, 1表示是
	IsCommandRelated int `form:"isCommandRelated" json:"isCommandRelated" yaml:"isCommandRelated" xml:"isCommandRelated" db:"isCommandRelated"`

	// 是否废弃，0表示否，1表示是
	Deprecated int `form:"deprecated" json:"deprecated" yaml:"deprecated" xml:"deprecated" db:"deprecated"`
}

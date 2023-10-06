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

// ExecutionType defining the type for command executing
type ExecutionType int
type AutomationKind int
type CommandTransactionScope int
type CommandType int
type ParaKind int

const (
	defaultObjectName                   = "command"
	defaultTableName                    = "commandDefined"
	defaultPkName                       = "id"
	defaultCommandTableName             = "command"
	defaultCommandPkName                = "commandID"
	defaultCommandParasDefinedTableName = "commandParasDefined"
	defaultCommandParasDefinedPkName    = "id"
	defaultParasTableName               = "commandParameters"
	defaultParasPkName                  = "parametersID"
	DefaultModuleName                   = "command"
	DefaultApiVersion                   = "1.0"

	// 对应的命令自动执行
	ExecutionTypeAuto ExecutionType = 0
	// 手动执行的命令
	ExecutionTypeHand ExecutionType = 1

	// 在对象创建时自动执行
	AutomationKindObjectCreate AutomationKind = 0
	// 在对象的配置信息修改时自动执行
	AutomationKindObjectConfChange AutomationKind = 1
	// 在对象的状态发生改变时自动执行
	AutomationKindObjectStatusChange AutomationKind = 2
	// 在对象被删除时自动执行
	AutomationKindObjectDelete AutomationKind = 3
	// 定时执行一次或周期性执行,支持linux下crontab格式定义执行的时间和周期
	AutomationKindCrontab AutomationKind = 4

	// 命令执行的先后顺序及相关性只限制在本节点范围内，即无需判断其它节点上是否有依赖命令
	TransationScopeHost CommandTransactionScope = 0
	// 命令执行的先后顺序及相关性限制在同一个集群内
	TransationScopeCluster CommandTransactionScope = 1

	// 内嵌命令
	CommandTypeBuiltin CommandType = 0
	// 系统命令
	CommandTypeSys CommandType = 1
	// 脚本或者批处理程序
	CommandTypeScript CommandType = 2

	// 无参数
	ParaKindNo ParaKind = 0
	// 固定值
	ParaKindFixed ParaKind = 1
	// 对象字段值
	ParaKindObjFieldValue ParaKind = 2
	// 通过另一个命令获取
	ParaKindGetByCommand ParaKind = 3

	// 废弃了的命令定义
	CommandDefinedDeprecated int = 1
	// 未被废弃的命令定义
	CommandDefinedUnDeprecated int = 0
)

var runData = runingData{}

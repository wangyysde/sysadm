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

package v1beta1

// 命令定义表结构
type Command struct {
	// 定义的命令ID值
	ID uint `form:"id" json:"id" yaml:"id" xml:"id" db:"id"`
	// 需执行的命令，如果是内嵌命令与命令名称相同，如果是脚本或系统系统则是命令的绝对路径，不能为空，且在本表内不能有重复.
	// agent 依据此名称进行路由
	Command string `form:"command" json:"command" yaml:"command" xml:"command" db:"command"`
	// 命令的名称，用于在前端显示时使用的
	Name string `form:"name" json:"name" yaml:"name" xml:"name" db:"name"`
	// 执行类型，具体类型含义见command对象中定义
	ExecutionType int `form:"executionType" json:"executionType" yaml:"executionType" xml:"executionType" db:"executionType"`
	// 自动执行类别，具体类别含义见command对象中定义
	AutomationKind int `form:"automationKind" json:"automationKind" yaml:"automationKind" xml:"automationKind" db:"automationKind"`
	// 当命令自动执行的类别与对象有关时，本字段指定命令是针对哪个对象的
	ObjectName string `form:"objectName" json:"objectName" yaml:"objectName" xml:"objectName" db:"objectName"`
	// 参数类别，具体含义参见command对象的定义文件
	ParaKind int `form:"paraKind" json:"paraKind" yaml:"paraKind" xml:"paraKind" db:"paraKind"`
	// 命令数据来自的对象名称，即对于本字段所指定的满足条件的每一个对象都执行本指令
	DataFromObject string `form:"dataFromObject" json:"dataFromObject" yaml:"dataFromObject" xml:"dataFromObject" db:"dataFromObject"`
	// 如果命令是crontab类别的自动执行命令，本字段指定crontab格式的执行时间和周期
	Crontab string `form:"crontab" json:"crontab" yaml:"crontab" xml:"crontab" db:"crontab"`
	// 指示命令是否是属地同步命令。所谓同步命令是指，命令能够快速执行完成，即能在一个HTTP会话请求超时之前（一般超时时间为几秒内）执行完成的命令。0表示同步命令，否则为异步命令
	Synchronized int `form:"synchronized" json:"synchronized" yaml:"synchronized" xml:"synchronized" db:"synchronized"`
	// 适应的操作系统类型.0表示适应于所有操作系统
	OSID int `form:"osID" json:"osID" yaml:"osID" xml:"osID" db:"osID"`
	// 适应的操作系统版本.0表示适应于所有版本
	OsVersionID int `form:"osversionid" json:"osversionid" yaml:"osversionid" xml:"osversionid" db:"osversionid"`
	// 命令是否依赖于其它命令，0表示不依赖于其它命令的独立命令,否则本字段记录被依赖命令的id
	Dependent int `form:"dependent" json:"dependent" yaml:"dependent" xml:"dependent" db:"dependent"`
	// 0表示内嵌命令,1系统命令, 2 脚本或批处理程序
	Type int `form:"type" json:"type" yaml:"type" xml:"type" db:"type"`
	// 命令事务范围,当dependent值为0时忽略本字段值.含义见command对象定义文件
	TransactionScope int `form:"transactionScope" json:"transactionScope" yaml:"transactionScope" xml:"transactionScope" db:"transactionScope"`
	// 命令的执行是否必须带至少一个参数，0表示否，1表示1
	MustParas int `form:"mustParas" json:"mustParas" yaml:"mustParas" xml:"mustParas" db:"mustParas"`
	// 命令描述，一般包含功能和执行方法
	Descriptions string `form:"descriptions" json:"descriptions" yaml:"descriptions" xml:"descriptions" db:"descriptions"`
	// 是否废弃，0表示否，1表示是
	Deprecated int `form:"deprecated" json:"deprecated" yaml:"deprecated" xml:"deprecated" db:"deprecated"`
}

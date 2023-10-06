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

type Host struct {
	// object name. the value should be "host"
	Name string
	// table name which hold host data in DB
	TableName string
	// field name of primary key in the table
	PkName string
}

// 节点IP地址关系表结构
type HostIP struct {
	// Identifier of a IP come from DB
	IpID int `form:"ipID" json:"ipID" yaml:"ipID" xml:"id" db:"ipID"`
	// interface name of which the ip set
	DevName string `form:"devName" json:"devName" yaml:"devName" xml:"devName" db:"devName"`
	// IP address for IPV4
	Ipv4 string `form:"ipv4" json:"ipv4" yaml:"ipv4" xml:"ipv4" db:"ipv4"`
	// netmask address for ipv4 address
	Maskv4 string `form:"maskv4" json:"maskv4" yaml:"maskv4" xml:"maskv4" db:"maskv4"`
	// IP address for IPV6
	Ipv6 string `form:"ipv6" json:"ipv6" yaml:"ipv6" xml:"ipv6" db:"ipv6"`
	// netmask address for ipv6 address
	Maskv6 string `form:"maskv6" json:"maskv6" yaml:"maskv6" xml:"maskv6" db:"maskv6"`
	// 接口的Mac地址
	Mac string `form:"mac" json:"mac" yaml:"mac" xml:"mac" db:"mac"`
	// 0 for offline, 1 for online
	Status int `form:"status" json:"status" yaml:"status" xml:"status" db:"status"`
	// 0 not management ip, 1 management IP
	isManage int `form:"isManage" json:"isManage" yaml:"isManage" xml:"isManage" db:"isManage"`
	// 所属的节点ID
	HostId int `form:"hostid" json:"hostid" yaml:"hostid" xml:"id" db:"hostid"`
}

// 节点Yum关系表结构
type HostYum struct {
	// relationID identified relation of host and yum
	RelationID int `form:"relationID" json:"relationID" yaml:"relationID" xml:"id" db:"relationID"`
	// 主机的ID值，自动增长型
	// TODO 下一版本尝试将本字段的值修改为通过雪花算法生成
	HostId int `form:"hostid" json:"hostid" yaml:"hostid" xml:"id" db:"hostid"`
	// yumid identified a yum configuration
	YumID int `form:"yumid" json:"yumid" yaml:"yumid" xml:"id" db:"yumid"`
}

// 节点信息表结构
type HostSchema struct {
	// 主机的ID值，自动增长型
	// TODO 下一版本尝试将本字段的值修改为通过雪花算法生成
	HostId int `form:"hostid" json:"hostid" yaml:"hostid" xml:"id" db:"hostid"`
	// 用户ID, 自动增长型
	UserId string `form:"userid" json:"userid" yaml:"userid" xml:"userid" db:"userid"`
	// 项目ID, 当前取消主机关联项目，增加这个字段是为了与之前的代码兼容
	ProjectID int `form:"projectid" json:"projectid" yaml:"projectid" xml:"projectid"`
	// 主机名
	Hostname string `form:"hostname" json:"hostname" yaml:"hostname" xml:"hostname" db:"hostname"`
	// 操作系统发行版ID
	OSID int `form:"osID" json:"osID" yaml:"osID" xml:"osID" db:"osID"`
	// 操作系统版本ID
	OSVersionID int `form:"osversionid" json:"osversionid" yaml:"osversionid" xml:"osversionid" db:"osversionid"`
	// 主机运行状态
	Status string `form:"status" json:"status" yaml:"status" xml:"status" db:"status"`
	// 当apiserver是以主动模式运行时，连接客户端agent的地址或域名
	AgentIP string `form:"ip" json:"ip" yaml:"ip" xml:"ip" db:"ip"`
	// 主机的IP地址类型4表示ipv4 6表示IPV6
	IpType int `form:"iptype" json:"iptype" yaml:"iptype" xml:"iptype" db:"iptype"`
	// Agent的运行模式，0表示主动模式 1表示被动模式。注意，这里的主被动是从apiserver角度来看的。即主动模式表示由apiserver主动连接agent
	PassiveMode int `form:"passiveMode" json:"passiveMode" yaml:"passiveMode" xml:"passiveMode" db:"passiveMode"`
	// 主动模式时，apiserver向agent发送命令时，请求发送的目标标路径。如果本字段为空，则apiserver默认向/receiveCommand请求命令的状态
	CommandUri string `form:"commandUri" json:"commandUri" yaml:"commandUri" xml:"commandUri" db:"commandUri"`
	// 主动模式时，apiserver向agent查询命令执行状态时，请求发送的目标路径。如果本字段为空，则apiserver默认向/getCommandStatus请求命令的状态
	CommandStatusUri string `form:"commandStatusUri" json:"commandStatusUri" yaml:"commandStatusUri" xml:"commandStatusUri" db:"commandStatusUri"`
	// 主动模式时,apiserver向agent查询命令的执行日志时，请求的发送目标路径。如果本字段为空，则apiserver默认向/getLogs请求命令的日志
	CommandLogsUri string `form:"commandLogsUri" json:"commandLogsUri" yaml:"commandLogsUri" xml:"commandLogsUri" db:"commandLogsUri"`
	// 主动模式时,apiserver连接agent是否使用TLS.0表示否，否则表示是
	AgentIsTls int `form:"agentIsTls" json:"agentIsTls" yaml:"agentIsTls" xml:"agentIsTls" db:"agentIsTls"`
	//主动模式时,apiserver连接agent时主动采用TLS,本子段是CA证书内容
	AgentCa string `form:"agentCa" json:"agentCa" yaml:"agentCa" xml:"agentCa" db:"agentCa"`
	// 主动模式时,apiserver连接agent时主动采用TLS,本子段是证书内容
	AgentCert string `form:"agentCert" json:"agentCert" yaml:"agentCert" xml:"agentCert" db:"agentCert"`
	// 主动模式时,apiserver连接agent时主动采用TLS,本子段是密钥内容
	AgentKey string `form:"agentKey" json:"agentKey" yaml:"agentKey" xml:"agentKey" db:"agentKey"`
	// 主动模式时,apiserver连接agent时主动采用TLS,指定是否跳过检查证书合法性检查，1表示是，否则为否
	InsecureSkipVerify int `form:"insecureSkipVerify" json:"insecureSkipVerify" yaml:"insecureSkipVerify" xml:"insecureSkipVerify" db:"insecureSkipVerify"`
	// 主动模式时,apiserver连接agent时的端口号
	AgentPort int `form:"agentPort" json:"agentPort" yaml:"agentPort" xml:"agentPort" db:"agentPort"`
	// 节点所属集群ID，如果为空表示节点不隶属于任何K8S集群
	K8sClusterID string `form:"k8sclusterid" json:"k8sclusterid" yaml:"k8sclusterid" xml:"k8sclusterid" db:"k8sclusterid"`
	// 节点创建时间截
	CreateTime string `form:"createTime" json:"createTime" yaml:"createTime" xml:"createTime" db:"createTime"`
	// 节点下线开始时间
	OfflineStartTime string `form:"offlineStartTime" json:"offlineStartTime" yaml:"offlineStartTime" xml:"offlineStartTime" db:"offlineStartTime"`
	// 节点删除时间
	DeleteTime string `form:"deletetime" json:"deletetime" yaml:"deletetime" xml:"deletetime" db:"deletetime"`
	// 节点所属的数据中心ID
	Dcid uint `form:"dcid" json:"dcid" yaml:"dcid" xml:"dcid" db:"dcid"`
	// 节点所属的可用区ID
	Azid uint `form:"azid" json:"azid" yaml:"azid" xml:"azid" db:"azid"`
	// 节点的machineID
	MachineID string `form:"machineID" json:"machineID" yaml:"machineID" xml:"machineID" db:"machineID"`
	// 节点systemID
	SystemID string `form:"systemID" json:"systemID" yaml:"systemID" xml:"systemID" db:"systemID"`
	// 节点架构
	Architecture string `form:"architecture" json:"architecture" yaml:"architecture" xml:"architecture" db:"architecture"`
	// 操作系统内核版本
	KernelVersion string `form:"kernelVersion" json:"kernelVersion" yaml:"kernelVersion" xml:"kernelVersion" db:"kernelVersion"`
	// 节点上IP地址信息列表
	HostIps []HostIP
	// 节点上配置的Yum信息列表
	HostYum []HostYum
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

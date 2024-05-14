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

import "sysadm/objectsUI"

var hostObjectName = "host"
var hostTableName = "host"
var hostTablePkName = "hostid"
var hostIPTableName = "hostIP"
var hostIPPkName = "ipID"
var hostYumTableName = "hostYum"
var hostYumPkName = "relationID"
var hostModuleName = "节点"

type HostStatusKind string

const (
	HostStatusRunning        HostStatusKind = "running"
	HostStatusMaintenance    HostStatusKind = "maintenance"
	HostStatusOffline        HostStatusKind = "offline"
	HostStatusDeleted        HostStatusKind = "deleted"
	HostStatusUnkown         HostStatusKind = "unkown"
	HostStatusReady          HostStatusKind = "Ready"
	HostStatusMemoryPressure HostStatusKind = "MemoryPressure"
	HostStatusDiskPressure   HostStatusKind = "DiskPressure"
	HostPIDPressure          HostStatusKind = "PIDPressure"
	HostNetworkUnavailable   HostStatusKind = "NetworkUnavailable"
	HostTypeIPTypeV4         int            = 4
	HostTypeIPTypeV6         int            = 6
)

var hostAllListItems = map[string]string{"TD1": "主机名", "TD2": "数据中心", "TD3": "可用区", "TD4": "操作系统/版本", "TD5": "所属K8S集群", "TD6": "状态"}
var hostAllPopMenuItems = []string{"查看详情,detail,GET,page", "编辑主机信息,edit,GET,page", "驱逐容器,drain,POST,tip", "停止调度,cordon,POST,tip",
	"恢复调度,uncordon,POST,tip", "移出集群,moveformcluster,POST,tip", "删除,del,POST,tip", "加入集群,joinCluster,page"}
var hostAllOrderFields = map[string]objectsUI.SortBy{"TD1": sortHostByName, "TD5": sortHostByClusterID}
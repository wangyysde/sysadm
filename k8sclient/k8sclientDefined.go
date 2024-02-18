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

package k8sclient

const (
	ObjectStatusPending    = "Pending"
	ObjectStatusRunning    = "Running"
	ObjectStatusSucceeded  = "Succeeded"
	ObjectStatusFailed     = "Failed"
	ObjectStatusUnknow     = "Unknow"
	FieldManager           = "k8sclient.sysadm.cn"
	defaultClusterUserName = "kubernetes-admin"
	// 当连接的上下文没有指定默认命名空间时，使用默认命名空间名
	defaultNamespace = "default"
)

var K8sClusterConnectType = map[string]string{"0": "证书方式连接", "1": "Token方式连接", "2": "KubeConfig"}

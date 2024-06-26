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

var DefaultObjectName = "k8scluster"
var DefaultTableName = "k8scluster"
var DefaultPkName = "id"
var DefaultModuleName = "k8scluster"
var DefaultApiVersion = "1.0"
var runData = runingData{}
var allStatus = map[int]string{0: "未启用", 1: "已启用", 2: "已停用"}
var k8sClusterConnectType = map[string]string{"0": "证书方式连接", "1": "Token方式连接", "2": "KubeConfig"}

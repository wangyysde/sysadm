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
*
* NOTE:
* 本文件定义了将来准备独立出来的sysadmApiServer组件所使用的导出的数据结构、常量和全局变量
*/

package app

// specifies a identifer of the node which agent running on it.
// It is any combination of the IP,HOSTNAME and MAC joined by commas  or a customize string what the leght of the string is less 63
// agent will get all IPs without not active and reponse these IPs in list to the server by nodeIdentifer.IPs filed if IP is included in NodeIdentifer
// agent will get hostname and reponse the hostname  to the server by nodeIdentifer.Hostname filed if hostname is included in NodeIdentifer
// agent will get all MACs without not active and reponse these MACs in list to the server by nodeIdentifer.MACs filed if MAC is included in NodeIdentifer
// customize string is reponse to the server directly .
// customize string is conflicted with IP,HOSTNAME and MAC. the nodeIdentifer can be changed by the server during agent communicate with the server
var DefaultNodeIdentifer string = "IP,HOSTNAME,MAC"

// max lenght of customize node identifier
var MaxCustomizeNodeIdentiferLen int = 64


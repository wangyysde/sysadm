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

var DefaultObjectName = "datacenter"
var DefaultTableName = "datacenter"
var DefaultPkName = "id"
var DefaultModuleName = "datacenter"
var DefaultApiVersion = "1.0"
var runData = runingData{}

const (
	StatusUnused   int = 0
	StatusEnabled  int = 1
	StatusDisabled int = 2

	LineTypeCT       int = 1
	LineTypeCUCC     int = 2
	LinetypeCMCC     int = 3
	LineTypeCBN      int = 4
	LineTypeBGP2     int = 5
	LineTypeBGP3     int = 6
	LineTypeBGP4     int = 7
	LineTypeOverseas int = 8
)

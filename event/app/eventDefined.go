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

var defaultObjectName = "event"
var defaultTableName = "event"
var defaultPkName = "id"
var defaultModuleName = "event"
var defaultApiVersion = "1.0"
var runData = runingData{}

const (
	ClassInfo    int = 0
	ClassWarning int = 1
	ClassError   int = 2
	ClassFatal   int = 3
)

const (
	ScopeHareware = iota
	ScopeCluster
	ScopeApp
	ScopeService
)

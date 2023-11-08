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
	"strings"
	sysadmObjects "sysadm/objects/app"
	sysadmUtils "sysadm/utils"
)

func validName(name string) bool {
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}

	// get total number of list objects
	var commandEntity sysadmObjects.ObjectEntity
	commandEntity, e := New(runData.dbConf, runData.workingRoot)
	if e != nil {
		return false
	}

	if len(name) > 255 {
		return false
	}

	commandConditions := make(map[string]string, 0)
	commandConditions["deprecated"] = "='0'"
	commandConditions["name"] = "='" + name + "'"
	var emptyString []string
	commandCount, e := commandEntity.GetObjectCount("", emptyString, emptyString, commandConditions)
	if e != nil || commandCount > 0 {
		return false
	}

	return true
}

func validCommand(command string) bool {
	command = strings.TrimSpace(command)
	if command == "" || len(command) > 64 {
		return false
	}

	if !sysadmUtils.ValidPath(command) {
		return false
	}

	// get total number of list objects
	var commandEntity sysadmObjects.ObjectEntity
	commandEntity, e := New(runData.dbConf, runData.workingRoot)
	if e != nil {
		return false
	}
	commandConditions := make(map[string]string, 0)
	commandConditions["deprecated"] = "='0'"
	commandConditions["name"] = "='" + command + "'"
	var emptyString []string
	commandCount, e := commandEntity.GetObjectCount("", emptyString, emptyString, commandConditions)
	if e != nil || commandCount > 0 {
		return false
	}

	return true
}

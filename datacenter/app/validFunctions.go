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
)

func validCnName(name string) bool {
	name = strings.TrimSpace(name)
	if name == "" {
		return false
	}

	// get total number of list objects
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = New()

	if len(name) > 255 {
		return false
	}

	dcConditions := make(map[string]string, 0)
	dcConditions["isDeleted"] = "='0'"
	dcConditions["cnName"] = "='" + name + "'"
	var emptyString []string
	commandCount, e := dcEntity.GetObjectCount("", emptyString, emptyString, dcConditions)
	if e != nil || commandCount > 0 {
		return false
	}

	return true
}

func validEnName(name string) bool {
	name = strings.TrimSpace(name)
	if name == "" {
		return true
	}

	// get total number of list objects
	var dcEntity sysadmObjects.ObjectEntity
	dcEntity = New()

	if len(name) > 255 {
		return false
	}

	dcConditions := make(map[string]string, 0)
	dcConditions["isDeleted"] = "='0'"
	dcConditions["enName"] = "='" + name + "'"
	var emptyString []string
	commandCount, e := dcEntity.GetObjectCount("", emptyString, emptyString, dcConditions)
	if e != nil || commandCount > 0 {
		return false
	}

	return true
}

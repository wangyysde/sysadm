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
)

func validTextValue(v string, minLen, maxLen int, allowEmpty bool) bool {
	v = strings.TrimSpace(v)
	if !allowEmpty && v == "" {
		return false
	}

	if len(v) < minLen || len(v) > maxLen {
		return false
	}

	return true
}

func validCNName(cnName string) bool {
	return validTextValue(cnName, 0, 255, false)
}

func validENName(enName string) bool {
	return validTextValue(enName, 0, 255, true)
}

func validApiserverAddress(address string) bool {
	if !validTextValue(address, 0, 255, false) {
		return false
	}

	address = strings.TrimSpace(address)
	apiserverArray := strings.Split(address, ":")
	if len(apiserverArray) != 2 {
		return false
	}

	return true
}

func validClusterUser(user string) bool {
	return validTextValue(user, 0, 255, false)
}

func validDutyTel(dutyTel string) bool {
	return validTextValue(dutyTel, 0, 20, true)
}

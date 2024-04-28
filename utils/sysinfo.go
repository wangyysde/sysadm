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

package utils

import (
	"os"
	"path"
	"strings"
)

func GetSystemUUID() (string, error) {
	if id, err := os.ReadFile(path.Join(dmiDir, "id", "product_uuid")); err == nil {
		return strings.TrimSpace(string(id)), nil
	} else if id, err = os.ReadFile(path.Join(ppcDevTree, "system-id")); err == nil {
		return strings.TrimSpace(string(id)), nil
	} else if id, err = os.ReadFile(path.Join(ppcDevTree, "vm,uuid")); err == nil {
		return strings.TrimSpace(string(id)), nil
	} else if id, err = os.ReadFile(path.Join(s390xDevTree, "machine-id")); err == nil {
		return strings.TrimSpace(string(id)), nil
	} else {
		return "", err
	}
}

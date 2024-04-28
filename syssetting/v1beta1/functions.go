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

package v1beta1

import (
	"fmt"
)

func GetCertAndKeyName(certType int) (string, string, error) {
	certKey := ""
	keyKey := ""

	switch certType {
	case CertTypeCa:
		certKey = SettingKeyForCA
		keyKey = SettingKeyForCaKey
	case CertTypeApiServer:
		certKey = SettingKeyForApiServerCert
		keyKey = SettingKeyForApiServerCertKey
	case CertTypeAgent:
		certKey = SettingKeyForAgentCert
		keyKey = SettingKeyFroAgentCertKey
	default:
		return "", "", fmt.Errorf("type of certificate was not found")
	}

	return certKey, keyKey, nil
}

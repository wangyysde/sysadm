/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
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

// settings for sysadm system.
// settings in this file is temporary.they will be moved into Database
// TODO
// move the settings in this file into Database
package app

import (
	"github.com/wangyysde/sshclient/sshkey"
)

// settings for infrastructure module
var infrastructure map[string]interface{} = map[string]interface{} {
	"publicKeyFile": "sshkeys/sshkey.pub",
	"publicKeyFingerprint": "",
	"keyType": sshkey.KEY_RSA,
	"keySize": 2048,
	"privateKeyFile": "sshkeys/id_rsa",
	"privateKeyFingerprint": "",
	"scriptsPath": "scripts",
	"scriptGetHostIPs": "gethHostIPs.sh",
}


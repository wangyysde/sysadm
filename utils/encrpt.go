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

Errorcode: 120xxx
 */

package utils

import(
	"crypto/md5"
	"encoding/hex"

)

/* 
   generating md5 data using data and salt. salt can empty.
   return "" if data is empty. 
   otherwise return string
*/
func Md5Encrypt(data string, salt string) string{
	if len(data) < 1 {
		return ""
	}
	md5Ctx := md5.New()
	md5Ctx.Write([]byte(data))
	if len(salt) > 0 {
		md5Ctx.Write([]byte(salt))
	}
	cipherStr := md5Ctx.Sum(nil)
	encryptedData := hex.EncodeToString(cipherStr)
	return encryptedData
}


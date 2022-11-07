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

package app

import (
	"fmt"
	"strings"

	"github.com/wangyysde/sysadm/utils"
)

/*
	gets the value of file path by module and setting name
	return the value just got and nil  or return "" and error
*/
func GetFilePath(moduleName string, settingName string ) (string, error){
	if strings.TrimSpace(moduleName) == "" || strings.TrimSpace(settingName) == "" {
		return "", fmt.Errorf("can not get setting value with empty module name(%s) or setting name(%s)",moduleName,settingName)
	}

	moduleName = strings.ToLower(strings.TrimSpace(moduleName))
	if moduleName == "infrastructure" {
		if v,ok := infrastructure[settingName]; ok {
			return utils.Interface2String(v),nil
		} 
		return "", fmt.Errorf("no found setting (%s) in module (%s) ",settingName, moduleName)
	}

	return "", fmt.Errorf("no found module %s", moduleName)  
}
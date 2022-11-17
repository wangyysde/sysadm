/* =============================================================
* @Author:  Wayne Wang <net_use@bzhy.com>
*
* @Copyright (c) 2022 Bzhy Network. All rights reserved.
* @HomePage http://www.sysadm.cn
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at:
* https://www.sysadm.cn/licenses/apache-2.0.txt
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and  limitations under the License.
*
 */

package app

import (
	"github.com/wangyysde/sysadm/config"
)

func SetVersion(version *config.Version){
	if version == nil {
		return
	}

	version.Version = ver
	version.Author = author

	CliOps.Version = *version 
}

func GetVersion() *config.Version {
	if CliOps.Version.Version != "" {
		return &CliOps.Version
	}

	return nil
}